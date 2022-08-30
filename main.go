package main

import (
	"fmt"
	"net/http"
	"strconv"
	"tsmedberg/te4-introprojekt/database"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	r.Static("/css", "templates/css")
	r.Static("/images", "templates/images")
	r.GET("/", func(c *gin.Context) {
		posts, err := database.Read()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"posts": posts,
		})
	})
	r.GET("/post/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		post, err := database.ReadOne(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "update.tmpl", gin.H{
			"Post": post,
		})
	})
	r.POST("/post/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var post database.Post
		err = c.ShouldBind(&post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		post.Id = id
		fmt.Println(post)
		fmt.Println(id)
		err = database.Update(post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		/* c.Redirect(http.StatusTemporaryRedirect, "/") */
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})

	})
	r.GET("/delete/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		post, err := database.ReadOne(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "delete.tmpl", gin.H{
			"Post": post,
		})
	})
	r.POST("/delete/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = database.Delete(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
	r.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.tmpl", gin.H{})
	})
	r.POST("/create", func(c *gin.Context) {
		var post database.Post
		err := c.ShouldBind(&post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = database.Create(post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	r.Run(":3000")
}
