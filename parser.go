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

func populateMemory(file []string) memory {
	var lineCount uint16 = 0
	labels := map[string]addr{}
	labeledValues := map[string]uint16{}
	var mem memory
	for lineNum, l := range file {
		isLabeledValue := false
		if l == "" || l[0] == ';' {
			continue
		}

		sects := strings.Split(l, ":")
		if strings.Contains(l, ":") {
			if len(sects) == 1 {
				// labeled section
				labels[sects[0]] = addr(lineCount)
				continue
			} else {
				// labeled value
				isLabeledValue = true
			}
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
		var arg uint16
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
				arg = uint16(i)
				argType = 1
			case '#':
				i, err := strconv.ParseUint(argString[1:], 10, 16)
				if err != nil {
					panic(fmt.Errorf("%d: Error parsing number: %v", lineNum, err))
				}
				arg = uint16(i)
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
						address, ok = labels[argString]
						if !ok {
							a, ok := labeledValues[argString]
							if !ok {
								panic(fmt.Errorf("%d: Error parsing address: %v", lineNum, err))
							}
							address = addr(a)
						}
					}
					arg = uint16(address)
				} else {
					arg = uint16(i)
				}
			}
			if isLabeledValue {
				labeledValues[sects[0]] = arg
				isLabeledValue = false
			}
		}
		if code != "IN" && code != "OUT" && code != "END" {
			if !hasArg {
				panic(fmt.Errorf("%d: No argument given when required", lineNum))
			}
			mem[lineCount] = (value(argType) << 32) + value(arg)
		}
		switch code {
		case "LDM":
			assertArgType(argType, 1, lineNum)
			mem[lineCount] += value(O_LDM) << 16
		case "LDD":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_LDD) << 16
		case "LDI":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_LDI) << 16
		case "LDX":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_LDX) << 16
		case "LDR":
			assertArgType(argType, 1, lineNum)
			mem[lineCount] += value(O_LDR) << 16
		case "STO":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_STO) << 16
		case "ADD":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_ADD) << 16
		case "INC":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_INC) << 16
		case "DEC":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_DEC) << 16
		case "JMP":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_JMP) << 16
		case "CMP":
			if argType == 0 {
				mem[lineCount] += value(O_CMPA) << 16
			} else {
				mem[lineCount] += value(O_CMPV) << 16
			}
		case "JPE":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_JPE) << 16
		case "JPN":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_JPN) << 16
		case "JGT":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_JGT) << 16
		case "JLT":
			assertArgType(argType, 0, lineNum)
			mem[lineCount] += value(O_JLT) << 16
		case "IN":
			mem[lineCount] += value(O_IN) << 16
		case "OUT":
			mem[lineCount] += value(O_OUT) << 16
		case "END":
			mem[lineCount] += value(O_END) << 16
		default:
			fmt.Printf("%d: Skipping line: %v", lineNum, l)
			continue
		}
		// Only increment lineCount if line was valid op.
		lineCount++
	}
	return mem
}
