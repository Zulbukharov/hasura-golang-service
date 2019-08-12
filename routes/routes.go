package routes

import (
	"github.com/Zulbukharov/hasura-golang-service/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/auth", handlers.Cors(), handlers.Auth)
	r.GET("/generateToken", handlers.GenerateToken)
}
