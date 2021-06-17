### *school-asm*

A compiler and bytecode interpeter for a made up instruction set inspired by whtat was used in my A level exam questions.
```shell
go build github.com/hrfee/schoolasm
./schoolasm run examples/hello.asm
```

```
Usage: ./schoolasm [arguments] [run/build/exec] filename
run: compile and execute a program.
build: compile and write to <filename.sch>.
exec: run a compiled binary.
  -debug
    	print extra info when parsing & instruction info as they are executed. Doesn't play well with the table.
  -step int
    	exec/run only. Wait this many milliseconds between each execution cycle.
  -table
    	exec/run only. show table of memory contents during execution. Enabling sets step to 500ms.
```

Instructions and encoded format are explained in `spec.txt`.

### *license*

no license, why would you want to modify/distribute this?
