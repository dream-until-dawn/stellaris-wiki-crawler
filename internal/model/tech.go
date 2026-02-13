package model

import (
	"database/sql/driver"
	"encoding/json"
)

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = StringSlice{} // 🔥 改这里
		return nil
	}

	switch v := value.(type) {
	case []byte:
		if len(v) == 0 {
			*s = StringSlice{}
			return nil
		}
		return json.Unmarshal(v, s)
	case string:
		if v == "" {
			*s = StringSlice{}
			return nil
		}
		return json.Unmarshal([]byte(v), s)
	default:
		*s = StringSlice{}
		return nil
	}
}

func (s StringSlice) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(s))
}

type CrawierTechnology struct {
	Title       string                  `json:"title"`
	Description StringSlice             `json:"description"`
	Table       []CrawierTechnologyItem `json:"table"`
}

type CrawierTechnologyItem struct {
	Icon       string `json:"Icon"`
	Technology struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"Technology"`
	Tier           string      `json:"Tier"`
	Cost           string      `json:"Cost"`
	EffectsUnlocks StringSlice `json:"Effects / unlocks"`
	Prerequisites  StringSlice `json:"Prerequisites"`
	DrawWeight     StringSlice `json:"Draw weight"`
	Empire         StringSlice `json:"Empire"`
	DLC            StringSlice `json:"DLC"`
}
