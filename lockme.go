package main

import (
	"fmt"
	"os"
	"time"

	FLAGS "github.com/jessevdk/go-flags"
	"github.com/weberr13/GoLock/lock"
)

type Opts struct {
	Time    string `short:"t" long:"time" description:"time to lock before exit"`
	Timeout string `short:"o" long:"timeout" description:"time to wait before giving up"`
}

var opts Opts
var args []string
var parser = FLAGS.NewParser(&opts, FLAGS.Default)

func main() {
	opts.Time = "1h"
	opts.Timeout = "1s"
	args, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(args) <= 0 {
		fmt.Println("Must specify filename")
		os.Exit(1)
	}
	fd, err := os.OpenFile(args[0], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fd.Close()
	d, err := time.ParseDuration(opts.Time)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	timeout, err := time.ParseDuration(opts.Timeout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = golock.WriteLockWithTimeout(fd, timeout)
	if err != nil {
		fmt.Println("Could not lock file ", args[0])
	}
	timer := time.NewTimer(d)
	select {
	case <-timer.C:
		return
	}

}
