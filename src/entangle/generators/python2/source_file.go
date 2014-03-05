package python2

import (
	"fmt"
	"bytes"
	"entangle/utils"
	"io"
	"sort"
)

// Python 2 source file.
type SourceFile struct {
	// Module name.
	moduleName string

	// Module imports.
	moduleImports utils.StringSet

	// Type imports.
	typeImports map[string]map[string]utils.StringSet

	// Blocks.
	blocks [][]byte

	// Exports.
	exports utils.StringSet
}

// New source file.
func NewSourceFile(moduleName string) *SourceFile {
	return &SourceFile {
		moduleName: moduleName,
		moduleImports: make(utils.StringSet),
		typeImports: make(map[string]map[string]utils.StringSet),
		blocks: make([][]byte, 0, 16),
		exports: make(utils.StringSet),
	}
}

// Make sure that a type or method is imported.
//
// Multiple calls for the same type has no implications.
func (s *SourceFile) Import(module, name string) {
	s.ImportAs(module, name, name)
}

// Make sure that a type or method is imported with an alias.
//
// Multiple calls for the same type has no implications.
func (s *SourceFile) ImportAs(module, name, as string) {
	if module == s.moduleName || module == fmt.Sprintf(".%s", s.moduleName) {
		return
	}

	if _, ok := s.typeImports[module]; !ok {
		s.typeImports[module] = make(map[string]utils.StringSet)
	}
	if _, ok := s.typeImports[module][name]; !ok {
		s.typeImports[module][name] = make(utils.StringSet)
	}

	s.typeImports[module][name].Add(as)
}

// Make sure that a module is imported.
func (s *SourceFile) ImportModule(name string) {
	if name == s.moduleName {
		return
	} else if name[0] == '.' {
		panic("You know that relative module imports are not valid")
	}

	s.moduleImports.Add(name)
}

// Export a type.
func (s *SourceFile) Export(name string) {
	s.exports.Add(name)
}

// Add a code block.
func (s *SourceFile) AddBlock(block []byte) {
	block = bytes.TrimSpace(block)
	if len(block) == 0 {
		return
	}
	block = append(block, '\n')

	s.blocks = append(s.blocks, block)
}

// Generate the source file.
func (s *SourceFile) Generate(w io.Writer) (err error) {
	// Generate imports.
	importStatements := make([]generatorImportStatement, 0, len(s.moduleImports) + len(s.typeImports))

	for name, _ := range s.moduleImports {
		importStatements = append(importStatements, generatorImportStatement {
			Name: name,
			Relativity: moduleRelativity(name),
		})
	}

	for name, types := range s.typeImports {
		importStatements = append(importStatements, generatorImportStatement {
			Name: name,
			Types: types,
			Relativity: moduleRelativity(name),
		})
	}

	sort.Sort(generatorImportStatementSorter(importStatements))

	if len(importStatements) > 0 {
		for _, stmt := range importStatements {
			if _, err = w.Write([]byte(stmt.Generate() + "\n")); err != nil {
				return
			}
		}
	}

	// Write blocks.
	for i, block := range s.blocks {
		if len(block) == 0 {
			continue
		}

		if i > 0 || len(importStatements) > 0 {
			if _, err = w.Write([]byte { '\n', '\n' }); err != nil {
				return
			}
		}

		if _, err = w.Write(block); err != nil {
			return
		}
	}

	// Write exports.
	exportsW := newCodeWriter()

	if len(importStatements) > 0 || len(s.blocks) > 0 {
		exportsW.BlankLine()
		exportsW.BlankLine()
	}

	exportsW.ParentherizedDefinition("__all__", stringifyStrings(s.exports.Sorted())...)
	_, err = w.Write(exportsW.Bytes())

	return
}
