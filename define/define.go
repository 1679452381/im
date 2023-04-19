package define

import "os"

var MailPassword = os.Getenv("MailPassword")

type MessageStruct struct {
	RoomIdentity string `json:"room_identity"`
	Message      string `json:"message"`
}

var RegisterPrefix = "TOKEN_"

type UserInfo struct {
	Account  string `json:"account"`
	NickName string `json:"nickName"`
	Sex      int    `json:"sex"`
	Avatar   string `json:"avatar"`
	IsFriend bool   `json:"is_friend"`
}
