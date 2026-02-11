#!/usr/bin/env python3
"""
Script para criar um arquivo Excel de exemplo para teste
"""

try:
    from openpyxl import Workbook
except ImportError:
    print("Instalando openpyxl...")
    import subprocess
    subprocess.check_call(['pip3', 'install', 'openpyxl'])
    from openpyxl import Workbook

# Dados de exemplo
dados = [
    ['IMBLOJA', 'CODIGOBARRAS'],
    ['0001002154', '7896050201756'],
    ['0001002154', '7898080070050'],
    ['0001006393', '070330717534'],
    ['0001006393', '0735202909010'],
    ['0001006393', '0736532327543'],
    ['0001006393', '0798190262291'],
    ['0001006393', '08000500121467'],
    ['0001006393', '095188794506'],
    ['0001006393', '4893993367528'],
]

# Criar workbook
wb = Workbook()
ws = wb.active
ws.title = "Dados"

# Adicionar dados
for row in dados:
    ws.append(row)

# Salvar arquivo
wb.save('dados_exemplo.xlsx')
print("✓ Arquivo dados_exemplo.xlsx criado com sucesso!")
print(f"  - {len(dados)-1} linhas de dados")
print(f"  - Colunas: IMBLOJA, CODIGOBARRAS")
