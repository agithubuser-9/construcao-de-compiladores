.section .data
buf: .space 20
nl:  .byte 10

.section .text
.global imprime_num
.global sair

imprime_num:
    mov %rax, %rdi
    mov $buf+19, %rsi
    mov $10, %rcx
    mov $0, %rdx

.conv_loop:
    xor %rdx, %rdx
    div %rcx
    add $48, %rdx
    dec %rsi
    mov %dl, (%rsi)
    test %rax, %rax
    jnz .conv_loop

    mov $1, %rax
    mov $1, %rdi
    mov %rsi, %rsi
    mov $buf+19, %rdx
    sub %rsi, %rdx
    syscall

    mov $1, %rax
    mov $1, %rdi
    mov $nl, %rsi
    mov $1, %rdx
    syscall

    ret

sair:
    mov $60, %rax
    xor %rdi, %rdi
    syscall
