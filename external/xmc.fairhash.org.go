package external

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

type XmcFairHashAnswer struct {
	Config struct {
		Ports []struct {
			Port       int    `json:"port"`
			Difficulty int    `json:"difficulty"`
			Desc       string `json:"desc"`
			Type       string `json:"type,omitempty"`
		} `json:"ports"`
		HashrateWindow       int    `json:"hashrateWindow"`
		Fee                  int    `json:"fee"`
		Coin                 string `json:"coin"`
		CoinUnits            int64  `json:"coinUnits"`
		CoinDifficultyTarget int    `json:"coinDifficultyTarget"`
		Symbol               string `json:"symbol"`
		Depth                int    `json:"depth"`
		Donation             struct {
		} `json:"donation"`
		Version             string `json:"version"`
		MinPaymentThreshold int64  `json:"minPaymentThreshold"`
		MinExchangeLevel    int64  `json:"minExchangeLevel"`
		DenominationUnit    int64  `json:"denominationUnit"`
	} `json:"config"`
	System struct {
		Load        int `json:"load"`
		NumberCores int `json:"number_cores"`
	} `json:"system"`
	Pool struct {
		Stats struct {
			LastBlockFound string `json:"lastBlockFound"`
		} `json:"stats"`
		Blocks             []string        `json:"blocks"`
		TotalBlocks        int             `json:"totalBlocks"`
		Payments           []string        `json:"payments"`
		TotalPayments      int             `json:"totalPayments"`
		TotalMinersPaid    int             `json:"totalMinersPaid"`
		Miners             int             `json:"miners"`
		Workers            int             `json:"workers"`
		Hashrate           int             `json:"hashrate"`
		RoundHashes        int64           `json:"roundHashes"`
		CurrentPriceBTC    float64         `json:"currentPriceBTC"`
		CurrentPriceUSD    float64         `json:"currentPriceUSD"`
		CurrentPriceBTCUSD float64         `json:"currentPriceBTCUSD"`
		LuckArray          [][]interface{} `json:"luckArray"`
		LastBlockFound     string          `json:"lastBlockFound"`
	} `json:"pool"`
	Charts struct {
		Hashrate   [][]int         `json:"hashrate"`
		Workers    [][]int         `json:"workers"`
		Difficulty [][]interface{} `json:"difficulty"`
		Price      [][]float64     `json:"price"`
		Profit     [][]float64     `json:"profit"`
	} `json:"charts"`
	Network struct {
		Difficulty int64  `json:"difficulty"`
		Height     int    `json:"height"`
		Timestamp  int    `json:"timestamp"`
		Reward     int64  `json:"reward"`
		Hash       string `json:"hash"`
	} `json:"network"`
}


func GetHeightFromXmcFairHashAnswer() (int, string) {
//   resp, err := http.Get(url + suffix)
  resp, err := http.Get("https://xmc.fairhash.org/api/stats")
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromXmcFairHashAnswer: Failed to run get request %s`, err)
  }
  
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromXmcFairHashAnswer: Failed to read body, err: '%s'`, err)
  }
  
  var v XmcFairHashAnswer
  err = json.Unmarshal(body, &v)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromXmcFairHashAnswer: json.Unmarshal, err: '%s'`, err)
  }
  return v.Network.Height, ""
  
}
