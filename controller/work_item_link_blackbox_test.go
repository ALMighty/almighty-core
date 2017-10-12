package controller_test

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/fabric8-services/fabric8-wit/app"
	"github.com/fabric8-services/fabric8-wit/app/test"
	"github.com/fabric8-services/fabric8-wit/application"
	. "github.com/fabric8-services/fabric8-wit/controller"
	"github.com/fabric8-services/fabric8-wit/gormapplication"
	"github.com/fabric8-services/fabric8-wit/gormtestsupport"
	"github.com/fabric8-services/fabric8-wit/jsonapi"
	"github.com/fabric8-services/fabric8-wit/resource"
	"github.com/fabric8-services/fabric8-wit/rest"
	"github.com/fabric8-services/fabric8-wit/space"
	testsupport "github.com/fabric8-services/fabric8-wit/test"
	tf "github.com/fabric8-services/fabric8-wit/test/testfixture"
	testtoken "github.com/fabric8-services/fabric8-wit/test/token"
	"github.com/fabric8-services/fabric8-wit/workitem"
	"github.com/fabric8-services/fabric8-wit/workitem/link"

	"github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

//-----------------------------------------------------------------------------
// Test Suite setup
//-----------------------------------------------------------------------------

// The workItemLinkSuite has state the is relevant to all tests.
// It implements these interfaces from the suite package: SetupAllSuite, SetupTestSuite, TearDownAllSuite, TearDownTestSuite
type workItemLinkSuite struct {
	gormtestsupport.DBTestSuite
	testDir          string
	workItemLinkCtrl *WorkItemLinkController

	svc                      *goa.Service
	workItemLinkTypeCtrl     *WorkItemLinkTypeController
	workItemLinkCategoryCtrl *WorkItemLinkCategoryController
	workItemCtrl             *WorkitemController
	workItemsCtrl            *WorkitemsController
	workItemRelsLinksCtrl    *WorkItemRelationshipsLinksController
	spaceCtrl                *SpaceController
	typeCtrl                 *WorkitemtypeController
	// These IDs can safely be used by all tests
	bug1ID               uuid.UUID
	bug2ID               uuid.UUID
	bug3ID               uuid.UUID
	feature1ID           uuid.UUID
	userLinkCategoryID   uuid.UUID
	bugBlockerLinkTypeID uuid.UUID
	userSpaceID          uuid.UUID
	appDB                application.DB
}

// cleanup removes all DB entries that will be created or have been created
// with this test suite. We need to remove them completely and not only set the
// "deleted_at" field, which is why we need the Unscoped() function.
func (s *workItemLinkSuite) cleanup() {
	// Delete all work item links for now
	db := s.DB.Unscoped().Delete(&link.WorkItemLink{})
	require.Nil(s.T(), db.Error)

	// Delete work item link types and categories by name.
	// They will be created during the tests but have to be deleted by name
	// rather than ID, unlike the work items or work item links.
	db = db.Unscoped().Delete(&link.WorkItemLinkType{Name: "test-bug-blocker"})
	require.Nil(s.T(), db.Error)
	db = db.Unscoped().Delete(&link.WorkItemLinkCategory{Name: "test-user"})
	require.Nil(s.T(), db.Error)
	if s.userSpaceID != uuid.Nil {
		db = db.Unscoped().Delete(&space.Space{ID: s.userSpaceID})
		require.Nil(s.T(), db.Error)
	}
}

// The SetupSuite method will run before the tests in the suite are run.
// It sets up a database connection for all the tests in this suite without polluting global space.
func (s *workItemLinkSuite) SetupSuite() {
	s.DBTestSuite.SetupSuite()
	s.testDir = filepath.Join("test-files", "work_item_link")
	s.appDB = gormapplication.NewGormDB(s.DB)
}

