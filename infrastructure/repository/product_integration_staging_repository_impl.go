package repository

import (
	"database/sql"
	"fmt"

	"github.thiagohmm.com.br/cargaparcial/domain/entities"
	"github.thiagohmm.com.br/cargaparcial/domain/repositories"
)

// ProductIntegrationStagingRepositoryImpl implementa o ProductIntegrationStagingRepository
type ProductIntegrationStagingRepositoryImpl struct {
	db *sql.DB
}

// NewProductIntegrationStagingRepository cria uma nova instância do repositório
func NewProductIntegrationStagingRepository(db *sql.DB) repositories.ProductIntegrationStagingRepository {
	return &ProductIntegrationStagingRepositoryImpl{
		db: db,
	}
}

// GetByProductAndDealer busca um registro de integração por produto e revendedor
func (r *ProductIntegrationStagingRepositoryImpl) GetByProductAndDealer(productID, dealerID int) (*entities.ProductIntegrationStaging, error) {
	query := `
		SELECT IdProduto, IdRevendedor 
		FROM IntegracaoProdutoStaging 
		WHERE IdProduto = :1 AND IdRevendedor = :2
	`

	var staging entities.ProductIntegrationStaging
	err := r.db.QueryRow(query, productID, dealerID).Scan(&staging.ProductID, &staging.DealerID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar ProductIntegrationStaging: %w", err)
	}

	return &staging, nil
}
