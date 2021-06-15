; input a string, and the uppercase version is returned.
shiftAmount: #32
E: #69
n: #110
t: #116
e: #101
r: #114
spc: #32
s: #115
o: #111
m: #109
x: #120
colon: #58
nl: #10

entersometext:
	LDD E
	OUT
	LDD n
	OUT
	LDD t
	OUT
	LDD e
	OUT
	LDD r
	OUT
	LDD spc
	OUT
	LDD s
	OUT
	LDD o
	OUT
	LDD m
	OUT
	LDD e
	OUT
	LDD spc
	OUT
	LDD t
	OUT
	LDD e
	OUT
	LDD x
	OUT
	LDD t
	OUT
	LDD colon
	OUT
	LDD spc
	OUT

LOADCHAR:
	; i := 0
	LDM #0
	STO #249
	; Get char and store in 250
	IN
	STO #250

; Valid characters (a-zA-Z) from 65-122
; Jump to end if outside, or already uppercase (<90)
INRANGE:
	; Check if newline
	CMPA nl
	JPE EXIT
	CMPV #65
	JLT OUTPUT
	CMPV #90
	JLT OUTPUT
	CMPV #122
	JGT OUTPUT

DECLOOP:
	; Load charcode, decrement, then store.
	LDD #250
	DEC ACC
	STO #250
	; Increment counter
	LDD #249
	INC ACC
	STO #249
	; Compare counter with 32. If not equal, continue looping.
	CMPA shiftAmount
	JPN DECLOOP

OUTPUT:
	LDD #250
	OUT
	JMP LOADCHAR

EXIT:
	LDD nl
	OUT
	END
