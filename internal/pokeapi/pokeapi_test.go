package pokeapi

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
)

func TestGetPokemon(t *testing.T) {
	expectedPokemonDetails := PokemonDetails{
		Id:             123,
		Name:           "scyther",
		BaseExperience: 100,
		Height:         15,
		IsDefault:      true,
		Order:          196,
		Weight:         560,
		Species: map[string]interface{}{
			"name": "scyther",
			"url":  "https://pokeapi.co/api/v2/pokemon-species/123/",
		},
	}

	actualPokemonDetails, err := GetPokemonDetails("scyther")
	if err != nil {
		t.Errorf("did not expect an error when calling the API: %s", err)
	}

	if !reflect.DeepEqual(expectedPokemonDetails, actualPokemonDetails) {
		t.Errorf("expected PokemonDetails: %#v, actual PokemonDetails: %#v", expectedPokemonDetails, actualPokemonDetails)
	}
}

func TestGetLocationValidLocation(t *testing.T) {
	expectedLocation := Location{
		Id:        126,
		Name:      "old-chateau-entrance",
		GameIndex: 126,
		PokemonEncounters: []PokemonEncounter{
			{
				Pokemon: Pokemon{
					Name: "gastly",
				},
			},
			{
				Pokemon: Pokemon{
					Name: "haunter",
				},
			},
		},
	}

	location, err := GetLocation("old-chateau-entrance")
	if err != nil {
		t.Errorf("did not expect an error when calling the API: %s", err)
	}

	if expectedLocation.Id != location.Id {
		t.Errorf("expected Id: %d, actual Id: %d", expectedLocation.Id, location.Id)
	}

	actualNames := getPokemonNames(location.PokemonEncounters)
	expectedNames := getPokemonNames(expectedLocation.PokemonEncounters)
	if slices.Compare(actualNames, expectedNames) != 0 {
		t.Errorf("expected Pokemon: %#v, actual Pokemon %#v", expectedNames, actualNames)
	}
}

// get the names of all Pokemon in a `[]PokemonEncounter`
func getPokemonNames(pe []PokemonEncounter) []string {
	var pokemonNames []string
	for _, encounter := range pe {
		pokemonNames = append(pokemonNames, encounter.Pokemon.Name)
	}
	return pokemonNames
}

func TestGetPokemonDetailsInvalidName(t *testing.T) {
	invalidName := "this-pokemon-does-not-exist"
	_, err := GetPokemonDetails(invalidName)
	if err == nil {
		t.Errorf("An invalid name `%s` should return an error: %s", invalidName, err)
	}
	expectedErr := fmt.Sprintf("Pokemon `%s` not found.", invalidName)
	if err.Error() != expectedErr {
		t.Errorf("Unexpected error message! Expected `%s`, actual: `%s`.", expectedErr, err)
	}
}

func TestGetPokemonDetailsValidName(t *testing.T) {
	validName := "scyther"
	_, err := GetPokemonDetails(validName)
	if err != nil {
		t.Errorf("A valid Pokemon name (`%s`) should not result in an error: %s", validName, err)
	}
}
