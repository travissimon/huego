package huego

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Base struct {
	Id         string `json:"id"`
	InternalIp string `json:"internalipaddress"`
}

func DiscoverBases() ([]Base, error) {
	response, err := http.Get("http://www.meethue.com/api/nupnp")

	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		var b []Base
		err = json.Unmarshal(contents, &b)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
}
