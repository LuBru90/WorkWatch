package main

// Imports
import "os"
import "time"
import "fmt"
import "io/ioutil"
import "strings"

// TODO: config parser
//import "github.com/bigkevmcd/go-configparser"


// Constants
// TODO: check if file exists
const PATH = "timelog.log"
const STOP = "<>--- STOP ---<>"

// Valid commands and argument counts
var CMDS = map[string]int{
                            // cmd: argcount
                            "log": 1,
                            "stop": 1,
                            "add": 2,
                            "rm": 1,
                            "init": 1,
                        }


func check(e error) {
    if e != nil {
        panic(e)
    }
}

func PrintT(output string) () {
    var out string
    var i int
    var spacer float32

    spacer = float32((50 - len(output)))
    spacer += 0.5
    spacer /= 2
    for i = 0; float32(i) < spacer; i++ {
        out += "-"
    }

    fmt.Print(out + " " + output + " ")

    if (int(spacer) % 2) == 0 {
        fmt.Println(out)
    } else {
        fmt.Println(out + "-")
    }
}

func getCurrentTime() (string) {
    return time.Now().Format("01.02.2006 15:04:05 Mon")
}

func _writeToFile(path string, content string, operation int) (error) {
    file, err := os.OpenFile(path, operation, 0644)
    check(err)

    _, err = file.WriteString(content + "\n")
    check(err)

    file.Sync()
    return err
}

func _getFileAsList(path string) ([]string) {
    data, err := ioutil.ReadFile(PATH)
    check(err)
    return strings.Split(string(data), "\n")
}

// Adds a new line to the logfile
func add(event string) (error) {
    output := getCurrentTime() + ": " + event
    err := _writeToFile(PATH, output, os.O_APPEND)
    return err
}

// shows the content of the logfile
func log() {
    data := _getFileAsList(PATH)
    PrintT("Log")
    for _, line := range data {
        fmt.Println(line)
    }
}

// TODO: create logfile if does not exits
func initFile() {
    //TODO: write path to config file
    add("INIT")
}

// shows last entry of logfile
func status() {
    PrintT("Status")
    temp := _getFileAsList(PATH)
    if len(temp) > 1 {
        fmt.Println(temp[len(temp)-2])
    } else {
        fmt.Println("File empty! Use ww init")
    }
}

// shows time difference between alle log entrys 
func showTimes() {
    fmt.Println("TODO: Calculate time diffs and print log with diffs")
}

// removes last entry of logfile
func remove() {
    data := _getFileAsList(PATH)
    err := os.Truncate(PATH, 0) // clear file
    check(err)

    err = _writeToFile(PATH, strings.Join(data[:len(data) - 2], "\n"), os.O_APPEND)
    check(err)
}

// adds a new stop-event to the logfile
func stop() {
    add(STOP)
}

func processUserCmd(args []string) {
    cmd := args[0]
    if argcount, ok := CMDS[cmd]; ok {
        if argcount == len(args) {
            switch {
            case cmd == "add":
                add(strings.Join(args[1:], " "))
                log()
            case cmd == "init":
                initFile()
                log()
                fmt.Println("Use: ww add '<text>'")
            case cmd == "stop":
                stop()
                log()
            case cmd == "times":
                showTimes()
            case cmd == "rm":
                remove()
                log()
            case cmd == "log":
                log()
            }
        }
    }
}

func main() {
    args := os.Args[1:]
    fmt.Println("Args:", strings.Join(args, " "))
    if len(args) != 0 {
        processUserCmd(args)
    } else {
        status()
    }
}
