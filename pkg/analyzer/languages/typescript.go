package languages

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// Parser interface for different languages
type Parser interface {
	Parse(filePath string) ([]string, error)
}

// TypeScriptParser handles .ts, .tsx, .js, .jsx files
type TypeScriptParser struct{}

func (p *TypeScriptParser) Parse(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var dependencies []string
	scanner := bufio.NewScanner(file)

	// Regex for: import ... from '...'
	// Regex for: require('...')
	// This is a simplified regex and might miss edge cases (e.g. multiline imports)
	importRegex := regexp.MustCompile(`(?:import|export)\s+.*\s+from\s+['"]([^'"]+)['"]`)
	requireRegex := regexp.MustCompile(`require\(['"]([^'"]+)['"]\)`)
	dynamicImportRegex := regexp.MustCompile(`import\(['"]([^'"]+)['"]\)`)

	for scanner.Scan() {
		line := scanner.Text()
		// Skip comments (very basic check)
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}

		matches := importRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			dependencies = append(dependencies, matches[1])
		}

		matches = requireRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			dependencies = append(dependencies, matches[1])
		}
		
		matches = dynamicImportRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			dependencies = append(dependencies, matches[1])
		}
	}

	return dependencies, scanner.Err()
}
