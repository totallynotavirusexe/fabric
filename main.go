package main

import (
    "fmt"
    "os"
    "os/exec"
    "time"
    "errors"
    "encoding/json"
    "strings"

    "github.com/fatih/color"
)

var VERSION string

type Build struct {
    Command string `json:"command"`
    Args []string `json:"args"`
}

type Fabric struct {
    Build Build `json:"build"`
}

func isFabric() bool {
    _, err := os.Stat("./.fabric")
    if err == nil {
        return true
    }
    if errors.Is(err, os.ErrNotExist) {
        return false
    }
    fmt.Printf("%s: '.fabric' is a quantum file that exists and doesn't exist simultaneously\r\n", color.YellowString("WARNING"))
    return false
}

func printVersion() {
    dev := false
    if VERSION == "" {
        YYYYMMDDhhmmss := "20060102030405"
        VERSION = fmt.Sprintf("dev-%s", time.Now().Format(YYYYMMDDhhmmss))
        dev = true
    }
    fmt.Printf("Fabric %s\r\n", VERSION)
    if (dev){
        fmt.Printf("%s: running dev version of Fabric\r\n", color.YellowString("WARNING"))
        fmt.Println("         fabric may be compiled without version info")
    }
}

func readFabric() (Fabric, bool) {
    data, err := os.ReadFile("./.fabric")
    var _fabric Fabric
    if err != nil {
        return _fabric, false
    }
    err = json.Unmarshal(data, &_fabric)
    if err != nil {
        return _fabric, false
    }
    return _fabric, true
}

func buildFabricProject(fabric Fabric) bool {
    cmd := exec.Command(fabric.Build.Command, fabric.Build.Args...)
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    err := cmd.Run()
    if err != nil {
        fmt.Printf("%s: %v\r\n", color.RedString("ERROR"), err)
        return false
    }
    return true
}

func main() {
    printVersion()

    if !isFabric() {
        fmt.Printf("%s: Current directory is not a fabric project\r\n", color.RedString("ERROR"))
        os.Exit(1)
    }
    fabric, status := readFabric()
    if !status {
        fmt.Printf("%s: Cannot read or parse ./.fabric\r\n", color.RedString("ERROR"))
        os.Exit(1)
    }

    switch status := buildFabricProject(fabric); status {
    case true:
        fmt.Println("Fabricated project into reality!")
        os.Exit(0)
    case false:
        fmt.Printf("%s: Could not fabricate project into reality\r\n", color.RedString("ERROR"))
        fmt.Printf("       Ran %s %s\r\n", fabric.Build.Command, strings.Join(fabric.Build.Args, " "))
        os.Exit(1)
    }
}
