package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.thiagohmm.com.br/cargaparcial/infrastructure/config"
	"github.thiagohmm.com.br/cargaparcial/infrastructure/database"
	"github.thiagohmm.com.br/cargaparcial/infrastructure/file"
	"github.thiagohmm.com.br/cargaparcial/infrastructure/queue"
	"github.thiagohmm.com.br/cargaparcial/infrastructure/repository"
	"github.thiagohmm.com.br/cargaparcial/usecase"
	"github.thiagohmm.com.br/cargaparcial/usecase/dto"
)

var (
	ibmFile    string
	codigoFile string
	excelFile  string
	outputFile string
	maxWorkers int
)

var rootCmd = &cobra.Command{
	Use:   "cargaparcial",
	Short: "Processador de carga parcial de produtos",
	Long: `Sistema de processamento paralelo de produtos e revendedores.
Lê códigos IBM e códigos de produtos de arquivos de entrada,
processa em paralelo e gera arquivo de resultado.`,
	Run: runProcess,
}

func init() {
	rootCmd.Flags().StringVarP(&ibmFile, "ibm", "i", "ibm.txt", "Arquivo com códigos IBM (um por linha)")
	rootCmd.Flags().StringVarP(&codigoFile, "codigo", "c", "codigo.txt", "Arquivo com códigos de produtos/EAN (um por linha)")
	rootCmd.Flags().StringVarP(&excelFile, "excel", "e", "", "Arquivo Excel (.xlsx) com colunas IMBLOJA e CODIGOBARRAS")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "resultado.json", "Arquivo de saída com resultados")
	rootCmd.Flags().IntVarP(&maxWorkers, "workers", "w", 0, "Número de workers paralelos (0 = auto, baseado em CPUs)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runProcess(cmd *cobra.Command, args []string) {
	log.Println("=== Carga Parcial - Processador de Produtos ===")

	// Verificar se está usando arquivo Excel ou arquivos TXT
	usingExcel := excelFile != ""

	if usingExcel {
		log.Printf("Arquivo Excel: %s", excelFile)
	} else {
		log.Printf("Arquivo IBM: %s", ibmFile)
		log.Printf("Arquivo Código: %s", codigoFile)
	}
	log.Printf("Arquivo Saída: %s", outputFile)

	// Carregar configurações usando Viper
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Criar configuração do banco de dados
	dbConfig := database.Config{
		Host:        cfg.Host,
		Port:        cfg.Port,
		ServiceName: cfg.ServiceName,
		User:        cfg.DBUser,
		Password:    cfg.DBPassword,
		Schema:      cfg.DBSchema,
		Driver:      cfg.DBDriver,
	}

	// Conectar ao banco de dados
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	log.Println("✓ Conexão com banco de dados estabelecida")

	// Inicializar repositórios
	dealerRepo := repository.NewDealerRepository(db)
	productRepo := repository.NewProductRepository(db)
	productDealerRepo := repository.NewProductDealerRepository(db)
	productIntegrationRepo := repository.NewProductIntegrationStagingRepository(db)

	// Inicializar serviço de fila RabbitMQ
	queueService, err := queue.NewQueueService(cfg.ENV_RABBITMQ)
	if err != nil {
		log.Printf("⚠️  Erro ao inicializar serviço de fila: %v", err)
		log.Println("Continuando com fila simulada...")
	}

	// Inicializar use case
	processProductsUseCase := usecase.NewProcessProductsUseCase(
		dealerRepo,
		productRepo,
		productDealerRepo,
		productIntegrationRepo,
		queueService,
	)

	// Configurar número de workers se especificado
	if maxWorkers > 0 {
		processProductsUseCase.SetMaxWorkers(maxWorkers)
		log.Printf("✓ Configurado para usar %d workers", maxWorkers)
	}

	var ibmCodes []string
	var productCodes []string
	var ibmToProducts map[string][]string
	var totalCombinations int

	// Ler arquivos de entrada
	if usingExcel {
		// Ler arquivo Excel
		log.Printf("Lendo arquivo Excel: %s", excelFile)
		xlsxData, err := file.ReadXLSX(excelFile)
		if err != nil {
			log.Fatalf("Erro ao ler arquivo Excel %s: %v", excelFile, err)
		}

		ibmCodes = xlsxData.IBMCodes
		productCodes = xlsxData.ProductCodes
		ibmToProducts = xlsxData.IBMToProducts

		log.Printf("✓ Lidos %d códigos IBM únicos", len(ibmCodes))
		log.Printf("✓ Lidos %d códigos de produto únicos", len(productCodes))

		// Calcular total real de combinações (apenas as do arquivo)
		totalCombinations = 0
		for _, products := range ibmToProducts {
			totalCombinations += len(products)
		}
		log.Printf("Total de combinações a processar: %d (relacionamento IBM → Produtos)", totalCombinations)
	} else {
		// Ler arquivos TXT tradicionais
		log.Printf("Lendo arquivo: %s", ibmFile)
		var err error
		ibmCodes, err = readLinesFromFile(ibmFile)
		if err != nil {
			log.Fatalf("Erro ao ler arquivo %s: %v", ibmFile, err)
		}
		log.Printf("✓ Lidos %d códigos IBM", len(ibmCodes))

		log.Printf("Lendo arquivo: %s", codigoFile)
		productCodes, err = readLinesFromFile(codigoFile)
		if err != nil {
			log.Fatalf("Erro ao ler arquivo %s: %v", codigoFile, err)
		}
		log.Printf("✓ Lidos %d códigos de produto", len(productCodes))

		totalCombinations = len(ibmCodes) * len(productCodes)
		log.Printf("Total de combinações a processar: %d", totalCombinations)
	}

	log.Println("Iniciando processamento paralelo...")

	// Processar produtos
	input := dto.ProcessProductsInput{
		IBMCodes:      ibmCodes,
		ProductCodes:  productCodes,
		IBMToProducts: ibmToProducts, // Passa o relacionamento correto
	}

	output, err := processProductsUseCase.Execute(input)
	if err != nil {
		log.Fatalf("Erro ao processar produtos: %v", err)
	}

	// Exibir resultados
	log.Println("=== Processamento Concluído ===")
	log.Printf("✓ Sucessos: %d", len(output.SuccessList))
	log.Printf("✗ Falhas: %d", len(output.FailureList))

	successRate := 0.0
	if totalCombinations > 0 {
		successRate = float64(len(output.SuccessList)) / float64(totalCombinations) * 100
	}
	log.Printf("Taxa de sucesso: %.2f%%", successRate)

	// Salvar resultado em arquivo JSON
	log.Printf("Salvando resultados em: %s", outputFile)
	resultJSON, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("Erro ao gerar JSON de resultado: %v", err)
	}

	if err := os.WriteFile(outputFile, resultJSON, 0644); err != nil {
		log.Fatalf("Erro ao salvar %s: %v", outputFile, err)
	}

	log.Printf("✓ Resultado salvo com sucesso em %s", outputFile)
	log.Println("=== Processo Finalizado ===")
}

// readLinesFromFile lê todas as linhas de um arquivo
func readLinesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Ignorar linhas vazias e comentários
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro ao ler linha %d: %w", lineNumber, err)
	}

	return lines, nil
}
