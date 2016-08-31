package golock

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

//Alarm will raise a SIGALRM if not stopped in time
type Alarm struct {
	t     time.Duration
	abort chan struct{}
	p     *os.Process
}

//NewAlarm creates an alarm for the given duration
func NewAlarm(t time.Duration) (a Alarm, err error) {
	a.p, err = getMyProcess()
	if err != nil {
		return a, err
	}
	a.t = t
	a.abort = make(chan struct{}, 1)
	return a, nil
}

//getMyProcess through os.Getpid
func getMyProcess() (p *os.Process, err error) {
	p, err = os.FindProcess(os.Getpid())
	if err != nil {
		return nil, err
	}
	return p, nil
}

//alarmAfter a duration or abort on a channel input
func (a Alarm) alarmAfter() {
	timer := time.NewTimer(a.t)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			fmt.Println("got timeout, alarm to giveup")
			a.p.Signal(syscall.SIGALRM)
		case <-a.abort:
			return
		}
	}
	return
}

//Start the countdown till an alarm is raised
func (a Alarm) Start() {
	go a.alarmAfter()
}

//Stop the countdown for the alarm
func (a Alarm) Stop() {
	a.abort <- struct{}{}
}

//WriteLockWithTimeout write lock fd but give up after t
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

//WriteUnLockWithTimeout un-lock fd but give up after t
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
