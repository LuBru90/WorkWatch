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
const TIMEFORMAT = "01.02.2006 15:04:05 Mon"

// print-colors:
const HEADER = "\033[95m"
const OKBLUE = "\033[94m"
const OKCYAN = "\033[96m"
const OKGREEN = "\033[92m"
const WARNING = "\033[93m"
const FAIL = "\033[91m"
const ENDC = "\033[0m"
const BOLD = "\033[1m"
const UNDERLINE = "\033[4m"

// Valid commands and argument counts
var CMDS = map[string]int{
                            // cmd: argcount
                            "log": 1,
                            "stop": 1,
                            "add": 2,
                            "rm": 1,
                            "init": 1,
                            "test": 1,
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

    fmt.Print("\n" + BOLD + out + " " + output + " ")

    if (int(spacer) % 2) == 0 {
        fmt.Println(out)
    } else {
        fmt.Println(out + "-")
    }
    fmt.Print(ENDC)
}

func getCurrentTime() (string) {
    return time.Now().Format(TIMEFORMAT)
}

func convertStringToTime(asctime string) (tTime time.Time) {
    tTime, err := time.Parse(TIMEFORMAT, asctime)
    check(err)
    return tTime
}

func getTimeDiff(t1, t0 time.Time) (tDiff time.Duration) {
    return t1.Sub(t0)
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
    output := getCurrentTime() + " " + event
    err := _writeToFile(PATH, output, os.O_APPEND)
    return err
}

func printDuration(dur time.Duration) {
    fmt.Print(OKGREEN + ">>>> ")
    fmt.Println(dur)
    fmt.Print(ENDC + "\n")
}

func printLogLine(line string) {
    fmt.Print(OKBLUE)
    fmt.Print(timeFromLog(line))
    fmt.Print(ENDC + "\n")
    fmt.Println(getMessageFromLog(line))
}

// shows the content of the logfile and times
func log() {
    var t1, t0 time.Time
    var dur time.Duration
    var skipDur bool

    data := _getFileAsList(PATH)

    PrintT("Log")
    for i, line := range data[:len(data)-1] {
        if len(line) != 0 {
            t1 = timeFromLog(line)
            if i != 0 && !skipDur {
                dur = getTimeDiff(t1, t0)
                printDuration(dur)
                skipDur = false
            } else {
                fmt.Println()
            }
            skipDur = getMessageFromLog(line) == STOP
            t0 = t1
        }
        if i != len(data) - 2 {
            printLogLine(line)
        }
    }
    status()
}

// TODO: create logfile if does not exits
func initFile() {
    //TODO: write path to config file
    add("INIT")
    fmt.Println("Use: ww add '<text>'")
}

func timeFromLog(line string) (t time.Time) {
    return convertStringToTime(strings.Join(strings.Split(line, " ")[:3], " "))
}

func getMessageFromLog(line string) (out string) {
    return strings.Join(strings.Split(line, " ")[3:], " ")
}

// shows last entry of logfile
func status() {
    PrintT("Current Status")
    temp := _getFileAsList(PATH)
    if len(temp) > 1 {
        t0 := timeFromLog(temp[len(temp) - 2])
        t1, err := time.Parse(TIMEFORMAT, getCurrentTime())
        check(err)
        if err == nil {
            //fmt.Println(temp[len(temp) - 2])
            printLogLine(temp[len(temp) - 2])
            printDuration(getTimeDiff(t1, t0))
        }
    } else {
        fmt.Println("File empty! Use ww init")
    }
    fmt.Println()
}

// removes last entry of logfile
func remove() {
    data := _getFileAsList(PATH)
    err := os.Truncate(PATH, 0) // clear file
    check(err)

    err = _writeToFile(PATH, strings.Join(data[:len(data) - 2], "\n"), os.O_APPEND)
    check(err)
}

func test() {
    //out :=_getFileAsList(PATH)
    //return getMessageFromLog(out[0])
    //return temp
    //fmt.Println("Log-Message:", temp)
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
                status()
            case cmd == "init":
                initFile()
                log()
            case cmd == "stop":
                stop()
                log()
            case cmd == "rm":
                remove()
                log()
            case cmd == "log":
                log()
            case cmd == "test":
                test()
            }
        }
    }
}

func main() {
    args := os.Args[1:]
    //fmt.Println("Args:", strings.Join(args, " "))
    if len(args) != 0 {
        processUserCmd(args)
    } else {
        status()
    }
}