// The SetupTest method will be run before every test in the suite.
// SetupTest ensures that none of the work item links that we will create already exist.
// It will also make sure that some resources that we rely on do exists.
func (s *workItemLinkSuite) SetupTest() {
	s.DBTestSuite.SetupTest()
	svc := goa.New("TestWorkItemLinkType-Service")
	require.NotNil(s.T(), svc)
	s.workItemLinkTypeCtrl = NewWorkItemLinkTypeController(svc, gormapplication.NewGormDB(s.DB), s.Configuration)
	require.NotNil(s.T(), s.workItemLinkTypeCtrl)

	svc = goa.New("TestWorkItemLinkCategory-Service")
	require.NotNil(s.T(), svc)
	s.workItemLinkCategoryCtrl = NewWorkItemLinkCategoryController(svc, gormapplication.NewGormDB(s.DB))
	require.NotNil(s.T(), s.workItemLinkCategoryCtrl)

	svc = goa.New("TestWorkItemLinkSpace-Service")
	require.NotNil(s.T(), svc)
	s.spaceCtrl = NewSpaceController(svc, gormapplication.NewGormDB(s.DB), s.Configuration, &DummyResourceManager{})
	require.NotNil(s.T(), s.spaceCtrl)

	svc = goa.New("TestWorkItemType-Service")
	s.typeCtrl = NewWorkitemtypeController(svc, gormapplication.NewGormDB(s.DB), s.Configuration)
	require.NotNil(s.T(), s.typeCtrl)

	svc = goa.New("TestWorkItemLink-Service")
	require.NotNil(s.T(), svc)
	s.workItemLinkCtrl = NewWorkItemLinkController(svc, gormapplication.NewGormDB(s.DB), s.Configuration)
	require.NotNil(s.T(), s.workItemLinkCtrl)

	svc = goa.New("TestWorkItemRelationshipsLinks-Service")
	require.NotNil(s.T(), svc)
	s.workItemRelsLinksCtrl = NewWorkItemRelationshipsLinksController(svc, gormapplication.NewGormDB(s.DB), s.Configuration)
	require.NotNil(s.T(), s.workItemRelsLinksCtrl)

	// create a test identity
	fxt := tf.NewTestFixture(s.T(), s.DB,
		tf.CreateWorkItemLinkEnvironment(),
		tf.WorkItemTypes(2, tf.SetWorkItemTypeNames([]string{"bug", "feature"})),
		tf.WorkItems(4,
			tf.SetWorkItemTitles([]string{"bug1", "bug2", "bug3", "feature1"}),
			func(fxt *tf.TestFixture, idx int) error {
				fxt.WorkItems[idx].Type = fxt.WorkItemTypeByName("bug").ID
				if idx == 3 {
					fxt.WorkItems[idx].Type = fxt.WorkItemTypeByName("feature").ID
				}
				return nil
			},
		),
	)

	s.svc = testsupport.ServiceAsUser("TestWorkItem-Service", *fxt.Identities[0])
	require.NotNil(s.T(), s.svc)
	s.workItemCtrl = NewWorkitemController(svc, gormapplication.NewGormDB(s.DB), s.Configuration)
	require.NotNil(s.T(), s.workItemCtrl)
	s.workItemsCtrl = NewWorkitemsController(svc, gormapplication.NewGormDB(s.DB), s.Configuration)
	require.NotNil(s.T(), s.workItemsCtrl)
	s.userSpaceID = fxt.Spaces[0].ID
	s.T().Logf("Created link space with ID: %s\n", s.userSpaceID)

	// Create 3 work items (bug1, bug2, and feature1)
	s.bug1ID = fxt.WorkItemByTitle("bug1").ID
	s.bug2ID = fxt.WorkItemByTitle("bug2").ID
	s.bug3ID = fxt.WorkItemByTitle("bug3").ID
	s.feature1ID = fxt.WorkItemByTitle("feature1").ID

	// Create a work item link category
	s.userLinkCategoryID = fxt.WorkItemLinkCategories[0].ID
	s.T().Logf("Created link category with ID: %s\n", s.userLinkCategoryID)

	s.bugBlockerLinkTypeID = fxt.WorkItemLinkTypes[0].ID
	s.T().Logf("Created link type with ID: %s\n", s.bugBlockerLinkTypeID)
}

// creates a work item with the given name and type and returns its ID
func (s *workItemLinkSuite) createWorkItem(typeID uuid.UUID, name string) uuid.UUID {
	payload := newCreateWorkItemPayload(s.userSpaceID, typeID, name)
	_, wi := test.CreateWorkitemsCreated(s.T(), s.svc.Context, s.svc, s.workItemsCtrl, *payload.Data.Relationships.Space.Data.ID, payload)
	require.NotNil(s.T(), wi)
	s.T().Logf("Created bug with ID: %d\n", *wi.Data.ID)
	return *wi.Data.ID
}

//-----------------------------------------------------------------------------
// helper method
//-----------------------------------------------------------------------------

// CreateWorkItemLinkCategory creates a work item link category
func newCreateWorkItemLinkCategoryPayload(name string) *app.CreateWorkItemLinkCategoryPayload {
	description := "This work item link category is managed by an admin user."
	// Use the goa generated code to create a work item link category
	return &app.CreateWorkItemLinkCategoryPayload{
		Data: &app.WorkItemLinkCategoryData{
			Type: link.EndpointWorkItemLinkCategories,
			Attributes: &app.WorkItemLinkCategoryAttributes{
				Name:        &name,
				Description: &description,
			},
		},
	}
}

// CreateWorkItem defines a work item link
func newCreateWorkItemPayload(spaceID uuid.UUID, workItemType uuid.UUID, title string) *app.CreateWorkitemsPayload {
	spaceRelatedURL := rest.AbsoluteURL(&http.Request{Host: "api.service.domain.org"}, app.SpaceHref(spaceID.String()))
	witRelatedURL := rest.AbsoluteURL(&http.Request{Host: "api.service.domain.org"}, app.WorkitemtypeHref(spaceID.String(), workItemType))
	payload := app.CreateWorkitemsPayload{
		Data: &app.WorkItem{
			Attributes: map[string]interface{}{
				workitem.SystemTitle: title,
				workitem.SystemState: workitem.SystemStateClosed,
			},
			Relationships: &app.WorkItemRelationships{
				BaseType: &app.RelationBaseType{
					Data: &app.BaseTypeData{
						ID:   workItemType,
						Type: "workitemtypes",
					},
					Links: &app.GenericLinks{
						Self:    &witRelatedURL,
						Related: &witRelatedURL,
					},
				},
				Space: app.NewSpaceRelation(spaceID, spaceRelatedURL),
			},
			Type: "workitems",
		},
	}
	return &payload
}

