package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/Dataman-Cloud/baker/config"
	"github.com/Dataman-Cloud/baker/model"
	"github.com/Dataman-Cloud/baker/store"
)

// Store is a middleware function that initializes the Store and attaches to
// the context of every http.Request.
func Store(cf *config.Config) gin.HandlerFunc {
	v := setupStaticUsersStore(cf)
	return func(c *gin.Context) {
		store.ToContext(c, v)
		c.Next()
	}
}

// helper function to create the Store from the CLI context config-path.
func setupStaticUsersStore(cf *config.Config) store.Store {
	// setup Store
	if cf.Users != nil {
		staticUsers := make(map[string]*model.StaticUser)
		for k, v := range cf.Users {
			staticUsers[k] = &model.StaticUser{
				Username: k,
				Password: v.Password,
			}
		}
		return store.NewStaticUsersStore(staticUsers)
	}
	return nil
}
