package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"time"
	"tonic-mediashare/structs"
)

var database *mongo.Database

func connectMongoDB() {
	fmt.Println("Connecting to MongoDB")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURL))
	if err != nil {
		log.Fatal("Failed connecting to MongoDB:", err)
	}
	database = client.Database("mediashare")
}

func AccountIsSaved(user string) bool {
	account := structs.Account{Name: user}
	err := database.Collection("accounts").FindOne(context.TODO(), bson.D{{"id", user}}).Decode(&account)
	return err == nil
}

func saveAccount(name string, code string, id string) {
	account := structs.Account{
		Name: name,
		Code: code,
		Id:   id,
	}
	database.Collection("accounts").InsertOne(context.TODO(), account)
}

func getAccount(user string) structs.Account {
	account := structs.Account{Name: user}
	database.Collection("accounts").FindOne(context.TODO(), bson.D{{"id", user}}).Decode(&account)
	return account
}

func getDate() string {
	// DD:MM:YYY HH.mm:ss
	now := time.Now()
	return fmt.Sprintf("%v.%v.%v %v:%v:%v", now.Day(), int(now.Month()), now.Year(), now.Hour(), now.Minute(), now.Second())
}

func IsAuthorized(authorization string) bool {
	if len(strings.Split(authorization, " ")) != 2 {
		return false
	}
	data := strings.Split(authorization, " ")
	if getAccount(data[0]).Code == data[1] {
		return true
	}
	return false
}

func saveURLShort(url string, user string, code string) structs.UrlData {
	data := structs.UrlData{
		User:        user,
		Destination: url,
		Code:        code,
	}
	database.Collection("urls").InsertOne(context.TODO(), data)
	return data
}

func getShortedUrl(code string) structs.UrlData {
	data := structs.UrlData{Code: code}
	database.Collection("urls").FindOne(context.TODO(), bson.D{{"code", code}}).Decode(&data)
	return data
}

func savePaste(paste string, user string, code string) structs.PasteData {
	data := structs.PasteData{
		User:       user,
		Data:       paste,
		Code:       code,
		UploadData: getDate(),
	}
	database.Collection("pastes").InsertOne(context.TODO(), data)
	return data
}

func saveImage(user string, filename string, code string, filetype string, filesize int64, uploaddate string) structs.FileData {
	data := structs.FileData{
		User:       user,
		Filename:   filename,
		Code:       code,
		Filetype:   filetype,
		Filesize:   filesize,
		UploadDate: uploaddate,
	}
	database.Collection("images").InsertOne(context.TODO(), data)
	return data
}
