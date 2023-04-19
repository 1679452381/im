package service

import (
	"Im/models"
	"Im/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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
