package external

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

type ExplorerZchaInAnswer struct {
	Name            string  `json:"name"`
	Accounts        int     `json:"accounts"`
	Transactions    int     `json:"transactions"`
	BlockHash       string  `json:"blockHash"`
	BlockNumber     int     `json:"blockNumber"`
	Difficulty      float64 `json:"difficulty"`
	Hashrate        int64   `json:"hashrate"`
	MeanBlockTime   float64 `json:"meanBlockTime"`
	PeerCount       int     `json:"peerCount"`
	ProtocolVersion int     `json:"protocolVersion"`
	RelayFee        float64 `json:"relayFee"`
	Version         int     `json:"version"`
	SubVersion      string  `json:"subVersion"`
	TotalAmount     float64 `json:"totalAmount"`
	SproutPool      float64 `json:"sproutPool"`
	SaplingPool     float64 `json:"saplingPool"`
}


func GetHeightFromExplorerZchaInAnswer() (int, string) {
  req, err := http.NewRequest("GET", "https://api.zcha.in/v2/mainnet/network", nil)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromExplorerZchaInAnswer: Failed to prepare request %s`, err)
  }
  req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:72.0) Gecko/20100101 Firefox/72.0")
  req.Header.Set("Accept", "*/*")
  req.Header.Set("Accept-Language", "ru,en-US;q=0.7,en;q=0.3")
  req.Header.Set("Referer", "https://explorer.zcha.in/")
  req.Header.Set("Origin", "https://explorer.zcha.in")
  req.Header.Set("Dnt", "1")
  req.Header.Set("Connection", "keep-alive")
  req.Header.Set("Pragma", "no-cache")
  req.Header.Set("Cache-Control", "no-cache")
  req.Header.Set("Te", "Trailers")
  
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromExplorerZchaInAnswer: Failed to run get request %s`, err)
  }
  defer resp.Body.Close()
  
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromExplorerZchaInAnswer: Failed to read body, err: '%s'`, err)
  }
  
  var v ExplorerZchaInAnswer
  err = json.Unmarshal(body, &v)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromExplorerZchaInAnswer: json.Unmarshal, err: '%s'`, err)
  }
  return v.BlockNumber, ""
  
}
