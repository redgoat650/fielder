package endtoend_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/redgoat650/fielder/tests/testenv"
)

func TestNewTeam(t *testing.T) {
	r := testenv.NewExeRunner(t)

	out, err := r.Run(t, false, "team", "create", "test-team")

	fmt.Println("OUT", out, err)

	out, _ = r.Run(t, false, "team", "list")
	if !strings.Contains(out, "test-team") {
		t.Fatal("Could not find newly created team name in list")
	}
}
