package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("Coinstats started up!")

	// A path to config json is expected as the first (and only) command line argument
	if len(os.Args) != 2 {
		log.Fatalf("Expecting the path to a JSON config file as the first and only command line argument")
	}

	configPath := os.Args[1]

	// Open a read-only file and check if it exists
	configFile, err := os.Open(configPath)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Specified file does not exist, please provide a valid path to JSON config file")
	}
	defer configFile.Close()

	configJsonRaw, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatalf("Error reading config JSON file: %s", err)
	}
	log.Println(string(configJsonRaw))
	config := Config{}
	err = json.Unmarshal(configJsonRaw, &config)
	if err != nil {
		log.Fatalf("Couldn't parse config JSON: %s", err)
	}

	log.Println(config)

	/**
	Alright, so I want to know how much money do I have
	There are several things that have to be checked:
	1. Waves Blockchain
		* wallet's on-chain balances
		* deposits to VIRES (smart-contract, but can interact with a service)
		* deposits to Puzzle.Swap (smart-contract, but can interact with a service)
	2. WE Blockchain
		* wallet's WEST balance
		* deposit to EAST service (smart-contract, but can interact with a service)

	And more to come...
	*/

	var (
		wavesNodeApi     = config.Waves.NodeUrl
		myWavesAddress   = config.Waves.Addresses[0]
		wavesUsdtAssetId = "34N9YcEETLWn93qYQ64EsP1x89tSruJU44RrEMSXXEPJ"
		// wavesWestAssetId = "4LHHvYGNKJUg5hj65aGD5vgScvCBmLpdRFtjokvCjSL8"
	)

	// Asking Waves' nodes for my USDT balance

	resp, err := http.Get(wavesNodeApi + "/assets/balance/" + myWavesAddress + "/" + wavesUsdtAssetId)
	if err != nil {
		log.Fatalf("Error requesting USDT balance from Waves: %s", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body from Waves: %s", err)
	}

	var wavesResponseParsed map[string]interface{}
	//could this explode also?
	json.Unmarshal([]byte(respBody), &wavesResponseParsed)

	log.Println("Whole data\n", wavesResponseParsed)
	log.Println("Just the balance:", wavesResponseParsed["balance"])

	wavesUsdtBalanceRaw, isFound := wavesResponseParsed["balance"]
	if !isFound {
		log.Fatalf(`Couldn't find "balance" field in response from Waves: %s`, wavesResponseParsed)
	}

	wavesUsdtBalanceFloat, isOk := wavesUsdtBalanceRaw.(float64)
	if !isOk {
		log.Fatalf(`Couldn't assert values of field "balance" to an integer: %s`, wavesResponseParsed)
	}

	wavesUsdtDecimalsFactor := 1e-6
	wavesUsdtBalance := wavesUsdtBalanceFloat * wavesUsdtDecimalsFactor

	log.Printf("Finally, my Waves USDT balance: %f USDT", wavesUsdtBalance)

	/////////////////////////////////////////////////

	//
}
