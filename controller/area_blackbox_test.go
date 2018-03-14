package controller_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/fabric8-services/fabric8-wit/account"
	"github.com/fabric8-services/fabric8-wit/app"
	"github.com/fabric8-services/fabric8-wit/app/test"
	"github.com/fabric8-services/fabric8-wit/application"
	"github.com/fabric8-services/fabric8-wit/area"
	. "github.com/fabric8-services/fabric8-wit/controller"
	"github.com/fabric8-services/fabric8-wit/gormapplication"
	"github.com/fabric8-services/fabric8-wit/gormsupport"
	"github.com/fabric8-services/fabric8-wit/gormtestsupport"
	"github.com/fabric8-services/fabric8-wit/log"
	"github.com/fabric8-services/fabric8-wit/resource"
	"github.com/fabric8-services/fabric8-wit/space"
	"github.com/fabric8-services/fabric8-wit/spacetemplate"
	testsupport "github.com/fabric8-services/fabric8-wit/test"
	"github.com/goadesign/goa"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestAreaREST struct {
	gormtestsupport.DBTestSuite
	db *gormapplication.GormDB
}

func TestRunAreaREST(t *testing.T) {
	resource.Require(t, resource.Database)
	pwd, err := os.Getwd()
	require.NoError(t, err)
	suite.Run(t, &TestAreaREST{DBTestSuite: gormtestsupport.NewDBTestSuite(pwd + "/../config.yaml")})
}

func (rest *TestAreaREST) SetupTest() {
	rest.DBTestSuite.SetupTest()
	rest.db = gormapplication.NewGormDB(rest.DB)
}

func (rest *TestAreaREST) SecuredController() (*goa.Service, *AreaController) {
	svc := testsupport.ServiceAsUser("Area-Service", testsupport.TestIdentity)
	return svc, NewAreaController(svc, rest.db, rest.Configuration)
}

func (rest *TestAreaREST) SecuredControllerWithIdentity(idn *account.Identity) (*goa.Service, *AreaController) {
	svc := testsupport.ServiceAsUser("Area-Service", *idn)
	return svc, NewAreaController(svc, rest.db, rest.Configuration)
}

func (rest *TestAreaREST) UnSecuredController() (*goa.Service, *AreaController) {
	svc := goa.New("Area-Service")
	return svc, NewAreaController(svc, rest.db, rest.Configuration)
}

func (rest *TestAreaREST) TestSuccessCreateChildArea() {
	// given
	sp, parentArea := createSpaceAndArea(rest.T(), rest.db)
	parentID := parentArea.ID
	ci := newCreateChildAreaPayload("TestSuccessCreateChildArea")
	owner, err := rest.db.Identities().Load(context.Background(), sp.OwnerID)
	require.NoError(rest.T(), err)
	svc, ctrl := rest.SecuredControllerWithIdentity(owner)
	// when
	_, created := test.CreateChildAreaCreated(rest.T(), svc.Context, svc, ctrl, parentID.String(), ci)
	// then
	assert.Equal(rest.T(), *ci.Data.Attributes.Name, *created.Data.Attributes.Name)
	fmt.Println(*created.Data.Relationships.Parent.Data.ID)
	assert.Equal(rest.T(), parentID.String(), *created.Data.Relationships.Parent.Data.ID)

	// try creating child area with different identity: should fail
	otherIdentity := &account.Identity{
		Username:     "non-space-owner-identity",
		ProviderType: account.KeycloakIDP,
	}
	errInCreateOther := rest.db.Identities().Create(context.Background(), otherIdentity)
	require.NoError(rest.T(), errInCreateOther)
	svc, ctrl = rest.SecuredControllerWithIdentity(otherIdentity)
	test.CreateChildAreaForbidden(rest.T(), svc.Context, svc, ctrl, parentID.String(), ci)
}

