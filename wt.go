package main

import (
	"fmt"
	"gopkg.in/fsnotify.v1"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const VERSION = "v0.0.1"

// usage for wt
const USAGE = `
Usage: %s filepath command [interval(sec default 1sec)]
`

func clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func runScript(watcher *fsnotify.Watcher, interval int, command string) {
	for {
		now := time.Now()
		select {
		case event := <-watcher.Events:
			if time.Since(now) > time.Duration(interval)*time.Millisecond {
				if event.Op >= 1 && !strings.HasPrefix(event.Name, ".") {
					clear()
					var cmd *exec.Cmd
					commands := strings.Split(command, " ")
					// TODO: windows support
					name := commands[0]
					cmd = exec.Command(name)
					cmd.Args = commands
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						log.Println("error:", err)
					}
				}
				now = time.Now()
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func contains(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

func main() {

	if len(os.Args) < 2 {
		log.Fatalf(USAGE, os.Args[0])
	}
	if contains(os.Args[1], []string{"--version", "-v", "version"}) {
		fmt.Printf("%s version: %s\n", os.Args[0], VERSION)
		return
	}
	if contains(os.Args[1], []string{"--help", "-h", "help"}) {
		fmt.Printf(USAGE, os.Args[0])
		return
	}
	if len(os.Args) < 3 {
		log.Fatalf(USAGE, os.Args[0])
	}
	filepath := os.Args[1]
	command := os.Args[2]
	interval := 1000
	if len(os.Args) > 3 {
		t, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		interval = t
	}

	clear()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go runScript(watcher, interval, command)

	err = watcher.Add(filepath)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
