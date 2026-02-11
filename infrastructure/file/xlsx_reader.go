package file

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// XLSXData representa os dados lidos do arquivo XLSX
type XLSXData struct {
	IBMCodes     []string
	ProductCodes []string
	// Novo: mantém o relacionamento IBM -> Produtos
	IBMToProducts map[string][]string
}

// ReadXLSX lê um arquivo XLSX e extrai os dados das colunas IMBLOJA e CODIGOBARRAS
func ReadXLSX(filename string) (*XLSXData, error) {
	// Abrir o arquivo XLSX
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo XLSX: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Erro ao fechar arquivo: %v\n", err)
		}
	}()

	// Obter a primeira planilha
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("nenhuma planilha encontrada no arquivo")
	}

	sheetName := sheets[0]

	// Ler todas as linhas
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler linhas: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("arquivo vazio")
	}

	// Encontrar índices das colunas IMBLOJA e CODIGOBARRAS
	header := rows[0]
	imbLojaIdx := -1
	codigoBarrasIdx := -1

	for i, col := range header {
		colUpper := strings.ToUpper(strings.TrimSpace(col))
		if colUpper == "IMBLOJA" {
			imbLojaIdx = i
		} else if colUpper == "CODIGOBARRAS" {
			codigoBarrasIdx = i
		}
	}

	if imbLojaIdx == -1 {
		return nil, fmt.Errorf("coluna IMBLOJA não encontrada no cabeçalho")
	}
	if codigoBarrasIdx == -1 {
		return nil, fmt.Errorf("coluna CODIGOBARRAS não encontrada no cabeçalho")
	}

	// Mapear IBM codes para seus produtos
	ibmToProducts := make(map[string][]string)

	// Processar linhas de dados (pular cabeçalho)
	for i := 1; i < len(rows); i++ {
		row := rows[i]

		// Verificar se a linha tem colunas suficientes
		if len(row) <= imbLojaIdx || len(row) <= codigoBarrasIdx {
			continue
		}

		ibmCode := strings.TrimSpace(row[imbLojaIdx])
		productCode := strings.TrimSpace(row[codigoBarrasIdx])

		// Ignorar linhas vazias
		if ibmCode == "" || productCode == "" {
			continue
		}

		// Adicionar ao mapa
		if _, exists := ibmToProducts[ibmCode]; !exists {
			ibmToProducts[ibmCode] = make([]string, 0)
		}
		ibmToProducts[ibmCode] = append(ibmToProducts[ibmCode], productCode)
	}

	// Converter mapa para listas
	// O usecase espera todas as combinações, então vamos criar listas únicas
	ibmCodesMap := make(map[string]bool)
	productCodesMap := make(map[string]bool)

	for ibm, products := range ibmToProducts {
		ibmCodesMap[ibm] = true
		for _, product := range products {
			productCodesMap[product] = true
		}
	}

	// Converter maps para slices
	ibmCodes := make([]string, 0, len(ibmCodesMap))
	for ibm := range ibmCodesMap {
		ibmCodes = append(ibmCodes, ibm)
	}

	productCodes := make([]string, 0, len(productCodesMap))
	for product := range productCodesMap {
		productCodes = append(productCodes, product)
	}

	return &XLSXData{
		IBMCodes:      ibmCodes,
		ProductCodes:  productCodes,
		IBMToProducts: ibmToProducts, // Mantém o relacionamento original
	}, nil
}

// ReadXLSXPairs lê um arquivo XLSX e retorna pares específicos de IBM e Produto
// Esta função mantém as combinações exatas do arquivo, sem criar todas as combinações
func ReadXLSXPairs(filename string) (map[string][]string, error) {
	// Abrir o arquivo XLSX
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo XLSX: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Erro ao fechar arquivo: %v\n", err)
		}
	}()

	// Obter a primeira planilha
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("nenhuma planilha encontrada no arquivo")
	}

	sheetName := sheets[0]

	// Ler todas as linhas
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler linhas: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("arquivo vazio")
	}

	// Encontrar índices das colunas IMBLOJA e CODIGOBARRAS
	header := rows[0]
	imbLojaIdx := -1
	codigoBarrasIdx := -1

	for i, col := range header {
		colUpper := strings.ToUpper(strings.TrimSpace(col))
		if colUpper == "IMBLOJA" {
			imbLojaIdx = i
		} else if colUpper == "CODIGOBARRAS" {
			codigoBarrasIdx = i
		}
	}

	if imbLojaIdx == -1 {
		return nil, fmt.Errorf("coluna IMBLOJA não encontrada no cabeçalho")
	}
	if codigoBarrasIdx == -1 {
		return nil, fmt.Errorf("coluna CODIGOBARRAS não encontrada no cabeçalho")
	}

	// Mapear IBM codes para seus produtos (mantendo duplicatas)
	ibmToProducts := make(map[string][]string)

	// Processar linhas de dados (pular cabeçalho)
	for i := 1; i < len(rows); i++ {
		row := rows[i]

		// Verificar se a linha tem colunas suficientes
		if len(row) <= imbLojaIdx || len(row) <= codigoBarrasIdx {
			continue
		}

		ibmCode := strings.TrimSpace(row[imbLojaIdx])
		productCode := strings.TrimSpace(row[codigoBarrasIdx])

		// Ignorar linhas vazias
		if ibmCode == "" || productCode == "" {
			continue
		}

		// Adicionar ao mapa
		if _, exists := ibmToProducts[ibmCode]; !exists {
			ibmToProducts[ibmCode] = make([]string, 0)
		}
		ibmToProducts[ibmCode] = append(ibmToProducts[ibmCode], productCode)
	}

	return ibmToProducts, nil
}
