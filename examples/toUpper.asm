shiftAmount: #32
E: #69
n: #110
t: #116
e: #101
r: #114
spc: #32
a: #97
c: #99
h: #104
colon: #58
nl: #10

enteracharacter:
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
	LDD a
	OUT
	LDD spc
	OUT
	LDD c
	OUT
	LDD h
	OUT
	LDD a
	OUT
	LDD r
	OUT
	LDD a
	OUT
	LDD c
	OUT
	LDD t
	OUT
	LDD e
	OUT
	LDD r
	OUT
	LDD colon
	OUT
	LDD spc
	OUT

; i := 0
LDM #0
STO #249
; Get char and store in 250
IN
STO #250

; Valid characters (a-zA-Z) from 65-122
; Jump to end if outside, or already uppercase (<90)
INRANGE:
	CMPV #65
	JLT EXIT
	CMPV #90
	JLT EXIT
	CMPV #122
	JGT EXIT

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

EXIT:
	LDD #250
	OUT
	END
