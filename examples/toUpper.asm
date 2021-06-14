shiftAmount: #32

; i := 0
LDM #0
STO #249
; Get char and store in 250
IN
STO #250

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

LDD #250
OUT
END
