  #
  # modelo de saida para o compilador
  #
  .section .text
  .globl _start

_start:

  mov $8, %rax
  mov $11, %rbx
  mul %rbx
  mov %rax, %r8

  mov $9, %rax
  mov $12, %rbx
  mul %rbx
  sub %rax, %r8

  mov $112, %r9
  sub $19, %r9
  add %r8, %r9

  mov %r9, %rax
  
  call imprime_num
  call sair

  .include "runtime.s"
