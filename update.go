package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"
)

func main() {

	existsDir := make(map[string]bool)

	existFiles := make(map[string]bool)
	scanDirs := []string{"."}
	for len(scanDirs) > 0 {
		dir := scanDirs[0]
		scanDirs = scanDirs[1:]
		files, err := os.ReadDir(dir)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if file.IsDir() {
				if file.Name()[0] == '.' {
					continue
				}
				scanDirs = append(scanDirs, filepath.Join(dir, file.Name()))
				existsDir[strings.ToLower(filepath.Dir(filepath.Join(dir, file.Name())))] = true
			} else {
				existFiles[strings.ToLower(filepath.Join(dir, file.Name()))] = true
			}
		}
	}

	fmt.Println("Downloading supported_devices.csv ...")
	resp, err := http.Get("https://storage.googleapis.com/play_public/supported_devices.csv")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	bodyU16 := make([]uint16, len(body)/2)
	for i := 0; i < len(body); i += 2 {
		bodyU16[i/2] = uint16(body[i]) | uint16(body[i+1])<<8
	}

	bodyStr := string(utf16.Decode([]uint16(bodyU16)))

	lines := strings.Split(bodyStr, "\n")
	lines = lines[1:]
	for _, line := range lines {
		if line == "" {
			break
		}
		flds := splitCsvLine(line)
		if len(flds) != 4 {
			continue
		}
		brand := strings.TrimSpace(flds[0])
		if brand == "" {
			continue
		}

		name := strings.TrimSpace(flds[1])
		device := strings.ReplaceAll(strings.TrimSpace(flds[2]), "/", " - ")
		model := strings.ReplaceAll(strings.TrimSpace(flds[3]), "/", "- ")

		firstBrandChar := brand[0]
		if len(device) > 0 {
			deviceFile := strings.ToLower(fmt.Sprintf("%c/%s/%s", firstBrandChar, brand, device))
			deviceDir := filepath.Dir(deviceFile)
			if !existsDir[deviceDir] {
				os.MkdirAll(deviceDir, 0755)
				existsDir[deviceDir] = true
			}
			if !existFiles[deviceFile] {
				os.WriteFile(deviceFile, []byte(name), 0644)
				fmt.Println("Created", deviceFile, "with name", name)
				existFiles[deviceFile] = true
			}
		}
		if len(model) > 0 {
			modelFile := strings.ToLower(fmt.Sprintf("%c/%s/%s", firstBrandChar, brand, model))
			modelDir := filepath.Dir(modelFile)
			if !existsDir[modelDir] {
				os.MkdirAll(modelDir, 0755)
				existsDir[modelDir] = true
			}
			if !existFiles[modelFile] {
				os.WriteFile(modelFile, []byte(name), 0644)
				fmt.Println("Created", modelFile, "with name", name)
				existFiles[modelFile] = true
			}
		}
	}
}

func splitCsvLine(line string) []string {
	var (
		flds []string
		fld  string
		inQ  bool
	)
	for _, r := range line {
		switch {
		case r == ',' && !inQ:
			flds = append(flds, fld)
			fld = ""
		case r == '"' && !inQ:
			inQ = true
		case r == '"' && inQ:
			inQ = false
		default:
			fld += string(r)
		}
	}
	flds = append(flds, fld)
	return flds
}
