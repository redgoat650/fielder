package fielder

// SeasonConfig contains options for initializing a season
type SeasonConfig struct {
	InningsPerGame  int
	NumGames        int
	MinGirlsPerGame int
	MinGuysPerGame  int
}

// NewSeasonFromCfg initializes and returns a pointer to a new Season,
// configured from the give SeasonConfig
func NewSeasonFromCfg(cfg *SeasonConfig) *Season {
	return NewSeason(cfg.NumGames, cfg.InningsPerGame)
}
