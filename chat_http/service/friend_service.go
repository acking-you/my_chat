package service

import (
	"errors"
	"go_http/constants"
	"go_http/models"
)

//拉取和添加列表的类型

type FriendResponse struct {
	UserInfos []*models.UserInfo `json:"user_infos,omitempty"`
}

type FriendService struct {
	UserId     int64
	UserToName *string

	response *FriendResponse
}

func (f *FriendService) DoAddFriend(actionType int) error {
	var err error
	if err = f.checkNum(actionType); err != nil {
		return err
	}
	if err = f.addFriend(actionType); err != nil {
		return err
	}
	return nil
}

func (f *FriendService) DoDeleteFriend(actionType int) error {
	var err error
	if err = f.checkNum(actionType); err != nil {
		return err
	}
	if err = f.deleteFriend(actionType); err != nil {
		return err
	}
	return nil
}

func (f *FriendService) DoQueryList(listType int) (*FriendResponse, error) {
	var err error
	if err = f.checkNum(listType); err != nil {
		return nil, err
	}
	//获取好友请求列表
	if err := f.getList(listType); err != nil {
		return nil, errors.New("用户不存在")
	}
	return f.response, nil
}

func (f *FriendService) getList(listType int) error {

	var list []*models.UserInfo
	var err error

	//具体的获取列表的逻辑
	switch listType {
	case constants.KFriendRequests:
		list, err = models.NewUserDAO().GetUserInfoListByType(f.UserId, constants.KFriendRequests)
	case constants.KFriends:
		list, err = models.NewUserDAO().GetUserInfoListByType(f.UserId, constants.KFriends)
	}

	if err != nil {
		return err
	}
	f.response = &FriendResponse{UserInfos: list}
	return nil
}

func (f *FriendService) checkNum(action int) error {
	if action != constants.KFriendRequests && action != constants.KFriends {
		return errors.New("未定义行为")
	}
	if f.UserToName != nil && models.NewUserDAO().GetUserIdByUserNamePassword(*f.UserToName, nil) == 0 {
		return errors.New("需要操作的username不存在")
	}
	if !models.NewUserDAO().IsUserExist(f.UserId) {
		return errors.New("userId不存在")
	}
	return nil
}

func (f *FriendService) addFriend(actionType int) error {
	if f.UserToName == nil {
		return errors.New("userToName nil")
	}
	err := models.NewUserDAO().AddFriend(f.UserId, *f.UserToName, actionType)
	if err != nil {
		return err
	}
	return nil
}

func (f *FriendService) deleteFriend(actionType int) error {
	if f.UserToName == nil {
		return errors.New("userToName nil")
	}
	err := models.NewUserDAO().DeleteFriend(f.UserId, *f.UserToName, actionType)
	if err != nil {
		return err
	}
	return nil
}
