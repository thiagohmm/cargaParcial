# Arquitetura do Sistema - Carga Parcial

## 📐 Visão Geral

Este documento descreve a arquitetura do sistema de processamento de produtos e revendedores, implementado seguindo os princípios de **Clean Architecture**.

## 🎯 Princípios da Clean Architecture

### 1. Independência de Frameworks

- O código de negócio não depende de frameworks específicos
- Frameworks são ferramentas, não arquitetura

### 2. Testabilidade

- Regras de negócio podem ser testadas sem UI, banco de dados ou servidor web

### 3. Independência de UI

- A UI pode mudar facilmente sem alterar o resto do sistema

### 4. Independência de Banco de Dados

- Regras de negócio não estão vinculadas ao banco de dados

### 5. Independência de Agentes Externos

- Regras de negócio não sabem nada sobre o mundo externo

## 🏗️ Camadas da Arquitetura

```
┌─────────────────────────────────────────────────────────┐
│                    Infrastructure                        │
│  (HTTP Handlers, Database, Queue, External Services)    │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                      Use Cases                           │
│         (Application Business Rules)                     │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                       Domain                             │
│    (Entities, Repository Interfaces, Service Interfaces) │
└─────────────────────────────────────────────────────────┘
```

### Domain Layer (Camada de Domínio)

**Responsabilidade**: Contém as regras de negócio fundamentais

**Componentes**:

- **Entities**: Objetos de negócio puros
  - `Product`: Representa um produto
  - `Dealer`: Representa um revendedor
  - `ProductDealer`: Relação produto-revendedor
  - `ProductIntegrationStaging`: Staging de integração

- **Repository Interfaces**: Contratos de acesso a dados
  - `DealerRepository`
  - `ProductRepository`
  - `ProductDealerRepository`
  - `ProductIntegrationStagingRepository`

- **Service Interfaces**: Contratos de serviços externos
  - `QueueService`

**Regras**:

- ✅ Não depende de nenhuma outra camada
- ✅ Contém apenas lógica de negócio pura
- ✅ Sem dependências externas

### Use Case Layer (Camada de Aplicação)

**Responsabilidade**: Orquestra o fluxo de dados e aplica regras de negócio específicas da aplicação

**Componentes**:

- **Use Cases**: Casos de uso da aplicação
  - `ProcessProductsUseCase`: Processa produtos para revendedores

- **DTOs**: Objetos de transferência de dados
  - `ProcessProductsInput`
  - `ProcessProductsOutput`
  - `ProductResultDTO`

**Regras**:

- ✅ Depende apenas da camada Domain
- ✅ Implementa lógica de aplicação
- ✅ Coordena entidades e repositórios

### Infrastructure Layer (Camada de Infraestrutura)

**Responsabilidade**: Implementa detalhes técnicos e integrações externas

**Componentes**:

1. **Database**:
   - `connection.go`: Gerenciamento de conexão com MySQL

2. **Repository Implementations**:
   - `DealerRepositoryImpl`
   - `ProductRepositoryImpl`
   - `ProductDealerRepositoryImpl`
   - `ProductIntegrationStagingRepositoryImpl`

3. **Queue**:
   - `QueueServiceImpl`: Implementação do serviço de fila

4. **HTTP**:
   - `ProcessProductsHandler`: Handler HTTP para processar produtos

**Regras**:

- ✅ Implementa interfaces definidas no Domain
- ✅ Contém detalhes de implementação
- ✅ Pode depender de frameworks e bibliotecas externas

## 🔄 Fluxo de Dados

### Processamento de Produtos

```
1. HTTP Request
   ↓
2. ProcessProductsHandler (Infrastructure)
   ↓
3. ProcessProductsUseCase (Use Case)
   ↓
4. Repositories (Domain Interfaces → Infrastructure Implementation)
   ↓
5. Database / External Services
   ↓
6. Response
```

### Exemplo Detalhado

```go
// 1. Request chega no Handler
POST /api/process-products
{
  "IBM": ["IBM001"],
  "codigo": ["7891234567890"]
}

// 2. Handler valida e chama Use Case
handler.Handle(w, r)
  → useCase.Execute(input)

// 3. Use Case orquestra a lógica
useCase.Execute(input)
  → dealerRepo.GetByIBM("IBM001")
  → productRepo.GetByEAN("7891234567890")
  → productDealerRepo.Exists(productID, dealerID)
  → productDealerRepo.Create(...)
  → productRepo.SaveIntegrationStaging(...)
  → queueService.Send("mover")

// 4. Response é retornada
{
  "arrayOk": [...],
  "arrayFail": [...]
}
```

