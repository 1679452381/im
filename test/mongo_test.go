package test

import (
	"Im/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

// 测试MongoDB链接
func TestFindOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username:    "admin",
		Password:    "admin",
		PasswordSet: false,
	}).ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		t.Fatal(err)
	}
	//连接数据库
	db := client.Database("im")

	cur := db.Collection("user_basic")

	ub := new(models.UserBasic)
	err = cur.FindOne(context.Background(), bson.D{}).Decode(ub)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ub)
}

// 测试MongoDB链接
func TestFind(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username:    "admin",
		Password:    "admin",
		PasswordSet: false,
	}).ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		t.Fatal(err)
	}
	//连接数据库
	db := client.Database("im")

	cur, err := db.Collection("user_room").Find(context.Background(), bson.D{})
	urs := make([]*models.UserRoom, 0)
	for cur.Next(context.Background()) {
		ur := new(models.UserRoom)
		err := cur.Decode(ur)
		if err != nil {
			t.Fatal(err)
		}
		urs = append(urs, ur)
	}
	for _, ur := range urs {
		fmt.Println(ur)
	}
}
