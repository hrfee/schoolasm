; run with -width 4 -height 4 -offset 200, and an appropriate step value.
; e.g: schoolasm -width 4 -height 4 -offset 200 -step 20 -scale 80 run examples/canvas.asm
ZERO:
	LDR #0

LOOP:
	INC IX
	LDD IX
	CMPV #17
	JPE ZERO
	LDX #199
	CMPV #0
	JPE WHITE
	JMP BLACK

WHITE:
	LDM #1
	STX #199
	JMP LOOP

BLACK:
	LDM #0
	STX #199
	JMP LOOP
