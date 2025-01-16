package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func mergeGoFiles(files ...string) (string, error) {
	all := ""
	imports := []string{}
	importRegex := regexp.MustCompile(`(?m)^import\s+(?:"[^"]+"|[\(\s][^)]*[\)\s])`)

	fmt.Printf("files: %v\n", files)
	for _, file := range files {
		f, err := os.ReadFile(file)
		if err != nil {
			return "", err
		}

		content := string(f)

		// Extract import statements and append to the imports slice
		matches := importRegex.FindAllString(content, -1)
		imports = append(imports, matches...)

		// Remove package and import statements from the content
		content = importRegex.ReplaceAllString(content, "")
		r := regexp.MustCompile(`(?m)^package\s+.*`)
		content = r.ReplaceAllString(content, "")

		all += content
	}

	// Deduplicate imports
	imports = deduplicateImports(imports)

	// Merge all imports at the top
	importSection := strings.Join(imports, "\n")
	finalOutput := importSection + "\n\n" + all

	return finalOutput, nil
}

func deduplicateImports(imports []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, imp := range imports {
		if !seen[imp] {
			seen[imp] = true
			result = append(result, imp)
		}
	}

	return result
}

func getFilesRec(files ...string) ([]string, error) {
	var result []string

	for _, file := range files {
		s, err := os.Stat(file)
		if err != nil {
			return nil, fmt.Errorf("failed to stat file %q: %v", file, err)
		}

		if !s.IsDir() {
			// Add the file if it is a Go file
			if strings.ToLower(filepath.Ext(file)) == ".go" {
				result = append(result, file)
			}
			continue
		}

		// Walk through the directory
		err = filepath.Walk(file, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error accessing path %q: %v", path, err)
			}
			// Add files with .go extension
			if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".go" {
				result = append(result, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
