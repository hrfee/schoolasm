START:
	IN
	STO CHAR1
	IN
	STO CHAR2
	LDD CHAR1
LOOP:
	OUT
	CMPA CHAR2
	JPE ENDFOR
	INC ACC
	JMP LOOP
ENDFOR:
	END
CHAR1: #0
CHAR2: #0
