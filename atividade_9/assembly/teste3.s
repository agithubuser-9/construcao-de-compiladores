.section .bss
.lcomm n, 8
.lcomm m, 8
.lcomm soma, 8

.section .text
.globl _start
_start:
mov $1, %rax
mov %rax, n

mov $10, %rax
mov %rax, m

mov $0, %rax
mov %rax, soma

inicio_while_1:
mov n, %rax
push %rax
mov m, %rax
pop %rbx
cmp %rax, %rbx
setl %al
movzx %al, %rax
cmp $0, %rax
je fim_while_1
mov soma, %rax
push %rax
mov n, %rax
pop %rbx
cmp %rax, %rbx
add %rax, %rbx
mov %rbx, %rax
mov %rax, soma
mov n, %rax
push %rax
mov $1, %rax
pop %rbx
cmp %rax, %rbx
add %rax, %rbx
mov %rbx, %rax
mov %rax, n
jmp inicio_while_1
fim_while_1:
mov soma, %rax
call imprime_num
call sair
