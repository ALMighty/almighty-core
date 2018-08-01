package kubernetes

import (
	"fmt"
	"net/http"

	"github.com/fabric8-services/fabric8-wit/log"
	errs "github.com/pkg/errors"
)

const (
	verbCreate           = "create"
	verbDelete           = "delete"
	verbDeleteCollection = "deletecollection"
	verbGet              = "get"
	verbList             = "list"
	verbPatch            = "patch"
	verbUpdate           = "update"
	verbWatch            = "watch"
)

// KubeAccessControl contains methods that answer whether the current user
// has sufficient authorization to call various methods of KubeClientInterface
type KubeAccessControl interface {
	CanGetSpace() (bool, error)
	CanGetApplication() (bool, error)
	CanGetDeployment(envName string) (bool, error)
	CanScaleDeployment(envName string) (bool, error)
	CanGetDeploymentStats(envName string) (bool, error)
	CanGetDeploymentStatSeries(envName string) (bool, error)
	CanDeleteDeployment(envName string) (bool, error)
	CanGetEnvironments() (bool, error)
	CanGetEnvironment(envName string) (bool, error)
}

// Actions on a resource type that are required by one of our API methods
type requestedAccess struct {
	resource qualifiedResource
	verbs    []string
}

// Maps resource types to authorized actions that may be performed by the user
type accessRules map[qualifiedResource]simpleAccessRule

// Names a resource type by group name and resource type
type qualifiedResource struct {
	apiGroup     string
	resourceType string
}

// Only handle rules that aren't qualified by resource name or URL
type simpleAccessRule map[string]struct{}

// Checks the subject rules review for the desired actions on resources
func (rulesMap accessRules) isAuthorized(reqs []*requestedAccess) bool {
	for _, req := range reqs {
		// Look up rules for resource type
		rules, pres := rulesMap[req.resource]
		if !pres {
			return false
		}
		// Check if all requested actions are permitted
		for _, verb := range req.verbs {
			_, pres := rules[verb]
			if !pres {
				return false
			}
		}
	}
	return true
}

func (kc *kubeClient) CanGetSpace() (bool, error) {
	// Also need access to build configs and builds in user namespace
	ok, err := kc.checkAuthorizedInEnv(getBuildConfigsAndBuildsRules, environmentTypeUser)
	if err != nil {
		return false, err
	} else if !ok {
		return false, nil
	}

	for envName := range kc.envMap {
		if kc.CanDeploy(envName) {
			ok, err := kc.checkAuthorizedInEnv(getDeploymentRules, envName)
			if err != nil {
				return false, err
			} else if !ok {
				return false, nil
			}
		}
	}
	return true, nil
}

func (kc *kubeClient) CanGetApplication() (bool, error) {
	// Also need access to builds in user namespace
	ok, err := kc.checkAuthorizedInEnv(getBuildsRules, environmentTypeUser)
	if err != nil {
		return false, err
	} else if !ok {
		return false, nil
	}

	for envName := range kc.envMap {
		if kc.CanDeploy(envName) {
			ok, err := kc.checkAuthorizedInEnv(getDeploymentRules, envName)
			if err != nil {
				return false, err
			} else if !ok {
				return false, nil
			}
		}
	}
	return true, nil
}

var getDeploymentRules = []*requestedAccess{
	{qualifiedResource{"", "deploymentconfigs"}, []string{verbGet}},
	{qualifiedResource{"", "replicationcontrollers"}, []string{verbList}},
	{qualifiedResource{"", "pods"}, []string{verbList}},
	{qualifiedResource{"", "services"}, []string{verbList}},
	{qualifiedResource{"", "routes"}, []string{verbList}},
}

func (kc *kubeClient) CanGetDeployment(envName string) (bool, error) {
	return kc.checkAuthorizedWithBuilds(envName, getDeploymentRules)
}

var scaleDeploymentRules = []*requestedAccess{
	{qualifiedResource{"", "deploymentconfigs"}, []string{verbGet}},
	{qualifiedResource{"", "deploymentconfigs/scale"}, []string{verbGet}},
	{qualifiedResource{"", "deploymentconfigs/scale"}, []string{verbUpdate}},
}

