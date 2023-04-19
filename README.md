# im 即时通讯

# 技术栈 
gin + websocket + MongoDB

## 核心包
https://github.com/gorilla/websocket

## 包安装
```shell
go get -u github.com/gin-gonic/gin
go get github.com/gorilla/websocket
go get go.mongodb.org/mongo-driver/mongo
go get github.com/dgrijalva/jwt-go
go get github.com/google/uuid
go get github.com/redis/go-redis/v9
# 发送email 验证码
go get  

```

## 安装docker
https://blog.csdn.net/weixin_51351637/article/details/128006765
## docker 安装 mo·ngoDB

```shell
docker run -d --name some-mongo -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=admin -p 27017:27017 mongo
docker run -p 6379:6379 -d --name redis01 --restart=always redis
```

## 数据库设计


##  接口实现
* 登录
* 注册
* 用户详情
* 验证码发送功能
* 发送，接受消息

* 查询
* 添加好友 
* 删除好友
* 建立群聊
* 解散群聊


* 用户详情
```go
//根据id查询用户信息时 要传ObjectId()

func GetUserBasicByIdentity(identity primitive.ObjectID) (*UserBasic, error) {
	ub := new(UserBasic)
	err := Mongo.Collection(UserBasic{}.CollectionName()).FindOne(context.Background(), bson.D{{"_id", identity}}).Decode(ub)
	return ub, err
}

// 用户详情
func UserDetail(c *gin.Context) {
u, _ := c.Get("user_claims")

//断言
uc := u.(*utils.UserClaims)
ub, err := models.GetUserBasicByIdentity(uc.Identity)
if err != nil {
c.JSON(http.StatusOK, gin.H{"msg": err.Error(), "code": -1})
return
}
c.JSON(http.StatusOK,
gin.H{
"msg":  "数据加载成功",
"code": 200,
"data": ub,
})
}

```

* 验证码发送功能
```go

// 发送邮箱验证码
func SendEmialCode(toEmial, code string) error {
	e := email.NewEmail()
	e.From = "GET <17700611471@163.com>"
	e.To = []string{toEmial}
	e.Subject = "验证码已发送"
	e.HTML = []byte("您的验证码：<b>" + code + "</b>")
	//err := e.Send("smtp.163.com:465", smtp.PlainAuth("", "15660589213@163.com", "DSBZHQSKFWQVDSVK", "smtp.163.com"))
	//返回EOF 关闭SSL重试
	return e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "17700611471@163.com", "DSBZHQSKFWQVDSVK", "smtp.163.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})

}

// 生成验证码
func GetCode() string {
    //初始化随机数种子
    rand.Seed(time.Now().UnixNano())
    
    var code string
    for i := 0; i < 4; i++ {
    //strconv.Itoa  将整数型数据转换为字符串型数据
    code += strconv.Itoa(rand.Intn(10))
    }
    return code
}

// 发送验证码
func SendCode(c *gin.Context) {
    email := c.PostForm("email")
    fmt.Println(email)
    if email == "" {
    c.JSON(http.StatusOK, gin.H{"msg": "邮箱不能为空", "code": -1})
    return
    }
	
    num, err := models.GetUserBasicByEmail(email)
    if err != nil {
    c.JSON(http.StatusOK, gin.H{"msg": "数据库查询错误", "code": -1})
    return
    }
    if num > 0 {
    c.JSON(http.StatusOK, gin.H{"msg": "该邮箱已被注册", "code": -1})
    return
    }
    //随机生成四位数验证码
    code := "121345"
    err = utils.SendEmialCode(email, code)
    if err != nil {
    c.JSON(http.StatusOK, gin.H{"msg": "验证码发送失败", "code": -1, "err": err.Error()})
    return
    }
    c.JSON(http.StatusOK, gin.H{"msg": "验证码发送成功", "code": 200})
}
```

* 使用http搭建websocket服务
```go

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options
var ws = map[*websocket.Conn]struct{}{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()
	ws[c] = struct{}{}
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		for conn := range ws {
			err = conn.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}

	}
}

func TestWebsocketServer(t *testing.T) {
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
```

* 使用gin搭建websocket服务
```go


func TestGinWebsocketServer(t *testing.T) {
	r := gin.Default()
	r.GET("/echo", func(c *gin.Context) {
		echo(c.Writer, c.Request)
	})
	r.Run(":8080")
}
```

