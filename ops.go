package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

const (
	// memRoot is squared to get the memory size
	memRoot = 16
	memSize = memRoot * memRoot // 256 addresses
)

type memory [memSize]value

type addr uint32
type value uint32
type Opcode uint16

const (
	IX addr = memSize - 1 - iota
	ACC
	PC
	COMP

	O_LDM  Opcode = 1
	O_LDD         = 2
	O_LDI         = 3
	O_LDX         = 4
	O_LDR         = 5
	O_STO         = 6
	O_STX         = 7
	O_ADD         = 8
	O_INC         = 9
	O_DEC         = 10
	O_JMP         = 11
	O_JMPA        = 12
	O_CMPA        = 13
	O_CMPV        = 14
	O_JPE         = 15
	O_JPN         = 16
	O_JGT         = 17
	O_JLT         = 18
	O_IN          = 19
	O_OUT         = 20
	O_END         = 21
	O_AND         = 22
	O_OR          = 23
)

type Op interface {
	Exec()
}

type LDM struct {
	val value
	mem *memory
}

func newLDM(val value, mem *memory) LDM {
	return LDM{val, mem}
}

func (op LDM) Exec() {
	op.mem[ACC] = op.val
}

type LDD struct {
	src addr
	mem *memory
}

func newLDD(src addr, mem *memory) LDD {
	return LDD{src, mem}
}

func (op LDD) Exec() {
	op.mem[ACC] = op.mem[op.src]
}

type LDI struct {
	addressSrc addr
	mem        *memory
}

func newLDI(addressSrc addr, mem *memory) LDI {
	return LDI{addressSrc, mem}
}

func (op LDI) Exec() {
	op.mem[ACC] = op.mem[op.mem[op.addressSrc]]
}

type LDX struct {
	index addr
	mem   *memory
}

func newLDX(index addr, mem *memory) LDX {
	return LDX{index, mem}
}

func (op LDX) Exec() {
	op.mem[ACC] = op.mem[op.index+addr(op.mem[IX])]
}

type LDR struct {
	val value
	mem *memory
}

func newLDR(val value, mem *memory) LDR {
	return LDR{val, mem}
}

func (op LDR) Exec() {
	op.mem[IX] = op.val
}

type STO struct {
	dest addr
	mem  *memory
}

func newSTO(dest addr, mem *memory) STO {
	return STO{dest, mem}
}

func (op STO) Exec() {
	op.mem[op.dest] = op.mem[ACC]
}

type STX struct {
	dest addr
	mem  *memory
}

func newSTX(dest addr, mem *memory) STX {
	return STX{dest, mem}
}

func (op STX) Exec() {
	op.mem[op.dest+addr(op.mem[IX])] = op.mem[ACC]
}

type ADD struct {
	src addr
	mem *memory
}

func newADD(src addr, mem *memory) ADD {
	return ADD{src, mem}
}

func (op ADD) Exec() {
	op.mem[ACC] += op.mem[op.src]
}

type INC struct {
	reg addr
	mem *memory
}

func newINC(reg addr, mem *memory) INC {
	if reg != ACC && reg != IX {
		reg = ACC
	}
	return INC{reg, mem}
}

func (op INC) Exec() {
	op.mem[op.reg] += 1
	if DEBUG {
		reg := "ACC"
		if op.reg == IX {
			reg = "IX"
		}
		Println("INC'd", reg, "to", op.mem[op.reg])
	}
}

type DEC struct {
	reg addr
	mem *memory
}

func newDEC(reg addr, mem *memory) DEC {
	if reg != ACC && reg != IX {
		reg = ACC
	}
	return DEC{reg, mem}
}

func (op DEC) Exec() {
	op.mem[op.reg] -= 1
	if DEBUG {
		reg := "ACC"
		if op.reg == IX {
			reg = "IX"
		}
		Println("DEC'd", reg, "to", op.mem[op.reg])
	}
}

type JMP struct {
	loc addr
	mem *memory
}

func newJMP(loc addr, mem *memory) JMP {
	return JMP{loc, mem}
}

func (op JMP) Exec() {
	op.mem[PC] = value(op.loc)
}

type JMPA struct {
	mem *memory
}

func newJMPA(mem *memory) JMPA {
	return JMPA{mem}
}

func (op JMPA) Exec() {
	op.mem[PC] = value(op.mem[ACC])
}

