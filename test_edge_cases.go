package main

import (
	"fmt"
	"log"

	"github.thiagohmm.com.br/cargaparcial/infrastructure/file"
)

func testFile(filename string, description string) bool {
	fmt.Printf("Testando: %s\n", description)
	fmt.Printf("Arquivo: %s\n", filename)

	xlsxData, err := file.ReadXLSX(filename)
	if err != nil {
		fmt.Printf("❌ FALHOU: %v\n\n", err)
		return false
	}

	fmt.Printf("✅ PASSOU!\n")
	fmt.Printf("   - IBM codes: %d\n", len(xlsxData.IBMCodes))
	fmt.Printf("   - Product codes: %d\n", len(xlsxData.ProductCodes))
	fmt.Printf("   - Total combinations: %d\n\n", len(xlsxData.IBMCodes)*len(xlsxData.ProductCodes))
	return true
}

func main() {
	fmt.Println("=== Teste de Edge Cases do Leitor XLSX ===")
	fmt.Println()

	passed := 0
	failed := 0

	// Teste 1: Ordem invertida
	if testFile("teste_ordem_invertida.xlsx", "Colunas em ordem invertida (CODIGOBARRAS, IMBLOJA)") {
		passed++
	} else {
		failed++
	}

	// Teste 2: Lowercase
	if testFile("teste_lowercase.xlsx", "Nomes de colunas em lowercase (imbloja, codigobarras)") {
		passed++
	} else {
		failed++
	}

	// Teste 3: Linhas vazias
	if testFile("teste_linhas_vazias.xlsx", "Arquivo com linhas vazias") {
		passed++
	} else {
		failed++
	}

	// Teste 4: Mixed case
	if testFile("teste_mixed_case.xlsx", "Nomes de colunas em mixed case (ImBLoJa, CoDiGoBarRaS)") {
		passed++
	} else {
		failed++
	}

	// Teste 5: Arquivo inexistente
	fmt.Println("Testando: Arquivo inexistente")
	fmt.Println("Arquivo: arquivo_que_nao_existe.xlsx")
	_, err := file.ReadXLSX("arquivo_que_nao_existe.xlsx")
	if err != nil {
		fmt.Printf("✅ PASSOU! (Erro esperado capturado: %v)\n\n", err)
		passed++
	} else {
		fmt.Println("❌ FALHOU: Deveria ter retornado erro\n")
		failed++
	}

	// Resumo
	fmt.Println("=== Resumo dos Testes ===")
	fmt.Printf("✅ Passaram: %d\n", passed)
	fmt.Printf("❌ Falharam: %d\n", failed)
	fmt.Printf("Total: %d\n", passed+failed)

	if failed == 0 {
		fmt.Println("\n🎉 Todos os testes de edge cases passaram!")
	} else {
		log.Fatalf("\n⚠️  Alguns testes falharam!")
	}
}
