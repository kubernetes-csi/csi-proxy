package generators

import (
	"io"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

// a typesGeneratedGenerator generates types_generated.go files - one per API group.
type typesGeneratedGenerator struct {
	generator.DefaultGen
	groupDefinition *groupDefinition
}

func (g *typesGeneratedGenerator) Filter(*generator.Context, *types.Type) bool {
	return false
}

func (g *typesGeneratedGenerator) Imports(*generator.Context) []string {
	return []string{
		"context",
		"google.golang.org/grpc",
		"github.com/kubernetes-csi/csi-proxy/client/apiversion",
	}
}

func (g *typesGeneratedGenerator) Init(context *generator.Context, writer io.Writer) error {
	snippetWriter := generator.NewSnippetWriter(writer, context, "$", "$")

	snippetWriter.Do(`type VersionedAPI interface {
Register(grpcServer *grpc.Server)
}

// All the functions this group's server needs to define.
type ServerInterface interface {
`, nil)

	for _, namedCallback := range g.groupDefinition.serverCallbacks {
		callback := replaceTypesPackage(namedCallback.callback, pkgPlaceholder, "")

		snippetWriter.Do(namedCallback.name+"(", nil)
		for _, param := range callback.Signature.Parameters {
			snippetWriter.Do("$.$, ", param)
		}
		// add the version parameter
		snippetWriter.Do("apiversion.Version) (", nil)
		for _, returnValue := range callback.Signature.Results {
			snippetWriter.Do("$.$, ", returnValue)
		}
		snippetWriter.Do(")\n", nil)
	}
	snippetWriter.Do("}\n", nil)

	return snippetWriter.Error()
}
