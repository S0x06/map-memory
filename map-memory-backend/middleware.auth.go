package main

import (
	"map-memory-backend/config"
	"map-memory-backend/entity"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func IownMemoryFromPost(ctx *gin.Context) {
	session := sessions.Default(ctx)
	uid := session.Get("uid").(int64)

	mid := ctx.PostForm("id")

	var mem entity.Memory

	config.RDB_CONN.Table("mp_memory").Where("id = ?", mid).Find(&mem)

	if mem.ID != 0 && mem.UserID == uid {
		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}

func AuthIsOwnerOfMemoryById() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid := session.Get("uid").(int64)

		mid := c.Param("id")

		var mem entity.Memory

		config.RDB_CONN.Table("mp_memory").Where("id = ?", mid).Find(&mem)

		if mem.ID != 0 && mem.UserID == uid {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// get memory by id
	}
}
func AuthNotRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			// You'd normally redirect to login page
			c.Next()
		} else {
			// Continue down the chain to handler etc
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			// You'd normally redirect to login page
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		} else {
			// Continue down the chain to handler etc
			c.Next()
		}
	}
}
