package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type TrafficSignal struct {
	ID              int `json:"id"`
	Congestion      int `json:"congestion"`
	RedLightTime    int `json:"redLightTime"`
	YellowLightTime int `json:"yellowLightTime"`
	GreenLightTime  int `json:"greenLightTime"`
	AverageFlowRate int `json:"averageFlowRate"`
}

type StatsResponse struct {
	TotalRequests     int64 `json:"totalRequests"`
	RequestsPerSecond int64 `json:"requestsPerSecond"`
}

func main() {
	// Gerar 100 requisições por segundo durante 10 segundos
	//ticker := time.Tick(time.Millisecond * 10)
	//stop := time.NewTimer(time.Second * 100)
	//
	//for {
	//	select {
	//	case <-ticker:
	//		go func() {
	//			trafficSignal := generateRandomTrafficSignal()
	//			err := sendTrafficSignalData(trafficSignal)
	//			if err != nil {
	//				log.Printf("Failed to send traffic signal data: %v\n", err)
	//			}
	//		}()
	//	case <-stop.C:
	//		return
	//	}
	//}
	// Gerar 10 requisições POST para cada sinal de trânsito
	for {
		for signalID := 1; signalID <= 3; signalID++ {
			for i := 1; i <= 10; i++ {
				trafficSignal := generateRandomTrafficSignal(signalID)
				err := sendTrafficSignalData(trafficSignal)
				if err != nil {
					log.Printf("Failed to send traffic signal data: %v\n", err)
				}
				time.Sleep(time.Millisecond * 100)
			}
		}

		// Fazer requisições GET para obter a média de engarrafamento de cada sinal de trânsito
		for signalID := 1; signalID <= 3; signalID++ {
			averageFlowRate, err := getAverageFlowRate(signalID)
			if err != nil {
				log.Printf("Failed to get average flow rate for traffic signal %d: %v\n", signalID, err)
			} else {
				log.Printf("Average flow rate for Traffic Signal %d: %d\n", signalID, averageFlowRate)
			}
		}
	}
}

func generateRandomTrafficSignal(signalID int) TrafficSignal {
	rand.Seed(time.Now().UnixNano())

	return TrafficSignal{
		ID:              signalID,
		Congestion:      rand.Intn(100),
		RedLightTime:    rand.Intn(10),
		YellowLightTime: rand.Intn(10),
		GreenLightTime:  rand.Intn(10),
	}
}

func sendTrafficSignalData(trafficSignal TrafficSignal) error {
	payload, err := json.Marshal(trafficSignal)
	if err != nil {
		return err
	}

	//resp, err := http.Post("http://processor-svc/traffic", "application/json", bytes.NewBuffer(payload))
	resp, err := http.Post("http://localhost:8080/traffic", "application/json", bytes.NewBuffer(payload))

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func getAverageFlowRate(signalID int) (int, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/traffic/info?id=%d", signalID))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response struct {
		AverageFlowRate int `json:"averageFlowRate"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	return response.AverageFlowRate, nil
}
