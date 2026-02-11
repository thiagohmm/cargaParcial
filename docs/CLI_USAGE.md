# Guia de Uso da CLI

## Visão Geral

O sistema utiliza **Cobra CLI** para fornecer uma interface de linha de comando flexível e intuitiva.

## Instalação

```bash
go build -o bin/cargaparcial cmd/api/main.go
```

## Uso Básico

### Executar com Arquivos Padrão (TXT)

```bash
./bin/cargaparcial
```

Isso irá:

- Ler `ibm.txt` (códigos IBM)
- Ler `codigo.txt` (códigos de produtos/EAN)
- Gerar `resultado.json` (arquivo de saída)
- Usar número automático de workers (baseado em CPUs)

### Executar com Arquivo Excel

```bash
./bin/cargaparcial --excel dados.xlsx
```

Isso irá:

- Ler `dados.xlsx` (arquivo Excel com colunas IMBLOJA e CODIGOBARRAS)
- Gerar `resultado.json` (arquivo de saída)
- Usar número automático de workers (baseado em CPUs)

## Flags Disponíveis

### Especificar Arquivos de Entrada

#### Arquivos TXT (Modo Tradicional)

```bash
# Arquivo de códigos IBM
./bin/cargaparcial --ibm meus_ibms.txt
./bin/cargaparcial -i meus_ibms.txt

# Arquivo de códigos de produtos
./bin/cargaparcial --codigo meus_codigos.txt
./bin/cargaparcial -c meus_codigos.txt

# Ambos
./bin/cargaparcial -i ibm_custom.txt -c codigo_custom.txt
```

#### Arquivo Excel (Modo XLSX)

```bash
# Arquivo Excel com colunas IMBLOJA e CODIGOBARRAS
./bin/cargaparcial --excel dados.xlsx
./bin/cargaparcial -e dados.xlsx

# Com arquivo de saída personalizado
./bin/cargaparcial -e dados.xlsx -o resultado_custom.json

# Com número específico de workers
./bin/cargaparcial -e dados.xlsx -w 16
```

**Nota:** Quando usar a flag `--excel`, as flags `--ibm` e `--codigo` são ignoradas.

### Especificar Arquivo de Saída

```bash
./bin/cargaparcial --output meu_resultado.json
./bin/cargaparcial -o meu_resultado.json
```

### Configurar Workers Paralelos

```bash
# Usar 10 workers
./bin/cargaparcial --workers 10
./bin/cargaparcial -w 10

# Usar 1 worker (processamento sequencial)
./bin/cargaparcial -w 1

# Usar 0 ou omitir (automático, baseado em CPUs)
./bin/cargaparcial -w 0
./bin/cargaparcial
```

## Exemplos Completos

### Exemplo 1: Processamento Padrão (TXT)

```bash
./bin/cargaparcial
```

**Saída esperada:**

```
=== Carga Parcial - Processador de Produtos ===
Arquivo IBM: ibm.txt
Arquivo Código: codigo.txt
Arquivo Saída: resultado.json
✓ Conexão com banco de dados estabelecida
Lendo arquivo: ibm.txt
✓ Lidos 3 códigos IBM
Lendo arquivo: codigo.txt
✓ Lidos 5 códigos de produto
Total de combinações a processar: 15
Iniciando processamento paralelo...
Iniciando processamento paralelo com 8 workers
Worker 1 finalizado: processou 2 itens no total
Worker 2 finalizado: processou 2 itens no total
...
=== Processamento Concluído ===
✓ Sucessos: 12
✗ Falhas: 3
Taxa de sucesso: 80.00%
✓ Resultado salvo com sucesso em resultado.json
=== Processo Finalizado ===
```

### Exemplo 2: Processamento com Excel

```bash
./bin/cargaparcial --excel dados_janeiro.xlsx
```

**Saída esperada:**

```
=== Carga Parcial - Processador de Produtos ===
Arquivo Excel: dados_janeiro.xlsx
Arquivo Saída: resultado.json
✓ Conexão com banco de dados estabelecida
Lendo arquivo Excel: dados_janeiro.xlsx
✓ Lidos 5 códigos IBM únicos
✓ Lidos 120 códigos de produto únicos
Total de combinações a processar: 600
Iniciando processamento paralelo...
Iniciando processamento paralelo com 8 workers
...
=== Processamento Concluído ===
✓ Sucessos: 580
✗ Falhas: 20
Taxa de sucesso: 96.67%
✓ Resultado salvo com sucesso em resultado.json
=== Processo Finalizado ===
```

