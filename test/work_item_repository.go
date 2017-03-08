// This file was generated by counterfeiter
package test

import (
	"sync"

	"github.com/almighty/almighty-core/app"
	"github.com/almighty/almighty-core/criteria"
	"github.com/almighty/almighty-core/workitem"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

type WorkItemRepository struct {
	LoadStub        func(ctx context.Context, ID string) (*app.WorkItem, error)
	loadMutex       sync.RWMutex
	loadArgsForCall []struct {
		ctx context.Context
		ID  string
	}
	loadReturns struct {
		result1 *app.WorkItem
		result2 error
	}
	SaveStub        func(ctx context.Context, wi app.WorkItem, modifierID uuid.UUID) (*app.WorkItem, error)
	saveMutex       sync.RWMutex
	saveArgsForCall []struct {
		ctx        context.Context
		wi         app.WorkItem
		modifierID uuid.UUID
	}
	saveReturns struct {
		result1 *app.WorkItem
		result2 error
	}
	ReorderStub        func(ctx context.Context, direction workitem.DirectionType, targetID *string, wi app.WorkItem, modifierID uuid.UUID) (*app.WorkItem, error)
	reorderMutex       sync.RWMutex
	reorderArgsForCall []struct {
		ctx        context.Context
		direction  workitem.DirectionType
		targetID   *string
		wi         app.WorkItem
		modifierID uuid.UUID
	}
	reorderReturns struct {
		wi  *app.WorkItem
		err error
	}
	DeleteStub        func(ctx context.Context, ID string, suppressorID uuid.UUID) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		ctx          context.Context
		ID           string
		suppressorID uuid.UUID
	}
	deleteReturns struct {
		result1 error
	}
	CreateStub        func(ctx context.Context, typeID uuid.UUID, fields map[string]interface{}, creatorID uuid.UUID) (*app.WorkItem, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		ctx       context.Context
		typeID    uuid.UUID
		fields    map[string]interface{}
		creatorID uuid.UUID
	}
	createReturns struct {
		result1 *app.WorkItem
		result2 error
	}
	ListStub        func(ctx context.Context, criteria criteria.Expression, start *int, length *int) ([]*app.WorkItem, uint64, error)
	listMutex       sync.RWMutex
	listArgsForCall []struct {
		ctx      context.Context
		criteria criteria.Expression
		start    *int
		length   *int
	}
	listReturns struct {
		result1 []*app.WorkItem
		result2 uint64
		result3 error
	}
	FetchStub        func(ctx context.Context, criteria criteria.Expression) (*app.WorkItem, error)
	fetchMutex       sync.RWMutex
	fetchArgsForCall []struct {
		ctx      context.Context
		criteria criteria.Expression
	}
	fetchReturns struct {
		result1 *app.WorkItem
		result2 error
	}
	GetCountsPerIterationStub        func(ctx context.Context, spaceID uuid.UUID) (map[string]workitem.WICountsPerIteration, error)
	getCountsPerIterationMutex       sync.RWMutex
	getCountsPerIterationArgsForCall []struct {
		ctx     context.Context
		spaceID uuid.UUID
	}
	getCountsPerIterationReturns struct {
		result1 map[string]workitem.WICountsPerIteration
		result2 error
	}
	GetCountsForIterationStub        func(ctx context.Context, iterationID uuid.UUID) (map[string]workitem.WICountsPerIteration, error)
	getCountsForIterationMutex       sync.RWMutex
	getCountsForIterationArgsForCall []struct {
		ctx         context.Context
		iterationID uuid.UUID
	}
	getCountsForIterationReturns struct {
		result1 map[string]workitem.WICountsPerIteration
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *WorkItemRepository) Load(ctx context.Context, ID string) (*app.WorkItem, error) {
	fake.loadMutex.Lock()
	fake.loadArgsForCall = append(fake.loadArgsForCall, struct {
		ctx context.Context
		ID  string
	}{ctx, ID})
	fake.recordInvocation("Load", []interface{}{ctx, ID})
	fake.loadMutex.Unlock()
	if fake.LoadStub != nil {
		return fake.LoadStub(ctx, ID)
	}
	return fake.loadReturns.result1, fake.loadReturns.result2
}

func (fake *WorkItemRepository) LoadCallCount() int {
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	return len(fake.loadArgsForCall)
}

func (fake *WorkItemRepository) LoadArgsForCall(i int) (context.Context, string) {
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	return fake.loadArgsForCall[i].ctx, fake.loadArgsForCall[i].ID
}

func (fake *WorkItemRepository) LoadReturns(result1 *app.WorkItem, result2 error) {
	fake.LoadStub = nil
	fake.loadReturns = struct {
		result1 *app.WorkItem
		result2 error
	}{result1, result2}
}

func (fake *WorkItemRepository) Save(ctx context.Context, wi app.WorkItem, modifierID uuid.UUID) (*app.WorkItem, error) {
	fake.saveMutex.Lock()
	fake.saveArgsForCall = append(fake.saveArgsForCall, struct {
		ctx        context.Context
		wi         app.WorkItem
		modifierID uuid.UUID
	}{ctx, wi, modifierID})
	fake.recordInvocation("Save", []interface{}{ctx, wi, modifierID})
	fake.saveMutex.Unlock()
	if fake.SaveStub != nil {
		return fake.SaveStub(ctx, wi, modifierID)
	}
	return fake.saveReturns.result1, fake.saveReturns.result2
}

func (fake *WorkItemRepository) SaveCallCount() int {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	return len(fake.saveArgsForCall)
}

func (fake *WorkItemRepository) SaveArgsForCall(i int) (context.Context, app.WorkItem, uuid.UUID) {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	return fake.saveArgsForCall[i].ctx, fake.saveArgsForCall[i].wi, fake.saveArgsForCall[i].modifierID
}

func (fake *WorkItemRepository) SaveReturns(result1 *app.WorkItem, result2 error) {
	fake.SaveStub = nil
	fake.saveReturns = struct {
		result1 *app.WorkItem
		result2 error
	}{result1, result2}
}

// Reorder is a fake function for reordering of workitems
// Used for testing purpose
func (fake *WorkItemRepository) Reorder(ctx context.Context, direction workitem.DirectionType, targetID *string, wi app.WorkItem, modifierID uuid.UUID) (*app.WorkItem, error) {
	fake.reorderMutex.Lock()
	fake.reorderArgsForCall = append(fake.reorderArgsForCall, struct {
		ctx        context.Context
		direction  workitem.DirectionType
		targetID   *string
		wi         app.WorkItem
		modifierID uuid.UUID
	}{ctx, direction, targetID, wi, modifierID})
	fake.recordInvocation("Reorder", []interface{}{ctx, direction, *targetID, wi, modifierID})
	fake.reorderMutex.Unlock()
	if fake.ReorderStub != nil {
		return fake.ReorderStub(ctx, direction, targetID, wi, modifierID)
	} else {
		return fake.reorderReturns.wi, fake.reorderReturns.err
	}
}

// ReorderCallCount returns the length of fake arguments
func (fake *WorkItemRepository) ReorderCallCount() int {
	fake.reorderMutex.RLock()
	defer fake.reorderMutex.RUnlock()
	return len(fake.reorderArgsForCall)
}

// ReorderArgsForCall returns fake arguments for Reorder function
func (fake *WorkItemRepository) ReorderArgsForCall(i int) (context.Context, workitem.DirectionType, string, app.WorkItem, uuid.UUID) {
	fake.reorderMutex.RLock()
	defer fake.reorderMutex.RUnlock()
	return fake.reorderArgsForCall[i].ctx, fake.reorderArgsForCall[i].direction, *fake.reorderArgsForCall[i].targetID, fake.reorderArgsForCall[i].wi, fake.reorderArgsForCall[i].modifierID
}

// ReorderReturns returns fake values for Reorder function
func (fake *WorkItemRepository) ReorderReturns(workItem *app.WorkItem, errr error) {
	fake.ReorderStub = nil
	type reorder struct {
		wi  *app.WorkItem
		err error
	}
	v := reorder{workItem, errr}
	fake.reorderReturns = v
}

func (fake *WorkItemRepository) Delete(ctx context.Context, ID string, suppressorID uuid.UUID) error {
	fake.deleteMutex.Lock()
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		ctx          context.Context
		ID           string
		suppressorID uuid.UUID
	}{ctx, ID, suppressorID})
	fake.recordInvocation("Delete", []interface{}{ctx, ID, suppressorID})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(ctx, ID, suppressorID)
	}
	return fake.deleteReturns.result1
}

