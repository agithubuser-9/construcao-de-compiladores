  #
  # modelo de saida para o compilador
  #
  .section .text
  .globl _start

_start:
  mov $72, %rax
  sub $101, %rax
  mov $4, %rbx
  mul %rbx
  mov %rax, %r8

  mov $14, %rax
  mov $77, %rbx
  mul %rbx
  mov %rax, %r9

  add %r8, %r9

  mov %r9, %rax
  
  call imprime_num
  call sair

  .include "runtime.s"
