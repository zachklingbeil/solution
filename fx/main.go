package fx

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

// key is the section header for charger availability reports in the input file.
const key = "[Charger Availability Reports]"

// Lines in the input file are treated as rows.
// Values within a row are separated by a single space.
// Each line has 4 values.
type Report struct {
	ChargerID uint32
	Start     uint64 // nano timestamp
	End       uint64 // nano timestamp
	Up        bool   // or down
}

// The Report struct matches both the order and type of values found in each row.
// Fields in the Report struct provide columns for working with our tabular data source.
// []Report scales automatically to match the number of lines within the [Charger Availability Reports] section.
type Era struct {
	Reports   []Report
	ReportMap map[uint32][]Report
	Uptime    map[uint32]int
}

// NewEra creates and returns a new Era instance with initialized fields.
func Electric() *Era {
	return &Era{
		ReportMap: make(map[uint32][]Report),
		Uptime:    make(map[uint32]int),
		Reports:   make([]Report, 0),
	}
}

// Input reads and parses charger reports from the specified file path.
// Returns an error if the file cannot be opened.
func (e *Era) Input(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var reports []Report
	scanner := bufio.NewScanner(file)
	inSection := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// Look for the section header
		if !inSection {
			if line == key {
				inSection = true
			}
			continue
		}
		// Stop if a new section starts
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			break
		}
		fields := strings.Fields(line)
		if len(fields) < 4 || len(fields)%4 != 0 {
			continue
		}
		// Parse each report in the line
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
			// Create a new report and append it to the reports slice
			reports = append(reports, Report{
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

// Fx maps reports by charger ID and calculates the uptime percentage for each charger.
// Uptime is an int representing the percentage of time a charger was reported "up" over the total reported time.
func (e *Era) Fx() {
	for _, r := range e.Reports {
		e.ReportMap[r.ChargerID] = append(e.ReportMap[r.ChargerID], r)
	}

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
			e.Uptime[id] = 0
		} else {
			// rounded down to the nearest percent
			e.Uptime[id] = int((totalUp * 100) / total)
		}
	}
}

// Output prints uptime percentages for each charger ID, in ascending order, as unformatted ints.
func (e *Era) Output(m map[uint32]int) {
	keys := make([]uint32, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Printf("%d %d\n", k, m[k])
	}
}
