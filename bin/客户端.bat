echo 开始启动客户端

@ping -n 5 127.0.0.1 > nul 2>nul
@echo 启动第【1】个客户端
@start go_build_QueueService_client.exe

@ping -n 5 127.0.0.1 > nul 2>nul
@echo 启动第【2】个客户端
@start go_build_QueueService_client.exe

@ping -n 5 127.0.0.1 > nul 2>nul
@echo 启动第【3】个客户端
@start go_build_QueueService_client.exe

@ping -n 5 127.0.0.1 > nul 2>nul
@echo 启动第【4】个客户端
@start go_build_QueueService_client.exe

@ping -n 5 127.0.0.1 > nul 2>nul
@echo 启动第【5】个客户端
@start go_build_QueueService_client.exe

@ping -n 5 127.0.0.1 > nul 2>nul
@echo 启动第【6】个客户端
@start go_build_QueueService_client.exe

@ping -n 5 127.0.0.1 > nul 2>nul
@echo 启动第【7】个客户端
@start go_build_QueueService_client.exe

@ping -n 5 127.0.0.1 > nul 2>nul
@echo 启动第【8】个客户端
@start go_build_QueueService_client.exe

@ping -n 5 127.0.0.1 > nul 2>nul
@echo 启动第【9】个客户端
@start go_build_QueueService_client.exe

@ping -n 5 127.0.0.1 > nul 2>nul
@echo 启动第【10】个客户端
@start go_build_QueueService_client.exe

@ping -n 5 127.0.0.1 > nul 2>nul
@echo 启动第【11】个客户端
@start go_build_QueueService_client.exe
pause
