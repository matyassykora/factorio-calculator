package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/matyassykora/factorio-calculator/crafting"
)

type UserSettings struct {
	DrillType   string
	FurnaceType string
	BeltType    string
	StackSize   float32
}

type Config struct {
	Furnace Furnace
	Belt    Belt
	Drill   Drill
}

type Furnace struct {
	BaseSmeltingSpeed   float32 `json:"baseSmeltingSpeed"`
	BeaconModifier      float32 `json:"beaconModifier"`
	FuelConsumption     float32 `json:"fuelConsumption"`
	EnergyConsumption   float32 `json:"energyConsumption"`
	PollutionProduction float32 `json:"pollutionProduction"`
}

type Belt struct {
	Speed float32 `json:"speed"`
}

type Drill struct {
	MiningSpeed float32 `json:"miningSpeed"`
	Pollution   float32 `json:"pollution"`
}

type Result struct {
	Output    float32
	Energy    float32
	Fuel      float32
	Pollution float32
}

type SmeltingResult struct {
	Furnaces float32
	Drills   float32
	Ore      float32
	Belts    float32
}

func FromDesired(ore *Smeltable, drill *Drill, furnace *Furnace, belt *Belt, stackSize float32, desired float32) *SmeltingResult {
	return &SmeltingResult{
		Furnaces: FurnaceCount(ore, furnace, desired),
		Drills:   DrillCount(ore, drill, desired),
		Ore:      OreCount(ore, drill, desired),
		Belts:    BeltCount(belt, stackSize, desired),
	}
}

func FromBuildingCount(ore *Smeltable, drill *Drill, furnace *Furnace, belt *Belt, stackSize float32, furnaceCount float32) *SmeltingResult {
	desired := ore.SmeltingSpeed * ore.MiningTime * furnace.BaseSmeltingSpeed * furnaceCount
	return FromDesired(ore, drill, furnace, belt, stackSize, desired)
}

func FurnaceCount(smeltable *Smeltable, furnace *Furnace, desired float32) float32 {
	return desired / (smeltable.SmeltingSpeed * furnace.BaseSmeltingSpeed)
}

func DrillCount(smeltable *Smeltable, drill *Drill, desired float32) float32 {
	return desired / drill.MiningSpeed / smeltable.MiningTime
}

func OreCount(smeltable *Smeltable, drill *Drill, desired float32) float32 {
	return DrillCount(smeltable, drill, desired) * drill.MiningSpeed
}

func BeltCount(belt *Belt, stackSize float32, desired float32) float32 {
	return desired / belt.Speed / stackSize
}

type Smeltable struct {
	SmeltingSpeed float32 `json:"smeltingSpeed"`
	MiningTime    float32 `json:"miningTime"`
}

func Mine(drill *Drill, ore *Smeltable, drillCount float32) float32 {
	return drillCount * drill.MiningSpeed
}

func MineDesired(drill *Drill, ore *Smeltable, desired float32) float32 {
	return desired * ore.MiningTime / drill.MiningSpeed
}

