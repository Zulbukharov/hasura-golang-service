package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type GithubAuthStruct struct {
	accessToken string `json: "access_token"`
	tokenType   string `json: "type"`
}

var clientID = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")

// Cors ...
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// GithubAuth ...
func GithubAuth(c *gin.Context) {
	url := "https://github.com/login/oauth/access_token"
	fmt.Println("URL:>", url)
	code := c.Query("code")
	var u = fmt.Sprintf(`%s?client_id=%s&client_secret=%s&code=%s&state=sup`, url,
		clientID, clientSecret, code)
	req, err := http.NewRequest("POST", u, nil)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	if resp.Status != "200 OK" {
		fmt.Println("hi")
		c.JSON(404, gin.H{"status": "Error"})
	}
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	var t GithubAuthStruct
	json.Unmarshal(body, &t)
	fmt.Println("[unmarshall]", t)
	fmt.Println("response Body:", string(body))
	c.JSON(200, string(body))
}
