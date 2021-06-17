TEST: "hello, "
TEST2: "world!"

one:
	LDM TEST
	STO STR
	LDM two
	STO RET
	JMP PRINT

two:
	LDM TEST2
	STO STR
	LDM end
	STO RET
	JMP PRINT

PRINT:
	; Stores start address of string to print.
	STR: #0
	; Stores address to jump to after print finished.
	RET: #0
	LOOP:
		LDI STR
		OUT
		CMPV #0
		JPE ENDPRINT
		LDD STR
		INC ACC
		STO STR
		JMP LOOP
	ENDPRINT:
		; newline
		LDM #10
		OUT
		LDD RET
		JMPA

end:
	END
