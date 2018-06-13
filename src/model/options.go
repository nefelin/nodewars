package model

// new model options

func CoinGoal(goal int) func(*GameModel) error {
	return func(m *GameModel) error {
		return m.SetCoinGoal(goal)
	}
}

func Languages(langs ...Language) func(*GameModel) error {
	return func(m *GameModel) error {
		return m.SetLanguages(langs)
	}
}

func TimeLimit(limit int) func(*GameModel) error {
	return func(m *GameModel) error {
		return m.SetTimeLimit(limit)
	}
}

func MapMode(m *GameModel) error {
	return m.SetMode(modes.Map)
}
