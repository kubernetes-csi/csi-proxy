package generators

import (
	"io"
	"sort"
	"strings"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// a typesGenerator generates types.go files - one per API group; only if it doesn't already exist,
// and the API group only has one version.
// This is simply meant to help bootstrapping new API groups.
type typesGenerator struct {
	generator.DefaultGen

	groupDefinition *groupDefinition
	version         *apiVersion

	importTracker namer.ImportTracker
}

type namedType struct {
	name string
	*types.Type
}

func (g *typesGenerator) Namers(*generator.Context) namer.NameSystems {
	g.importTracker = generator.NewImportTracker()

	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.version.Path, g.importTracker),
	}
}

func (g *typesGenerator) Filter(*generator.Context, *types.Type) bool {
	return false
}

func (g *typesGenerator) Imports(context *generator.Context) (imports []string) {
	return g.importTracker.ImportLines()
}

func (g *typesGenerator) Init(context *generator.Context, writer io.Writer) error {
	snippetWriter := generator.NewSnippetWriter(writer, context, "$", "$")

	protoPkgPath := g.groupDefinition.versionedAPIPkg(g.version.Name)
	protoPkg := context.Universe[protoPkgPath]
	if protoPkg == nil {
		// shouldn't happen
		klog.Fatalf("proto package %s for API group %s version %s not loaded",
			protoPkgPath, g.groupDefinition.name, g.version.Name)
	}

	protoMessages := make([]namedType, 0)

	for name, t := range protoPkg.Types {
		if isProtobufMessage(t) {
			protoMessages = append(protoMessages, namedType{
				name: name,
				Type: t,
			})
		}
	}

	// re-order alphabetically to ensure deterministic builds
	sort.Slice(protoMessages, func(i, j int) bool {
		return protoMessages[i].name < protoMessages[j].name
	})

	for _, protoMessage := range protoMessages {
		g.generateStruct(protoMessage.name, protoMessage.Type, snippetWriter)
	}

	return snippetWriter.Error()
}

// isProtobufMessage returns true iff t is a protobuf message; it determines that
// by looking for a ProtoMessage method with no parameters and no results.
func isProtobufMessage(t *types.Type) bool {
	if t == nil || t.Methods == nil {
		return false
	}
	protoMsgFunc, present := t.Methods["ProtoMessage"]
	if !present || protoMsgFunc.Signature == nil {
		return false
	}

	sig := protoMsgFunc.Signature

	return len(sig.Parameters) == 0 && len(sig.Results) == 0 && !sig.Variadic
}

func (g *typesGenerator) generateStruct(typeName string, t *types.Type, snippetWriter *generator.SnippetWriter) {
	snippetWriter.Do("type "+typeName+" struct {\n", nil)

	for _, member := range t.Members {
		if isInternalProtobufField(&member) {
			// internal protobuf field
			continue
		}

		for _, commentLine := range member.CommentLines {
			commentLine = strings.TrimSpace(commentLine)
			if commentLine == "" {
				continue
			}
			snippetWriter.Do("// $.$\n", commentLine)
		}
		snippetWriter.Do(member.Name+" $.|raw$\n", member.Type)
	}

	snippetWriter.Do("}\n\n", nil)
}
