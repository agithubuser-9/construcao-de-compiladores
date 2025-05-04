# Compilador YZY – Documentação de Execução
``` 
Alunos: 
Yuri Gabriel da Silva Alves
Cássio Andrêzza de Almeida
```

# 🛠️ Comandos principais

## 🔹 go run main.go
A execução da main irá retornar o resultado dos 5 arquivos de teste dentro da pasta 'resultados' e gerar os codigos intermediarios dentro da pasta 'assembly'
---

## 🔹 make

Deve ser executado na pasta RAIZ do projeto.
compila os arquivos assembly `.s` dentro na pasta `assembly/`, usando `runtime.s` como suporte. Os binários são salvos em `bin/`. 
---

## 🔹make tests

Executa e compara a saída de cada binário com o conteúdo de `resultados/NOME.txt` pra validação de que o código assembly está correto também
---

Exemplo de saída:
```
✅ teste1 passou! Resultado: 10
❌ teste2 falhou! Esperado: 42, Obtido: 0
```
---

### 🔹 make clean

Remove todos os arquivos dentro da pasta `bin/`.
---

## 📄 Formato de entrada YZY

Exemplo de um programa `.yzy`:

```yzy
a = 5;
b = 3;
c = 0;
{
  while (a > 0) {
    if (b < 5) {
      c = c + b;
      b = b + 1;
    } else {
      c = c + 1;
    }
    a = a - 1;
  }
  return c;
}
```

Esse código soma valores de `b` enquanto `a` for maior que 0.

---

## 🧠 Observações

- Todas as variáveis são inteiras (`int64`)
- A linguagem ainda não possui funções
- A comparação de valores retorna `1` pra true ou `0` pra false