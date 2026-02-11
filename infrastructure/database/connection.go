package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
)

// Config contém as configurações de conexão com o banco de dados Oracle
type Config struct {
	Host        string
	Port        int
	ServiceName string
	User        string
	Password    string
	Schema      string
	Driver      string
}

// NewConnection cria uma nova conexão com o banco de dados Oracle
func NewConnection(config Config) (*sql.DB, error) {
	// Configurar as opções de URL com timeout mais longo e SSL
	urlOptions := map[string]string{

		"SSL":        "enable", // Habilitar SSL
		"SSL Verify": "false",  // Não verificar certificado (ajustar conforme necessário)
	}

	// Construir a string de conexão
	connStr := go_ora.BuildUrl(
		config.Host,
		config.Port,
		config.ServiceName,
		config.User,
		config.Password,
		urlOptions,
	)

	log.Printf("Tentando conectar ao banco de dados Oracle em %s:%d (service: %s)...", config.Host, config.Port, config.ServiceName)

	// Abrir a conexão com o banco de dados
	db, err := sql.Open(config.Driver, connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com banco de dados: %w", err)
	}

	// Verificar a conexão com timeout maior
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	log.Println("Verificando conexão com o banco de dados...")
	if err = db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("erro ao verificar a conexão: %w", err)
	}

	// Se um schema foi especificado, definir como schema padrão
	if config.Schema != "" {
		_, err = db.Exec(fmt.Sprintf("ALTER SESSION SET CURRENT_SCHEMA = %s", config.Schema))
		if err != nil {
			return nil, fmt.Errorf("erro ao definir schema padrão: %w", err)
		}
	}

	// Configurar pool de conexões para alta concorrência
	db.SetMaxOpenConns(100) // Aumentado de 25 para 100
	db.SetMaxIdleConns(20)  // Aumentado de 5 para 20
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	log.Println("Conexão com banco de dados Oracle estabelecida com sucesso!")
	return db, nil
}
