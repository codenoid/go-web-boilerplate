package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/codenoid/go-web-boilerplate/web-base/structs/user"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	loadConfig()

	appName = os.Getenv("APP_NAME")

	initMongoDB()
	initRedis()
}

func loadConfig() {
	err := godotenv.Load("config.env")
	if err != nil {
		panic(err)
	}
}

func initMongoDB() {
	if client, err := qmgo.NewClient(context.TODO(), &qmgo.Config{Uri: os.Getenv("MONGO_URI")}); err != nil {
		panic(err)
	} else {
		mainDB = client.Database(appName + "_web")
	}

	mainDB.Collection("users").CreateOneIndex(context.TODO(), options.IndexModel{
		Key:    []string{"username"},
		Unique: true,
	})

	if err := mainDB.Collection("users").Find(context.TODO(), bson.M{"username": "admin"}).One(&user.User{}); err == qmgo.ErrNoSuchDocuments {
		// Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("@admin1234"), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		mainDB.Collection("users").InsertOne(context.TODO(), user.User{
			Name:      "admin",
			Username:  "admin",
			Email:     "admin@localhost",
			Password:  string(hashedPassword),
			PrivateID: uuid.New().String(),
			PublicID:  uuid.New().String(),
			CreatedAt: time.Now().Unix(),
		})
	}
}

func initRedis() {
	var rdbnum int = 0 // default database

	if envVal := os.Getenv("REDIS_DB"); envVal != "" {
		customNum, err := strconv.Atoi(envVal)
		if err != nil {
			panic(err)
		}
		rdbnum = customNum
	}

	sessionDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWD"),
		DB:       rdbnum,
	})
	if _, err := sessionDB.Ping().Result(); err != nil {
		panic(err)
	}
}
