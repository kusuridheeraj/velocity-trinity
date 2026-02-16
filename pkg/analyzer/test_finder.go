package analyzer

import (
	"path/filepath"
	"strings"
)

// FindTestFiles returns a list of test files that *might* be affected by the changes in `files`
func FindTestFiles(files []string) ([]string, error) {
	var tests []string
	
	for _, file := range files {
		// Heuristic 1: If file is itself a test, add it
		if isTestFile(file) {
			tests = append(tests, file)
			continue
		}

		// Heuristic 2: Check for co-located test files
		// file: src/foo.ts
		// test: src/foo.test.ts or src/foo.spec.ts
		ext := filepath.Ext(file)
		base := strings.TrimSuffix(file, ext)
		
		candidates := []string{
			base + ".test" + ext,
			base + ".spec" + ext,
			base + "_test" + ext,
		}

		for _, candidate := range candidates {
			// In a real system, use os.Stat to check if file exists
			// For MVP, we'll assume it exists if the naming matches our convention
			// or just return the candidate and let the test runner handle "file not found"
			tests = append(tests, candidate)
		}
	}

	return unique(tests), nil
}

func isTestFile(path string) bool {
	return strings.Contains(path, ".test.") || 
		strings.Contains(path, ".spec.") || 
		strings.Contains(path, "_test.py")
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
