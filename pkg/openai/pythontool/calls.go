package pythontool

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//go:embed counttokens.py
var counttokensPython []byte
var counttokenspath string

func init() {
	f, err := os.CreateTemp(os.TempDir(), "counttokens*.py")
	if err != nil {
		panic(err)
	}
	n, err := f.Write(counttokensPython)
	if err != nil {
		panic(err)
	}
	if n != len(counttokensPython) {
		panic("")
	}
	counttokenspath = f.Name()
	fileInfo, err := os.Stat(f.Name())
	if err != nil {
		panic(err)
	}
	currentMode := fileInfo.Mode()
	newMode := currentMode | 0111 // Set executable bits for user, group, and others.
	err = os.Chmod(f.Name(), newMode)
	if err != nil {
		panic(err)
	}
	f.Close()
}

// cl100k_base for gpt-4, gpt-3.5-turbo, text-embedding-ada-002 ____ p50k_base for Codex models, text-davinci-002, text-davinci-003
func CountTokens(encoding, prompt string) (int, error) {
	out, err := exec.Command(counttokenspath, "--encoding", encoding, "--prompt", fmt.Sprintf(`"%s"`, prompt)).CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("%v: %v (you may have to install python3)", err, string(out))
	}

	n, err := strconv.Atoi(strings.Trim(string(out), " \n\r\t"))
	if err != nil {
		return 0, err
	}

	return n, nil
}
