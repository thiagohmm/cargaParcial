#!/bin/bash

# Script para testar a conexão com RabbitMQ
# Uso: ./test_rabbitmq.sh

echo "🐰 Testando conexão com RabbitMQ..."
echo ""

# Verifica se o RabbitMQ está rodando
echo "1. Verificando se RabbitMQ está rodando..."
if command -v docker &> /dev/null; then
    if docker ps | grep -q rabbitmq; then
        echo "✅ RabbitMQ encontrado no Docker"
        docker ps | grep rabbitmq
    else
        echo "❌ RabbitMQ não está rodando no Docker"
        echo ""
        echo "Para iniciar RabbitMQ com Docker:"
        echo ""
        echo "docker run -d --name rabbitmq \\"
        echo "  -p 5672:5672 \\"
        echo "  -p 15672:15672 \\"
        echo "  -e RABBITMQ_DEFAULT_USER=admin \\"
        echo "  -e RABBITMQ_DEFAULT_PASS=admin123 \\"
        echo "  rabbitmq:3-management"
        exit 1
    fi
else
    echo "⚠️  Docker não encontrado, verificando instalação local..."
    if systemctl is-active --quiet rabbitmq-server; then
        echo "✅ RabbitMQ está rodando localmente"
    else
        echo "❌ RabbitMQ não está rodando"
        exit 1
    fi
fi

echo ""
echo "2. Testando porta 5672 (AMQP)..."
if nc -z localhost 5672 2>/dev/null; then
    echo "✅ Porta 5672 está acessível"
else
    echo "❌ Porta 5672 não está acessível"
    exit 1
fi

echo ""
echo "3. Testando porta 15672 (Management UI)..."
if nc -z localhost 15672 2>/dev/null; then
    echo "✅ Porta 15672 está acessível"
    echo "   Acesse: http://localhost:15672"
else
    echo "⚠️  Porta 15672 não está acessível (Management UI pode não estar habilitado)"
fi

echo ""
echo "4. Verificando arquivo .env..."
if [ -f .env ]; then
    if grep -q "ENV_RABBITMQ" .env; then
        echo "✅ ENV_RABBITMQ encontrado no .env"
        grep "ENV_RABBITMQ" .env
    else
        echo "⚠️  ENV_RABBITMQ não encontrado no .env"
        echo "   Adicione: ENV_RABBITMQ=amqp://admin:admin123@localhost:5672/"
    fi
else
    echo "⚠️  Arquivo .env não encontrado"
    echo "   Copie config.example para .env e configure"
fi

echo ""
echo "5. Testando conexão HTTP com Management API..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -u admin:admin123 http://localhost:15672/api/overview)
if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ API de gerenciamento acessível"
    
    # Tentar obter informações sobre a fila
    echo ""
    echo "6. Verificando fila 'integracao'..."
    QUEUE_INFO=$(curl -s -u admin:admin123 http://localhost:15672/api/queues/%2F/integracao)
    if echo "$QUEUE_INFO" | grep -q '"name":"integracao"'; then
        echo "✅ Fila 'integracao' existe"
        echo "$QUEUE_INFO" | grep -o '"messages":[0-9]*' | head -1
    else
        echo "⚠️  Fila 'integracao' ainda não foi criada (será criada automaticamente ao executar o programa)"
    fi
elif [ "$HTTP_CODE" = "401" ]; then
    echo "❌ Credenciais inválidas (usuário/senha incorretos)"
else
    echo "⚠️  API não acessível (código HTTP: $HTTP_CODE)"
fi

echo ""
echo "=========================================="
echo "✅ Testes concluídos!"
echo ""
echo "Para executar o programa:"
echo "  ./bin/cargaparcial -e lojas_produtos.xlsx"
echo ""
echo "Para monitorar o RabbitMQ:"
echo "  http://localhost:15672 (admin/admin123)"
echo "=========================================="
