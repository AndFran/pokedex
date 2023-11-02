package main

import (
	"pokedexcli_mine/internal/config"
	"pokedexcli_mine/internal/pokeapi"
	"pokedexcli_mine/internal/pokecache"
	"pokedexcli_mine/internal/pokedex"
)

func main() {

	app := config.ApplicationData{
		Cache: pokecache.NewMemoryCache(config.CacheDuration),
	}
	app.ApiClient = pokeapi.NewPokeClient(app.Cache)
	app.Db = pokedex.Pokedex

	repl(&app)
}
