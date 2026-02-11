arqui# Carga Parcial - Clean Architecture

Este projeto implementa um sistema de processamento de produtos e revendedores seguindo os princípios de Clean Architecture.

## 📁 Estrutura do Projeto

```
.
├── cmd/
│   └── api/
│       └── main.go                 # Ponto de entrada da aplicação
├── domain/
│   ├── entities/                   # Entidades de negócio
│   │   └── product.go
│   ├── repositories/               # Interfaces dos repositórios
│   │   ├── dealer_repository.go
│   │   ├── product_repository.go
│   │   ├── product_dealer_repository.go
│   │   └── product_integration_staging_repository.go
│   └── services/                   # Interfaces dos serviços
│       └── queue_service.go
├── usecase/
│   ├── dto/                        # Data Transfer Objects
│   │   └── process_products_dto.go
│   └── process_products_usecase.go # Lógica de negócio
├── infrastructure/
│   ├── database/                   # Configuração do banco de dados
│   │   └── connection.go
│   ├── repository/                 # Implementações dos repositórios
│   │   ├── dealer_repository_impl.go
│   │   ├── product_repository_impl.go
│   │   ├── product_dealer_repository_impl.go
│   │   └── product_integration_staging_repository_impl.go
│   ├── queue/                      # Implementação do serviço de fila
│   │   └── queue_service_impl.go
│   └── http/
│       └── handler/                # Handlers HTTP
│           └── process_products_handler.go
├── go.mod
└── README.md
```

## 🏗️ Arquitetura

O projeto segue os princípios de **Clean Architecture**, dividido em camadas:

### 1. Domain (Camada de Domínio)

- **Entities**: Modelos de negócio puros
- **Repositories**: Interfaces que definem contratos de acesso a dados
- **Services**: Interfaces de serviços externos

### 2. Use Cases (Camada de Aplicação)

- Contém a lógica de negócio da aplicação
- Orquestra o fluxo de dados entre as camadas
- Independente de frameworks e detalhes de implementação

### 3. Infrastructure (Camada de Infraestrutura)

- **Database**: Configuração e conexão com banco de dados
- **Repository**: Implementações concretas dos repositórios
- **Queue**: Implementação do serviço de filas
- **HTTP**: Handlers e rotas HTTP

### 4. CMD (Camada de Interface)

- Ponto de entrada da aplicação
- Configuração e inicialização de dependências

## 🚀 Como Executar

### Pré-requisitos

- Go 1.25.3 ou superior
- Oracle Database
- Variáveis de ambiente configuradas

### Configuração

Crie um arquivo `.env` baseado no `config.example`:

```bash
cp config.example .env
```

Edite o arquivo `.env` com suas configurações:

```bash
# Database Configuration (Oracle)
DB_DIALECT=oracle
DB_USER=STAGE
DB_PASSWD=sua_senha
DB_SCHEMA=STAGE
DB_CONNECTSTRING=(description= (retry_count=20)(retry_delay=3)(address=(protocol=tcps)(port=1522)(host=seu_host))(connect_data=(service_name=seu_service_name))(security=(ssl_server_dn_match=no)))

# RabbitMQ Configuration
ENV_RABBITMQ=amqp://guest:guest@localhost:5672/
QUEUE_NAME=integracao
```

**Nota:** A string de conexão (`DB_CONNECTSTRING`) deve seguir o formato TNS do Oracle:

```
(description= (retry_count=20)(retry_delay=3)(address=(protocol=tcps)(port=PORT)(host=HOST))(connect_data=(service_name=SERVICE_NAME))(security=(ssl_server_dn_match=no)))
```

O sistema extrai automaticamente `host`, `port` e `service_name` da string TNS.

### Preparar Arquivos de Entrada

Crie os arquivos de entrada com os dados a serem processados:

```bash
# Copiar exemplos
cp ibm.txt.example ibm.txt
cp codigo.txt.example codigo.txt
```

Edite os arquivos com seus dados:

**ibm.txt** - Um código IBM por linha:

```
IBM001
IBM002
IBM003
```

**codigo.txt** - Um código EAN por linha:

```
7891234567890
7891234567891
7891234567892
```

### Instalação de Dependências

```bash
go mod download
```

### Compilação

```bash
go build -o bin/cargaparcial cmd/api/main.go
```

### Executar a Aplicação

```bash
./bin/cargaparcial
```

O programa irá:


O programa irá:

1. Ler os arquivos `ibm.txt` e `codigo.txt`
2. Processar cada combinação de IBM + Código
3. Salvar o resultado em `resultado.json`
4. Enviar mensagem "mover" para a fila "integracao" do RabbitMQ

## 📄 Arquivos de Saída

### resultado.json

Contém o resultado do processamento:

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

## 🗄️ Estrutura do Banco de Dados

O sistema espera as seguintes tabelas no Oracle:

- `Revendedor`: Armazena informações dos revendedores
  - `IdRevendedor` (NUMBER)
  - `IBM` (VARCHAR2)
- `Produto`: Armazena informações dos produtos
  - `IDPRODUTO` (NUMBER)
  - `EAN` (VARCHAR2)
- `ProdutoRevendedor`: Relacionamento entre produtos e revendedores
  - `IdProduto` (NUMBER)
  - `IdRevendedor` (NUMBER)
  - `StatusProdutoRevendedor` (NUMBER/BOOLEAN)`Produto`: Armazena informações dos produtos
  - `IDPRODUTO` (NUMBER)
  - `EAN` (VARCHAR2)
- `ProdutoRevendedor`: Relacionamento entre produtos e revendedores
  - `IdProduto` (NUMBER)
  - `IdRevendedor` (NUMBER)
  - `StatusProdutoRevendedor` (NUMBER/BOOLEAN)- `StatusProdutoRevendedor` (NUMBER/BOOLEAN)- `StatusProdutoRevendedor` (NUMBER/BOOLEAN)- `StatusProdutoRevendedor` (NUMBER/BOOLEAN)
- `ProdutoIntegracaoStaging`: Staging de integração de produtos
  - `IdProduto` (NUMBER)
  - `IdRevendedor` (NUMBER)
  - `DataCriacao` (DATE)
  - `DataAtualizacao` (DATE)

## 🧪 Testes

Para executar os testes:

```bash
go test ./...
```

## 📝 Princípios Aplicados

- **Dependency Inversion**: As camadas internas não dependem das externas
- **Single Responsibility**: Cada componente tem uma única responsabilidade
- **Open/Closed**: Aberto para extensão, fechado para modificação
- **Interface Segregation**: Interfaces específicas e coesas
- **Dependency Injection**: Dependências injetadas via construtor

## 🔄 Fluxo de Dados

1. **HTTP Handler** recebe a requisição
2. **Use Case** processa a lógica de negócio
3. **Repositories** acessam o banco de dados
4. **Queue Service** envia mensagens para fila
5. **HTTP Handler** retorna a resposta

## 🛠️ Tecnologias

- **Go 1.25.3**: Linguagem de programação
- **Oracle Database**: Banco de dados relacional
- **go-ora/v2**: Driver Oracle para Go
- **database/sql**: Driver SQL padrão do Go
- **net/http**: Servidor HTTP padrão do Go

## 📄 Licença

Este projeto é proprietário da Raizen.
