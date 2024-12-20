.section .text
.globl _start

_start:
  mov $7, %rax
  imul $6, %rax

  imul $5, %rax
  mov %rax, %r8

  mov $4, %rax
  imul $3, %rax

  imul $2, %rax

  imul $1, %rax
  mov %rax, %r9

  mov %r8, %rax
  xor %rdx, %rdx
  div %r9

  # imprime 8 ou o .75 esta no rdx  
  call imprime_num
  call sair

  .include "runtime.s"
