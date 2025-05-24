package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	allPerm    = 0777
	lessPerm   = 0666
	umask      = 0022
	FilePerm   = lessPerm &^ umask
	FolderPerm = allPerm &^ umask
)

var (
	// flags for command line arguments
	SrcF, DstF    *string
	DeleteMissing *bool

	// logger configuration
	logPath = "log/"
)

// ResetFlags is useful for testing to clear previous flag values.
func ResetFlags() {
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
}

func ParseArgv() error {
	SrcF = flag.String("src", "", "Use --src source-folder")
	DstF = flag.String("dst", "", "Use --dst destination-folder")
	DeleteMissing = flag.Bool("delete-missing", false, "Use --delete-missing true | false")
	flag.Parse()
	if *SrcF == "" || *DstF == "" {
		return fmt.Errorf("both --src and --dst options are mandatory")
	}
	return nil
}

func InitLogger() {
	info, err := os.Stat(logPath)
	if err == nil && !info.IsDir() || os.IsNotExist(err) {
		if err = os.MkdirAll(logPath, FolderPerm); err != nil {
			logrus.Warnf("Failed to create log folder %s: %v", logPath, err)
		}
	}

	logFileName := time.Now().Format("2006-01-02_15-04-05") + ".log"
	logFile, err := os.OpenFile(logPath+logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, FilePerm)
	if err == nil {
		logrus.SetOutput(logFile)
	} else {
		logrus.Warnf("Log information will be printed to stdout")
		logrus.Errorf("Could not open log file %s: %v", logFileName, err)
	}
}
