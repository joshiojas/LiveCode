package main

import (
	"flag"
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

func parseArgs() (string, []string) {
    cmdPtr := flag.String("cmd", "", "Command to run (required)")
    flag.Parse()

    if *cmdPtr == "" {
        fmt.Println("Error: -cmd flag is required")
        flag.Usage()
        os.Exit(1)
    }

    return *cmdPtr, flag.Args()
}

func main() {
    run, args := parseArgs()

    var wg sync.WaitGroup
    wg.Add(1)
    restart := make(chan bool)
    fmt.Println("Starting command...")
    fmt.Println("Press Ctrl+C to exit")
    // print the watching directory
    curdir, _ := os.Getwd()
    fmt.Println("Watching directory: ", curdir)

    go runCommand(&wg, run, args, restart)
    go eventListener(restart, &wg)

    wg.Wait()
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