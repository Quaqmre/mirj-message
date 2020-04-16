package user

import (
	"math"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Quaqmre/mırjmessage/logger"
	"github.com/Quaqmre/mırjmessage/mock"
)

var mockedlogger logger.Service = mock.NewMockedLogger()

var u *user = newUserService(mockedlogger)

func TestNewUser(t *testing.T) {
	ex, err := u.NewUser("ali", "arat")
	if err != nil {
		t.Error("expected nil error but returned:", err)

	}
	if ex.UniqId != 1 {
		t.Error("expected uniqname ali1 but returned:", ex.UniqId)
	}
}
func TestMakeUniqName_with_max_int32(t *testing.T) {
	a := int32(math.MaxInt32)
	u.atomicCounter = &a
	_, err := u.NewUser("ali", "arat")
	if err == nil {
		t.Error("expected error but returned:", err)
	}
}
func TestAtomic_Increase_with_multiple_goroutine(t *testing.T) {

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
	if lastNewUser.UniqId != 1001 {
		t.Error("expected count 1001 but returned:", lastNewUser.UniqId)
	}
}

func TestAtomic_Increase_generete_uniq_Id(t *testing.T) {
	count := int32(0)
	loopcount := 10000
	func() {
		var wg sync.WaitGroup
		for i := 1; i < loopcount; i++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				user, _ := u.NewUser("user"+string(j), "pass")
				atomic.AddInt32(&count, user.UniqId)
			}(i)
		}
		wg.Wait()
	}()

	expectedCount := loopcount * (loopcount - 1) / 2

	if count != int32(expectedCount) {
		t.Error("expected total count:", expectedCount, "but returned total count:", count)
	}
}