// CMP #n
type CMPval struct {
	val value
	mem *memory
}

func newCMPval(val value, mem *memory) CMPval {
	return CMPval{val, mem}
}

func (op CMPval) Exec() {
	if op.mem[ACC] > op.val {
		op.mem[COMP] = 2
	} else if op.mem[ACC] < op.val {
		op.mem[COMP] = 0
	} else {
		op.mem[COMP] = 1
	}
}

type CMPaddr struct {
	src addr
	mem *memory
}

func newCMPaddr(src addr, mem *memory) CMPaddr {
	return CMPaddr{src, mem}
}

func (op CMPaddr) Exec() {
	(&CMPval{
		val: op.mem[op.src],
		mem: op.mem,
	}).Exec()
}

type JPE struct {
	loc addr
	mem *memory
}

func newJPE(loc addr, mem *memory) JPE {
	return JPE{loc, mem}
}

func (op JPE) Exec() {
	if op.mem[COMP] == 1 {
		op.mem[PC] = value(op.loc)
	}
}

type JPN struct {
	loc addr
	mem *memory
}

func newJPN(loc addr, mem *memory) JPN {
	return JPN{loc, mem}
}

func (op JPN) Exec() {
	if op.mem[COMP] != 1 {
		op.mem[PC] = value(op.loc)
	}
}

type JGT struct {
	loc addr
	mem *memory
}

func newJGT(loc addr, mem *memory) JGT {
	return JGT{loc, mem}
}

func (op JGT) Exec() {
	if op.mem[COMP] == 2 {
		op.mem[PC] = value(op.loc)
	}
}

type JLT struct {
	loc addr
	mem *memory
}

func newJLT(loc addr, mem *memory) JLT {
	return JLT{loc, mem}
}

func (op JLT) Exec() {
	if op.mem[COMP] == 0 {
		op.mem[PC] = value(op.loc)
	}
}

// When IN is called, the character from the position (pos) in the buffer is loaded into the ACC. pos is then incremented.
// The buffer stores the most recent line.
// When a new line is entered, pos is set to zero and the line is written to buffer.
type stdinBuffer struct {
	buffer []byte
	pos    int
}

// Capture continously writes new stdin input to the buffer.
func (b *stdinBuffer) Capture() {
	var err error
	reader := bufio.NewReader(os.Stdin)
	for {
		b.buffer, err = reader.ReadBytes('\n')
		b.pos = 0
		Println("NEW INPUT:", strings.ReplaceAll(string(b.buffer), "\n", "\\n"))
		if err != nil {
			panic(fmt.Sprintf("Failed to read from os.Stdin: %v", err))
		}
	}
}

type IN struct {
	mem *memory
}

func newIN(mem *memory) IN {
	return IN{mem}
}

func (op IN) Exec() {
	for StdinBuffer.buffer == nil || len(StdinBuffer.buffer) == 0 {
		continue
	}
	char := string(StdinBuffer.buffer)[StdinBuffer.pos]
	if char > unicode.MaxASCII {
		panic(fmt.Errorf("Character was outside ASCII range"))
	}
	op.mem[ACC] = value(char)
	StdinBuffer.pos++
}

type stdout struct{}

func (w stdout) Write(p []byte) (n int, err error) {
	if TABLE {
		outContent += string(p)
	} else {
		return os.Stdout.Write(p)
	}
	return len(p), nil
}

type OUT struct {
	mem *memory
}

func newOUT(mem *memory) OUT {
	return OUT{mem}
}

func (op OUT) Exec() {
	Println("OUTING", op.mem[ACC])
	out := []byte{byte(op.mem[ACC])}
	n, err := Out.Write(out)
	if n != 1 || err != nil {
		panic(err)
	}
}

type END struct{}

func newEND() END { return END{} }

func (op END) Exec() {
	os.Exit(0)
}

type AND struct {
	src addr
	mem *memory
}

func newAND(src addr, mem *memory) AND {
	return AND{src, mem}
}

func (op AND) Exec() {
	op.mem[ACC] = (op.mem[ACC]) & (op.mem[op.src])
}

type OR struct {
	src addr
	mem *memory
}

func newOR(src addr, mem *memory) OR {
	return OR{src, mem}
}

func (op OR) Exec() {
	op.mem[ACC] = (op.mem[ACC]) | (op.mem[op.src])
}
