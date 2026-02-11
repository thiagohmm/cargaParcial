.PHONY: help build run test clean deps

# Variáveis
APP_NAME=cargaparcial
BUILD_DIR=bin
MAIN_PATH=cmd/api/main.go

help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## Instala as dependências do projeto
	@echo "Instalando dependências..."
	go mod download
	go mod tidy

build: ## Compila a aplicação
	@echo "Compilando aplicação..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build concluído: $(BUILD_DIR)/$(APP_NAME)"

run: ## Executa a aplicação
	@echo "Executando aplicação..."
	go run $(MAIN_PATH)

test: ## Executa os testes
	@echo "Executando testes..."
	go test -v ./...

test-coverage: ## Executa os testes com cobertura
	@echo "Executando testes com cobertura..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Relatório de cobertura gerado: coverage.html"

clean: ## Remove arquivos de build
	@echo "Limpando arquivos de build..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Limpeza concluída"

fmt: ## Formata o código
	@echo "Formatando código..."
	go fmt ./...

lint: ## Executa o linter
	@echo "Executando linter..."
	golangci-lint run

docker-build: ## Constrói a imagem Docker
	@echo "Construindo imagem Docker..."
	docker build -t $(APP_NAME):latest .

docker-run: ## Executa o container Docker
	@echo "Executando container Docker..."
	docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

.DEFAULT_GOAL := help
