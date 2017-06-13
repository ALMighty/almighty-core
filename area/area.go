package area

import (
	"context"
	"fmt"
	"time"

	"github.com/almighty/almighty-core/application/repository"
	"github.com/almighty/almighty-core/errors"
	"github.com/almighty/almighty-core/gormsupport"
	"github.com/almighty/almighty-core/log"
	"github.com/almighty/almighty-core/path"

	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
	errs "github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const APIStringTypeAreas = "areas"

// Area describes a single Area
type Area struct {
	gormsupport.Lifecycle
	ID      uuid.UUID `sql:"type:uuid default uuid_generate_v4()" gorm:"primary_key"` // This is the ID PK field
	SpaceID uuid.UUID `sql:"type:uuid"`
	Path    path.Path
	Name    string
	Version int
}

// GetETagData returns the field values to use to generate the ETag
func (m Area) GetETagData() []interface{} {
	return []interface{}{m.ID, m.Version}
}

// GetLastModified returns the last modification time
func (m Area) GetLastModified() time.Time {
	return m.UpdatedAt.Truncate(time.Second)
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m *GormAreaRepository) TableName() string {
	return "areas"
}

// Repository describes interactions with Areas
type Repository interface {
	repository.Exister
	Create(ctx context.Context, u *Area) error
	List(ctx context.Context, spaceID uuid.UUID) ([]Area, error)
	Load(ctx context.Context, id uuid.UUID) (*Area, error)
	LoadMultiple(ctx context.Context, ids []uuid.UUID) ([]Area, error)
	ListChildren(ctx context.Context, parentArea *Area) ([]Area, error)
	Query(funcs ...func(*gorm.DB) *gorm.DB) ([]Area, error)
	Root(ctx context.Context, spaceID uuid.UUID) (*Area, error)
}

// NewAreaRepository creates a new storage type.
func NewAreaRepository(db *gorm.DB) Repository {
	return &GormAreaRepository{db: db}
}

// GormAreaRepository is the implementation of the storage interface for Areas.
type GormAreaRepository struct {
	db *gorm.DB
}

// Create creates a new record.
func (m *GormAreaRepository) Create(ctx context.Context, u *Area) error {
	defer goa.MeasureSince([]string{"goa", "db", "area", "create"}, time.Now())
	u.ID = uuid.NewV4()
	err := m.db.Create(u).Error
	if err != nil {
		// ( name, spaceID ,path ) needs to be unique
		if gormsupport.IsUniqueViolation(err, "areas_name_space_id_path_unique") {
			return errors.NewBadParameterError("name & space_id & path", u.Name+" & "+u.SpaceID.String()+" & "+u.Path.String()).Expected("unique")
		}
		log.Error(ctx, map[string]interface{}{}, "error adding Area: %s", err.Error())
		return err
	}
	return nil
}

// List all Areas related to a single item
func (m *GormAreaRepository) List(ctx context.Context, spaceID uuid.UUID) ([]Area, error) {
	defer goa.MeasureSince([]string{"goa", "db", "Area", "query"}, time.Now())
	var objs []Area
	err := m.db.Where("space_id = ?", spaceID).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return objs, nil
}

// Load a single Area regardless of parent
func (m *GormAreaRepository) Load(ctx context.Context, id uuid.UUID) (*Area, error) {
	defer goa.MeasureSince([]string{"goa", "db", "Area", "get"}, time.Now())
	var obj Area

	tx := m.db.Where("id = ?", id).First(&obj)
	if tx.RecordNotFound() {
		return nil, errors.NewNotFoundError("Area", id.String())
	}
	if tx.Error != nil {
		return nil, errors.NewInternalError(tx.Error.Error())
	}
	return &obj, nil
}

// Exists returns true|false where an object exists with an identifier
func (m *GormAreaRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	defer goa.MeasureSince([]string{"goa", "db", "area", "exists"}, time.Now())
	var exists bool
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1 FROM %[1]s
			WHERE
				id=$1
				AND deleted_at IS NULL
		)`, m.TableName())

	if err := m.db.CommonDB().QueryRow(query, id).Scan(&exists); err != nil {
		return false, errs.Wrapf(err, "failed to check if an area exists with this id %v", id)
	}
	return exists, nil
}

// Load multiple areas
func (m *GormAreaRepository) LoadMultiple(ctx context.Context, ids []uuid.UUID) ([]Area, error) {
	defer goa.MeasureSince([]string{"goa", "db", "Area", "getmultiple"}, time.Now())
	var objs []Area

	for i := 0; i < len(ids); i++ {
		m.db = m.db.Or("id = ?", ids[i])
	}
	tx := m.db.Find(&objs)
	if tx.Error != nil {
		return nil, errors.NewInternalError(tx.Error.Error())
	}
	return objs, nil
}

// ListChildren fetches all Areas belonging to a parent - list all child areas.
func (m *GormAreaRepository) ListChildren(ctx context.Context, parentArea *Area) ([]Area, error) {
	defer goa.MeasureSince([]string{"goa", "db", "Area", "querychild"}, time.Now())
	var objs []Area

	tx := m.db.Where("path ~ ?", path.ToExpression(parentArea.Path, parentArea.ID)).Find(&objs)
	if tx.RecordNotFound() {
		return nil, errors.NewNotFoundError("Area", parentArea.ID.String())
	}
	if tx.Error != nil {
		return nil, errors.NewInternalError(tx.Error.Error())
	}
	return objs, nil
}

// Root fetches the Root Areas inside a space.
func (m *GormAreaRepository) Root(ctx context.Context, spaceID uuid.UUID) (*Area, error) {
	defer goa.MeasureSince([]string{"goa", "db", "Area", "root"}, time.Now())
	var rootArea []Area

	parentPathOfRootArea := path.Path{}
	rootArea, err := m.Query(FilterBySpaceID(spaceID), FilterByPath(parentPathOfRootArea))
	if len(rootArea) != 1 {
		return nil, errors.NewInternalError(fmt.Sprintf("Single Root area not found for space %s", spaceID.String()))
	}
	if err != nil {
		return nil, err
	}
	return &rootArea[0], nil
}

// Query exposes an open ended Query model for Area
func (m *GormAreaRepository) Query(funcs ...func(*gorm.DB) *gorm.DB) ([]Area, error) {
	defer goa.MeasureSince([]string{"goa", "db", "area", "query"}, time.Now())
	var objs []Area

	err := m.db.Scopes(funcs...).Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errs.WithStack(err)
	}

	log.Debug(nil, map[string]interface{}{
		"area_list": objs,
	}, "Area query executed successfully!")

	return objs, nil
}

// FilterBySpaceID is a gorm filter for a Belongs To relationship.
func FilterBySpaceID(spaceID uuid.UUID) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("space_id = ?", spaceID)
	}
}

// FilterByName is a gorm filter by 'name'
func FilterByName(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name).Limit(1)
	}
}

// FilterByPath is a gorm filter by 'path' of the parent area for any given area.
func FilterByPath(pathOfParent path.Path) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("path = ?", pathOfParent.Convert()).Limit(1)
	}
}
