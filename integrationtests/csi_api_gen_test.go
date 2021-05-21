package integrationtests

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This tests API generator; more specifically, its main goal is to ensure
// that the API generator works as expected when creating a new API group.
// On top of this, the regular build checks that all checked-in generated files
// are up-to-date (i.e. consistent with the current generator).

// TestNewAPIGroup tests that bootstraping a new group works as intended.
func TestNewAPIGroup(t *testing.T) {
	skipTestOnCondition(t, isRunningOnGhActions())

	// clean slate
	require.Nil(t, os.RemoveAll("csiapigen/new_group/actual_output"))

	// check that the csi-proxy-api-gen binary exists
	_, b, _, _ := runtime.Caller(0)
	csiAPIGenPath := filepath.Join(filepath.Dir(b), "../build/csi-proxy-api-gen")
	_, err := os.Lstat(csiAPIGenPath)
	require.Truef(t, err == nil, "expected err=nil, instead got err=%+v", err)

	// run the generator
	stdout, _ := runGenerator(t, csiAPIGenPath, "TestNewAPIGroup",
		"--input-dirs", "github.com/kubernetes-csi/csi-proxy/integrationtests/csiapigen/new_group/api",
		// might as well check that logging CLI args work as expected
		"-v=3")

	assert.Contains(t, stdout, "Verbosity level set to 3")
	assert.Contains(t, stdout, "Generation successful!")

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
func runGenerator(t *testing.T, csiAPIGenPath string, testName string, cliArgs ...string) (string, string) {
	stdoutFile, err := ioutil.TempFile("", "test-csi-proxy-api-gen-stdout-"+testName)
	require.Nil(t, err)
	stderrFile, err := ioutil.TempFile("", "test-csi-proxy-api-gen-stderr-"+testName)
	require.Nil(t, err)

	defer func() {

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

	// run generator through powershell
	cmd := exec.Command(csiAPIGenPath)
	out, err := cmd.CombinedOutput()
	t.Logf("%s", out)
	if err := cmd.Run(); err != nil {
		t.Fatalf("command %q failed with err=%+v", cmd.String(), err)
	}

	return readFile(t, stdoutFile.Name()),
		readFile(t, stderrFile.Name())
}
