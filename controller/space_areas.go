package controller

import (
	"github.com/fabric8-services/fabric8-wit/app"
	"github.com/fabric8-services/fabric8-wit/application"
	"github.com/fabric8-services/fabric8-wit/jsonapi"
	"github.com/fabric8-services/fabric8-wit/log"

	"github.com/goadesign/goa"
)

// SpaceAreasController implements the space-Areas resource.
type SpaceAreasController struct {
	*goa.Controller
	db     application.DB
	config SpaceAreasControllerConfig
}

//SpaceAreasControllerConfig the configuration for the SpaceAreasController
type SpaceAreasControllerConfig interface {
	GetCacheControlAreas() string
}

// NewSpaceAreasController creates a space-Areas controller.
func NewSpaceAreasController(service *goa.Service, db application.DB, config SpaceAreasControllerConfig) *SpaceAreasController {
	return &SpaceAreasController{
		Controller: service.NewController("SpaceAreasController"),
		db:         db,
		config:     config,
	}
}

// List runs the list action.
func (c *SpaceAreasController) List(ctx *app.ListSpaceAreasContext) error {
	err := application.Transactional(c.db, func(appl application.Application) error {
		if err := appl.Spaces().CheckExists(ctx, ctx.SpaceID); err != nil {
			return err
		}
		areas, err := appl.Areas().List(ctx, ctx.SpaceID)
		if err != nil {
			return err
		}
		for _, a := range areas {
			log.Info(ctx, map[string]interface{}{"area_id": a.ID}, "Found space area with id %s", a.ID)
		}
		return ctx.ConditionalEntities(areas, c.config.GetCacheControlAreas, func() error {
			res := &app.AreaList{}
			res.Data = ConvertAreas(appl, ctx.Request, areas, addResolvedPath)
			return ctx.OK(res)
		})
	})
	if err != nil {
		return jsonapi.JSONErrorResponse(ctx, err)
	}
	return nil

}
