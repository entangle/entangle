package python2

import (
	"entangle/utils"
	"fmt"
	"math"
	"strings"
)

// Generator import statement.
type generatorImportStatement struct {
	// Module name.
	Name string

	// Types.
	//
	// If nil, this is a module import statement.
	Types map[string]utils.StringSet

	// Relativity to source file.
	Relativity uint32
}

func (s generatorImportStatement) Generate() string {
	if s.Types == nil {
		return fmt.Sprintf("import %s", s.Name)
	}

	typeStatements := make([]string, 0, len(s.Types))
	for name, aliases := range s.Types {
		for alias, _ := range aliases {
			if alias == name {
				typeStatements = append(typeStatements, name)
			} else {
				typeStatements = append(typeStatements, fmt.Sprintf("%s as %s", name, alias))
			}
		}
	}

	joined := strings.Join(typeStatements, ", ")

	if len(s.Name) + 13 + len(joined) <= 79 {
		return fmt.Sprintf("from %s import %s", s.Name, joined)
	}

	statementLines := make([]string, 0, len(typeStatements))
	for _, stmt := range typeStatements {
		if len(stmt) > 79 - 4 - 1 {
			lastSpace := strings.LastIndex(stmt, " ")
			if lastSpace != -1 {
				statementLines = append(statementLines, fmt.Sprintf("%s\n    %s", stmt[:lastSpace], stmt[lastSpace + 1:]))
				continue
			}
		}

		statementLines = append(statementLines, stmt)
	}

	return fmt.Sprintf(`from %s import (
    %s,
)`, s.Name, strings.Join(statementLines, ",\n    "))
}

// Sorted generator import statement.
type generatorImportStatementSorter []generatorImportStatement

func (l generatorImportStatementSorter) Len() int {
	return len(l)
}

func (l generatorImportStatementSorter) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l generatorImportStatementSorter) Less(i, j int) bool {
	return l[i].Relativity > l[j].Relativity || (l[i].Relativity == l[j].Relativity && l[i].Name < l[j].Name)
}

// Determine module relativity.
func moduleRelativity(name string) uint32 {
	for i, r := range name {
		if r == '.' {
			continue
		}

		if i == 0 {
			return math.MaxUint32
		}
		return uint32(i - 1)
	}

	return uint32(len(name) - 1)
}
