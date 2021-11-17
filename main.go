package main

import (
    "fmt"
    "os"
    "os/exec"
    "time"
    "errors"
    "encoding/json"
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
    fmt.Println("WARNING: '.fabric' is a quantum file that exists and doesn't exist simultaneously")
    return false
}

func printVersion() {
    dev := false
    if VERSION == "" {
        YYYYMMDDhhmmss := "20060102030405"
        VERSION = fmt.Sprintf("dev-%s", time.Now().Format(YYYYMMDDhhmmss))
        dev = true
    }
    fmt.Printf("fabric %s\r\n", VERSION)
    if (dev){
        fmt.Println("WARNING: running dev version of Fabric")
        fmt.Println("         fabric may be compiled without version info")
    }
}

func readFabric() (Fabric, bool) {
    data, err := os.ReadFile("./.fabric")
    fabric := Fabric{}
    var _fabric Fabric
    if err != nil {
        return fabric, false
    }
    err = json.Unmarshal(data, &_fabric)
    if err != nil {
        return fabric, false
    }
    return _fabric, true
}

func debugFabric(fabric Fabric, status bool) {
    fmt.Printf("fabric.build.command: %s\r\n", fabric.Build.Command)
    switch numArgs := len(fabric.Build.Args); numArgs {
    case 0:
    default:
        fmt.Printf("fabric.build.args: %s\r\n", fabric.Build.Args[0])
        for index := range fabric.Build.Args {
            if index != 0 {
                fmt.Printf("                   %s\r\n", fabric.Build.Args[index])
            }
        }
    }
    fmt.Println("")
}

func buildFabricProject(fabric Fabric) int {
    out, err := exec.Command(fabric.Build.Command, fabric.Build.Args...).CombinedOutput()
    if err != nil {
        fmt.Println("Error running fabric instructions");
        fmt.Println(err)
        return 1
    }
    if len(out) > 0 {
        fmt.Println(out)
    }
    fmt.Println("done.")
    return 0
}

func main() {
    printVersion()

    if !isFabric() {
        fmt.Println("ERROR: Current directory is not a fabric project")
        os.Exit(1)
    }
    fabric, status := readFabric()
    if !status {
        fmt.Println("ERROR: Cannot read or parse ./.fabric")
        os.Exit(1)
    }
    //debugFabric(fabric, status)
    os.Exit(buildFabricProject(fabric))
}
