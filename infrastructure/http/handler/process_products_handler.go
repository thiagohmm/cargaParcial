package handler

import (
	"encoding/json"
	"net/http"

	"github.thiagohmm.com.br/cargaparcial/usecase"
	"github.thiagohmm.com.br/cargaparcial/usecase/dto"
)

// ProcessProductsHandler gerencia as requisições HTTP para processar produtos
type ProcessProductsHandler struct {
	useCase *usecase.ProcessProductsUseCase
}

// NewProcessProductsHandler cria uma nova instância do handler
func NewProcessProductsHandler(useCase *usecase.ProcessProductsUseCase) *ProcessProductsHandler {
	return &ProcessProductsHandler{
		useCase: useCase,
	}
}

// Handle processa a requisição HTTP
func (h *ProcessProductsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var input dto.ProcessProductsInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validar entrada
	if len(input.IBMCodes) == 0 {
		http.Error(w, "Lista de códigos IBM não pode estar vazia", http.StatusBadRequest)
		return
	}

	if len(input.ProductCodes) == 0 {
		http.Error(w, "Lista de códigos de produto não pode estar vazia", http.StatusBadRequest)
		return
	}

	// Executar use case
	output, err := h.useCase.Execute(input)
	if err != nil {
		http.Error(w, "Erro ao processar produtos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retornar resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, "Erro ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
