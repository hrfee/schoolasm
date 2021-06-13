package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

type memory [256]value

type addr uint32
type value uint32
type Opcode uint16

const (
	IX   addr = 255
	ACC       = 254
	PC        = 253
	COMP      = 252

	O_LDM Opcode = 1
	O_LDD        = 2
	O_LDI        = 3
	O_LDX        = 4
	O_LDR        = 5
	O_STO        = 6
	O_ADD        = 7
	O_INC        = 8
	O_DEC        = 9
	O_JMP        = 10
	O_CMP        = 11
	O_JPE        = 12
	O_JPN        = 13
	O_JGT        = 14
	O_JLT        = 15
	O_IN         = 16
	O_OUT        = 17
	O_END        = 18
)

type Op interface {
	Exec()
}

//

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

type IN struct {
	mem *memory
}

func newIN(mem *memory) IN {
	return IN{mem}
}

// Currently, the enter keypress is required for the op to unblock, so entering strings isn't possible.
// if multiple are entered, the first character is taken.
func (op IN) Exec() {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		panic(err)
	}
	if char > unicode.MaxASCII {
		panic(fmt.Errorf("Character was outside ASCII range"))
	}
	op.mem[ACC] = value(char)
}

type OUT struct {
	mem *memory
}

func newOUT(mem *memory) OUT {
	return OUT{mem}
}

func (op OUT) Exec() {
	out := []byte{byte(op.mem[ACC])}
	n, err := os.Stdout.Write(out)
	if n != 1 || err != nil {
		panic(err)
	}
}

type END struct{}

func newEND() END { return END{} }

func (op END) Exec() {
	os.Exit(0)
}