### Exemplo 3: Arquivos Personalizados (TXT)

```bash
./bin/cargaparcial \
  --ibm dados/ibm_janeiro.txt \
  --codigo dados/produtos_novos.txt \
  --output resultados/janeiro_2024.json \
  --workers 16
```

### Exemplo 4: Excel com Configurações Personalizadas

```bash
./bin/cargaparcial \
  --excel dados/carga_completa.xlsx \
  --output resultados/resultado_completo.json \
  --workers 32
```

### Exemplo 5: Processamento Sequencial (Debug)

```bash
# Com arquivos TXT
./bin/cargaparcial -w 1

# Com arquivo Excel
./bin/cargaparcial -e dados.xlsx -w 1
```

Útil para debugging, pois processa um item por vez.

### Exemplo 6: Alta Performance

```bash
# Com arquivos TXT
./bin/cargaparcial -w 32

# Com arquivo Excel
./bin/cargaparcial -e dados_grandes.xlsx -w 32
```

Para servidores com muitos núcleos e arquivos muito grandes.

## Tabela de Flags

| Flag        | Forma Curta | Valor Padrão     | Descrição                                                     |
| ----------- | ----------- | ---------------- | ------------------------------------------------------------- |
| `--ibm`     | `-i`        | `ibm.txt`        | Arquivo com códigos IBM (um por linha)                        |
| `--codigo`  | `-c`        | `codigo.txt`     | Arquivo com códigos de produtos/EAN (um por linha)            |
| `--excel`   | `-e`        | -                | Arquivo Excel (.xlsx) com colunas IMBLOJA e CODIGOBARRAS      |
| `--output`  | `-o`        | `resultado.json` | Arquivo de saída com resultados JSON                          |
| `--workers` | `-w`        | `0` (auto)       | Número de workers paralelos (0 = baseado em CPUs disponíveis) |
| `--help`    | `-h`        | -                | Exibe ajuda e sai                                             |

## Formato dos Arquivos de Entrada

### Arquivos TXT

#### ibm.txt

```
IBM001
IBM002
IBM003
```

- Um código IBM por linha
- Linhas vazias são ignoradas
- Linhas começando com `#` são tratadas como comentários

#### codigo.txt

```
7891234567890
7891234567891
7891234567892
```

- Um código EAN/produto por linha
- Linhas vazias são ignoradas
- Linhas começando com `#` são tratadas como comentários

#### Exemplo com Comentários

```
# Códigos IBM - Janeiro 2024
IBM001
IBM002

# Novos revendedores
IBM003
IBM004
```

### Arquivo Excel (.xlsx)

O arquivo Excel deve conter as seguintes colunas (a ordem não importa):

- **IMBLOJA**: Código IBM da loja/revendedor
- **CODIGOBARRAS**: Código de barras do produto (EAN)

#### Exemplo de Estrutura

| IMBLOJA    | CODIGOBARRAS  |
| ---------- | ------------- |
| 0001002154 | 7896050201756 |
| 0001002154 | 7898080070050 |
| 0001006393 | 070330717534  |
| 0001006393 | 0735202909010 |

**Características:**

- A primeira linha deve conter o cabeçalho com os nomes das colunas
- Os nomes das colunas não são case-sensitive (IMBLOJA, imbloja, ImBLoJa são aceitos)
- As colunas podem estar em qualquer ordem
- Linhas vazias são ignoradas
- O sistema extrai todos os códigos IBM e produtos únicos e processa todas as combinações

## Formato do Arquivo de Saída

### resultado.json

```json
{
  "arrayOk": [
    {
      "IdRevendedor": 1,
      "IdProduto": 100,
      "Status": "ok"
    }
  ],
  "arrayFail": [
    {
      "IdRevendedor": 2,
      "IdProduto": null,
      "EAN": "7891234567891",
      "Status": "fail",
      "Motivo": "Produto não encontrado pelo EAN"
    }
  ]
}
```

## Ajuda Integrada

```bash
./bin/cargaparcial --help
```

**Saída:**

