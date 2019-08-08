package main

import (
	"os"

	"github.com/Zulbukharov/hasura-golang-service/handlers"
	"github.com/Zulbukharov/hasura-golang-service/services"
	gql "github.com/btubbs/garphunql"

	"github.com/gin-gonic/gin"
)

var HASURA_GRAPHQL_ADDRESS = os.Getenv("HASURA_GRAPHQL_ADDRESS")
var HASURA_GRAPHQL_ADMIN_SECRET = os.Getenv("HASURA_GRAPHQL_ADMIN_SECRET")
var CLIENT_ID = os.Getenv("CLIENT_ID")
var CLIENT_SECRET = os.Getenv("CLIENT_SECRET")

func main() {

	client := gql.NewClient(
		HASURA_GRAPHQL_ADDRESS,
		gql.Header("x-hasura-admin-secret", HASURA_GRAPHQL_ADMIN_SECRET),
	)
	services.SimpleQuery(client)
	services.GetUser(client)

	router := gin.Default()
	router.GET("/auth", handlers.Cors(), handlers.GithubAuth)
	router.Run(":3001")
}
