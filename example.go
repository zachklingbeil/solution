package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"sort"
// 	"strconv"
// 	"strings"
// )

// // --- Section Constants ---
// const (
// 	SectionStations = "[Stations]"
// 	SectionReports  = "[Charger Availability Reports]"
// )

// // --- Data Structures ---

// type Electric struct {
// 	Filename string
// 	Map      map[string]map[uint32]any // section -> id -> *Station or *Charger
// }

// type Station struct {
// 	ID         uint32
// 	ChargerIDs []uint32
// }

// type Charger struct {
// 	ID      uint32
// 	Reports []Report
// }

// type Report struct {
// 	Start uint64
// 	End   uint64
// 	Up    bool
// }

// func NewElectric(filename string) *Electric {
// 	return &Electric{
// 		Filename: filename,
// 		Map: map[string]map[uint32]any{
// 			SectionStations: make(map[uint32]any),
// 			SectionReports:  make(map[uint32]any),
// 		},
// 	}
// }
// func (e *Electric) ReadRawFile() ([]byte, error) {
// 	return os.ReadFile(e.Filename)
// }

// // --- Parsing Methods ---

// func (e *Electric) ParseFile() error {
// 	file, err := os.Open(e.Filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	var currentSection string
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := strings.TrimSpace(scanner.Text())
// 		if line == "" || strings.HasPrefix(line, "//") {
// 			continue
// 		}
// 		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
// 			currentSection = line
// 			continue
// 		}
// 		switch currentSection {
// 		case SectionStations:
// 			station := parseStationLine(line)
// 			if station != nil {
// 				e.Map[SectionStations][station.ID] = station
// 			}
// 		case SectionReports:
// 			report, chargerID, ok := parseChargerReportLine(line)
// 			if ok {
// 				// Get or create Charger
// 				var ch *Charger
// 				val, exists := e.Map[SectionReports][chargerID]
// 				if !exists {
// 					ch = &Charger{ID: chargerID}
// 					e.Map[SectionReports][chargerID] = ch
// 				} else {
// 					ch = val.(*Charger)
// 				}
// 				ch.Reports = append(ch.Reports, report)
// 			}
// 		}
// 	}
// 	return scanner.Err()
// }

// // Returns a Station struct or nil if the line is invalid.
// func parseStationLine(line string) *Station {
// 	parts := strings.Fields(line)
// 	if len(parts) < 2 {
// 		return nil
// 	}
// 	stationID, err := strconv.ParseUint(parts[0], 10, 32)
// 	if err != nil {
// 		return nil
// 	}
// 	chargerIDs := make([]uint32, 0, len(parts)-1)
// 	for _, chargerStr := range parts[1:] {
// 		chargerID, err := strconv.ParseUint(chargerStr, 10, 32)
// 		if err != nil {
// 			return nil
// 		}
// 		chargerIDs = append(chargerIDs, uint32(chargerID))
// 	}
// 	return &Station{
// 		ID:         uint32(stationID),
// 		ChargerIDs: chargerIDs,
// 	}
// }

// // Returns (Report, chargerID, ok) or (zero, 0, false) if invalid.
// func parseChargerReportLine(line string) (Report, uint32, bool) {
// 	parts := strings.Fields(line)
// 	if len(parts) < 4 {
// 		return Report{}, 0, false
// 	}
// 	chargerID, err1 := strconv.ParseUint(parts[0], 10, 32)
// 	start, err2 := strconv.ParseUint(parts[1], 10, 64)
// 	end, err3 := strconv.ParseUint(parts[2], 10, 64)
// 	up, err4 := strconv.ParseBool(parts[3])
// 	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
// 		return Report{}, 0, false
// 	}
// 	return Report{Start: start, End: end, Up: up}, uint32(chargerID), true
// }

// // --- Uptime Calculation ---

// func (e *Electric) UptimeByStation() map[uint32]int {
// 	uptimeByStation := make(map[uint32]int)
// 	stations := e.Map[SectionStations]
// 	chargers := e.Map[SectionReports]
// 	for stationID, s := range stations {
// 		station := s.(*Station)
// 		var totalTime, upTime uint64
// 		for _, chargerID := range station.ChargerIDs {
// 			chVal, ok := chargers[chargerID]
// 			if !ok {
// 				continue
// 			}
// 			ch := chVal.(*Charger)
// 			if len(ch.Reports) == 0 {
// 				continue
// 			}
// 			// Sort reports by start time
// 			reports := ch.Reports
// 			sort.Slice(reports, func(i, j int) bool { return reports[i].Start < reports[j].Start })

