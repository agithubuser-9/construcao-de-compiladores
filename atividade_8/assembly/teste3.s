.section .bss

.section .text
.globl _start
_start:
# return
mov $1, %rax
call imprime_num
call sair

.include "runtime.s"
