package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	DEBUG     = false
	TABLE     = false
	SHOWMEM   = []string{}
	STEP      = 0
	WIDTH     = 0
	HEIGHT    = 0
	GUIOFFSET = 0
	SCALE     = 10
)

var Out stdout
var outContent string

var StdinBuffer = stdinBuffer{}

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

func run(table *memTable, mem *memory, showAddresses map[string]addr) {
	mem[PC] = 0
	var lastInstruction value
	addressOrder := make([]string, len(showAddresses))
	if len(showAddresses) != 0 {
		i := 0
		for name := range showAddresses {
			addressOrder[i] = name
			i++
		}
		sort.Strings(addressOrder)
	}
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
		if len(showAddresses) != 0 {
			out := ""
			for _, name := range addressOrder {
				out += fmt.Sprintf("%s (%d): %08b ", name, showAddresses[name], mem[showAddresses[name]])
			}
			fmt.Println(out)
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
	if CANVAS {
		fmt.Println("Built with canvas support")
	} else {
		fmt.Println("Built without canvas support")
	}
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
	var showmem string
	flag.StringVar(&showmem, "showmem", showmem, "comma-separated list of named/decimal addresses to show the value of on each cycle. named addresses only available with run.")
	if CANVAS {
		flag.IntVar(&WIDTH, "width", WIDTH, "width of canvas window. Disabled if blank.")
		flag.IntVar(&HEIGHT, "height", HEIGHT, "height of canvas window. Disabled if blank.")
		flag.IntVar(&GUIOFFSET, "offset", GUIOFFSET, "starting address of memory used to set pixels for the canvas window. Goes by row, then column.")
		flag.IntVar(&SCALE, "scale", SCALE, "scale pixel size for canvas.")
	}
	flag.Usage = argUsage
	flag.Parse()
	if showmem != "" {
		SHOWMEM = strings.Split(showmem, ",")
	}
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
	showAddresses := map[string]addr{}
	if runType == "run" || runType == "build" {
		lines := strings.Split(string(content), "\n")
		mem, showAddresses = populateMemory(lines)
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
	for _, v := range SHOWMEM {
		switch v {
		case "IX":
			showAddresses["IX"] = IX
		case "ACC":
			showAddresses["ACC"] = ACC
		case "PC":
			showAddresses["PC"] = PC
		case "COMP":
			showAddresses["COM"] = COMP
		default:
			n, err := strconv.Atoi(v)
			if err == nil {
				showAddresses[v] = addr(n)
			}
		}
	}
	if runType == "run" || runType == "exec" {
		if CANVAS && WIDTH != 0 && HEIGHT != 0 {
			go newCanvas(WIDTH, HEIGHT, SCALE, addr(GUIOFFSET), mem)
		}
		var table *memTable
		if TABLE {
			table = NewTable(mem)
			Clear()
		}
		go StdinBuffer.Capture()
		run(table, mem, showAddresses)
	}
}
