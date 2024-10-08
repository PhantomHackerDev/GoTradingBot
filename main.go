package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// CMCResponse represents the response structure from CoinMarketCap API
type CMCResponse struct {
	Data []struct {
		Name  string `json:"name"`
		Quote map[string]struct {
			Price float64 `json:"price"`
		} `json:"quote"`
	} `json:"data"`
}

type TokenPrice struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Response struct {
	Tokens []TokenPrice `json:"tokens"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cmcAPIKey := os.Getenv("COINMARKETCAP_API_KEY")

	for {
		prices, err := getPrices(cmcAPIKey)
		if err != nil {
			log.Fatal(err)
		}

		response := Response{Tokens: prices}
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(jsonData))

		// Wait for 60 seconds before fetching data again
		time.Sleep(60 * time.Second)
	}
}

// Fetch the current prices of multiple cryptocurrencies from CoinMarketCap
func getPrices(apiKey string) ([]TokenPrice, error) {
	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-CMC_PRO_API_KEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching prices: %s", resp.Status)
	}

	var cmcResponse CMCResponse
	if err := json.NewDecoder(resp.Body).Decode(&cmcResponse); err != nil {
		return nil, err
	}

	var prices []TokenPrice
	for _, token := range cmcResponse.Data {
		prices = append(prices, TokenPrice{
			Name:  token.Name,
			Price: token.Quote["USD"].Price,
		})
	}

	return prices, nil
}
