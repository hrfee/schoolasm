package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

const (
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

type memTable struct {
	table  *tablewriter.Table
	mem    *memory
	cache  [][]string
	colors [][]tablewriter.Colors
}

func (t *memTable) genTable() {
	emptyCache := false
	if len(t.cache) == 0 || len(t.cache[0]) == 0 {
		emptyCache = true
	}
	t.table.ClearRows()
	data := make([][]string, 16)
	for i := range data {
		data[i] = make([]string, 16)
	}
	col := 0
	row := 0
	if len(t.colors) == 0 || len(t.colors[0]) == 0 {
		t.colors = make([][]tablewriter.Colors, 16)
		for i := range t.colors {
			t.colors[i] = make([]tablewriter.Colors, 16)
			for j := range t.colors[i] {
				t.colors[i][j] = tablewriter.Colors{0, tablewriter.FgWhiteColor}
			}
		}
	}
	for i, v := range t.mem {
		data[row][col] = fmt.Sprintf("%08X", v)
		if !emptyCache && data[row][col] != t.cache[row][col] {
			nc := tablewriter.BgWhiteColor
			if prevColor := t.colors[row][col]; prevColor != nil {
				if prevColor[0] == tablewriter.BgWhiteColor {
					nc = tablewriter.BgRedColor
				}
			}
			t.colors[row][col] = tablewriter.Colors{nc, tablewriter.FgBlackColor}
		} else {
			t.colors[row][col] = tablewriter.Colors{0, tablewriter.FgWhiteColor}
		}
		if (i+1)%16 == 0 {
			t.table.Rich(data[row], t.colors[row])
			row++
			col = 0
			continue
		}
		col++
	}
	t.cache = data
	Clear()
	t.table.Render()
	fmt.Println("--------OUTPUT--------")
	fmt.Println(outContent)
	fmt.Println("----------------------")
}

func Clear() {
	cmd := exec.Command("/usr/bin/clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
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
		time.Sleep(STEP * time.Millisecond)
	}
}

func main() {
	Out = stdout{}
	fname := os.Args[1]
	content, err := os.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(content), "\n")
	mem := populateMemory(lines)
	var table *memTable
	if TABLE {
		header := make([]string, 16)
		for i := range header {
			header[i] = strconv.Itoa(i)
		}
		footer := make([]string, 16)
		footer[15] = "IX"
		footer[14] = "ACC"
		footer[13] = "PC"
		footer[12] = "COMP"
		table = &memTable{
			table: tablewriter.NewWriter(os.Stdout),
			mem:   mem,
		}
		table.table.SetHeader(header)
		table.table.SetFooter(footer)
		Clear()
	}
	run(table, mem)
}