// CreateWorkItemLinkType defines a work item link type
func newCreateWorkItemLinkTypePayload(name string, categoryID, spaceID uuid.UUID) *app.CreateWorkItemLinkTypePayload {
	description := "Specify that one bug blocks another one."
	lt := link.WorkItemLinkType{
		Name:           name,
		Description:    &description,
		Topology:       link.TopologyNetwork,
		ForwardName:    "forward name string for " + name,
		ReverseName:    "reverse name string for " + name,
		LinkCategoryID: categoryID,
		SpaceID:        spaceID,
	}
	reqLong := &http.Request{Host: "api.service.domain.org"}
	payload := ConvertWorkItemLinkTypeFromModel(reqLong, lt)
	// The create payload is required during creation. Simply copy data over.
	return &app.CreateWorkItemLinkTypePayload{
		Data: payload.Data,
	}
}

// newCreateWorkItemLinkPayload returns the payload to create a work item link
func newCreateWorkItemLinkPayload(sourceID, targetID, linkTypeID uuid.UUID) *app.CreateWorkItemLinkPayload {
	lt := link.WorkItemLink{
		SourceID:   sourceID,
		TargetID:   targetID,
		LinkTypeID: linkTypeID,
	}
	payload := ConvertLinkFromModel(lt)
	// The create payload is required during creation. Simply copy data over.
	return &app.CreateWorkItemLinkPayload{
		Data: payload.Data,
	}
}

// newUpdateWorkItemLinkPayload returns the payload to update a work item link
func newUpdateWorkItemLinkPayload(linkID, sourceID, targetID, linkTypeID uuid.UUID) *app.UpdateWorkItemLinkPayload {
	lt := link.WorkItemLink{
		ID:         linkID,
		SourceID:   sourceID,
		TargetID:   targetID,
		LinkTypeID: linkTypeID,
	}
	payload := ConvertLinkFromModel(lt)
	// The create payload is required during creation. Simply copy data over.
	return &app.UpdateWorkItemLinkPayload{
		Data: payload.Data,
	}
}

//-----------------------------------------------------------------------------
// Actual tests
//-----------------------------------------------------------------------------

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSuiteWorkItemLinks(t *testing.T) {
	resource.Require(t, resource.Database)
	suite.Run(t, new(workItemLinkSuite))
}

func (s *workItemLinkSuite) TestCreate() {
	resetFn := s.DisableGormCallbacks()
	defer resetFn()

	s.T().Run("topology", func(t *testing.T) {
		topologies := []string{
			link.TopologyDependency,
			link.TopologyDirectedNetwork,
			link.TopologyNetwork,
			link.TopologyTree,
		}
		for _, topo := range topologies {
			t.Run(topo, func(t *testing.T) {
				t.Run("ok", func(t *testing.T) {
					// given
					fxt := tf.NewTestFixture(t, s.DB,
						tf.CreateWorkItemLinkEnvironment(),
						tf.WorkItems(2, tf.SetWorkItemTitles([]string{"bug1", "bug2"})),
						tf.WorkItemLinkTypes(1, tf.Topology(topo)),
					)
					createPayload := newCreateWorkItemLinkPayload(fxt.WorkItemByTitle("bug1").ID, fxt.WorkItemByTitle("bug2").ID, fxt.WorkItemLinkTypes[0].ID)
					// when
					res, workItemLink := test.CreateWorkItemLinkCreated(t, s.svc.Context, s.svc, s.workItemLinkCtrl, createPayload)
					// then
					require.NotNil(t, workItemLink)
					compareWithGoldenUUIDAgnostic(t, filepath.Join(s.testDir, "create", "ok.topology."+topo+".golden.json"), workItemLink)
					compareWithGoldenUUIDAgnostic(t, filepath.Join(s.testDir, "delete", "ok.topology."+topo+".headers.golden.json"), res.Header())
				})
			})
		}
		t.Run("network", func(t *testing.T) {
			t.Run("disallow creation of the same link", func(t *testing.T) {
				// Check if #586 is fixed.
				// given
				fxt := tf.NewTestFixture(t, s.DB,
					tf.CreateWorkItemLinkEnvironment(),
					tf.WorkItems(2, tf.SetWorkItemTitles([]string{"bug1", "bug2"})),
					tf.WorkItemLinkTypes(1, tf.TopologyNetwork()),
				)
				// The create payload is required during creation. Simply copy data over.
				payload := ConvertLinkFromModel(link.WorkItemLink{
					SourceID:   fxt.WorkItemByTitle("bug1").ID,
					TargetID:   fxt.WorkItemByTitle("bug2").ID,
					LinkTypeID: fxt.WorkItemLinkTypes[0].ID,
				})
				payload1 := app.CreateWorkItemLinkPayload{Data: payload.Data}
				payload2 := app.CreateWorkItemLinkPayload{Data: payload.Data}

				//createPayload1 := newCreateWorkItemLinkPayload(fxt.WorkItemByTitle("bug1").ID, fxt.WorkItemByTitle("bug1").ID, fxt.WorkItemLinkTypes[0].ID)
				_, workItemLink1 := test.CreateWorkItemLinkCreated(t, s.svc.Context, s.svc, s.workItemLinkCtrl, &payload1)
				require.NotNil(t, workItemLink1)
				// then
				res, jerrs := test.CreateWorkItemLinkConflict(t, s.svc.Context, s.svc, s.workItemLinkCtrl, &payload2)
				require.NotNil(t, jerrs)
				require.Len(t, jerrs.Errors, 1)
				ignoreString := "IGNORE_ME"
				jerrs.Errors[0].ID = &ignoreString
				compareWithGoldenUUIDAgnostic(t, filepath.Join(s.testDir, "create", "disallow_creation_of_the_same_link.error.golden.json"), jerrs)
				compareWithGoldenUUIDAgnostic(t, filepath.Join(s.testDir, "create", "disallow_creation_of_the_same_link.headers.golden.json"), res.Header())
			})
		})
	})
}

