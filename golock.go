package golock

import (
	"os"
	"sync"
	"syscall"
	"time"
)

//Alarm will raise a SIGALRM if not stopped in time
type Alarm struct {
	t       time.Duration
	abort   chan struct{}
	p       *os.Process
	started *sync.WaitGroup
}

func NewAlarm(t time.Duration) (a Alarm, err error) {
	a.p, err = getMyProcess()
	if err != nil {
		return a, err
	}
	a.abort = make(chan struct{}, 1)
	a.started = &sync.WaitGroup{}
	return a, nil
}

//GetMyProcess through os.Getpid
func getMyProcess() (p *os.Process, err error) {
	p, err = os.FindProcess(os.Getpid())
	if err != nil {
		return nil, err
	}
	return p, nil
}

//AlarmAfter a duration or abort on a channel input
func (a Alarm) alarmAfter() {
	timer := time.NewTimer(a.t)
	defer timer.Stop()
	a.started.Done()
	for {
		select {
		case <-timer.C:
			a.p.Signal(syscall.SIGALRM)
		case <-a.abort:
			return
		}
	}
	return
}

func (a Alarm) Start() {
	a.started.Add(1)
	go a.alarmAfter()
	a.started.Wait()
}

func (a Alarm) Stop() {
	a.abort <- struct{}{}
}

func WriteLockWithTimeout(fd *os.File, t time.Duration) (err error) {
	a, err := NewAlarm(t)
	if err != nil {
		return err
	}
	flock := syscall.Flock_t{
		Type: syscall.F_WRLCK,
	}
	a.Start()
	err = syscall.FcntlFlock(fd.Fd(), syscall.F_SETLKW, &flock)
	a.Stop()
	return err

}
func WriteUnLockWithTimeout(fd *os.File, t time.Duration) (err error) {
	a, err := NewAlarm(t)
	if err != nil {
		return err
	}
	flock := syscall.Flock_t{
		Type: syscall.F_UNLCK,
	}
	a.Start()
	err = syscall.FcntlFlock(fd.Fd(), syscall.F_SETLKW, &flock)
	a.Stop()
	return err

}
