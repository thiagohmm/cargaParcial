#!/bin/bash

# Script de execução otimizada
# Configura o ambiente para máxima performance

echo "⚡ Execução Otimizada - Carga Parcial"
echo "====================================="
echo ""

# Configuração padrão
WORKERS=96
EXCEL_FILE="lojas_produtos.xlsx"
OUTPUT_FILE="resultado.json"

# Processar argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--excel)
            EXCEL_FILE="$2"
            shift 2
            ;;
        -w|--workers)
            WORKERS="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        --fast)
            WORKERS=96
            shift
            ;;
        --turbo)
            WORKERS=150
            shift
            ;;
        --max)
            WORKERS=200
            shift
            ;;
        -h|--help)
            echo "Uso: $0 [opções]"
            echo ""
            echo "Opções:"
            echo "  -e, --excel FILE    Arquivo Excel (padrão: lojas_produtos.xlsx)"
            echo "  -o, --output FILE   Arquivo de saída (padrão: resultado.json)"
            echo "  -w, --workers NUM   Número de workers (padrão: 96)"
            echo "  --fast              Preset rápido (96 workers)"
            echo "  --turbo             Preset turbo (150 workers)"
            echo "  --max               Preset máximo (200 workers)"
            echo "  -h, --help          Mostrar esta ajuda"
            echo ""
            echo "Exemplos:"
            echo "  $0 --fast"
            echo "  $0 -e dados.xlsx --turbo"
            echo "  $0 -e dados.xlsx -w 120 -o saida.json"
            exit 0
            ;;
        *)
            echo "Opção desconhecida: $1"
            echo "Use -h para ajuda"
            exit 1
            ;;
    esac
done

# Verificações
if [ ! -f "./bin/cargaparcial" ]; then
    echo "❌ Binário não encontrado. Execute 'make build' primeiro."
    exit 1
fi

if [ ! -f "$EXCEL_FILE" ]; then
    echo "❌ Arquivo não encontrado: $EXCEL_FILE"
    exit 1
fi

# Configurar ambiente
echo "🔧 Configurando ambiente para máxima performance..."
echo ""

# Aumentar limite de arquivos
ORIGINAL_ULIMIT=$(ulimit -n)
ulimit -n 65536
echo "✅ ulimit -n: $ORIGINAL_ULIMIT → 65536"

# Informações do sistema
echo "✅ CPUs: $(nproc) cores"
echo "✅ Memória: $(free -h | grep Mem | awk '{print $7}') disponível"
echo ""

# Configuração da execução
echo "📊 Configuração:"
echo "  • Arquivo: $EXCEL_FILE"
echo "  • Workers: $WORKERS"
echo "  • Saída: $OUTPUT_FILE"
echo ""

# Determinar preset
PRESET="customizado"
if [ "$WORKERS" -eq 96 ]; then
    PRESET="🚀 FAST (recomendado)"
elif [ "$WORKERS" -eq 150 ]; then
    PRESET="⚡ TURBO"
elif [ "$WORKERS" -eq 200 ]; then
    PRESET="🔥 MÁXIMO"
elif [ "$WORKERS" -eq 24 ]; then
    PRESET="📊 PADRÃO"
fi

echo "  • Preset: $PRESET"
echo ""

# Confirmar execução
read -p "Continuar? [Y/n] " -n 1 -r
echo
if [[ $REPLY =~ ^[Nn]$ ]]; then
    echo "Cancelado."
    exit 0
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 Iniciando processamento..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Executar
START_TIME=$(date +%s)
./bin/cargaparcial -e "$EXCEL_FILE" -w "$WORKERS" -o "$OUTPUT_FILE"
EXIT_CODE=$?
END_TIME=$(date +%s)

ELAPSED=$((END_TIME - START_TIME))

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ $EXIT_CODE -eq 0 ]; then
    echo "✅ Processamento concluído com sucesso!"
else
    echo "❌ Processamento falhou (exit code: $EXIT_CODE)"
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "⏱️  Tempo total: ${ELAPSED}s"

# Estatísticas do resultado
if [ -f "$OUTPUT_FILE" ]; then
    SUCCESS_COUNT=$(grep -o '"Status":"ok"' "$OUTPUT_FILE" | wc -l)
    FAIL_COUNT=$(grep -o '"Status":"fail"' "$OUTPUT_FILE" | wc -l)
    TOTAL=$((SUCCESS_COUNT + FAIL_COUNT))
    
    echo "📊 Estatísticas:"
    echo "  • Total: $TOTAL itens"
    echo "  • Sucessos: $SUCCESS_COUNT"
    echo "  • Falhas: $FAIL_COUNT"
    
    if [ $TOTAL -gt 0 ] && [ $ELAPSED -gt 0 ]; then
        RATE=$((TOTAL / ELAPSED))
        SUCCESS_RATE=$((SUCCESS_COUNT * 100 / TOTAL))
        echo "  • Taxa de sucesso: ${SUCCESS_RATE}%"
        echo "  • Velocidade: ~${RATE} itens/seg"
    fi
    
    echo ""
    echo "💾 Resultado salvo em: $OUTPUT_FILE"
fi

echo ""

# Restaurar ulimit
ulimit -n "$ORIGINAL_ULIMIT"

exit $EXIT_CODE