func (rest *TestAreaREST) TestSuccessCreateMultiChildArea() {
	/*
		TestAreaREST ---> TestSuccessCreateMultiChildArea-0 ----> TestSuccessCreateMultiChildArea-0-0
	*/
	// given
	sp, parentArea := createSpaceAndArea(rest.T(), rest.db)
	parentID := parentArea.ID
	ci := newCreateChildAreaPayload("TestSuccessCreateMultiChildArea-0")
	owner, err := rest.db.Identities().Load(context.Background(), sp.OwnerID)
	require.NoError(rest.T(), err)
	svc, ctrl := rest.SecuredControllerWithIdentity(owner)
	// when
	_, created := test.CreateChildAreaCreated(rest.T(), svc.Context, svc, ctrl, parentID.String(), ci)
	// then
	assert.Equal(rest.T(), *ci.Data.Attributes.Name, *created.Data.Attributes.Name)
	assert.Equal(rest.T(), parentID.String(), *created.Data.Relationships.Parent.Data.ID)
	// Create a child of the child created above.
	ci = newCreateChildAreaPayload("TestSuccessCreateMultiChildArea-0-0")
	newParentID := *created.Data.Relationships.Parent.Data.ID
	// when
	_, created = test.CreateChildAreaCreated(rest.T(), svc.Context, svc, ctrl, newParentID, ci)
	// then
	assert.Equal(rest.T(), *ci.Data.Attributes.Name, *created.Data.Attributes.Name)
	assert.NotNil(rest.T(), *created.Data.Attributes.CreatedAt)
	assert.NotNil(rest.T(), *created.Data.Attributes.Version)
	assert.Equal(rest.T(), newParentID, *created.Data.Relationships.Parent.Data.ID)
	assert.Contains(rest.T(), *created.Data.Relationships.Children.Links.Self, "children")
}

func (rest *TestAreaREST) TestConflictCreatDuplicateChildArea() {
	// given
	sp, parentArea := createSpaceAndArea(rest.T(), rest.db)
	parentID := parentArea.ID
	ci := newCreateChildAreaPayload(uuid.NewV4().String())
	owner, err := rest.db.Identities().Load(context.Background(), sp.OwnerID)
	require.NoError(rest.T(), err)
	svc, ctrl := rest.SecuredControllerWithIdentity(owner)
	// when
	_, created := test.CreateChildAreaCreated(rest.T(), svc.Context, svc, ctrl, parentID.String(), ci)
	// then
	assert.Equal(rest.T(), *ci.Data.Attributes.Name, *created.Data.Attributes.Name)
	assert.Equal(rest.T(), parentID.String(), *created.Data.Relationships.Parent.Data.ID)

	// try creating the same area again
	test.CreateChildAreaConflict(rest.T(), svc.Context, svc, ctrl, parentID.String(), ci)

}

func (rest *TestAreaREST) TestFailCreateChildAreaMissingName() {
	// given
	sp, parentArea := createSpaceAndArea(rest.T(), rest.db)
	parentID := parentArea.ID
	createChildAreaPayload := newCreateChildAreaPayload("will remove below")
	createChildAreaPayload.Data.Attributes.Name = nil
	owner, err := rest.db.Identities().Load(context.Background(), sp.OwnerID)
	require.NoError(rest.T(), err)
	svc, ctrl := rest.SecuredControllerWithIdentity(owner)
	// when/then
	test.CreateChildAreaBadRequest(rest.T(), svc.Context, svc, ctrl, parentID.String(), createChildAreaPayload)
}

func (rest *TestAreaREST) TestFailCreateChildAreaWithInvalidsParent() {
	// given
	createChildAreaPayload := newCreateChildAreaPayload("TestFailCreateChildAreaWithInvalidsParent")
	svc, ctrl := rest.SecuredController()
	// when/then
	test.CreateChildAreaNotFound(rest.T(), svc.Context, svc, ctrl, uuid.NewV4().String(), createChildAreaPayload)
}

func (rest *TestAreaREST) TestFailCreateChildAreaNotAuthorized() {
	// given
	_, parentArea := createSpaceAndArea(rest.T(), rest.db)
	parentID := parentArea.ID
	createChildAreaPayload := newCreateChildAreaPayload("TestFailCreateChildAreaNotAuthorized")
	svc, ctrl := rest.UnSecuredController()
	// when/then
	test.CreateChildAreaUnauthorized(rest.T(), svc.Context, svc, ctrl, parentID.String(), createChildAreaPayload)
}

