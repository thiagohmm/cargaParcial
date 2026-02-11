package repositories

import "github.thiagohmm.com.br/cargaparcial/domain/entities"

// ProductRepository define as operações de acesso a dados para Product
type ProductRepository interface {
	GetByEAN(ean string) ([]entities.Product, error)
	SaveIntegrationStaging(dealerID, productID int) error
}
