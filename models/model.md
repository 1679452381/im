# 集合列表

## 用户表
```json
{
    "account":"账号",
		"password":"密码",
		"nickname":"昵称",
		"sex":1, //0 未知  1-男 2-女
		"email":"邮箱",
		"avatar":"头像",
		"created_at":1,
		"updated_at":1,
}

```

## 消息表

```json
{
    "user_identity":"用户的唯一标识",
		"room_identity":"房间的唯一标识",
		"data":"发送的数据",
		"created_at":1,
		"updated_at":1,
}
```

## 房间表
```json

{
    "number":"房间号",
		"name":"房间名",
		"info":"房间简介",
		"user_identity":"创建者唯一标识",
		"created_at":1,
		"updated_at":1,
}
```

## 用户房间关联表

```json

{
    "user_identity":"用户唯一标识",
		"room_identity":"房间唯一标识",
		"message_identity":"消息唯一标识",
		"created_at":1,
		"updated_at":1,
}
```