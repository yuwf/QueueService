package utils

import (
	"container/list"
	"github.com/panjf2000/ants"
	"sync"
)

// 消息处理使用协成池 解决大规模消息 per msg per gorutinue的问题
// 协成池任务队列 保证任务顺序执行
// 可以解决消息的顺序和数据安全问题
type Sequence struct {
	mutex sync.Mutex
	tasks list.List
}

func (self *Sequence) Submit(task func()) {
	if task == nil {
		return
	}
	self.mutex.Lock()         // 加锁
	defer self.mutex.Unlock() // 退出时解锁

	// 添加任务
	self.tasks.PushBack(task)

	// 当前只有一个刚加入的任务，开启协成池调用handle
	if self.tasks.Len() == 1 {
		ants.Submit(self.handle)
	}
}

func (self *Sequence) handle() {
	//取出一个任务
	self.mutex.Lock() // 加锁
	if self.tasks.Len() == 0 {
		self.mutex.Unlock() // 解锁
		return              // 退出
	}
	task := self.tasks.Front().Value.(func())
	self.mutex.Unlock() // 解锁

	// 执行task
	if task != nil {
		task()
	}

	self.mutex.Lock() // 加锁
	// 移除当前完成的任务
	self.tasks.Remove(self.tasks.Front())
	// 如果任务列表不为空继续开启下一个handle
	if self.tasks.Len() > 0 {
		ants.Submit(self.handle)
	}
	self.mutex.Unlock() // 解锁
}
