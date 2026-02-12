# 🔄 Guia de Migração e Deploy

## ✅ Verificação Pré-Deploy

### 1. Build e Teste
```bash
# Build
make build

# Teste com arquivo pequeno primeiro
./cargaparcial --ibm ibm_test.txt --codigo codigo_test.txt

# Verifique os logs
# Deve aparecer:
# 🚀 Fazendo batch insert de X ProductDealers
# ⚡ Progresso: X itens | X items/seg
```

### 2. Compatibilidade
✅ **Nenhuma mudança de schema necessária**  
✅ **Nenhuma mudança nas SPs necessária**  
✅ **100% compatível com código anterior**  
✅ **Zero breaking changes**

---

## 🚀 Deploy em Produção

### Passo 1: Backup
```bash
# Backup do binário atual
cp /caminho/prod/cargaparcial /caminho/prod/cargaparcial.backup

# Backup da config
cp /caminho/prod/config /caminho/prod/config.backup
```

### Passo 2: Deploy do Novo Binário
```bash
# Build do novo código
make build

# Copy para produção
scp cargaparcial usuario@servidor:/caminho/prod/

# SSH no servidor
ssh usuario@servidor

# Dar permissão de execução
chmod +x /caminho/prod/cargaparcial
```

### Passo 3: Teste em Produção (Dry Run)
```bash
# Execute com volume pequeno primeiro
./cargaparcial --ibm ibm_test.txt --codigo codigo_test.txt

# Monitore os logs
tail -f logs/cargaparcial.log

# Verifique métricas:
# - Items/seg deve estar entre 3000-5000
# - Batch inserts devem aparecer nos logs
# - SP Stats deve mostrar média < 15ms
```

### Passo 4: Rollout Completo
```bash
# Se tudo OK, execute com volume completo
./cargaparcial --ibm ibm.txt --codigo codigo.txt
```

---

## 📊 Monitoramento Pós-Deploy

### Métricas para Observar

#### 1. Throughput
```bash
# Nos logs, procure por:
⚡ Progresso: 10000 itens | 3500 items/seg | Tempo: 2.9s

# Esperado: 3000-5000 items/seg
# Se abaixo: investigar banco de dados (CPU, I/O)
```

#### 2. Batch Inserts
```bash
# Deve aparecer a cada ~100 items:
🚀 Fazendo batch insert de 100 ProductDealers

# Se não aparecer: problema com acumulação
```

#### 3. Stored Procedure
```bash
# A cada 1000 chamadas:
📊 SP Stats: 5000 chamadas | Média: 12.45ms | Erros: 0

# Esperado: Média < 15ms, Erros = 0
```

#### 4. Database
```sql
-- Monitor de conexões ativas
SELECT COUNT(*) FROM v$session WHERE username = 'SEU_USER';
-- Esperado: ~20-50 conexões (pool configurado para 100)

-- Monitor de locks
SELECT * FROM v$lock WHERE type = 'TX';
-- Esperado: Poucos locks, sem deadlocks

-- Monitor de CPU
SELECT value FROM v$sysmetric 
WHERE metric_name = 'Database CPU Time Ratio';
-- Esperado: 40-60% (antes era 70-85%)
```

---

## 🔍 Troubleshooting

### Problema 1: Performance não melhorou

**Sintomas:**
- Items/seg ainda em 1500-2000
- Batch inserts não aparecem nos logs

**Diagnóstico:**
```bash
# Verifique se está usando a versão nova
./cargaparcial --version

# Verifique os logs em detalhes
grep "batch insert" logs/cargaparcial.log
grep "Prepared Statement" logs/cargaparcial.log
```

**Solução:**
- Rebuild com `make clean && make build`
- Verifique se todos os arquivos foram atualizados

---

### Problema 2: Erro "too many bind variables"

**Sintomas:**
```
Erro ao criar batch de ProductDealers: ORA-01745: invalid host/bind variable name
```

**Solução:**
Reduzir tamanho do batch em `usecase/process_products_usecase.go`:
```go
batchSize: 50,  // Reduzir de 100 para 50
```

---

### Problema 3: Deadlocks no banco

**Sintomas:**
```sql
ORA-00060: deadlock detected while waiting for resource
```

**Solução:**
1. Verificar índices nas tabelas:
```sql
-- ProdutoRevendedor deve ter índice único em (IdProduto, IdRevendedor)
CREATE UNIQUE INDEX idx_produto_revendedor 
ON ProdutoRevendedor(IdProduto, IdRevendedor);
```

2. Reduzir número de workers:
```go
maxWorkers := runtime.NumCPU()  // Ao invés de NumCPU() * 2
```

---

### Problema 4: Consumo alto de memória

**Sintomas:**
- OOM (Out of Memory)
- Processo morto pelo sistema

