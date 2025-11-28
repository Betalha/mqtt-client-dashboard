package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message)
	mutex     = sync.RWMutex{}
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	csvFile   *os.File
	csvWriter *csv.Writer
	csvMutex  sync.Mutex // evitar escritas concorrentes
)

type Message struct {
	ID          string  `json:"ID"`
	Temperatura float64 `json:"temperatura"`
	Umidade     float64 `json:"umidade"`
	Timestamp   string  `json:"timestamp"`
	Controle    bool    `json:"controle"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(26))
	}
	return string(bytes)
}

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://broker.hivemq.com:1883")
	opts.SetClientID("go-web-mqtt-client-" + randString(5))
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(true)

	opts.OnConnect = func(c mqtt.Client) {
		topic := "sensor/client"
		token := c.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
			var data Message
			if err := json.Unmarshal(msg.Payload(), &data); err != nil {
				log.Printf("Erro ao decodificar payload MQTT: %v", err)
				return
			}
			csvMutex.Lock()
			csvWriter.Write([]string{
				data.Timestamp,
				fmt.Sprintf("%.2f", data.Temperatura),
				fmt.Sprintf("%.2f", data.Umidade),
			})
			csvWriter.Flush()
			csvMutex.Unlock()

			log.Printf("Recebido: ID=%s, T=%.1f°C, U=%.1f%%, TS=%s, Controle=%s",
				data.ID, data.Temperatura, data.Umidade, data.Timestamp, data.Controle)

			broadcast <- data
		})
		token.Wait()
		if token.Error() != nil {
			log.Printf("Erro ao subscrever: %v", token.Error())
		}
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Erro ao conectar ao broker MQTT:", token.Error())
	}

	var err error
	csvFile, err = os.OpenFile("sensor_data.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Erro ao criar/abrir arquivo CSV:", err)
	}
	csvWriter = csv.NewWriter(csvFile)

	if info, _ := csvFile.Stat(); info.Size() == 0 {
		csvWriter.Write([]string{"timestamp", "temperatura", "umidade"})
		csvWriter.Flush()
	}

	defer csvFile.Close()

	// Rotas HTTP 
	http.HandleFunc("/ws", handleConnections)
	http.Handle("/", http.FileServer(http.Dir("./static/")))

	go handleMessages()

	log.Println("frontend: http://localhost:8080")
	log.Println("escutando topico: sensor/client em broker.hivemq.com")
	log.Println("dados salvos em: sensor_data.csv")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, ws)
			mutex.Unlock()
			break
		}
	}
}

func handleMessages() {
	for msg := range broadcast {
		mutex.RLock()
		for client := range clients {
			if err := client.WriteJSON(msg); err != nil {
				client.Close()
				mutex.Lock()
				delete(clients, client)
				mutex.Unlock()
			}
		}
		mutex.RUnlock()
	}
}