package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		pwd, ok := db[user]
		_ = pwd

		if ok {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "failed", "message": "User not found"})
		}
	})

	// Post mqtt auth
	r.GET("/mqtt/auth", func(c *gin.Context) {
		username := c.Query("username")
		password := c.Query("password")

		pwd, ok := db[username]

		if ok && (strings.Compare(pwd, password) == 0) {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed"})
		}
	})

	// Post mqtt auth
	r.POST("/mqtt/auth", func(c *gin.Context) {

		if err := c.Request.ParseForm(); err != nil {
			fmt.Println("resolve param error:", err)
		}

		username := c.Request.FormValue("username")
		password := c.Request.FormValue("password")

		pwd, ok := db[username]

		if ok && (strings.Compare(pwd, password) == 0) {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed"})
		}

	})

	// Get mqtt superuser
	r.GET("/mqtt/superuser", func(c *gin.Context) {
		username := c.Query("username")
		_ = username
		clientid := c.Query("clientid")
		_ = clientid

		c.JSON(http.StatusBadRequest, gin.H{"status": "failed"})
	})

	// Get mqtt superuser
	r.POST("/mqtt/superuser", func(c *gin.Context) {

		if err := c.Request.ParseForm(); err != nil {
			fmt.Println("resolve param error:", err)
		}

		username := c.Request.FormValue("username")
		_ = username
		clientid := c.Request.FormValue("clientid")
		_ = clientid

		c.JSON(http.StatusBadRequest, gin.H{"status": "failed"})
	})

	// Get mqtt superuser
	r.GET("/mqtt/acl", func(c *gin.Context) {
		access := c.Query("access")
		_ = access
		username := c.Query("username")
		_ = username
		clientid := c.Query("clientid")
		_ = clientid
		ipaddr := c.Query("ipaddr")
		_ = ipaddr
		topic := c.Query("topic")
		_ = topic

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	// 	"admin":  "password",
	// }))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "password", // user:admin password:password
	}))

	/* example curl for /admin with basicauth header
	   YWRtaW46cGFzc3dvcmQK is base64("admin:password")
	   Zm9vOmJhcgo= is base64("foo:bar")

		curl -X POST \
	    -u "admin:password" \
	  	-H 'content-type: application/json' \
	  	-d '{"user":"foo", "pwd":"bar"}' \
		http://localhost:8080/admin
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		_ = user

		// Parse JSON
		var json struct {
			Pwd  string `json:"pwd" binding:"required"`
			User string `json:"user" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[json.User] = json.Pwd
			// db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
