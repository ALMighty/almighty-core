package gormsupport_test

import (
	"testing"

	"github.com/fabric8-services/fabric8-wit/convert"
	"github.com/fabric8-services/fabric8-wit/gormsupport"
	"github.com/fabric8-services/fabric8-wit/gormtestsupport"
	"github.com/fabric8-services/fabric8-wit/resource"
	"github.com/fabric8-services/fabric8-wit/workitem/link"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestVersioning_Equal(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)

	a := gormsupport.Versioning{
		Version: 42,
	}

	t.Run("equality", func(t *testing.T) {
		b := gormsupport.Versioning{
			Version: 42,
		}
		require.True(t, a.Equal(b))
	})
	t.Run("type difference", func(t *testing.T) {
		b := convert.DummyEqualer{}
		require.False(t, a.Equal(b))
	})
	t.Run("version difference", func(t *testing.T) {
		b := gormsupport.Versioning{
			Version: 123,
		}
		require.False(t, a.Equal(b))
	})
}

type VersioningSuite struct {
	gormtestsupport.DBTestSuite
}

func TestVersioning(t *testing.T) {
	suite.Run(t, &VersioningSuite{DBTestSuite: gormtestsupport.NewDBTestSuite()})
}

func (s *VersioningSuite) TestCallbacks() {
	// given a work item link category that embeds the Versioning struct is a
	// good way to demonstrate the way the callbacks work.
	cat := link.WorkItemLinkCategory{
		Name: "foo",
	}
	// set the version to something else than 0 the BeforeCreate callback will
	// automatically set this to 0 before entering it in the DB.
	cat.Version = 123

	s.T().Run("before create", func(t *testing.T) {
		// when
		err := s.DB.Create(&cat).Error
		// then
		require.NoError(t, err)
		require.Equal(t, 0, cat.Version, "initial version of entity must be 0 nomatter what the given version was")
	})
	s.T().Run("before update", func(t *testing.T) {
		// given
		newName := "new name"
		cat.Name = newName
		// when
		err := s.DB.Save(&cat).Error
		// then
		require.NoError(t, err)
		require.Equal(t, 1, cat.Version, "followup version of entity must be 1")
		require.Equal(t, newName, cat.Name)
	})
}