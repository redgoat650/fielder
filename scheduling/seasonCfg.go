package fielder

type SeasonConfig struct {
	InningsPerGame  int
	NumGames        int
	MinGirlsPerGame int
	MinGuysPerGame  int
}

func NewSeasonFromCfg(cfg *SeasonConfig) *Season {
	return NewSeason(cfg.NumGames, cfg.InningsPerGame)
}
