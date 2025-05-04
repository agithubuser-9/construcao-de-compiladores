.section .bss
.lcomm a, 8
.lcomm C, 8
.lcomm B, 8

.section .text
.globl _start
_start:
mov $1, %rax
push %rax
mov $1, %rax
pop %rbx
mov $0, %rdx
idiv %rbx
push %rax
mov $2, %rax
pop %rbx
add %rbx, %rax
mov %rax, a

mov $3, %rax
push %rax
mov $2, %rax
push %rax
mov $2, %rax
pop %rbx
mul %rbx
pop %rbx
add %rbx, %rax
mov %rax, C

mov $4, %rax
mov %rax, B

# return
mov a, %rax
push %rax
mov B, %rax
pop %rbx
add %rbx, %rax
push %rax
mov C, %rax
pop %rbx
add %rbx, %rax
call imprime_num
call sair

.include "runtime.s"
