package search

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/fabric8-services/fabric8-wit/criteria"
	"github.com/fabric8-services/fabric8-wit/resource"
	"github.com/fabric8-services/fabric8-wit/workitem"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// While registering URLs do not include protocol because it will be removed before scanning starts
	// Please do not include trailing slashes because it will be removed before scanning starts
	RegisterAsKnownURL("test-work-item-list-details", `(?P<domain>demo.openshift.io)(?P<path>/work-item/list/detail/)(?P<id>\d*)`)
	RegisterAsKnownURL("test-work-item-board-details", `(?P<domain>demo.openshift.io)(?P<path>/work-item/board/detail/)(?P<id>\d*)`)
}

func TestGenerateSQLSearchStringText(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)
	input := searchKeyword{
		number: []string{"10", "99"},
		words:  []string{"username", "title_substr", "desc_substr"},
	}
	expectedSQLParameter := "10 & 99 & username & title_substr & desc_substr"

	actualSQLParameter := generateSQLSearchInfo(input)
	assert.Equal(t, expectedSQLParameter, actualSQLParameter)
}

func TestGenerateSQLSearchStringIdOnly(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)
	input := searchKeyword{
		number: []string{"10"},
		words:  []string{},
	}
	expectedSQLParameter := "10"

	actualSQLParameter := generateSQLSearchInfo(input)
	assert.Equal(t, expectedSQLParameter, actualSQLParameter)
}

func TestParseSearchString(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)
	input := "user input for search string with some ids like number:99 and number:400 but this is not id like 800"
	op, _ := parseSearchString(context.Background(), input)
	expectedSearchRes := searchKeyword{
		number: []string{"99:*A", "400:*A"},
		words:  []string{"user:*", "input:*", "for:*", "search:*", "string:*", "with:*", "some:*", "ids:*", "like:*", "and:*", "but:*", "this:*", "is:*", "not:*", "id:*", "like:*", "800:*"},
	}
	t.Log("Parsed search string: ", op)
	assert.True(t, assert.ObjectsAreEqualValues(expectedSearchRes, op))
}

type searchTestData struct {
	query    string
	expected searchKeyword
}

func TestParseSearchStringURL(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)
	inputSet := []searchTestData{{
		query: "http://demo.openshift.io/work-item/list/detail/100",
		expected: searchKeyword{
			number: nil,
			words:  []string{"(100:*A | demo.openshift.io/work-item/list/detail/100:*)"},
		},
	}, {
		query: "http://demo.openshift.io/work-item/board/detail/100",
		expected: searchKeyword{
			number: nil,
			words:  []string{"(100:*A | demo.openshift.io/work-item/board/detail/100:*)"},
		},
	}}

	for _, input := range inputSet {
		op, _ := parseSearchString(context.Background(), input.query)
		assert.Equal(t, input.expected, op)
	}
}

func TestParseSearchStringURLWithouID(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)
	inputSet := []searchTestData{{
		query: "http://demo.openshift.io/work-item/list/detail/",
		expected: searchKeyword{
			number: nil,
			words:  []string{"demo.openshift.io/work-item/list/detail:*"},
		},
	}, {
		query: "http://demo.openshift.io/work-item/board/detail/",
		expected: searchKeyword{
			number: nil,
			words:  []string{"demo.openshift.io/work-item/board/detail:*"},
		},
	}}

	for _, input := range inputSet {
		op, _ := parseSearchString(context.Background(), input.query)
		assert.True(t, assert.ObjectsAreEqualValues(input.expected, op))
	}

}

func TestParseSearchStringDifferentURL(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)
	input := "http://demo.redhat.io"
	op, _ := parseSearchString(context.Background(), input)
	expectedSearchRes := searchKeyword{
		number: nil,
		words:  []string{"demo.redhat.io:*"},
	}
	assert.True(t, assert.ObjectsAreEqualValues(expectedSearchRes, op))
}

func TestParseSearchStringCombination(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)
	// do combination of ID, full text and URLs
	// check if it works as expected.
	input := "http://general.url.io http://demo.openshift.io/work-item/list/detail/100 number:300 golang book and           number:900 \t \n unwanted"
	op, _ := parseSearchString(context.Background(), input)
	expectedSearchRes := searchKeyword{
		number: []string{"300:*A", "900:*A"},
		words:  []string{"general.url.io:*", "(100:*A | demo.openshift.io/work-item/list/detail/100:*)", "golang:*", "book:*", "and:*", "unwanted:*"},
	}
	assert.True(t, assert.ObjectsAreEqualValues(expectedSearchRes, op))
}

