package repository

import (
	"database/sql"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.thiagohmm.com.br/cargaparcial/domain/entities"
	"github.thiagohmm.com.br/cargaparcial/domain/repositories"
)

var (
	spCallCount  int64
	spTotalTime  int64
	spErrorCount int64
	lastLogTime  int64
)

// ProductRepositoryImpl implementa o ProductRepository
type ProductRepositoryImpl struct {
	db *sql.DB
}

// NewProductRepository cria uma nova instância do repositório
func NewProductRepository(db *sql.DB) repositories.ProductRepository {
	return &ProductRepositoryImpl{
		db: db,
	}
}

// GetByEAN busca produtos pelo código EAN (código de barras)
func (r *ProductRepositoryImpl) GetByEAN(ean string) ([]entities.Product, error) {
	// O código de barras está na tabela EmbalagemProduto, não na tabela Produto
	query := `
		SELECT DISTINCT p.IDPRODUTO, e.CODIGOBARRAS 
		FROM Produto p
		INNER JOIN EmbalagemProduto e ON p.IDPRODUTO = e.IDPRODUTO
		WHERE e.CODIGOBARRAS = :1
	`

	rows, err := r.db.Query(query, ean)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produto por EAN: %w", err)
	}
	defer rows.Close()

	var products []entities.Product
	for rows.Next() {
		var product entities.Product
		if err := rows.Scan(&product.ID, &product.EAN); err != nil {
			return nil, fmt.Errorf("erro ao escanear produto: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar produtos: %w", err)
	}

	return products, nil
}

// SaveIntegrationStaging grava a integração do produto no staging chamando a stored procedure
func (r *ProductRepositoryImpl) SaveIntegrationStaging(dealerID, productID int) error {
	start := time.Now()

	// Incrementa contador de chamadas
	count := atomic.AddInt64(&spCallCount, 1)

	// Log a cada 1000 chamadas ou a cada 10 segundos
	now := time.Now().Unix()
	lastLog := atomic.LoadInt64(&lastLogTime)
	if count%1000 == 0 || (now-lastLog) >= 10 {
		atomic.StoreInt64(&lastLogTime, now)
		avgTime := float64(atomic.LoadInt64(&spTotalTime)) / float64(count) / 1000000.0
		errors := atomic.LoadInt64(&spErrorCount)
		log.Printf("📊 SP Stats: %d chamadas | Média: %.2fms | Erros: %d", count, avgTime, errors)
	}

	// Chama a stored procedure SP_GRAVARINTEGRACAOPRODUTOSTAGING
	query := `
		BEGIN
			SP_GRAVARINTEGRACAOPRODUTOSTAGING(:p_idRevendedor, :p_idProduto);
		END;
	`

	_, err := r.db.Exec(query, dealerID, productID)

	// Registra tempo de execução
	elapsed := time.Since(start).Nanoseconds()
	atomic.AddInt64(&spTotalTime, elapsed)

	if err != nil {
		atomic.AddInt64(&spErrorCount, 1)
		return fmt.Errorf("erro ao gravar integração produto staging: %w", err)
	}

	return nil
}