## 🔌 Dependency Injection

O sistema utiliza **Dependency Injection** para manter o baixo acoplamento:

```go
// main.go
func main() {
    // 1. Criar dependências externas
    db := database.NewConnection(config)

    // 2. Criar repositórios (implementações)
    dealerRepo := repository.NewDealerRepository(db)
    productRepo := repository.NewProductRepository(db)

    // 3. Criar serviços
    queueService := queue.NewQueueService()

    // 4. Injetar no Use Case
    useCase := usecase.NewProcessProductsUseCase(
        dealerRepo,
        productRepo,
        productDealerRepo,
        productIntegrationRepo,
        queueService,
    )

    // 5. Injetar no Handler
    handler := handler.NewProcessProductsHandler(useCase)
}
```

## 🧪 Testabilidade

### Vantagens da Arquitetura para Testes

1. **Mocks Fáceis**: Interfaces permitem criar mocks facilmente
2. **Testes Isolados**: Cada camada pode ser testada independentemente
3. **Testes Rápidos**: Use cases podem ser testados sem banco de dados

### Exemplo de Teste

```go
// Mock do Repository
type MockDealerRepository struct {
    mock.Mock
}

func (m *MockDealerRepository) GetByIBM(ibm string) (*entities.Dealer, error) {
    args := m.Called(ibm)
    return args.Get(0).(*entities.Dealer), args.Error(1)
}

// Teste do Use Case
func TestProcessProductsUseCase(t *testing.T) {
    // Arrange
    mockDealerRepo := new(MockDealerRepository)
    mockDealerRepo.On("GetByIBM", "IBM001").Return(&entities.Dealer{ID: 1}, nil)

    useCase := usecase.NewProcessProductsUseCase(mockDealerRepo, ...)

    // Act
    result, err := useCase.Execute(input)
este do Use Case
func TestProcessProductsUseCase(t *testing.T) {
    // Arrange
    mockDealerRepo := new(MockDealerRepository)
    mockDealerRepo.On("GetByIBM", "IBM001").Return(&entities.Dealer{ID: 1}, nil)
    
    useCase := usecase.NewProcessProductsUseCase(mockDealerRepo, ...)
    
    // Act
    result, err := useCase.Execute(input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

## 📊 Diagrama de Dependências

```
cmd/api/main.go
    ↓
infrastructure/http/handler
    ↓
usecase/process_products_usecase
    ↓
domain/repositories (interfaces)
    ↑
infrastructure/repository (implementations)
```

## 🎨 Padrões de Design Utilizados

1. **Repository Pattern**: Abstração de acesso a dados
2. **Dependency Injection**: Inversão de controle
3. **DTO Pattern**: Transferência de dados entre camadas
4. **Use Case Pattern**: Encapsulamento de lógica de negócio
5. **Interface Segregation**: Interfaces específicas e coesas

## 🔐 Princípios SOLID

- **S**ingle Responsibility: Cada componente tem uma única responsabilidade
- **O**pen/Closed: Aberto para extensão, fechado para modificação
- **L**iskov Substitution: Implementações podem substituir interfaces
- **I**nterface Segregation: Interfaces específicas e focadas
- **D**ependency Inversion: Dependa de abstrações, não de implementações

## 🚀 Benefícios da Arquitetura

1. **Manutenibilidade**: Código organizado e fácil de entender
2. **Escalabilidade**: Fácil adicionar novos recursos
3. **Testabilidade**: Testes isolados e rápidos
4. **Flexibilidade**: Fácil trocar implementações
5. **Independência**: Camadas desacopladas

## 📝 Convenções de Código

1. **Nomenclatura**:
   - Interfaces: `XxxRepository`, `XxxService`
   - Implementações: `XxxRepositoryImpl`, `XxxServiceImpl`
   - DTOs: `XxxInput`, `XxxOutput`, `XxxDTO`

2. **Organização**:
   - Um arquivo por tipo/interface
   - Pacotes organizados por responsabilidade
   - Testes ao lado do código

3. **Comentários**:
   - Documentar interfaces públicas
   - Explicar lógica complexa
   - Usar godoc format
