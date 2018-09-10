package cache

import (
	"sync"

	"github.com/materials-commons/mc/internal/store"
)

type UsersCacheByAPIKey struct {
	sync.RWMutex
	Cache map[string]store.UserSchema
	store.UsersStore
}

func (c *UsersCacheByAPIKey) Get(apikey string) (*store.UserSchema, error) {
	c.RLock()
	if user, ok := c.Cache[apikey]; ok {
		c.RUnlock()
		return &user, nil
	}

	c.RUnlock()

	c.Lock()
	defer c.Unlock()

	user, err := c.UsersStore.GetUserByAPIKey(apikey)
	if err != nil {
		return nil, err
	}

	c.Cache[apikey] = user
	return &user, nil
}

func (c *UsersCacheByAPIKey) Delete(apikey string) {
	c.Lock()
	delete(c.Cache, apikey)
	c.Unlock()
}

func (c *UsersCacheByAPIKey) ReloadUser(apikey, id string) error {
	c.Lock()
	defer c.Unlock()

	delete(c.Cache, apikey)

	user, err := c.UsersStore.GetUserByID(id)
	if err != nil {
		return err
	}

	c.Cache[user.APIKey] = user
	return nil
}
