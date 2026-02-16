package languages

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// PythonParser handles .py files
type PythonParser struct{}

func (p *PythonParser) Parse(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var dependencies []string
	scanner := bufio.NewScanner(file)

	// Regex for: import module
	// Regex for: from module import ...
	importRegex := regexp.MustCompile(`^\s*import\s+([\w.]+)`)
	fromRegex := regexp.MustCompile(`^\s*from\s+([\w.]+)\s+import`)

	for scanner.Scan() {
		line := scanner.Text()
		// Skip comments
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		matches := importRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			dependencies = append(dependencies, matches[1])
		}

		matches = fromRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			dependencies = append(dependencies, matches[1])
		}
	}

	return dependencies, scanner.Err()
}