func (fake *WorkItemRepository) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *WorkItemRepository) DeleteArgsForCall(i int) (context.Context, string, uuid.UUID) {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return fake.deleteArgsForCall[i].ctx, fake.deleteArgsForCall[i].ID, fake.deleteArgsForCall[i].suppressorID
}

func (fake *WorkItemRepository) DeleteReturns(result1 error) {
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *WorkItemRepository) Create(ctx context.Context, typeID uuid.UUID, fields map[string]interface{}, creatorID uuid.UUID) (*app.WorkItem, error) {
	fake.createMutex.Lock()
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		ctx       context.Context
		typeID    uuid.UUID
		fields    map[string]interface{}
		creatorID uuid.UUID
	}{ctx, typeID, fields, creatorID})
	fake.recordInvocation("Create", []interface{}{ctx, typeID, fields, creatorID})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(ctx, typeID, fields, creatorID)
	}
	return fake.createReturns.result1, fake.createReturns.result2
}

func (fake *WorkItemRepository) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *WorkItemRepository) CreateArgsForCall(i int) (context.Context, uuid.UUID, map[string]interface{}, uuid.UUID) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].ctx, fake.createArgsForCall[i].typeID, fake.createArgsForCall[i].fields, fake.createArgsForCall[i].creatorID
}

