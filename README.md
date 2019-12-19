## gorobot说明

本项目是从python的wxbot项目"翻译"过来的，加了一些俏皮的群聊功能。

## 使用

根目录下创建一个config.ini，就可以跑了，mac/linux/windows都支持

config.ini的内容如下，需要说明的是，你得先去http://www.tuling123.com 注册个机器人。
“熊本君”是你的机器人名字，key从网站获取
~~~
[server]
port = 8004

[turing]
base_url = http://openapi.tuling123.com/openapi/api/v2

[熊本君]
key = xxxxxxxxxxx
~~~

项目目录下运行go build

然后运行二进制文件就可以了
 ./go-wxbot.exe(win)
 ./go-wxbot(unix)
