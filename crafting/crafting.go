package crafting

type Recipe struct {
	Parts        []Part  `json:"parts"`
	OutputAmount float32 `json:"outputAmount"`
	CraftingTime float32 `json:"craftingTime"`
}

type Part struct {
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
}
