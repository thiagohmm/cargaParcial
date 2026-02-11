# Processamento Paralelo

## Visão Geral

O sistema utiliza **goroutines** (threads leves do Go) para processar grandes volumes de dados de forma eficiente e paralela.

## Arquitetura de Paralelização

### Worker Pool Pattern

O sistema implementa o padrão **Worker Pool** com as seguintes características:

1. **Workers**: Goroutines que processam jobs do canal
2. **Job Channel**: Canal buffered que recebe os trabalhos a serem processados
3. **Result Channel**: Canal buffered que coleta os resultados
4. **WaitGroups**: Sincronização para aguardar conclusão de todos os workers

### Fluxo de Processamento

```
┌─────────────┐
│   Main      │
│  Goroutine  │
└──────┬──────┘
       │
       ├─────────────────────────────────────┐
       │                                     │
       ▼                                     ▼
┌─────────────┐                      ┌─────────────┐
│ Job Channel │                      │   Result    │
│  (buffered) │                      │   Channel   │
└──────┬──────┘                      └──────▲──────┘
       │                                     │
       ├──────┬──────┬──────┬──────┐        │
       │      │      │      │      │        │
       ▼      ▼      ▼      ▼      ▼        │
    ┌────┐ ┌────┐ ┌────┐ ┌────┐ ┌────┐     │
    │ W1 │ │ W2 │ │ W3 │ │ W4 │ │ WN │─────┤
    └────┘ └────┘ └────┘ └────┘ └────┘     │
                                            │
                                            ▼
                                    ┌───────────────┐
                                    │   Collector   │
                                    │   Goroutine   │
                                    └───────────────┘
```

## Configuração

### Número de Workers

Por padrão, o sistema utiliza `runtime.NumCPU()` workers (número de CPUs disponíveis).

Para configurar manualmente:

```go
processProductsUseCase := usecase.NewProcessProductsUseCase(...)
processProductsUseCase.SetMaxWorkers(10) // Define 10 workers
```

### Tamanho dos Buffers

- **Job Channel**: Buffer de 100 jobs
- **Result Channel**: Buffer de 100 resultados

Esses valores podem ser ajustados no código se necessário para otimizar o uso de memória.

## Características de Performance

### Vantagens

1. **Escalabilidade**: Aproveita todos os núcleos da CPU
2. **Throughput**: Processa múltiplos itens simultaneamente
3. **Eficiência**: Goroutines são leves (2KB de stack inicial)
4. **Resiliência**: Erros em um worker não afetam os outros

### Monitoramento

O sistema fornece logs de progresso:

```
Worker 1: processou 100 itens
Worker 2: processou 100 itens
Worker 3: processou 100 itens
...
Worker 1 finalizado: processou 523 itens no total
Worker 2 finalizado: processou 498 itens no total
```

### Métricas Finais

Ao final do processamento:

```
Processamento concluído: 5000 jobs processados
Sucessos: 4850, Falhas: 150
```

## Sincronização e Thread-Safety

### WaitGroups

- **Worker WaitGroup**: Aguarda todos os workers finalizarem
- **Result WaitGroup**: Aguarda coleta de todos os resultados

### Channels

- Channels são thread-safe por natureza em Go
- Fechamento adequado dos channels previne deadlocks

### Mutex (se necessário)

Para operações que requerem exclusão mútua, utilize `sync.Mutex`:

```go
var mu sync.Mutex
mu.Lock()
// operação crítica
mu.Unlock()
```

## Otimizações

### Para Arquivos Muito Grandes

1. **Aumentar Buffer dos Channels**:

   ```go
   jobs := make(chan JobInput, 1000)
   results := make(chan dto.ProductResultDTO, 1000)
   ```

2. **Ajustar Número de Workers**:

   ```go
   // Para I/O intensivo
   processProductsUseCase.SetMaxWorkers(runtime.NumCPU() * 2)

   // Para CPU intensivo
   processProductsUseCase.SetMaxWorkers(runtime.NumCPU())
   ```

3. **Processamento em Lotes**:
   - Considere processar arquivos em chunks se a memória for limitada

### Logs Reduzidos

Para evitar sobrecarga de logs em grandes volumes:

```go
// Log apenas a cada 1000 itens
if productID%1000 == 0 {
    fmt.Printf("Debug: %v, %d, %d\n", productIntegrationStaging, dealerID, productID)
}
```

## Exemplo de Uso

```go
// Criar use case
uc := usecase.NewProcessProductsUseCase(
    dealerRepo,
    productRepo,
    productDealerRepo,
    productIntegrationRepo,
    queueService,
)

// Configurar workers (opcional)
uc.SetMaxWorkers(8)

// Executar processamento paralelo
output, err := uc.Execute(input)
if err != nil {
    log.Fatal(err)
}

// Resultados
fmt.Printf("Sucessos: %d\n", len(output.SuccessList))
fmt.Printf("Falhas: %d\n", len(output.FailureList))
```

## Considerações de Banco de Dados

### Connection Pool

Certifique-se de configurar adequadamente o pool de conexões:

```go
db.SetMaxOpenConns(25)  // Máximo de conexões abertas
db.SetMaxIdleConns(5)   // Conexões idle no pool
db.SetConnMaxLifetime(5 * time.Minute)
```

### Transações

Para operações que requerem transações, considere:

```go
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

// operações...

return tx.Commit()
```

## Troubleshooting

### Deadlocks

Se o programa travar:


Se o programa travar:

1. Verifique se todos os channels estão sendo fechados
2. Confirme que WaitGroups estão balanceados (Add/Done)

### Uso Excessivo de Memória

1. Reduza o buffer dos channels
2. Diminua o número de workers
3. Processe em lotes menores

### Performance Baixa

1. Aumente o número de workers
2. Verifique gargalos no banco de dados
3. Profile com `pprof`:

   ```bash
   go test -cpuprofile=cpu.prof
   go tool pprof cpu.prof
   ```

## Referências

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Worker Pool Pattern](https://gobyexample.com/worker-pools)