// 			// Calculate total span (including gaps)
// 			spanStart := reports[0].Start
// 			spanEnd := reports[0].End
// 			for _, r := range reports {
// 				if r.End > spanEnd {
// 					spanEnd = r.End
// 				}
// 			}
// 			chargerTotal := spanEnd - spanStart

// 			// Calculate up time and fill gaps as downtime
// 			var lastEnd = spanStart
// 			var chargerUp uint64
// 			for _, r := range reports {
// 				// Gap between lastEnd and this start is downtime
// 				if r.Start > lastEnd {
// 					// downtime, do nothing (not added to up time)
// 				}
// 				// Up time only if marked up
// 				if r.Up {
// 					chargerUp += r.End - r.Start
// 				}
// 				if r.End > lastEnd {
// 					lastEnd = r.End
// 				}
// 			}
// 			totalTime += chargerTotal
// 			upTime += chargerUp
// 		}
// 		if totalTime == 0 {
// 			uptimeByStation[stationID] = 0
// 		} else {
// 			uptimeByStation[stationID] = int((upTime * 100) / totalTime)
// 		}
// 	}
// 	return uptimeByStation
// }

// // --- Main ---

// func main() {
// 	inputFiles := []string{"input1.txt", "input2.txt"}
// 	for _, fname := range inputFiles {
// 		fmt.Printf("Results for %s:\n", fname)
// 		electric := NewElectric(fname)
// 		if err := electric.ParseFile(); err != nil {
// 			fmt.Println("Error:", err)
// 			continue
// 		}
// 		uptimeByStation := electric.UptimeByStation()
// 		var stationIDs []int
// 		for id := range uptimeByStation {
// 			stationIDs = append(stationIDs, int(id))
// 		}
// 		sort.Ints(stationIDs)
// 		for _, id := range stationIDs {
// 			fmt.Printf("Station ID %d uptime %d%%\n", id, uptimeByStation[uint32(id)])
// 		}
// 		fmt.Println()
// 	}
// }

// func DefineHaystack(filename string, keyword string) map[string][][]any {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return nil
// 	}
// 	defer file.Close()

// 	result := make(map[string][][]any)
// 	scanner := bufio.NewScanner(file)
// 	inSection := false
// 	for scanner.Scan() {
// 		line := strings.TrimSpace(scanner.Text())
// 		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
// 			inSection = strings.Contains(strings.ToLower(line), strings.ToLower(keyword))
// 			continue
// 		}
// 		if !inSection || line == "" {
// 			continue
// 		}
// 		fields := strings.Fields(line)
// 		if len(fields) == 0 {
// 			continue
// 		}
// 		key := fields[0]
// 		var row []any
// 		for _, f := range fields[1:] {
// 			row = append(row, f)
// 		}
// 		result[key] = append(result[key], row)
// 	}
// 	return result
// }

// func main() {
// 	if len(os.Args) < 3 {
// 		fmt.Println("Usage: <filepath> <section-keyword>")
// 		return
// 	}
// 	filename := os.Args[1]
// 	keyword := os.Args[2]

// 	haystack := DefineHaystack(filename, keyword)

// 	for key, rows := range haystack {
// 		fmt.Printf("%s: %v\n", key, rows)
// 	}
// }

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"sort"
// 	"strconv"
// 	"strings"
// )

// type ChargerReport struct {
// 	start uint64
// 	end   uint64
// 	up    bool
// }

// type StationUptime struct {
// 	StationID uint32
// 	Uptime    int // percentage [0-100]
// }

// type Electric struct {
// 	stations map[uint32][]uint32        // stationID -> []chargerID
// 	reports  map[uint32][]ChargerReport // chargerID -> []ChargerReport
// 	handlers map[string]func(string) error
// }

// func NewElectric() *Electric {
// 	e := &Electric{
// 		stations: make(map[uint32][]uint32),
// 		reports:  make(map[uint32][]ChargerReport),
// 	}
// 	e.handlers = map[string]func(string) error{
// 		"stations": func(line string) error {
// 			parts := strings.Fields(line)
// 			if len(parts) < 2 {
// 				return fmt.Errorf("invalid stations line: %s", line)
// 			}
// 			stationID, err := strconv.ParseUint(parts[0], 10, 32)
// 			if err != nil {
// 				return err
// 			}
// 			var chargers []uint32
// 			for _, cid := range parts[1:] {
// 				chid, err := strconv.ParseUint(cid, 10, 32)
// 				if err != nil {
// 					return err
// 				}
// 				chargers = append(chargers, uint32(chid))
// 			}
// 			e.AddStation(uint32(stationID), chargers)
// 			return nil
// 		},
// 		"charger availability reports": func(line string) error {
// 			parts := strings.Fields(line)
// 			if len(parts) != 4 {
// 				return fmt.Errorf("invalid report line: %s", line)
// 			}
// 			chid, err := strconv.ParseUint(parts[0], 10, 32)
// 			if err != nil {
// 				return err
// 			}
// 			start, err := strconv.ParseUint(parts[1], 10, 64)
// 			if err != nil {
// 				return err
// 			}
// 			end, err := strconv.ParseUint(parts[2], 10, 64)
// 			if err != nil {
// 				return err
// 			}
// 			if end < start {
// 				return fmt.Errorf("end before start")
// 			}
// 			up := false
// 			if parts[3] == "true" {
// 				up = true
// 			} else if parts[3] == "false" {
// 				up = false
// 			} else {
// 				return fmt.Errorf("invalid up value: %s", parts[3])
// 			}
// 			e.AddChargerReport(uint32(chid), start, end, up)
// 			return nil
// 		},
// 	}
// 	return e
// }

