package config

import (
	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokecache"
	"pokedexcli/internal/pokedex"
	"time"
)

const CacheDuration = 20 * time.Second

type ApplicationData struct {
	Cache     pokecache.PokeCache
	ApiClient *pokeapi.PokeClient
	Db        pokedex.PokemonDb
}
