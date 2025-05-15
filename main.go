package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR")
		return
	}
	NewElectric(os.Args[1])
}

func Output(m map[uint32]int) {
	keys := make([]uint32, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Printf("%d %d\n", k, m[k])
	}
}

const key = "[Charger Availability Reports]"

type Electric struct {
	Map map[uint32]int
}

type AvailabilityReport struct {
	ChargerID uint32
	Start     uint64
	End       uint64
	Up        bool
}

func NewElectric(path string) *Electric {
	elec := &Electric{
		Map: make(map[uint32]int),
	}
	reports := Input(path)
	elec.ComputeUptime(reports)
	Output(elec.Map)
	return elec
}

func Input(path string) []AvailabilityReport {
	file, err := os.Open(path)
	if err != nil {
		return nil
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
			return nil
		}
		chargerID, err := strconv.ParseUint(fields[0], 10, 32)
		if err != nil {
			return nil
		}
		start, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return nil
		}
		end, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			return nil
		}
		up, err := strconv.ParseBool(fields[3])
		if err != nil {
			return nil
		}
		reports = append(reports, AvailabilityReport{
			ChargerID: uint32(chargerID),
			Start:     start,
			End:       end,
			Up:        up,
		})
	}
	return reports
}

func (e *Electric) ComputeUptime(reports []AvailabilityReport) map[uint32]int {
	chargerReports := make(map[uint32][]AvailabilityReport)
	for _, r := range reports {
		chargerReports[r.ChargerID] = append(chargerReports[r.ChargerID], r)
	}
	for id, reps := range chargerReports {
		var totalUp, total uint64
		for _, rep := range reps {
			duration := rep.End - rep.Start
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
