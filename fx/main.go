package fx

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	Stations = "[Stations]"
	Chargers = "[Charger Availability Reports]"
)

type Era struct {
	Source map[uint32]map[uint32][]Report // StationID -> ChargerID -> []Report
	Uptime []Uptime                       // StationID -> Uptime percentag
}

type Report struct {
	ChargerID uint32
	Start     uint64 // nanos
	End       uint64 // nanos
	Up        bool
}

type Uptime struct {
	StationID uint32
	Percent   int
}

func Electric() *Era {
	return &Era{
		Uptime: make([]Uptime, 0),
		Source: make(map[uint32]map[uint32][]Report),
	}
}

func (e *Era) Input(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// Process the [Stations] section
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == Chargers {
			break // Exit when reaching the [Chargers] section
		}
		if line != Stations {
			e.parseStations(line)
		}
	}
	// Process the [Chargers] section
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		e.parseChargerLine(line)
	}
	return nil
}

// parseStations parses a line describing a station and its chargers.
// It adds the station and its chargers to the Source map.
func (e *Era) parseStations(line string) error {
	value := strings.Fields(line)
	stationID, _ := strconv.ParseUint(value[0], 10, 32) // Parse the station ID from the first field
	station := uint32(stationID)
	if _, exists := e.Source[station]; !exists {
		e.Source[station] = make(map[uint32][]Report) // Initialize the chargers map for this station if not present
	}

	for _, chargerStr := range value[1:] {
		id, _ := strconv.ParseUint(chargerStr, 10, 32) // Parse each charger ID
		charger := uint32(id)
		if _, exists := e.Source[station][charger]; !exists {
			e.Source[station][charger] = []Report{} // Initialize an empty report slice for each charger
		}
	}
	return nil
}

// parseChargerLine parses a line describing a charger report and adds it to the correct station/charger.
// Returns an error if the charger is not found in any station.
func (e *Era) parseChargerLine(line string) error {
	value := strings.Fields(line)
	id, _ := strconv.ParseUint(value[0], 10, 32)    // Parse charger IDw
	start, _ := strconv.ParseUint(value[1], 10, 64) // Parse start time
	end, _ := strconv.ParseUint(value[2], 10, 64)   // Parse end time
	up, _ := strconv.ParseBool(value[3])            // Parse up status (true/false)

	charger := uint32(id)
	for stationID, chargers := range e.Source {
		if _, exists := chargers[charger]; exists {
			// Append the report to the correct charger under the correct station
			e.Source[stationID][charger] = append(e.Source[stationID][charger], Report{
				ChargerID: charger,
				Start:     start,
				End:       end,
				Up:        up,
			})
			return nil // Successfully added, exit
		}
	}
	return fmt.Errorf("charger ID %d does not belong to any station", id) // Error if charger not found
}

func (e *Era) Fx() {
	for stationID, chargers := range e.Source {
		var totalUp, total uint64
		for _, reports := range chargers {
			minStart, maxEnd, upTime := reports[0].Start, reports[0].End, uint64(0)
			for _, r := range reports {
				if r.Start < minStart {
					minStart = r.Start
				}
				if r.End > maxEnd {
					maxEnd = r.End
				}
				if r.Up {
					upTime += r.End - r.Start
				}
			}
			total += maxEnd - minStart
			totalUp += upTime
		}
		percent := 0
		if total > 0 {
			percent = int((totalUp * 100) / total)
		}
		e.Uptime = append(e.Uptime, Uptime{
			StationID: stationID,
			Percent:   percent,
		})
		sort.Slice(e.Uptime, func(i, j int) bool {
			return e.Uptime[i].StationID < e.Uptime[j].StationID
		})
	}
}

func (e *Era) Output() {
	for _, uptime := range e.Uptime {
		fmt.Println(uptime.StationID, uptime.Percent)
	}
}