* 发送，接受消息
```go

var upGrader = websocket.Upgrader{} // use default options
var wc = make(map[string]*websocket.Conn)

func WebsocketMessage(c *gin.Context) {
	//建立websocket连接
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "系统错误", "code": -1})
		return
	}
	defer conn.Close()
	uc := c.MustGet("user_claims").(*utils.UserClaims)
	wc[uc.Identity] = conn
	for {
		ms := new(define.MessageStruct)
		//读数据
		err := conn.ReadJSON(ms)
		if err != nil {
			log.Printf("[ERROR] %v", err.Error())
			return
		}
		//判断用户是否属于房间
		_, err = models.GetUserRoomByUserIdentityRoomIdentity(uc.Identity, ms.RoomIdentity)
		if err != nil {
			log.Printf("user_identity:%v,room_identity:%v is not exits", uc.Identity, ms.RoomIdentity)
			return
		}
		//TODO 保存消息
		mb := models.MessageBasic{
			UserIdentity: uc.Identity,
			RoomIdentity: ms.RoomIdentity,
			Data:         ms.Message,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		err = models.InsertMessageBasic(mb)
		if err != nil {
			log.Printf("DB ERROR:%v", err.Error())
		}
		//TODO 获取房间在线用户
		urs, err := models.GetUserRoomByRoomIdentity(ms.RoomIdentity)
		if err != nil {
			log.Printf("[ERROR]:%v", err.Error())
			return
		}
		for _, room := range urs {
			if cc, ok := wc[room.UserIdentity]; ok {
				err = cc.WriteMessage(websocket.TextMessage, []byte(ms.Message))
				if err != nil {
					log.Println("write:", err)
					break
				}
			}
		}

	}
}

```

* 获取消息列表
```go

func ChatList(c *gin.Context) {
	roomIdentity := c.Query("room_identity")
	if roomIdentity == "" {
		c.JSON(http.StatusOK, gin.H{"msg": "room_identoty不能为空", "code": -1})
		return
	}
	uc := c.MustGet("user_claims").(*utils.UserClaims)
	// 判断用户是否属于该房间
	_, err := models.GetUserRoomByUserIdentityRoomIdentity(uc.Identity, roomIdentity)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "非法访问", "code": -1, "err": err.Error()})
		return
	}
	//分页
	//s：需要转换为整数的字符串。
	//base：指定进制，比如10表示十进制，16表示十六进制。
	//bitSize：指定整数的位大小，通常使用0表示按s中的数值返回一个合适的位大小。
	pageIndex, _ := strconv.ParseInt(c.Query("page_index"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 32)
	skip := (pageIndex - 1) * pageSize

	mbs, err := models.GetMessageBasicByRoomIdentity(roomIdentity, &pageSize, &skip)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "数据库错误", "code": -1, "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "聊天列表", "code": -1, "data": gin.H{
		"list": mbs,
	}})
}
// 获取房间信息列表
func GetMessageBasicByRoomIdentity(roomIdentity string, limit, skip *int64) ([]*MessageBasic, error) {
mbs := make([]*MessageBasic, 0)
cursor, err := Mongo.Collection(MessageBasic{}.CollectionName()).Find(context.Background(),
bson.D{{"room_identity", roomIdentity}},
&options.FindOptions{
Limit: limit,
Skip:  skip,
Sort: bson.D{
{"created_at", -1},
},
})
if err != nil {
return nil, err
}
for cursor.Next(context.Background()) {
mb := new(MessageBasic)
err = cursor.Decode(mb)
if err != nil {
return nil, err
}
mbs = append(mbs, mb)
}
return mbs, nil
}
```

