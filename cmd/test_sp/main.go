package main

import (
	"log"
	"os"

	"github.thiagohmm.com.br/cargaparcial/infrastructure/config"
	"github.thiagohmm.com.br/cargaparcial/infrastructure/database"
)

func main() {
	log.Println("=== Teste de Stored Procedure ===")

	// Carregar configurações
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Conectar ao banco
	dbConfig := database.Config{
		Host:        cfg.Host,
		Port:        cfg.Port,
		ServiceName: cfg.ServiceName,
		User:        cfg.DBUser,
		Password:    cfg.DBPassword,
		Schema:      cfg.DBSchema,
		Driver:      cfg.DBDriver,
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}
	defer db.Close()

	log.Println("✅ Conexão estabelecida com sucesso!")

	// Teste 1: Verificar se a SP existe
	log.Println("\n--- Teste 1: Verificando se a SP existe ---")
	var objectName, objectType, status string
	err = db.QueryRow(`
		SELECT object_name, object_type, status 
		FROM user_objects 
		WHERE object_name = 'SP_GRAVARINTEGRACAOPRODUTOSTAGING' 
		AND object_type = 'PROCEDURE'
	`).Scan(&objectName, &objectType, &status)

	if err != nil {
		log.Printf("❌ SP não encontrada: %v", err)
		log.Println("⚠️  A stored procedure SP_GRAVARINTEGRACAOPRODUTOSTAGING não existe no schema!")
		log.Println("    Você precisa criar a SP ou usar INSERT direto.")
		os.Exit(1)
	}

	log.Printf("✅ SP encontrada: %s (%s) - Status: %s", objectName, objectType, status)

	// Teste 2: Verificar parâmetros da SP
	log.Println("\n--- Teste 2: Verificando parâmetros da SP ---")
	rows, err := db.Query(`
		SELECT argument_name, data_type, in_out, position
		FROM user_arguments
		WHERE object_name = 'SP_GRAVARINTEGRACAOPRODUTOSTAGING'
		ORDER BY position
	`)
	if err != nil {
		log.Printf("❌ Erro ao buscar parâmetros: %v", err)
	} else {
		defer rows.Close()
		log.Println("Parâmetros:")
		for rows.Next() {
			var argName, dataType, inOut string
			var position int
			if err := rows.Scan(&argName, &dataType, &inOut, &position); err != nil {
				log.Printf("Erro: %v", err)
				continue
			}
			log.Printf("  - %s: %s (%s) - Posição: %d", argName, dataType, inOut, position)
		}
	}

	// Teste 3: Verificar se a tabela existe
	log.Println("\n--- Teste 3: Verificando tabela ProdutoIntegracaoStaging ---")
	var tableName string
	var numRows int
	err = db.QueryRow(`
		SELECT table_name, NVL(num_rows, 0) 
		FROM user_tables 
		WHERE table_name = 'PRODUTOINTEGRACAOSTAGING'
	`).Scan(&tableName, &numRows)

	if err != nil {
		log.Printf("❌ Tabela não encontrada: %v", err)
	} else {
		log.Printf("✅ Tabela encontrada: %s - Registros: %d", tableName, numRows)
	}

	// Teste 4: Ver colunas da tabela
	log.Println("\n--- Teste 4: Estrutura da tabela ---")
	rows2, err := db.Query(`
		SELECT column_name, data_type, nullable
		FROM user_tab_columns
		WHERE table_name = 'PRODUTOINTEGRACAOSTAGING'
		ORDER BY column_id
	`)
	if err != nil {
		log.Printf("❌ Erro ao buscar colunas: %v", err)
	} else {
		defer rows2.Close()
		log.Println("Colunas:")
		for rows2.Next() {
			var colName, dataType, nullable string
			if err := rows2.Scan(&colName, &dataType, &nullable); err != nil {
				log.Printf("Erro: %v", err)
				continue
			}
			log.Printf("  - %s: %s (Nullable: %s)", colName, dataType, nullable)
		}
	}

	// Teste 5: Chamar a SP (comentado por segurança)
	log.Println("\n--- Teste 5: Executando SP de teste ---")
	log.Println("⚠️  Descomente o código abaixo para testar a execução da SP")

	// DESCOMENTE PARA TESTAR:
	/*
		query := `
			BEGIN
				SP_GRAVARINTEGRACAOPRODUTOSTAGING(:p_idRevendedor, :p_idProduto);
			END;
		`
		_, err = db.Exec(query, 999999, 999999)
		if err != nil {
			log.Printf("❌ Erro ao executar SP: %v", err)
		} else {
			log.Println("✅ SP executada com sucesso!")
		}
	*/

	log.Println("\n=== Teste concluído ===")
}
