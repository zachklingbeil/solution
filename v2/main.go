package v2

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

type Fx struct {
	Uptime []Uptime                       // StationID -> Uptime percentage
	Data   map[uint32]map[uint32][]Report // StationID -> ChargerID -> []Report
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

func NewFx() *Fx {
	return &Fx{
		Uptime: make([]Uptime, 0),
		Data:   make(map[uint32]map[uint32][]Report),
	}
}
func (fx *Fx) parseStationsLine(line string) error {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return fmt.Errorf("invalid station line: %s (expected at least 2 fields)", line)
	}

	stationID, err := strconv.ParseUint(fields[0], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid station ID: %s (must be a uint32)", fields[0])
	}

	station := uint32(stationID)
	if _, exists := fx.Data[station]; !exists {
		fx.Data[station] = make(map[uint32][]Report)
	}

	for _, chargerStr := range fields[1:] {
		chargerID, err := strconv.ParseUint(chargerStr, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid charger ID: %s (must be a uint32)", chargerStr)
		}
		charger := uint32(chargerID)
		if _, exists := fx.Data[station][charger]; !exists {
			fx.Data[station][charger] = []Report{}
		}
	}

	return nil
}

func (fx *Fx) parseChargerLine(line string) error {
	fields := strings.Fields(line)
	if len(fields) != 4 {
		return fmt.Errorf("invalid charger line: %s (expected exactly 4 fields)", line)
	}

	chargerID, err1 := strconv.ParseUint(fields[0], 10, 32)
	start, err2 := strconv.ParseUint(fields[1], 10, 64)
	end, err3 := strconv.ParseUint(fields[2], 10, 64)
	up, err4 := strconv.ParseBool(fields[3])

	if err1 != nil {
		return fmt.Errorf("invalid charger ID: %s (must be a uint32)", fields[0])
	}
	if err2 != nil {
		return fmt.Errorf("invalid start time: %s (must be a uint64)", fields[1])
	}
	if err3 != nil {
		return fmt.Errorf("invalid end time: %s (must be a uint64)", fields[2])
	}
	if err4 != nil {
		return fmt.Errorf("invalid up value: %s (must be a boolean)", fields[3])
	}

	charger := uint32(chargerID)
	for stationID, chargers := range fx.Data {
		if _, exists := chargers[charger]; exists {
			fx.Data[stationID][charger] = append(fx.Data[stationID][charger], Report{
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

func (fx *Fx) ParseReportsFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var parseFunc func(string) error

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		switch line {
		case Stations:
			parseFunc = fx.parseStationsLine
			continue
		case Chargers:
			parseFunc = fx.parseChargerLine
			continue
		}
		if parseFunc != nil {
			if err := parseFunc(line); err != nil {
				return fmt.Errorf("error parsing line: %s: %w", line, err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", path, err)
	}
	return nil
}

func (fx *Fx) CalculateUptime() {
	for stationID, chargers := range fx.Data {
		var totalUp, total uint64
		for _, reports := range chargers {
			if len(reports) == 0 {
				continue
			}
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
		fx.Uptime = append(fx.Uptime, Uptime{
			StationID: stationID,
			Percent:   percent,
		})
	}
}
func (fx *Fx) PrintStationUptime() {
	for _, uptime := range fx.Uptime {
		println(uptime.StationID, uptime.Percent)
	}
}
