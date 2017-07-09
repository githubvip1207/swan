
# 小天鹅热编译

swan是一个自动检测工程中代码变化，触发编译、起停等动作的工具

## 使用方法：
 - Usage: swan -p <path>
 - Param:
	- -p: 工程根目录，默认为运行目录
 - Example:
	- swan -p /home/chenyu/swan

## 配置文件
swan会在工程目录，用户家目录寻找.swanconfig配置文件
```
[basic]
suffixes = go, conf # 检测的文件后缀

[command]
build = ./build.sh # 工程编译命令
stop  = ./control.sh stop # 工程停止命令
start = ./control.sh start # 工程启动命令
```
