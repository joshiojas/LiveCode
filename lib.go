package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

func Run(c string) {
	fmt.Println("Running base command...")
	cmd := exec.Command("bash", "-c", c)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Running command: ", cmd.String())
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatal(err)
	}
	defer ptmx.Close()

	go io.Copy(os.Stdout, ptmx)
	go io.Copy(ptmx, os.Stdin)

	done := make(chan error, 1) 
        
	// Wait for the command to finish in a goroutine
	go func() {
		done <- cmd.Wait() // Send the error (if any) when done
	}()

	for <-done != nil {
		fmt.Println("Error running base command")
		os.Exit(1)
	}
}