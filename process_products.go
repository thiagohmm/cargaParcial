package main

import (
	"fmt"
	"log"
)

// Estruturas de dados
type ProductResult struct {
	IdRevendedor *int   `json:"IdRevendedor"`
	IdProduto    *int   `json:"IdProduto"`
	EAN          string `json:"EAN,omitempty"`
	Status       string `json:"Status"`
	Motivo       string `json:"Motivo,omitempty"`
}

type Result struct {
	ArrayOk   []ProductResult `json:"arrayOk"`
	ArrayFail []ProductResult `json:"arrayFail"`
}

type Data struct {
	IBM    []string `json:"IBM"`
	Codigo []string `json:"codigo"`
}

type Dealer struct {
	IdRevendedor int
}

type Product struct {
	IDPRODUTO int
}

type ProductDealerCreate struct {
	IdProduto               int
	IdRevendedor            int
	StatusProdutoRevendedor bool
}

// Interfaces para queries (você precisará implementar essas funções)
type DealerQuery interface {
	GetDealerByIbm(ibm string) (*Dealer, error)
}

type ProductQuery interface {
	GravarIntegracaoProdutoStaging(idRevendedor, idProduto int) error
}

type ProductDealerQuery interface {
	ExistsDealer(idProduto, idRevendedor int) (bool, error)
	Create(data ProductDealerCreate) error
}

type ProductIntegrationStagingQuery interface {
	GetByProductIntegrationStaging(idProduto, idRevendedor int) (interface{}, error)
}

// Função principal
func ProcessProducts(
	data Data,
	dealerQuery DealerQuery,
	productQuery ProductQuery,
	productDealerQuery ProductDealerQuery,
	productIntegrationStagingQuery ProductIntegrationStagingQuery,
	listProductByEAN func(ean string) ([]Product, error),
	sendToQueue func(msg string),
) (Result, error) {
	resultProduct := Result{
		ArrayOk:   make([]ProductResult, 0),
		ArrayFail: make([]ProductResult, 0),
	}

	for _, ibm := range data.IBM {
		ibmValue := ibm
		if ibmValue == "" {
			ibmValue = "0"
		}

		revendedor, err := dealerQuery.GetDealerByIbm(ibmValue)
		if err != nil {
			log.Printf("Erro ao buscar revendedor por IBM %s: %v", ibmValue, err)
			continue
		}

		for _, codigo := range data.Codigo {
			produto := codigo

			// Buscar produto por EAN
			idProduct, err := listProductByEAN(produto)

			// Validar se o produto foi encontrado
			if err != nil || len(idProduct) == 0 {
				idRev := revendedor.IdRevendedor
				resultProduct.ArrayFail = append(resultProduct.ArrayFail, ProductResult{
					IdRevendedor: &idRev,
					IdProduto:    nil,
					EAN:          produto,
					Status:       "fail",
					Motivo:       "Produto não encontrado pelo EAN",
				})
				continue
			}

			// Verificar se existe ProductDealer
			checkProductDealer, err := productDealerQuery.ExistsDealer(idProduct[0].IDPRODUTO, revendedor.IdRevendedor)
			if err != nil {
				log.Printf("Erro ao verificar ProductDealer: %v", err)
				continue
			}

			if !checkProductDealer {
				err := productDealerQuery.Create(ProductDealerCreate{
					IdProduto:               idProduct[0].IDPRODUTO,
					IdRevendedor:            revendedor.IdRevendedor,
					StatusProdutoRevendedor: true,
				})
				if err != nil {
					log.Printf("Erro ao criar ProductDealer: %v", err)
					continue
				}
			}

			// Gravar integração produto staging
			err = productQuery.GravarIntegracaoProdutoStaging(revendedor.IdRevendedor, idProduct[0].IDPRODUTO)
			if err != nil {
				log.Printf("Erro ao gravar integração produto staging: %v", err)
			}

			// Buscar productIntegrationStaging
			productIntegrationStaging, err := productIntegrationStagingQuery.GetByProductIntegrationStaging(
				idProduct[0].IDPRODUTO,
				revendedor.IdRevendedor,
			)

			fmt.Printf("%v, %d, %d\n", productIntegrationStaging, revendedor.IdRevendedor, idProduct[0].IDPRODUTO)

			idRev := revendedor.IdRevendedor
			idProd := idProduct[0].IDPRODUTO

			if err == nil && productIntegrationStaging != nil {
				resultProduct.ArrayOk = append(resultProduct.ArrayOk, ProductResult{
					IdRevendedor: &idRev,
					IdProduto:    &idProd,
					Status:       "ok",
				})
			} else {
				resultProduct.ArrayFail = append(resultProduct.ArrayFail, ProductResult{
					IdRevendedor: &idRev,
					IdProduto:    &idProd,
					Status:       "fail",
				})
			}
		}
	}

	sendToQueue("mover")
	return resultProduct, nil
}
