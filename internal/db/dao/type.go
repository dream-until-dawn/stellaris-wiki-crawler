package dao

import (
	"stellarisWikiCrawler/internal/model"
)

type TechnologyTreeItem struct {
	Name           string            `json:"name"`
	Classify       string            `json:"classify"`
	Technology     string            `json:"technology"`
	Icon           string            `json:"icon"`
	Description    string            `json:"description"`
	Tier           string            `json:"tier"`
	Cost           string            `json:"cost"`
	EffectsUnlocks model.StringSlice `json:"effectsUnlocks"`
	Prerequisites  model.StringSlice `json:"prerequisites"`
	DrawWeight     model.StringSlice `json:"draw weight"`
	Empire         model.StringSlice `json:"empire"`
	DLC            model.StringSlice `json:"dlc"`
	Depth          int               `json:"depth"`
}

type GraphItem struct {
	Name           string            `json:"name"`
	Classify       string            `json:"classify"`
	Technology     string            `json:"technology"`
	Icon           string            `json:"icon"`
	Description    string            `json:"description"`
	Tier           string            `json:"tier"`
	Cost           string            `json:"cost"`
	EffectsUnlocks model.StringSlice `json:"effectsUnlocks"`
	Prerequisites  model.StringSlice `json:"prerequisites"`
	DrawWeight     model.StringSlice `json:"draw weight"`
	Empire         model.StringSlice `json:"empire"`
	DLC            model.StringSlice `json:"dlc"`
	Source         string            `json:"source"`
	Target         string            `json:"target"`
}
