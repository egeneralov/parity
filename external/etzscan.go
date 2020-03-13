package external

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "bytes"
  "encoding/json"
)


type ETZPayload struct {
	Action string `json:"action"`
}

type ETZAnswer struct {
	Blocks []struct {
		ID        string        `json:"_id"`
		ExtraData string        `json:"extraData"`
		Miner     string        `json:"miner"`
		Number    int           `json:"number"`
		Timestamp int           `json:"timestamp"`
		Txs       []interface{} `json:"txs"`
	} `json:"blocks"`
	BlockHeight    int     `json:"blockHeight"`
	BlockTime      int     `json:"blockTime"`
	TPS            float64     `json:"TPS"`
	MeanDayRewards float64 `json:"meanDayRewards"`
}

func GetEthHeightFromEtzscanCom () (int, string) {
  data := ETZPayload{
    Action: "latest_blocks",
  }
  payloadBytes, err := json.Marshal(data)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtzscanCom: json.Marshal: '%s'`, err)
  }
  payloadReader := bytes.NewReader(payloadBytes)
  
  req, err := http.NewRequest("POST", "https://etzscan.com/data", payloadReader)
  if err != nil {
    return 0, fmt.Sprintf(`error: http.NewRequest: '%s'`, err)
  }
  req.Header.Set("Accept", "application/json, text/plain, */*")
  req.Header.Set("Accept-Language", "ru,en-US;q=0.7,en;q=0.3")
  req.Header.Set("Content-Type", "application/json;charset=utf-8")
  req.Header.Set("Origin", "https://etzscan.com")
  req.Header.Set("Dnt", "1")
  req.Header.Set("Connection", "keep-alive")
  req.Header.Set("Referer", "https://etzscan.com/home")
  req.Header.Set("Pragma", "no-cache")
  req.Header.Set("Cache-Control", "no-cache")
  req.Header.Set("Te", "Trailers")
  
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return 0, fmt.Sprintf(`error: http.DefaultClient.Do(req): '%s'`, err)
  }
  defer resp.Body.Close()
  
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return 0, fmt.Sprintf(`error: ioutil.ReadAll(resp.Body): '%s'`, err)
  }
  
  var v ETZAnswer
  err = json.Unmarshal(body, &v)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromParityRpc: json.Unmarshal, err: '%s'`, err)
  }
  return v.BlockHeight, ""
}
