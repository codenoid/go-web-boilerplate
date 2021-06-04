package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/codenoid/go-web-boilerplate/web-base/structs/user"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func loginHTML(c *gin.Context) {
	// token := csrf.GetToken(c)
	c.HTML(http.StatusOK, "auth_login.html", gin.H{
		// "token": token,
		"app_name": appName,
	})
}

func loginVerify(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user user.User
	err := mainDB.Collection("users").Find(ctx, bson.M{"username": username}).One(&user)
	if err != nil {
		c.Redirect(302, "/auth/login?error=true&msg=invalid username or password")
		return
	}

	hashedPassword := user.Password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		c.Redirect(302, "/auth/login?error=True&msg=invalid username or password")
		return
	}

	sessionDB.Set("auth:public_id_of:"+c.GetString("guid"), user.PublicID, 7*24*time.Hour)
	if user.SingleSession {
		sessionDB.Set("auth:guid_of:"+user.PublicID, c.GetString("guid"), 7*24*time.Hour)
	}

	c.Redirect(302, "/")
}

func logout(c *gin.Context) {
	sessionDB.Del("auth:public_id_of:" + c.GetString("guid"))
	c.Redirect(http.StatusFound, "/auth/login?msg=Logged Out Successfully")
}

func mainMiddleware(c *gin.Context) {

	guid, err := c.Cookie("guid")
	if err != nil {
		newGUID := uuid.New().String()
		if newGUID == "" {
			c.String(500, "Internal Server Error")
			c.Abort()
			return
		}

		c.SetCookie("guid", newGUID, 86400*6, "/", "", false, true)
		c.Redirect(302, "/auth/login?error=true&msg=invalid session 1")
		c.Abort()
		return
	}

	c.Set("guid", guid)

	paths := strings.Split(c.Request.URL.Path, "/")
	if paths[1] != "auth" {
		publicID, err := sessionDB.Get("auth:public_id_of:" + guid).Result()
		if err != nil {
			c.Redirect(302, "/auth/login?error=true&msg=invalid session 2")
			c.Abort()
			return
		}

		var user user.User

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = mainDB.Collection("users").Find(ctx, bson.M{"public_id": publicID}).One(&user)
		if err != nil {
			c.Redirect(302, "/auth/login?error=true&msg=invalid user, please relogin")
			c.Abort()
			return
		}

		if user.SingleSession {
			authorizedGUID, err := sessionDB.Get("auth:guid_of:" + user.PublicID).Result()
			if err != nil {
				c.Redirect(302, "/auth/logout?error=true&msg=invalid system check, please relogin")
				c.Abort()
				return
			}
			if authorizedGUID != guid {
				c.Redirect(302, "/auth/logout?error=true&msg=unauthorized, please re-login")
				c.Abort()
				return
			}
		}

		c.Set("private_id", user.PrivateID)
		c.Set("username", user.Username)
		c.Set("userData", user)
	}
}
