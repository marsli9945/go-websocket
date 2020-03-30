# go-websocket

## 路由
- /websocket ws连接地址
- /websocket/list 查看在线用户
- /websocket/test 测试用html页面
- /websocket/push 消息推送接口

## 使用
- 使用docker进行部署
- 默认使用7777端口
- 首次部署在项目目录make build && make run
- 更新重新构建make rebuild

## 数据格式

```
{
    "name": "Tom",
    "data": {}
}
```