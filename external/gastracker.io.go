// package main
package external

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
)


type EtcGastrackerMessage struct {
	Offset int `json:"offset"`
	Items []struct {
		Gas          int    `json:"gas"`
		Hash         string `json:"hash"`
		Height       int    `json:"height"`
		Miner        string `json:"miner"`
		MiningTime   int    `json:"miningTime"`
		Timestamp    string `json:"timestamp"`
		Transactions int    `json:"transactions"`
		Value        struct {
			Ether float64 `json:"ether"`
			Hex   string  `json:"hex"`
			Wei   string  `json:"wei"`
		} `json:"value"`
	} `json:"items"`
}


func GetEtcHeightFromGastrackerIo () (int, string) {
  url := "https://api.gastracker.io/v1/blocks/latest"
  
  resp, err := http.Get(url)
  if err != nil {
    return 0, fmt.Sprintf(`error: getEtcHeightFromGastrackerIo: Failed to get url: '%s', err: '%s'`, url, err)
  }
  defer resp.Body.Close()
  
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return 0, fmt.Sprintf(`error: getEtcHeightFromGastrackerIo: Failed to read body, err: '%s'`, err)
  }
  
  var msg EtcGastrackerMessage
  err = json.Unmarshal(body, &msg)
  if err != nil {
    return 0, fmt.Sprintf(`error: getEtcHeightFromGastrackerIo: Failed to json.Unmarshal, err: '%s'`, err)
  }
  
  if (len(msg.Items) == 0) {
    return 0, fmt.Sprintf(`error: getEtcHeightFromGastrackerIo: len(Items) == 0`)
  }
  
  return msg.Items[0].Height, ""
}
