package main

import (
	"os"

	"github.com/Zulbukharov/hasura-golang-service/db"
	"github.com/Zulbukharov/hasura-golang-service/routes"

	"github.com/gin-gonic/gin"
)

var CLIENT_ID = os.Getenv("CLIENT_ID")
var CLIENT_SECRET = os.Getenv("CLIENT_SECRET")

func main() {

	db.Init()
	router := gin.Default()
	routes.SetupRoutes(router)
	router.Run(":3001")
}
