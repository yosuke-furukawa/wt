package main

import (
	"fmt"
	"gopkg.in/fsnotify.v1"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// usage for wt
const USAGE = `
Usage: %s filepath command [interval(sec)]
`

func main() {

	fmt.Println(os.Args)
	if len(os.Args) < 3 {
		fmt.Printf(USAGE, os.Args[0])
		os.Exit(1)
	}
	filepath := os.Args[1]
	command := os.Args[2]
	interval := 1000
	fmt.Println(filepath)
	fmt.Println(command)
	if len(os.Args) > 3 {
		t, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		interval = t
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			now := time.Now()
			select {
			case event := <-watcher.Events:
				if time.Since(now) > time.Duration(interval)*time.Millisecond {
					if event.Op >= 1 {
						log.Println("file ", event.Name)
						var cmd *exec.Cmd
						cmd = exec.Command(command)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						err := cmd.Run()
						if err != nil {
							log.Fatal(err)
						}
					}
					now = time.Now()
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(filepath)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
