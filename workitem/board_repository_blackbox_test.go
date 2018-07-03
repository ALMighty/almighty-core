package workitem_test

import (
	"testing"
	//"time"

	"github.com/fabric8-services/fabric8-wit/convert"
	"github.com/fabric8-services/fabric8-wit/errors"
	"github.com/fabric8-services/fabric8-wit/gormtestsupport"
	"github.com/fabric8-services/fabric8-wit/resource"
	tf "github.com/fabric8-services/fabric8-wit/test/testfixture"
	"github.com/fabric8-services/fabric8-wit/workitem"
	errs "github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type workItemBoardRepoTest struct {
	gormtestsupport.DBTestSuite
	repo workitem.BoardRepository
}

func TestWorkItemBoardRepository(t *testing.T) {
	suite.Run(t, &workItemBoardRepoTest{DBTestSuite: gormtestsupport.NewDBTestSuite()})
}

func (s *workItemBoardRepoTest) SetupTest() {
	s.DBTestSuite.SetupTest()
	s.repo = workitem.NewBoardRepository(s.DB)
}

func (s *workItemBoardRepoTest) TestExists() {
	s.T().Run("board exists", func(t *testing.T) {
		// given
		fxt := tf.NewTestFixture(t, s.DB, tf.WorkItemBoards(1))
		// when
		err := s.repo.CheckExists(s.Ctx, fxt.WorkItemBoards[0].ID)
		// then
		require.NoError(s.T(), err)
	})

	s.T().Run("group doesn't exist", func(t *testing.T) {
		// given
		nonExistingWorkItemBoardID := uuid.NewV4()
		// when
		err := s.repo.CheckExists(s.Ctx, nonExistingWorkItemBoardID)
		// then
		require.IsType(t, errors.NotFoundError{}, err)
	})
}

func (s *workItemBoardRepoTest) TestCreate() {
	// given
	fxt := tf.NewTestFixture(s.T(), s.DB, tf.WorkItemBoards(3))
	ID := uuid.NewV4()
	expected := workitem.Board{
		ID:              ID,
		SpaceTemplateID: fxt.SpaceTemplates[0].ID,
		Name:            "Some Board Name",
		Description:     "Some Board Description ",
		ContextType:     "TypeLevelContext",
		Context:         uuid.NewV4().String(),
		Columns: []workitem.BoardColumn{
			{
				ID:                uuid.NewV4(),
				Name:              "New",
				ColumnOrder:       0,
				TransRuleKey:      "updateStateFromColumnMove",
				TransRuleArgument: "{ 'metastate': 'mNew' }",
				BoardID:           ID,
			},
			{
				ID:                uuid.NewV4(),
				Name:              "Done",
				ColumnOrder:       1,
				TransRuleKey:      "updateStateFromColumnMove",
				TransRuleArgument: "{ 'metastate': 'mDone' }",
				BoardID:           ID,
			},
		},
	}

	s.T().Run("ok", func(t *testing.T) {
		actual, err := s.repo.Create(s.Ctx, expected)
		require.NoError(t, err)
		require.True(t, expected.Equal(*actual))
		require.True(t, expected.Columns[0].Equal(actual.Columns[0]))
		require.True(t, expected.Columns[1].Equal(actual.Columns[1]))
		t.Run("load same work item board and check it is the same", func(t *testing.T) {
			actual, err := s.repo.Load(s.Ctx, ID)
			require.NoError(t, err)
			require.True(t, expected.Equal(*actual))
			require.True(t, expected.Columns[0].Equal(actual.Columns[0]))
			require.True(t, expected.Columns[1].Equal(actual.Columns[1]))
		})
	})
	s.T().Run("invalid", func(t *testing.T) {
		t.Run("unknown space template", func(t *testing.T) {
			g := expected
			g.ID = uuid.NewV4()
			g.SpaceTemplateID = uuid.NewV4()
			_, err := s.repo.Create(s.Ctx, g)
			require.Error(t, err)
		})
	})
}

