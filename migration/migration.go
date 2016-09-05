package migration

import (
	"github.com/almighty/almighty-core/app"
	"github.com/almighty/almighty-core/models"
	"github.com/almighty/almighty-core/remoteworkitem"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

// Perform executes the required migration of the database on startup
func Perform(ctx context.Context, db *gorm.DB, witr models.WorkItemTypeRepository) error {
	db.AutoMigrate(
		&models.WorkItem{},
		&models.WorkItemType{},
		&remoteworkitem.Tracker{},
		&remoteworkitem.TrackerQuery{},
		&remoteworkitem.TrackerItem{})
	if db.Error != nil {
		return db.Error
	}

	db.Model(&remoteworkitem.TrackerQuery{}).AddForeignKey("tracker_id", "trackers(id)", "RESTRICT", "RESTRICT")

	_, err := witr.Load(ctx, "system.issue")
	switch err.(type) {
	case models.NotFoundError:
		_, err := witr.Create(ctx, nil, "system.issue", map[string]app.FieldDefinition{
			"system.title": app.FieldDefinition{Type: &app.FieldType{Kind: "string"}, Required: true},
			"system.owner": app.FieldDefinition{Type: &app.FieldType{Kind: "user"}, Required: true},
			"system.state": app.FieldDefinition{Type: &app.FieldType{Kind: "string"}, Required: true},
		})
		if err != nil {
			return err
		}
	}
	return db.Error
}