**Solução:**
Reduzir capacidade do batch em `usecase/process_products_usecase.go`:
```go
batchProductDealers: make([]*entities.ProductDealer, 0, 100),  // Reduzir de 500 para 100
```

---

## 🎛️ Tuning de Performance

### Se performance ainda não é ideal:

#### 1. Aumentar Pool de Conexões
`infrastructure/database/connection.go`:
```go
db.SetMaxOpenConns(150)    // Aumentar de 100 para 150
db.SetMaxIdleConns(30)     // Aumentar de 20 para 30
```

#### 2. Aumentar Tamanho do Batch
`usecase/process_products_usecase.go`:
```go
batchSize: 200,  // Aumentar de 100 para 200
```

**⚠️ Atenção**: Oracle tem limite de ~1000 bind variables  
200 items × 3 campos = 600 binds (OK)  
400 items × 3 campos = 1200 binds (ERRO)

#### 3. Aumentar Workers
`usecase/process_products_usecase.go`:
```go
maxWorkers := runtime.NumCPU() * 3  // Aumentar de 2x para 3x
```

**⚠️ Atenção**: Mais workers = mais conexões ao banco

#### 4. Otimizar Banco de Dados

```sql
-- Gather statistics (Oracle)
EXEC DBMS_STATS.GATHER_TABLE_STATS('SCHEMA', 'ProdutoRevendedor');
EXEC DBMS_STATS.GATHER_TABLE_STATS('SCHEMA', 'IntegracaoProdutoStaging');

-- Verificar índices
SELECT * FROM user_indexes WHERE table_name IN ('PRODUTOREVENDEDOR', 'PRODUTO', 'REVENDEDOR');

-- Criar índice se não existir
CREATE INDEX idx_produto_ean ON EmbalagemProduto(CODIGOBARRAS);
CREATE INDEX idx_revendedor_ibm ON Revendedor(CodigoIBM);
```

---

## 📈 Benchmark Comparativo

### Antes do Deploy
```bash
# Execute com versão antiga e anote métricas
./cargaparcial.backup --ibm ibm_benchmark.txt --codigo codigo_benchmark.txt

# Anote:
# - Tempo total: _____
# - Items/seg: _____
# - CPU DB: _____
```

### Depois do Deploy
```bash
# Execute com versão nova
./cargaparcial --ibm ibm_benchmark.txt --codigo codigo_benchmark.txt

# Compare:
# - Tempo total: _____ (esperado: 2-3x menor)
# - Items/seg: _____ (esperado: 2-3x maior)
# - CPU DB: _____ (esperado: 30-50% menor)
```

---

## 🔐 Rollback (se necessário)

### Se algo der errado:

```bash
# Parar processo atual
pkill -9 cargaparcial

# Restaurar binário anterior
cp /caminho/prod/cargaparcial.backup /caminho/prod/cargaparcial

# Restaurar config
cp /caminho/prod/config.backup /caminho/prod/config

# Executar versão anterior
./cargaparcial --ibm ibm.txt --codigo codigo.txt
```

**Importante**: Não há alterações no banco de dados, então rollback é seguro!

---

## ✅ Checklist de Deploy

- [ ] Build executado com sucesso
- [ ] Testes unitários passando
- [ ] Teste com arquivo pequeno OK
- [ ] Backup do binário atual feito
- [ ] Backup da config feita
- [ ] Deploy do novo binário feito
- [ ] Permissões corretas (chmod +x)
- [ ] Dry run em produção OK
- [ ] Métricas de baseline capturadas
- [ ] Monitoramento configurado
- [ ] Logs sendo capturados
- [ ] Plano de rollback documentado
- [ ] Equipe notificada

---

## 📞 Suporte

### Logs para Debugging
```bash
# Aumentar verbosidade (se necessário)
export LOG_LEVEL=DEBUG
./cargaparcial --ibm ibm.txt --codigo codigo.txt

# Capturar logs detalhados
./cargaparcial --ibm ibm.txt --codigo codigo.txt 2>&1 | tee cargaparcial_debug.log
```

### Informações para Reportar Issues

Se encontrar problemas, colete:
1. Versão do Go: `go version`
2. Versão do Oracle: `SELECT * FROM v$version`
3. Logs completos
4. Métricas de CPU/Memória do servidor
5. Número de produtos processados
6. Tamanho dos arquivos de entrada

---

## 🎉 Sucesso!

Se os logs mostrarem:
```
✅ 🚀 Fazendo batch insert de 100 ProductDealers
✅ ⚡ Progresso: 10000 itens | 3500 items/seg | Tempo: 2.9s
✅ 📊 SP Stats: 10000 chamadas | Média: 12.45ms | Erros: 0
```

**Parabéns! Deploy foi um sucesso! 🎊**

Performance esperada: **2-3x mais rápida** que a versão anterior.
