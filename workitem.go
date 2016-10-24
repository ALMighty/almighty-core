package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/goadesign/goa"

	"github.com/almighty/almighty-core/app"
	"github.com/almighty/almighty-core/application"
	"github.com/almighty/almighty-core/models"
	"github.com/almighty/almighty-core/query/simple"
)

// WorkitemController implements the workitem resource.
type WorkitemController struct {
	*goa.Controller
	db application.DB
}

// NewWorkitemController creates a workitem controller.
func NewWorkitemController(service *goa.Service, db application.DB) *WorkitemController {
	return &WorkitemController{Controller: service.NewController("WorkitemController"), db: db}
}

// Show runs the show action.
func (c *WorkitemController) Show(ctx *app.ShowWorkitemContext) error {
	return application.Transactional(c.db, func(appl application.Application) error {
		wi, err := appl.WorkItems().Load(ctx.Context, ctx.ID)
		if err != nil {
			switch err.(type) {
			case models.NotFoundError:
				log.Printf("not found, id=%s", ctx.ID)
				return goa.ErrNotFound(err.Error())
			default:
				return err
			}
		}
		return ctx.OK(wi)
	})
}

func parseInts(s *string) ([]int, error) {
	if s == nil || len(*s) == 0 {
		return []int{}, nil
	}
	split := strings.Split(*s, ",")
	result := make([]int, len(split))
	for index, value := range split {
		converted, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		result[index] = converted
	}
	return result, nil
}

func parseLimit(pageParameter *string) (s *int, l int, e error) {
	params, err := parseInts(pageParameter)
	if err != nil {
		return nil, 0, err
	}

	if len(params) > 1 {
		return &params[0], params[1], nil
	}
	if len(params) > 0 {
		return nil, params[0], nil
	}
	return nil, 100, nil
}

// List runs the list action
func (c *WorkitemController) List(ctx *app.ListWorkitemContext) error {
	exp, err := query.Parse(ctx.Filter)
	if err != nil {
		return goa.ErrBadRequest(fmt.Sprintf("could not parse filter: %s", err.Error()))
	}
	start, limit, err := parseLimit(ctx.Page)
	if err != nil {
		return goa.ErrBadRequest(fmt.Sprintf("could not parse paging: %s", err.Error()))
	}
	return application.Transactional(c.db, func(appl application.Application) error {
		result, err := appl.WorkItems().List(ctx.Context, exp, start, &limit)
		if err != nil {
			return goa.ErrInternal(fmt.Sprintf("Error listing work items: %s", err.Error()))
		}
		return ctx.OK(result)
	})
}

// Create runs the create action.
func (c *WorkitemController) Create(ctx *app.CreateWorkitemContext) error {
	return application.Transactional(c.db, func(appl application.Application) error {
		wi, err := appl.WorkItems().Create(ctx.Context, ctx.Payload.Type, ctx.Payload.Fields)

		if err != nil {
			switch err := err.(type) {
			case models.BadParameterError, models.ConversionError:
				return goa.ErrBadRequest(err.Error())
			default:
				return goa.ErrInternal(err.Error())
			}
		}
		ctx.ResponseData.Header().Set("Location", app.WorkitemHref(wi.ID))
		return ctx.Created(wi)
	})
}

// Delete runs the delete action.
func (c *WorkitemController) Delete(ctx *app.DeleteWorkitemContext) error {
	return application.Transactional(c.db, func(appl application.Application) error {
		err := appl.WorkItems().Delete(ctx.Context, ctx.ID)
		if err != nil {
			switch err.(type) {
			case models.NotFoundError:
				return goa.ErrNotFound(err.Error())
			default:
				return goa.ErrInternal(err.Error())
			}
		}
		return ctx.OK([]byte{})
	})
}

// Update runs the update action.
func (c *WorkitemController) Update(ctx *app.UpdateWorkitemContext) error {
	return application.Transactional(c.db, func(appl application.Application) error {

		toSave := app.WorkItem{
			ID:      ctx.ID,
			Type:    ctx.Payload.Type,
			Version: ctx.Payload.Version,
			Fields:  ctx.Payload.Fields,
		}
		wi, err := appl.WorkItems().Save(ctx.Context, toSave)

		if err != nil {
			switch err := err.(type) {
			case models.BadParameterError, models.ConversionError:
				return goa.ErrBadRequest(err.Error())
			default:
				return goa.ErrInternal(err.Error())
			}
		}
		return ctx.OK(wi)
	})
}
