package util

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestExits(t *testing.T, f func(), testName string) {
	if os.Getenv("TEST_EXITS") == "1" {
		f()
		return
	}
	
	cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=%s", testName))
	cmd.Env = append(os.Environ(), "TEST_EXITS=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}

	t.Fatalf("process ran with err %v, want exit status 1", err)
}
