package main

import (
	"os"
	"strings"
)

const DEBUG = false

func Printf(a string, b ...interface{}) {
	if DEBUG {
		Printf(a, b...)
	}
}

func Println(a ...interface{}) {
	if DEBUG {
		Println(a...)
	}
}

func run(mem *memory) {
	mem[PC] = 0
	var lastInstruction value
	for {
		lastInstruction = mem[PC]
		Printf("ADDR %d ", mem[PC])
		op, ok := parseInstruction(mem[addr(mem[PC])], mem)
		if ok && op != nil {
			(*op).Exec()
		}
		if mem[PC] == lastInstruction {
			mem[PC]++
		} else {
			Println("JUMPED TO", mem[PC])
		}
	}
}

func main() {
	fname := os.Args[1]
	content, err := os.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(content), "\n")
	mem := populateMemory(lines)
	run(mem)
}
