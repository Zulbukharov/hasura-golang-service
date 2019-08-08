package main

import (
	"os"

	"github.com/Zulbukharov/hasura-golang-service/db"

	"github.com/Zulbukharov/hasura-golang-service/handlers"
	"github.com/gin-gonic/gin"
)

var CLIENT_ID = os.Getenv("CLIENT_ID")
var CLIENT_SECRET = os.Getenv("CLIENT_SECRET")

func main() {

	// services.SimpleQuery(client)
	// services.GetUser(client)
	db.Init()
	router := gin.Default()
	router.GET("/auth", handlers.Cors(), handlers.Auth)
	router.Run(":3001")
}
