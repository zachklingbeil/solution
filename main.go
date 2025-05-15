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
	elec := NewElectric(os.Args[1])
	Output(elec.UptimeMap)
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
	Reports   []AvailabilityReport
	ReportMap map[uint32][]AvailabilityReport
	UptimeMap map[uint32]int
}

type AvailabilityReport struct {
	ChargerID uint32
	Start     uint64
	End       uint64
	Up        bool
}

func NewElectric(path string) *Electric {
	elec := &Electric{
		ReportMap: make(map[uint32][]AvailabilityReport),
		UptimeMap: make(map[uint32]int),
		Reports:   make([]AvailabilityReport, 0),
	}
	elec.ParseAvailabilityReports(path)
	elec.BuildReportMap()
	elec.ComputeUptime()
	return elec
}

func (e *Electric) ParseAvailabilityReports(path string) error {
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
		if len(fields) < 4 || len(fields)%4 != 0 {
			// Skip lines that don't have at least one complete record
			continue
		}
		for i := 0; i+3 < len(fields); i += 4 {
			chargerID, err := strconv.ParseUint(fields[i], 10, 32)
			if err != nil {
				continue
			}
			start, err := strconv.ParseUint(fields[i+1], 10, 64)
			if err != nil {
				continue
			}
			end, err := strconv.ParseUint(fields[i+2], 10, 64)
			if err != nil {
				continue
			}
			up, err := strconv.ParseBool(fields[i+3])
			if err != nil {
				continue
			}
			reports = append(reports, AvailabilityReport{
				ChargerID: uint32(chargerID),
				Start:     start,
				End:       end,
				Up:        up,
			})
		}
	}
	e.Reports = reports
	return nil
}

func (e *Electric) BuildReportMap() error {
	for _, r := range e.Reports {
		e.ReportMap[r.ChargerID] = append(e.ReportMap[r.ChargerID], r)
	}
	return nil
}

func (e *Electric) ComputeUptime() error {
	for id, reps := range e.ReportMap {
		var totalUp, total uint64
		for _, rep := range reps {
			duration := rep.End - rep.Start
			total += duration
			if rep.Up {
				totalUp += duration
			}
		}
		if total == 0 {
			e.UptimeMap[id] = 0
		} else {
			e.UptimeMap[id] = int((totalUp * 100) / total)
		}
	}
	return nil
}
