// Code generated by entc, DO NOT EDIT.

package hook

import (
	"context"
	"fmt"
	"go-ent-demo/ent/gen"

	"github.com/facebookincubator/ent"
)

// The ArticleFunc type is an adapter to allow the use of ordinary
// function as Article mutator.
type ArticleFunc func(context.Context, *gen.ArticleMutation) (gen.Value, error)

// Mutate calls f(ctx, m).
func (f ArticleFunc) Mutate(ctx context.Context, m gen.Mutation) (gen.Value, error) {
	mv, ok := m.(*gen.ArticleMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *gen.ArticleMutation", m)
	}
	return f(ctx, mv)
}

// The GroupFunc type is an adapter to allow the use of ordinary
// function as Group mutator.
type GroupFunc func(context.Context, *gen.GroupMutation) (gen.Value, error)

// Mutate calls f(ctx, m).
func (f GroupFunc) Mutate(ctx context.Context, m gen.Mutation) (gen.Value, error) {
	mv, ok := m.(*gen.GroupMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *gen.GroupMutation", m)
	}
	return f(ctx, mv)
}

// The UserFunc type is an adapter to allow the use of ordinary
// function as User mutator.
type UserFunc func(context.Context, *gen.UserMutation) (gen.Value, error)

// Mutate calls f(ctx, m).
func (f UserFunc) Mutate(ctx context.Context, m gen.Mutation) (gen.Value, error) {
	mv, ok := m.(*gen.UserMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *gen.UserMutation", m)
	}
	return f(ctx, mv)
}

// On executes the given hook only of the given operation.
//
//	hook.On(Log, gen.Delete|gen.Create)
//
func On(hk gen.Hook, op gen.Op) gen.Hook {
	return func(next gen.Mutator) gen.Mutator {
		return gen.MutateFunc(func(ctx context.Context, m gen.Mutation) (gen.Value, error) {
			if m.Op().Is(op) {
				return hk(next).Mutate(ctx, m)
			}
			return next.Mutate(ctx, m)
		})
	}
}

// Reject returns a hook that rejects all operations that match op.
//
//	func (T) Hooks() []gen.Hook {
//		return []gen.Hook{
//			Reject(gen.Delete|gen.Update),
//		}
//	}
//
func Reject(op gen.Op) ent.Hook {
	return func(next gen.Mutator) gen.Mutator {
		return gen.MutateFunc(func(ctx context.Context, m gen.Mutation) (gen.Value, error) {
			if m.Op().Is(op) {
				return nil, fmt.Errorf("%s operation is not allowed", m.Op())
			}
			return next.Mutate(ctx, m)
		})
	}
}
