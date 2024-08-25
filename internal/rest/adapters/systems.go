package adapters

import "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"

type System struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	X    int32  `json:"x"`
	Y    int32  `json:"y"`
}

func AdaptSystems(models []model.Systems) []System {
	systems := make([]System, len(models))
	for i := range models {
		model := models[i]
		systems[i] = System{
			Id:   model.ID,
			Type: model.Type,
			X:    model.X,
			Y:    model.Y,
		}
	}

	return systems
}
