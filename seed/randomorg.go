package seed

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// RandomOrg returns a Seeder that gets values from random.org using the api key
func RandomOrg(apiKey string) Seeder {
	type Params struct {
		APIKey      string `json:"apiKey"`
		N           int    `json:"n"`
		Length      int    `json:"length"`
		Min         int    `json:"min"`
		Max         int    `json:"max"`
		Replacement bool   `json:"replacement"`
		Base        int    `json:"base"`
	}

	type Random struct {
		Data [][]byte `json:"data"`
	}

	type Result struct {
		Random Random `json:"random"`
	}

	type API struct {
		JSONRPC string  `json:"jsonrpc"`
		Method  string  `json:"method,omitempty"`
		ID      int     `json:"id"`
		Params  *Params `json:"params,omitempty"`
		Result  *Result `json:"result,omitempty"`
	}

	id := new(int)

	return func() (uint64, error) {
		(*id)++

		request := &API{
			JSONRPC: "2.0",
			Method:  "generateIntegerSequences",
			ID:      *id,
			Params: &Params{
				APIKey:      apiKey,
				N:           1,
				Length:      8,
				Min:         0,
				Max:         255,
				Replacement: false,
				Base:        10,
			},
		}

		requestJSON, _ := json.Marshal(request)
		httpReq, _ := http.NewRequest("POST", "https://api.random.org/json-rpc/2/invoke", bytes.NewReader(requestJSON))
		httpReq.Header.Add("Content-Type", "application/json")
		httpReq.Header.Add("Accept", "application/json")

		httpResp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			return 0, err
		}
		defer httpResp.Body.Close()

		responseJSON, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return 0, err
		}

		var response API
		err = json.Unmarshal(responseJSON, &response)
		if err != nil {
			return 0, err
		}

		if response.Result == nil {
			return 0, fmt.Errorf("no result")
		}

		return binary.LittleEndian.Uint64(response.Result.Random.Data[0]), nil
	}
}
