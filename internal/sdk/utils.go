package sdk

import "math"

func GetDistance(x1, y1, x2, y2 int32) float64 {
	deltaX := float64(x1 - x2)
	deltaY := float64(y1 - y2)
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}

// Source: https://github.com/SpaceTradersAPI/api-docs/wiki/Travel-Fuel-and-Time#travel-time
func GetFuelCost(x1, y1, x2, y2 int32) int32 {
	return max(1, GetFuelFromDistance(GetDistance(x1, y1, x2, y2)))
}

func GetFuelFromDistance(distance float64) int32 {
	return int32(math.Round(distance))
}
