package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bitvavo/go-bitvavo-api"
	"github.com/pubkraal/amirich/internal/pkg/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	api := getAPI(cfg)

	totalInvested := 0.0
	totalNow := 0.0

	for _, buyin := range cfg.Buyins {
		// Pull rates
		options := make(map[string]string)
		options["market"] = fmt.Sprintf("%s-EUR", buyin.Ticker)
		response, err := api.TickerPrice(options)
		if err != nil {
			log.Fatal("Could not retrieve market:", err)
		}
		if len(response) != 1 {
			log.Fatal("Response has the wrong number of prices")
		}
		respPrice := response[0].Price
		price, err := strconv.ParseFloat(respPrice, 64)
		if err != nil {
			log.Fatal("Invalid price for market", response, err)
		}

		buyinCost := buyin.Num * buyin.Price
		current := buyin.Num * price

		totalInvested = totalInvested + buyinCost
		totalNow = totalNow + current

		fmt.Printf("%6s: %.2f (%.7f): %.5f%%\n", buyin.Ticker, current, price, ((current/buyinCost)-1)*100.0)
	}

	fmt.Printf("Totals: %.2f (%.2f): %.5f%%\n", totalNow, totalInvested, ((totalNow/totalInvested)-1)*100.0)
}

func getAPI(cfg config.Cfg) *bitvavo.Bitvavo {
	return &bitvavo.Bitvavo{
		ApiKey:       cfg.API.Key,
		ApiSecret:    cfg.API.Secret,
		RestUrl:      "https://api.bitvavo.com/v2",
		WsUrl:        "wss://ws.bitvavo.com/v2/",
		AccessWindow: 10000,
		Debugging:    false,
	}
}
