package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func assertArgType(got, wanted, lineNum int) {
	if got == 1 && wanted == 0 {
		panic(fmt.Errorf("%d: Wanted address, got value constant", lineNum))
	} else if got == 0 && wanted == 1 {
		panic(fmt.Errorf("%d: Wanted value constant, got address", lineNum))
	}
}

func populateMemory(file []string) *memory {
	var lineCount uint16 = 1
	labels := map[string]addr{}
	labeledValues := map[string]addr{}
	firstInstruction := true
	var mem memory
	for lineNum, l := range file {
		isLabeledValue := false
		if l == "" || l[0] == ';' {
			continue
		}
		labelSects := strings.Split(l, ":")
		if strings.Contains(l, ":") {
			if len(labelSects) == 1 {
				// labeled section
				labels[labelSects[0]] = addr(lineCount)
				continue
			} else {
				// labeled value
				isLabeledValue = true
			}
		}
		sects := strings.Split(l, " ")
		var code string
		i := 0
		for _, c := range sects[0] {
			if !unicode.IsUpper(c) {
				break
			}
			i++
		}
		code = sects[0][0:i]
		var arg uint32
		hasArg := false
		// 0 == address, 1 == constant.
		argType := 0
		if len(sects) != 1 {
			hasArg = true
			argString := ""
			for _, c := range sects[1] {
				if c != ' ' {
					argString += string(c)
				}
			}
			switch argString[0] {
			case 'B':
				i, err := strconv.ParseUint(argString[1:], 2, 16)
				if err != nil {
					panic(fmt.Errorf("%d: Error parsing binary constant: %v", lineNum, err))
				}
				arg = uint32(i)
				argType = 1
			case '#':
				i, err := strconv.ParseUint(argString[1:], 10, 16)
				if err != nil {
					panic(fmt.Errorf("%d: Error parsing number: %v", lineNum, err))
				}
				arg = uint32(i)
				argType = 1
			default:
				i, err := strconv.ParseUint(argString, 2, 16)
				argType = 0
				if err != nil {
					var address addr
					if argString == "ACC" {
						address = ACC
					} else if argString == "IX" {
						address = IX
					} else {
						var ok bool
						// Println("Looking up label")
						// Printf("labels: %+v\n", labels)
						// Printf("labeled: %+v\n", labeledValues)
						address, ok = labels[argString]
						if !ok {
							a, ok := labeledValues[argString]
							if !ok {
								panic(fmt.Errorf("%d: Error parsing address: %v", lineNum, err))
							} else {
								argType = 0
							}
							address = addr(a)
						}
					}
					arg = uint32(address)
				} else {
					arg = uint32(i)
				}
			}
		}
		if code != "IN" && code != "OUT" && code != "END" {
			if !hasArg {
				panic(fmt.Errorf("%d: No argument given when required: %s", lineNum, l))
			}
			mem[lineCount] = (value(argType) << 31) + value(arg)
		}
		isInstruction := true
		switch code {
		case "LDM":
			assertArgType(argType, 1, lineNum)
			mem[lineCount] += value(O_LDM) << 15
		case "LDD":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_LDD) << 15
		case "LDI":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_LDI) << 15
		case "LDX":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_LDX) << 15
		case "LDR":
			assertArgType(argType, 1, lineNum)
			mem[lineCount] += value(O_LDR) << 15
		case "STO":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_STO) << 15
		case "ADD":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_ADD) << 15
		case "INC":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_INC) << 15
		case "DEC":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_DEC) << 15
		case "JMP":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_JMP) << 15
		case "CMP":
			mem[lineCount] += value(O_CMP) << 15
		case "JPE":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_JPE) << 15
		case "JPN":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_JPN) << 15
		case "JGT":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_JGT) << 15
		case "JLT":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_JLT) << 15
		case "IN":
			mem[lineCount] += value(O_IN) << 15
		case "OUT":
			mem[lineCount] += value(O_OUT) << 15
		case "END":
			mem[lineCount] += value(O_END) << 15
		default:
			isInstruction = false
			if !isLabeledValue {
				Printf("%d: Skipping line: %v\n", lineNum, l)
				continue
			}
		}
		if firstInstruction && isInstruction {
			Println("First instruction @", lineCount)
			mem[0] = value(O_JMP<<15) + value(lineCount)
			firstInstruction = false
		}
		if isLabeledValue {
			labeledValues[labelSects[0]] = addr(lineCount)
			mem[lineCount] = value(arg)
			isLabeledValue = false
		}
		// Only increment lineCount if line was valid op.
		lineCount++
	}
	return &mem
}

func parseInstruction(val value, mem *memory) (*Op, bool) {
	Println("COD", strconv.FormatUint(uint64(val), 2))
	opc := uint16(val >> 15)
	isConstant := opc&uint16(1<<15) == uint16(1<<15)
	if isConstant {
		opc -= uint16(1 << 15)
	}
	arg := uint16(val) - (opc << 15)
	opcode := Opcode(opc)
	var op Op
	ok := true
	switch opcode {
	case O_LDM:
		op = newLDM(value(arg), mem)
		Println("LDM")
	case O_LDD:
		op = newLDD(addr(arg), mem)
		Println("LDD")
	case O_LDI:
		op = newLDI(addr(arg), mem)
		Println("LDI")
	case O_LDX:
		op = newLDX(addr(arg), mem)
		Println("LDX")
	case O_LDR:
		op = newLDR(value(arg), mem)
		Println("LDR")
	case O_STO:
		op = newSTO(addr(arg), mem)
		Println("STO")
	case O_ADD:
		op = newADD(addr(arg), mem)
		Println("ADD")
	case O_INC:
		op = newINC(addr(arg), mem)
		Println("INC")
	case O_DEC:
		op = newDEC(addr(arg), mem)
		Println("DEC")
	case O_JMP:
		op = newJMP(addr(arg), mem)
		Println("JMP")
	case O_CMP:
		Println("CMP")
		if isConstant {
			op = newCMPval(value(arg), mem)
		} else {
			op = newCMPaddr(addr(arg), mem)
		}
	case O_JPE:
		op = newJPE(addr(arg), mem)
		Println("JPE")
	case O_JPN:
		op = newJPN(addr(arg), mem)
		Println("JPN")
	case O_JGT:
		op = newJGT(addr(arg), mem)
		Println("JGT")
	case O_JLT:
		op = newJLT(addr(arg), mem)
		Println("JLT")
	case O_IN:
		op = newIN(mem)
		Println("IN")
	case O_OUT:
		op = newOUT(mem)
		Println("OUT")
	case O_END:
		op = newEND()
		Println("END")
	default:
		ok = false
	}
	return &op, ok
}
