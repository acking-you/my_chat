package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	Ok      = 0
	Invalid = 1
	Online  = 2
	Offline = 3
)

type CommonResponse struct {
	StatusCode int32
	StatusMsg  string
}

type Response[T any] struct {
	CommonResponse
	ResponseData T
}

func GetId(ctx *gin.Context) int64 {
	raw, _ := ctx.Get("user_id")
	userId, ok := raw.(int64)
	if !ok {
		SendError(ctx, "id获取出错")
		return -1
	}
	return userId
}

func SendError(ctx *gin.Context, errorText string) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.JSON(http.StatusOK, Response[*int]{
		CommonResponse{
			StatusCode: Invalid,
			StatusMsg:  errorText,
		},
		nil,
	})
}

func SendOk[T any](ctx *gin.Context, data T) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.JSON(http.StatusOK, Response[T]{
		CommonResponse{
			StatusCode: Ok,
			StatusMsg:  "请求成功",
		},
		data,
	})
}

func SendUserStatus(ctx *gin.Context, userState int32) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.JSON(http.StatusOK, &CommonResponse{
		StatusCode: userState,
		StatusMsg:  "请求成功",
	},
	)
}
