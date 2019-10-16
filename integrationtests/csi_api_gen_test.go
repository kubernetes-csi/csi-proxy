package integrationtests

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/cmd/csi-proxy-api-gen/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This tests API generator; more specifically, its main goal is to ensure
// that the API generator works as expected when creating a new API group.
// On top of this, the regular build checks that all checked-in generated files
// are up-to-date (i.e. consistent with the current generator).

// TestNewAPIGroup tests that bootstraping a new group works as intended.
func TestNewAPIGroup(t *testing.T) {
	// clean slate
	require.Nil(t, os.RemoveAll("csiapigen/new_group/actual_output"))

	logLevel := "3"
	stdout, _ := runGenerator(t, "TestNewAPIGroup",
		"--input-dirs", "github.com/kubernetes-csi/csi-proxy/integrationtests/csiapigen/new_group/api",
		// might as well check that logging CLI args work as expected
		"-v", logLevel)

	assert.Contains(t, stdout, "Verbosity level set to "+logLevel)

	// now check the generated files are exactly what we expect
	// the files in expected_output have had their `.go` extension changed to `go_code` so that one
	// can still build all subpackages in one command.
	// If you need to regenerate them, removing the extension can be done in bash with:
	// ```
	// find integrationtests/csiapigen/new_group/expected_output -name '*.go' -exec mv -v {}{,_code} \;
	// ```
	// or in powershell with:
	// ```
	// Get-ChildItem -Path integrationtests/csiapigen/new_group/expected_output -Filter '*.go' -Recurse | % FullName | ForEach-Object { mv -Verbose $_ ${_}_code }
	// ```
	recursiveDiff(t, "csiapigen/new_group/expected_output", "csiapigen/new_group/actual_output", "_code")
}

// runGenerator runs csi-proxy-api-gen with the given CLI args, and returns stdout and stderr.
// It will also fail the test immediately if there any panics during the generation (but
// will handle those graciously).
func runGenerator(t *testing.T, testName string, cliArgs ...string) (string, string) {
	stdoutFile, err := ioutil.TempFile("", "test-csi-proxy-api-gen-stdout-"+testName)
	require.Nil(t, err)
	stderrFile, err := ioutil.TempFile("", "test-csi-proxy-api-gen-stderr-"+testName)
	require.Nil(t, err)

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	os.Stdout = stdoutFile
	os.Stderr = stderrFile

	restored := false

	restoreStdOutAndErr := func() {
		if restored {
			return
		}
		restored = true

		os.Stdout = oldStdout
		os.Stderr = oldStderr

		assert.Nil(t, stdoutFile.Close())
		assert.Nil(t, stderrFile.Close())
	}

	defer func() {
		restoreStdOutAndErr()

		panicErr := recover()
		failedErrorMsg := ""
		if panicErr != nil {
			failedErrorMsg = fmt.Sprintf("panic when generating code: %v\n", panicErr)

			readLogFile := func(logFile *os.File) string {
				contents, err := ioutil.ReadFile(logFile.Name())
				if err != nil {
					return fmt.Sprintf("<unable to read: %v>", err)
				}
				return string(contents)
			}

			failedErrorMsg += fmt.Sprintf("stdout:\n%s\n", readLogFile(stdoutFile))
			failedErrorMsg += fmt.Sprintf("stderr:\n%s\n", readLogFile(stderrFile))
		}

		assert.Nil(t, os.Remove(stdoutFile.Name()))
		assert.Nil(t, os.Remove(stderrFile.Name()))

		require.Fail(t, failedErrorMsg)
	}()

	// show time
	generators.Execute(testName, cliArgs...)

	// to flush & close the log files
	restoreStdOutAndErr()

	return readFile(t, readFile(t, stdoutFile.Name())),
		readFile(t, readFile(t, stderrFile.Name()))
}
