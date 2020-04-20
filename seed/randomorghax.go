package seed

import (
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strings"
)

// RandomOrgHax returns a Seeder that gets values from random.org using an unofficial method
func RandomOrgHax() Seeder {
	return func() (uint64, error) {
		req, _ := http.NewRequest("GET", "https://www.random.org/cgi-bin/randbyte?nbytes=8&format=h", nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()

		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}

		respHex := strings.ReplaceAll(string(respBytes), " ", "")
		respVal, err := hex.DecodeString(respHex)
		if err != nil {
			return 0, err
		}

		return binary.LittleEndian.Uint64(respVal), nil
	}
}
