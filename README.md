# X_UGC

一个简单的UGC平台.

## 技术选型

Gin,GORM,MySQL,Redis,RabbitMQ

## 项目结构

```
├─assets
├─biz
│  ├─dal
│  │  ├─mysql
│  │  ├─rabbitmq
│  │  └─redis
│  ├─handler
│  ├─model
│  ├─mw
│  ├─router
│  └─service
├─cmd
├─conf
├─pkg
│  ├─common
│  │  ├─bloomfilter
│  │  ├─email
│  │  ├─ffmpeg
│  │  ├─jwt
│  │  └─pwd
│  └─websocket
```

## 功能展示

![主页测试.jpg](assets/主页测试.jpg)

![朋友内容测试.jpg](assets/朋友内容测试.jpg)

![消息测试.jpg](assets/消息测试.jpg)

![点赞评论测试.jpg](assets/点赞评论测试.jpg)