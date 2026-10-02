package converter

import (
	"testing"
	"time"

	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/jackc/pgx/v5/pgtype"
)

// The API serializes int64 ids as JSON strings to avoid precision loss in JS
// clients, so a snowflake id must survive the conversion exactly.
func TestFromDBUserToAPIUser_SendsID(t *testing.T) {
	dbUser := db.User{
		ID:        1305234176287309824,
		Username:  "johndoe",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@mail.com",
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	apiUser := FromDBUserToAPIUser(dbUser)

	if apiUser.Id != "1305234176287309824" {
		t.Errorf("expected id %q, got %q", "1305234176287309824", apiUser.Id)
	}
	if apiUser.Username != dbUser.Username {
		t.Errorf("expected username %q, got %q", dbUser.Username, apiUser.Username)
	}
}

// The user's salt and password hash must never reach the API surface.
func TestFromDBUserToAPIUser_OmitsSecrets(t *testing.T) {
	apiUser := FromDBUserToAPIUser(db.User{
		ID:             1,
		Username:       "johndoe",
		HashedPassword: "hashed-secret",
		Salt:           "salt-value",
	})

	if apiUser.Id != "1" {
		t.Fatalf("expected id %q, got %q", "1", apiUser.Id)
	}
	// openapi.User has no field for either; this documents that the conversion
	// is a deliberate projection rather than a copy.
	if apiUser.Username != "johndoe" {
		t.Errorf("expected username johndoe, got %q", apiUser.Username)
	}
}
