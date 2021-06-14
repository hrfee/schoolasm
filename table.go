package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

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

func NewTable(mem *memory) *memTable {
	header := make([]string, memRoot)
	for i := range header {
		header[i] = strconv.Itoa(i)
	}
	footer := make([]string, memRoot)
	footer[memRoot-1] = "IX"
	footer[memRoot-2] = "ACC"
	footer[memRoot-3] = "PC"
	footer[memRoot-4] = "COMP"
	table := &memTable{
		table: tablewriter.NewWriter(os.Stdout),
		mem:   mem,
	}
	table.table.SetHeader(header)
	table.table.SetFooter(footer)
	return table
}
