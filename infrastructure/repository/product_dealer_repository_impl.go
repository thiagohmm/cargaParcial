package repository

import (
	"database/sql"
	"fmt"

	"github.thiagohmm.com.br/cargaparcial/domain/entities"
	"github.thiagohmm.com.br/cargaparcial/domain/repositories"
)

// ProductDealerRepositoryImpl implementa o ProductDealerRepository
type ProductDealerRepositoryImpl struct {
	db *sql.DB
}

// NewProductDealerRepository cria uma nova instância do repositório
func NewProductDealerRepository(db *sql.DB) repositories.ProductDealerRepository {
	return &ProductDealerRepositoryImpl{
		db: db,
	}
}

// Exists verifica se existe uma relação entre produto e revendedor
func (r *ProductDealerRepositoryImpl) Exists(productID, dealerID int) (bool, error) {
	query := `SELECT COUNT(*) FROM ProdutoRevendedor WHERE IdProduto = :1 AND IdRevendedor = :2`

	var count int
	err := r.db.QueryRow(query, productID, dealerID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência de ProductDealer: %w", err)
	}

	return count > 0, nil
}

// Create cria uma nova relação entre produto e revendedor
func (r *ProductDealerRepositoryImpl) Create(productDealer *entities.ProductDealer) error {
	query := `
		INSERT INTO ProdutoRevendedor (IdProduto, IdRevendedor, StatusProdutoRevendedor)
		VALUES (:1, :2, :3)
	`

	_, err := r.db.Exec(query, productDealer.ProductID, productDealer.DealerID, productDealer.IsActive)
	if err != nil {
		return fmt.Errorf("erro ao criar ProductDealer: %w", err)
	}

	return nil
}
