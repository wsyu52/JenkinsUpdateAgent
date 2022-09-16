# JenkinsUpdateAgent
Jenkins插件源镜像小工具
### 遇到的问题
1. Jenkins官方的插件源有速度慢或被墙的问题，导致在下载插件时出现速度缓慢或无法链接的情况。 
虽然有国内的几个镜像， 但是更新信息json文件里面的下载地址还是指向官方源的下载地址，简单更改json文件路径并不能起到真正地加速作用
2. 将插件信息json文件里面的下载路径直接替换成国内镜像地址，可以暂时解决，但每次都要手动一遍实属不够经济
3. 插件信息json文件版本问题：jenkins版本跟插件版本有一些兼容要求，需要根据不同版本找对应版本的更新信息，官方更新服务器有自动跳转对应版本的功能，国内镜像是纯粹的文件服务，并没有根据版本跳转的功能，换一次jenkins版本就需要手动指定，且并不是每个jenkins版本都有对应版本的json文件，需要做一个匹配
### 如何解决
恰巧最近刚刚开始接触Go语言，简洁效率高，干脆写一个小代理工具一次性解决上面的问题，将jenkins插件更新地址指向这个小工具，让小工具完成这些版本适配，地址替换等工作就好了， 于是就有了这么个项目

## 如何使用
1. 启动此工具：直接部署 or docker容器部署 都可
2. Jenkins启动参数添加` -Dhudson.model.DownloadService.noSignatureCheck=true`
3. 将Jenkins插件更新地址改为 `http://<此工具地址>:8888/update-center.json`
4. 好了，享受飞一般地下载速度

## 如何部署本工具
### 使用已做好的docker镜像
```shell
docker run -d -p 8888:8888 --name jenkins-update-agent wsyu52/jenkins-update-agent
```
### 自己构建docker镜像
```shell
docker build -t wsyu52/jenkins-update-agent .
docker run -d -p 8888:8888 --name jenkins-update-agent wsyu52/jenkins-update-agent
```
### 搭配jenkins一起用docker compose
```shell
# 注意修改jenkins的jenkins_home路径
docker compose up
```
### 纯自编译
```shell
go env -w GOPROXY=https://goproxy.cn
go mod tidy
go build src/JenkinsUpdateAgent.go
# start service
./JenkinsUpdateAgent
```


## 关于本工具
* 本工具镜像容量和运行时占用资源都非常小，可以随jenkins一直运行，无需担心资源占用
* 目前只是初步做了实现，未来会进一步完善
* 如果有任何建议或问题，欢迎提issue或PR
