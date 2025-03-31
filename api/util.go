package api

import (
	"man/util"

	"github.com/gin-gonic/gin"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func getUserID(ctx *gin.Context) int32 {
	return int32(ctx.MustGet(authorizationPayloadKey).(*util.Payload).UserID)
}

type Pagination struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=50"`
}