func (rest *TestAreaREST) TestShowArea() {
	rest.T().Run("Success", func(t *testing.T) {
		// Setup
		_, a := createSpaceAndArea(t, rest.db)
		svc, ctrl := rest.SecuredController()
		t.Run("OK", func(t *testing.T) {
			// when
			res, _ := test.ShowAreaOK(t, svc.Context, svc, ctrl, a.ID.String(), nil, nil)
			//then
			assertResponseHeaders(t, res)
		})

		t.Run("Using ExpiredIfModifedSince Header", func(t *testing.T) {
			// when
			ifModifiedSince := app.ToHTTPTime(a.UpdatedAt.Add(-1 * time.Hour))
			res, _ := test.ShowAreaOK(t, svc.Context, svc, ctrl, a.ID.String(), &ifModifiedSince, nil)
			//then
			assertResponseHeaders(t, res)
		})

		t.Run("Using ExpiredIfNoneMatch Header", func(t *testing.T) {
			// when
			ifNoneMatch := "foo"
			res, _ := test.ShowAreaOK(t, svc.Context, svc, ctrl, a.ID.String(), nil, &ifNoneMatch)
			//then
			assertResponseHeaders(t, res)
		})

		t.Run("Not Modified Using IfModifedSince Header", func(t *testing.T) {
			// when
			ifModifiedSince := app.ToHTTPTime(a.UpdatedAt)
			res := test.ShowAreaNotModified(t, svc.Context, svc, ctrl, a.ID.String(), &ifModifiedSince, nil)
			//then
			assertResponseHeaders(t, res)
		})

		t.Run("Not Modified IfNoneMatch Header", func(t *testing.T) {
			// when
			ifNoneMatch := app.GenerateEntityTag(a)
			res := test.ShowAreaNotModified(t, svc.Context, svc, ctrl, a.ID.String(), nil, &ifNoneMatch)
			//then
			assertResponseHeaders(t, res)
		})
	})

	rest.T().Run("Failure", func(t *testing.T) {
		// Setup
		svc, ctrl := rest.SecuredController()
		t.Run("Not Found", func(t *testing.T) {
			// when/then
			test.ShowAreaNotFound(t, svc.Context, svc, ctrl, uuid.NewV4().String(), nil, nil)
		})
	})
}

func (rest *TestAreaREST) TestShowAreaOKUsingExpiredIfNoneMatchHeader() {
	// given
	_, a := createSpaceAndArea(rest.T(), rest.db)
	svc, ctrl := rest.SecuredController()
	// when
	ifNoneMatch := "foo"
	res, _ := test.ShowAreaOK(rest.T(), svc.Context, svc, ctrl, a.ID.String(), nil, &ifNoneMatch)
	//then
	assertResponseHeaders(rest.T(), res)
}

func (rest *TestAreaREST) TestShowAreaNotModifiedUsingIfModifedSinceHeader() {
	// given
	_, a := createSpaceAndArea(rest.T(), rest.db)
	svc, ctrl := rest.SecuredController()
	// when
	ifModifiedSince := app.ToHTTPTime(a.UpdatedAt)
	res := test.ShowAreaNotModified(rest.T(), svc.Context, svc, ctrl, a.ID.String(), &ifModifiedSince, nil)
	//then
	assertResponseHeaders(rest.T(), res)
}

func (rest *TestAreaREST) TestShowAreaNotModifiedIfNoneMatchHeader() {
	// given
	_, a := createSpaceAndArea(rest.T(), rest.db)
	svc, ctrl := rest.SecuredController()
	// when
	ifNoneMatch := app.GenerateEntityTag(a)
	res := test.ShowAreaNotModified(rest.T(), svc.Context, svc, ctrl, a.ID.String(), nil, &ifNoneMatch)
	//then
	assertResponseHeaders(rest.T(), res)
}

func (rest *TestAreaREST) createChildArea(name string, parent area.Area, svc *goa.Service, ctrl *AreaController) *app.AreaSingle {
	ci := newCreateChildAreaPayload(name)
	// when
	_, created := test.CreateChildAreaCreated(rest.T(), svc.Context, svc, ctrl, parent.ID.String(), ci)
	return created
}

