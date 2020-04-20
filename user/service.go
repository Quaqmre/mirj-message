package user

import (
	"encoding/json"
	"errors"
	"math"
	"sync"
	"sync/atomic"

	"github.com/Quaqmre/mırjmessage/logger"
)

// // Service OutDoor tho user service
// type Service interface {
// 	NewUser(name, password string) (*User, error)
// }

// ErrorInvalidContext using  when user or password is nil
var ErrorInvalidContext = errors.New("username or password cannot be nil")

type Service interface {
	NewUser(name, password string) (*User, error)
	Get(userID int32) *User
	Marshall(source []byte) (*User, error)
}

// User hold information most tiny way
type User struct {
	Name     string
	UniqID   int32
	Password string
}

type UserService struct {
	Dict          map[int32]*User
	mutex         sync.RWMutex
	atomicCounter *int32
	logger        logger.Service
}

//NewUserService return new Service
func NewUserService(logger logger.Service) Service {
	return newUserService(logger)
}

//NewUserService return new Service
func newUserService(logger logger.Service) *UserService {
	var t int32 = 0
	return &UserService{
		Dict:          make(map[int32]*User),
		atomicCounter: &t,
		logger:        logger,
	}
}

func (u *UserService) NewUser(name, password string) (*User, error) {
	if name == "" || password == "" {
		return nil, ErrorInvalidContext
	}
	newUser := &User{Name: name, Password: password}
	newUser, err := u.makeUniqName(newUser)
	if err != nil {
		return nil, err
	}
	u.mutex.Lock()
	u.Dict[newUser.UniqID] = newUser
	u.mutex.Unlock()
	u.logger.Info("cmp", "user", "method", "newuser", "name", newUser.Name)
	return newUser, nil
}

func (u *UserService) makeUniqName(user *User) (ru *User, e error) {

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
	nUser.UniqID = i

	u.logger.Info("cmp", "user", "method", "makeuniqname", "err", nil)

	return nUser, nil
}

// Get user in map
func (u *UserService) Get(userID int32) *User {
	return u.Dict[userID]
}

// TODO : give a change to admin parse json or protobuff

// Marshall bytes to user
func (u *UserService) Marshall(source []byte) (*User, error) {
	user := &User{}
	err := json.Unmarshal(source, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// TODO implement in memory store
