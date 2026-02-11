package main

import (
	"fmt"
	"log"

	"github.thiagohmm.com.br/cargaparcial/infrastructure/file"
)

func main() {
	fmt.Println("=== Teste do Leitor de XLSX ===")
	fmt.Println()

	// Testar leitura do arquivo
	fmt.Println("1. Testando leitura do arquivo dados_exemplo.xlsx...")
	xlsxData, err := file.ReadXLSX("dados_exemplo.xlsx")
	if err != nil {
		log.Fatalf("❌ Erro ao ler arquivo: %v", err)
	}

	fmt.Println("✅ Arquivo lido com sucesso!")
	fmt.Println()

	// Exibir resultados
	fmt.Printf("2. Códigos IBM únicos encontrados: %d\n", len(xlsxData.IBMCodes))
	fmt.Println("   Códigos IBM:")
	for i, ibm := range xlsxData.IBMCodes {
		fmt.Printf("   [%d] %s\n", i+1, ibm)
	}
	fmt.Println()

	fmt.Printf("3. Códigos de produto únicos encontrados: %d\n", len(xlsxData.ProductCodes))
	fmt.Println("   Códigos de Produto:")
	for i, code := range xlsxData.ProductCodes {
		fmt.Printf("   [%d] %s\n", i+1, code)
	}
	fmt.Println()

	// Calcular combinações
	totalCombinations := len(xlsxData.IBMCodes) * len(xlsxData.ProductCodes)
	fmt.Printf("4. Total de combinações a processar: %d\n", totalCombinations)
	fmt.Println()

	// Testar ReadXLSXPairs
	fmt.Println("5. Testando leitura de pares específicos...")
	pairs, err := file.ReadXLSXPairs("dados_exemplo.xlsx")
	if err != nil {
		log.Fatalf("❌ Erro ao ler pares: %v", err)
	}

	fmt.Println("✅ Pares lidos com sucesso!")
	fmt.Println("   Mapeamento IBM -> Produtos:")
	for ibm, products := range pairs {
		fmt.Printf("   %s -> %d produtos\n", ibm, len(products))
		for _, product := range products {
			fmt.Printf("      - %s\n", product)
		}
	}
	fmt.Println()

	fmt.Println("=== ✅ Todos os testes passaram! ===")
}