func (fake *WorkItemRepository) CreateReturns(result1 *app.WorkItem, result2 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 *app.WorkItem
		result2 error
	}{result1, result2}
}

func (fake *WorkItemRepository) List(ctx context.Context, c criteria.Expression, start *int, length *int) ([]*app.WorkItem, uint64, error) {
	fake.listMutex.Lock()
	fake.listArgsForCall = append(fake.listArgsForCall, struct {
		ctx      context.Context
		criteria criteria.Expression
		start    *int
		length   *int
	}{ctx, c, start, length})
	fake.recordInvocation("List", []interface{}{ctx, c, start, length})
	fake.listMutex.Unlock()
	if fake.ListStub != nil {
		return fake.ListStub(ctx, c, start, length)
	}
	return fake.listReturns.result1, fake.listReturns.result2, fake.listReturns.result3
}

func (fake *WorkItemRepository) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *WorkItemRepository) ListArgsForCall(i int) (context.Context, criteria.Expression, *int, *int) {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return fake.listArgsForCall[i].ctx, fake.listArgsForCall[i].criteria, fake.listArgsForCall[i].start, fake.listArgsForCall[i].length
}

func (fake *WorkItemRepository) ListReturns(result1 []*app.WorkItem, result2 uint64, result3 error) {
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 []*app.WorkItem
		result2 uint64
		result3 error
	}{result1, result2, result3}
}

func (fake *WorkItemRepository) Fetch(ctx context.Context, c criteria.Expression) (*app.WorkItem, error) {
	fake.fetchMutex.Lock()
	fake.fetchArgsForCall = append(fake.fetchArgsForCall, struct {
		ctx      context.Context
		criteria criteria.Expression
	}{ctx, c})
	fake.recordInvocation("Fetch", []interface{}{ctx, c})
	fake.fetchMutex.Unlock()
	if fake.FetchStub != nil {
		return fake.FetchStub(ctx, c)
	}
	return fake.fetchReturns.result1, fake.fetchReturns.result2
}

func (fake *WorkItemRepository) FetchCallCount() int {
	fake.fetchMutex.RLock()
	defer fake.fetchMutex.RUnlock()
	return len(fake.fetchArgsForCall)
}

func (fake *WorkItemRepository) FetchArgsForCall(i int) (context.Context, criteria.Expression) {
	fake.fetchMutex.RLock()
	defer fake.fetchMutex.RUnlock()
	return fake.fetchArgsForCall[i].ctx, fake.fetchArgsForCall[i].criteria
}

func (fake *WorkItemRepository) FetchReturns(result1 *app.WorkItem, result2 error) {
	fake.FetchStub = nil
	fake.fetchReturns = struct {
		result1 *app.WorkItem
		result2 error
	}{result1, result2}
}

