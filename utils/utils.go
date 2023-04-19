package utils

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// 生成md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var myKey = []byte("im")

// 生成token
func GenerateToken(identity, email string) (string, error) {
	//将string 转化为ObjectID
	//objectId, err := primitive.ObjectIDFromHex(identity)
	//if err != nil {
	//	return "", err
	//}
	UserClaim := &UserClaims{
		identity,
		email,
		jwt.StandardClaims{}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析token
func AnalyToken(tokenString string) (*UserClaims, error) {
	UserClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, UserClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analy token error:%v", err)
	}
	return UserClaim, nil
}

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

//生成唯一标识

type UUID [16]byte

func GetUUID() string {
	return fmt.Sprintf("%x", uuid.New())

}
