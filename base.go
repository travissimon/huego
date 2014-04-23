package huego

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var httpClient = &http.Client{}

type Base struct {
	Id         string `json:"id"`
	InternalIp string `json:"internalipaddress"`
	Username   string
	lights     []*Light
}

type LightName struct {
	Id   string
	Name string
}

func NewLightName(id, name string) LightName {
	return LightName{id, name}
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
		var bases []Base
		err = json.Unmarshal(contents, &bases)
		if err != nil {
			return nil, err
		}
		for _, b := range bases {
			b.Username = "newdeveloper"
		}
		return bases, nil
	}
}

func (b *Base) GetApi(methodPart string) string {
	url := fmt.Sprintf("http://%v/api/%v%v", b.InternalIp, "newdeveloper", methodPart)
	return url
}

func (b *Base) doGet(urlPart string) ([]byte, error) {
	url := b.GetApi(urlPart)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return contents, nil
	}
}

func (b *Base) doPost(urlPart string, body string) ([]byte, error) {
	url := b.GetApi(urlPart)
	response, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return contents, nil
	}
}

func (b *Base) doPut(urlPart string, body string) ([]byte, error) {
	url := b.GetApi(urlPart)
	req, err := http.NewRequest("PUT", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func (b *Base) GetLightNames() ([]LightName, error) {
	response, err := http.Get(b.GetApi("/lights"))
	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var object interface{}
		err = json.Unmarshal(contents, &object)
		if err != nil {
			return nil, err
		}
		objectMap := object.(map[string]interface{})
		names := make([]LightName, 0, 3)
		for k, v := range objectMap {
			nameObj := v.(map[string]interface{})
			name := NewLightName(k, nameObj["name"].(string))
			names = append(names, name)
		}
		return names, nil
	}
}

func (b *Base) GetLights() ([]*Light, error) {
	if b.lights != nil && len(b.lights) > 0 {
		return b.lights, nil
	}

	names, err := b.GetLightNames()
	if err != nil {
		return nil, err
	}
	lights := make([]*Light, 0, len(names))
	for _, name := range names {
		light, err := b.GetLightById(name.Id)
		if err != nil {
			return nil, err
		}
		lights = append(lights, light)
	}

	return lights, nil
}

func (b *Base) GetLightById(id string) (*Light, error) {
	url := fmt.Sprintf("/lights/%v", id)
	response, err := http.Get(b.GetApi(url))
	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var light = new(Light)
		err = json.Unmarshal(contents, light)
		if err != nil {
			return nil, err
		}
		light.Id = id
		light.base = b
		light.prevState = new(State)
		light.ResetState()

		return light, nil
	}
}
