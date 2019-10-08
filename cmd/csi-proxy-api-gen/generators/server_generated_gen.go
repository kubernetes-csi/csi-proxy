package generators

import (
	"io"
	"strings"

	"github.com/iancoleman/strcase"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

// a serverGeneratedGenerator generates server_generated.go files - one per API version.
type serverGeneratedGenerator struct {
	generator.DefaultGen
	groupDefinition *groupDefinition
	version         *apiVersion
}

func (g *serverGeneratedGenerator) Namers(*generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"shortenVersionPackage": &shortenVersionPackageNamer{
			version: g.version,
		},
		"versionedVariable": &versionedVariableNamer{
			version: g.version,
		},
	}
}

func (g *serverGeneratedGenerator) Filter(*generator.Context, *types.Type) bool {
	return false
}

func (g *serverGeneratedGenerator) Imports(*generator.Context) []string {
	return []string{
		"context",
		"google.golang.org/grpc",
		"github.com/kubernetes-csi/csi-proxy/client/apiversion",
		g.groupDefinition.internalServerPkg(),
		g.groupDefinition.versionedAPIPkg(g.version.Name),
	}
}

func (g *serverGeneratedGenerator) Init(context *generator.Context, writer io.Writer) error {
	snippetWriter := generator.NewSnippetWriter(writer, context, "$", "$")

	snippetWriter.Do(`var version = apiversion.NewVersionOrPanic("$.version$")

type versionedAPI struct {
	apiGroupServer internal.ServerInterface
}

func NewVersionedServer(apiGroupServer internal.ServerInterface) internal.VersionedAPI {
	return &versionedAPI{
		apiGroupServer: apiGroupServer,
	}
}

func (s *versionedAPI) Register(grpcServer *grpc.Server) {
	$.version$.Register$.camelGroupName$Server(grpcServer, s)
}

	`, map[string]string{
		"camelGroupName": strcase.ToCamel(g.groupDefinition.name),
		"version":        g.version.Name,
	})

	// write a request handler for each server callback
	for _, namedCallback := range g.version.serverCallbacks {
		g.writeWrapperFunction(namedCallback.name, namedCallback.callback, snippetWriter)
	}

	return snippetWriter.Error()
}

func (g *serverGeneratedGenerator) writeWrapperFunction(callbackName string, callback *types.Type, snippetWriter *generator.SnippetWriter) {
	// write the func signature
	snippetWriter.Do("func (s *versionedAPI) $.$(", callbackName)
	for _, param := range callback.Signature.Parameters {
		snippetWriter.Do("$.|versionedVariable$ $.|shortenVersionPackage$, ", param)
	}
	snippetWriter.Do(") (", nil)
	for _, returnValue := range callback.Signature.Results {
		snippetWriter.Do("$.|shortenVersionPackage$, ", returnValue)
	}
	snippetWriter.Do(") {\n", nil)

	// when returning errors from conversion
	returnErrLine := "return " + strings.Repeat("nil, ", len(callback.Signature.Results)-1) + "err"

	// then convert all versioned arguments to internal structs
	for _, param := range callback.Signature.Parameters {
		if !isVersionedVariable(param, g.version) {
			continue
		}
		snippetWriter.Do("$.|short$ := &internal.$.|removePackage${}\n", param)
		snippetWriter.Do("if err := Convert_"+g.version.Name+"_$.|removePackage$_To_internal_$.|removePackage$($.|versionedVariable$, $.|short$); err != nil {\n", param)
		snippetWriter.Do(returnErrLine+"\n}\n", nil)
	}
	snippetWriter.Do("\n", nil)

	// call the internal server
	for i, returnValue := range callback.Signature.Results {
		if i != 0 {
			snippetWriter.Do(", ", nil)
		}
		snippetWriter.Do("$.|short$", returnValue)
	}
	snippetWriter.Do(" := s.apiGroupServer."+callbackName+"(", nil)
	for _, param := range callback.Signature.Parameters {
		snippetWriter.Do("$.|short$, ", param)
	}
	snippetWriter.Do("version)\nif err != nil {\n"+returnErrLine+"\n}\n\n", nil)

	// convert all internal return values to versioned structs
	for _, returnValue := range callback.Signature.Results {
		if !isVersionedVariable(returnValue, g.version) {
			continue
		}
		snippetWriter.Do("$.|versionedVariable$ := &"+g.version.Name+".$.|removePackage${}\n", returnValue)
		snippetWriter.Do("if err := Convert_internal_$.|removePackage$_To_"+g.version.Name+"_$.|removePackage$($.|short$, $.|versionedVariable$); err != nil {\n", returnValue)
		snippetWriter.Do(returnErrLine+"\n}\n", nil)
	}
	snippetWriter.Do("\n", nil)

	// return values
	snippetWriter.Do("return ", nil)
	for i, returnValue := range callback.Signature.Results {
		if i != 0 {
			snippetWriter.Do(", ", nil)
		}
		snippetWriter.Do("$.|versionedVariable$", returnValue)
	}

	// end of the request handler
	snippetWriter.Do("\n}\n\n", nil)
}
