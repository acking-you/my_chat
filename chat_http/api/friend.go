package api

import (
	"chat_socket/core"
	"github.com/gin-gonic/gin"
	"go_http/models"
	"go_http/pkg/utils"
	"go_http/service"
	"strconv"
)

func AddFriend(c *gin.Context) {
	userId := utils.GetId(c)
	if userId == -1 {
		return
	}
	action, err := strconv.ParseInt(c.Query("action"), 10, 32)
	if err != nil {
		utils.SendError(c, "请求参数缺失")
		return
	}
	userToName := c.Query("username")
	friendService := service.FriendService{UserId: userId, UserToName: &userToName}
	err = friendService.DoAddFriend(int(action))
	if err != nil {
		utils.SendError(c, err.Error())
		return
	}
	utils.SendOk[*int](c, nil)
}

func DeleteFriend(c *gin.Context) {
	userId := utils.GetId(c)
	if userId == -1 {
		return
	}
	action, err := strconv.ParseInt(c.Query("action"), 10, 32)
	if err != nil {
		utils.SendError(c, "请求参数缺失")
		return
	}
	userToName := c.Query("username")
	friendService := service.FriendService{UserId: userId, UserToName: &userToName}
	err = friendService.DoDeleteFriend(int(action))
	if err != nil {
		utils.SendError(c, err.Error())
		return
	}
	utils.SendOk[*int](c, nil)
}

func QueryList(c *gin.Context) {
	userId := utils.GetId(c)
	if userId == -1 {
		return
	}
	action, err := strconv.ParseInt(c.Query("action"), 10, 32)
	if err != nil {
		utils.SendError(c, "请求参数缺失")
		return
	}
	friendService := service.FriendService{UserId: userId}
	response, err := friendService.DoQueryList(int(action))
	if err != nil {
		utils.SendError(c, err.Error())
	}

	utils.SendOk[*service.FriendResponse](c, response)
}

func GetUserInfoByName(c *gin.Context) {
	username := c.Query("username")
	if len(username) == 0 {
		utils.SendError(c, "参数请求缺失")
		return
	}
	userInfo, err := models.NewUserDAO().GetUserInfoByName(username)
	if err != nil {
		utils.SendError(c, "用户名不存在")
		return
	}
	utils.SendOk[*models.UserInfo](c, userInfo)
}

func CheckUserState(c *gin.Context) {
	username := c.Query("username")
	if len(username) == 0 {
		utils.SendError(c, "参数请求缺失")
		return
	}
	userId := models.NewUserDAO().GetUserIdByUserNamePassword(username, nil)
	if userId == 0 {
		utils.SendError(c, "用户名不存在")
		return
	}
	if core.Manager.GetClient(userId) == nil {
		utils.SendUserStatus(c, utils.Offline) //不在线
		return
	}
	utils.SendUserStatus(c, utils.Online) //在线
}
