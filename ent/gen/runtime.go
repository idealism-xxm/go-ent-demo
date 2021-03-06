// Code generated by entc, DO NOT EDIT.

package gen

import (
	"go-ent-demo/ent/gen/article"
	"go-ent-demo/ent/gen/user"
	"go-ent-demo/ent/schema"
	"time"

	"github.com/facebookincubator/ent"
)

// The init function reads all schema descriptors with runtime
// code (default values, validators or hooks) and stitches it
// to their package variables.
func init() {
	articleMixin := schema.Article{}.Mixin()
	articleMixinFields := [...][]ent.Field{
		articleMixin[0].Fields(),
	}
	articleFields := schema.Article{}.Fields()
	_ = articleFields
	// articleDescCreatedAt is the schema descriptor for created_at field.
	articleDescCreatedAt := articleMixinFields[0][0].Descriptor()
	// article.DefaultCreatedAt holds the default value on creation for the created_at field.
	article.DefaultCreatedAt = articleDescCreatedAt.Default.(func() time.Time)
	// articleDescUpdatedAt is the schema descriptor for updated_at field.
	articleDescUpdatedAt := articleMixinFields[0][1].Descriptor()
	// article.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	article.DefaultUpdatedAt = articleDescUpdatedAt.Default.(func() time.Time)
	// article.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	article.UpdateDefaultUpdatedAt = articleDescUpdatedAt.UpdateDefault.(func() time.Time)
	// articleDescTitle is the schema descriptor for title field.
	articleDescTitle := articleFields[0].Descriptor()
	// article.TitleValidator is a validator for the "title" field. It is called by the builders before save.
	article.TitleValidator = articleDescTitle.Validators[0].(func(string) error)
	userMixin := schema.User{}.Mixin()
	userMixinFields := [...][]ent.Field{
		userMixin[0].Fields(),
	}
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userMixinFields[0][0].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescUpdatedAt is the schema descriptor for updated_at field.
	userDescUpdatedAt := userMixinFields[0][1].Descriptor()
	// user.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	user.DefaultUpdatedAt = userDescUpdatedAt.Default.(func() time.Time)
	// user.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	user.UpdateDefaultUpdatedAt = userDescUpdatedAt.UpdateDefault.(func() time.Time)
	// userDescNickname is the schema descriptor for nickname field.
	userDescNickname := userFields[0].Descriptor()
	// user.NicknameValidator is a validator for the "nickname" field. It is called by the builders before save.
	user.NicknameValidator = userDescNickname.Validators[0].(func(string) error)
	// userDescUsername is the schema descriptor for username field.
	userDescUsername := userFields[1].Descriptor()
	// user.UsernameValidator is a validator for the "username" field. It is called by the builders before save.
	user.UsernameValidator = func() func(string) error {
		validators := userDescUsername.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(username string) error {
			for _, fn := range fns {
				if err := fn(username); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescEmail is the schema descriptor for email field.
	userDescEmail := userFields[2].Descriptor()
	// user.EmailValidator is a validator for the "email" field. It is called by the builders before save.
	user.EmailValidator = userDescEmail.Validators[0].(func(string) error)
}
