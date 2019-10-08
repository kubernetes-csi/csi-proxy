package generators

import (
	"io"
	"sort"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
)

// a apiGroupGeneratedGenerator generates api_group_generated.go files - one per API group.
type apiGroupGeneratedGenerator struct {
	generator.DefaultGen
	groupDefinition *groupDefinition
}

func (g *apiGroupGeneratedGenerator) Filter(*generator.Context, *types.Type) bool {
	return false
}

func (g *apiGroupGeneratedGenerator) Imports(*generator.Context) []string {
	imports := []string{
		"github.com/kubernetes-csi/csi-proxy/client/apiversion",
		"github.com/kubernetes-csi/csi-proxy/internal/server",
		g.groupDefinition.internalServerPkg(),
	}

	for _, version := range g.groupDefinition.versions {
		imports = append(imports, g.groupDefinition.versionedServerPkg(version.Name))
	}

	return imports
}

func (g *apiGroupGeneratedGenerator) Init(context *generator.Context, writer io.Writer) error {
	snippetWriter := generator.NewSnippetWriter(writer, context, "$", "$")

	snippetWriter.Do(`const name = "$.$"`, g.groupDefinition.name)

	snippetWriter.Do(`

// ensure the server defines all the required methods
var _ internal.ServerInterface = &Server{}

func (s *Server) VersionedAPIs() []*server.VersionedAPI {
`, nil)

	versions := make([]apiversion.Version, len(g.groupDefinition.versions))
	for i, vsn := range g.groupDefinition.versions {
		versions[i] = apiversion.NewVersionOrPanic(vsn.Name)
	}
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Compare(versions[j]) == apiversion.Lesser
	})

	for _, version := range versions {
		snippetWriter.Do("$.$Server := $.$.NewVersionedServer(s)\n", version)
	}

	snippetWriter.Do("\n\nreturn []*server.VersionedAPI{\n", nil)
	for _, version := range versions {
		snippetWriter.Do(`{
				Group:      name,
				Version:    apiversion.NewVersionOrPanic("$.$"),
				Registrant: $.$Server.Register,
			},
			`, version.String())
	}
	snippetWriter.Do("\n}\n}\n", nil)

	return snippetWriter.Error()
}
