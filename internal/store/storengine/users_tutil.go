package storengine

import (
	"testing"

	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/tutils/assert"
	r "gopkg.in/gorethink/gorethink.v4"
)

func AddDefaultUsersToStoreEngine(t *testing.T, e UsersStoreEngine) {
	users := []model.UserSchema{
		{ID: "tuser@test.com", APIKey: "tuser@test.com apikey", Fullname: "tuser", Password: "tuser-password"},
	}

	for _, user := range users {
		_, err := e.AddUser(user)
		assert.Okf(t, err, "Failed to add user %s", user.ID)
	}
}

func CleanupUsersStoreEngine(e UsersStoreEngine) {
	if re, ok := e.(*UsersRethinkdb); ok {
		session := re.Session
		_, _ = r.Table("users").Delete().RunWrite(session)
	}
}
