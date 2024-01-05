package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

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

func GetLocations(locationOffset int) ([]LocationArea, error) {
	res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=%d&limit=20", locationOffset))
	locations := []LocationArea{}
	if err != nil {
		return locations, errors.New(fmt.Sprintf("Request failed with error: %s", err))
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return locations, errors.New(fmt.Sprintf("Response failed with status code: %d and body %s", res.StatusCode, body))
	}
	if err != nil {
		return locations, errors.New(fmt.Sprintf("Response failed with error %s", err))
	}

	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return locations, err
	}

	return response.Results, nil
}
