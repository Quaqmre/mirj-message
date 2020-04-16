package user

import (
	"errors"
	"math"
	"sync"
	"sync/atomic"

	"github.com/Quaqmre/mırjmessage/logger"
)

type Service interface {
}

var ErrorInvalidContext = errors.New("username or password cannot be null")

type User struct {
	Name     string
	UniqId   int32
	Password string
}

type user struct {
	dict          map[string]User
	mutex         sync.RWMutex
	atomicCounter *int32
	logger        logger.Service
}

//NewUserService return new Service
func NewUserService(logger logger.Service) Service {
	return newUserService(logger)
}

//NewUserService return new Service
func newUserService(logger logger.Service) *user {
	var t int32 = 0
	return &user{
		dict:          make(map[string]User),
		atomicCounter: &t,
		logger:        logger,
	}
}

func (u *user) NewUser(name, password string) (*User, error) {
	if name == "" || password == "" {
		return nil, ErrorInvalidContext
	}
	newUser := &User{Name: name, Password: password}
	newUser, err := u.makeUniqName(newUser)
	if err != nil {
		return nil, err
	}
	u.logger.Info("cmp", "user", "method", "newuser", "name", newUser.Name)
	return newUser, nil
}

func (u *user) makeUniqName(user *User) (ru *User, e error) {

	if math.MaxInt32 <= *u.atomicCounter {
		defer func() {
			u.logger.Fatal("cmp", "user", "method", "makeuniqname", "err", e)
		}()

		return nil, errors.New("cant accept new user any more")
	}

	nUser := &User{}
	i := atomic.AddInt32(u.atomicCounter, 1)

	nUser.Name = user.Name
	nUser.Password = user.Password
	nUser.UniqId = i

	u.logger.Info("cmp", "user", "method", "makeuniqname", "err", nil)

	return nUser, nil
}