```
Sistema de processamento paralelo de produtos e revendedores.
Lê códigos IBM e códigos de produtos de arquivos de entrada,
processa em paralelo e gera arquivo de resultado.

Usage:
  cargaparcial [flags]

Flags:
  -c, --codigo string   Arquivo com códigos de produtos/EAN (um por linha) (default "codigo.txt")
  -e, --excel string    Arquivo Excel (.xlsx) com colunas IMBLOJA e CODIGOBARRAS
  -h, --help            help for cargaparcial
  -i, --ibm string      Arquivo com códigos IBM (um por linha) (default "ibm.txt")
  -o, --output string   Arquivo de saída com resultados (default "resultado.json")
  -w, --workers int     Número de workers paralelos (0 = auto, baseado em CPUs)
```

## Tratamento de Erros

### Arquivo Não Encontrado

```bash
$ ./bin/cargaparcial -i arquivo_inexistente.txt
Lendo arquivo: arquivo_inexistente.txt
Erro ao ler arquivo arquivo_inexistente.txt: open arquivo_inexistente.txt: no such file or directory
```

### Erro de Conexão com Banco

```bash
Erro ao conectar ao banco de dados: ORA-12154: TNS:could not resolve the connect identifier specified
```

### Erro ao Salvar Resultado

```bash
Erro ao salvar resultado.json: permission denied
```

## Dicas de Performance

### Escolher Número de Workers

```bash
# Para I/O intensivo (muitas queries ao banco)
./bin/cargaparcial -w $(nproc --all)  # Número de CPUs

# Para CPU intensivo
./bin/cargaparcial -w $(($(nproc --all) * 2))  # 2x número de CPUs

# Para arquivos pequenos
./bin/cargaparcial -w 4

# Para arquivos gigantes
./bin/cargaparcial -w 32
```

### Monitorar Progresso

Os logs mostram o progresso em tempo real:

```
Worker 1: processou 100 itens
Worker 2: processou 100 itens
Worker 3: processou 100 itens
...
```

### Redirecionar Logs

```bash
# Salvar logs em arquivo
./bin/cargaparcial 2>&1 | tee processamento.log

# Apenas erros
./bin/cargaparcial 2> erros.log

# Sem logs (apenas resultado final)
./bin/cargaparcial > /dev/null 2>&1
```

## Integração com Scripts

### Bash Script - Arquivos TXT

```bash
#!/bin/bash

# Processar múltiplos arquivos TXT
for mes in janeiro fevereiro marco; do
  echo "Processando $mes..."
  ./bin/cargaparcial \
    -i "dados/${mes}_ibm.txt" \
    -c "dados/${mes}_codigo.txt" \
    -o "resultados/${mes}_resultado.json" \
    -w 16
done
```

### Bash Script - Arquivos Excel

```bash
#!/bin/bash

# Processar múltiplos arquivos Excel
for arquivo in dados/*.xlsx; do
  nome=$(basename "$arquivo" .xlsx)
  echo "Processando $nome..."
  ./bin/cargaparcial \
    -e "$arquivo" \
    -o "resultados/${nome}_resultado.json" \
    -w 16
done
```

### Makefile

```makefile
.PHONY: process
process:
 ./bin/cargaparcial -i ibm.txt -c codigo.txt -o resultado.json -w 8

.PHONY: process-excel
process-excel:
 ./bin/cargaparcial -e dados.xlsx -o resultado.json -w 8

.PHONY: process-prod
process-prod:
 ./bin/cargaparcial \
  -i /data/ibm_producao.txt \
  -c /data/codigo_producao.txt \
  -o /output/resultado_producao.json \
  -w 32

.PHONY: process-prod-excel
process-prod-excel:
 ./bin/cargaparcial \
  -e /data/producao.xlsx \
  -o /output/resultado_producao.json \
  -w 32
```

## Troubleshooting

### Problema: Processamento Muito Lento

**Solução**: Aumente o número de workers

```bash
./bin/cargaparcial -w 16
```

### Problema: Uso Excessivo de Memória

**Solução**: Reduza o número de workers

```bash
./bin/cargaparcial -w 4
```

### Problema: Arquivo de Saída Não Criado

**Verificar**:
cial -w 4

````

### Problema: Arquivo de Saída Não Criado

**Verificar**:

1. Permissões de escrita no diretório
2. Espaço em disco disponível
3. Caminho do arquivo de saída válido

```bash
# Verificar permissões
ls -la resultado.json

# Verificar espaço
df -h .
````

## Referências

- [Cobra CLI Documentation](https://github.com/spf13/cobra)
- [Processamento Paralelo](PARALLEL_PROCESSING.md)
- [Arquitetura do Sistema](ARCHITECTURE.md)
