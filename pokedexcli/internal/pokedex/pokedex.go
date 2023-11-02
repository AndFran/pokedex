package pokedex

import (
	"pokedexcli/internal/pokeapi"
)

type PokemonDb interface {
	Add(p pokeapi.PokemonDetailsResponse)
	Find(name string) (pokeapi.PokemonDetailsResponse, bool)
	GetAll() map[string]pokeapi.PokemonDetailsResponse
	Count() int
}

type PokemonMapDb struct {
	pokemons map[string]pokeapi.PokemonDetailsResponse
}

var Pokedex *PokemonMapDb

func init() {
	Pokedex = &PokemonMapDb{
		pokemons: make(map[string]pokeapi.PokemonDetailsResponse),
	}
}

func (pc *PokemonMapDb) Add(p pokeapi.PokemonDetailsResponse) {
	pc.pokemons[p.Name] = p
}

func (pc *PokemonMapDb) Find(name string) (pokeapi.PokemonDetailsResponse, bool) {
	pokemon, ok := pc.pokemons[name]
	return pokemon, ok
}

func (pc *PokemonMapDb) GetAll() map[string]pokeapi.PokemonDetailsResponse {
	return pc.pokemons
}

func (pc *PokemonMapDb) Count() int {
	return len(pc.pokemons)
}
