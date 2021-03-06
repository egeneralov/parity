package external

import (
  "fmt"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

type zecRpcAnswer struct {
	Result struct {
		Chain                string  `json:"chain"`
		Blocks               int     `json:"blocks"`
		Headers              int     `json:"headers"`
		Bestblockhash        string  `json:"bestblockhash"`
		Difficulty           float64 `json:"difficulty"`
		Verificationprogress float64 `json:"verificationprogress"`
		Chainwork            string  `json:"chainwork"`
		Pruned               bool    `json:"pruned"`
		SizeOnDisk           int64   `json:"size_on_disk"`
		Commitments          int     `json:"commitments"`
		ValuePools           []struct {
			ID            string  `json:"id"`
			Monitored     bool    `json:"monitored"`
			ChainValue    float64 `json:"chainValue"`
			ChainValueZat int64   `json:"chainValueZat"`
		} `json:"valuePools"`
		Softforks []struct {
			ID      string `json:"id"`
			Version int    `json:"version"`
			Enforce struct {
				Status   bool `json:"status"`
				Found    int  `json:"found"`
				Required int  `json:"required"`
				Window   int  `json:"window"`
			} `json:"enforce"`
			Reject struct {
				Status   bool `json:"status"`
				Found    int  `json:"found"`
				Required int  `json:"required"`
				Window   int  `json:"window"`
			} `json:"reject"`
		} `json:"softforks"`
		Upgrades struct {
			FiveBa81B19 struct {
				Name             string `json:"name"`
				Activationheight int    `json:"activationheight"`
				Status           string `json:"status"`
				Info             string `json:"info"`
			} `json:"5ba81b19"`
			Seven6B809Bb struct {
				Name             string `json:"name"`
				Activationheight int    `json:"activationheight"`
				Status           string `json:"status"`
				Info             string `json:"info"`
			} `json:"76b809bb"`
			TwoBb40E60 struct {
				Name             string `json:"name"`
				Activationheight int    `json:"activationheight"`
				Status           string `json:"status"`
				Info             string `json:"info"`
			} `json:"2bb40e60"`
		} `json:"upgrades"`
		Consensus struct {
			Chaintip  string `json:"chaintip"`
			Nextblock string `json:"nextblock"`
		} `json:"consensus"`
	} `json:"result"`
	Error interface{} `json:"error"`
	ID    string      `json:"id"`
}

func GetHeightFromZecRpc(url string, rpcUser string, rpcPassword string) (int, string) {
  // Generated by curl-to-Go: https://mholt.github.io/curl-to-go
  
  payload := strings.NewReader("{\"jsonrpc\": \"1.0\", \"id\":\"curltest\", \"method\": \"getblockchaininfo\", \"params\": [] }")
  req, err := http.NewRequest("POST", url, payload)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromZecRpc: Failed to prepare request %s`, err)
  }
  req.SetBasicAuth(rpcUser, rpcPassword)
  req.Header.Set("Content-Type", "text/plain;")
  
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromZecRpc: Failed to run get request %s`, err)
  }
  defer resp.Body.Close()
  
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromZecRpc: Failed to read body, err: '%s'`, err)
  }
  
  var v zecRpcAnswer
  err = json.Unmarshal(body, &v)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromZecRpc: json.Unmarshal, err: '%s'`, err)
  }
  return v.Result.Blocks, ""  
}
