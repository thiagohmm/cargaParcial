package usecase

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.thiagohmm.com.br/cargaparcial/domain/entities"
	"github.thiagohmm.com.br/cargaparcial/domain/repositories"
	"github.thiagohmm.com.br/cargaparcial/domain/services"
	"github.thiagohmm.com.br/cargaparcial/usecase/dto"
)

var (
	processedItems  int64
	startTime       time.Time
	lastProgressLog time.Time
)

// ProcessProductsUseCase implementa a lógica de negócio para processar produtos
type ProcessProductsUseCase struct {
	dealerRepo             repositories.DealerRepository
	productRepo            repositories.ProductRepository
	productDealerRepo      repositories.ProductDealerRepository
	productIntegrationRepo repositories.ProductIntegrationStagingRepository
	queueService           services.QueueService
	maxWorkers             int
	dealerCache            map[string]*entities.Dealer // Cache de dealers
	dealerCacheMutex       sync.RWMutex                // Mutex para acesso seguro ao cache
	
	// Batch processing
	batchProductDealers      []*entities.ProductDealer
	batchProductDealersMutex sync.Mutex
	batchSize                int
}

// NewProcessProductsUseCase cria uma nova instância do use case
func NewProcessProductsUseCase(
	dealerRepo repositories.DealerRepository,
	productRepo repositories.ProductRepository,
	productDealerRepo repositories.ProductDealerRepository,
	productIntegrationRepo repositories.ProductIntegrationStagingRepository,
	queueService services.QueueService,
) *ProcessProductsUseCase {
	// Define o número de workers baseado no número de CPUs disponíveis
	// Multiplica por 2 para melhor aproveitamento de I/O bound operations
	maxWorkers := runtime.NumCPU() * 2
	if maxWorkers < 4 {
		maxWorkers = 4
	}

	return &ProcessProductsUseCase{
		dealerRepo:             dealerRepo,
		productRepo:            productRepo,
		productDealerRepo:      productDealerRepo,
		productIntegrationRepo: productIntegrationRepo,
		queueService:           queueService,
		maxWorkers:             maxWorkers,
		dealerCache:            make(map[string]*entities.Dealer),
		batchProductDealers:    make([]*entities.ProductDealer, 0, 500),
		batchSize:              100, // Flush a cada 100 items
	}
}

// SetMaxWorkers permite configurar o número máximo de workers
func (uc *ProcessProductsUseCase) SetMaxWorkers(workers int) {
	if workers > 0 {
		uc.maxWorkers = workers
	}
}

// JobInput representa um trabalho a ser processado
type JobInput struct {
	Dealer      *entities.Dealer
	ProductCode string
}

