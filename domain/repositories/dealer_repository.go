package repositories

import "github.thiagohmm.com.br/cargaparcial/domain/entities"

// DealerRepository define as operações de acesso a dados para Dealer
type DealerRepository interface {
	GetByIBM(ibm string) (*entities.Dealer, error)
}
