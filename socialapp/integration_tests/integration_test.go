package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func TestMain(m *testing.M) {
	// run tests
	// add jitter
	if os.Getenv("ADD_TEST_JITTER") != "" {
		jitterInSeconds := uuid.New().ID() % 60
		log.Printf("Adding test jitter of %d seconds", jitterInSeconds)
		time.Sleep(time.Duration(jitterInSeconds))
	}

	Setup()
	code := m.Run()
	os.Exit(code)
}

func TestListUsers(t *testing.T) {
	if err := ListUsersLifecycle(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestCreateUser(t *testing.T) {
	if err := CreateUserLifecycle(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestFollowCycle(t *testing.T) {
	if err := FollowLifeCycle(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestGetExpectedFeed(t *testing.T) {
	if err := GetExpectedFeed(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestGetAccessToken(t *testing.T) {
	if err := GetAccessToken(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterUserFlow(t *testing.T) {
	if err := RegisterUserFlow(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestChangePassword(t *testing.T) {
	if err := ChangePassword(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestRoleLifecycle(t *testing.T) {
	if err := RoleLifecycle(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestScopeLifeCycle(t *testing.T) {
	if err := ScopeLifecycle(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestUserRoleLifeCycle(t *testing.T) {
	if err := UserRoleLifeCycle(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestCacheRequestSameUser(t *testing.T) {
	if err := CacheRequestSameUser(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestURLLifeCycle(t *testing.T) {
	if err := URLLifeCycle(context.Background()); err != nil {
		t.Fatal(err)
	}
}
