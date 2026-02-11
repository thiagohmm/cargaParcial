package dto

// ProcessProductsInput representa os dados de entrada para processar produtos
type ProcessProductsInput struct {
	IBMCodes      []string            `json:"IBM"`
	ProductCodes  []string            `json:"codigo"`
	IBMToProducts map[string][]string `json:"-"` // Relacionamento IBM -> Produtos (não vem do JSON)
}

// ProductResultDTO representa o resultado do processamento de um produto
type ProductResultDTO struct {
	DealerID  *int   `json:"IdRevendedor"`
	ProductID *int   `json:"IdProduto"`
	EAN       string `json:"EAN,omitempty"`
	Status    string `json:"Status"`
	Reason    string `json:"Motivo,omitempty"`
}

// ProcessProductsOutput representa o resultado do processamento
type ProcessProductsOutput struct {
	SuccessList []ProductResultDTO `json:"arrayOk"`
	FailureList []ProductResultDTO `json:"arrayFail"`
}
