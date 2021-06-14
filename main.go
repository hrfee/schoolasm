package main

import (
	"flag"
	"fmt"
	"os"
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
	fmt.Printf("Usage: %s [arguments] filename.asm\n", os.Args[0])
	flag.PrintDefaults()
}

func loadArgs() {
	flag.IntVar(&STEP, "step", STEP, "Wait this many milliseconds between each execution cycle.")
	flag.BoolVar(&DEBUG, "debug", DEBUG, "print extra info when parsing & instruction info as they are executed. Doesn't play well with the table.")
	flag.BoolVar(&TABLE, "table", TABLE, "show table of memory contents during execution. Enabling sets step to 500ms.")
	flag.Usage = argUsage
	flag.Parse()
}

func main() {
	loadArgs()
	Out = stdout{}
	fname := os.Args[len(os.Args)-1]
	if fname == "" || len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}
	content, err := os.ReadFile(fname)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(content), "\n")
	mem := populateMemory(lines)
	var table *memTable
	if TABLE {
		table = NewTable(mem)
		Clear()
	}
	run(table, mem)
}
