/*
** The dal (Data Access Layer) package applies business logic to retrieving and creating objects in the system.
** Where the store does validation and the storengine provides access to different underlying databases, the
** dal checks access, performs any needed transformations and where appropriate performs writes to multiple
** underlying stores.
 */
package dal

import "github.com/materials-commons/mc/internal/store"

type DAL struct {
	db store.DB
}
