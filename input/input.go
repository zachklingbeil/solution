package input

import (
	"bufio"
	"bytes"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Input struct {
	Map map[string][]any
}

// TextToMap reads the Input.File, finds sections, maps the data in each section,
// and stores the result in Input.Map. It also returns the map.
func (i *Input) TextToMap(path string) error {
	lines, err := i.readTxtFileToMap(path)
	if err != nil {
		return err
	}
	sections := i.findSections(lines)
	sectionMap := i.mapSections(sections, lines)
	i.Map = sectionMap
	return nil
}

func (i *Input) readTxtFileToMap(path string) (map[int][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make(map[int][]byte)
	scanner := bufio.NewScanner(file)
	for lineNum := 1; scanner.Scan(); lineNum++ {
		lines[lineNum] = slices.Clone(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func (i *Input) findSections(lines map[int][]byte) map[string]int {
	sections := make(map[string]int)
	for lineNum, line := range lines {
		start := bytes.IndexByte(line, '[')
		end := bytes.IndexByte(line, ']')
		if start != -1 && end > start {
			section := string(line[start+1 : end])
			sections[section] = lineNum
		}
	}
	return sections
}

func (i *Input) mapSections(sections map[string]int, lines map[int][]byte) map[string][]any {
	result := make(map[string][]any)

	for section, lineNum := range sections {
		line := lines[lineNum]
		fields := strings.Fields(string(line))
		var parsedFields []any
		for _, field := range fields {
			if val, err := strconv.ParseUint(field, 10, 32); err == nil {
				parsedFields = append(parsedFields, uint32(val))
			} else {
				parsedFields = append(parsedFields, field)
			}
		}
		if len(parsedFields) > 0 {
			result[section] = parsedFields
		}
	}
	return result
}
