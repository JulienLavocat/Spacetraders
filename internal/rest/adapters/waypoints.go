package adapters

import (
	"time"

	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
)

type Waypoint struct {
	SubmittedOn       *time.Time `json:"submittedOn"`
	Orbits            *string    `json:"orbits"`
	SubmittedBy       *string    `json:"submittedBy"`
	Id                string     `json:"id"`
	Type              string     `json:"type"`
	Faction           string     `json:"faction"`
	X                 int32      `json:"x"`
	Y                 int32      `json:"y"`
	UnderConstruction bool       `json:"isUnderConstruction"`
}

func adaptWaypoint(wp model.Waypoints) Waypoint {
	return Waypoint{
		Id:                wp.ID,
		SubmittedOn:       wp.SubmittedOn,
		Type:              wp.Type,
		Faction:           wp.Faction,
		Orbits:            wp.Orbits,
		SubmittedBy:       wp.SubmittedBy,
		X:                 wp.X,
		Y:                 wp.Y,
		UnderConstruction: wp.UnderConstruction,
	}
}

func AdaptWaypoints(wps []model.Waypoints) []Waypoint {
	waypoints := make([]Waypoint, len(wps))
	for i := range waypoints {
		waypoints[i] = adaptWaypoint(wps[i])
	}
	return waypoints
}
