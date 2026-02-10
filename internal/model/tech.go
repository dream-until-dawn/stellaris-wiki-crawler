package model

type Technology struct {
	ID            string
	Name          string
	Area          string // Physics / Society / Engineering
	Tier          int
	Category      string
	Cost          int
	Prerequisites []string
	Weight        int
	Effects       []string
	DLC           []string
	WikiURL       string
}
