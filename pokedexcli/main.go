package main

import (
	"pokedexcli/internal/config"
	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokecache"
	"pokedexcli/internal/pokedex"
)

func main() {

	app := config.ApplicationData{
		Cache: pokecache.NewMemoryCache(config.CacheDuration),
	}
	app.ApiClient = pokeapi.NewPokeClient(app.Cache)
	app.Db = pokedex.Pokedex

	repl(&app)
}
