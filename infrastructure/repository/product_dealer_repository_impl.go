package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.thiagohmm.com.br/cargaparcial/domain/entities"
	"github.thiagohmm.com.br/cargaparcial/domain/repositories"
)

// ProductDealerRepositoryImpl implementa o ProductDealerRepository
type ProductDealerRepositoryImpl struct {
	db           *sql.DB
	stmtExists   *sql.Stmt
	stmtCreate   *sql.Stmt
}

// NewProductDealerRepository cria uma nova instância do repositório
func NewProductDealerRepository(db *sql.DB) repositories.ProductDealerRepository {
	repo := &ProductDealerRepositoryImpl{
		db: db,
	}
	
	// Pré-compilar query de verificação de existência
	var err error
	repo.stmtExists, err = db.Prepare(`SELECT COUNT(*) FROM ProdutoRevendedor WHERE IdProduto = :1 AND IdRevendedor = :2`)
	if err != nil {
		panic(fmt.Sprintf("Erro ao preparar statement Exists: %v", err))
	}
	
	// Pré-compilar query de insert único
	repo.stmtCreate, err = db.Prepare(`
		INSERT INTO ProdutoRevendedor (IdProduto, IdRevendedor, StatusProdutoRevendedor)
		VALUES (:1, :2, :3)
	`)
	if err != nil {
		panic(fmt.Sprintf("Erro ao preparar statement Create: %v", err))
	}
	
	return repo
}

// Exists verifica se existe uma relação entre produto e revendedor
func (r *ProductDealerRepositoryImpl) Exists(productID, dealerID int) (bool, error) {
	var count int
	err := r.stmtExists.QueryRow(productID, dealerID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência de ProductDealer: %w", err)
	}

	return count > 0, nil
}

// Create cria uma nova relação entre produto e revendedor
func (r *ProductDealerRepositoryImpl) Create(productDealer *entities.ProductDealer) error {
	_, err := r.stmtCreate.Exec(productDealer.ProductID, productDealer.DealerID, productDealer.IsActive)
	if err != nil {
		return fmt.Errorf("erro ao criar ProductDealer: %w", err)
	}

	return nil
}

// CreateBatch cria múltiplas relações em batch (mais eficiente)
func (r *ProductDealerRepositoryImpl) CreateBatch(productDealers []*entities.ProductDealer) error {
	if len(productDealers) == 0 {
		return nil
	}

	// Para Oracle, vamos usar INSERT ALL
	// INSERT ALL
	//   INTO ProdutoRevendedor VALUES (?, ?, ?)
	//   INTO ProdutoRevendedor VALUES (?, ?, ?)
	// SELECT 1 FROM DUAL
	
	const batchSize = 100 // Oracle tem limite de 1000 binds, 100 * 3 = 300 é seguro
	
	for i := 0; i < len(productDealers); i += batchSize {
		end := i + batchSize
		if end > len(productDealers) {
			end = len(productDealers)
		}
		
		batch := productDealers[i:end]
		
		var query strings.Builder
		query.WriteString("INSERT ALL\n")
		
		args := make([]interface{}, 0, len(batch)*3)
		for idx, pd := range batch {
			offset := idx * 3
			query.WriteString(fmt.Sprintf("  INTO ProdutoRevendedor (IdProduto, IdRevendedor, StatusProdutoRevendedor) VALUES (:%d, :%d, :%d)\n", 
				offset+1, offset+2, offset+3))
			args = append(args, pd.ProductID, pd.DealerID, pd.IsActive)
		}
		
		query.WriteString("SELECT 1 FROM DUAL")
		
		_, err := r.db.Exec(query.String(), args...)
		if err != nil {
			return fmt.Errorf("erro ao criar ProductDealers em batch: %w", err)
		}
	}

	return nil
}
