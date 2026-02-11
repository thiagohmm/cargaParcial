# 🐰 Configuração do RabbitMQ

## 📋 Visão Geral

O sistema agora está integrado com **RabbitMQ** para enviar mensagens de integração após processar os produtos. Quando o processamento é concluído com sucesso, uma mensagem `"mover"` é enviada para a fila `integracao`.

## 🔧 Configuração

### 1. Instalar RabbitMQ

#### Docker (Recomendado)
```bash
docker run -d --name rabbitmq \
  -p 5672:5672 \
  -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=admin \
  -e RABBITMQ_DEFAULT_PASS=admin123 \
  rabbitmq:3-management
```

#### Instalação Local
- **Ubuntu/Debian**: 
  ```bash
  sudo apt-get install rabbitmq-server
  sudo systemctl start rabbitmq-server
  ```
- **macOS**: 
  ```bash
  brew install rabbitmq
  brew services start rabbitmq
  ```
- **Windows**: Baixar de https://www.rabbitmq.com/download.html

### 2. Configurar Variável de Ambiente

Adicione no arquivo `.env`:

```env
# RabbitMQ Configuration
ENV_RABBITMQ=amqp://admin:admin123@localhost:5672/
```

**Formato da URL:**
```
amqp://username:password@host:port/vhost
```

**Exemplos:**
- Local padrão: `amqp://guest:guest@localhost:5672/`
- Com autenticação: `amqp://admin:admin123@localhost:5672/`
- RabbitMQ na nuvem: `amqp://user:pass@my-rabbit-server.com:5672/`
- Com vhost específico: `amqp://user:pass@localhost:5672/myvhost`

### 3. Verificar Conexão

Acesse o painel de gerenciamento do RabbitMQ:
```
http://localhost:15672
```

Credenciais padrão (se usando Docker acima):
- **Usuário**: admin
- **Senha**: admin123

## 🚀 Como Funciona

### Fluxo de Processamento

1. **Processamento de Produtos**: O sistema processa produtos e revendedores
2. **Chamada da Stored Procedure**: `SP_GRAVARINTEGRACAOPRODUTOSTAGING` é executada
3. **Verificação**: Confirma que o registro foi inserido em `IntegracaoProdutoStaging`
4. **Envio para Fila**: Ao final, envia mensagem `"mover"` para a fila `integracao`

### Código de Envio

```go
// Ao final do processamento (process_products_usecase.go linha 215-217)
if err := uc.queueService.Send("mover"); err != nil {
    log.Printf("Erro ao enviar mensagem para fila: %v", err)
}
```

### Consumir Mensagens da Fila

Para consumir as mensagens da fila `integracao`, você pode criar um consumer:

```go
package main

import (
    "log"
    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
    // Conectar ao RabbitMQ
    conn, err := amqp.Dial("amqp://admin:admin123@localhost:5672/")
    if err != nil {
        log.Fatalf("Erro ao conectar: %v", err)
    }
    defer conn.Close()

    // Criar canal
    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Erro ao criar canal: %v", err)
    }
    defer ch.Close()

    // Consumir mensagens
    msgs, err := ch.Consume(
        "integracao", // fila
        "",           // consumer
        true,         // auto-ack
        false,        // exclusive
        false,        // no-local
        false,        // no-wait
        nil,          // args
    )
    if err != nil {
        log.Fatalf("Erro ao registrar consumer: %v", err)
    }

    log.Println("Aguardando mensagens...")

    for msg := range msgs {
        log.Printf("📬 Mensagem recebida: %s", msg.Body)
        
        // Processar a mensagem "mover" aqui
        // Exemplo: iniciar integração com sistema externo
    }
}
```

## 🐛 Troubleshooting

### Erro: "Erro ao conectar ao RabbitMQ"

**Solução:**
1. Verifique se o RabbitMQ está rodando:
   ```bash
   # Docker
   docker ps | grep rabbitmq
   
   # Linux
   sudo systemctl status rabbitmq-server
   ```

2. Verifique se a porta 5672 está aberta:
   ```bash
   telnet localhost 5672
   ```

3. Verifique as credenciais na URL do `.env`

### Modo Degradado (Sem RabbitMQ)

Se o RabbitMQ não estiver disponível, o sistema **continua funcionando normalmente** em modo simulado:

```
⚠️  RabbitMQ URL não configurada, fila será simulada
📤 [SIMULADO] Mensagem para fila 'integracao': mover
```

O processamento de produtos **não é afetado**, apenas o envio da mensagem é logado ao invés de enviado para a fila real.

## 📊 Monitoramento

### Ver mensagens na fila

1. Acesse o painel web: http://localhost:15672
2. Vá em **Queues** → **integracao**
3. Visualize:
   - Total de mensagens
   - Taxa de mensagens/segundo
   - Consumidores ativos

### Logs do Sistema

O sistema loga todas as operações da fila:

```
✅ Conectado ao RabbitMQ. Fila 'integracao' pronta para uso.
✅ Mensagem enviada para fila 'integracao': mover
```

## 🔒 Segurança em Produção

Para ambientes de produção, considere:

1. **Não usar credenciais padrão** (guest:guest)
2. **Usar TLS/SSL**: `amqps://` ao invés de `amqp://`
3. **Criar usuário dedicado** com permissões limitadas
4. **Usar variáveis de ambiente** para credenciais sensíveis
5. **Configurar vhost separado** para isolamento

Exemplo de URL segura:
```
amqps://prod_user:strong_password@rabbitmq.mycompany.com:5671/prod_vhost
```

## 📚 Referências

- [RabbitMQ Documentation](https://www.rabbitmq.com/documentation.html)
- [amqp091-go (Go Client)](https://github.com/rabbitmq/amqp091-go)
- [RabbitMQ Management Plugin](https://www.rabbitmq.com/management.html)