func (fake *WorkItemRepository) GetCountsPerIteration(ctx context.Context, spaceID uuid.UUID) (map[string]workitem.WICountsPerIteration, error) {
	fake.getCountsPerIterationMutex.Lock()
	fake.getCountsPerIterationArgsForCall = append(fake.getCountsPerIterationArgsForCall, struct {
		ctx     context.Context
		spaceID uuid.UUID
	}{ctx, spaceID})
	fake.recordInvocation("GetCountsPerIteration", []interface{}{ctx, spaceID})
	fake.getCountsPerIterationMutex.Unlock()
	if fake.GetCountsPerIterationStub != nil {
		return fake.GetCountsPerIterationStub(ctx, spaceID)
	}
	return fake.getCountsPerIterationReturns.result1, fake.getCountsPerIterationReturns.result2
}

func (fake *WorkItemRepository) GetCountsPerIterationCallCount() int {
	fake.getCountsPerIterationMutex.RLock()
	defer fake.getCountsPerIterationMutex.RUnlock()
	return len(fake.getCountsPerIterationArgsForCall)
}

func (fake *WorkItemRepository) GetCountsPerIterationArgsForCall(i int) (context.Context, uuid.UUID) {
	fake.getCountsPerIterationMutex.RLock()
	defer fake.getCountsPerIterationMutex.RUnlock()
	return fake.getCountsPerIterationArgsForCall[i].ctx, fake.getCountsPerIterationArgsForCall[i].spaceID
}

func (fake *WorkItemRepository) GetCountsPerIterationReturns(result1 map[string]workitem.WICountsPerIteration, result2 error) {
	fake.GetCountsPerIterationStub = nil
	fake.getCountsPerIterationReturns = struct {
		result1 map[string]workitem.WICountsPerIteration
		result2 error
	}{result1, result2}
}

func (fake *WorkItemRepository) GetCountsForIteration(ctx context.Context, iterationID uuid.UUID) (map[string]workitem.WICountsPerIteration, error) {
	fake.getCountsForIterationMutex.Lock()
	fake.getCountsForIterationArgsForCall = append(fake.getCountsForIterationArgsForCall, struct {
		ctx         context.Context
		iterationID uuid.UUID
	}{ctx, iterationID})
	fake.recordInvocation("GetCountsForIteration", []interface{}{ctx, iterationID})
	fake.getCountsForIterationMutex.Unlock()
	if fake.GetCountsForIterationStub != nil {
		return fake.GetCountsForIterationStub(ctx, iterationID)
	}
	return fake.getCountsForIterationReturns.result1, fake.getCountsForIterationReturns.result2
}

func (fake *WorkItemRepository) GetCountsForIterationCallCount() int {
	fake.getCountsForIterationMutex.RLock()
	defer fake.getCountsForIterationMutex.RUnlock()
	return len(fake.getCountsForIterationArgsForCall)
}

func (fake *WorkItemRepository) GetCountsForIterationArgsForCall(i int) (context.Context, uuid.UUID) {
	fake.getCountsForIterationMutex.RLock()
	defer fake.getCountsForIterationMutex.RUnlock()
	return fake.getCountsForIterationArgsForCall[i].ctx, fake.getCountsForIterationArgsForCall[i].iterationID
}

func (fake *WorkItemRepository) GetCountsForIterationReturns(result1 map[string]workitem.WICountsPerIteration, result2 error) {
	fake.GetCountsForIterationStub = nil
	fake.getCountsForIterationReturns = struct {
		result1 map[string]workitem.WICountsPerIteration
		result2 error
	}{result1, result2}
}

func (fake *WorkItemRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	fake.fetchMutex.RLock()
	defer fake.fetchMutex.RUnlock()
	fake.getCountsPerIterationMutex.RLock()
	defer fake.getCountsPerIterationMutex.RUnlock()
	fake.getCountsForIterationMutex.RLock()
	defer fake.getCountsForIterationMutex.RUnlock()
	return fake.invocations
}

func (fake *WorkItemRepository) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ workitem.WorkItemRepository = new(WorkItemRepository)
