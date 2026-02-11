# TODO - Implementação de Suporte a Arquivos XLSX

## Tarefas Concluídas

- [x] 1. Adicionar dependência excelize ao go.mod
- [x] 2. Criar infrastructure/file/xlsx_reader.go com função de leitura
- [x] 3. Atualizar cmd/api/main.go para suportar flag --excel
- [x] 4. Atualizar docs/CLI_USAGE.md com exemplos de uso XLSX
- [x] 5. Compilar o projeto com sucesso

## Arquivos Criados/Modificados

### Criados

- `infrastructure/file/xlsx_reader.go` - Leitor de arquivos XLSX com suporte às colunas IMBLOJA e CODIGOBARRAS
- `criar_excel_exemplo.py` - Script para criar arquivo Excel de exemplo para testes

### Modificados

- `go.mod` / `go.sum` - Adicionada dependência github.com/xuri/excelize/v2
- `cmd/api/main.go` - Adicionada flag --excel (-e) e lógica para processar arquivos XLSX
- `docs/CLI_USAGE.md` - Documentação completa com exemplos de uso do arquivo Excel

## Como Usar

### Modo Tradicional (Arquivos TXT)

```bash
./bin/cargaparcial -i ibm.txt -c codigo.txt -o resultado.json
```

### Novo Modo (Arquivo Excel)

```bash
./bin/cargaparcial --excel dados.xlsx -o resultado.json
```

ou

```bash
./bin/cargaparcial -e dados.xlsx -o resultado.json
```

## Formato do Arquivo Excel

O arquivo deve ter as seguintes colunas (ordem não importa):

- **IMBLOJA**: Código IBM da loja/revendedor
- **CODIGOBARRAS**: Código de barras do produto (EAN)

Exemplo:

| IMBLOJA | CODIGOBARRAS |
|------------|---------------|
| 0001002154 | 7896050201756 |
| 0001002154 | 7898080070050 |
| 0001006393 | 070330717534 |

## Status: ✅ Implementação Completa
