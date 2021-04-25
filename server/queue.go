package main

import (
	"QueueService/utils"
	"container/list"
	"github.com/rs/zerolog/log"
	"sync"
	"sync/atomic"
	"time"
)

var theQueue = &Queue{clientPos: make(map[string]*QueueItem), exitChan: make(chan bool), changeChann: make(chan *ChangeInfo, 2048)}

// 队列类 用户排队管理
type Queue struct {
	// 读写锁
	mutex sync.RWMutex
	// 排队用户位置 [uid:*QueueItem]
	clientPos map[string]*QueueItem
	// 排队的用户列表 *QueueItem
	clientList list.List

	// 待添加或者删除对象 协成安全
	// 数据每隔一段时间更新到排队位置中
	changeChann    chan *ChangeInfo
	changeList     list.List
	changeAddCount int32 // changeList中添加的数量 原子操作 为了能及时反应出玩家排队位置

	// 统计使用 原子调用
	LastAddNum    int32
	LastRemoveNum int32
	LastLoginNum  int32

	exitChan chan bool
}

type ChangeInfo struct {
	op  int // 0添加 1删除
	uid string
	seq *utils.Sequence
}

type QueueItem struct {
	uid string
	seq *utils.Sequence
	pos int
	e   *list.Element // 在列表中的位置 用于快速删除
}

func (self *Queue) Init() bool {
	log.Info().Msg("Queue Init")

	go self.Run_Change()
	go self.Run_Login()
	return true
}

func (self *Queue) Run_Change() {
	for {
		select {
		case <-self.exitChan:
			return
		case change := <-self.changeChann:
			self.changeList.PushBack(change)
		case <-time.After(time.Millisecond * 10):
			add, remove := self.UpdateChange()
			if add > 0 {
				atomic.AddInt32(&self.LastAddNum, int32(add))
			}
			if remove > 0 {
				atomic.AddInt32(&self.LastRemoveNum, int32(remove))
			}
		}
	}
}

func (self *Queue) Run_Login() {
	for {
		select {
		case <-self.exitChan:
			return
		case <-time.After(time.Second):
			login := self.CheckLogin()
			atomic.AddInt32(&self.LastLoginNum, int32(login))
		}
	}
}

func (self *Queue) UpdateChange() (add int, remove int) {
	if self.changeList.Len() == 0 {
		return
	}

	// 写锁加锁
	self.mutex.Lock()

	for change := self.changeList.Front(); change != nil; change = change.Next() {
		info, _ := change.Value.(*ChangeInfo)
		if info.op == 0 {
			// 添加操作 判断是否已经在队列里面了
			_, ok := self.clientPos[info.uid]
			if ok {
				continue
			}
			item := &QueueItem{info.uid, info.seq, self.clientList.Len(), nil}
			// 添加到队列里面
			item.e = self.clientList.PushBack(item)
			self.clientPos[info.uid] = item
			add++
		} else if info.op == 1 {
			// 添加操作 判断是否已经在队列里面了
			item, ok := self.clientPos[info.uid]
			if !ok {
				continue
			}
			delete(self.clientPos, info.uid)
			self.clientList.Remove(item.e)
			remove++
		}
	}

	// 根据clientList重新生成排队
	index := 1
	for client := self.clientList.Front(); client != nil; client = client.Next() {
		item, _ := client.Value.(*QueueItem)
		item.pos = index
		index++
	}

	// 解锁
	self.mutex.Unlock()

	// 清空
	self.changeList = list.List{}
	atomic.StoreInt32(&self.changeAddCount, 0)
	return
}

// 判断排队的用户是否可以登录了 返回登录的数量
func (self *Queue) CheckLogin() (num int) {
	// 获取在线人数和最近一秒登录的人数
	onlinenum := int(atomic.LoadInt32(&theServer.OnlineNum))
	lastloginnum := int(atomic.LoadInt32(&theServer.LastLoginNum))
	if onlinenum >= theConf.MaxOnlineNum || lastloginnum >= theConf.LoginNumPreSec {
		return
	}
	// 计算可以登录的人数
	loginnum := theConf.LoginNumPreSec - lastloginnum
	if loginnum > theConf.MaxOnlineNum-onlinenum {
		loginnum = theConf.MaxOnlineNum - onlinenum
	}
	if loginnum <= 0 {
		return
	}

	// 可以登录 写锁加锁
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if self.clientList.Len() == 0 {
		return
	}

	// 给玩家发成功登录逻辑
	var next *list.Element
	for client := self.clientList.Front(); client != nil; client = next {
		next = client.Next()
		item, _ := client.Value.(*QueueItem)
		item.seq.Submit(func() {
			c, ok := theServer.ClientMgr.Load(item.uid)
			if ok {
				c.(*Client).FinishQueue()
			} else {
			}
		})
		delete(self.clientPos, item.uid)
		self.clientList.Remove(client)
		num++
		if num >= loginnum {
			break
		}
	}

	// 根据clientList重新生成排队
	index := 1
	for client := self.clientList.Front(); client != nil; client = client.Next() {
		item, _ := client.Value.(*QueueItem)
		item.pos = index
		index++
	}
	return
}

// 添加用户
func (self *Queue) AddClient(uid string, s *utils.Sequence) {
	change := &ChangeInfo{0, uid, s}
	self.changeChann <- change
	atomic.AddInt32(&self.changeAddCount, 1)
}

// 移除用户
func (self *Queue) RemoveClient(uid string) {
	change := &ChangeInfo{1, uid, nil}
	self.changeChann <- change
}

// 获取用户排队信息
// 若用户没有排队或者还没放到排队列表中返回的排队位置和排队人数相等
// 用户有没有排队需要根据client中的LoginState来判断
func (self *Queue) GetClientPos(uid string) (int, int) {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	item, ok := self.clientPos[uid]
	waitAddCount := int(atomic.LoadInt32(&self.changeAddCount))
	if !ok {
		return self.clientList.Len() + waitAddCount, self.clientList.Len() + waitAddCount
	}
	return item.pos, self.clientList.Len() + waitAddCount
}

// 获取排队人数
func (self *Queue) GetClientNum() int {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	return self.clientList.Len() + int(atomic.LoadInt32(&self.changeAddCount))
}

func (self *Queue) Quit() {
	log.Info().Msg("Queue Quit")

	//广播 关闭
	close(self.exitChan)
	close(self.changeChann)
}
