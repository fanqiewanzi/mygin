package main

import (
	"mygin/gin"
	"net/http"
)

// FIXME:别这么写测试，去看看接口测试怎么做，顺便看看单元测试。
func main() {
	r := gin.New()
	r.Use(gin.Logger())
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "das", "<h1>gin</h1>")
		})

		v1.GET("/hello", func(c *gin.Context) {
			c.String(http.StatusOK, "hello %s\n", c.Query("name"))
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gin.Context) {
			c.String(http.StatusOK, "hello %s\n", c.Param("name"))
		})
		v2.POST("/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":8080")
}
