.section .bss

.section .text
.globl _start
_start:
mov $7, %rax
push %rax
mov $6, %rax
pop %rbx
cmp %rax, %rbx
imul %rax, %rbx
mov %rbx, %rax
call imprime_num
call sair
