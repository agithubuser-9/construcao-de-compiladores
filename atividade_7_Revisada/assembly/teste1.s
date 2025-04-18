.section .bss
.lcomm x, 8
.lcomm y, 8

.section .text
.globl _start
_start:
mov $2, %rax
mov %rax, x

mov $3, %rax
mov %rax, y

# return
mov x, %rax
push %rax
mov y, %rax
pop %rbx
add %rbx, %rax
call imprime_num
call sair

.include "runtime.s"
