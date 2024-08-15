package adapters

type GetAgenResult struct {
	Balance int64 `json:"balance"`
	Ships   int64 `json:"ships"`
	Probes  int64 `json:"probes"`
}
