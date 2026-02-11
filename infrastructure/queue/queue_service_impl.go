package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.thiagohmm.com.br/cargaparcial/domain/services"
)

// QueueServiceImpl implementa o QueueService com RabbitMQ
type QueueServiceImpl struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	queueName   string
	isConnected bool
}

// NewQueueService cria uma nova instância do serviço de fila RabbitMQ
func NewQueueService(rabbitURL string) (services.QueueService, error) {
	if rabbitURL == "" {
		log.Println("⚠️  RabbitMQ URL não configurada, fila será simulada")
		return &QueueServiceImpl{
			isConnected: false,
			queueName:   "integracao",
		}, nil
	}

	// Conectar ao RabbitMQ
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Printf("⚠️  Erro ao conectar ao RabbitMQ: %v. Fila será simulada.", err)
		return &QueueServiceImpl{
			isConnected: false,
			queueName:   "integracao",
		}, nil
	}

	// Criar canal
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Printf("⚠️  Erro ao criar canal RabbitMQ: %v. Fila será simulada.", err)
		return &QueueServiceImpl{
			isConnected: false,
			queueName:   "integracao",
		}, nil
	}

	queueName := "integracao"

	// Declarar a fila
	_, err = channel.QueueDeclare(
		queueName, // nome da fila
		true,      // durable (sobrevive a reinicializações)
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // argumentos
	)
	if err != nil {
		channel.Close()
		conn.Close()
		log.Printf("⚠️  Erro ao declarar fila RabbitMQ: %v. Fila será simulada.", err)
		return &QueueServiceImpl{
			isConnected: false,
			queueName:   queueName,
		}, nil
	}

	log.Printf("✅ Conectado ao RabbitMQ. Fila '%s' pronta para uso.", queueName)

	return &QueueServiceImpl{
		conn:        conn,
		channel:     channel,
		queueName:   queueName,
		isConnected: true,
	}, nil
}

// Send envia uma mensagem para a fila RabbitMQ
func (s *QueueServiceImpl) Send(message string) error {
	if message == "" {
		return fmt.Errorf("mensagem vazia não pode ser enviada")
	}

	// Se não está conectado, apenas loga
	if !s.isConnected {
		log.Printf("📤 [SIMULADO] Mensagem para fila '%s': %s", s.queueName, message)
		return nil
	}

	// Criar contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Publicar mensagem na fila
	err := s.channel.PublishWithContext(
		ctx,
		"",          // exchange (vazio = default)
		s.queueName, // routing key (nome da fila)
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // mensagem persistente
			ContentType:  "text/plain",
			Body:         []byte(message),
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("erro ao publicar mensagem no RabbitMQ: %w", err)
	}

	log.Printf("✅ Mensagem enviada para fila '%s': %s", s.queueName, message)
	return nil
}

// Close fecha a conexão com o RabbitMQ
func (s *QueueServiceImpl) Close() error {
	if !s.isConnected {
		return nil
	}

	var errs []error

	if s.channel != nil {
		if err := s.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("erro ao fechar canal: %w", err))
		}
	}

	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("erro ao fechar conexão: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("erros ao fechar RabbitMQ: %v", errs)
	}

	log.Println("✅ Conexão RabbitMQ fechada")
	return nil
}
