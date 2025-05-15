package fx

import "github.com/zachklingbeil/electric/input"

type Fx struct {
	Input *input.Input
}

func NewFx(inputFile string) *Fx {
	return &Fx{
		Input: &input.Input{
			Map: make(map[string][]any),
		},
	}
}

// Station represents a station with a unique StationID and a list of ChargerIDs.
type Station struct {
	StationID uint32
	Map       []map[uint32]AvailabilityReport
}

// AvailabilityReport represents a single availability entry for a charger.
type AvailabilityReport struct {
	ChargerID uint32 // Unique charger ID
	Start     uint64 // Start time in nanoseconds
	End       uint64 // End time in nanoseconds
	Up        bool   // True if charger was up (available)
}