func (s *workItemLinkSuite) TestDelete() {
	resetFn := s.DisableGormCallbacks()
	defer resetFn()

	s.T().Run("ok", func(t *testing.T) {
		// given
		fxt := tf.NewTestFixture(t, s.DB, tf.WorkItemLinks(1))
		// when
		res := test.DeleteWorkItemLinkOK(t, s.svc.Context, s.svc, s.workItemLinkCtrl, fxt.WorkItemLinks[0].ID)
		// then
		compareWithGoldenUUIDAgnostic(t, filepath.Join(s.testDir, "delete", "ok.headers.golden.json"), res.Header())
	})
}

// Same for /api/workitems/:id/relationships/links
func (s *workItemLinkSuite) TestCreateAndDeleteWorkItemRelationshipsLink() {
	createPayload := newCreateWorkItemLinkPayload(s.bug1ID, s.bug2ID, s.bugBlockerLinkTypeID)
	_, workItemLink := test.CreateWorkItemRelationshipsLinksCreated(s.T(), s.svc.Context, s.svc, s.workItemRelsLinksCtrl, s.bug1ID, createPayload)
	require.NotNil(s.T(), workItemLink)
}

func (s *workItemLinkSuite) TestCreateWorkItemLinkBadRequestDueToInvalidLinkTypeID() {
	createPayload := newCreateWorkItemLinkPayload(s.bug1ID, s.bug2ID, uuid.Nil)
	_, _ = test.CreateWorkItemLinkBadRequest(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, createPayload)
}

// Same for /api/workitems/:id/relationships/links
func (s *workItemLinkSuite) TestCreateWorkItemRelationshipsLinksBadRequestDueToInvalidLinkTypeID() {
	createPayload := newCreateWorkItemLinkPayload(s.bug1ID, s.bug2ID, uuid.Nil)
	_, _ = test.CreateWorkItemRelationshipsLinksBadRequest(s.T(), s.svc.Context, s.svc, s.workItemRelsLinksCtrl, s.bug1ID, createPayload)
}

func (s *workItemLinkSuite) TestCreateWorkItemLinkBadRequestDueToNotFoundLinkType() {
	createPayload := newCreateWorkItemLinkPayload(s.bug1ID, s.bug2ID, uuid.FromStringOrNil("11122233-871b-43a6-9166-0c4bd573e333"))
	_, _ = test.CreateWorkItemLinkNotFound(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, createPayload)
}

// Same for /api/workitems/:id/relationships/links
func (s *workItemLinkSuite) TestCreateWorkItemRelationshipLinksBadRequestDueToNotFoundLinkType() {
	createPayload := newCreateWorkItemLinkPayload(s.bug1ID, s.bug2ID, uuid.FromStringOrNil("11122233-871b-43a6-9166-0c4bd573e333"))
	_, _ = test.CreateWorkItemRelationshipsLinksNotFound(s.T(), s.svc.Context, s.svc, s.workItemRelsLinksCtrl, s.bug1ID, createPayload)
}

func (s *workItemLinkSuite) TestCreateWorkItemLinkBadRequestDueToNotFoundSource() {
	createPayload := newCreateWorkItemLinkPayload(uuid.NewV4(), s.bug2ID, s.bugBlockerLinkTypeID)
	_, _ = test.CreateWorkItemLinkNotFound(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, createPayload)
}

// Same for /api/workitems/:id/relationships/links
func (s *workItemLinkSuite) TestCreateWorkItemRelationshipsLinksBadRequestDueToNotFoundSource() {
	createPayload := newCreateWorkItemLinkPayload(uuid.NewV4(), s.bug2ID, s.bugBlockerLinkTypeID)
	_, _ = test.CreateWorkItemRelationshipsLinksBadRequest(s.T(), s.svc.Context, s.svc, s.workItemRelsLinksCtrl, s.bug2ID, createPayload)
}

func (s *workItemLinkSuite) TestCreateWorkItemLinkBadRequestDueToNotFoundTarget() {
	createPayload := newCreateWorkItemLinkPayload(s.bug1ID, uuid.NewV4(), s.bugBlockerLinkTypeID)
	_, _ = test.CreateWorkItemLinkNotFound(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, createPayload)
}

// Same for /api/workitems/:id/relationships/links
func (s *workItemLinkSuite) TestCreateWorkItemRelationshipsLinksBadRequestDueToNotFoundTarget() {
	createPayload := newCreateWorkItemLinkPayload(s.bug1ID, uuid.NewV4(), s.bugBlockerLinkTypeID)
	_, _ = test.CreateWorkItemRelationshipsLinksNotFound(s.T(), s.svc.Context, s.svc, s.workItemRelsLinksCtrl, s.bug1ID, createPayload)
}

func (s *workItemLinkSuite) TestDeleteWorkItemLinkNotFound() {
	test.DeleteWorkItemLinkNotFound(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, uuid.FromStringOrNil("1e9a8b53-73a6-40de-b028-5177add79ffa"))
}