func TestRegisterAsKnownURL(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	// build 2 fake urls and cross check against RegisterAsKnownURL
	urlRegex := `(?P<domain>google.me.io)(?P<path>/everything/)(?P<param>.*)`
	routeName := "custom-test-route"
	RegisterAsKnownURL(routeName, urlRegex)
	compiledRegex := regexp.MustCompile(urlRegex)
	groupNames := compiledRegex.SubexpNames()
	var expected = make(map[string]KnownURL)
	expected[routeName] = KnownURL{
		URLRegex:          urlRegex,
		compiledRegex:     regexp.MustCompile(urlRegex),
		groupNamesInRegex: groupNames,
	}
	assert.True(t, assert.ObjectsAreEqualValues(expected[routeName], knownURLs[routeName]))
	//cleanup
	delete(knownURLs, routeName)
}

func TestIsKnownURL(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	// register few URLs and cross check is knwon or not one by one
	urlRegex := `(?P<domain>google.me.io)(?P<path>/everything/)(?P<param>.*)`
	routeName := "custom-test-route"
	RegisterAsKnownURL(routeName, urlRegex)
	known, patternName := isKnownURL("google.me.io/everything/v1/v2/q=1")
	assert.True(t, known)
	assert.Equal(t, routeName, patternName)

	known, patternName = isKnownURL("google.different.io/everything/v1/v2/q=1")
	assert.False(t, known)
	assert.Equal(t, "", patternName)

	// cleanup
	delete(knownURLs, routeName)
}

func TestGetSearchQueryFromURLPattern(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	// getSearchQueryFromURLPattern
	// register urls
	// select pattern and pass search string
	// validate output with different scenarios like ID present not present
	urlRegex := `(?P<domain>google.me.io)(?P<path>/everything/)(?P<id>\d*)`
	routeName := "custom-test-route"
	RegisterAsKnownURL(routeName, urlRegex)

	searchQuery := getSearchQueryFromURLPattern(routeName, "google.me.io/everything/100")
	assert.Equal(t, "(100:*A | google.me.io/everything/100:*)", searchQuery)

	searchQuery = getSearchQueryFromURLPattern(routeName, "google.me.io/everything/")
	assert.Equal(t, "google.me.io/everything/:*", searchQuery)

	// cleanup
	delete(knownURLs, routeName)
}

func TestGetSearchQueryFromURLString(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	// register few urls
	// call getSearchQueryFromURLString with different urls - both registered and non-registered
	searchQuery := getSearchQueryFromURLString("abcd.something.com")
	assert.Equal(t, "abcd.something.com:*", searchQuery)

	urlRegex := `(?P<domain>google.me.io)(?P<path>/everything/)(?P<id>\d*)`
	routeName := "custom-test-route"
	RegisterAsKnownURL(routeName, urlRegex)

	searchQuery = getSearchQueryFromURLString("google.me.io/everything/")
	assert.Equal(t, "google.me.io/everything/:*", searchQuery)

	searchQuery = getSearchQueryFromURLString("google.me.io/everything/100")
	assert.Equal(t, "(100:*A | google.me.io/everything/100:*)", searchQuery)
}

func TestIsOperator(t *testing.T) {
	testData := map[string]bool{
		AND:   true,
		OR:    true,
		OPTS:  false,
		"":    false,
		"   ": false,
		"foo": false,
		uuid.NewV4().String(): false,
		EQ:       false,
		NE:       false,
		NOT:      false,
		IN:       false,
		SUBSTR:   false,
		WITGROUP: false,
	}
	for k, v := range testData {
		t.Run(k, func(t *testing.T) {
			if v {
				require.True(t, isOperator(k), "isOperator(%s) should be true", k)
			} else {
				require.False(t, isOperator(k), "isOperator(%s) should be false", k)
			}
		})
	}
}

