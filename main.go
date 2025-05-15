package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type AvailabilityReport struct {
	ChargerID      uint32
	StartTimeNanos uint64
	EndTimeNanos   uint64
	Up             bool
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR")
		return
	}
	reports, err := ParseAvailabilityReports(os.Args[1])
	if err != nil {
		fmt.Println("ERROR")
		return
	}

	// Group reports by charger
	chargerReports := make(map[uint32][]AvailabilityReport)
	for _, r := range reports {
		chargerReports[r.ChargerID] = append(chargerReports[r.ChargerID], r)
	}

	// Get sorted charger IDs
	var chargerIDs []uint32
	for id := range chargerReports {
		chargerIDs = append(chargerIDs, id)
	}
	slices.Sort(chargerIDs)

	for _, cid := range chargerIDs {
		uptime := CalculateChargerUptime(chargerReports[cid])
		fmt.Printf("%d %d\n", cid, uptime)
	}
}

// ParseAvailabilityReports parses only the [Charger Availability Reports] section.
func ParseAvailabilityReports(path string) ([]AvailabilityReport, error) {
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
		if line == "[Charger Availability Reports]" {
			inSection = true
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			inSection = false
			continue
		}
		if inSection {
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
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

// CalculateChargerUptime computes the uptime percentage for a charger.
func CalculateChargerUptime(reports []AvailabilityReport) int {
	var totalUp, total uint64
	for _, r := range reports {
		duration := r.EndTimeNanos - r.StartTimeNanos
		total += duration
		if r.Up {
			totalUp += duration
		}
	}
	if total == 0 {
		return 0
	}
	return int((totalUp * 100) / total)
}
