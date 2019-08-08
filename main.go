package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	gql "github.com/btubbs/garphunql"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"io/ioutil"
)

type User struct {
	ID   string `json: "id"`
	Name string `json: "name"`
}

type Arguments struct {
	Eq string `json: "_eq"`
}

func (u User) Format() (string, error) {
	v := fmt.Sprintf(`{id: "%s", name: "%s"}`, u.ID, u.Name)
	return v, nil
}

func (a Arguments) Format() (string, error) {
	v := fmt.Sprintf(`{id: {_eq: "%s"}}`, a.Eq)
	return v, nil
}

func InsertUser(client *gql.Client) (err error) {
	//	mutation {
	//			insert_users(objects: {id: "", name: ""}) {
	//				affected_rows
	//		}
	//	}

	var w gql.ArgumentFormatter = User{"22", "abl"}
	d := map[string]interface{}{"objects": w}
	mutationQuery := gql.GraphQLField{
		Name:      "insert_users",
		Arguments: d,
		Fields: []gql.GraphQLField{
			{Name: "affected_rows"},
		},
	}
	err = client.Mutation(mutationQuery)
	fmt.Println(err)
	return err
}

func GetUser(client *gql.Client) (err error) {
	//	{
	//		users(where: {id: {_eq: "0"}}) {
	//			id
	//			name
	//		}
	//	}

	var me []User
	var eq gql.ArgumentFormatter = Arguments{"0"}
	arguments := map[string]interface{}{"where": eq}
	getQuery := gql.GraphQLField{
		Name:      "users",
		Arguments: arguments,
		Fields: []gql.GraphQLField{
			{Name: "id"},
			{Name: "name"},
		},
		Dest: &me,
	}

	err = client.Query(getQuery)
	fmt.Println(err)
	fmt.Println(me)
	return err
}

func SimpleQuery(client *gql.Client) {
	var me []User

	/*
		{
		users {
				id
				name
			}
		}
	*/
	myField := gql.Field("users", gql.Field("id"), gql.Field("name"))
	err := client.Query(
		myField(gql.Dest(&me)),
	)
	fmt.Println(err, me)
}

var HASURA_GRAPHQL_ADDRESS = os.Getenv("HASURA_GRAPHQL_ADDRESS")
var HASURA_GRAPHQL_ADMIN_SECRET = os.Getenv("HASURA_GRAPHQL_ADMIN_SECRET")
var CLIENT_ID = os.Getenv("CLIENT_ID")
var CLIENT_SECRET = os.Getenv("CLIENT_SECRET")
var jwtSecret = []byte("opa")

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

type Claims struct {
	HasuraClaims map[string]interface{} `json:"https://hasura.io/jwt/claims"`
	jwt.StandardClaims
}

func GenerateToken(id string) (string, error) {
	nowTime := time.Now()
	// expireTime := nowTime.Add(3 * time.Hour)
	expireTime := nowTime.Add(1 * time.Hour)
	a := make([]string, 1)
	a[0] = "user"
	claims := &Claims{
		map[string]interface{}{
			"x-hasura-default-role":  "user",
			"x-hasura-allowed-roles": a,
			"x-hasura-user-id":       "1",
		},
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "test",
		},
	}
	fmt.Println(claims.HasuraClaims)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	// fmt.Println("[jwt.Valid]", tokenClaims.Valid)
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		} else {
			fmt.Println("[parse error]", err)
		}
	}

	return nil, err
}

type GithubAuthStruct struct {
	accessToken string `json: "access_token"`
	tokenType   string `json: "type"`
}

func GithubAuth(c *gin.Context) {
	url := "https://github.com/login/oauth/access_token"
	fmt.Println("URL:>", url)
	code := c.Query("code")
	var u = fmt.Sprintf(`%s?client_id=%s&client_secret=%s&code=%s&state=sup`, url,
		CLIENT_ID, CLIENT_SECRET, code)
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

func main() {

	client := gql.NewClient(
		HASURA_GRAPHQL_ADDRESS,
		gql.Header("x-hasura-admin-secret", HASURA_GRAPHQL_ADMIN_SECRET),
	)
	SimpleQuery(client)
	GetUser(client)

	router := gin.Default()
	router.GET("/auth", Cors(), GithubAuth)
	router.Run(":3001")
}
