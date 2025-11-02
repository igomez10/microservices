package main

import (
	"fmt"

	"github.com/igomez10/microservices/socialapp/pkg/scopes"
)

// Example demonstrating how to use the scopes package
func main() {
	fmt.Println("=== OAuth2 Scopes Demo ===")
	fmt.Println()

	// Example 1: Using type-safe constants
	fmt.Println("1. Using type-safe scope constants:")
	userListScope := scopes.SocialappUsersList
	fmt.Printf("   Scope: %s\n", userListScope.String())
	fmt.Printf("   Description: %s\n", userListScope.Description())
	fmt.Printf("   Is Valid: %v\n\n", userListScope.IsValid())

	// Example 2: Converting string to scope with validation
	fmt.Println("2. Converting string to scope:")
	scopeStr := "socialapp.users.create"
	scope, valid := scopes.FromString(scopeStr)
	if valid {
		fmt.Printf("   ✓ Valid scope: %s - %s\n\n", scope, scope.Description())
	} else {
		fmt.Printf("   ✗ Invalid scope: %s\n\n", scopeStr)
	}

	// Example 3: Handling invalid scopes
	fmt.Println("3. Handling invalid scopes:")
	invalidScopeStr := "socialapp.users.invalid"
	_, valid = scopes.FromString(invalidScopeStr)
	if !valid {
		fmt.Printf("   ✗ Invalid scope detected: %s\n\n", invalidScopeStr)
	}

	// Example 4: Listing all available scopes
	fmt.Println("4. All available scopes:")
	allScopes := scopes.AllScopes()
	fmt.Printf("   Total: %d scopes\n", len(allScopes))
	fmt.Println("   First 5 scopes:")
	for i, scope := range allScopes[:5] {
		fmt.Printf("      %d. %s - %s\n", i+1, scope, scope.Description())
	}
	fmt.Println()

	// Example 5: Using scopes in a slice
	fmt.Println("5. Working with multiple scopes:")
	requiredScopes := []scopes.Scope{
		scopes.SocialappUsersList,
		scopes.SocialappUsersRead,
		scopes.SocialappUsersUpdate,
	}
	fmt.Println("   Required scopes for user management:")
	for _, scope := range requiredScopes {
		fmt.Printf("      - %s\n", scope.String())
	}
	fmt.Println()

	// Example 6: Scope validation in action
	fmt.Println("6. Validating user-provided scopes:")
	userProvidedScopes := []string{
		"socialapp.users.list",
		"socialapp.users.create",
		"invalid.scope.here",
		"socialapp.users.delete",
	}
	validatedScopes := []scopes.Scope{}
	for _, scopeStr := range userProvidedScopes {
		scope, valid := scopes.FromString(scopeStr)
		if valid {
			fmt.Printf("   ✓ %s\n", scopeStr)
			validatedScopes = append(validatedScopes, scope)
		} else {
			fmt.Printf("   ✗ %s (invalid)\n", scopeStr)
		}
	}
	fmt.Printf("   Valid scopes: %d/%d\n\n", len(validatedScopes), len(userProvidedScopes))

	// Example 7: Using MustFromString (safe because we know it's valid)
	fmt.Println("7. Using MustFromString for known valid scopes:")
	knownValidScope := scopes.MustFromString("socialapp.roles.list")
	fmt.Printf("   Scope: %s - %s\n", knownValidScope, knownValidScope.Description())
}
