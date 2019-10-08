package generators

import (
	"strings"

	"github.com/iancoleman/strcase"
	"k8s.io/gengo/types"
)

// a removePackageNamer removes the package from a type's name; e.g.
// "v1.FooRequest" becomes just "FooRequest"
type removePackageNamer struct{}

func (n *removePackageNamer) Name(t *types.Type) string {
	parts := strings.Split(t.Name.String(), ".")
	return parts[len(parts)-1]
}

// a shortNamer returns suitable short variable names, derived from the type's name.
// e.g. a "FooBarRequest" type will yield "request", an "error" will be "err", etc.
type shortNamer struct{}

func (*shortNamer) Name(t *types.Type) string {
	return shortName(t)
}

// a shortenVersionPackageNamer replaces the long package name from type names
// to just the package name.
type shortenVersionPackageNamer struct {
	version *apiVersion
}

func (n *shortenVersionPackageNamer) Name(t *types.Type) string {
	return strings.ReplaceAll(t.Name.String(), n.version.Package.Path, n.version.Package.Name)
}

// a versionedVariableNamer returns suitable short variable names, derived from the type's name,
// akin to a shortNamer; it also prefixes variables from types belonging to version-specific
// packages with the string "versioned".
type versionedVariableNamer struct {
	version *apiVersion
}

func (n *versionedVariableNamer) Name(t *types.Type) string {
	varName := shortName(t)
	if isVersionedVariable(t, n.version) {
		varName = "versioned" + strcase.ToCamel(varName)
	}
	return varName
}

func shortName(t *types.Type) string {
	snake := strcase.ToSnake(t.Name.Name)
	parts := strings.Split(snake, "_")
	result := parts[len(parts)-1]
	if result == t.Name.Name {
		result = result[:3]
	}
	return result
}
