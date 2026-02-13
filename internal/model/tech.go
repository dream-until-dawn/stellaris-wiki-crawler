package model

type CrawierTechnology struct {
	Title       string                  `json:"title"`
	Description []string                `json:"description"`
	Table       []CrawierTechnologyItem `json:"table"`
}

type CrawierTechnologyItem struct {
	Icon       string `json:"Icon"`
	Technology struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"Technology"`
	Tier           string   `json:"Tier"`
	Cost           string   `json:"Cost"`
	EffectsUnlocks []string `json:"Effects / unlocks"`
	Prerequisites  []string `json:"Prerequisites"`
	DrawWeight     []string `json:"Draw weight"`
	Empire         []string `json:"Empire"`
	DLC            []string `json:"DLC"`
}
