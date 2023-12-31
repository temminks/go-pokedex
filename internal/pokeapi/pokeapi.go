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

type Location struct {
	Id                   int
	Name                 string
	GameIndex            int                      `json:"game_index"`
	EncounterMethodRates []map[string]interface{} `json:"encounter_method_rates"`
	Location             map[string]string
	Names                []map[string]interface{}
	PokemonEncounters    []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon        Pokemon
	VersionDetails []map[string]interface{} `json:"version_details"`
}

// BaseExperience for blissey: 635, mewtwo: 340, rayquaza: 340
// dialga: 340, palkia 340, giratina 340, arceus 360, audino 390
// eternatus: 345
type PokemonDetails struct {
	Id             int
	Name           string
	BaseExperience int `json:"base_experience"`
	Height         int
	IsDefault      bool `json:"is_default"`
	Order          int
	Weight         int
	Species        Pokemon
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int
		Stat     Pokemon
	}
	Types []struct {
		Slot int
		Type Pokemon
	}
}

type Pokemon struct {
	Name string
	Url  string
}

type Response struct {
	Count    int
	Next     string
	Previous string
	Results  []Pokemon
}

type RequestError struct {
	StatusCode int
	Body       []byte
	Err        error
}

func (r *RequestError) Error() string {
	return r.Err.Error()
}

func loadOrRetrieveFromCache(url string) (body []byte, err error) {
	body, exists := cache.Get(url)
	if !exists {
		responseBody, err := getData(url)
		if err != nil {
			return nil, err
		}
		cache.Add(url, responseBody)
		return responseBody, nil
	}
	return body, nil
}

func getData(url string) (body []byte, err error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", err))
	}
	body, err = io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, &RequestError{
			StatusCode: res.StatusCode,
			Body:       body,
			Err:        fmt.Errorf("Response failed with status code: `%d` and body `%s`", res.StatusCode, body),
		}
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Response failed with error %s", err))
	}

	return body, nil
}

func GetLocation(location string) (Location, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location)
	body, err := loadOrRetrieveFromCache(url)
	response := Location{}
	if err != nil {
		if re, ok := err.(*RequestError); ok {
			if re.StatusCode == 404 {
				return response, fmt.Errorf("Location `%s` not found.", location)
			} else {
				return response, re.Err
			}
		}
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func GetLocations(locationOffset int) ([]Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=%d&limit=20", locationOffset)
	body, err := loadOrRetrieveFromCache(url)
	if err != nil {
		return []Pokemon{}, err
	}

	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Results, nil
}

func GetPokemonDetails(name string) (PokemonDetails, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	body, err := loadOrRetrieveFromCache(url)
	response := PokemonDetails{}
	if err != nil {
		if re, ok := err.(*RequestError); ok {
			if re.StatusCode == 404 {
				return response, fmt.Errorf("Pokemon `%s` not found.", name)
			} else {
				return response, re.Err
			}
		}
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
