package repository

import (
	"database/sql"
	"fmt"

	"github.thiagohmm.com.br/cargaparcial/domain/entities"
	"github.thiagohmm.com.br/cargaparcial/domain/repositories"
)

// DealerRepositoryImpl implementa o DealerRepository
type DealerRepositoryImpl struct {
	db            *sql.DB
	stmtGetByIBM  *sql.Stmt
}

// NewDealerRepository cria uma nova instância do repositório
func NewDealerRepository(db *sql.DB) repositories.DealerRepository {
	repo := &DealerRepositoryImpl{
		db: db,
	}
	
	// Pré-compilar query de busca por IBM
	var err error
	repo.stmtGetByIBM, err = db.Prepare(`SELECT IdRevendedor, CodigoIBM FROM Revendedor WHERE CodigoIBM = :1`)
	if err != nil {
		panic(fmt.Sprintf("Erro ao preparar statement GetByIBM: %v", err))
	}
	
	return repo
}

// GetByIBM busca um revendedor pelo código IBM
func (r *DealerRepositoryImpl) GetByIBM(ibm string) (*entities.Dealer, error) {
	var dealer entities.Dealer
	err := r.stmtGetByIBM.QueryRow(ibm).Scan(&dealer.ID, &dealer.IBM)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("revendedor não encontrado para IBM: %s", ibm)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar revendedor: %w", err)
	}

	return &dealer, nil
}
