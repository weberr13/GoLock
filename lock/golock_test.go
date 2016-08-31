package golock

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAlarm(t *testing.T) {

	Convey("NewAlarm", t, func() {
		a, err := NewAlarm(1 * time.Millisecond)
		So(err, ShouldBeNil)
		So(a.p, ShouldNotBeNil)
		So(a.abort, ShouldNotBeNil)
		So(a.t, ShouldEqual, 1*time.Millisecond)
	})
	Convey("NewAlarm timesout", t, func() {
		a, err := NewAlarm(1 * time.Millisecond)
		So(err, ShouldBeNil)
		So(a.p, ShouldNotBeNil)
		So(a.abort, ShouldNotBeNil)
		So(a.t, ShouldEqual, 1*time.Millisecond)
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGALRM)
		t2 := time.NewTimer(2 * time.Millisecond)
		defer t2.Stop()
		a.Start()
		select {
		case s := <-c:
			So(s, ShouldEqual, syscall.SIGALRM)
		case <-t2.C:
			t.FailNow()
		}

	})
	Convey("NewAlarm stopped in time", t, func() {
		a, err := NewAlarm(2 * time.Millisecond)
		So(err, ShouldBeNil)
		So(a.p, ShouldNotBeNil)
		So(a.abort, ShouldNotBeNil)
		So(a.t, ShouldEqual, 2*time.Millisecond)
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGALRM)
		t2 := time.NewTimer(1 * time.Millisecond)
		defer t2.Stop()
		a.Start()
		select {
		case <-c:
			t.FailNow()
		case <-t2.C:
			a.Stop()
		}

		t2.Reset(2 * time.Millisecond)
		select {
		case <-c:
			t.FailNow()
		case <-t2.C:
			a.Stop()
		}

	})
}
