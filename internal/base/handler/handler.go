package handler

import (
	"net/http"
	"wstester/pkg/log"

	"github.com/gin-gonic/gin"
)

func BindAndCheck(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBind(req); err != nil {
		log.Errorf("http_handler bind and check error: %s", err.Error())
		HandleResponse(ctx, err, nil)
		return true
	}

	return false
}

func HandleResponse(ctx *gin.Context, err error, data interface{}) {
	if err == nil {
		ctx.JSON(http.StatusOK, data)
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"reason": err.Error(),
	})
	return
}
