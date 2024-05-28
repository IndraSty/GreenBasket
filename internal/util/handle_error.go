package util

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HandleError(ctx *gin.Context, err error, status int, message string) {
	log.Println(err.Error())
	ctx.JSON(status, gin.H{"error": message})
}
