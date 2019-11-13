package external

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

type BtgRpcAnswer struct {
	Chain                string  `json:"chain"`
	Blocks               int     `json:"blocks"`
	Headers              int     `json:"headers"`
	Bestblockhash        string  `json:"bestblockhash"`
	Difficulty           float64     `json:"difficulty"`
	Mediantime           int     `json:"mediantime"`
	Verificationprogress float64 `json:"verificationprogress"`
	Chainwork            string  `json:"chainwork"`
	Pruned               bool    `json:"pruned"`
	Softforks            []struct {
		ID      string `json:"id"`
		Version int    `json:"version"`
		Reject  struct {
			Status bool `json:"status"`
		} `json:"reject"`
	} `json:"softforks"`
	Bip9Softforks struct {
		Csv struct {
			Status    string `json:"status"`
			StartTime int    `json:"startTime"`
			Timeout   int    `json:"timeout"`
			Since     int    `json:"since"`
		} `json:"csv"`
		Segwit struct {
			Status    string `json:"status"`
			StartTime int    `json:"startTime"`
			Timeout   int    `json:"timeout"`
			Since     int    `json:"since"`
		} `json:"segwit"`
	} `json:"bip9_softforks"`
}


func GetHeightFromBtgRpc(url string) (int, string) {
  resp, err := http.Get(url + "/rest/chaininfo.json")
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromBtgRpc: Failed to run get request %s`, err)
  }
  
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromBtgRpc: Failed to read body, err: '%s'`, err)
  }
  
  var v BtgRpcAnswer
  err = json.Unmarshal(body, &v)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromBtgRpc: json.Unmarshal, err: '%s'`, err)
  }
  return v.Blocks, ""
  
}