* 注册
```go

func Register(c *gin.Context) {
	//	获取注册信息
	account := c.PostForm("account")
	password := c.PostForm("password")
	email := c.PostForm("email")
	code := c.PostForm("code")
	if account == "" || password == "" || email == "" || code == "" {
		c.JSON(http.StatusOK, gin.H{"msg": "参数有误", "code": -1})
		return
	}
	//	检查账号是否已被注册
	num, err := models.GetUserBasicByAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "[DB ERROR]", "code": -1})
		return
	}
	if num > 0 {
		c.JSON(http.StatusOK, gin.H{"msg": "该账号已被注册", "code": -1})
		return
	}
	//  校验验证码 从redis中取验证码
	registerCode, err := models.RDB.Get(context.Background(), define.RegisterPrefix+email).Result()
	fmt.Println(code, registerCode)
	if err != nil {
		log.Printf("%v", err.Error())
		c.JSON(http.StatusOK, gin.H{"msg": "验证码超时,请重新获取", "code": -1})
		return
	}
	if registerCode != code {
		c.JSON(http.StatusOK, gin.H{"msg": "验证码不正确", "code": -1})
		return
	}
	//生成数据
	ub := &models.UserBasic{
		Identity:  utils.GetUUID(),
		Account:   account,
		Password:  utils.GetMd5(password),
		NickName:  "",
		Sex:       0,
		Email:     email,
		Avatar:    "",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	//	校验通过，存储数据
	err = models.InsertUserBasic(ub)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "系统错误", "code": -1})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "添加成功", "code": 200, "data": ub})
}
```
* 通过账号查询用户信息
```go
// 通过账号查询用户信息
func SearchByAccount(c *gin.Context) {
	account := c.PostForm("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{"msg": "参数有误", "code": -1})
		return
	}
	//根据账号查询信息
	ub, err := models.GetUserBasicByAccount(account)
	if err != nil {
		log.Printf("%v", err.Error())
		c.JSON(http.StatusOK, gin.H{"msg": "数据库错误", "code": -1})
		return
	}
	ui := &UserInfo{
		Account:  ub.Account,
		NickName: ub.NickName,
		Sex:      ub.Sex,
		Avatar:   ub.Avatar,
		IsFriend: false,
	}
	//判断是否是好友
	uc := c.MustGet("user_claims").(*utils.UserClaims)
	isFriend := models.IsFriend(uc.Identity, ub.Identity)
	ui.IsFriend = isFriend
	c.JSON(http.StatusOK, gin.H{"msg": "查询成功", "code": 200, "data": ui})
}
// 判断是否是好友
func IsFriend(identity1, identity2 string) bool {
    cursor, err := Mongo.Collection(UserRoom{}.CollectionName()).Find(context.Background(), bson.D{{"user_identity", identity1}})
    if err != nil {
        log.Printf("%v", err.Error())
        return false
    }
    roomIdentities := make([]string, 0)
    for cursor.Next(context.Background()) {
        ur := new(UserRoom)
        err := cursor.Decode(ur)
        if err != nil {
            log.Printf("[DB ERROR]%v", err.Error())
            return false
        }
        //房间类型有两种 1私聊 2 群聊 只有是两个人的私聊房间才算好友
        if ur.RoomType == 1 {
            roomIdentities = append(roomIdentities, ur.RoomIdentity)
        }
    }
    count, err := Mongo.Collection(UserRoom{}.CollectionName()).CountDocuments(context.Background(),
    bson.M{"user_identity": identity2, "room_identity": bson.M{"$in": roomIdentities}})
    if err != nil {
        log.Printf("%v", err.Error())
        return false
    }
    if count > 0 {
        return true
    }
    return false
}

```

* 添加好友 
```go
func AddFriend(c *gin.Context) {
	account := c.PostForm("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{"msg": "参数有误", "code": -1})
		return
	}
	//根据账号查询信息
	ub, err := models.GetUserBasicByAccount(account)
	if err != nil {
		log.Printf("%v", err.Error())
		c.JSON(http.StatusOK, gin.H{"msg": "数据库错误", "code": -1})
		return
	}
	uc := c.MustGet("user_claims").(*utils.UserClaims)
	//查看是否是好友
	isFriend := models.IsFriend(uc.Identity, ub.Identity)
	if isFriend {
		c.JSON(http.StatusOK, gin.H{"msg": "已经是好友，无需添加", "code": -1})
		return
	}
	//创建一个新的聊天室
	roomIdentity := utils.GetUUID()
	rb := &models.RoomBasic{
		Identity:     roomIdentity,
		Info:         "好友",
		UserIdentity: uc.Identity,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	err = models.InsertOneRoomBasic(rb)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "添加好友失败", "code": -1})
		return
	}
	rooms := make([]interface{}, 0)
	ur1 := &models.UserRoom{
		UserIdentity: uc.Identity,
		RoomIdentity: roomIdentity,
		RoomType:     1,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	rooms = append(rooms, ur1)
	ur2 := &models.UserRoom{
		UserIdentity: ub.Identity,
		RoomIdentity: roomIdentity,
		RoomType:     1,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	rooms = append(rooms, ur2)
	//	添加聊天室
	err = models.InsertUserRooms(rooms)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "添加好友失败", "code": -1})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "添加好友成功", "code": -1})
}
```
* 删除好友
```go
func DeleteFriend(c *gin.Context) {
	account := c.PostForm("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{"msg": "参数有误", "code": -1})
		return
	}
	//根据账号查询信息
	ub, err := models.GetUserBasicByAccount(account)
	if err != nil {
		log.Printf("%v", err.Error())
		c.JSON(http.StatusOK, gin.H{"msg": "数据库错误", "code": -1})
		return
	}
	uc := c.MustGet("user_claims").(*utils.UserClaims)
	//查看是否是好友
	isFriend := models.IsFriend(uc.Identity, ub.Identity)
	if !isFriend {
		c.JSON(http.StatusOK, gin.H{"msg": "不是好友，非法操作", "code": -1})
		return
	}
	//查好友房间号
	rbIdentity := models.UserRoomSearchFriend(uc.Identity, ub.Identity)
	fmt.Println(rbIdentity)
	if rbIdentity == "" {
		c.JSON(http.StatusOK, gin.H{"msg": "数据库错误", "code": -1})
		return
	}
	//删除房间
	err = models.DeleteOneRoomBasic(rbIdentity)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "系统错误，删除失败", "code": -1})
		return
	}
	//删除ur关系
	err = models.DeleteUserRoom(rbIdentity)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "系统错误，删除失败", "code": -1})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "删除好友成功", "code": -1})
}

```