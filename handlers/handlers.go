package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Zulbukharov/hasura-golang-service/db"
	"github.com/Zulbukharov/hasura-golang-service/services"
	"github.com/Zulbukharov/hasura-golang-service/tools"
	"github.com/gin-gonic/gin"
)

// GithubAuthStruct ...
type GithubAuthStruct struct {
	AccessToken string `json:"access_token"`
	JwtToken    string `json:"jwt_token"`
}

type GithubUser struct {
	UserLogin string `json:"login"`
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

func getGithubUser(token string) ([]byte, int) {
	// https://api.github.com
	url := "https://api.github.com/user"
	tokenHeader := fmt.Sprintf("token %s", token)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", tokenHeader)

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return []byte{'0'}, 404
	}
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return body, 200
}

// GithubAuth ...
func getGithubToken(code string) ([]byte, int) {
	url := "https://github.com/login/oauth/access_token"
	fmt.Println("URL:>", url)
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
		return []byte{'0'}, 404
	}
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return (body), 200
}

func Auth(c *gin.Context) {
	code := c.Query("code")
	res, status := getGithubToken(code)
	if status != 200 {
		c.JSON(404, nil)
		return
	}
	var t GithubAuthStruct
	json.Unmarshal(res, &t)
	fmt.Println("[unmarshall]", t)
	// var user GithubUser
	var user GithubUser
	u, s := getGithubUser(t.AccessToken)
	if s != 200 {
		c.JSON(404, nil)
		return
	}
	err := json.Unmarshal(u, &user)
	fmt.Println("[unmarshall err]", err)
	fmt.Println("[user]", user)
	client := db.GetClient()
	// get user from hasura
	_, err = services.GetUser(client, user.UserLogin)
	if err != nil {
		fmt.Println(err)
		// if not exist add user
		er := services.InsertUser(client, user.UserLogin)
		if er != nil {
			c.JSON(404, nil)
			return
		}
	}
	//send acess token to jwt generator
	token, err := tools.GenerateToken(user.UserLogin)
	t.JwtToken = token
	js, err := json.Marshal(t)
	if err != nil {
		c.JSON(404, nil)
		return
	}
	c.JSON(200, string(js))
}
