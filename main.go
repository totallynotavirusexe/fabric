package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "time"
    "reflect"

    "github.com/fatih/color"
)

var VERSION string

type Instructions struct {
    Command string `json:"command"`
    Args []string `json:"args"`
}

type Fabric struct {
    Build Instructions `json:"build"`
    Install Instructions `json:"install,omitempty"`
    Clean Instructions `json:"clean,omitempty"`
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

func run(instructions Instructions) error {
    cmd := exec.Command(instructions.Command, instructions.Args...)
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    return cmd.Run()
}

func buildFabricProject(fabric Fabric) bool {
    if reflect.DeepEqual(fabric.Build, Instructions{}) {
        fmt.Printf("%s: No build instructions available for project\r\n", color.RedString("ERROR"))
        return false
    }
    err := run(fabric.Build)
    if err != nil {
        fmt.Printf("%s: %v\r\n", color.RedString("ERROR"), err)
        fmt.Printf("%s: Could not fabricate project into reality\r\n", color.RedString("ERROR"))
        fmt.Printf("       Ran %s %s\r\n", fabric.Build.Command, strings.Join(fabric.Build.Args, " "))
        return false
    }
    fmt.Println("Fabricated project into reality!")
    return true
}

func installFabricProject(fabric Fabric) bool {
    if reflect.DeepEqual(fabric.Install, Instructions{}) {
        fmt.Printf("%s: No install instructions available for project\r\n", color.RedString("ERROR"))
        return false
    }
    err := run(fabric.Install)
    if err != nil {
        fmt.Printf("%s: %v\r\n", color.RedString("ERROR"), err)
        fmt.Printf("%s: Could not fabricate project into reality\r\n", color.RedString("ERROR"))
        fmt.Printf("%s: Failed to install project\r\n", color.RedString("ERROR"))
        fmt.Printf("       Ran %s %s\r\n", fabric.Build.Command, strings.Join(fabric.Build.Args, " "))
        return false
    }
    fmt.Println("Fabricated project into reality!")
    fmt.Println("Installed project!")
    return true
}

func cleanFabricProject(fabric Fabric) bool {
    if reflect.DeepEqual(fabric.Clean, Instructions{}) {
        fmt.Printf("%s: No clean instructions available for project\r\n", color.RedString("ERROR"))
        return false
    }
    err := run(fabric.Clean)
    if err != nil {
        fmt.Printf("%s: %v\r\n", color.RedString("ERROR"), err)
        fmt.Printf("%s: Could not clean project\r\n", color.RedString("ERROR"))
        fmt.Printf("       Ran %s %s\r\n", fabric.Build.Command, strings.Join(fabric.Build.Args, " "))
        return false
    }
    fmt.Println("Cleaned project!")
    return true
}

func main() {
    printVersion()
    args := os.Args[1:]

    if !isFabric() {
        fmt.Printf("%s: Current directory is not a fabric project\r\n", color.RedString("ERROR"))
        os.Exit(1)
    }
    fabric, status := readFabric()
    if !status {
        fmt.Printf("%s: Cannot read or parse ./.fabric\r\n", color.RedString("ERROR"))
        os.Exit(1)
    }

    if len(args) == 0 {
        args = []string {"help"}
    }

    completed := true
    switch v := args[0]; v {
    case "build":
        completed = buildFabricProject(fabric);
    case "install":
        completed = installFabricProject(fabric)
    case "clean":
        completed = cleanFabricProject(fabric)
    default:
        fmt.Println()
        fmt.Printf("%s: command '%s' not supported\r\n", color.RedString("ERROR"), v)
        fallthrough
    case "help":
        fmt.Println()
        fmt.Println("Available commands:")
        fmt.Println("help\t\tshow this message")
        fmt.Println("build\t\tbuild project")
        fmt.Println("install\t\tbuild and install project")
        fmt.Println("clean\t\tclean project")
        fmt.Println()
        fmt.Printf("%s: availability of these commands is based on the project .fabric file\r\n", color.CyanString("NOTE"))
    }

    switch completed {
    case true:
        os.Exit(0)
    case false:
        os.Exit(1)
    }
}
