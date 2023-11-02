package commands

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"pokedexcli/internal/config"
	"pokedexcli/internal/pokeapi"
)

type MenuCommand struct {
	Name            string
	Description     string
	CommandCallback func(subcommand []string) error // need to get state to the callbacks
}

type MenuState struct {
	NextUrl *string
	PrevUrl *string
	Explore string
}

type Menu struct {
	State        MenuState
	MenuCommands map[string]MenuCommand
	Config       *config.ApplicationData
}

func NewMenu(app *config.ApplicationData) Menu {
	menu := Menu{
		MenuCommands: make(map[string]MenuCommand),
		Config:       app,
	}
	menu.setMenuCommands()
	return menu
}

func (m *Menu) helpCb(fields []string) error {
	for k, v := range m.MenuCommands {
		fmt.Printf("%s - %s\n", k, v.Description)
	}
	return nil
}

func (m *Menu) exitCb(fields []string) error {
	fmt.Println("Exiting system")
	os.Exit(0)
	return nil
}

func (m *Menu) mapCb(fields []string) error {
	client := m.Config.ApiClient
	var locations pokeapi.LocationAreasResponse
	var err error

	if m.State.NextUrl != nil {
		locations, err = client.GetAreaLocationAreasByPageUrl(*m.State.NextUrl)
	} else {
		locations, err = client.GetAreaLocationsAreas()
	}

	if err != nil {
		return err
	}
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	m.State.NextUrl = locations.Next
	m.State.PrevUrl = locations.Previous
	return nil
}

func (m *Menu) mapbCb(fields []string) error {
	client := m.Config.ApiClient
	var locations pokeapi.LocationAreasResponse
	var err error

	if m.State.PrevUrl != nil {
		locations, err = client.GetAreaLocationAreasByPageUrl(*m.State.PrevUrl)
	} else {
		fmt.Println("There are currently no previous results")
	}

	if err != nil {
		return err
	}
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	m.State.NextUrl = locations.Next
	m.State.PrevUrl = locations.Previous
	return nil
}

func (m *Menu) addMenuCommand(command MenuCommand) {
	m.MenuCommands[command.Name] = command
}

func (m *Menu) stateCb(fields []string) error {
	fmt.Println(m.State)
	return nil
}

func (m *Menu) cacheKeyCb(fields []string) error {
	for i, v := range m.Config.Cache.ViewKeys() {
		fmt.Println(i, "-", v)
	}
	return nil
}

func (m *Menu) exploreCb(fields []string) error {
	if len(fields) != 2 {
		return errors.New("usage: explore area-name")
	}
	fmt.Printf("Exploring %s...", fields[1])
	client := m.Config.ApiClient
	locationDetails, err := client.GetLocationArea(fields[1])
	if err != nil {
		return err
	}
	if len(locationDetails.PokemonEncounters) > 0 {
		fmt.Println("Found Pokemon:")
		for _, v := range locationDetails.PokemonEncounters {
			fmt.Println(v.Pokemon.Name)
		}
	} else {
		fmt.Println("No pokemon found")
	}

	return nil
}

func (m *Menu) catchCb(fields []string) error {
	if len(fields) != 2 {
		return errors.New("usage: catch pokemon-name")
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", fields[1])
	client := m.Config.ApiClient
	pokemonDetails, err := client.GetPokemon(fields[1])
	if err != nil {
		return err
	}

	computerScore := (1000 - pokemonDetails.BaseExperience) / 10 //60
	playerChance := rand.Intn(100) + 1
	if playerChance <= computerScore {
		fmt.Printf("%s was caught!\n", fields[1])
		m.Config.Db.Add(pokemonDetails)
	} else {
		fmt.Printf("%s escaped!\n", fields[1])
	}

	return nil
}

func (m *Menu) pokedexCb(fields []string) error {
	if m.Config.Db.Count() == 0 {
		fmt.Println("You do not have any Pokemons")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, v := range m.Config.Db.GetAll() {
		fmt.Println("-", v.Name)
	}
	return nil
}

func (m *Menu) inspectCb(fields []string) error {
	if len(fields) != 2 {
		return errors.New("usage: inspect pokemon-name")
	}
	pokemon, ok := m.Config.Db.Find(fields[1])
	if !ok {
		return errors.New("you have not caught that pokemon")
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf("\t-%s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Println("\t-", t.Type.Name)
	}
	return nil
}

func (m *Menu) setMenuCommands() {
	help := MenuCommand{
		Name:            "help",
		Description:     "The help command",
		CommandCallback: m.helpCb,
	}
	m.addMenuCommand(help)

	exit := MenuCommand{
		Name:            "exit",
		Description:     "The exit command",
		CommandCallback: m.exitCb,
	}
	m.addMenuCommand(exit)

	mapCommand := MenuCommand{
		Name:            "map",
		Description:     "The map command fetches a page of information and advances",
		CommandCallback: m.mapCb,
	}
	m.addMenuCommand(mapCommand)

	mapbCommand := MenuCommand{
		Name:            "mapb",
		Description:     "The mapb command goes back a page",
		CommandCallback: m.mapbCb,
	}
	m.addMenuCommand(mapbCommand)

	state := MenuCommand{
		Name:            "state",
		Description:     "The state command to view the menu state",
		CommandCallback: m.stateCb,
	}
	m.addMenuCommand(state)

	cacheKeys := MenuCommand{
		Name:            "keys",
		Description:     "See all the keys in the cache",
		CommandCallback: m.cacheKeyCb,
	}
	m.addMenuCommand(cacheKeys)

	explore := MenuCommand{
		Name:            "explore",
		Description:     "explore an area, usage: explore area-name",
		CommandCallback: m.exploreCb,
	}
	m.addMenuCommand(explore)

	catch := MenuCommand{
		Name:            "catch",
		Description:     "catch a Pokemon, usage: catch pokemon-name",
		CommandCallback: m.catchCb,
	}
	m.addMenuCommand(catch)

	pokedex := MenuCommand{
		Name:            "pokedex",
		Description:     "list the pokedex",
		CommandCallback: m.pokedexCb,
	}
	m.addMenuCommand(pokedex)

	inspect := MenuCommand{
		Name:            "inspect",
		Description:     "Inspect a pokemon you have in your pokedex",
		CommandCallback: m.inspectCb,
	}
	m.addMenuCommand(inspect)

}
