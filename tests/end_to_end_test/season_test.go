package endtoend_test

import (
	"strings"
	"testing"

	"github.com/redgoat650/fielder/tests/testenv"
)

func TestAddSeason(t *testing.T) {
	// Create a team
	r := testenv.NewExeRunner(t)

	out, err := r.Run(t, false, "team", "create", "test-team")

	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	if !strings.Contains(out, "Creating team") {
		t.Fatal("Did not find expected output for new team")
	}

	// Add a season without a name
	out, _ = r.Run(t, true, "season", "add")
	if !strings.Contains(out, "expecting one season name argument") {
		t.Fatal("Season name error did not print")
	}

	// Add a season with a name
	out, _ = r.Run(t, false, "season", "add", "test-season")
	if !strings.Contains(out, "Adding season called \"test-season\"") {
		t.Fatal("Success feedback was not found")
	}

	// Add a season with a name that already exists
	out, _ = r.Run(t, true, "season", "add", "test-season")
	if !strings.Contains(out, "season already exists with this name") {
		t.Fatal("Expecting error message indicating season name already exists")
	}

	// List seasons
	out, _ = r.Run(t, false, "season", "list")
	if !strings.Contains(out, "Listing seasons for team \"test-team\"") {
		t.Fatal("Season list did not indicate the correct team")
	}

	if !strings.Contains(out, "* test-season") {
		t.Fatal("Season did not appear in the list")
	}

	// Add a second season
	_, _ = r.Run(t, false, "season", "add", "test-season-2")

	// List seasons
	out, _ = r.Run(t, false, "season", "list")
	if !strings.Contains(out, "Listing seasons for team \"test-team\"") {
		t.Fatal("Season list did not indicate the correct team")
	}

	if !strings.Contains(out, "* test-season-2") {
		t.Fatal("Second season did not appear selected in the list")
	}

	if !strings.Contains(out, "- test-season") {
		t.Fatal("First season did not appear in the list")
	}

	// Switch back to select the first season
	out, _ = r.Run(t, false, "season", "switch", "test-season")
	if !strings.Contains(out, "Switching to season \"test-season\"") {
		t.Fatal("Did not indicate successful team switch")
	}

	if !strings.Contains(out, "* test-season") {
		t.Fatal("List did not show correct season selected")
	}

	if !strings.Contains(out, "- test-season-2") {
		t.Fatal("List did not show previously selected season")
	}

	// Switch to the currently selected season
	out, _ = r.Run(t, false, "season", "switch", "test-season")
	if !strings.Contains(out, "Season \"test-season\" was already selected.") {
		t.Fatal("Did not indicate season was already selected.")
	}

	// Unselect the season
	out, _ = r.Run(t, false, "season", "switch", "--none")
	if !strings.Contains(out, "- test-season") {
		t.Fatal("Expecting test-season to be unselected")
	}

	// Reselect the season
	out, _ = r.Run(t, false, "season", "switch", "test-season")
	if !strings.Contains(out, "* test-season") {
		t.Fatal("Expected test-season to be selected")
	}

	// Delete a season that doesn't exist
	out, _ = r.Run(t, true, "season", "delete", "test-season-nonexistent")
	if !strings.Contains(out, "season not found") {
		t.Fatal("Expecting indication that season was not found")
	}

	// Delete the selected season without confirmation
	out, _ = r.Run(t, false, "season", "delete", "test-season")
	if !strings.Contains(out, "Rerun with '--delete' flag to confirm deletion.") {
		t.Fatal("Expecting confirmation help message")
	}

	// Delete the non-selected season - expect current season to remain selected
	out, _ = r.Run(t, false, "season", "delete", "test-season-2", "-D")
	if !strings.Contains(out, "Will permanently delete season: \"test-season-2\"") {
		t.Fatal("Expected output confirmation of deleted season")
	}

	out, _ = r.Run(t, false, "season", "list")
	if !strings.Contains(out, "* test-season") {
		t.Fatal("Expected current season to still be selected")
	}

	// Create another season
	out, _ = r.Run(t, false, "season", "add", "test-season-3")
	if !strings.Contains(out, "Adding season called \"test-season-3\"") {
		t.Fatal("Expecting output indicating successful season creation")
	}

	out, _ = r.Run(t, false, "season", "list")
	if !strings.Contains(out, "* test-season-3") {
		t.Fatal("Expecting active season to be set to the new season")
	}

	// Delete the selected season - expect current season to be deselected
	out, _ = r.Run(t, false, "season", "delete", "test-season-3", "-D")
	if !strings.Contains(out, "Will permanently delete season: \"test-season-3\"") {
		t.Fatal("Expected output confirmation of deleted season")
	}

	out, _ = r.Run(t, false, "season", "list")
	if !strings.Contains(out, "- test-season") {
		t.Fatal("Expecting no active season set")
	}
}
