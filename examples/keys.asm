; print arrow key events in the canvas window.
; schoolasm -width 4 -height 4 -offset 200 -step 20 -scale 80 run examples/keys.asm

upCache: #0
downCache: #0
leftCache: #0
rightCache: #0

LUP:
	LDD #217
	CMPA upCache
	STO upCache
	JPE LDOWN
	CMPV #1
	JPE up
LDOWN:
	LDD #218
	CMPA downCache
	STO downCache
	JPE LLEFT
	CMPV #1
	JPE down
LLEFT:
	LDD #219
	CMPA leftCache
	STO leftCache
	JPE LRIGHT
	CMPV #1
	JPE left
LRIGHT:
	LDD #220
	CMPA rightCache
	STO rightCache
	JPE LUP
	CMPV #1
	JPE right
	JMP LUP


keyUp: "arrow key up\n"
keyDown: "arrow key down\n"
keyLeft: "arrow key left\n"
keyRight: "arrow key right\n"

up:
	PRINT keyUp
	JMP LUP

down:
	PRINT keyDown
	JMP LUP

left:
	PRINT keyLeft
	JMP LUP

right:
	PRINT keyRight
	JMP LUP