func loadDefs[T any](file string) (map[string]T, error) {
	defBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	values := map[string]T{}
	err = json.Unmarshal(defBytes, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	smeltableDefs, err := loadDefs[Smeltable]("defs/smeltables.json")
	check(err)

	furnaceDefs, err := loadDefs[Furnace]("defs/furnaces.json")
	check(err)

	beltDefs, err := loadDefs[Belt]("defs/belts.json")
	check(err)

	drillDefs, err := loadDefs[Drill]("defs/drills.json")
	check(err)

	userSettings := UserSettings{
		DrillType:   "electric-mining-drill",
		BeltType:    "fast-transport-belt",
		FurnaceType: "stone-furnace",
		StackSize:   1,
	}

	oreKey := "iron-ore"

	ore, ok := smeltableDefs[oreKey]
	if ok != true {
		panic("key not found: " + oreKey)
	}

	furnace, ok := furnaceDefs[userSettings.FurnaceType]
	if ok != true {
		panic("key not found: " + userSettings.FurnaceType)
	}

	belt, ok := beltDefs[userSettings.BeltType]
	if ok != true {
		panic("key not found: " + userSettings.BeltType)
	}

	drill, ok := drillDefs[userSettings.DrillType]
	if ok != true {
		panic("key not found: " + userSettings.DrillType)
	}

	var desired float32 = 10

	fmt.Printf("want '%.3f' '%+v'\n", desired, oreKey)
	fmt.Printf("smelted in '%s'\n", userSettings.FurnaceType)
	fmt.Printf("mined with '%s'\n", userSettings.DrillType)
	fmt.Printf("transported on '%s' stacked '%.0f' high'\n", userSettings.BeltType, userSettings.StackSize)

	res := FromDesired(&ore, &drill, &furnace, &belt, userSettings.StackSize, desired)

	fmt.Printf("amount of '%s' needed: %+v\n", userSettings.FurnaceType, res.Furnaces)
	fmt.Printf("amount of '%s' needed: %.4f\n", userSettings.DrillType, res.Drills)
	fmt.Printf("amount of '%s' needed: %.4f\n", oreKey, res.Ore)
	fmt.Printf("amount of '%s' stacked '%.0f' tall needed: %+v\n", userSettings.BeltType, userSettings.StackSize, res.Belts)

	fmt.Println("")

	var buildingCount float32 = 3
	fmt.Printf("want '%.0f' '%s'\n", buildingCount, userSettings.FurnaceType)
	fmt.Printf("smelting '%s'\n", oreKey)
	fmt.Printf("mined with '%s'\n", userSettings.DrillType)
	fmt.Printf("transported on '%s' stacked '%.0f' high'\n", userSettings.BeltType, userSettings.StackSize)

	res2 := FromBuildingCount(&ore, &drill, &furnace, &belt, userSettings.StackSize, buildingCount)
	fmt.Printf("amount of '%s' needed: %+v\n", userSettings.FurnaceType, res2.Furnaces)
	fmt.Printf("amount of '%s' needed: %.4f\n", userSettings.DrillType, res2.Drills)
	fmt.Printf("amount of '%s' needed: %.4f\n", oreKey, res2.Ore)
	fmt.Printf("amount of '%s' stacked '%.0f' tall needed: %+v\n", userSettings.BeltType, userSettings.StackSize, res2.Belts)

	recipes, err := loadDefs[crafting.Recipe]("defs/recipes.json")
	check(err)

	for k, recipe := range recipes {
		fmt.Println("\n######################################")
		fmt.Printf("recipe: %s\n", k)
		i := 0
		printParts(recipes, recipe, &i)
		fmt.Println("######################################")
	}
}

func printParts(recipes map[string]crafting.Recipe, recipe crafting.Recipe, index *int) {
	indent := func() string {
		s := "\t"
		res := ""
		for range *index {
			res += s
		}
		return res
	}
	fmt.Printf("%s crafting time: %.4f\n", indent(), recipe.CraftingTime)
	fmt.Printf("%s output amount: %.4f\n", indent(), recipe.OutputAmount)
	for _, a2 := range recipe.Parts {
		fmt.Printf("%s part: %s\n", indent(), a2.Name)
		fmt.Printf("%s amount: %.4f\n", indent(), a2.Amount)
		recipe, ok := recipes[a2.Name]
		if !ok {
			panic("part doesn't exist as a recipe: " + a2.Name)
		}
		if recipe.Parts == nil {
			if *index != 0 {
				*index -= 1
			}
		}
		if recipe.Parts != nil {
			fmt.Printf("%s to craft this part:\n", indent())
			*index += 1
			printParts(recipes, recipe, index)
		}
	}
}

type Graph struct {
	Vertices map[string]*Vertex
}

func NewGraph(opts ...GraphOption) *Graph {
	g := &Graph{Vertices: map[string]*Vertex{}}
	for _, opt := range opts {
		opt(g)
	}

	return g
}

type GraphOption func(g *Graph)

func WithAdjacencyList(list map[string][]string) GraphOption {
	return func(g *Graph) {
		for vertex, edges := range list {
			if _, ok := g.Vertices[vertex]; !ok {
				g.AddVertex(vertex, vertex)
			}

			for _, edge := range edges {
				if _, ok := g.Vertices[edge]; !ok {
					g.AddVertex(edge, edge)
				}
				g.AddEdge(vertex, edge, 0)
			}
		}
	}
}

func (g *Graph) AddVertex(key, val string) {
	g.Vertices[key] = &Vertex{Val: val, Edges: map[string]*Edge{}}
}

func (g *Graph) AddEdge(sourceKey string, destinationKey string, weight int) {
	if _, ok := g.Vertices[sourceKey]; !ok {
		return
	}
	if _, ok := g.Vertices[destinationKey]; !ok {
		return
	}

	g.Vertices[sourceKey].Edges[destinationKey] = &Edge{Weight: weight, Vertex: g.Vertices[destinationKey]}
}

func (g *Graph) Neighbors(sourceKey string) []string {
	res := []string{}

	for _, edge := range g.Vertices[sourceKey].Edges {
		res = append(res, edge.Vertex.Val)
	}

	return res
}

type Vertex struct {
	Val   string
	Edges map[string]*Edge
}

type Edge struct {
	Weight int
	Vertex *Vertex
}

func findRecipeTree(recipes map[string]crafting.Recipe, recipeName string) {

}
