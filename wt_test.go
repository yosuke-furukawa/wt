package main

import (
	"log"
	"os/exec"
	"strings"
	"testing"
	"time"
)

type FakeWriter struct {
	t        *testing.T
	expected string
	passed   bool
}

func (fw *FakeWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	log.Println(message)
	if fw.passed {
		fw.t.Log("succeeded")
		return len(p), nil
	}
	if strings.Contains(message, fw.expected) {
		fw.t.Logf("success: %s", message)
		fw.passed = true
	}
	return len(p), nil
}

func TestWtWithHelp(t *testing.T) {
	c := exec.Command("go", "run", "./wt.go", "-h")
	fw := &FakeWriter{t, "Usage", false}
	c.Stdout = fw
	c.Run()
	if !fw.passed {
		t.Errorf("failed, it is not contained %s", fw.expected)
	}
}

func TestWtWithVersion(t *testing.T) {
	c := exec.Command("go", "run", "./wt.go", "-v")
	fw := &FakeWriter{t, "version", false}
	c.Stdout = fw
	c.Run()
	if !fw.passed {
		t.Errorf("failed, it is not contained %s", fw.expected)
	}
}

func TestWtls(t *testing.T) {
	cmd := exec.Command("go", "run", "./wt.go", ".", "ls")
	fw := &FakeWriter{t, "abcdef", false}
	lsDone := make(chan bool)
	touchDone := make(chan bool)
	go func() {
		go func() {
			cmd.Stdout = fw
			cmd.Run()
		}()
		<-touchDone
		lsDone <- true
	}()
	go func() {
		time.Sleep(3 * time.Second)
		c := exec.Command("touch", "abcdef")
		c.Run()
		time.Sleep(3 * time.Second)
		c = exec.Command("rm", "abcdef")
		c.Run()
		touchDone <- true
	}()
	<-lsDone
	if !fw.passed {
		t.Errorf("failed, it is not contained %s", fw.expected)
	}
}
