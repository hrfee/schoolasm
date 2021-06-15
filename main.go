package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	DEBUG = false
	TABLE = false
	STEP  = 0
)

type stdout struct{}

func (w stdout) Write(p []byte) (n int, err error) {
	if TABLE {
		outContent += string(p)
	} else {
		return os.Stdout.Write(p)
	}
	return len(p), nil
}

var Out stdout
var outContent string

func Printf(a string, b ...interface{}) {
	if DEBUG {
		fmt.Printf(a, b...)
	}
}

func Println(a ...interface{}) {
	if DEBUG {
		fmt.Println(a...)
	}
}

func run(table *memTable, mem *memory) {
	mem[PC] = 0
	var lastInstruction value
	for {
		lastInstruction = mem[PC]
		Printf("ADDR %d ", mem[PC])
		op, ok := parseInstruction(mem[addr(mem[PC])], mem)
		if ok && op != nil {
			(*op).Exec()
		}
		if TABLE {
			table.genTable()
		}
		if mem[PC] == lastInstruction {
			mem[PC]++
		} else {
			Println("JUMPED TO", mem[PC])
		}
		time.Sleep(time.Duration(STEP) * time.Millisecond)
	}
}

func argUsage() {
	fmt.Printf("Usage: %s [arguments] [run/build/exec] filename\n", os.Args[0])
	fmt.Println(`run: compile and execute a program.
build: compile and write to <filename.sch>.
exec: run a compiled binary.`)
	flag.PrintDefaults()
}

func loadArgs() {
	flag.BoolVar(&DEBUG, "debug", DEBUG, "print extra info when parsing & instruction info as they are executed. Doesn't play well with the table.")
	flag.IntVar(&STEP, "step", STEP, "exec/run only. Wait this many milliseconds between each execution cycle.")
	flag.BoolVar(&TABLE, "table", TABLE, "exec/run only. show table of memory contents during execution. Enabling sets step to 500ms.")
	flag.Usage = argUsage
	flag.Parse()
}

func main() {
	loadArgs()
	Out = stdout{}
	fname := os.Args[len(os.Args)-1]
	if fname == "" || len(os.Args) == 1 || fname == "run" || fname == "build" || fname == "exec" {
		flag.Usage()
		os.Exit(1)
	}
	runType := os.Args[len(os.Args)-2]
	if len(os.Args) == 2 {
		flag.Usage()
		os.Exit(1)
	}
	content, err := os.ReadFile(fname)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}
	var mem *memory
	if runType == "run" || runType == "build" {
		lines := strings.Split(string(content), "\n")
		mem = populateMemory(lines)
		if runType == "build" {
			name := strings.TrimSuffix(fname, filepath.Ext(fname)) + ".sch"
			out := MarshalMemory(mem)
			err := os.WriteFile(name, out, 0666)
			if err != nil {
				fmt.Printf("Failed to write file: %v", err)
				os.Exit(1)
			}
			fmt.Println("Written to", name)
			os.Exit(0)
		}
	} else if runType == "exec" {
		mem = UnmarshalMemory(content)
	}
	if runType == "run" || runType == "exec" {
		var table *memTable
		if TABLE {
			table = NewTable(mem)
			Clear()
		}
		run(table, mem)
	}
}
