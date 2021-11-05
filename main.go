package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"strconv"

	"github.com/bitvavo/go-bitvavo-api"
	"gopkg.in/yaml.v2"
)

type APICfg struct {
	Key    string
	Secret string
}

type Buyin struct {
	Ticker string
	Num    float64
	Price  float64
}

type Cfg struct {
	API    APICfg
	Buyins []Buyin
}

func main() {
	cfgPath, err := getConfigPath()
	if err != nil {
		log.Fatal("Could not determine current user. Like. What.", err)
	}

	dat, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Fatal("Couldn't read configuration file", err)
	}

	cfg := Cfg{}
	err = yaml.Unmarshal(dat, &cfg)
	if err != nil {
		log.Fatal("rip", err)
	}

	// Set up api
	bitv := bitvavo.Bitvavo{
		ApiKey:       cfg.API.Key,
		ApiSecret:    cfg.API.Secret,
		RestUrl:      "https://api.bitvavo.com/v2",
		WsUrl:        "wss://ws.bitvavo.com/v2/",
		AccessWindow: 10000,
		Debugging:    false,
	}

	totalInvested := 0.0
	totalNow := 0.0

	for _, buyin := range cfg.Buyins {
		// Pull rates
		options := make(map[string]string)
		options["market"] = fmt.Sprintf("%s-EUR", buyin.Ticker)
		response, err := bitv.TickerPrice(options)
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

		fmt.Printf("%6s: %.2f (%.2f): %.5f%%\n", buyin.Ticker, current, buyinCost, (current/buyinCost)*100.0)
	}

	fmt.Printf("Totals: %.2f (%.2f): %.5f%%\n", totalNow, totalInvested, (totalNow/totalInvested)*100.0)
}

func getConfigPath() (string, error) {
	// Read $HOME/.amirich.yaml
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	dir := usr.HomeDir
	return path.Join(dir, ".amirich.yaml"), nil
}
