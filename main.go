package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

func parseArgs() (string, []string) {
    cmdPtr := flag.String("cmd", "", "Command to run ")
    updatePtr := flag.Bool("update", false, "Update LiveCode")
    baseCmd := flag.String("base", "", "Base command to run")
    flag.Parse()

    if *updatePtr  {
        fmt.Println("Updating...")
        Run("curl -fsSL https://raw.githubusercontent.com/joshiojas/LiveCode/main/uninstall.sh | bash")
        fmt.Println("Uninstall complete")
        Run("curl -fsSL https://raw.githubusercontent.com/joshiojas/LiveCode/main/install.sh | bash")
        fmt.Println("Update complete")
        os.Exit(0)
    } 
    if *cmdPtr == "" {
        fmt.Println("Error: -cmd flag is required")
        flag.Usage()
        os.Exit(1)
    }

    if *baseCmd != "" {
        c := *baseCmd
        fmt.Println("Running base command...")
        Run(c)
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
    
    curdir, _ := os.Getwd()
    fmt.Println("Watching directory: ", curdir)

    go runCommand(&wg, run, args, restart)
    go eventListener(restart, &wg)

    wg.Wait()
}
