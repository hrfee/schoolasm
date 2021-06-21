; move a dot around a 8x8 grid with the arrow keys.

; run with -width 8 -height 8 -offset 180
; e.g: schoolasm -width 8 -height 8 -offset 180 -scale 80 run examples/move.asm

x: #0
y: #0

startOffset: #180
currentOffset: #0
LDM #1
; initial position
STO #180

width: #8
; maximum x & y coords (0-indexed, so 8-1=7)
maxX: #0
maxY: #0
LDD width
DEC ACC
STO maxX
STO maxY

; size = width^2 (width assumed equal to height)
size: #0
count: #0
SIZELOOP:
	LDD count
	CMPA width
	JPE SETKEYS
	INC ACC
	STO count
	LDD size
	ADD width
	STO size
	JMP SIZELOOP

upAddr: #0
downAddr: #0
leftAddr: #0
rightAddr: #0

SETKEYS:
	LDM #0
	ADD startOffset
	ADD size
	INC ACC
	STO upAddr
	INC ACC
	STO downAddr
	INC ACC
	STO leftAddr
	INC ACC
	STO rightAddr

JMP EVLOOP

; convert offset to x,y form
XY:
	cmp: #0
	LDM #0
	STO x
	STO y
	STO cmp
	STO count
	; add the width (4) to cmp, check if greater than current offset, if not increment the y-value.
	Y:
		LDD cmp
		ADD width
		STO cmp
		CMPA currentOffset
		JGT Xpre
		LDD y
		INC ACC
		STO y
		JMP Y
	; add (y*width) to cmp before calculating x.
	Xpre:
		LDM #0
		STO cmp
		STO count
		; stuck here
		Xpre2:
			LDD count
			CMPA y
			JPE X
			INC ACC
			STO count
			LDD cmp
			ADD width
			STO cmp
			JMP Xpre2
	; increment x and cmp until cmp = current offset
	X:
		LDD cmp
		CMPA currentOffset
		JPE LUP
		INC ACC
		STO cmp
		LDD x
		INC ACC
		STO x
		JMP X

; mem[a] - mem[b] = mem[c]
; only used to move upwards, so return address is hardcoded.
a: #0
b: #0
c: #0
subcount: #0
SUB:
	LDD a
	STO c
	LDM #0
	STO subcount
	LOOP:
		LDD c
		DEC ACC
		STO c
		LDD subcount
		INC ACC
		STO subcount
		CMPA b
		JPN LOOP
		JMP upret

; cache for previous key states.
upCache: #0
downCache: #0
leftCache: #0
rightCache: #0

; loop and c
EVLOOP:
	; convert offset to x-y coords
	JMP XY
	LUP:
		LDI upAddr
		CMPA upCache
		STO upCache
		JPE LDOWN
		CMPV #1
		JPE up
	LDOWN:
		LDI downAddr
		CMPA downCache
		STO downCache
		JPE LLEFT
		CMPV #1
		JPE down
	LLEFT:
		LDI leftAddr
		CMPA leftCache
		STO leftCache
		JPE LRIGHT
		CMPV #1
		JPE left
	LRIGHT:
		LDI rightAddr
		CMPA rightCache
		STO rightCache
		JPE LUP
		CMPV #1
		JPE right
	JMP EVLOOP

up:
	LDD y
	CMPV #0
	JPE LUP
	LDD currentOffset
	STO a
	LDD width
	STO b
	JMP SUB
	upret:
	LDD currentOffset
	ADD startOffset
	STA #0
	LDD c
	STO currentOffset
	ADD startOffset
	STA #1
	JMP EVLOOP

down:
	LDD y
	CMPA maxY
	JPE LUP
	LDD currentOffset
	ADD startOffset
	STA #0
	ADD width
	STA #1
	LDD currentOffset
	ADD width
	STO currentOffset
	JMP EVLOOP

left:
	LDD x
	CMPV #0
	JPE LUP
	LDD currentOffset
	ADD startOffset
	STA #0
	LDD currentOffset
	DEC ACC
	STO currentOffset
	ADD startOffset
	STA #1
	JMP EVLOOP

right:
	LDD x
	CMPA maxX
	JPE LUP
	LDD currentOffset
	ADD startOffset
	STA #0
	LDD currentOffset
	INC ACC
	STO currentOffset
	ADD startOffset
	STA #1
	JMP EVLOOP
