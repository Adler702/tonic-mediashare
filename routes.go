package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"tonic-mediashare/structs"
)

var DiscordConfig *oauth2.Config

func addRoutes() {

	DiscordConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%v/user", config.BaseURL),
		ClientID:     config.DiscordClientId,
		ClientSecret: config.DiscordClientSecret,
		Scopes:       []string{discord.ScopeIdentify},
		Endpoint:     discord.Endpoint,
	}

	GETIndex()
	GETAuth(DiscordConfig)
	GETUser(DiscordConfig)
	POSTImageUpload()
	POSTFileUpload()
	POSTPasteUpload()
	POSTURLShortener()
	GETUrl()
}

func GETAuth(conf *oauth2.Config) {
	engine.GET("/auth", func(context *gin.Context) {
		context.Redirect(http.StatusTemporaryRedirect, conf.AuthCodeURL(config.State))
	})
}

func GETIndex() {
	engine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{
			"title":      "Tonic-MediaShare",
			"discordUrl": fmt.Sprintf("%v/auth", config.BaseURL),
		})
	})
}

func GETUser(conf *oauth2.Config) {
	engine.GET("/user", func(c *gin.Context) {
		if c.Request.FormValue("state") != config.State {
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.Writer.Write([]byte("State does not match."))
			return
		}
		token, err := conf.Exchange(context.Background(), c.Request.FormValue("code"))

		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		client := conf.Client(context.Background(), token)
		res, err := client.Get("https://discord.com/api/users/@me")

		if err != nil || res.StatusCode != 200 {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			if err != nil {
				c.Writer.Write([]byte(err.Error()))
			} else {
				c.Writer.Write([]byte(res.Status))
			}
			return
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Write([]byte(err.Error()))
			return
		}
		discorduser := structs.DiscordUser{}
		json.Unmarshal(body, &discorduser)
		var code string

		if AccountIsSaved(discorduser.ID) {
			code = getAccount(discorduser.ID).Code
		} else {
			code = randomString(discorduser.ID)
			saveAccount(discorduser.Username, code, discorduser.ID)
		}
		c.HTML(http.StatusOK, "user.html", gin.H{
			"title": "Tonic-MediaShare",
			"user":  discorduser.Username,
			"code":  code,
		})
	})
}

func POSTImageUpload() {
	engine.POST("/imageupload", func(context *gin.Context) {
		if !IsAuthorized(context.GetHeader("Authorization")) {
			context.HTML(http.StatusOK, "denied.html", gin.H{
				"title": "Tonic-MediaShare",
			})
			return
		}
		//
		//user := strings.Split(context.Request.Header.Get("Authorization"), " ")[0]
		//data := saveImage(user)
	})
}

func POSTFileUpload() {
	engine.POST("/fileupload", func(context *gin.Context) {
		if !IsAuthorized(context.GetHeader("Authorization")) {
			context.HTML(http.StatusOK, "denied.html", gin.H{
				"title": "Tonic-MediaShare",
			})
			return
		}
		//
	})
}

func POSTPasteUpload() {
	engine.POST("/pasteupload", func(context *gin.Context) {
		if !IsAuthorized(context.GetHeader("Authorization")) {
			context.HTML(http.StatusOK, "denied.html", gin.H{
				"title": "Tonic-MediaShare",
			})
			return
		}
		user := strings.Split(context.Request.Header.Get("Authorization"), " ")[0]
		data := savePaste(context.PostForm("paste"), user, randomString(strconv.Itoa(int(time.Now().UnixMicro()))))
		context.Writer.WriteString(fmt.Sprintf("%v/paste/%v", config.BaseURL, data.Code))
	})
}

func POSTURLShortener() {
	engine.POST("/urlshortener", func(context *gin.Context) {
		if !IsAuthorized(context.GetHeader("Authorization")) {
			context.HTML(http.StatusOK, "denied.html", gin.H{
				"title": "Tonic-MediaShare",
			})
			return
		}
		user := strings.Split(context.Request.Header.Get("Authorization"), " ")[0]
		data := saveURLShort(context.PostForm("link"), user, randomString(strconv.Itoa(int(time.Now().UnixMicro()))))
		context.Writer.WriteString(fmt.Sprintf("%v/url/%v", config.BaseURL, data.Code))
	})
}

func GETUrl() {
	engine.GET("/url/:code", func(c *gin.Context) {
		code := c.Params.ByName("code")
		destination := getShortedUrl(code).Destination
		parse, err := url.Parse(destination)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "denied.html", gin.H{
				"title": "Tonic-MediaShare",
			})
			return
		}
		c.HTML(http.StatusOK, "url.html", gin.H{
			"title":       "Tonic-MediaShare",
			"destination": destination,
			"domain":      strings.TrimPrefix(parse.Host, "www."),
		})
	})
}
