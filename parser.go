package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type unsatisfiedLabel struct {
	address addr
	lineNum int
}

func populateMemory(file []string) (*memory, map[string]addr) {
	var lineCount uint16 = 1
	labels := map[string]addr{}
	labeledValues := map[string]addr{}
	// If usage of a label is encountered before creation,
	// unsatisfiedLabels[name] = line num
	// At end, if any remaining 65536's, error.
	unsatisfiedLabels := map[string][]unsatisfiedLabel{}
	firstInstruction := true
	showAddresses := map[string]addr{}
	var mem memory
	for lineNum, l := range file {
		isLabeledValue := false
		if l == "" {
			continue
		}
		for l[0] == ' ' || l[0] == '\t' {
			l = l[1:]
		}
		if l[0] == ';' || l[0] == '#' {
			continue
		}
		labelSects := strings.Split(l, ":")
		if strings.Contains(l, ":") {
			if len(labelSects) == 1 || labelSects[1] == "" {
				// labeled section
				labels[labelSects[0]] = addr(lineCount)
				if labels, ok := unsatisfiedLabels[labelSects[0]]; ok {
					for _, label := range labels {
						Printf("Satisfying unsatisfied label \"%s\" on line %d with address %d\n", labelSects[0], label.lineNum, label.address)
						mem[label.address] += value(lineCount)
						delete(unsatisfiedLabels, labelSects[0])
					}
				}
				for _, v := range SHOWMEM {
					if labelSects[0] == v {
						showAddresses[labelSects[0]] = addr(lineCount)
						break
					}
				}
				continue
			} else {
				// labeled value
				isLabeledValue = true
			}
		}
		var sects []string
		out := ""
		inQuote := false
		for _, c := range l {
			if c == '"' {
				inQuote = !inQuote
			}
			if c == ' ' && !inQuote {
				sects = append(sects, out)
				out = ""
				continue
			}
			out += string(c)
		}
		if out != "" {
			sects = append(sects, out)
		}
		var code string
		i := 0
		for _, c := range sects[0] {
			if !unicode.IsUpper(c) {
				break
			}
			i++
		}
		code = sects[0][0:i]
		var arg []uint32
		hasArg := false
		if len(sects) != 1 {
			hasArg = true
			argString := ""
			inQuote := false
			for _, c := range sects[1] {
				if inQuote || c != ' ' {
					argString += string(c)
					if c == '"' {
						inQuote = !inQuote
					}
				}
			}
			switch argString[0] {
			case '#':
				i, err := strconv.ParseUint(argString[1:], 10, 16)
				if err != nil {
					panic(fmt.Errorf("%d: Error parsing number: %v", lineNum, err))
				}
				arg = []uint32{uint32(i)}
			case '"':
				i := len(argString) - 1
				for argString[i] != '"' && i >= 0 {
					i--
				}
				if i == 0 {
					continue
				}
				replacer := strings.NewReplacer("\\n", "\n", "\\\"", "\"")
				content := replacer.Replace(argString[1:i])
				arg = make([]uint32, len(content)+1)
				for i := range content {
					arg[i] = uint32(content[i])
				}
			default:
				if argString[0] == 'B' && (argString[1] == '0' || argString[1] == '1') {
					i, err := strconv.ParseUint(argString[1:], 2, 16)
					if err != nil {
						panic(fmt.Errorf("%d: Error parsing binary constant: %v", lineNum, err))
					}
					arg = []uint32{uint32(i)}
				} else {
					i, err := strconv.ParseUint(argString, 2, 16)
					if err != nil {
						i, err = strconv.ParseUint(argString, 10, 16)
					}
					if err != nil {
						var address addr
						if argString == "ACC" {
							address = ACC
						} else if argString == "IX" {
							address = IX
						} else {
							var ok bool
							address, ok = labels[argString]
							if !ok {
								a, ok := labeledValues[argString]
								if !ok {
									unsatisfiedLabels[argString] = append(unsatisfiedLabels[argString], unsatisfiedLabel{address: addr(lineCount), lineNum: lineNum})
									address = 0
								} else {
									address = addr(a)
								}
							}
						}
						arg = []uint32{uint32(address)}
					} else {
						arg = []uint32{uint32(i)}
					}
				}
			}
		}
		if code != "IN" && code != "OUT" && code != "END" && code != "JMPA" {
			if !hasArg {
				panic(fmt.Errorf("%d: No argument given when required: %s", lineNum, l))
			}
			i := lineCount
			for _, v := range arg {
				mem[i] = value(v)
				i++
			}
		}
		isInstruction := true
		switch code {
		case "LDM":
			mem[lineCount] += value(O_LDM) << 15
		case "LDD":
			mem[lineCount] += value(O_LDD) << 15
		case "LDI":
			mem[lineCount] += value(O_LDI) << 15
		case "LDX":
			mem[lineCount] += value(O_LDX) << 15
		case "LDR":
			mem[lineCount] += value(O_LDR) << 15
		case "STO":
			mem[lineCount] += value(O_STO) << 15
		case "STA":
			mem[lineCount] += value(O_STA) << 15
		case "STX":
			mem[lineCount] += value(O_STX) << 15
		case "ADD":
			mem[lineCount] += value(O_ADD) << 15
		case "INC":
			mem[lineCount] += value(O_INC) << 15
		case "DEC":
			mem[lineCount] += value(O_DEC) << 15
		case "JMP":
			mem[lineCount] += value(O_JMP) << 15
		case "JMPA":
			mem[lineCount] += value(O_JMPA) << 15
		case "CMPA":
			mem[lineCount] += value(O_CMPA) << 15
		case "CMPV":
			mem[lineCount] += value(O_CMPV) << 15
		case "JPE":
			mem[lineCount] += value(O_JPE) << 15
		case "JPN":
			mem[lineCount] += value(O_JPN) << 15
		case "JGT":
			mem[lineCount] += value(O_JGT) << 15
		case "JLT":
			mem[lineCount] += value(O_JLT) << 15
		case "IN":
			mem[lineCount] += value(O_IN) << 15
		case "OUT":
			mem[lineCount] += value(O_OUT) << 15
		case "END":
			mem[lineCount] += value(O_END) << 15
		case "AND":
			mem[lineCount] += value(O_AND) << 15
		case "OR":
			mem[lineCount] += value(O_OR) << 15
		case "XOR":
			mem[lineCount] += value(O_XOR) << 15
		case "WMI":
			mem[lineCount] += value(O_WMI) << 15
		case "PRINT":
			/*
				LDR #0
				LOOP:
					LDI <address>
					OUT
					INC IX)
					CMPV #10
					JPN LOOP
			*/
			mem[lineCount] -= value(arg[0])
			mem[lineCount] += value(O_LDR) << 15
			lineCount++
			loop := lineCount
			mem[loop] += (value(O_LDX) << 15) + value(arg[0])
			lineCount++
			mem[lineCount] += value(O_OUT) << 15
			lineCount++
			mem[lineCount] += (value(O_INC) << 15) + value(IX)
			lineCount++
			mem[lineCount] += (value(O_CMPV) << 15) + 0
			lineCount++
			mem[lineCount] += (value(O_JPN) << 15) + value(loop)
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
			if labels, ok := unsatisfiedLabels[labelSects[0]]; ok {
				for _, label := range labels {
					Printf("Satisfying unsatisfied label \"%s\" on line %d with address %d\n", labelSects[0], label.lineNum, label.address)
					mem[label.address] += value(lineCount)
					delete(unsatisfiedLabels, labelSects[0])
				}
			}
			for _, v := range SHOWMEM {
				if labelSects[0] == v {
					showAddresses[labelSects[0]] = addr(lineCount)
					break
				}
			}
			if len(arg) > 1 {
				lineCount += uint16(len(arg) - 1)
			}
			isLabeledValue = false
		}
		// Only increment lineCount if line was valid op.
		lineCount++
	}
	errOut := ""
	for label, line := range unsatisfiedLabels {
		errOut += fmt.Sprintf("%d: Label \"%s\" not defined\n", line, label)
	}
	if errOut != "" {
		panic(errors.New(errOut))
	}
	return &mem, showAddresses
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
	case O_STA:
		op = newSTA(value(arg), mem)
		Println("STA")
	case O_STX:
		op = newSTX(addr(arg), mem)
		Println("STX")
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
	case O_JMPA:
		op = newJMPA(mem)
		Println("JMPA")
	case O_CMPA:
		Println("CMP Address")
		op = newCMPaddr(addr(arg), mem)
	case O_CMPV:
		Println("CMP Value")
		op = newCMPval(value(arg), mem)
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
	case O_AND:
		op = newAND(addr(arg), mem)
		Println("AND")
	case O_OR:
		op = newOR(addr(arg), mem)
		Println("OR")
	case O_XOR:
		op = newXOR(addr(arg), mem)
		Println("XOR")
	case O_WMI:
		op = newWMI(value(arg))
		Println("WMI")
	default:
		ok = false
	}
	return &op, ok
}

func MarshalMemory(mem *memory) []byte {
	out := make([]byte, 4*memSize)
	for j, v := range mem {
		i := j * 4
		b1 := byte((v & 0b11111111000000000000000000000000) >> 24)
		b2 := byte((v & 0b00000000111111110000000000000000) >> 16)
		b3 := byte((v & 0b00000000000000001111111100000000) >> 8)
		b4 := byte((v & 0b00000000000000000000000011111111))
		out[i] = b1
		out[i+1] = b2
		out[i+2] = b3
		out[i+3] = b4
	}
	return out
}

func UnmarshalMemory(mem []byte) *memory {
	var out memory
	for j := range out {
		i := j * 4
		val := value(mem[i]) << 24
		val += value(mem[i+1]) << 16
		val += value(mem[i+2]) << 8
		val += value(mem[i+3])
		out[j] = val
	}
	return &out
}
