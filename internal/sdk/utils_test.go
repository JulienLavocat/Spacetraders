package sdk_test

import (
	"testing"

	"github.com/julienlavocat/spacetraders/internal/sdk"
	"github.com/stretchr/testify/assert"
)

func TestGetDistance(t *testing.T) {
	dst := sdk.GetDistance(0, 0, 0, 0)
	assert.Equal(t, float64(0), dst, "expected dst to be 0")
	dst = sdk.GetDistance(40, 150, -112, 1228)
	assert.Equal(t, float64(1088.6634006891202), dst, "expected dst to be 1088.663")
	dst = sdk.GetDistance(-112, 1228, -112, 1228)
	assert.Equal(t, float64(0), dst, "expected dst to be 0")
}

func TestGetFuelCost(t *testing.T) {
	dst := sdk.GetFuelCost(0, 0, 0, 0)
	assert.Equal(t, int32(1), dst, "expected cost to be rounded to be at least 1")
	dst = sdk.GetFuelCost(40, 150, -112, 1228)
	assert.Equal(t, int32(1089), dst, "expected cost to be 1089")
	dst = sdk.GetFuelCost(-112, 1228, -112, 1228)
	assert.Equal(t, int32(1), dst, "expected cost to be 1")
}
