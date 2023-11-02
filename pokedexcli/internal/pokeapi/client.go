package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
)

type PokeClient struct {
	Cache pokecache.PokeCache
}

type cachable interface {
	LocationAreasResponse | LocationAreaResponse | PokemonDetailsResponse
}

func NewPokeClient(cache pokecache.PokeCache) *PokeClient {
	pc := &PokeClient{
		Cache: cache,
	}
	return pc
}

const BaseUrl = "https://pokeapi.co/api/v2"
const AreaUrl = BaseUrl + "/location-area"
const PokemonUrl = "/pokemon"

func marshallResponse[T cachable](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func performAPICall(url string, data *[]byte) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	*data, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("bad request status %s", err)
	}
	return nil
}

func (p *PokeClient) GetAreaLocationAreasByPageUrl(pageUrl string) (LocationAreasResponse, error) {

	data, ok := p.Cache.Get(pageUrl)
	if !ok {
		err := performAPICall(pageUrl, &data)
		if err != nil {
			return LocationAreasResponse{}, err
		}
		p.Cache.Add(pageUrl, data)
	}

	return marshallResponse[LocationAreasResponse](data)
}

func (p *PokeClient) GetAreaLocationsAreas() (LocationAreasResponse, error) {
	return p.GetAreaLocationAreasByPageUrl(AreaUrl)
}

func (p *PokeClient) GetLocationArea(location string) (LocationAreaResponse, error) {
	data, ok := p.Cache.Get(location)
	if !ok {
		err := performAPICall(fmt.Sprintf("%s/%s", AreaUrl, location), &data)
		if err != nil {
			return LocationAreaResponse{}, err
		}
	}
	p.Cache.Add(location, data)
	return marshallResponse[LocationAreaResponse](data)
}

func (p *PokeClient) GetPokemon(pokemonName string) (PokemonDetailsResponse, error) {
	data, ok := p.Cache.Get(pokemonName)
	if !ok {
		err := performAPICall(fmt.Sprintf("%s/%s/%s", BaseUrl, PokemonUrl, pokemonName), &data)
		if err != nil {
			return PokemonDetailsResponse{}, err
		}
	}
	p.Cache.Add(pokemonName, data)
	return marshallResponse[PokemonDetailsResponse](data)
}
