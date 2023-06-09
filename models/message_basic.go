package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageBasic struct {
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	Data         string `bson:"data"`
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

func (MessageBasic) CollectionName() string {
	return "message_basic"
}

// 插入一条信息
func InsertMessageBasic(mb MessageBasic) error {
	_, err := Mongo.Collection(MessageBasic{}.CollectionName()).InsertOne(context.Background(), mb)
	return err
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
