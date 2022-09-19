package service

import (
	"errors"
	"fmt"
	"go_http/models"
	"go_http/pkg/middleware"
	chatlog "logger/log"
)

const (
	NameMax       = 20
	NameAtLeast   = 5
	LoginState    = 0
	RegisterState = 1
)

var (
	ErrCheckNumLen      = fmt.Errorf("用户名长度应该在%d - %d之间", NameAtLeast, NameMax)
	ErrCheckNumEmpty    = errors.New("用户名或密码不能为空")
	ErrCheckNumNotExist = errors.New("用户名或密码不存在")
	ErrToken            = errors.New("token 颁发失败")
	ErrRegisterRepeat   = errors.New("注册失败,用户名已经存在")
)

type LoginResponse struct {
	UserId int64
	Token  string
}

type LoginService struct {
	state    int
	Name     string
	Password string

	Token  string
	UserId int64
}

func NewLoginService(name string, password string) *LoginService {
	return &LoginService{Name: name, Password: password}
}

func (l *LoginService) DoLogin() (*LoginResponse, error) {
	var err error
	l.state = LoginState
	if err = l.checkNum(); err != nil {
		return nil, err
	}

	if err := l.getToken(); err != nil {
		return nil, err
	}
	return &LoginResponse{
		UserId: l.UserId,
		Token:  l.Token,
	}, nil
}

func (l *LoginService) DoRegister() (*LoginResponse, error) {
	l.state = RegisterState
	var err error
	if err = l.checkNum(); err != nil {
		return nil, err
	}
	if err = l.register(); err != nil {
		return nil, err
	}
	return &LoginResponse{
		UserId: l.UserId,
		Token:  l.Token,
	}, nil
}

func (l *LoginService) checkNum() error {
	if len(l.Name) == 0 || len(l.Password) == 0 {
		return ErrCheckNumEmpty
	}
	if len(l.Name) < NameAtLeast || len(l.Name) > NameMax {
		return ErrCheckNumLen
	}
	if l.state == RegisterState {
		return nil
	}

	userId := models.NewUserDAO().GetUserIdByUserNamePassword(l.Name, &l.Password)
	if userId == 0 {
		return ErrCheckNumNotExist
	}
	l.UserId = userId
	return nil
}

func (l *LoginService) getToken() error {
	token, err := middleware.ReleaseToken(l.UserId)
	if err != nil {
		return ErrToken
	}
	l.Token = token
	return nil
}

func (l *LoginService) register() error {
	user, err := models.NewUserDAO().AddUser(l.Name, l.Password)
	if err != nil {
		chatlog.Lg().Errorln(err)
		return ErrRegisterRepeat
	}
	l.UserId = user.Id
	//颁发token
	err = l.getToken()
	if err != nil {
		return err
	}
	return nil
}
