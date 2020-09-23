package main

import (
	"Crewl/Handler"
	"Crewl/Model"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println(context.Request.URL.Path)

		if context.Request.URL.Path == "/Login" {
			if Handler.Start(context) == true {
				context.Redirect(http.StatusMovedPermanently, "/Normal")
			} else {
				context.Next()
			}
		} else  {
			fmt.Println("second ..")
			if Handler.Start(context) == true {
				context.Next()
			} else {
				fmt.Println("last..")
				context.Redirect(http.StatusMovedPermanently, "/Login")
			}
		}

	}
}


func main() {
	// 初始化数据库连接
	fmt.Println("Starting...")
	Model.Init()
	router := gin.Default()
	router.LoadHTMLGlob("Templates/*")

	router.Use(Auth())
	router.GET("/Login", Handler.Login)
	router.POST("/LoginAction", Handler.LoginAction)
	router.GET("/Normal", Handler.ShowNormal)
	router.GET("/Exception", Handler.ShowExcept)

	router.Run(":9999")
}
