package catching

func GetProbability(baseExperience int) float64 {
	if baseExperience <= 255 {
		return 0.7 - float64(baseExperience)*(0.7-0.2)/255.0
	}
	// 635 is the maximum baseExperience (for blissey). Then 635 - 256 = 379
	return 0.2 - float64(baseExperience-256)*(0.2-0.15)/379.0
}
