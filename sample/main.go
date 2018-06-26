package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	user := map[string]interface{}{
		"name": "jeremaihloo",
		"age":  23,
	}
	users := []map[string]interface{}{
		user,
	}
	var newUser map[string]interface{}
	r.GET("/api/users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"items": users,
		})
	})
	r.POST("/api/users", func(c *gin.Context) {
		c.BindJSON(&newUser)
		users = append(users, newUser)
		c.JSON(200, newUser)
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
