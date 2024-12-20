  #
  # modelo de saida para o compilador
  #
  .section .text
  .globl _start

_start:

  mov $8, %rax
  imul $11, %rax
  mov %rax, %r8
  mov $9, %rax
  imul $12, %rax
  sub %rax, %r8
  mov $112, %r9
  sub $19, %r9
  add %r8, %r9

  mov %r9, %rax
  
  call imprime_num
  call sair

  .include "runtime.s"
