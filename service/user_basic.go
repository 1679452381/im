package service

import (
	"Im/define"
	"Im/models"
	"Im/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// 登录
func Login(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
	if account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{"msg": "账号或密码不能为空", "code": -1})
	}
	ub, err := models.GetUserBasicByAccountPassword(account, utils.GetMd5(password))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "账号或密码不正确", "code": -1, "err": err.Error()})
		return
	}
	token, err := utils.GenerateToken(ub.Identity, ub.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "系统错误", "code": -1})
		return
	}
	c.JSON(http.StatusOK,
		gin.H{
			"msg":  "登陆成功",
			"code": 200,
			"data": gin.H{
				"token": token,
			},
		})
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

// 发送验证码
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{"msg": "邮箱不能为空", "code": -1})
		return
	}
	//查询数据库
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
	code := utils.GetCode()
	fmt.Println(code)
	//将验证码存到redis
	err = models.RDB.Set(context.Background(), define.RegisterPrefix+email, code, time.Second*60).Err()
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	err = utils.SendEmialCode(email, code)
	if err != nil {
		log.Printf("[ERROR]:%v", err)
		c.JSON(http.StatusOK, gin.H{"msg": "验证码发送失败", "code": -1, "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "验证码发送成功", "code": 200})
}

// 注册
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
	num, err := models.GetNumUserBasicByAccount(account)
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
	ui := &define.UserInfo{
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

// 添加好友
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

// 删除好友
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
