// Code generated by lab.weave.nl/weave/generator, DO NOT EDIT.

package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/nID-sourcecode/nid-core/svc/luarunner/models"
	generr "lab.weave.nl/weave/generator/pkg/errors"
	"lab.weave.nl/weave/generator/utils/database"
	"lab.weave.nl/weave/generator/utils/mutation"
	"lab.weave.nl/weave/generator/utils/sanitize"
	"lab.weave.nl/weave/generator/utils/sqlconv"
)

type scriptResolver struct {
	*Resolver
	Hooks *CustomScriptHooks
}

var _ ScriptResolver = &scriptResolver{nil, nil}

type customScriptHooks interface {
}

type CustomScriptHooks struct{ *Resolver }

var _ customScriptHooks = &CustomScriptHooks{nil}

func (r *mutationResolver) CreateScript(ctx context.Context, input CreateScript) (*models.Script, error) {
	var err error
	var m models.Script

	hasCreateAccess := false

	// Scope: open, Relation: None
	if !input.containsField("DeletedAt") {
		hasCreateAccess = true
	}

	if !hasCreateAccess {
		return nil, generr.ErrAccessDenied
	}

	err = database.Transact(r.Resolver.DB, func(tx *gorm.DB) error {
		if m, err = input.ToModel(ctx, r.Resolver); err != nil {
			return errors.Wrap(err, "converting input CreateScript to model")
		}
		create := tx.Create(&m)
		if err := create.Error; err != nil {
			return generr.WrapAsInternal(err, "creating Script")
		}
		if err := create.First(&m).Error; err != nil {
			return generr.WrapAsInternal(err, "getting result from create of Script")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *mutationResolver) UpdateScript(ctx context.Context, id uuid.UUID, input UpdateScript) (*models.Script, error) {
	var err error
	var m models.Script
	fieldCtx := graphql.GetFieldContext(ctx)

	filters := []string{fmt.Sprintf("/*path: %v*/ /*fallback*/false", fieldCtx.Path())}
	var values []interface{}
	var restrictedFields []string

	// Scope: open, Relation: None
	filters = append(filters, "/*ScopeOpen*/ true")

	for _, f := range restrictedFields {
		if input.containsField(f) {
			return nil, fmt.Errorf("%w: %s", generr.ErrFieldAccessDenied, sanitize.GraphFieldName(f))
		}
	}

	err = database.Transact(r.Resolver.DB, func(tx *gorm.DB) error {
		m, err = input.ToModel(ctx, r.Resolver)
		if err != nil {
			return errors.Wrap(err, "converting input UpdateScript to model")
		}
		m.ID = id
		changes := mutation.ExtractChanges(ctx, m)
		update := tx.Model(&m).Where(strings.Join(filters, " OR "), values...).Updates(changes)
		if err := update.Error; err != nil {
			// FIXME check whether this can error on user error. If so, we should handle that elsewhere or somehow convert it here. https://lab.weave.nl/weave/generator/-/issues/105
			return generr.WrapAsInternal(err, fmt.Sprintf("updating Script %v", m.ID))
		}
		if update.RowsAffected == 0 {
			return fmt.Errorf("%w (type=Script,id=%v)", generr.ErrRecordNotFound, m.ID)
		}
		if err := update.First(&m).Error; err != nil {
			return generr.WrapAsInternal(err, "getting update result")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *queryResolver) constructFilterScript(ctx context.Context, joins map[string]string) (filter string, values []interface{}, restrictedFields []string) {
	fieldCtx := graphql.GetFieldContext(ctx)
	filters := []string{fmt.Sprintf("/*path: %v*/ /*fallback*/false", fieldCtx.Path())}

	// Scope: open, Relation: None
	filters = append(filters, "/*ScopeOpen*/ true")

	return strings.Join(filters, " OR "), values, restrictedFields
}

func (r *queryResolver) readFilterScript(ctx context.Context, filter *ScriptFilterInput, joins map[string]string) (*gorm.DB, error) {
	db := r.Resolver.DB.Model(&models.Script{})

	for _, v := range joins {
		db = db.Joins(v)
	}

	filters, values, restrictedFields := r.constructFilterScript(ctx, joins)
	db = db.Where(filters, values...)

	if filter != nil {
		for _, f := range restrictedFields {
			if filter.containsField(f) {
				return nil, fmt.Errorf("%w: %s", generr.ErrFieldAccessDenied, sanitize.GraphFieldName(f))
			}
		}
		expr, args := filter.parse()
		db = db.Where(expr, args...)
		if filter.Unscoped != nil && *filter.Unscoped {
			db = db.Unscoped()
		}
	}

	return db, nil
}

func (r *queryResolver) Script(ctx context.Context, id uuid.UUID) (*models.Script, error) {
	m := models.Script{}
	db, err := r.readFilterScript(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	err = db.Where(`"scripts"."id" = ?`, id).First(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, generr.ErrRecordNotFound
		}
		return nil, generr.WrapAsInternal(err, fmt.Sprintf("getting Scriptrecord with id %s from db", id))
	}
	return &m, nil
}

func (r *queryResolver) Scripts(ctx context.Context, limit *int, offset *int, filter *ScriptFilterInput, orderBy *string, order ScriptFieldName, orderDirection OrderDirection) ([]*models.Script, error) {
	var m []*models.Script
	db, err := r.readFilterScript(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	orderString, err := sqlconv.ParseOrderBy(ctx, orderBy, &order, &orderDirection)
	if err != nil {
		return nil, generr.WrapAsInternal(err, "parsing order by")
	}

	err = db.Limit(*limit).Offset(*offset).Order(orderString).Find(&m).Error
	if err != nil {
		return nil, generr.WrapAsInternal(err, "getting Scripts from db")
	}
	return m, nil
}

func (r *scriptResolver) ID(ctx context.Context, obj *models.Script) (uuid.UUID, error) {
	return obj.ID, nil
}

var _ Directive = &scriptResolver{nil, nil}

func (r *scriptResolver) HasFieldAccess(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	model := castScript(obj)
	if model == nil {
		return nil, generr.ErrAccessDenied
	}
	fieldName := graphql.GetFieldContext(ctx).Field.Name

	// Scope: open, Relation: None
	if map[string]bool{
		"createdAt":      true,
		"eventType":      true,
		"id":             true,
		"organisationId": true,
		"script":         true,
		"updatedAt":      true,
	}[fieldName] {
		return next(ctx)
	}

	return nil, generr.ErrAccessDenied
}

func castScript(obj interface{}) *models.Script {
	switch res := obj.(type) {
	case **models.Script:
		return *res
	case *models.Script:
		return res
	case models.Script:
		return &res
	default:
		return nil
	}
}
