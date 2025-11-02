package scopes

import (
	"testing"
)

func TestScopeString(t *testing.T) {
	tests := []struct {
		name     string
		scope    Scope
		expected string
	}{
		{
			name:     "SocialappUsersList",
			scope:    SocialappUsersList,
			expected: "socialapp.users.list",
		},
		{
			name:     "SocialappUsersCreate",
			scope:    SocialappUsersCreate,
			expected: "socialapp.users.create",
		},
		{
			name:     "ShortlyUrlDelete",
			scope:    ShortlyUrlDelete,
			expected: "shortly.url.delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scope.String(); got != tt.expected {
				t.Errorf("Scope.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestScopeIsValid(t *testing.T) {
	tests := []struct {
		name  string
		scope Scope
		valid bool
	}{
		{
			name:  "valid scope - SocialappUsersList",
			scope: SocialappUsersList,
			valid: true,
		},
		{
			name:  "valid scope - ShortlyUrlCreate",
			scope: ShortlyUrlCreate,
			valid: true,
		},
		{
			name:  "invalid scope - empty",
			scope: Scope(""),
			valid: false,
		},
		{
			name:  "invalid scope - non-existent",
			scope: Scope("invalid.scope.name"),
			valid: false,
		},
		{
			name:  "invalid scope - typo",
			scope: Scope("socialapp.users.lists"),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scope.IsValid(); got != tt.valid {
				t.Errorf("Scope.IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestScopeDescription(t *testing.T) {
	tests := []struct {
		name        string
		scope       Scope
		description string
	}{
		{
			name:        "SocialappUsersList",
			scope:       SocialappUsersList,
			description: "List users",
		},
		{
			name:        "SocialappUsersCreate",
			scope:       SocialappUsersCreate,
			description: "Create users",
		},
		{
			name:        "ShortlyUrlDelete",
			scope:       ShortlyUrlDelete,
			description: "Delete a url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scope.Description(); got != tt.description {
				t.Errorf("Scope.Description() = %v, want %v", got, tt.description)
			}
		})
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedScope Scope
		expectedValid bool
	}{
		{
			name:         "valid scope - socialapp.users.list",
			input:        "socialapp.users.list",
			expectedScope: SocialappUsersList,
			expectedValid: true,
		},
		{
			name:         "valid scope - shortly.url.create",
			input:        "shortly.url.create",
			expectedScope: ShortlyUrlCreate,
			expectedValid: true,
		},
		{
			name:         "invalid scope - empty",
			input:        "",
			expectedScope: Scope(""),
			expectedValid: false,
		},
		{
			name:         "invalid scope - non-existent",
			input:        "invalid.scope",
			expectedScope: Scope("invalid.scope"),
			expectedValid: false,
		},
		{
			name:         "invalid scope - typo",
			input:        "socialapp.users.lists",
			expectedScope: Scope("socialapp.users.lists"),
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope, valid := FromString(tt.input)
			if scope != tt.expectedScope {
				t.Errorf("FromString() scope = %v, want %v", scope, tt.expectedScope)
			}
			if valid != tt.expectedValid {
				t.Errorf("FromString() valid = %v, want %v", valid, tt.expectedValid)
			}
		})
	}
}

func TestMustFromString(t *testing.T) {
	t.Run("valid scope", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustFromString() panicked unexpectedly: %v", r)
			}
		}()

		scope := MustFromString("socialapp.users.list")
		if scope != SocialappUsersList {
			t.Errorf("MustFromString() = %v, want %v", scope, SocialappUsersList)
		}
	})

	t.Run("invalid scope should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustFromString() should have panicked but didn't")
			}
		}()

		MustFromString("invalid.scope")
	})
}

func TestAllScopes(t *testing.T) {
	scopes := AllScopes()

	// Check that we have at least some scopes
	if len(scopes) == 0 {
		t.Error("AllScopes() returned empty slice")
	}

	// Check that all returned scopes are valid
	for i, scope := range scopes {
		if !scope.IsValid() {
			t.Errorf("AllScopes()[%d] = %v is not valid", i, scope)
		}
	}

	// Check for duplicates
	seen := make(map[Scope]bool)
	for i, scope := range scopes {
		if seen[scope] {
			t.Errorf("AllScopes()[%d] = %v is a duplicate", i, scope)
		}
		seen[scope] = true
	}

	// Verify specific known scopes are present
	expectedScopes := []Scope{
		SocialappUsersList,
		SocialappUsersCreate,
		SocialappUsersUpdate,
		SocialappUsersDelete,
		SocialappUsersRead,
		ShortlyUrlCreate,
		ShortlyUrlDelete,
	}

	for _, expected := range expectedScopes {
		found := false
		for _, scope := range scopes {
			if scope == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("AllScopes() missing expected scope: %v", expected)
		}
	}
}

func TestStringConversion(t *testing.T) {
	// Test that converting to string and back works
	originalScope := SocialappUsersList
	asString := originalScope.String()
	backToScope, valid := FromString(asString)

	if !valid {
		t.Errorf("FromString(%q) returned invalid", asString)
	}
	if backToScope != originalScope {
		t.Errorf("Round-trip conversion failed: %v -> %v -> %v", originalScope, asString, backToScope)
	}
}

func TestAllScopesCount(t *testing.T) {
	// This test will need to be updated when scopes are added/removed
	// Currently there are 31 scopes in the OpenAPI spec
	scopes := AllScopes()
	expectedCount := 31 // Update this if scopes are added/removed in openapi.yaml

	if len(scopes) != expectedCount {
		t.Errorf("AllScopes() returned %d scopes, expected %d. If scopes were added/removed in openapi.yaml, update this test.", len(scopes), expectedCount)
	}
}

func BenchmarkScopeIsValid(b *testing.B) {
	scope := SocialappUsersList
	for i := 0; i < b.N; i++ {
		scope.IsValid()
	}
}

func BenchmarkFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FromString("socialapp.users.list")
	}
}

func BenchmarkScopeString(b *testing.B) {
	scope := SocialappUsersList
	for i := 0; i < b.N; i++ {
		_ = scope.String()
	}
}

