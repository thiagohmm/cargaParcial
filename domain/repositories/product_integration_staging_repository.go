package repositories

import "github.thiagohmm.com.br/cargaparcial/domain/entities"

// ProductIntegrationStagingRepository define as operações de acesso a dados para ProductIntegrationStaging
type ProductIntegrationStagingRepository interface {
	GetByProductAndDealer(productID, dealerID int) (*entities.ProductIntegrationStaging, error)
}
