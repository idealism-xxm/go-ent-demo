// Code generated by entc, DO NOT EDIT.

package gen

import (
	"context"
	"errors"
	"fmt"
	"go-ent-demo/ent/gen/article"
	"go-ent-demo/ent/gen/user"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
)

// ArticleCreate is the builder for creating a Article entity.
type ArticleCreate struct {
	config
	mutation *ArticleMutation
	hooks    []Hook
}

// SetCreatedAt sets the created_at field.
func (ac *ArticleCreate) SetCreatedAt(t time.Time) *ArticleCreate {
	ac.mutation.SetCreatedAt(t)
	return ac
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (ac *ArticleCreate) SetNillableCreatedAt(t *time.Time) *ArticleCreate {
	if t != nil {
		ac.SetCreatedAt(*t)
	}
	return ac
}

// SetUpdatedAt sets the updated_at field.
func (ac *ArticleCreate) SetUpdatedAt(t time.Time) *ArticleCreate {
	ac.mutation.SetUpdatedAt(t)
	return ac
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (ac *ArticleCreate) SetNillableUpdatedAt(t *time.Time) *ArticleCreate {
	if t != nil {
		ac.SetUpdatedAt(*t)
	}
	return ac
}

// SetTitle sets the title field.
func (ac *ArticleCreate) SetTitle(s string) *ArticleCreate {
	ac.mutation.SetTitle(s)
	return ac
}

// SetContent sets the content field.
func (ac *ArticleCreate) SetContent(s string) *ArticleCreate {
	ac.mutation.SetContent(s)
	return ac
}

// SetAuthorID sets the author edge to User by id.
func (ac *ArticleCreate) SetAuthorID(id int) *ArticleCreate {
	ac.mutation.SetAuthorID(id)
	return ac
}

// SetNillableAuthorID sets the author edge to User by id if the given value is not nil.
func (ac *ArticleCreate) SetNillableAuthorID(id *int) *ArticleCreate {
	if id != nil {
		ac = ac.SetAuthorID(*id)
	}
	return ac
}

// SetAuthor sets the author edge to User.
func (ac *ArticleCreate) SetAuthor(u *User) *ArticleCreate {
	return ac.SetAuthorID(u.ID)
}

// Save creates the Article in the database.
func (ac *ArticleCreate) Save(ctx context.Context) (*Article, error) {
	if _, ok := ac.mutation.CreatedAt(); !ok {
		v := article.DefaultCreatedAt()
		ac.mutation.SetCreatedAt(v)
	}
	if _, ok := ac.mutation.UpdatedAt(); !ok {
		v := article.DefaultUpdatedAt()
		ac.mutation.SetUpdatedAt(v)
	}
	if _, ok := ac.mutation.Title(); !ok {
		return nil, errors.New("gen: missing required field \"title\"")
	}
	if v, ok := ac.mutation.Title(); ok {
		if err := article.TitleValidator(v); err != nil {
			return nil, fmt.Errorf("gen: validator failed for field \"title\": %v", err)
		}
	}
	if _, ok := ac.mutation.Content(); !ok {
		return nil, errors.New("gen: missing required field \"content\"")
	}
	var (
		err  error
		node *Article
	)
	if len(ac.hooks) == 0 {
		node, err = ac.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ArticleMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ac.mutation = mutation
			node, err = ac.sqlSave(ctx)
			return node, err
		})
		for i := len(ac.hooks) - 1; i >= 0; i-- {
			mut = ac.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ac.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ac *ArticleCreate) SaveX(ctx context.Context) *Article {
	v, err := ac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ac *ArticleCreate) sqlSave(ctx context.Context) (*Article, error) {
	var (
		a     = &Article{config: ac.config}
		_spec = &sqlgraph.CreateSpec{
			Table: article.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: article.FieldID,
			},
		}
	)
	if value, ok := ac.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: article.FieldCreatedAt,
		})
		a.CreatedAt = value
	}
	if value, ok := ac.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: article.FieldUpdatedAt,
		})
		a.UpdatedAt = value
	}
	if value, ok := ac.mutation.Title(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: article.FieldTitle,
		})
		a.Title = value
	}
	if value, ok := ac.mutation.Content(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: article.FieldContent,
		})
		a.Content = value
	}
	if nodes := ac.mutation.AuthorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   article.AuthorTable,
			Columns: []string{article.AuthorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, ac.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	a.ID = int(id)
	return a, nil
}
