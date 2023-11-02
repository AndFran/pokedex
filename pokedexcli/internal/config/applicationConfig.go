package config

import (
	"pokedexcli_mine/internal/pokeapi"
	"pokedexcli_mine/internal/pokecache"
	"pokedexcli_mine/internal/pokedex"
	"time"
)

const CacheDuration = 20 * time.Second

type ApplicationData struct {
	Cache     pokecache.PokeCache
	ApiClient *pokeapi.PokeClient
	Db        pokedex.PokemonDb
}