func (s *workItemLinkSuite) TestUpdateWorkItemLinkNotFound() {
	createPayload := newCreateWorkItemLinkPayload(s.bug1ID, s.bug2ID, s.userLinkCategoryID)
	notExistingId := uuid.FromStringOrNil("46bbce9c-8219-4364-a450-dfd1b501654e")
	createPayload.Data.ID = &notExistingId
	// Wrap data portion in an update payload instead of a create payload
	updateLinkPayload := &app.UpdateWorkItemLinkPayload{
		Data: createPayload.Data,
	}
	test.UpdateWorkItemLinkNotFound(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, *updateLinkPayload.Data.ID, updateLinkPayload)
}

func (s *workItemLinkSuite) TestUpdateWorkItemLinkOK() {
	// given
	createPayload := newCreateWorkItemLinkPayload(s.bug1ID, s.bug2ID, s.bugBlockerLinkTypeID)
	_, workItemLink := test.CreateWorkItemLinkCreated(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, createPayload)
	require.NotNil(s.T(), workItemLink)
	// Specify new description for link type that we just created
	// Wrap data portion in an update payload instead of a create payload
	updateLinkPayload := newUpdateWorkItemLinkPayload(*workItemLink.Data.ID, s.bug1ID, s.bug3ID, s.bugBlockerLinkTypeID)
	// when
	_, l := test.UpdateWorkItemLinkOK(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, *updateLinkPayload.Data.ID, updateLinkPayload)
	// then
	require.NotNil(s.T(), l.Data)
	require.Equal(s.T(), workItemLink.Data.Attributes.CreatedAt.UTC(), l.Data.Attributes.CreatedAt.UTC())
	require.NotNil(s.T(), l.Data.Attributes.CreatedAt)
	require.NotNil(s.T(), l.Data.Relationships)
	require.NotNil(s.T(), l.Data.Relationships.Target.Data)
	assert.Equal(s.T(), s.bug3ID, l.Data.Relationships.Target.Data.ID)
}

func (s *workItemLinkSuite) TestUpdateWorkItemLinkVersionConflict() {
	// given
	createPayload := newCreateWorkItemLinkPayload(s.bug1ID, s.bug2ID, s.bugBlockerLinkTypeID)
	_, workItemLink := test.CreateWorkItemLinkCreated(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, createPayload)
	require.NotNil(s.T(), workItemLink)
	// Specify new description for link type that we just created
	// Wrap data portion in an update payload instead of a create payload
	updateLinkPayload := &app.UpdateWorkItemLinkPayload{
		Data: workItemLink.Data,
	}
	updateLinkPayload.Data.Relationships.Target.Data.ID = s.bug3ID
	// force a different version of the entity
	previousVersion := *updateLinkPayload.Data.Attributes.Version - 1
	updateLinkPayload.Data.Attributes.Version = &previousVersion
	// when/then
	test.UpdateWorkItemLinkConflict(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, *updateLinkPayload.Data.ID, updateLinkPayload)
	// then
}

func (s *workItemLinkSuite) newCreateWorkItemLinkPayload() *app.WorkItemLinkSingle {
	createPayload := newCreateWorkItemLinkPayload(s.bug1ID, s.bug2ID, s.bugBlockerLinkTypeID)
	_, workItemLink := test.CreateWorkItemLinkCreated(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, createPayload)
	require.NotNil(s.T(), workItemLink)
	require.NotNil(s.T(), workItemLink.Data.Attributes.UpdatedAt)
	return workItemLink
}

func assertWorkItemLink(t *testing.T, expectedWorkItemLink *app.WorkItemLinkSingle, wiLink *app.WorkItemLinkSingle) {
	require.NotNil(t, wiLink)
	expected, err := ConvertLinkToModel(*expectedWorkItemLink)
	require.Nil(t, err)
	// Convert to model space and use equal function
	actual, err := ConvertLinkToModel(*wiLink)
	require.Nil(t, err)
	require.Equal(t, *expected, *actual)
	require.NotNil(t, wiLink.Data.Links, "The link MUST include a self link")
	require.NotEmpty(t, wiLink.Data.Links.Self, "The link MUST include a self link that's not empty")
}

// TestShowWorkItemLinkOK tests if we can fetch the "system" work item link
func (s *workItemLinkSuite) TestShowWorkItemLinkOK() {
	// given
	createdWorkItemLink := s.newCreateWorkItemLinkPayload()
	// when
	_, retrievedWorkItemLink := test.ShowWorkItemLinkOK(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, *createdWorkItemLink.Data.ID, nil, nil)
	// then
	assertWorkItemLink(s.T(), createdWorkItemLink, retrievedWorkItemLink)
}

// TestShowWorkItemLinkOKUsingExpiredIfModifiedSinceHeader
func (s *workItemLinkSuite) TestShowWorkItemLinkOKUsingExpiredIfModifiedSinceHeader() {
	// given
	createdWorkItemLink := s.newCreateWorkItemLinkPayload()
	// when
	ifModifiedSince := app.ToHTTPTime(createdWorkItemLink.Data.Attributes.UpdatedAt.Add(-1 * time.Hour))
	res, retrievedWorkItemLink := test.ShowWorkItemLinkOK(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, *createdWorkItemLink.Data.ID, &ifModifiedSince, nil)
	// then
	assertWorkItemLink(s.T(), createdWorkItemLink, retrievedWorkItemLink)
	assertResponseHeaders(s.T(), res)
}