// Execute executa o processamento de produtos com paralelização
func (uc *ProcessProductsUseCase) Execute(input dto.ProcessProductsInput) (*dto.ProcessProductsOutput, error) {
	// Reset contadores globais
	atomic.StoreInt64(&processedItems, 0)
	startTime = time.Now()
	lastProgressLog = time.Now()

	log.Printf("Iniciando processamento paralelo com %d workers", uc.maxWorkers)

	// Calcular tamanho do buffer baseado no volume de trabalho
	totalItems := len(input.IBMCodes) * len(input.ProductCodes)
	bufferSize := 1000
	if totalItems < bufferSize {
		bufferSize = totalItems
	}

	// Canais para comunicação entre goroutines com buffer maior
	jobs := make(chan JobInput, bufferSize)
	results := make(chan dto.ProductResultDTO, bufferSize)

	// WaitGroup para aguardar conclusão de todos os workers
	var wg sync.WaitGroup

	// Iniciar workers
	for w := 1; w <= uc.maxWorkers; w++ {
		wg.Add(1)
		go uc.worker(w, jobs, results, &wg)
	}

	// Goroutine para coletar resultados
	output := &dto.ProcessProductsOutput{
		SuccessList: make([]dto.ProductResultDTO, 0, totalItems/2),
		FailureList: make([]dto.ProductResultDTO, 0, totalItems/10),
	}

	var resultWg sync.WaitGroup
	resultWg.Add(1)
	go func() {
		defer resultWg.Done()
		for result := range results {
			if result.Status == "ok" {
				output.SuccessList = append(output.SuccessList, result)
			} else {
				output.FailureList = append(output.FailureList, result)
			}
		}
	}()

	// Pré-carregar dealers no cache para evitar consultas repetidas
	dealerMap := make(map[string]*entities.Dealer)
	for _, ibmCode := range input.IBMCodes {
		if ibmCode == "" {
			ibmCode = "0"
		}

		// Verificar cache primeiro
		uc.dealerCacheMutex.RLock()
		dealer, cached := uc.dealerCache[ibmCode]
		uc.dealerCacheMutex.RUnlock()

		if !cached {
			// Buscar revendedor por IBM
			var err error
			dealer, err = uc.dealerRepo.GetByIBM(ibmCode)
			if err != nil {
				log.Printf("Erro ao buscar revendedor por IBM %s: %v", ibmCode, err)
				continue
			}

			if dealer == nil {
				log.Printf("Revendedor não encontrado para IBM: %s", ibmCode)
				continue
			}

			// Adicionar ao cache
			uc.dealerCacheMutex.Lock()
			uc.dealerCache[ibmCode] = dealer
			uc.dealerCacheMutex.Unlock()
		}

		dealerMap[ibmCode] = dealer
	}

	// Enviar jobs para processamento
	totalJobs := 0

	// Se temos o mapeamento IBM -> Produtos, usar ele
	if len(input.IBMToProducts) > 0 {
		log.Println("📋 Usando relacionamento IBM → Produtos do arquivo")

		for ibmCode, dealer := range dealerMap {
			// Pegar apenas os produtos associados a este IBM
			products, exists := input.IBMToProducts[ibmCode]
			if !exists || len(products) == 0 {
				log.Printf("⚠️  IBM %s não tem produtos associados no arquivo", ibmCode)
				continue
			}

			// Enviar jobs apenas para os produtos deste IBM
			for _, productCode := range products {
				jobs <- JobInput{
					Dealer:      dealer,
					ProductCode: productCode,
				}
				totalJobs++
			}
		}
	} else {
		// Modo legado: produto cartesiano (todas as combinações)
		log.Println("⚠️  Usando modo legado: todas as combinações IBM × Produtos")

		for ibmCode, dealer := range dealerMap {
			// Enviar jobs para cada produto
			for _, productCode := range input.ProductCodes {
				jobs <- JobInput{
					Dealer:      dealer,
					ProductCode: productCode,
				}
				totalJobs++
			}
			_ = ibmCode // evita warning unused
		}
	}

	log.Printf("Total de %d jobs enviados para processamento", totalJobs)

	// Fechar canal de jobs (não haverá mais trabalhos)
	close(jobs)

	// Aguardar todos os workers terminarem
	wg.Wait()

	// Fechar canal de resultados
	close(results)

	// Aguardar coleta de todos os resultados
	resultWg.Wait()

	log.Printf("Processamento concluído: %d jobs processados", totalJobs)
	log.Printf("Sucessos: %d, Falhas: %d", len(output.SuccessList), len(output.FailureList))

	// Flush final do batch de ProductDealers
	if err := uc.flushProductDealerBatch(); err != nil {
		log.Printf("Erro ao fazer flush final do batch: %v", err)
	}

	// Enviar mensagem "mover" para a fila "integracao"
	if err := uc.queueService.Send("mover"); err != nil {
		log.Printf("Erro ao enviar mensagem para fila: %v", err)
	}

	return output, nil
}

// worker processa jobs do canal
func (uc *ProcessProductsUseCase) worker(id int, jobs <-chan JobInput, results chan<- dto.ProductResultDTO, wg *sync.WaitGroup) {
	defer wg.Done()

	processedCount := 0
	for job := range jobs {
		result := uc.processProduct(job.Dealer, job.ProductCode)
		results <- result
		processedCount++

		// Incrementa contador global
		total := atomic.AddInt64(&processedItems, 1)

		// Log de progresso a cada 5 segundos
		if time.Since(lastProgressLog) >= 5*time.Second {
			lastProgressLog = time.Now()
			elapsed := time.Since(startTime).Seconds()
			rate := float64(total) / elapsed
			log.Printf("⚡ Progresso: %d itens | %.0f items/seg | Tempo: %.1fs", total, rate, elapsed)
		}
	}

	log.Printf("Worker %d finalizado: processou %d itens no total", id, processedCount)
}

