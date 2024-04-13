// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/dstgo/maxwell/app/data/ent/schema"
	"github.com/dstgo/maxwell/app/data/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userFields[4].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(int64)
	// userDescUpdatedAt is the schema descriptor for updated_at field.
	userDescUpdatedAt := userFields[5].Descriptor()
	// user.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	user.DefaultUpdatedAt = userDescUpdatedAt.Default.(int64)
	// user.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	user.UpdateDefaultUpdatedAt = userDescUpdatedAt.UpdateDefault.(func() int64)
}