// TestShowWorkItemLinkOKUsingExpiredIfModifiedSinceHeader
func (s *workItemLinkSuite) TestShowWorkItemLinkOKUsingExpiredIfNoneMatchHeader() {
	// given
	createdWorkItemLink := s.newCreateWorkItemLinkPayload()
	// when
	ifNoneMatch := "foo"
	res, retrievedWorkItemLink := test.ShowWorkItemLinkOK(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, *createdWorkItemLink.Data.ID, nil, &ifNoneMatch)
	// then
	assertWorkItemLink(s.T(), createdWorkItemLink, retrievedWorkItemLink)
	assertResponseHeaders(s.T(), res)
}

// TestShowWorkItemLinkOKUsingExpiredIfModifiedSinceHeader
func (s *workItemLinkSuite) TestShowWorkItemLinkNotModifiedUsingIfModifiedSinceHeader() {
	// given
	createdWorkItemLink := s.newCreateWorkItemLinkPayload()
	// when
	ifModifiedSince := app.ToHTTPTime(*createdWorkItemLink.Data.Attributes.UpdatedAt)
	res := test.ShowWorkItemLinkNotModified(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, *createdWorkItemLink.Data.ID, &ifModifiedSince, nil)
	// then
	assertResponseHeaders(s.T(), res)
}

// TestShowWorkItemLinkOKUsingExpiredIfModifiedSinceHeader
func (s *workItemLinkSuite) TestShowWorkItemLinkNotModifiedUsingIfNoneMatchHeader() {
	// given
	createdWorkItemLink := s.newCreateWorkItemLinkPayload()
	// when
	modelWorkItemLink, err := ConvertLinkToModel(*createdWorkItemLink)
	require.Nil(s.T(), err)
	ifNoneMatch := app.GenerateEntityTag(modelWorkItemLink)
	res := test.ShowWorkItemLinkNotModified(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, *createdWorkItemLink.Data.ID, nil, &ifNoneMatch)
	// then
	assertResponseHeaders(s.T(), res)
}

// TestShowWorkItemLinkNotFound tests if we can fetch a non existing work item link
func (s *workItemLinkSuite) TestShowWorkItemLinkNotFound() {
	test.ShowWorkItemLinkNotFound(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, uuid.FromStringOrNil("88727441-4a21-4b35-aabe-007f8273cd19"), nil, nil)
}

func (s *workItemLinkSuite) createSomeLinks() (*app.WorkItemLinkSingle, *app.WorkItemLinkSingle) {
	createPayload1 := newCreateWorkItemLinkPayload(s.bug1ID, s.bug2ID, s.bugBlockerLinkTypeID)
	_, workItemLink1 := test.CreateWorkItemLinkCreated(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, createPayload1)
	require.NotNil(s.T(), workItemLink1)
	_, err := ConvertLinkToModel(*workItemLink1)
	require.Nil(s.T(), err)

	createPayload2 := newCreateWorkItemLinkPayload(s.bug2ID, s.bug3ID, s.bugBlockerLinkTypeID)
	_, workItemLink2 := test.CreateWorkItemLinkCreated(s.T(), s.svc.Context, s.svc, s.workItemLinkCtrl, createPayload2)
	require.NotNil(s.T(), workItemLink2)
	_, err = ConvertLinkToModel(*workItemLink2)
	require.Nil(s.T(), err)

	return workItemLink1, workItemLink2
}

// validateSomeLinks validates that workItemLink1 and workItemLink2 are in the
// linkCollection and that all resources are included
func (s *workItemLinkSuite) validateSomeLinks(linkCollection *app.WorkItemLinkList, workItemLink1, workItemLink2 *app.WorkItemLinkSingle) {
	require.NotNil(s.T(), linkCollection)
	require.Nil(s.T(), linkCollection.Validate())
	// Check the number of found work item links
	require.NotNil(s.T(), linkCollection.Data)
	require.Condition(s.T(), func() bool {
		return (len(linkCollection.Data) >= 2)
	}, "At least two work item links must exist (%s and %s), but only %d exist.", *workItemLink1.Data.ID, *workItemLink2.Data.ID, len(linkCollection.Data))
	// Search for the work item types that must exist at minimum
	toBeFound := 2
	for i := 0; i < len(linkCollection.Data) && toBeFound > 0; i++ {
		actualLink := *linkCollection.Data[i]
		var expectedLink *app.WorkItemLinkData

		switch *actualLink.ID {
		case *workItemLink1.Data.ID:
			expectedLink = workItemLink1.Data
		case *workItemLink2.Data.ID:
			expectedLink = workItemLink2.Data
		}

		if expectedLink != nil {
			s.T().Log("Found work item link in collection: ", *expectedLink.ID)
			toBeFound--

			// Check JSONAPI "type"" field (should be "workitemlinks")
			require.Equal(s.T(), expectedLink.Type, actualLink.Type)

			// Check work item link type
			require.Equal(s.T(), expectedLink.Relationships.LinkType.Data.ID, actualLink.Relationships.LinkType.Data.ID)
			require.Equal(s.T(), expectedLink.Relationships.LinkType.Data.Type, actualLink.Relationships.LinkType.Data.Type)

			// Check source type
			require.Equal(s.T(), expectedLink.Relationships.Source.Data.ID, actualLink.Relationships.Source.Data.ID, "Wrong source ID for the link")
			require.Equal(s.T(), expectedLink.Relationships.Source.Data.Type, actualLink.Relationships.Source.Data.Type, "Wrong source JSONAPI type for the link")

			// Check target type
			require.Equal(s.T(), expectedLink.Relationships.Target.Data.ID, actualLink.Relationships.Target.Data.ID, "Wrong target ID for the link")
			require.Equal(s.T(), expectedLink.Relationships.Target.Data.Type, actualLink.Relationships.Target.Data.Type, "Wrong target JSONAPI type for the link")
		}
	}
	require.Exactly(s.T(), 0, toBeFound, "Not all required work item links (%s and %s) where found.", *workItemLink1.Data.ID, *workItemLink2.Data.ID)

	toBeFound = 5 // 1 x link category, 1 x link type, 3 x work items
	for i := 0; i < len(linkCollection.Included) && toBeFound > 0; i++ {
		switch v := linkCollection.Included[i].(type) {
		case *app.WorkItemLinkCategoryData:
			if *v.ID == s.userLinkCategoryID {
				s.T().Log("Found work item link category in \"included\" element: ", *v.ID)
				toBeFound--
			}
		case *app.Space:
			if *v.ID == s.userSpaceID {
				s.T().Log("Found work item link space in \"included\" element: ", *v.ID)
				toBeFound--
			}
		case *app.WorkItemLinkTypeData:
			if *v.ID == s.bugBlockerLinkTypeID {
				s.T().Log("Found work item link type in \"included\" element: ", *v.ID)
				toBeFound--
			}
		case *app.WorkItem:
			wid := *v.ID
			if wid == s.bug1ID || wid == s.bug2ID || wid == s.bug3ID {
				s.T().Log("Found work item in \"included\" element: ", *v.ID)
				toBeFound--
			}
		// TODO(kwk): Check for WITs (once #559 is merged)
		// case *app.WorkItemTypeData:
		default:
			s.T().Errorf("Object of unknown type included in work item link list response: %T", linkCollection.Included[i])
		}
	}
	require.Exactly(s.T(), 0, toBeFound, "Not all required included elements where found.")
}

