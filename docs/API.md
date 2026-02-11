# API Documentation

## Base URL

```
http://localhost:8080
```

## Endpoints

### POST /api/process-products

Processa produtos para múltiplos revendedores baseado em códigos IBM e EAN.

#### Request

**Headers:**

```
Content-Type: application/json
```

**Body:**

```json
{
  "IBM": ["string"], // Array de códigos IBM dos revendedores
  "codigo": ["string"] // Array de códigos EAN dos produtos
}
```

**Exemplo:**

```json
{
  "IBM": ["IBM001", "IBM002"],
  "codigo": ["7891234567890", "7891234567891", "7891234567892"]
}
```

#### Response

**Status Code:** `200 OK`

**Body:**

```json
{
  "arrayOk": [
    {
      "IdRevendedor": 1,
      "IdProduto": 100,
      "Status": "ok"
    }
  ],
  "arrayFail": [
    {
      "IdRevendedor": 2,
      "IdProduto": null,
      "EAN": "7891234567891",
      "Status": "fail",
      "Motivo": "Produto não encontrado pelo EAN"
    }
  ]
}
```

#### Response Fields

**arrayOk** - Array de produtos processados com sucesso

- `IdRevendedor` (int): ID do revendedor
- `IdProduto` (int): ID do produto
- `Status` (string): Status do processamento ("ok")

**arrayFail** - Array de produtos que falhou no processamento

- `IdRevendedor` (int|null): ID do revendedor
- `IdProduto` (int|null): ID do produto
- `EAN` (string): Código EAN do produto (quando aplicável)
- `Status` (string): Status do processamento ("fail")
- `Motivo` (string): Motivo da falha

#### Possíveis Motivos de Falha

1. `"Produto não encontrado pelo EAN"` - O código EAN não existe no banco de dados
2. `"Erro ao verificar relação produto-revendedor"` - Erro ao consultar ProductDealer
3. `"Erro ao criar relação produto-revendedor"` - Erro ao criar registro ProductDealer
4. `"Erro ao processar integração"` - Erro geral no processamento da integração

#### Error Responses

**400 Bad Request**

```json
{
  "error": "Lista de códigos IBM não pode estar vazia"
}
```

**400 Bad Request**

```json
{
  "error": "Lista de códigos de produto não pode estar vazia"
}
```

**405 Method Not Allowed**

```json
{
  "error": "Método não permitido"
}
```

**500 Internal Server Error**

```json
{
  "error": "Erro ao processar produtos: [detalhes do erro]"
}
```

## Exemplos de Uso

### cURL

```bash
curl -X POST http://localhost:8080/api/process-products \
  -H "Content-Type: application/json" \
  -d '{
    "IBM": ["IBM001"],
    "codigo": ["7891234567890"]
  }'
```

### JavaScript (Fetch)

```javascript
fetch('http://localhost:8080/api/process-products', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    IBM: ['IBM001'],
    codigo: ['7891234567890'],
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => console.error('Error:', error));
```

### Python (requests)

```python
import requests

url = 'http://localhost:8080/api/process-products'
data = {
    'IBM': ['IBM001'],
    'codigo': ['7891234567890']
}

response = requests.post(url, json=data)
print(response.json())
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

func main() {
    url := "http://localhost:8080/api/process-products"

    payload := map[string]interface{}{
        "IBM":    []string{"IBM001"},
        "codigo": []string{"7891234567890"},
    }

    jsonData, _ := json.Marshal(payload)

    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    "IBM":    []string{"IBM001"},
        "codigo": []string{"7891234567890"},
    }
    
    jsonData, _ := json.Marshal(payload)
    
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
}
```

## Fluxo de Processamento

1. **Validação de Entrada**
   - Verifica se os arrays IBM e codigo não estão vazios

2. **Para cada código IBM:**
   - Busca o revendedor no banco de dados
   - Se não encontrado, pula para o próximo

3. **Para cada código de produto (EAN):**
   - Busca o produto pelo EAN
   - Se não encontrado, adiciona em `arrayFail` com motivo
   - Verifica se já existe relação ProductDealer
   - Se não existe, cria a relação
   - Grava no staging de integração
   - Verifica se a integração foi bem-sucedida
   - Adiciona em `arrayOk` ou `arrayFail` conforme resultado

4. **Pós-processamento:**
   - Envia mensagem "mover" para a fila
   - Retorna resultado com arrays de sucesso e falha

## Notas Importantes

- O processamento é feito de forma síncrona
- Cada combinação IBM x Código é processada independentemente
- Falhas em produtos individuais não interrompem o processamento dos demais
- A mensagem é enviada para a fila independentemente de sucessos ou falhas
- Logs detalhados são gerados durante o processamento
