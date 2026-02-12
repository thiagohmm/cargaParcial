package repository

import (
	"database/sql"
	"fmt"

	"github.thiagohmm.com.br/cargaparcial/domain/entities"
	"github.thiagohmm.com.br/cargaparcial/domain/repositories"
)

// ProductIntegrationStagingRepositoryImpl implementa o ProductIntegrationStagingRepository
type ProductIntegrationStagingRepositoryImpl struct {
	db                         *sql.DB
	stmtGetByProductAndDealer  *sql.Stmt
}

// NewProductIntegrationStagingRepository cria uma nova instância do repositório
func NewProductIntegrationStagingRepository(db *sql.DB) repositories.ProductIntegrationStagingRepository {
	repo := &ProductIntegrationStagingRepositoryImpl{
		db: db,
	}
	
	// Pré-compilar query de busca
	var err error
	repo.stmtGetByProductAndDealer, err = db.Prepare(`
		SELECT IdProduto, IdRevendedor 
		FROM IntegracaoProdutoStaging 
		WHERE IdProduto = :1 AND IdRevendedor = :2
	`)
	if err != nil {
		panic(fmt.Sprintf("Erro ao preparar statement GetByProductAndDealer: %v", err))
	}
	
	return repo
}

// GetByProductAndDealer busca um registro de integração por produto e revendedor
func (r *ProductIntegrationStagingRepositoryImpl) GetByProductAndDealer(productID, dealerID int) (*entities.ProductIntegrationStaging, error) {
	var staging entities.ProductIntegrationStaging
	err := r.stmtGetByProductAndDealer.QueryRow(productID, dealerID).Scan(&staging.ProductID, &staging.DealerID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar ProductIntegrationStaging: %w", err)
	}

	return &staging, nil
}
