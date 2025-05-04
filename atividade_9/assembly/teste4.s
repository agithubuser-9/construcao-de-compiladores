.section .bss
.lcomm a, 8
.lcomm b, 8
.lcomm c, 8
.lcomm soma, 8

.section .text
.globl _start
_start:
mov $10, %rax
mov %rax, a

mov $20, %rax
mov %rax, b

mov $30, %rax
mov %rax, c

mov $0, %rax
mov %rax, soma

mov a, %rax
push %rax
mov $10, %rax
pop %rbx
cmp %rax, %rbx
sete %al
movzx %al, %rax
cmp $0, %rax
je else_2
mov soma, %rax
push %rax
mov a, %rax
pop %rbx
cmp %rax, %rbx
add %rax, %rbx
mov %rbx, %rax
mov %rax, soma
jmp fim_if_2
else_2:
fim_if_2:
mov b, %rax
push %rax
mov $20, %rax
pop %rbx
cmp %rax, %rbx
setle %al
movzx %al, %rax
cmp $0, %rax
je else_3
mov soma, %rax
push %rax
mov b, %rax
pop %rbx
cmp %rax, %rbx
add %rax, %rbx
mov %rbx, %rax
mov %rax, soma
jmp fim_if_3
else_3:
fim_if_3:
mov c, %rax
push %rax
mov $30, %rax
pop %rbx
cmp %rax, %rbx
setge %al
movzx %al, %rax
cmp $0, %rax
je else_4
mov soma, %rax
push %rax
mov c, %rax
pop %rbx
cmp %rax, %rbx
add %rax, %rbx
mov %rbx, %rax
mov %rax, soma
jmp fim_if_4
else_4:
fim_if_4:
mov soma, %rax
call imprime_num
call sair
