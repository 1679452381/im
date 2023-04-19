package service

import (
	"Im/define"
	"Im/models"
	"Im/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

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
	//这两个方法的主要区别在于，当获取的键值对不存在时，MustGet方法会直接抛出异常，而Get方法则只会返回一个空值
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

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6IjY0M2Q3YWUxZWU2ZTE3NzVmYzA4ODYxMSIsImVtYWlsIjoiMTc3MDA2MTE0NzFAMTYzLmNvbSJ9.EfbSLom8Z52kJstWvqIedt7XtIN70zG96UFjcvCloT0

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6IjY0M2U1ZDE3N2QyOWJhZmFhMDBlNWYwMSIsImVtYWlsIjoiMTY3OTQ1MjM4MUBxcS5jb20ifQ.vpxzCToH5eZkix_v7Rv5m8txUEf-lAU1Gayl2oN6sOo
