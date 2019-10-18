package generators

import (
	"strings"

	goflag "flag"

	"github.com/spf13/pflag"
	"k8s.io/gengo/args"
	"k8s.io/klog"
)

// Execute runs csi-proxy-api-gen. It's exposed as a public function
// to be able to easily run it from integration tests.
func Execute(executableName string, cliArgs ...string) {
	if err := buildArgs(executableName, cliArgs).Execute(
		nameSystems(),
		defaultNameSystem(),
		packages,
	); err != nil {
		klog.Fatalf("Error: %v", err)
	}

	klog.Infof("Generation successful!")
}

func buildArgs(executableName string, cliArgs []string) *args.GeneratorArgs {
	goFlagSet := goflag.NewFlagSet(executableName, goflag.ExitOnError)
	klog.InitFlags(goFlagSet)

	genericArgs := args.Default().WithoutDefaultFlagParsing()

	pflagFlagSet := pflag.NewFlagSet(executableName, pflag.ExitOnError)
	genericArgs.AddFlags(pflagFlagSet)
	pflagFlagSet.AddGoFlagSet(goFlagSet)
	if err := pflagFlagSet.Parse(cliArgs); err != nil {
		klog.Fatalf("Unable to parse CLI args: %v", err)
	}

	klog.Infof("Verbosity level set to %d", verbosityLevel())

	// if no package argument, default to processing canonical API groups, under csiProxyAPIPath
	if len(genericArgs.InputDirs) == 0 {
		genericArgs.InputDirs = append(genericArgs.InputDirs, csiProxyAPIPath)
	}

	// it doesn't really make sense to consider a package in isolation, since an API group is
	// always a collection of subpackages (its versions)
	// so we consider all inputs recursively
	for i, inputDir := range genericArgs.InputDirs {
		if !strings.HasSuffix(inputDir, "...") {
			genericArgs.InputDirs[i] = canonicalizePkgPath(inputDir) + "/..."
		}
	}

	return genericArgs
}

// verbosityLevel returns the current verbosity level
func verbosityLevel() klog.Level {
	level := klog.Level(0)
	for klog.V(level) {
		level++
	}
	return level - 1
}
