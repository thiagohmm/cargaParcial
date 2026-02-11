package services

// QueueService define as operações para envio de mensagens para fila
type QueueService interface {
	Send(message string) error
}
