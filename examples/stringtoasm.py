#!/usr/bin/env python3

specials = {
    " ": "spc",
    "\n": "nl",
    "\r": "cr",
    ":": "colon",
    "!": "excl"
}

input = input("String to convert: ") + "\n"
values = {}
print("\n")
for c in input:
    if c not in values:
        out = c
        if c in specials:
            out = specials[c]
        values[c] = ord(c)
        print(out + ": #" + str(values[c]))
print("")
for c in input:
    if c in specials:
        c = specials[c]
    print("LDD " + c)
    print("OUT")
