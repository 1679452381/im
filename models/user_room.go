package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

type UserRoom struct {
	UserIdentity    string `bson:"user_identity"`
	RoomIdentity    string `bson:"room_identity"`
	MessageIdentity string `bson:"message_identity"`
	RoomType        int    `bson:"room_type"` //0-删除好友 1-私人聊天 2-群聊
	CreatedAt       int64  `bson:"created_at"`
	UpdatedAt       int64  `bson:"updated_at"`
}

func (UserRoom) CollectionName() string {
	return "user_room"
}

func GetUserRoomByUserIdentityRoomIdentity(userIdentity, roomIdentity string) (*UserRoom, error) {
	ur := new(UserRoom)
	err := Mongo.Collection(UserRoom{}.CollectionName()).FindOne(context.Background(),
		bson.D{{"user_identity", userIdentity}, {"room_identity", roomIdentity}}).Decode(ur)
	return ur, err
}

func GetUserRoomByRoomIdentity(roomIdentity string) ([]*UserRoom, error) {
	cursor, err := Mongo.Collection(UserRoom{}.CollectionName()).Find(context.Background(), bson.D{{"room_identity", roomIdentity}})
	if err != nil {
		return nil, err
	}
	urs := make([]*UserRoom, 0)
	for cursor.Next(context.Background()) {
		ur := new(UserRoom)

		err := cursor.Decode(ur)
		if err != nil {
			return nil, err
		}
		urs = append(urs, ur)
	}
	return urs, nil
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

// 添加好友
func InsertUserRooms(rooms []interface{}) error {
	_, err := Mongo.Collection(UserRoom{}.CollectionName()).InsertMany(context.Background(), rooms)
	if err != nil {
		return err
	}
	return nil
}

// 好友房间号
func UserRoomSearchFriend(identity1, identity2 string) string {
	cursor, err := Mongo.Collection(UserRoom{}.CollectionName()).Find(context.Background(), bson.D{{"user_identity", identity1}})
	if err != nil {
		log.Printf("%v", err.Error())
		return ""
	}
	roomIdentities := make([]string, 0)
	for cursor.Next(context.Background()) {
		ur := new(UserRoom)
		err := cursor.Decode(ur)
		if err != nil {
			log.Printf("[DB ERROR]%v", err.Error())
			return ""
		}
		//房间类型有两种 1私聊 2 群聊 只有是两个人的私聊房间才算好友
		if ur.RoomType == 1 {
			roomIdentities = append(roomIdentities, ur.RoomIdentity)
		}
	}
	urData := new(UserRoom)
	err = Mongo.Collection(UserRoom{}.CollectionName()).FindOne(context.Background(),
		bson.M{"user_identity": identity2, "room_identity": bson.M{"$in": roomIdentities}}).Decode(urData)
	if err != nil {
		return ""
	}
	return urData.RoomIdentity
}

//删除好友 -删除关联房间

func DeleteUserRoom(rbIdentity string) error {
	_, err := Mongo.Collection(UserRoom{}.CollectionName()).DeleteMany(context.Background(), bson.D{{"room_identity", rbIdentity}})
	return err
}
