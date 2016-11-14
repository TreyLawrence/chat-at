package main

import (
	"chatat/controllers"
	"chatat/rest"

	gin "gopkg.in/gin-gonic/gin.v1"
)

func main() {
	engine := gin.Default()
	api := engine.Group("")

	rest.CRUD(api, "/conversations", controllers.Conversations{})
	rest.CRUD(api, "/conversations/:/messages", controllers.Messages{})

	engine.Run()
}
