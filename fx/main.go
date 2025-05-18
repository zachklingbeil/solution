package fx

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	Stations = "[Stations]"
	Chargers = "[Charger Availability Reports]"
)

type Era struct {
	Uptime []Uptime                       // StationID -> Uptime percentage
	Source map[uint32]map[uint32][]Report // StationID -> ChargerID -> []Report
}

type Report struct {
	ChargerID uint32
	Start     uint64
	End       uint64
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

func (e *Era) parseStations(line string) error {
	fields := strings.Fields(line)
	stationID, _ := strconv.ParseUint(fields[0], 10, 32)
	station := uint32(stationID)
	if _, exists := e.Source[station]; !exists {
		e.Source[station] = make(map[uint32][]Report)
	}

	for _, chargerStr := range fields[1:] {
		chargerID, _ := strconv.ParseUint(chargerStr, 10, 32)
		charger := uint32(chargerID)
		if _, exists := e.Source[station][charger]; !exists {
			e.Source[station][charger] = []Report{}
		}
	}
	return nil
}

func (e *Era) parseChargerLine(line string) error {
	fields := strings.Fields(line)
	if len(fields) != 4 {
		return fmt.Errorf("invalid charger line: %s (expected exactly 4 fields)", line)
	}

	chargerID, _ := strconv.ParseUint(fields[0], 10, 32)
	start, _ := strconv.ParseUint(fields[1], 10, 64)
	end, _ := strconv.ParseUint(fields[2], 10, 64)
	up, _ := strconv.ParseBool(fields[3])

	charger := uint32(chargerID)
	for stationID, chargers := range e.Source {
		if _, exists := chargers[charger]; exists {
			e.Source[stationID][charger] = append(e.Source[stationID][charger], Report{
				ChargerID: charger,
				Start:     start,
				End:       end,
				Up:        up,
			})
			return nil
		}
	}
	return fmt.Errorf("charger ID %d does not belong to any station", chargerID)
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
	}
}

func (e *Era) Output() {
	for _, uptime := range e.Uptime {
		println(uptime.StationID, uptime.Percent)
	}
}