// // ParseInput parses the input file and populates the Electric struct.
// func (e *Electric) ParseInput(path string) error {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	section := ""
// 	for scanner.Scan() {
// 		line := strings.TrimSpace(scanner.Text())
// 		if line == "" {
// 			continue
// 		}
// 		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
// 			section = strings.ToLower(strings.Trim(line, "[]"))
// 			continue
// 		}
// 		if handler, ok := e.handlers[section]; ok {
// 			if err := handler(line); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	if err := scanner.Err(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (e *Electric) AddStation(stationID uint32, chargerIDs []uint32) {
// 	e.stations[stationID] = chargerIDs
// }

// func (e *Electric) AddChargerReport(chargerID uint32, start, end uint64, up bool) {
// 	e.reports[chargerID] = append(e.reports[chargerID], ChargerReport{start, end, up})
// }

// func (e *Electric) StationUptimes() []StationUptime {
// 	stationIDs := make([]uint32, 0, len(e.stations))
// 	for sid := range e.stations {
// 		stationIDs = append(stationIDs, sid)
// 	}
// 	sort.Slice(stationIDs, func(i, j int) bool { return stationIDs[i] < stationIDs[j] })

// 	var results []StationUptime
// 	for _, sid := range stationIDs {
// 		chargerIDs := e.stations[sid]
// 		var reportingIntervals [][2]uint64
// 		var upIntervals [][2]uint64
// 		for _, cid := range chargerIDs {
// 			for _, rep := range e.reports[cid] {
// 				reportingIntervals = append(reportingIntervals, [2]uint64{rep.start, rep.end})
// 				if rep.up {
// 					upIntervals = append(upIntervals, [2]uint64{rep.start, rep.end})
// 				}
// 			}
// 		}
// 		reportingIntervals = mergeIntervalsSimple(reportingIntervals)
// 		upIntervals = mergeIntervalsSimple(upIntervals)
// 		totalReporting := sumIntervalsSimple(reportingIntervals)
// 		totalUp := sumIntervalsSimple(upIntervals)
// 		percent := 0
// 		if totalReporting > 0 {
// 			percent = int((totalUp * 100) / totalReporting)
// 		}
// 		results = append(results, StationUptime{StationID: sid, Uptime: percent})
// 	}
// 	return results
// }

// // mergeIntervalsSimple merges overlapping intervals represented as [2]uint64 slices.
// func mergeIntervalsSimple(intervals [][2]uint64) [][2]uint64 {
// 	if len(intervals) == 0 {
// 		return nil
// 	}
// 	sort.Slice(intervals, func(i, j int) bool {
// 		return intervals[i][0] < intervals[j][0]
// 	})
// 	merged := make([][2]uint64, 0, len(intervals))
// 	merged = append(merged, intervals[0])
// 	for _, curr := range intervals[1:] {
// 		last := &merged[len(merged)-1]
// 		if curr[0] <= last[1] {
// 			if curr[1] > last[1] {
// 				last[1] = curr[1]
// 			}
// 		} else {
// 			merged = append(merged, curr)
// 		}
// 	}
// 	return merged
// }

// func sumIntervalsSimple(intervals [][2]uint64) uint64 {
// 	var sum uint64
// 	for _, iv := range intervals {
// 		if iv[1] > iv[0] {
// 			sum += iv[1] - iv[0]
// 		}
// 	}
// 	return sum
// }

// func main() {
// 	if len(os.Args) != 2 {
// 		fmt.Println("ERROR")
// 		return
// 	}
// 	e := NewElectric()
// 	if err := e.ParseInput(os.Args[1]); err != nil {
// 		fmt.Println("ERROR")
// 		return
// 	}
// 	for _, uptime := range e.StationUptimes() {
// 		fmt.Printf("%d %d\n", uptime.StationID, uptime.Uptime)
// 	}
// }
