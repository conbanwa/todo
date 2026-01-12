package db

import (
	"testing"

	"github.com/conbanwa/todo/internal/dao/cache"
	"github.com/conbanwa/todo/internal/model"
)

func TestTeamStore_Create(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	team := &model.Team{
		Name:        "Test Team",
		Description: "A test team",
	}

	id, err := store.CreateTeam(team)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero ID")
	}
	if team.ID != id {
		t.Errorf("expected team.ID to be set to %d, got %d", id, team.ID)
	}
}

func TestTeamStore_GetByID(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	team := &model.Team{
		Name:        "Test Team",
		Description: "A test team",
	}

	id, err := store.CreateTeam(team)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}

	got, err := store.GetTeamByID(id)
	if err != nil {
		t.Fatalf("GetTeamByID failed: %v", err)
	}

	if got.ID != id {
		t.Errorf("expected ID %d, got %d", id, got.ID)
	}
	if got.Name != "Test Team" {
		t.Errorf("expected name 'Test Team', got %q", got.Name)
	}
	if got.Description != "A test team" {
		t.Errorf("expected description 'A test team', got %q", got.Description)
	}
}

func TestTeamStore_AddUserToTeam(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	// Create user
	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
	}
	userID, err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Create team
	team := &model.Team{
		Name: "Test Team",
	}
	teamID, err := store.CreateTeam(team)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}

	// Add user to team
	err = store.AddUserToTeam(userID, teamID, model.RoleAdmin)
	if err != nil {
		t.Fatalf("AddUserToTeam failed: %v", err)
	}

	// Verify user is in team
	role, err := store.GetUserRoleInTeam(userID, teamID)
	if err != nil {
		t.Fatalf("GetUserRoleInTeam failed: %v", err)
	}
	if role != model.RoleAdmin {
		t.Errorf("expected role %v, got %v", model.RoleAdmin, role)
	}
}

func TestTeamStore_GetUserTeams(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	// Create user
	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
	}
	userID, err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Create teams
	team1 := &model.Team{Name: "Team 1"}
	team1ID, err := store.CreateTeam(team1)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}

	team2 := &model.Team{Name: "Team 2"}
	team2ID, err := store.CreateTeam(team2)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}

	// Add user to teams
	err = store.AddUserToTeam(userID, team1ID, model.RoleAdmin)
	if err != nil {
		t.Fatalf("AddUserToTeam failed: %v", err)
	}
	err = store.AddUserToTeam(userID, team2ID, model.RoleMember)
	if err != nil {
		t.Fatalf("AddUserToTeam failed: %v", err)
	}

	// Get user teams
	teams, err := store.GetUserTeams(userID)
	if err != nil {
		t.Fatalf("GetUserTeams failed: %v", err)
	}

	if len(teams) != 2 {
		t.Fatalf("expected 2 teams, got %d", len(teams))
	}
}

func TestTeamStore_GetTeamMembers(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	// Create users
	user1 := &model.User{
		Username:     "user1",
		Email:        "user1@example.com",
		PasswordHash: "hash1",
	}
	user1ID, err := store.CreateUser(user1)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user2 := &model.User{
		Username:     "user2",
		Email:        "user2@example.com",
		PasswordHash: "hash2",
	}
	user2ID, err := store.CreateUser(user2)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Create team
	team := &model.Team{Name: "Test Team"}
	teamID, err := store.CreateTeam(team)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}

	// Add users to team
	err = store.AddUserToTeam(user1ID, teamID, model.RoleAdmin)
	if err != nil {
		t.Fatalf("AddUserToTeam failed: %v", err)
	}
	err = store.AddUserToTeam(user2ID, teamID, model.RoleMember)
	if err != nil {
		t.Fatalf("AddUserToTeam failed: %v", err)
	}

	// Get team members
	members, err := store.GetTeamMembers(teamID)
	if err != nil {
		t.Fatalf("GetTeamMembers failed: %v", err)
	}

	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}
}

func TestTeamStore_RemoveUserFromTeam(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	// Create user and team
	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
	}
	userID, err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	team := &model.Team{Name: "Test Team"}
	teamID, err := store.CreateTeam(team)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}

	// Add user to team
	err = store.AddUserToTeam(userID, teamID, model.RoleAdmin)
	if err != nil {
		t.Fatalf("AddUserToTeam failed: %v", err)
	}

	// Remove user from team
	err = store.RemoveUserFromTeam(userID, teamID)
	if err != nil {
		t.Fatalf("RemoveUserFromTeam failed: %v", err)
	}

	// Verify user is not in team
	_, err = store.GetUserRoleInTeam(userID, teamID)
	if err != cache.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestTeamStore_GetTeamByID_NotFound(t *testing.T) {
	store, cleanup := setupTestUserDB(t)
	defer cleanup()

	_, err := store.GetTeamByID(999)
	if err != cache.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
