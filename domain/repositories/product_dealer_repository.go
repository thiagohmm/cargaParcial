package repositories

import "github.thiagohmm.com.br/cargaparcial/domain/entities"

// ProductDealerRepository define as operações de acesso a dados para ProductDealer
type ProductDealerRepository interface {
	Exists(productID, dealerID int) (bool, error)
	Create(productDealer *entities.ProductDealer) error
}
