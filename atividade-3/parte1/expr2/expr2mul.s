.section .text
.globl _start

_start:
  mov $7, %rax
  mov $6, %rbx
  mul %rbx

  mov $5, %rbx
  mul %rbx
  mov %rax, %r8

  mov $4, %rax
  mov $3, %rbx
  mul %rbx

  mov $2, %rbx
  mul %rbx

  mov $1, %rbx
  mul %rbx
  mov %rax, %r9

  mov %r8, %rax
  xor %rdx, %rdx
  div %r9

  # imprime 8 ou o .75 esta no rdx  
  call imprime_num
  call sair

  .include "runtime.s"
