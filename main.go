package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"tonic-mediashare/structs"
)

var config structs.Config
var engine *gin.Engine

func main() {
	fmt.Println("starting tonic-mediashare")
	loadConfig()
	engine = gin.Default()
	connectMongoDB()
	addRoutes()
	engine.Static("/assets", "./assets")
	engine.LoadHTMLGlob("./templates/*")
	engine.SetTrustedProxies(nil)
	fmt.Println("Starting Gin Default Engine on Port", config.Port)
	err := engine.Run(fmt.Sprintf(":%v", config.Port))
	if err != nil {
		fmt.Println("Engine failed to Start")
		return
	}
}

func loadConfig() {
	raw, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("Error occurred while reading config")
		return
	}
	json.Unmarshal(raw, &config)
}

func randomString(seed string) string {
	i, _ := strconv.Atoi(seed)
	rand.Seed(int64(i))
	length := 8
	letters := "abcdefghijklmnopqrstuvwxyz0123456789"
	//letters := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func formatSizeString(size int64) string {
	types := map[int8]string{0: "B", 1: "KB", 2: "MB", 3: "GB"}
	var result int64 = 1
	for true {
		var filetype int8 = 0
		if size >= 1 && size < 1024 {
			for size >= 1 {
				result = size
				return fmt.Sprintf("%v%v", result, types[filetype])
			}
		}
		for size >= 1024 {
			size -= 1024
			result++
		}
		filetype++
		for result >= 1000 {
			result -= 1000
			filetype++
		}
		return fmt.Sprintf("%v%v", result, types[filetype])
	}
	return "0B"
}
