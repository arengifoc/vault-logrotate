package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/go-co-op/gocron"
)

func logInfo(message string) {
	timestamp := time.Now().UTC().Format("Jan 02 15:04:05")
	hostname, _ := os.Hostname()
	log.Printf("%s %s %s", timestamp, hostname, message)
}

func logError(message string) {
	timestamp := time.Now().UTC().Format("Jan 02 15:04:05")
	hostname, _ := os.Hostname()
	log.Printf("%s %s [ERROR] %s", timestamp, hostname, message)
}

func main() {
	// Open log file
	logFile, err := os.OpenFile("/var/log/logrotate/logrotate.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Set log output to file
	log.SetOutput(logFile)
	log.SetFlags(0) // Disable default timestamp since we handle it ourselves

	crond := gocron.NewScheduler(time.UTC)
	var crontab string
	var ok bool
	if crontab, ok = os.LookupEnv("CRONTAB"); !ok {
		crontab = "0 * * * *"
	}
	logInfo(fmt.Sprintf("Starting vault logrotation with \"%s\"", crontab))
	crond.Cron(crontab).Do(run_logrotate)
	crond.StartBlocking()
}

func run_logrotate() {
	cmd := exec.Command("/usr/sbin/logrotate", "--state=/tmp/logrotate.status", "/etc/logrotate.conf")
	logInfo("Starting logrotation")
	if err := cmd.Run(); err != nil {
		logError(fmt.Sprintf("Finished logrotation with error: %v", err))
	} else {
		logInfo("Finished logrotation")
	}
}