func (s *workItemBoardRepoTest) TestLoad() {
	s.T().Run("board exists", func(t *testing.T) {
		// given
		fxt := tf.NewTestFixture(t, s.DB, tf.WorkItemBoards(1))
		// when
		actual, err := s.repo.Load(s.Ctx, fxt.WorkItemBoards[0].ID)
		require.NoError(t, err)
		require.True(t, fxt.WorkItemBoards[0].Equal(*actual))
	})
	s.T().Run("board doesn't exist", func(t *testing.T) {
		// when
		_, err := s.repo.Load(s.Ctx, uuid.NewV4())
		// then
		require.Error(t, err)
	})
}

func (s *workItemBoardRepoTest) TestList() {
	s.T().Run("ok", func(t *testing.T) {
		// given
		fxt := tf.NewTestFixture(t, s.DB, tf.WorkItemBoards(3))
		// when
		actual, err := s.repo.List(s.Ctx, fxt.SpaceTemplates[0].ID)
		// then
		require.NoError(t, err)
		require.Len(t, actual, len(fxt.WorkItemBoards))
		for idx := range fxt.WorkItemBoards {
			require.True(t, fxt.WorkItemBoards[idx].Equal(*actual[idx]))
		}
	})
	s.T().Run("space template not found", func(t *testing.T) {
		// when
		groups, err := s.repo.List(s.Ctx, uuid.NewV4())
		// then
		require.Error(t, err)
		require.IsType(t, errors.NotFoundError{}, errs.Cause(err))
		require.Empty(t, groups)
	})
}

func TestWorkItemBoard_Equal(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)
	// given
	ID := uuid.NewV4()
	a := workitem.Board{
		ID:              ID,
		SpaceTemplateID: uuid.NewV4(),
		Name:            "Some Board Name",
		Description:     "Some Board Description ",
		ContextType:     "TypeLevelContext",
		Context:         uuid.NewV4().String(),
		Columns: []workitem.BoardColumn{
			{
				ID:                uuid.NewV4(),
				Name:              "New",
				ColumnOrder:       0,
				TransRuleKey:      "updateStateFromColumnMove",
				TransRuleArgument: "{ 'metastate': 'mNew' }",
				BoardID:           ID,
			},
			{
				ID:                uuid.NewV4(),
				Name:              "Done",
				ColumnOrder:       1,
				TransRuleKey:      "updateStateFromColumnMove",
				TransRuleArgument: "{ 'metastate': 'mDone' }",
				BoardID:           ID,
			},
		},
	}
	t.Run("equality", func(t *testing.T) {
		t.Parallel()
		b := a
		assert.True(t, a.Equal(b))
	})
	t.Run("types", func(t *testing.T) {
		t.Parallel()
		b := convert.DummyEqualer{}
		assert.False(t, a.Equal(b))
	})
	t.Run("Name", func(t *testing.T) {
		t.Parallel()
		b := a
		b.Name = "bar"
		assert.False(t, a.Equal(b))
	})
	t.Run("SpaceTemplateID", func(t *testing.T) {
		t.Parallel()
		b := a
		b.SpaceTemplateID = uuid.NewV4()
		assert.False(t, a.Equal(b))
	})
	t.Run("Columns", func(t *testing.T) {
		t.Parallel()
		b := a
		// different IDs
		b.Columns = []workitem.BoardColumn{
			{
				ID:                uuid.NewV4(),
				Name:              "New",
				ColumnOrder:       0,
				TransRuleKey:      "updateStateFromColumnMove",
				TransRuleArgument: "{ 'metastate': 'mNew' }",
				BoardID:           ID,
			},
			{
				ID:                uuid.NewV4(),
				Name:              "Done",
				ColumnOrder:       1,
				TransRuleKey:      "updateStateFromColumnMove",
				TransRuleArgument: "{ 'metastate': 'mDone' }",
				BoardID:           ID,
			},
		}
		assert.False(t, a.Equal(b))
		// different length
		b.Columns = []workitem.BoardColumn{
			{
				ID:                uuid.NewV4(),
				Name:              "New",
				ColumnOrder:       0,
				TransRuleKey:      "updateStateFromColumnMove",
				TransRuleArgument: "{ 'metastate': 'mNew' }",
				BoardID:           ID,
			},
		}
		assert.False(t, a.Equal(b))
	})
}
