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

	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	if !strings.Contains(out, "Creating team") {
		t.Fatal("Did not find expected output for new team")
	}

	out, _ = r.Run(t, false, "team", "list")
	if !strings.Contains(out, "test-team") {
		t.Fatal("Could not find newly created team name in list")
	}
}

func TestNewTeamDelete(t *testing.T) {
	r := testenv.NewExeRunner(t)

	out, err := r.Run(t, false, "team", "create", "test-team")

	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	if !strings.Contains(out, "Creating team") {
		t.Fatal("Did not find expected output for new team")
	}

	// No argument
	out, _ = r.Run(t, true, "team", "delete")

	if !strings.Contains(out, "expecting one argument, team name") {
		t.Fatal("Expected argument error")
	}

	out, _ = r.Run(t, false, "team", "list")

	if !strings.Contains(out, "test-team") {
		t.Fatal("Team was expected to still be in the team list")
	}

	// Team doesn't exist

	out, _ = r.Run(t, true, "team", "delete", "nonexistent")

	if !strings.Contains(out, "team not found") {
		t.Fatal("Expected team not found error")
	}

	out, _ = r.Run(t, false, "team", "list")

	if !strings.Contains(out, "test-team") {
		t.Fatal("Team was expected to still be in the team list")
	}

	// Delete team without confirmation
	out, _ = r.Run(t, false, "team", "delete", "test-team")

	if !strings.Contains(out, "Rerun with '--delete' flag to confirm.") {
		t.Fatal("Expected the command to ask for confirmation")
	}

	// Delete team with confirmation
	out, _ = r.Run(t, false, "team", "delete", "--delete", "test-team")

	if !strings.Contains(out, "Removing") || !strings.Contains(out, "test-team") {
		fmt.Println("OUT", out)
		t.Fatal("Expected confirmation of deleted team")
	}

	out, _ = r.Run(t, false, "team", "list")

	if strings.Contains(out, "test-team") {
		t.Fatal("Team should no longer be in the team list")
	}
}

func TestNewTeamSwitch(t *testing.T) {
	r := testenv.NewExeRunner(t)

	out, err := r.Run(t, false, "team", "create", "test-team-1")

	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	if !strings.Contains(out, "Creating team") {
		t.Fatal("Did not find expected output for new team")
	}

	// Check first team was automatically selected
	out, _ = r.Run(t, false, "team", "list")

	if !strings.Contains(out, "* test-team-1") {
		t.Fatal("Team was expected to appear as selected when listing teams")
	}

	out, err = r.Run(t, false, "team", "create", "test-team-2")

	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	if !strings.Contains(out, "Creating team") {
		t.Fatal("Did not find expected output for new team")
	}

	// Check second team was automatically selected on creation and first team is still listed
	out, _ = r.Run(t, false, "team", "list")

	if !strings.Contains(out, "* test-team-2") {
		t.Fatal("Team was expected to appear as selected when listing teams")
	}

	if !strings.Contains(out, "- test-team-1") {
		t.Fatal("Team was expected to appear as selected when listing teams")
	}

	// Check able to switch back to test-team-1
	out, _ = r.Run(t, false, "team", "switch", "test-team-1")

	if !strings.Contains(out, "* test-team-1") {
		t.Fatal("Team was expected to appear as selected when listing teams")
	}

	if !strings.Contains(out, "- test-team-2") {
		t.Fatal("Team was expected to appear as selected when listing teams")
	}

}
