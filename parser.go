package inparser

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type IniFile struct {
	Sections map[string]*Section
}

type Section struct {
	Name string
	Properties []*Property
}

type Property struct {
	Key string
	Value string
	Mapping map[string]string
}

func NewIni() *IniFile {
	mapping := make(map[string]*Section)
	return &IniFile{Sections:mapping}
}

func (ini *IniFile) GetSection(name string) *Section {
	val, ok := ini.Sections[name]
	if ok {
		return val
	}
	sect := &Section{Name:name}
	ini.Sections[name] = sect
	return sect
}

func commentOrEmpty(line string) bool {
	if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#"){
		return true
	}
	if line == "" {
		return true
	}
	return false
}

func isReadSection(line string) bool {
	if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
		return true
	}
	return false
}

func complexParse(line string, indexes [][]int) *Property {
	// Rewrite to cover string with space edge cases
	prop := &Property{}
	prop.Mapping = make(map[string]string)
	prop.Key = line[0:indexes[0][0]]
	newLine := strings.TrimSpace(line[indexes[0][1]:])
	keyVals := strings.Split(newLine, " ")
	for _, val := range keyVals {
		kv := strings.Split(val, "=")
		if len(kv) > 1 {
			prop.Mapping[kv[0]] = kv[1]
		}
	}
	return prop
}

func simpleParse(line string) *Property {
	values := strings.Split(line, "=")
	if len(values) > 1 {
		return &Property{Key: strings.TrimSpace(values[0]), Value: strings.TrimSpace(values[1])}
	}
	return nil
}

func parseProperty(line string) *Property {
	assignR := regexp.MustCompile("=")
	indexes := assignR.FindAllStringIndex(line, 2)
	if len(indexes) > 1 {
		return complexParse(line, indexes)
	} else {
		return simpleParse(line)
	}
}

func ParseIni(filename string) (*IniFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	iniFile := NewIni()
	scanner := bufio.NewScanner(file)
	name := ""
	for scanner.Scan() {
		text := scanner.Text()
		trimmed := strings.TrimSpace(text)
		if commentOrEmpty(trimmed) {
			continue
		}
		if isReadSection(trimmed) {
			name = trimmed[1:len(trimmed)-1]
		} else {
			res := parseProperty(trimmed)
			sect := iniFile.GetSection(name)
			sect.Properties = append(sect.Properties, res)
		}
	}
	return iniFile, nil
}