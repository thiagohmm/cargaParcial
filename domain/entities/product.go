package entities

// Product representa um produto no sistema
type Product struct {
	ID  int
	EAN string
}

// Dealer representa um revendedor no sistema
type Dealer struct {
	ID  int
	IBM string
}

// ProductDealer representa a relação entre produto e revendedor
type ProductDealer struct {
	ProductID int
	DealerID  int
	IsActive  bool
}

// ProductIntegrationStaging representa o staging de integração de produto
type ProductIntegrationStaging struct {
	ProductID int
	DealerID  int
}
