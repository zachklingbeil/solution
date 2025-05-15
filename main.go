package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const key = "[Charger Availability Reports]"

type Electric struct {
	Map map[uint32]int
}

func NewElectric(path string) (*Electric, error) {
	elec := &Electric{
		Map: make(map[uint32]int),
	}
	reports, err := elec.Input(path)
	if err != nil {
		return nil, err
	}
	elec.ComputeUptime(reports)
	return elec, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR")
		return
	}
	e, err := NewElectric(os.Args[1])
	if err != nil {
		fmt.Println("ERROR")
		return
	}
	for key, value := range e.Map {
		fmt.Printf("Charger %d: %d%% uptime\n", key, value)
	}
}

type AvailabilityReport struct {
	ChargerID      uint32
	StartTimeNanos uint64
	EndTimeNanos   uint64
	Up             bool
}

func (e *Electric) ComputeUptime(reports []AvailabilityReport) map[uint32]int {
	chargerReports := make(map[uint32][]AvailabilityReport)
	for _, r := range reports {
		chargerReports[r.ChargerID] = append(chargerReports[r.ChargerID], r)
	}
	for id, reps := range chargerReports {
		var totalUp, total uint64
		for _, rep := range reps {
			duration := rep.EndTimeNanos - rep.StartTimeNanos
			total += duration
			if rep.Up {
				totalUp += duration
			}
		}
		if total == 0 {
			e.Map[id] = 0
		} else {
			e.Map[id] = int((totalUp * 100) / total)
		}
	}
	return e.Map
}

func (e *Electric) Input(path string) ([]AvailabilityReport, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var reports []AvailabilityReport
	scanner := bufio.NewScanner(file)
	inSection := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if !inSection {
			if line == key {
				inSection = true
			}
			continue
		}
		// If we hit another section, stop processing
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			break
		}
		fields := strings.Fields(line)
		if len(fields) != 4 {
			return nil, fmt.Errorf("invalid availability line")
		}
		chargerID, err := strconv.ParseUint(fields[0], 10, 32)
		if err != nil {
			return nil, err
		}
		start, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return nil, err
		}
		end, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			return nil, err
		}
		up, err := strconv.ParseBool(fields[3])
		if err != nil {
			return nil, err
		}
		reports = append(reports, AvailabilityReport{
			ChargerID:      uint32(chargerID),
			StartTimeNanos: start,
			EndTimeNanos:   end,
			Up:             up,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}