// Same as TestListWorkItemLinkOK, for /api/workitems/:id/relationships/links
func (s *workItemLinkSuite) TestListWorkItemRelationshipsLinksOK() {
	link1, link2 := s.createSomeLinks()
	_, linkCollection := test.ListWorkItemRelationshipsLinksOK(s.T(), s.svc.Context, s.svc, s.workItemRelsLinksCtrl, s.bug2ID, nil, nil)
	s.validateSomeLinks(linkCollection, link1, link2)
}

func (s *workItemLinkSuite) TestListWorkItemRelationshipsLinksNotFound() {
	_, _ = test.ListWorkItemRelationshipsLinksNotFound(s.T(), s.svc.Context, s.svc, s.workItemRelsLinksCtrl, uuid.NewV4(), nil, nil)
}

func (s *workItemLinkSuite) getWorkItemLinkTestDataFunc() func(t *testing.T) []testSecureAPI {
	return func(t *testing.T) []testSecureAPI {
		privatekey := testtoken.PrivateKey()
		differentPrivatekey, err := jwt.ParseRSAPrivateKeyFromPEM(([]byte(RSADifferentPrivateKeyTest)))
		require.Nil(t, err, "Could not parse private key")
		createWorkItemLinkPayloadString := bytes.NewBuffer([]byte(`
		{
			"data": {
				"attributes": {
					"version": 0
				},
				"id": "40bbdd3d-8b5d-4fd6-ac90-7236b669af04",
				"relationships": {
					"link_type": {
						"data": {
						"id": "6c5610be-30b2-4880-9fec-81e4f8e4fd76",
						"type": "workitemlinktypes"
						}
					},
					"source": {
						"data": {
						"id": "1234",
						"type": "workitems"
						}
					},
					"target": {
						"data": {
						"id": "1234",
						"type": "workitems"
						}
					}
				},
				"type": "workitemlinks"
			}
		}
  		`))

		testWorkItemLinksAPI := []testSecureAPI{
			// Create Work Item API with different parameters
			{
				method:             http.MethodPost,
				url:                endpointWorkItemLinks,
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            createWorkItemLinkPayloadString,
				jwtToken:           getExpiredAuthHeader(t, privatekey),
			}, {
				method:             http.MethodPost,
				url:                endpointWorkItemLinks,
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            createWorkItemLinkPayloadString,
				jwtToken:           getMalformedAuthHeader(t, privatekey),
			}, {
				method:             http.MethodPost,
				url:                endpointWorkItemLinks,
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            createWorkItemLinkPayloadString,
				jwtToken:           getValidAuthHeader(t, differentPrivatekey),
			}, {
				method:             http.MethodPost,
				url:                endpointWorkItemLinks,
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            createWorkItemLinkPayloadString,
				jwtToken:           "",
			},
			// Update Work Item API with different parameters
			{
				method:             http.MethodPatch,
				url:                endpointWorkItemLinks + "/6c5610be-30b2-4880-9fec-81e4f8e4fd76",
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            createWorkItemLinkPayloadString,
				jwtToken:           getExpiredAuthHeader(t, privatekey),
			}, {
				method:             http.MethodPatch,
				url:                endpointWorkItemLinks + "/6c5610be-30b2-4880-9fec-81e4f8e4fd76",
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            createWorkItemLinkPayloadString,
				jwtToken:           getMalformedAuthHeader(t, privatekey),
			}, {
				method:             http.MethodPatch,
				url:                endpointWorkItemLinks + "/6c5610be-30b2-4880-9fec-81e4f8e4fd76",
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            createWorkItemLinkPayloadString,
				jwtToken:           getValidAuthHeader(t, differentPrivatekey),
			}, {
				method:             http.MethodPatch,
				url:                endpointWorkItemLinks + "/6c5610be-30b2-4880-9fec-81e4f8e4fd76",
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            createWorkItemLinkPayloadString,
				jwtToken:           "",
			},
			// Delete Work Item API with different parameters
			{
				method:             http.MethodDelete,
				url:                endpointWorkItemLinks + "/6c5610be-30b2-4880-9fec-81e4f8e4fd76",
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            nil,
				jwtToken:           getExpiredAuthHeader(t, privatekey),
			}, {
				method:             http.MethodDelete,
				url:                endpointWorkItemLinks + "/6c5610be-30b2-4880-9fec-81e4f8e4fd76",
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            nil,
				jwtToken:           getMalformedAuthHeader(t, privatekey),
			}, {
				method:             http.MethodDelete,
				url:                endpointWorkItemLinks + "/6c5610be-30b2-4880-9fec-81e4f8e4fd76",
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            nil,
				jwtToken:           getValidAuthHeader(t, differentPrivatekey),
			}, {
				method:             http.MethodDelete,
				url:                endpointWorkItemLinks + "/6c5610be-30b2-4880-9fec-81e4f8e4fd76",
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            nil,
				jwtToken:           "",
			},
			// Try fetching a random work item link
			// We do not have security on GET hence this should return 404 not found
			{
				method:             http.MethodGet,
				url:                endpointWorkItemLinks + "/fc591f38-a805-4abd-bfce-2460e49d8cc4",
				expectedStatusCode: http.StatusNotFound,
				expectedErrorCode:  jsonapi.ErrorCodeNotFound,
				payload:            nil,
				jwtToken:           "",
			},
		}
		return testWorkItemLinksAPI
	}
}

