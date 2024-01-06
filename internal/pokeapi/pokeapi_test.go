package pokeapi

import "testing"

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
}
