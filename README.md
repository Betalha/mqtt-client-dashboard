# Sistema de Monitoramento de Sensores (Temperatura e Umidade)

## Descrição

Este é um sistema de monitoramento de sensores de temperatura e umidade desenvolvido em Go. O sistema utiliza o protocolo MQTT para receber dados de sensores, armazena os dados em um arquivo CSV e exibe os dados em tempo real em uma interface web com gráficos interativos.

## Funcionalidades

- Recebimento de dados de sensores via protocolo MQTT
- Armazenamento dos dados em arquivo CSV
- Interface web em tempo real com gráficos usando Chart.js
- WebSocket para atualização instantânea dos dados na interface
- Visualização de temperatura e umidade em tempo real
- Histórico de leituras com gráficos lineares

## Arquitetura

O sistema é composto por:

- **Backend (Go)**: Servidor HTTP que também atua como cliente MQTT para receber dados de sensores
- **Frontend (HTML/CSS/JS)**: Interface web para visualização dos dados em tempo real
- **Protocolo MQTT**: Utilizado para comunicação entre sensores e servidor
- **WebSocket**: Conexão em tempo real entre servidor e frontend
- **CSV**: Armazenamento persistente dos dados recebidos

## Estrutura de Dados

O sistema trabalha com mensagens JSON no seguinte formato:
```json
{
  "id": "string",
  "temperatura": 23.5,
  "umidade": 65.0,
  "timestamp": "2025-11-27T10:15:30Z",
  "controle": "string"
}
```

## Estrutura de Arquivos

```
sis cyberfis/
├── main.go              # Código Go do servidor backend
├── sensor_data.csv      # Arquivo CSV com dados armazenados
├── go.mod               # Dependências Go
├── go.sum               # Checksums das dependências
└── static/
    ├── index.html       # Página principal da interface web
    ├── script.js        # Lógica JavaScript para gráficos e WebSocket
    └── style.css        # Estilos da interface web
```

## Tecnologias Utilizadas

- **Go** (v1.24+)
- **MQTT** (Cliente Eclipse Paho)
- **WebSocket** (Gorilla WebSocket)
- **HTML5/CSS3**
- **JavaScript** (Chart.js para gráficos)

## Configurações

### Dependências Go

- `github.com/eclipse/paho.mqtt.golang`: Cliente MQTT
- `github.com/gorilla/websocket`: WebSocket
- `encoding/csv`: Manipulação de CSV
- `net/http`: Servidor HTTP
- `encoding/json`: Manipulação de JSON

### Broker MQTT

O sistema está configurado para se conectar ao broker público: `tcp://broker.hivemq.com:1883`
Topico de escuta: `sensor/client`

## Instalação e Execução

1. Clone ou baixe o repositório

2. Instale as dependências Go:
```bash
go mod tidy
```

3. Execute o servidor:
```bash
go run main.go
```

4. Acesse a interface web no navegador:
```
http://localhost:8080
```

## Como Funciona

1. O servidor Go inicia e se conecta ao broker MQTT
2. O servidor se inscreve no tópico `sensor/client`
3. O servidor inicia o servidor HTTP para servir a interface web
4. O servidor abre ou cria o arquivo CSV `sensor_data.csv`
5. Quando o frontend se conecta, estabelece uma conexão WebSocket
6. Quando dados de sensores chegam via MQTT, são:
   - Armazenados no CSV
   - Enviados via WebSocket para todos os clientes conectados
7. O frontend atualiza os dados em tempo real e os exibe graficamente

## Configuração de Sensores

Para enviar dados para este sistema, configure seus sensores para publicar mensagens JSON no tópico MQTT `sensor/client` no broker `broker.hivemq.com`.

## Arquivo de Dados

Os dados recebidos são armazenados no arquivo `sensor_data.csv` com o formato:
```
timestamp,temperatura,umidade
2025-11-27T10:15:30Z,24.70,58.30
```

## Interface Web

A interface web oferece:

- Leitura em tempo real de temperatura e umidade
- Gráficos dinâmicos com histórico das últimas 25 leituras
- Design responsivo e visualização intuitiva dos dados