func TestHandleWitGroup(t *testing.T) {
	type testData struct {
		Name                string
		Value               string
		Negate              bool
		ExpectError         bool
		ExpectedExrpessions []criteria.Expression
	}
	td := []testData{
		{"foo", "bar", false, false, []criteria.Expression{}},
		{WITGROUP, "unknown", false, true, []criteria.Expression{}},
		{WITGROUP, "Scenarios", false, false, []criteria.Expression{
			criteria.Or(
				criteria.Or(
					criteria.Equals(
						criteria.Field("Type"),
						criteria.Literal(workitem.SystemScenario.String()),
					),
					criteria.Equals(
						criteria.Field("Type"),
						criteria.Literal(workitem.SystemFundamental.String()),
					),
				),
				criteria.Equals(
					criteria.Field("Type"),
					criteria.Literal(workitem.SystemPapercuts.String()),
				),
			)},
		},
		{WITGROUP, "Experiences", false, false, []criteria.Expression{
			criteria.Or(
				criteria.Equals(
					criteria.Field("Type"),
					criteria.Literal(workitem.SystemExperience.String()),
				),
				criteria.Equals(
					criteria.Field("Type"),
					criteria.Literal(workitem.SystemValueProposition.String()),
				),
			)},
		},
		{WITGROUP, "Requirements", false, false, []criteria.Expression{
			criteria.Or(
				criteria.Equals(
					criteria.Field("Type"),
					criteria.Literal(workitem.SystemFeature.String()),
				),
				criteria.Equals(
					criteria.Field("Type"),
					criteria.Literal(workitem.SystemBug.String()),
				),
			)},
		},
		{WITGROUP, "Execution", false, false, []criteria.Expression{
			criteria.Or(
				criteria.Or(
					criteria.Equals(
						criteria.Field("Type"),
						criteria.Literal(workitem.SystemTask.String()),
					),
					criteria.Equals(
						criteria.Field("Type"),
						criteria.Literal(workitem.SystemBug.String()),
					),
				),
				criteria.Equals(
					criteria.Field("Type"),
					criteria.Literal(workitem.SystemFeature.String()),
				),
			),
		}},
		// // same with negation
		{"foo", "bar", true, false, []criteria.Expression{}},
		{WITGROUP, "unknown", true, true, []criteria.Expression{}},
		{WITGROUP, "Scenarios", true, false, []criteria.Expression{
			criteria.And(
				criteria.And(
					criteria.Not(
						criteria.Field("Type"),
						criteria.Literal(workitem.SystemScenario.String()),
					),
					criteria.Not(
						criteria.Field("Type"),
						criteria.Literal(workitem.SystemFundamental.String()),
					),
				),
				criteria.Not(
					criteria.Field("Type"),
					criteria.Literal(workitem.SystemPapercuts.String()),
				),
			)},
		},
		{WITGROUP, "Experiences", true, false, []criteria.Expression{
			criteria.And(
				criteria.Not(
					criteria.Field("Type"),
					criteria.Literal(workitem.SystemExperience.String()),
				),
				criteria.Not(
					criteria.Field("Type"),
					criteria.Literal(workitem.SystemValueProposition.String()),
				),
			)},
		},
		{WITGROUP, "Requirements", true, false, []criteria.Expression{
			criteria.And(
				criteria.Not(
					criteria.Field("Type"),
					criteria.Literal(workitem.SystemFeature.String()),
				),
				criteria.Not(
					criteria.Field("Type"),
					criteria.Literal(workitem.SystemBug.String()),
				),
			)},
		},
		{WITGROUP, "Execution", true, false, []criteria.Expression{
			criteria.And(
				criteria.And(
					criteria.Not(
						criteria.Field("Type"),
						criteria.Literal(workitem.SystemTask.String()),
					),
					criteria.Not(
						criteria.Field("Type"),
						criteria.Literal(workitem.SystemBug.String()),
					),
				),
				criteria.Not(
					criteria.Field("Type"),
					criteria.Literal(workitem.SystemFeature.String()),
				),
			),
		}},
	}
	for _, d := range td {
		format := "%s = %s"
		if d.Negate {
			format = "%s != %s"
		}
		t.Run(fmt.Sprintf(format, d.Name, d.Value), func(t *testing.T) {
			exp := []criteria.Expression{}
			err := handleWitGroup(Query{Name: d.Name, Value: &d.Value, Negate: d.Negate}, &exp)
			if d.ExpectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, exp, d.ExpectedExrpessions)
		})
	}

	t.Run("value is nil", func(t *testing.T) {
		// given
		var v *string
		exp := []criteria.Expression{}
		// when
		err := handleWitGroup(Query{Name: WITGROUP, Value: v}, &exp)
		// then
		require.Error(t, err)
	})
	t.Run("expression array is nil", func(t *testing.T) {
		// given
		v := "Scenarios"
		var exp *[]criteria.Expression
		// when
		err := handleWitGroup(Query{Name: WITGROUP, Value: &v}, exp)
		// then
		require.Error(t, err)
	})
}