func (rest *TestAreaREST) TestShowChildrenArea() {
	// Setup
	sp, parentArea := createSpaceAndArea(rest.T(), rest.db)
	owner, err := rest.db.Identities().Load(context.Background(), sp.OwnerID)
	require.NoError(rest.T(), err)
	svc, ctrl := rest.SecuredControllerWithIdentity(owner)
	childArea := rest.createChildArea("TestShowChildrenArea", parentArea, svc, ctrl)
	rest.T().Run("Success", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			res, result := test.ShowChildrenAreaOK(rest.T(), svc.Context, svc, ctrl, parentArea.ID.String(), nil, nil)
			assert.Equal(rest.T(), 1, len(result.Data))
			assertResponseHeaders(rest.T(), res)
		})
		t.Run("Using ExpiredIfModifedSince Header", func(t *testing.T) {
			ifModifiedSince := app.ToHTTPTime(parentArea.UpdatedAt.Add(-1 * time.Hour))
			res, result := test.ShowChildrenAreaOK(rest.T(), svc.Context, svc, ctrl, parentArea.ID.String(), &ifModifiedSince, nil)
			assert.Equal(rest.T(), 1, len(result.Data))
			assertResponseHeaders(rest.T(), res)
		})

		t.Run("Using ExpiredIfNoneMatch Header", func(t *testing.T) {
			ifNoneMatch := "foo"
			res, result := test.ShowChildrenAreaOK(rest.T(), svc.Context, svc, ctrl, parentArea.ID.String(), nil, &ifNoneMatch)
			assert.Equal(rest.T(), 1, len(result.Data))
			assertResponseHeaders(rest.T(), res)
		})

		t.Run("Not Modified Using IfModifedSince Header", func(t *testing.T) {
			ifModifiedSince := app.ToHTTPTime(*childArea.Data.Attributes.UpdatedAt)
			res := test.ShowChildrenAreaNotModified(rest.T(), svc.Context, svc, ctrl, parentArea.ID.String(), &ifModifiedSince, nil)
			assertResponseHeaders(rest.T(), res)
		})

		t.Run("Not Modified IfNoneMatch Header", func(t *testing.T) {
			modelChildArea := convertAreaToModel(*childArea)
			ifNoneMatch := app.GenerateEntityTag(modelChildArea)
			res := test.ShowChildrenAreaNotModified(rest.T(), svc.Context, svc, ctrl, parentArea.ID.String(), nil, &ifNoneMatch)
			assertResponseHeaders(rest.T(), res)
		})
	})
}

func ConvertAreaToModel(appArea app.AreaSingle) area.Area {
	return area.Area{
		ID:      *appArea.Data.ID,
		Version: *appArea.Data.Attributes.Version,
		Lifecycle: gormsupport.Lifecycle{
			UpdatedAt: *appArea.Data.Attributes.UpdatedAt,
		},
	}
}

func newCreateChildAreaPayload(name string) *app.CreateChildAreaPayload {
	areaType := area.APIStringTypeAreas
	return &app.CreateChildAreaPayload{
		Data: &app.Area{
			Type: areaType,
			Attributes: &app.AreaAttributes{
				Name: &name,
			},
		},
	}
}

func createSpaceAndArea(t *testing.T, db *gormapplication.GormDB) (space.Space, area.Area) {
	var areaObj area.Area
	var spaceObj space.Space
	application.Transactional(db, func(app application.Application) error {
		owner := &account.Identity{
			Username:     "new-space-owner-identity",
			ProviderType: account.KeycloakIDP,
		}
		errCreateOwner := app.Identities().Create(context.Background(), owner)
		require.NoError(t, errCreateOwner)

		spaceObj = space.Space{
			Name:            "TestAreaREST-" + uuid.NewV4().String(),
			OwnerID:         owner.ID,
			SpaceTemplateID: spacetemplate.SystemLegacyTemplateID,
		}
		_, err := app.Spaces().Create(context.Background(), &spaceObj)
		require.NoError(t, err)
		name := "Main Area-" + uuid.NewV4().String()
		areaObj = area.Area{
			Name:    name,
			SpaceID: spaceObj.ID,
		}
		err = app.Areas().Create(context.Background(), &areaObj)
		require.NoError(t, err)
		return nil
	})
	log.Info(nil, nil, "Space and root area created")
	return spaceObj, areaObj
}
