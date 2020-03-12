package external

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "bytes"
//   "strings"
//   "net/url"
//   "strconv"
  "encoding/json"
//   "github.com/gorilla/websocket"
//   "github.com/antchfx/htmlquery"
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
  	// handle err
    return 0, fmt.Sprintf(`error: http.NewRequest: '%s'`, err)
  }
//   req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:75.0) Gecko/20100101 Firefox/75.0")
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
/*
  
  doc, err := htmlquery.Parse(strings.NewReader(string(body)))
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtzscanCom: htmlquery.Parse err: '%s'`, err)
  }

  list, err := htmlquery.QueryAll(doc, `//span[@id="lastblock"]`)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtzscanCom: htmlquery.QueryAll err: '%s'`, err)
  }
  
  var str string
  for _, n := range list { str = htmlquery.InnerText(n) }
  
  LastEtherBlock, err := strconv.Atoi(str)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtzscanCom: strconv.Atoi err: '%s'`, err)
  }
  
*/
  
  var v ETZAnswer
  err = json.Unmarshal(body, &v)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromParityRpc: json.Unmarshal, err: '%s'`, err)
  }
/*
  mm := v.(map[string]interface{})
  vvv := fmt.Sprintf(`%s`, mm["result"])
  vvvv := hex2int(vvv)
  
*/
  return v.BlockHeight, ""
}



















