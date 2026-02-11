package main

import (
	"fmt"
	"log"

	"github.thiagohmm.com.br/cargaparcial/infrastructure/config"
	"github.thiagohmm.com.br/cargaparcial/infrastructure/database"
	"github.thiagohmm.com.br/cargaparcial/infrastructure/file"
)

func main() {
	log.Println("=== Validador de IBMs do Excel ===\n")

	// Ler arquivo Excel
	excelFile := "lojas_produtos.xlsx"
	data, err := file.ReadXLSX(excelFile)
	if err != nil {
		log.Fatalf("Erro ao ler Excel: %v", err)
	}

	log.Printf("📋 Excel lido: %d IBMs únicos, %d produtos únicos\n", len(data.IBMCodes), len(data.ProductCodes))

	// Conectar ao banco
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Erro ao carregar config: %v", err)
	}

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

	log.Println("✅ Conectado ao banco\n")

	// Verificar cada IBM
	found := 0
	notFound := 0
	notFoundList := make([]string, 0)

	log.Println("🔍 Verificando IBMs do Excel no banco de dados...")
	log.Println("------------------------------------------------------------")

	for i, ibm := range data.IBMCodes {
		var id int
		var codigo string

		err := db.QueryRow("SELECT IdRevendedor, CodigoIBM FROM Revendedor WHERE CodigoIBM = :1", ibm).Scan(&id, &codigo)

		if err == nil {
			found++
			if found <= 5 {
				log.Printf("✅ [%3d] IBM %s → ID %d (encontrado)", i+1, ibm, id)
			}
		} else {
			notFound++
			notFoundList = append(notFoundList, ibm)
			if notFound <= 10 {
				log.Printf("❌ [%3d] IBM %s → NÃO ENCONTRADO", i+1, ibm)
			}
		}
	}

	log.Println("------------------------------------------------------------")
	log.Printf("\n📊 RESUMO:\n")
	log.Printf("   Total de IBMs no Excel: %d\n", len(data.IBMCodes))
	log.Printf("   ✅ Encontrados no banco: %d (%.1f%%)\n", found, float64(found)/float64(len(data.IBMCodes))*100)
	log.Printf("   ❌ NÃO encontrados:      %d (%.1f%%)\n", notFound, float64(notFound)/float64(len(data.IBMCodes))*100)

	if notFound > 0 {
		log.Printf("\n⚠️  IBMs NÃO encontrados no banco:\n")
		for i, ibm := range notFoundList {
			if i < 20 {
				log.Printf("   - %s\n", ibm)
			}
		}
		if len(notFoundList) > 20 {
			log.Printf("   ... e mais %d IBMs\n", len(notFoundList)-20)
		}
	}

	// Sugestão: buscar IBMs similares
	if notFound > 0 {
		log.Printf("\n💡 Sugestões:\n")
		log.Println("   1. Verifique se os IBMs do Excel têm o formato correto")
		log.Println("   2. Talvez seja necessário remover ou adicionar zeros à esquerda")
		log.Println("   3. Verifique se os revendedores existem no banco com outro código")

		// Tentar buscar IBMs similares (sem zeros à esquerda)
		testIBM := notFoundList[0]
		log.Printf("\n🔎 Testando variações para IBM '%s':\n", testIBM)

		variations := []string{
			testIBM,
			fmt.Sprintf("%010s", testIBM),        // Com 10 dígitos
			fmt.Sprintf("%d", mustAtoi(testIBM)), // Sem zeros à esquerda
		}

		for _, v := range variations {
			var id int
			err := db.QueryRow("SELECT IdRevendedor FROM Revendedor WHERE CodigoIBM = :1", v).Scan(&id)
			if err == nil {
				log.Printf("   ✅ Variação '%s' ENCONTRADA! ID=%d\n", v, id)
			} else {
				log.Printf("   ❌ Variação '%s' não encontrada\n", v)
			}
		}
	}
}

func mustAtoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
