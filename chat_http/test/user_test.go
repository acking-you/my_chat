package test

import (
	"fmt"
	"go_http/models"
	"testing"
)

func TestAddUser(t *testing.T) {
	for i := 0; i < 100; i++ {
		user, err := models.NewUserDAO().AddUser(fmt.Sprintf("test%d", i), "dfskfjskfsdfsf")
		if err != nil {
			panic(err)
		}
		fmt.Println(*user)
	}

}

func TestGetFriendRequest(t *testing.T) {
	list, err := models.NewUserDAO().GetUserInfoListByType(1, 0)
	if err != nil {
		panic(err)
	}
	for _, info := range list {
		fmt.Println(*info)
	}
}
