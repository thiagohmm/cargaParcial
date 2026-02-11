#!/bin/bash

# Script de Benchmark de Performance
# Testa diferentes números de workers para encontrar a configuração ideal

echo "🚀 Benchmark de Performance - Carga Parcial"
echo "==========================================="
echo ""

# Verificar se o binário existe
if [ ! -f "./bin/cargaparcial" ]; then
    echo "❌ Binário não encontrado. Execute 'make build' primeiro."
    exit 1
fi

# Verificar se o arquivo Excel existe
EXCEL_FILE="lojas_produtos.xlsx"
if [ ! -f "$EXCEL_FILE" ]; then
    echo "⚠️  Arquivo $EXCEL_FILE não encontrado."
    echo "Especifique o arquivo: $0 <arquivo.xlsx>"
    exit 1
fi

# Arquivo de entrada pode ser passado como argumento
if [ -n "$1" ]; then
    EXCEL_FILE="$1"
fi

echo "📁 Arquivo: $EXCEL_FILE"
echo ""

# Aumentar limite de arquivos
echo "📊 Configurando sistema..."
ulimit -n 65536
echo "✅ ulimit -n: $(ulimit -n)"
echo ""

# Array de configurações para testar
CONFIGS=(
    "24:padrão"
    "48:dobro"
    "96:recomendado"
    "150:turbo"
    "200:máximo"
)

RESULTS_FILE="benchmark_results_$(date +%Y%m%d_%H%M%S).txt"

echo "📝 Resultados serão salvos em: $RESULTS_FILE"
echo "" | tee "$RESULTS_FILE"
echo "=== BENCHMARK DE PERFORMANCE ===" | tee -a "$RESULTS_FILE"
echo "Data: $(date)" | tee -a "$RESULTS_FILE"
echo "Arquivo: $EXCEL_FILE" | tee -a "$RESULTS_FILE"
echo "CPUs: $(nproc)" | tee -a "$RESULTS_FILE"
echo "" | tee -a "$RESULTS_FILE"

# Executar testes
for config in "${CONFIGS[@]}"; do
    WORKERS="${config%%:*}"
    LABEL="${config##*:}"
    OUTPUT_FILE="resultado_${WORKERS}w.json"
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" | tee -a "$RESULTS_FILE"
    echo "🧪 Teste: $WORKERS workers ($LABEL)" | tee -a "$RESULTS_FILE"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" | tee -a "$RESULTS_FILE"
    echo "" | tee -a "$RESULTS_FILE"
    
    # Executar e medir tempo
    START_TIME=$(date +%s)
    ./bin/cargaparcial -e "$EXCEL_FILE" -w "$WORKERS" -o "$OUTPUT_FILE" 2>&1 | tee -a benchmark_log_${WORKERS}w.txt
    EXIT_CODE=${PIPESTATUS[0]}
    END_TIME=$(date +%s)
    
    ELAPSED=$((END_TIME - START_TIME))
    
    if [ $EXIT_CODE -eq 0 ]; then
        echo "" | tee -a "$RESULTS_FILE"
        echo "✅ Sucesso!" | tee -a "$RESULTS_FILE"
        echo "⏱️  Tempo: ${ELAPSED}s" | tee -a "$RESULTS_FILE"
        
        # Extrair estatísticas do resultado
        if [ -f "$OUTPUT_FILE" ]; then
            SUCCESS_COUNT=$(grep -o '"Status":"ok"' "$OUTPUT_FILE" | wc -l)
            FAIL_COUNT=$(grep -o '"Status":"fail"' "$OUTPUT_FILE" | wc -l)
            TOTAL=$((SUCCESS_COUNT + FAIL_COUNT))
            
            if [ $TOTAL -gt 0 ] && [ $ELAPSED -gt 0 ]; then
                RATE=$((TOTAL / ELAPSED))
                echo "📊 Processados: $TOTAL itens" | tee -a "$RESULTS_FILE"
                echo "✓  Sucessos: $SUCCESS_COUNT" | tee -a "$RESULTS_FILE"
                echo "✗  Falhas: $FAIL_COUNT" | tee -a "$RESULTS_FILE"
                echo "⚡ Velocidade: ~${RATE} itens/seg" | tee -a "$RESULTS_FILE"
            fi
        fi
    else
        echo "" | tee -a "$RESULTS_FILE"
        echo "❌ Falhou (exit code: $EXIT_CODE)" | tee -a "$RESULTS_FILE"
    fi
    
    echo "" | tee -a "$RESULTS_FILE"
    
    # Pequena pausa entre testes
    sleep 2
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" | tee -a "$RESULTS_FILE"
echo "🏁 Benchmark Concluído!" | tee -a "$RESULTS_FILE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" | tee -a "$RESULTS_FILE"
echo "" | tee -a "$RESULTS_FILE"
echo "📊 Relatório completo salvo em: $RESULTS_FILE"
echo "📝 Logs individuais: benchmark_log_*w.txt"
echo ""
echo "💡 Dica: Use a configuração com melhor velocidade (itens/seg)"
echo ""
