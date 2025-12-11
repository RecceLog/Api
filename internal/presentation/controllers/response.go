package controllers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(ctx *gin.Context, status int, msg string, err error) {
	slog.Error(msg, "error", err.Error())
	ctx.JSON(status, gin.H{
		"error": msg,
	})
}
