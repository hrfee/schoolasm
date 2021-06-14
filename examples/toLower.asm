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

IN
STO #250
LDD shiftAmount
ADD #250
OUT
END

