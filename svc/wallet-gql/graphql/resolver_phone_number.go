// Code generated by lab.weave.nl/weave/generator, DO NOT EDIT.

package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gofrs/uuid"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/auth"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/models"
	generr "lab.weave.nl/weave/generator/pkg/errors"
	"lab.weave.nl/weave/generator/utils/database"
	"lab.weave.nl/weave/generator/utils/mutation"
	"lab.weave.nl/weave/generator/utils/sanitize"
	"lab.weave.nl/weave/generator/utils/sqlconv"
)

type phoneNumberResolver struct {
	*Resolver
	Hooks *CustomPhoneNumberHooks
}

var _ PhoneNumberResolver = &phoneNumberResolver{nil, nil}

type customPhoneNumberHooks interface {
	BeforeCreateHook(ctx context.Context, tx *gorm.DB, input *CreatePhoneNumber) error
	AfterCreateHook(ctx context.Context, tx *gorm.DB, model *models.PhoneNumber) error
}

type CustomPhoneNumberHooks struct{ *Resolver }

var _ customPhoneNumberHooks = &CustomPhoneNumberHooks{nil}

func (r *mutationResolver) CreatePhoneNumber(ctx context.Context, input CreatePhoneNumber) (*models.PhoneNumber, error) {
	var err error
	var m models.PhoneNumber
	user := auth.GetUser(ctx)
	if input.UserID == nil {
		if user != nil {
			input.UserID = &user.ID
		} else {
			return nil, fmt.Errorf("%w: userId", generr.ErrFieldNotProvided)
		}
	}

	hasCreateAccess := false

	// Scope: 'api:access', Relation: HasMyUserID, UserIDField: UserID
	if auth.UserHasScope(ctx, "api:access") && cmp.Equal(input.UserID, &user.ID) {
		if !input.containsField("DeletedAt") {
			hasCreateAccess = true
		}
	}

	if !hasCreateAccess {
		return nil, generr.ErrAccessDenied
	}

	err = database.Transact(r.Resolver.DB, func(tx *gorm.DB) error {
		if err := r.Resolver.PhoneNumber().(*phoneNumberResolver).Hooks.BeforeCreateHook(ctx, tx, &input); err != nil {
			return errors.Wrap(err, "error in BeforeCreateHook")
		}
		if m, err = input.ToModel(ctx, r.Resolver); err != nil {
			return errors.Wrap(err, "converting input CreatePhoneNumber to model")
		}
		create := tx.Create(&m)
		if err := create.Error; err != nil {
			return generr.WrapAsInternal(err, "creating PhoneNumber")
		}
		if err := create.First(&m).Error; err != nil {
			return generr.WrapAsInternal(err, "getting result from create of PhoneNumber")
		}
		if err := r.Resolver.PhoneNumber().(*phoneNumberResolver).Hooks.AfterCreateHook(ctx, tx, &m); err != nil {
			return errors.Wrap(err, "error in AfterCreateHook")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *mutationResolver) UpdatePhoneNumber(ctx context.Context, id uuid.UUID, input UpdatePhoneNumber) (*models.PhoneNumber, error) {
	var err error
	var m models.PhoneNumber
	user := auth.GetUser(ctx)
	if input.UserID == nil {
		if user != nil {
			input.UserID = &user.ID
		} else {
			return nil, fmt.Errorf("%w: userId", generr.ErrFieldNotProvided)
		}
	}

	fieldCtx := graphql.GetFieldContext(ctx)

	filters := []string{fmt.Sprintf("/*path: %v*/ /*fallback*/false", fieldCtx.Path())}
	var values []interface{}
	var restrictedFields []string

	// Scope: 'api:access', Relation: HasMyUserID, UserIDField: UserID
	if auth.UserHasScope(ctx, "api:access") {
		filters = append(filters, "/*api:access HasMyUserID*/ \"phone_numbers\".\"user_id\" = ?")
		values = append(values, user.ID)
		restrictedFields = append(restrictedFields, "DeletedAt", "Unscoped")
	}

	for _, f := range restrictedFields {
		if input.containsField(f) {
			return nil, fmt.Errorf("%w: %s", generr.ErrFieldAccessDenied, sanitize.GraphFieldName(f))
		}
	}

	err = database.Transact(r.Resolver.DB, func(tx *gorm.DB) error {
		m, err = input.ToModel(ctx, r.Resolver)
		if err != nil {
			return errors.Wrap(err, "converting input UpdatePhoneNumber to model")
		}
		m.ID = id
		changes := mutation.ExtractChanges(ctx, m)
		update := tx.Model(&m).Where(strings.Join(filters, " OR "), values...).Updates(changes)
		if err := update.Error; err != nil {
			// FIXME check whether this can error on user error. If so, we should handle that elsewhere or somehow convert it here. https://lab.weave.nl/weave/generator/-/issues/105
			return generr.WrapAsInternal(err, fmt.Sprintf("updating PhoneNumber %v", m.ID))
		}
		if update.RowsAffected == 0 {
			return fmt.Errorf("%w (type=PhoneNumber,id=%v)", generr.ErrRecordNotFound, m.ID)
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

func (r *queryResolver) constructFilterPhoneNumber(ctx context.Context, joins map[string]string) (filter string, values []interface{}, restrictedFields []string) {
	fieldCtx := graphql.GetFieldContext(ctx)
	filters := []string{fmt.Sprintf("/*path: %v*/ /*fallback*/false", fieldCtx.Path())}
	user := auth.GetUser(ctx)

	// Scope: 'api:access', Relation: HasMyUserID, UserIDField: UserID
	if auth.UserHasScope(ctx, "api:access") {
		filters = append(filters, "/*api:access HasMyUserID*/ \"phone_numbers\".\"user_id\" = ?")
		values = append(values, user.ID)
		restrictedFields = append(restrictedFields, "DeletedAt", "Unscoped")
	}

	return strings.Join(filters, " OR "), values, restrictedFields
}

func (r *queryResolver) readFilterPhoneNumber(ctx context.Context, filter *PhoneNumberFilterInput, joins map[string]string) (*gorm.DB, error) {
	db := r.Resolver.DB.Model(&models.PhoneNumber{})

	for _, v := range joins {
		db = db.Joins(v)
	}

	filters, values, restrictedFields := r.constructFilterPhoneNumber(ctx, joins)
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

func (r *queryResolver) PhoneNumber(ctx context.Context, id uuid.UUID) (*models.PhoneNumber, error) {
	m := models.PhoneNumber{}
	db, err := r.readFilterPhoneNumber(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	err = db.Where(`"phone_numbers"."id" = ?`, id).First(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, generr.ErrRecordNotFound
		}
		return nil, generr.WrapAsInternal(err, fmt.Sprintf("getting PhoneNumberrecord with id %s from db", id))
	}
	return &m, nil
}

func (r *queryResolver) PhoneNumbers(ctx context.Context, limit *int, offset *int, filter *PhoneNumberFilterInput, orderBy *string, order PhoneNumberFieldName, orderDirection OrderDirection) ([]*models.PhoneNumber, error) {
	var m []*models.PhoneNumber
	db, err := r.readFilterPhoneNumber(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	orderString, err := sqlconv.ParseOrderBy(ctx, orderBy, &order, &orderDirection)
	if err != nil {
		return nil, generr.WrapAsInternal(err, "parsing order by")
	}

	err = db.Limit(*limit).Offset(*offset).Order(orderString).Find(&m).Error
	if err != nil {
		return nil, generr.WrapAsInternal(err, "getting PhoneNumbers from db")
	}
	return m, nil
}

func (r *phoneNumberResolver) ID(ctx context.Context, obj *models.PhoneNumber) (uuid.UUID, error) {
	return obj.ID, nil
}

func (r *phoneNumberResolver) User(ctx context.Context, obj *models.PhoneNumber) (*models.User, error) {
	return r.Query().User(ctx, obj.UserID)
}

var _ Directive = &phoneNumberResolver{nil, nil}

func (r *phoneNumberResolver) HasFieldAccess(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	model := castPhoneNumber(obj)
	if model == nil {
		return nil, generr.ErrAccessDenied
	}
	fieldName := graphql.GetFieldContext(ctx).Field.Name
	user := auth.GetUser(ctx)

	// Scope: 'api:access', Relation: HasMyUserID, UserIDField: UserID
	if auth.UserHasScope(ctx, "api:access") && cmp.Equal(&model.UserID, &user.ID) {
		if map[string]bool{
			"createdAt":         true,
			"id":                true,
			"phoneNumber":       true,
			"updatedAt":         true,
			"user":              true,
			"userId":            true,
			"verificationToken": true,
			"verificationType":  true,
			"verified":          true,
		}[fieldName] {
			return next(ctx)
		}
	}

	return nil, generr.ErrAccessDenied
}

func castPhoneNumber(obj interface{}) *models.PhoneNumber {
	switch res := obj.(type) {
	case **models.PhoneNumber:
		return *res
	case *models.PhoneNumber:
		return res
	case models.PhoneNumber:
		return &res
	default:
		return nil
	}
}
