package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

var stdout = os.Stdout
var stderr = os.Stderr

const (
	// man console_codes
	colorRed    = iota + 31 // 31
	colorGreen              // 32
	colorYellow             // 33
)

func main() {
	log.SetFlags(log.Lshortfile)
	if err := run(os.Args[1:]); err != nil {
		// log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(args []string) error {
	// defer は 宣言順と逆に実行される
	// reader.ReadLine() は w.Close() を実行しないとロックを解除しないため
	// w.Closeより先に<-cを宣言する必要がある
	c := make(chan struct{}, 1)
	defer func() {
		<-c
	}()
	var wg sync.WaitGroup
	wg.Add(2)
	defer wg.Wait()

	args = append([]string{"test"}, args...)
	r, w := io.Pipe()
	defer w.Close()

	cmd := exec.Command("go", args...)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Env = os.Environ()

	if err := cmd.Start(); err != nil {
		log.Print(err)
		return err
	}

	go func() {
		defer func() {
			c <- struct{}{}
		}()
		reader := bufio.NewReader(r)
		for {
			l, _, err := reader.ReadLine()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Print(err)
				return
			}
			s := string(l)

			if strings.Contains(s, "PASS") {
				fmt.Fprintf(stdout, "\x1b[%vm%s\x1b[0m\n", colorGreen, s)
			} else if strings.Contains(s, "FAIL") {
				fmt.Fprintf(stderr, "\x1b[%vm%s\x1b[0m\n", colorRed, s)
			} else {
				fmt.Fprintf(stdout, "%s\n", s)
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			if ws.ExitStatus() != 0 {
				return fmt.Errorf("process status %d", ws.ExitStatus())
			}
			return nil
		}
		return err
	}
	return nil
}
