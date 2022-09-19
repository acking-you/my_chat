
- [MyChat](#mychat)
  - [项目运行](#项目运行)
  - [项目愿景](#项目愿景)
  - [目前进度](#目前进度)
  - [后端架构](#后端架构)
    - [项目目录说明](#项目目录说明)
    - [整体业务流程](#整体业务流程)
    - [模块划分](#模块划分)
    - [类图](#类图)
  - [前端效果展示](#前端效果展示)
# MyChat
>注意：本项目使用go1.18进行开发，且使用到了go workspace进行内部依赖管理，以及go泛型简化业务，故需要的golang版本至少达到1.18
## 项目运行
>温馨提示：运行项目前请修改配置文件
```shell
go run main/main.go
```
## 项目愿景
* 实现实时聊天系统（对应使用websocket协议）
* 实现用户信息和类似于QQ空间的功能（对应http协议）
* 实现音视频通话（对应webrtc协议）
* 实现远程控制（对应webrtc协议）

## 目前进度
- [x] 登录注册
- [x] 一对一实时聊天
- [ ] 聊天群组
- [x] 用户信息拉取
- [ ] 用户信息自定义（头像设置、个性签名）
- [x] 好友申请相关
- [x] 消息本地化存储（获取历史消息）
- [x] 账号记录（本地存储）
- [ ] 消息云端存储（离线消息）
- [x] 发送文本消息
- [ ] 发送二进制消息（文件、图片等等）
- [ ] 多设备账号登录（windows、Linux、macOS、安卓、IOS等等）
- [ ] 多设备跨端文件传输
- [ ] 好友动态的发送和获取（类似于朋友圈）


以上为暂定的功能模块，后续会加入视频通话和远程控制模块。

## 后端架构

### 项目目录说明
> 注意由于本项目采用go1.18版本，使用的是go work来管理多个项目依赖，下面的所有文件夹都是一个单独的go mod项目。
```
├───chat_http    使用http相关服务的业务（如不需实时主动转发的）
├───chat_socket  使用websocket相关服务的业务（如需要实时主动转发的）
├───conf         整个项目的配置文件（包含数据库连接等）
├───logger       整个项目的logger配置
└───main         项目的运行main函数文件

```

### 整体业务流程

![image.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/0e91269725a64710a97344379c841379~tplv-k3u1fbpfcp-watermark.image?)

### 模块划分


![image.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/32d8edfa30934db3800808af8bd66fa6~tplv-k3u1fbpfcp-watermark.image?)

### 类图
> 以下是实时聊天系统的类图，目前只支持文字消息，也就是对应的TextMessage的转发，如果需要增加其他大型文件的转发，可以采取分块的方式。

![chat.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/142dc8ed549e44b6ba537de90bb628ac~tplv-k3u1fbpfcp-watermark.image?)

## 前端效果展示

### 登录和注册
![my_chat_login](https://user-images.githubusercontent.com/73544345/191010055-39e32a2d-5160-4fa5-9b01-6133ab07d5e5.gif)

### 用户信息
![image.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/41a94b33c79941c490d346c01e8b275f~tplv-k3u1fbpfcp-watermark.image?)

### 添加好友
![my_chat_add](https://user-images.githubusercontent.com/73544345/191010164-1af29169-ae6e-4097-9d54-dc094aba56b4.gif)

### 好友实时聊天以及历史聊天记录获取
![my_chat_talk](https://user-images.githubusercontent.com/73544345/191010342-b795d840-0597-4bdd-a61c-d61eb19a5950.gif)
