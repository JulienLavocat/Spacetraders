package sdk

type Coordinates struct {
	X int32
	Y int32
}

func NewCoordinates(x int32, y int32) Coordinates {
	return Coordinates{X: x, Y: y}
}