func (kc *kubeClient) CanScaleDeployment(envName string) (bool, error) {
	return kc.checkAuthorizedWithBuilds(envName, scaleDeploymentRules)
}

var deleteDeploymentRules = []*requestedAccess{
	{qualifiedResource{"", "services"}, []string{verbList, verbDelete}},
	{qualifiedResource{"", "routes"}, []string{verbList, verbDelete}},
	{qualifiedResource{"", "deploymentconfigs"}, []string{verbGet, verbDelete}},
}

// TODO Do we want resource name-specific checks. Can delete DC named X?
func (kc *kubeClient) CanDeleteDeployment(envName string) (bool, error) {
	return kc.checkAuthorizedWithBuilds(envName, deleteDeploymentRules)
}

var getDeploymentStatsRules = []*requestedAccess{
	{qualifiedResource{"", "deploymentconfigs"}, []string{verbGet}},
	{qualifiedResource{"", "replicationcontrollers"}, []string{verbList}},
	{qualifiedResource{"", "pods"}, []string{verbList}},
}

func (kc *kubeClient) CanGetDeploymentStats(envName string) (bool, error) {
	return kc.checkAuthorizedWithBuilds(envName, getDeploymentStatsRules)
}

func (kc *kubeClient) CanGetDeploymentStatSeries(envName string) (bool, error) {
	return kc.checkAuthorizedWithBuilds(envName, getDeploymentStatsRules)
}

func (kc *kubeClient) checkAuthorizedWithBuilds(envName string, reqs []*requestedAccess) (bool, error) {
	// Also need access to builds in user namespace
	ok, err := kc.checkAuthorizedInEnv(getBuildsRules, environmentTypeUser)
	if err != nil {
		return false, err
	} else if !ok {
		return false, nil
	}

	return kc.checkAuthorizedInEnv(reqs, envName)
}

const environmentTypeUser = "user"

var getBuildConfigsAndBuildsRules = []*requestedAccess{
	{qualifiedResource{"", "buildconfigs"}, []string{verbList}},
	{qualifiedResource{"", "builds"}, []string{verbList}},
}

var getBuildsRules = []*requestedAccess{
	{qualifiedResource{"", "builds"}, []string{verbList}},
}

func (kc *kubeClient) checkAuthorizedInEnv(reqs []*requestedAccess, envName string) (bool, error) {
	rules, err := kc.getRulesForEnvironment(envName)
	if err != nil {
		return false, err
	}

	return rules.isAuthorized(reqs), nil
}

var getEnvironmentRules = []*requestedAccess{
	{qualifiedResource{"", "resourcequotas"}, []string{verbList}},
}

func (kc *kubeClient) CanGetEnvironments() (bool, error) {
	for envName := range kc.envMap {
		if kc.CanDeploy(envName) {
			ok, err := kc.CanGetEnvironment(envName)
			if err != nil {
				return false, err
			} else if !ok {
				return false, nil
			}
		}
	}
	return true, nil
}

func (kc *kubeClient) CanGetEnvironment(envName string) (bool, error) {
	return kc.checkAuthorizedInEnv(getEnvironmentRules, envName)
}

// Gets the authorization rules for the current user in a given environment
func (kc *kubeClient) getRulesForEnvironment(envName string) (*accessRules, error) {
	// Check if we have a cached copy
	rules, pres := kc.rulesMap[envName]
	if pres {
		return rules, nil
	}

	// Lookup authorization rules for this environment
	envNS, err := kc.getEnvironmentNamespace(envName, true)
	if err != nil {
		return nil, err
	}
	rules, err = kc.lookupAllRules(envNS)
	if err != nil {
		return nil, err
	}

	// Cache rules, so subsequent calls by this kubeClient don't
	// trigger lookup over network
	kc.rulesMap[envName] = rules
	return rules, nil
}

