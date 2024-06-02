package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/creack/pty"
	"github.com/fsnotify/fsnotify"
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


func eventListener(restart chan bool, wg *sync.WaitGroup) {
    defer wg.Done()

    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Println("Error creating watcher:", err)
        return
    }
    defer watcher.Close()
    curdir, _ := os.Getwd()
    err = watcher.Add(curdir)
    if err != nil {
        log.Println("Error adding path to watcher:", err)
        return
    }

    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }

            // Restart if any file changes occur
            if event.Op&fsnotify.Write == fsnotify.Write || 
               event.Op&fsnotify.Create == fsnotify.Create ||
               event.Op&fsnotify.Remove == fsnotify.Remove ||
               event.Op&fsnotify.Rename == fsnotify.Rename {

                restart <- true
            }
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            log.Println("Watcher error:", err)
        }
    }
}

func runCommand(wg *sync.WaitGroup, run string, args []string, restart chan bool) {
    defer wg.Done()

    for {
        cmd := exec.Command(run, args...)
        ptmx, err := pty.Start(cmd)
        if err != nil {
            log.Fatal(err)
        }
        defer ptmx.Close()

        go io.Copy(os.Stdout, ptmx)
        go io.Copy(ptmx, os.Stdin)
        // Create a channel to signal when the command has finished
        done := make(chan error, 1) 
        
        // Wait for the command to finish in a goroutine
        go func() {
            done <- cmd.Wait() // Send the error (if any) when done
        }()

        select {
        case err := <-done: 
            if err != nil {
                log.Println("Command exited with error:", err)
            } else {
                log.Println("Command exited successfully")
            }
            restart <- false
            
        case <-restart:
            log.Println("Restarting command...")
            // Gracefully stop the command
            if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
                log.Println("Error sending SIGTERM:", err)
                if err := cmd.Process.Kill(); err != nil {
                    log.Fatal("Error killing process:", err)
                }
            }
            // Wait for the command to actually terminate
            cmd.Wait()
        }
    }
}
