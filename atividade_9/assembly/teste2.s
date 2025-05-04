.section .bss
.lcomm a, 8
.lcomm b, 8
.lcomm c, 8
.lcomm delta, 8

.section .text
.globl _start
_start:
mov $1, %rax
mov %rax, a

mov $2, %rax
mov %rax, b

mov $3, %rax
mov %rax, c

mov b, %rax
push %rax
mov b, %rax
pop %rbx
cmp %rax, %rbx
imul %rax, %rbx
mov %rbx, %rax
push %rax
mov $4, %rax
push %rax
mov a, %rax
pop %rbx
cmp %rax, %rbx
imul %rax, %rbx
mov %rbx, %rax
push %rax
mov c, %rax
pop %rbx
cmp %rax, %rbx
imul %rax, %rbx
mov %rbx, %rax
pop %rbx
cmp %rax, %rbx
sub %rax, %rbx
mov %rbx, %rax
mov %rax, delta

mov delta, %rax
push %rax
mov $0, %rax
pop %rbx
cmp %rax, %rbx
setl %al
movzx %al, %rax
cmp $0, %rax
je else_0
mov $0, %rax
push %rax
mov delta, %rax
pop %rbx
cmp %rax, %rbx
sub %rax, %rbx
mov %rbx, %rax
mov %rax, delta
jmp fim_if_0
else_0:
mov delta, %rax
mov %rax, delta
fim_if_0:
mov delta, %rax
call imprime_num
call sair