func (kc *kubeClient) lookupAllRules(namespace string) (*accessRules, error) {
	review := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "SelfSubjectRulesReview",
	}
	reviewResult, err := kc.CreateSelfSubjectRulesReview(namespace, review)
	if err != nil {
		return nil, err
	}

	// TODO Parse using info from https://github.com/openshift/api/blob/master/authorization/v1/types.go
	status, ok := reviewResult["status"].(map[string]interface{})
	if !ok {
		log.Error(nil, map[string]interface{}{
			"err":       err,
			"namespace": namespace,
			"response":  reviewResult,
		}, "status missing from SelfSubjectRulesReview")
		return nil, errs.Errorf("status missing from SelfSubjectRulesReview returned from %s", namespace)
	}
	rules, ok := status["rules"].([]interface{})
	if !ok {
		log.Error(nil, map[string]interface{}{
			"err":       err,
			"namespace": namespace,
			"response":  reviewResult,
		}, "rules missing from SelfSubjectRulesReview")
		return nil, errs.Errorf("rules missing from SelfSubjectRulesReview returned from %s", namespace)
	}

	result := make(accessRules)
	for _, rawRule := range rules {
		rule, ok := rawRule.(map[string]interface{})
		if !ok {
			log.Error(nil, map[string]interface{}{
				"err":       err,
				"namespace": namespace,
				"response":  reviewResult,
			}, "rules missing from SelfSubjectRulesReview")
			return nil, errs.Errorf("rule returned from %s is not a JSON object", namespace)
		}

		processRule(result, rule)
	}
	return &result, nil
}

func processRule(rules accessRules, rule map[string]interface{}) {
	// For now, only consider rules that don't specify particular resource names or URLs
	resourceNames, ok := rule["resourceNames"].([]interface{})
	if ok && len(resourceNames) > 0 {
		return
	}
	nonResourceURLs, ok := rule["nonResourceURLs"].([]interface{})
	if ok && len(nonResourceURLs) > 0 {
		return
	}

	verbs := getStringSetFromJSON(rule, "verbs")
	groups := getStringSliceFromJSON(rule, "apiGroups")
	resources := getStringSliceFromJSON(rule, "resources")

	// Add verbs for each group/resource in rule
	for _, resource := range resources {
		// If no groups are specified, the rule is for the default k8s/OpenShift API group
		if len(groups) == 0 {
			key := qualifiedResource{"", resource}
			rules[key] = verbs
		} else {
			for _, group := range groups {
				key := qualifiedResource{group, resource}
				rules[key] = verbs
			}
		}
	}
}

func getStringSliceFromJSON(jsonObj map[string]interface{}, name string) []string {
	var items []string
	jsonArray, ok := jsonObj[name].([]interface{})
	if ok {
		items = make([]string, 0, len(jsonArray))
		for _, jsonItem := range jsonArray {
			item, ok := jsonItem.(string)
			if !ok {
				log.Error(nil, map[string]interface{}{
					"item":        jsonItem,
					"json_object": jsonObj,
				}, "item in %s array is not a string", name)
			}
			items = append(items, item)
		}
	}
	return items
}

func getStringSetFromJSON(jsonObj map[string]interface{}, name string) map[string]struct{} {
	var items map[string]struct{}
	jsonArray, ok := jsonObj[name].([]interface{})
	if ok {
		items = make(map[string]struct{}, len(jsonArray))
		for _, jsonItem := range jsonArray {
			item, ok := jsonItem.(string)
			if !ok {
				log.Error(nil, map[string]interface{}{
					"item":        jsonItem,
					"json_object": jsonObj,
				}, "item in %s array is not a string", name)
			}
			items[item] = struct{}{}
		}
	}
	return items
}

func (oc *openShiftAPIClient) CreateSelfSubjectRulesReview(namespace string,
	review map[string]interface{}) (map[string]interface{}, error) {
	reviewPath := fmt.Sprintf("/oapi/v1/namespaces/%s/selfsubjectrulesreviews", namespace)
	return oc.sendResource(reviewPath, http.MethodPost, review)
}