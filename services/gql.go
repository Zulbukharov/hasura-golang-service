package services

import (
	"fmt"

	gql "github.com/btubbs/garphunql"
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
