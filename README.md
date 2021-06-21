### *school-asm*

A compiler and bytecode interpeter for a made up instruction set inspired by whtat was used in my A level exam questions.
```shell
go build github.com/hrfee/schoolasm
./schoolasm run examples/hello.asm
```
more cool examples in the [examples](https://github.com/hrfee/schoolasm/tree/main/examples) folder.

built with support for a simple black&white canvas with arrowkey input by default. This depends on SDL, and can take a long time to build at first. It can be disabled with `go build -tags=nocanvas`.

```
Usage: ./schoolasm [arguments] [run/build/exec] filename
run: compile and execute a program.
build: compile and write to <filename.sch>.
exec: run a compiled binary.
  -debug
    	print extra info when parsing & instruction info as they are executed. Doesn't play well with the table.
  -width int
    	width of canvas window. Disabled if blank.
  -height int
    	height of canvas window. Disabled if blank.
  -offset int
    	starting address of memory used to set pixels for the canvas window. Goes by row, then column.
  -scale int
    	scale pixel size for canvas. (default 10)
  -showmem string
    	comma-separated list of named/decimal addresses to show the value of on each cycle. named addresses only available with run.
  -step int
    	exec/run only. Wait this many milliseconds between each execution cycle.
  -table
    	exec/run only. show table of memory contents during execution. Enabling sets step to 500ms.
```

Instructions and encoded format are explained in `spec.txt`.

### *license*

no license, why would you want to modify/distribute this?
