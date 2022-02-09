Jim相关文档
===
###运行模式
* local
* grpc
* kafka
### 支持连接方式
* websocket
* tcp
###消息协议
* 认证消息：
```json 
{"cmd":"auth.req","data":{"token":"client_1"}}
```
* 认证失败：
```json 
{"cmd":"auth.fail","data":{"tip":"auth fail"}}
```
* 未认证：
```json 
{"cmd":"auth.not","data":{"tip":"not auth"}}
```
* 单聊文本：
```json 
{"cmd":"chat.c2c.txt","data":{"id":0,"content":"hello client 2","to_receiver_id":"client_2"}}  
```

###实际效果
|说明|截图|
|-----|------|
|单聊发送消息|![alt 单聊发送消息](./c2c_send_txt.png "单聊发送消息")|
|单聊接收消息|![alt 单聊接收消息](./c2c_receive_txt.png "单聊接收消息")|