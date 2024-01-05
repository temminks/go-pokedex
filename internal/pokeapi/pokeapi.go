package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/temminks/go-pokedex/internal/pokecache"
)

var cache = pokecache.NewCache(5 * time.Second)

type LocationArea struct {
	Name string
	Url  string
}

type Response struct {
	Count    int
	Next     string
	Previous string
	Results  []LocationArea
}

func getData(url string) (body []byte, err error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", err))
	}
	body, err = io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("Response failed with status code: %d and body %s", res.StatusCode, body))
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Response failed with error %s", err))
	}

	return body, nil
}

func GetLocations(locationOffset int) ([]LocationArea, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=%d&limit=20", locationOffset)
	body, exists := cache.Get(url)
	if !exists {
		responseBody, err := getData(url)
		if err != nil {
			return nil, err
		}
		cache.Add(url, responseBody)
		body = responseBody
	}

	response := Response{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Results, nil
}
