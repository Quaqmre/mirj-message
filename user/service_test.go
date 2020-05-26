package user

import (
	"encoding/json"
	"math"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Quaqmre/mırjmessage/logger"
	"github.com/Quaqmre/mırjmessage/mock"
)

func TestNewUser(t *testing.T) {
	var mockedlogger logger.Service = mock.NewMockedLogger()
	var u *UserService = newUserService(mockedlogger)
	tests := []struct {
		name           string
		input          string
		expectedResult int32
	}{
		{
			name:           "firstuser",
			input:          "user1",
			expectedResult: 1,
		},
		{
			name:           "seconduser",
			input:          "user2",
			expectedResult: 2,
		},
	}
	for _, test := range tests {
		ex, err := u.NewUser(test.input, "arat")
		if err != nil {
			t.Error("expected nil error but returned:", err)
		}
		if ex.UniqID != test.expectedResult {
			t.Error("expected uniqname ali1 but returned:", ex.UniqID)
		}
	}
}

func TestNewUser_return_error(t *testing.T) {
	var mockedlogger logger.Service = mock.NewMockedLogger()
	var u *UserService = newUserService(mockedlogger)

	expected := ErrorInvalidContext
	_, err := u.NewUser("", "")
	if err != expected {
		t.Errorf("expected %s but returned %s", expected.Error(), err.Error())
	}

}

func TestMakeUniqName_with_max_int32(t *testing.T) {
	var mockedlogger logger.Service = mock.NewMockedLogger()
	var u *UserService = newUserService(mockedlogger)

	a := int32(math.MaxInt32)
	u.atomicCounter = &a
	_, err := u.NewUser("ali", "arat")
	if err == nil {
		t.Error("expected error but returned:", err)
	}
}
func TestAtomic_Increase_with_multiple_goroutine(t *testing.T) {
	var mockedlogger logger.Service = mock.NewMockedLogger()
	var u *UserService = newUserService(mockedlogger)

	func() {
		var wg sync.WaitGroup
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				u.NewUser("user"+string(j), "pass")
			}(i)
		}
		wg.Wait()
	}()
	lastNewUser, _ := u.NewUser("test", "deneme")
	if lastNewUser.UniqID != 1001 {
		t.Error("expected count 1001 but returned:", lastNewUser.UniqID)
	}
}

func TestAtomic_Increase_generete_uniq_Id(t *testing.T) {
	var mockedlogger logger.Service = mock.NewMockedLogger()
	var u *UserService = newUserService(mockedlogger)

	count := int32(0)
	loopcount := 10000
	func() {
		var wg sync.WaitGroup
		for i := 1; i < loopcount; i++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				user, _ := u.NewUser("user"+string(j), "pass")
				atomic.AddInt32(&count, user.UniqID)
			}(i)
		}
		wg.Wait()
	}()

	expectedCount := loopcount * (loopcount - 1) / 2

	if count != int32(expectedCount) {
		t.Error("expected total count:", expectedCount, "but returned total count:", count)
	}
}

func TestUserService_Get(t *testing.T) {
	var mockedlogger logger.Service = mock.NewMockedLogger()
	var u *UserService = newUserService(mockedlogger)
	expected, _ := u.NewUser("deneme", "sifre")

	newUser := u.Get(expected.UniqID)
	if !reflect.DeepEqual(expected, newUser) {
		t.Error("genereted user and getting user is not equal")
	}

}

func TestUserService_Marshall(t *testing.T) {
	var mockedlogger logger.Service = mock.NewMockedLogger()
	var u *UserService = newUserService(mockedlogger)

	expected := User{Name: "akif", Password: "123", UniqID: 1}
	testMarshal, _ := json.Marshal(expected)

	parsedUser, _ := u.Marshall(testMarshal)

	if expected != *parsedUser {
		t.Errorf("expected User string %v,returned %v", expected, parsedUser)
	}
}
func TestUserService_Marshall_return_error(t *testing.T) {
	var mockedlogger logger.Service = mock.NewMockedLogger()
	var u *UserService = newUserService(mockedlogger)

	testMarshal := `{"User":"akif"`

	_, err := u.Marshall([]byte(testMarshal))

	if err == nil {
		t.Error("expected err but returned err=nil")
	}
}
