# Suporte a Arquivos Excel (.xlsx)

## 📋 Visão Geral

O sistema agora suporta a leitura de arquivos Excel (.xlsx) contendo dados de produtos e revendedores, além do formato tradicional de arquivos TXT.

## 🚀 Como Usar

### Comando Básico

```bash
./bin/cargaparcial --excel seu_arquivo.xlsx
```

ou usando a forma curta:

```bash
./bin/cargaparcial -e seu_arquivo.xlsx
```

### Com Opções Adicionais

```bash
./bin/cargaparcial -e dados.xlsx -o resultado.json -w 16
```

## 📊 Formato do Arquivo Excel

O arquivo Excel deve conter as seguintes colunas:

- **IMBLOJA**: Código IBM da loja/revendedor
- **CODIGOBARRAS**: Código de barras do produto (EAN)

### Exemplo de Estrutura

| IMBLOJA    | CODIGOBARRAS  |
| ---------- | ------------- |
| 0001002154 | 7896050201756 |
| 0001002154 | 7898080070050 |
| 0001006393 | 070330717534  |
| 0001006393 | 0735202909010 |

### Características

✅ A primeira linha deve conter o cabeçalho com os nomes das colunas  
✅ Os nomes das colunas não são case-sensitive (IMBLOJA, imbloja, ImBLoJa são aceitos)  
✅ As colunas podem estar em qualquer ordem  
✅ Linhas vazias são automaticamente ignoradas  
✅ O sistema extrai todos os códigos IBM e produtos únicos e processa todas as combinações

## 🔄 Comparação: TXT vs Excel

### Modo TXT (Tradicional)

```bash
# Requer dois arquivos separados
./bin/cargaparcial -i ibm.txt -c codigo.txt
```

**Vantagens:**

- Simples e direto
- Fácil de editar manualmente
- Processa todas as combinações de IBM × Produtos

### Modo Excel (Novo)

```bash
# Um único arquivo com tudo
./bin/cargaparcial -e dados.xlsx
```

**Vantagens:**

- Dados organizados em uma única planilha
- Fácil de exportar de outros sistemas
- Suporta grandes volumes de dados
- Formato familiar para usuários de negócio

## 📝 Exemplo Prático

### 1. Criar um arquivo Excel de exemplo

Execute o script fornecido:

```bash
python3 criar_excel_exemplo.py
```

Isso criará um arquivo `dados_exemplo.xlsx` com dados de teste.

### 2. Processar o arquivo

```bash
./bin/cargaparcial -e dados_exemplo.xlsx -o resultado.json
```

### 3. Verificar o resultado

O arquivo `resultado.json` conterá:

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

## 🛠️ Flags Disponíveis

| Flag        | Forma Curta | Descrição                                                |
| ----------- | ----------- | -------------------------------------------------------- |
| `--excel`   | `-e`        | Arquivo Excel (.xlsx) com colunas IMBLOJA e CODIGOBARRAS |
| `--output`  | `-o`        | Arquivo de saída com resultados (padrão: resultado.json) |
| `--workers` | `-w`        | Número de workers paralelos (0 = auto)                   |

## 📚 Documentação Completa

Para mais detalhes, consulte:

- [docs/CLI_USAGE.md](docs/CLI_USAGE.md) - Guia completo de uso da CLI
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) - Arquitetura do sistema
- [docs/PARALLEL_PROCESSING.md](docs/PARALLEL_PROCESSING.md) - Processamento paralelo

## 🔧 Implementação Técnica

### Arquivos Criados/Modificados

- **`infrastructure/file/xlsx_reader.go`**: Leitor de arquivos XLSX
- **`cmd/api/main.go`**: Integração da flag --excel
- **`go.mod`**: Dependência github.com/xuri/excelize/v2

### Biblioteca Utilizada

- [excelize](https://github.com/xuri/excelize) - Biblioteca Go para leitura/escrita de arquivos Excel

## ❓ Troubleshooting

### Erro: "coluna IMBLOJA não encontrada"

**Solução**: Verifique se a primeira linha do arquivo Excel contém o cabeçalho com os nomes corretos das colunas.

### Erro: "arquivo vazio"

**Solução**: Certifique-se de que o arquivo Excel contém dados além do cabeçalho.

### Erro ao abrir arquivo

**Solução**: Verifique se:

1. O arquivo tem extensão .xlsx
2. O arquivo não está corrompido
3. Você tem permissão de leitura no arquivo

## 📞 Suporte

Para mais informações ou problemas, consulte a documentação completa ou entre em contato com a equipe de desenvolvimento.
