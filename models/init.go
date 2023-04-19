package models

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var Mongo = MongoDBInit()

// MongoDB
func MongoDBInit() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username:    "admin",
		Password:    "admin",
		PasswordSet: false,
	}).ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		fmt.Println("MongoDB 连接失败")
	}
	//连接数据库
	db := client.Database("im")

	return db
}

var RDB = RedisInit()

// Redis
func RedisInit() *redis.Client {
	//建立连接
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
