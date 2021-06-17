; From Q 6c) of M/J/16.
; SENSORS stores sensor readings for 4 sensors. 
; Smallest 4 bits stores whether the each sensor has triggered (1) or not (0)
; The system only activated if 2 or more sensors are triggered.

SENSORS: B00001010
COUNT: #0
VALUE: #1

ALARMON: B11111111
ALARMOFF: B00000000

LOOP:
	LDD SENSORS
	AND VALUE
	CMPV #0
	JPE ZERO
	LDD COUNT
	INC ACC
	STO COUNT
ZERO:
	LDD VALUE
	CMPV #8
	JPE EXIT
	ADD VALUE
	STO VALUE
	JMP LOOP
	LDD COUNT
EXIT:
	LDD COUNT
TEST:
	CMPV #1
	JGT ALARM

; ALARM added so the program makes sense, prints 'h'.
ALARM:
	LDM #104
	OUT
	END
