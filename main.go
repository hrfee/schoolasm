package main

import (
	"os"
	"strings"
)

func run(mem *memory) {
	mem[PC] = 0
	var lastInstruction value
	for {
		lastInstruction = mem[PC]
		op, ok := parseInstruction(mem[addr(mem[PC])], mem)
		if ok && op != nil {
			(*op).Exec()
		}
		if mem[PC] == lastInstruction {
			mem[PC]++
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
