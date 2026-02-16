package analyzer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/velocity-trinity/core/pkg/analyzer/languages"
)

// LanguageParser is a generic interface for language parsers
type LanguageParser interface {
	Parse(filePath string) ([]string, error)
}

// GetParser returns the appropriate parser for the file extension
func GetParser(filePath string) (LanguageParser, error) {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".ts", ".tsx", ".js", ".jsx":
		return &languages.TypeScriptParser{}, nil
	case ".py":
		return &languages.PythonParser{}, nil
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

// BuildDependencyGraph walks a directory and builds a dependency map
// Map format: File -> List of Files that import it
func BuildDependencyGraph(root string) (map[string][]string, error) {
	// Not implemented yet - this will be the core of Phase 2
	// For now, let's just return a placeholder
	return nil, nil
}

// AnalyzeFile returns the direct dependencies of a single file
func AnalyzeFile(filePath string) ([]string, error) {
	parser, err := GetParser(filePath)
	if err != nil {
		return nil, err
	}
	
	deps, err := parser.Parse(filePath)
	if err != nil {
		return nil, err
	}

	// Filter for relative imports only (start with ./ or ../)
	// In a real implementation, we would resolve absolute paths based on tsconfig/python path
	var relativeDeps []string
	for _, dep := range deps {
		if strings.HasPrefix(dep, ".") {
			relativeDeps = append(relativeDeps, dep)
		}
	}
	return relativeDeps, nil
}
