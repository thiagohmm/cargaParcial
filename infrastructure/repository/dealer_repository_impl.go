package repository

import (
	"database/sql"
	"fmt"

	"github.thiagohmm.com.br/cargaparcial/domain/entities"
	"github.thiagohmm.com.br/cargaparcial/domain/repositories"
)

// DealerRepositoryImpl implementa o DealerRepository
type DealerRepositoryImpl struct {
	db *sql.DB
}

// NewDealerRepository cria uma nova instância do repositório
func NewDealerRepository(db *sql.DB) repositories.DealerRepository {
	return &DealerRepositoryImpl{
		db: db,
	}
}

// GetByIBM busca um revendedor pelo código IBM
func (r *DealerRepositoryImpl) GetByIBM(ibm string) (*entities.Dealer, error) {
	query := `SELECT IdRevendedor, CodigoIBM FROM Revendedor WHERE CodigoIBM = :1`

	var dealer entities.Dealer
	err := r.db.QueryRow(query, ibm).Scan(&dealer.ID, &dealer.IBM)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("revendedor não encontrado para IBM: %s", ibm)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar revendedor: %w", err)
	}

	return &dealer, nil
}
