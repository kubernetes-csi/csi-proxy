package generators

import (
	"strings"

	"k8s.io/gengo/types"
)

// canonicalizePkgPath ensures package paths are consistent.
func canonicalizePkgPath(pkgPath string) string {
	return strings.TrimSuffix(pkgPath, "/")
}

// snakeCaseToPackageName turns a snake case string into a go package name.
func snakeCaseToPackageName(name string) string {
	return strings.ReplaceAll(name, "_", "")
}

// replaceTypesPackage return a new type, equal to t except moved from package
// pkg to package newPkg (and same for other types referenced by t).
// t itself remains unchanged.
func replaceTypesPackage(t *types.Type, pkg, newPkg string) *types.Type {
	return recursiveReplaceTypesPackage(t, normalizePkg(pkg), normalizePkg(newPkg), make(map[*types.Type]*types.Type))
}

func normalizePkg(pkg string) string {
	if pkg != "" && !strings.HasSuffix(pkg, ".") {
		pkg += "."
	}
	return pkg
}

func recursiveReplaceTypesPackage(t *types.Type, pkg, newPkg string, visited map[*types.Type]*types.Type) *types.Type {
	if t == nil {
		return nil
	}
	if result, present := visited[t]; present {
		return result
	}

	result := &types.Type{
		Name: types.Name{
			Name:    strings.ReplaceAll(t.Name.Name, pkg, newPkg),
			Package: strings.ReplaceAll(t.Name.Package, pkg, newPkg),
		},
		Kind:                      t.Kind,
		CommentLines:              t.CommentLines,
		SecondClosestCommentLines: t.SecondClosestCommentLines,
	}
	visited[t] = result

	if len(t.Members) != 0 {
		members := make([]types.Member, len(t.Members))
		for i, member := range t.Members {
			members[i] = types.Member{
				Name:         member.Name,
				Embedded:     member.Embedded,
				CommentLines: member.CommentLines,
				Tags:         member.Tags,
				Type:         recursiveReplaceTypesPackage(member.Type, pkg, newPkg, visited),
			}
		}
		result.Members = members
	}

	result.Elem = recursiveReplaceTypesPackage(t.Elem, pkg, newPkg, visited)
	result.Key = recursiveReplaceTypesPackage(t.Key, pkg, newPkg, visited)
	result.Underlying = recursiveReplaceTypesPackage(t.Underlying, pkg, newPkg, visited)

	if len(t.Methods) != 0 {
		methods := make(map[string]*types.Type)
		for k, v := range t.Methods {
			methods[k] = recursiveReplaceTypesPackage(v, pkg, newPkg, visited)
		}
		result.Methods = methods
	}

	var signature *types.Signature
	if t.Signature != nil {
		signature = &types.Signature{
			Receiver:     recursiveReplaceTypesPackage(t.Signature.Receiver, pkg, newPkg, visited),
			Parameters:   replaceTypesSlicePackage(t.Signature.Parameters, pkg, newPkg, visited),
			Results:      replaceTypesSlicePackage(t.Signature.Results, pkg, newPkg, visited),
			Variadic:     t.Signature.Variadic,
			CommentLines: t.Signature.CommentLines,
		}
		result.Signature = signature
	}

	return result
}

func replaceTypesSlicePackage(ts []*types.Type, pkg, newPkg string, visited map[*types.Type]*types.Type) []*types.Type {
	result := make([]*types.Type, len(ts))
	for i, t := range ts {
		result[i] = recursiveReplaceTypesPackage(t, pkg, newPkg, visited)
	}
	return result
}

// isInternalProtobufField returns true iff the given member is internal
// to protobuf, and can be safely ignored by any other piece of code.
func isInternalProtobufField(member *types.Member) bool {
	return strings.HasPrefix(member.Name, "XXX_")
}
