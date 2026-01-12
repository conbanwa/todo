package model

import (
	"testing"
	"time"
)

func TestTeam_Fields(t *testing.T) {
	now := time.Now()
	team := Team{
		ID:          1,
		Name:        "Test Team",
		Description: "A test team",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if team.ID != 1 {
		t.Errorf("Team.ID = %v, want 1", team.ID)
	}
	if team.Name != "Test Team" {
		t.Errorf("Team.Name = %v, want Test Team", team.Name)
	}
	if team.Description != "A test team" {
		t.Errorf("Team.Description = %v, want A test team", team.Description)
	}
}

func TestTeam_ZeroValue(t *testing.T) {
	var team Team
	if team.ID != 0 {
		t.Errorf("Team.ID should be 0 for zero value, got %v", team.ID)
	}
	if team.Name != "" {
		t.Errorf("Team.Name should be empty for zero value, got %v", team.Name)
	}
	if !team.CreatedAt.IsZero() {
		t.Errorf("Team.CreatedAt should be zero for zero value")
	}
	if !team.UpdatedAt.IsZero() {
		t.Errorf("Team.UpdatedAt should be zero for zero value")
	}
}
