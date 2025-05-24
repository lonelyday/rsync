package main

import (
	"github.com/lonelyday/rsync/config"
	"github.com/lonelyday/rsync/rc"

	"github.com/sirupsen/logrus"
)

func main() {
	config.InitLogger()
	if err := config.ParseArgv(); err != nil {
		logrus.Fatal(err.Error())
	}

	err := rc.Sync()
	if err != nil {
		logrus.Fatal(err.Error())
	}
}