// processProduct processa um único produto para um revendedor
func (uc *ProcessProductsUseCase) processProduct(dealer *entities.Dealer, productCode string) dto.ProductResultDTO {
	dealerID := dealer.ID

	// Buscar produto por EAN
	products, err := uc.productRepo.GetByEAN(productCode)
	if err != nil || len(products) == 0 {
		return dto.ProductResultDTO{
			DealerID:  &dealerID,
			ProductID: nil,
			EAN:       productCode,
			Status:    "fail",
			Reason:    "Produto não encontrado pelo EAN",
		}
	}

	product := products[0]
	productID := product.ID

	// Verificar se já existe relação ProductDealer
	exists, err := uc.productDealerRepo.Exists(productID, dealerID)
	if err != nil {
		log.Printf("Erro ao verificar ProductDealer: %v", err)
		return dto.ProductResultDTO{
			DealerID:  &dealerID,
			ProductID: &productID,
			Status:    "fail",
			Reason:    "Erro ao verificar relação produto-revendedor",
		}
	}

	// Criar relação se não existir - usando BATCH
	if !exists {
		productDealer := &entities.ProductDealer{
			ProductID: productID,
			DealerID:  dealerID,
			IsActive:  true,
		}

		// Adiciona ao batch (faz flush automático se necessário)
		if err := uc.addToProductDealerBatch(productDealer); err != nil {
			log.Printf("Erro ao adicionar ProductDealer ao batch: %v", err)
			return dto.ProductResultDTO{
				DealerID:  &dealerID,
				ProductID: &productID,
				Status:    "fail",
				Reason:    "Erro ao criar relação produto-revendedor (batch)",
			}
		}
	}

	// Gravar integração produto staging (chama a stored procedure)
	if err := uc.productRepo.SaveIntegrationStaging(dealerID, productID); err != nil {
		log.Printf("Erro ao gravar integração produto staging: %v", err)
		return dto.ProductResultDTO{
			DealerID:  &dealerID,
			ProductID: &productID,
			Status:    "fail",
			Reason:    "Erro ao gravar integração produto staging",
		}
	}

	// Verificar se o registro foi realmente inserido na tabela IntegracaoProdutoStaging
	// (igual ao código TypeScript que faz productIntegrationStagingQuery.getByProductIntegrationStaging)
	staging, err := uc.productIntegrationRepo.GetByProductAndDealer(productID, dealerID)
	if err != nil {
		log.Printf("Erro ao verificar ProductIntegrationStaging: %v", err)
		return dto.ProductResultDTO{
			DealerID:  &dealerID,
			ProductID: &productID,
			Status:    "fail",
			Reason:    "Erro ao verificar integração produto staging",
		}
	}

	// Se o registro existe, retorna sucesso. Caso contrário, falha.
	if staging != nil {
		return dto.ProductResultDTO{
			DealerID:  &dealerID,
			ProductID: &productID,
			Status:    "ok",
		}
	}

	return dto.ProductResultDTO{
		DealerID:  &dealerID,
		ProductID: &productID,
		Status:    "fail",
		Reason:    "Registro não encontrado após chamada da procedure",
	}
}

// addToProductDealerBatch adiciona um ProductDealer ao batch e faz flush se necessário
func (uc *ProcessProductsUseCase) addToProductDealerBatch(productDealer *entities.ProductDealer) error {
	uc.batchProductDealersMutex.Lock()
	defer uc.batchProductDealersMutex.Unlock()

	uc.batchProductDealers = append(uc.batchProductDealers, productDealer)

	// Se atingiu o tamanho do batch, faz o flush
	if len(uc.batchProductDealers) >= uc.batchSize {
		return uc.flushProductDealerBatchUnsafe()
	}

	return nil
}

// flushProductDealerBatch faz o flush do batch com lock
func (uc *ProcessProductsUseCase) flushProductDealerBatch() error {
	uc.batchProductDealersMutex.Lock()
	defer uc.batchProductDealersMutex.Unlock()

	return uc.flushProductDealerBatchUnsafe()
}

// flushProductDealerBatchUnsafe faz o flush sem lock (deve ser chamado com lock já adquirido)
func (uc *ProcessProductsUseCase) flushProductDealerBatchUnsafe() error {
	if len(uc.batchProductDealers) == 0 {
		return nil
	}

	log.Printf("🚀 Fazendo batch insert de %d ProductDealers", len(uc.batchProductDealers))

	err := uc.productDealerRepo.CreateBatch(uc.batchProductDealers)
	if err != nil {
		return fmt.Errorf("erro ao criar batch de ProductDealers: %w", err)
	}

	// Limpar o batch
	uc.batchProductDealers = uc.batchProductDealers[:0]

	return nil
}
