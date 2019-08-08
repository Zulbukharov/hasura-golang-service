package db

import (
	"os"

	gql "github.com/btubbs/garphunql"
)

var HASURA_GRAPHQL_ADDRESS = os.Getenv("HASURA_GRAPHQL_ADDRESS")
var HASURA_GRAPHQL_ADMIN_SECRET = os.Getenv("HASURA_GRAPHQL_ADMIN_SECRET")

var client *gql.Client

func Init() {
	client = gql.NewClient(
		HASURA_GRAPHQL_ADDRESS,
		gql.Header("x-hasura-admin-secret", HASURA_GRAPHQL_ADMIN_SECRET),
	)
}

func GetClient() *gql.Client {
	return client
}
