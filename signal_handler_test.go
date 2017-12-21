package signalhandler

import (
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestSHandle_RegisterNoSignal(t *testing.T) {

	detect := false

	sH := NewNotifyHandler()

	err := sH.Register(func(sig os.Signal) {
		detect = true
		assert.Equal(t, sig, syscall.SIGUSR1)
	})

	assert.NotNil(t, err)

}
func TestSHandle_Register(t *testing.T) {

	detect := false

	sH := NewNotifyHandler()

	sH.Register(func(sig os.Signal) {
		detect = true
		assert.Equal(t, sig, syscall.SIGUSR1)
	}, syscall.SIGUSR1)

	pid := os.Getpid()
	process, _ := os.FindProcess(pid)
	process.Signal(syscall.SIGUSR1)
	time.Sleep(100 * time.Millisecond)
	assert.True(t, detect)

}

func TestSHandle_RegisterNotNotify(t *testing.T) {

	var fnc NotifyFnc = func(sig os.Signal) {
		assert.Fail(t, "Cannot go here")
	}

	sH := NewNotifyHandler()

	sH.Register(fnc, syscall.SIGKILL)
	pid := os.Getpid()

	process, _ := os.FindProcess(pid)
	process.Signal(syscall.SIGUSR1)
	time.Sleep(10 * time.Millisecond)
	assert.True(t, true)
}

func TestSHandle_UnRegister(t *testing.T) {
	registered := true

	var doNothing NotifyFnc = func(sig os.Signal) {
		if registered {

			assert.True(t, registered)
			return
		}
		assert.Fail(t, "Cannot go here")
	}

	sH := NewNotifyHandler()
	sH.Register(doNothing, syscall.SIGUSR1)

	pid := os.Getpid()

	process, _ := os.FindProcess(pid)
	process.Signal(syscall.SIGUSR1)
	time.Sleep(10 * time.Millisecond)
	assert.True(t, registered)

	//Now unregister it
	registered = false
	sH.UnRegister(syscall.SIGUSR1)
	process.Signal(syscall.SIGUSR1)
	time.Sleep(10 * time.Millisecond)
	assert.False(t, registered)

}

func TestSHandle_UnRegisterAll(t *testing.T) {
	registered := true

	var doNothing NotifyFnc = func(sig os.Signal) {
		if registered {

			assert.True(t, registered)
			return
		}
		assert.Fail(t, "Cannot go here")
	}

	sH := NewNotifyHandler()
	sH.Register(doNothing, syscall.SIGUSR1, syscall.SIGUSR2)

	pid := os.Getpid()

	process, _ := os.FindProcess(pid)
	process.Signal(syscall.SIGUSR1)
	time.Sleep(10 * time.Millisecond)
	assert.True(t, registered)

	//Now unregister it
	registered = false
	sH.UnRegister()
	process.Signal(syscall.SIGUSR1)
	time.Sleep(10 * time.Millisecond)
	assert.False(t, registered)

}
