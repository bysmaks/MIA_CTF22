#!/bin/bash

nasm -f elf32 asm_lerner.asm -o asm_lerner.o && ld asm_lerner.o -o asm_lerner -m elf_i386 -n
rm asm_lerner.o
chmod +x asm_lerner

#nasm -f elf32 elf.asm -o elf.o && ld elf.o -o elf -m elf_i386 -n
#rm elf.o
#chmod +x elf
