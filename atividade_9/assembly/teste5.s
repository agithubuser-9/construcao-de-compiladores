.section .bss
.lcomm a, 8
.lcomm b, 8
.lcomm c, 8
.lcomm d, 8

.section .text
.globl _start
_start:
mov $5, %rax
mov %rax, a

mov $3, %rax
mov %rax, b

mov $0, %rax
mov %rax, c

mov $1, %rax
mov %rax, d

inicio_while_5:
mov a, %rax
push %rax
mov $0, %rax
pop %rbx
cmp %rax, %rbx
setg %al
movzx %al, %rax
cmp $0, %rax
je fim_while_5
mov b, %rax
push %rax
mov $5, %rax
pop %rbx
cmp %rax, %rbx
setl %al
movzx %al, %rax
cmp $0, %rax
je else_6
mov c, %rax
push %rax
mov b, %rax
pop %rbx
cmp %rax, %rbx
add %rax, %rbx
mov %rbx, %rax
mov %rax, c
mov b, %rax
push %rax
mov $1, %rax
pop %rbx
cmp %rax, %rbx
add %rax, %rbx
mov %rbx, %rax
mov %rax, b
jmp fim_if_6
else_6:
mov c, %rax
push %rax
mov $1, %rax
pop %rbx
cmp %rax, %rbx
add %rax, %rbx
mov %rbx, %rax
mov %rax, c
fim_if_6:
mov a, %rax
push %rax
mov $1, %rax
pop %rbx
cmp %rax, %rbx
sub %rax, %rbx
mov %rbx, %rax
mov %rax, a
jmp inicio_while_5
fim_while_5:
mov c, %rax
call imprime_num
call sair
