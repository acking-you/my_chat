package models

import (
	"go_http/constants"
	"gorm.io/gorm"
)

type UserInfo struct {
	Id       int64
	UserId   int64
	UserName string
	Name     string
}

type User struct {
	Id             int64
	Username       string `gorm:"unique"`
	Password       string
	UserInfo       *UserInfo //用户信息一对一
	FriendRequests []*User   `gorm:"many2many:user_requests"`  // 好友请求多对多
	Friends        []*User   `gorm:"many2many:user_relations"` //好友列表多对多
}

type UserDAO struct {
}

var userDao = UserDAO{}

func NewUserDAO() UserDAO {
	return userDao
}

func (UserDAO) IsUserExist(userId int64) bool {
	var count int64
	DB.Model(&User{}).Where("id = ?", userId).Count(&count)
	return count != 0
}

func (u UserDAO) GetUserIdByUserNamePassword(username string, password *string) int64 {
	var user User
	if password == nil {
		DB.Select("id").Where("username = ?", username).Find(&user)
	} else {
		DB.Select("id").Where("username = ? and password = ?", username, *password).Find(&user)
	}
	return user.Id
}

func (u UserDAO) GetUserInfoByName(username string) (*UserInfo, error) {
	user := &UserInfo{}

	err := DB.Where("name=?", username).First(user).Error

	if user.Id == 0 {
		return nil, err
	}

	return user, nil
}

// GetUserInfoListByType 逻辑相似整合到一起
func (u UserDAO) GetUserInfoListByType(userId int64, listType int) ([]*UserInfo, error) {
	var list []*UserInfo
	var err error
	switch listType {
	case constants.KFriendRequests:
		err = DB.Raw(`SELECT ui.* FROM user_requests ur INNER JOIN user_infos ui 
            WHERE ur.user_id = ? AND ur.friend_request_id = ui.user_id`, userId).Scan(&list).Error
	case constants.KFriends:
		err = DB.Raw(`SELECT ui.* FROM user_relations ur INNER JOIN user_infos ui 
            WHERE ur.user_id = ? AND ur.friend_id = ui.user_id`, userId).Scan(&list).Error
	}
	return list, err
}

func (u UserDAO) AddUser(username string, password string) (*User, error) {

	user := &User{
		Username: username,
		Password: password,
		UserInfo: &UserInfo{UserName: username, Name: username}, //默认userinfo的名字为用户名
	}

	err := DB.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AddFriend 添加好友的所有动作合集（包括添加FriendRequestList和FriendList
func (u UserDAO) AddFriend(userId int64, userToName string, actionType int) error {
	var err error

	//注意下面两类操作中userId的添加顺序不同
	switch actionType {
	case constants.KFriendRequests:
		err = DB.Exec(`INSERT IGNORE INTO user_requests(user_id,friend_request_id)  
    SELECT u2.id , u1.id FROM users u1,users u2 WHERE u1.id = ? AND u2.username = ?`,
			userId, userToName).Error
	case constants.KFriends:
		//执行事务：需要插入两次好友列表，请求者和同意请求的好友列表
		err = DB.Transaction(func(tx *gorm.DB) error {
			err := tx.Exec(`INSERT IGNORE INTO user_relations(user_id,friend_id)  
		    SELECT u1.id , u2.id FROM users u1,users u2 WHERE u1.id = ? AND u2.username = ?`,
				userId, userToName).Error
			if err != nil {
				return err
			}
			err = tx.Exec(`INSERT IGNORE INTO user_relations(user_id,friend_id)  
		    SELECT u2.id , u1.id FROM users u1,users u2 WHERE u1.id = ? AND u2.username = ?`,
				userId, userToName).Error
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (u UserDAO) DeleteFriend(userId int64, username string, actionType int) error {
	var err error
	switch actionType {
	case constants.KFriendRequests:
		err = DB.Exec(`DELETE ur FROM  user_requests ur INNER JOIN users us 
WHERE us.username = ? AND us.id = ur.friend_request_id AND ur.user_id = ?`, username, userId).Error
	case constants.KFriends:
		err = DB.Exec(`DELETE ur FROM  user_relations ur INNER JOIN users us 
WHERE us.username = ? AND us.id = ur.friend_id AND ur.user_id = ?`, username, userId).Error
	}
	return err
}
