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
	"lab.weave.nl/nid/nid-core/svc/info-manager/models"
	generr "lab.weave.nl/weave/generator/pkg/errors"
	"lab.weave.nl/weave/generator/utils/sanitize"
	"lab.weave.nl/weave/generator/utils/sqlconv"
)

type scriptSourceResolver struct {
	*Resolver
	Hooks *CustomScriptSourceHooks
}

var _ ScriptSourceResolver = &scriptSourceResolver{nil, nil}

type customScriptSourceHooks interface {
	AfterReadGetSignedURL(ctx context.Context, db *gorm.DB, model *models.ScriptSource) error
}

type CustomScriptSourceHooks struct{ *Resolver }

var _ customScriptSourceHooks = &CustomScriptSourceHooks{nil}

func (r *queryResolver) constructFilterScriptSource(ctx context.Context, joins map[string]string) (filter string, values []interface{}, restrictedFields []string) {
	fieldCtx := graphql.GetFieldContext(ctx)
	filters := []string{fmt.Sprintf("/*path: %v*/ /*fallback*/false", fieldCtx.Path())}

	// Scope: open, Relation: None
	filters = append(filters, "/*ScopeOpen*/ true")

	return strings.Join(filters, " OR "), values, restrictedFields
}

func (r *queryResolver) readFilterScriptSource(ctx context.Context, filter *ScriptSourceFilterInput, joins map[string]string) (*gorm.DB, error) {
	db := r.Resolver.DB.Model(&models.ScriptSource{})

	for _, v := range joins {
		db = db.Joins(v)
	}

	filters, values, restrictedFields := r.constructFilterScriptSource(ctx, joins)
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

func (r *queryResolver) ScriptSource(ctx context.Context, id uuid.UUID) (*models.ScriptSource, error) {
	m := models.ScriptSource{}
	db, err := r.readFilterScriptSource(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	err = db.Where(`"script_sources"."id" = ?`, id).First(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, generr.ErrRecordNotFound
		}
		return nil, generr.WrapAsInternal(err, fmt.Sprintf("getting ScriptSourcerecord with id %s from db", id))
	}
	if err := r.Resolver.ScriptSource().(*scriptSourceResolver).Hooks.AfterReadGetSignedURL(ctx, db, &m); err != nil {
		return nil, errors.Wrap(err, "error in AfterReadGetSignedURL")
	}
	return &m, nil
}

func (r *queryResolver) ScriptSources(ctx context.Context, limit *int, offset *int, filter *ScriptSourceFilterInput, orderBy *string, order ScriptSourceFieldName, orderDirection OrderDirection) ([]*models.ScriptSource, error) {
	var m []*models.ScriptSource
	db, err := r.readFilterScriptSource(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	orderString, err := sqlconv.ParseOrderBy(ctx, orderBy, &order, &orderDirection)
	if err != nil {
		return nil, generr.WrapAsInternal(err, "parsing order by")
	}

	err = db.Limit(*limit).Offset(*offset).Order(orderString).Find(&m).Error
	if err != nil {
		return nil, generr.WrapAsInternal(err, "getting ScriptSources from db")
	}
	return m, nil
}

func (r *scriptSourceResolver) ID(ctx context.Context, obj *models.ScriptSource) (uuid.UUID, error) {
	return obj.ID, nil
}

var _ Directive = &scriptSourceResolver{nil, nil}

func (r *scriptSourceResolver) HasFieldAccess(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	model := castScriptSource(obj)
	if model == nil {
		return nil, generr.ErrAccessDenied
	}
	fieldName := graphql.GetFieldContext(ctx).Field.Name

	// Scope: open, Relation: None
	if map[string]bool{
		"changeDescription": true,
		"checksum":          true,
		"createdAt":         true,
		"id":                true,
		"rawScript":         true,
		"scriptId":          true,
		"signedUrl":         true,
		"updatedAt":         true,
		"version":           true,
	}[fieldName] {
		return next(ctx)
	}

	return nil, generr.ErrAccessDenied
}

func castScriptSource(obj interface{}) *models.ScriptSource {
	switch res := obj.(type) {
	case **models.ScriptSource:
		return *res
	case *models.ScriptSource:
		return res
	case models.ScriptSource:
		return &res
	default:
		return nil
	}
}
