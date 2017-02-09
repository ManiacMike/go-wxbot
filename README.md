## gorobot说明

本项目是从python的wxbot项目"翻译"过来的，加了一些俏皮的群聊功能。

## 使用

根目录下创建一个config.ini，就可以跑了，mac/linux是肯定没问题的

config.ini的内容如下，需要说明的是，你得先去http://www.tuling123.com 注册个机器人。
“熊本君”是注册是机器人名字，key从网站获取
~~~
[server]
port = 8004

[turing]
base_url = http://www.tuling123.com/openapi/api

[熊本君]
key = xxxxxxxxxxx
~~~

## 其它

代码基本以实现功能为主，有点乱，欢迎各位fork重构，想学习go的推荐下我这个项目
https://github.com/ManiacMike/gwork