// This test case will check authorized access to Create/Update/Delete APIs
func (s *workItemLinkSuite) TestUnauthorizeWorkItemLinkCUD() {
	UnauthorizeCreateUpdateDeleteTest(s.T(), s.getWorkItemLinkTestDataFunc(), func() *goa.Service {
		return goa.New("TestUnauthorizedCreateWorkItemLink-Service")
	}, func(service *goa.Service) error {
		controller := NewWorkItemLinkController(service, gormapplication.NewGormDB(s.DB), s.Configuration)
		app.MountWorkItemLinkController(service, controller)
		return nil
	})
}

// The work item ID will be used to construct /api/workitems/:id/relationships/links endpoints
func (s *workItemLinkSuite) getWorkItemRelationshipLinksTestData(spaceID, wiID uuid.UUID) func(t *testing.T) []testSecureAPI {
	return func(t *testing.T) []testSecureAPI {
		privatekey := testtoken.PrivateKey()
		differentPrivatekey, err := jwt.ParseRSAPrivateKeyFromPEM(([]byte(RSADifferentPrivateKeyTest)))
		if err != nil {
			t.Fatal("Could not parse different private key ", err)
		}

		createWorkItemLinkPayloadString := bytes.NewBuffer([]byte(`
		{
			"data": {
				"attributes": {
					"version": 0
				},
				"id": "40bbdd3d-8b5d-4fd6-ac90-7236b669af04",
				"relationships": {
					"link_type": {
						"data": {
						"id": "6c5610be-30b2-4880-9fec-81e4f8e4fd76",
						"type": "workitemlinktypes"
						}
					},
					"source": {
						"data": {
						"id": "1234",
						"type": "workitems"
						}
					},
					"target": {
						"data": {
						"id": "1234",
						"type": "workitems"
						}
					}
				},
				"type": "workitemlinks"
			}
		}
  		`))

		relationshipsEndpoint := fmt.Sprintf(endpointWorkItemRelationshipsLinks, wiID)
		testWorkItemLinksAPI := []testSecureAPI{
			// Create Work Item API with different parameters
			{
				method:             http.MethodPost,
				url:                relationshipsEndpoint,
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            createWorkItemLinkPayloadString,
				jwtToken:           getExpiredAuthHeader(t, privatekey),
			}, {
				method:             http.MethodPost,
				url:                relationshipsEndpoint,
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            createWorkItemLinkPayloadString,
				jwtToken:           getMalformedAuthHeader(t, privatekey),
			}, {
				method:             http.MethodPost,
				url:                relationshipsEndpoint,
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            createWorkItemLinkPayloadString,
				jwtToken:           getValidAuthHeader(t, differentPrivatekey),
			}, {
				method:             http.MethodPost,
				url:                relationshipsEndpoint,
				expectedStatusCode: http.StatusUnauthorized,
				expectedErrorCode:  jsonapi.ErrorCodeJWTSecurityError,
				payload:            createWorkItemLinkPayloadString,
				jwtToken:           "",
			},
		}
		return testWorkItemLinksAPI
	}
}

func (s *workItemLinkSuite) TestUnauthorizeWorkItemRelationshipsLinksCUD() {
	UnauthorizeCreateUpdateDeleteTest(s.T(), s.getWorkItemRelationshipLinksTestData(space.SystemSpace, s.bug1ID), func() *goa.Service {
		return goa.New("TestUnauthorizedCreateWorkItemRelationshipsLinks-Service")
	}, func(service *goa.Service) error {
		controller := NewWorkItemRelationshipsLinksController(service, gormapplication.NewGormDB(s.DB), s.Configuration)
		app.MountWorkItemRelationshipsLinksController(service, controller)
		return nil
	})
}
