// This file was generated by counterfeiter
package test

import (
	"sync"

	"github.com/almighty/almighty-core/app"
	"github.com/almighty/almighty-core/criteria"
	"github.com/almighty/almighty-core/workitem"
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
	SaveStub        func(ctx context.Context, wi app.WorkItem) (*app.WorkItem, error)
	saveMutex       sync.RWMutex
	saveArgsForCall []struct {
		ctx context.Context
		wi  app.WorkItem
	}
	saveReturns struct {
		result1 *app.WorkItem
		result2 error
	}
	ReorderStub        func(ctx context.Context, before string, wi app.WorkItem) (*app.WorkItem, error)
	reorderMutex       sync.RWMutex
	reorderArgsForCall []struct {
		ctx    context.Context
		before string
		wi     app.WorkItem
	}
	reorderReturns struct {
		result1 *app.WorkItem
		result2 error
	}
	DeleteStub        func(ctx context.Context, ID string) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		ctx context.Context
		ID  string
	}
	deleteReturns struct {
		result1 error
	}
	CreateStub        func(ctx context.Context, typeID string, fields map[string]interface{}) (*app.WorkItem, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		ctx    context.Context
		typeID string
		fields map[string]interface{}
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
	} else {
		return fake.loadReturns.result1, fake.loadReturns.result2
	}
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

func (fake *WorkItemRepository) Save(ctx context.Context, wi app.WorkItem) (*app.WorkItem, error) {
	fake.saveMutex.Lock()
	fake.saveArgsForCall = append(fake.saveArgsForCall, struct {
		ctx context.Context
		wi  app.WorkItem
	}{ctx, wi})
	fake.recordInvocation("Save", []interface{}{ctx, wi})
	fake.saveMutex.Unlock()
	if fake.SaveStub != nil {
		return fake.SaveStub(ctx, wi)
	} else {
		return fake.saveReturns.result1, fake.saveReturns.result2
	}
}

func (fake *WorkItemRepository) SaveCallCount() int {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	return len(fake.saveArgsForCall)
}

func (fake *WorkItemRepository) SaveArgsForCall(i int) (context.Context, app.WorkItem) {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	return fake.saveArgsForCall[i].ctx, fake.saveArgsForCall[i].wi
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
func (fake *WorkItemRepository) Reorder(ctx context.Context, before string, wi app.WorkItem) (*app.WorkItem, error) {
	fake.reorderMutex.Lock()
	fake.reorderArgsForCall = append(fake.reorderArgsForCall, struct {
		ctx    context.Context
		before string
		wi     app.WorkItem
	}{ctx, before, wi})
	fake.recordInvocation("Reorder", []interface{}{ctx, before, wi})
	fake.reorderMutex.Unlock()
	if fake.ReorderStub != nil {
		return fake.ReorderStub(ctx, before, wi)
	} else {
		return fake.reorderReturns.result1, fake.reorderReturns.result2
	}
}

// ReorderCallCount returns the length of fake arguments
func (fake *WorkItemRepository) ReorderCallCount() int {
	fake.reorderMutex.RLock()
	defer fake.reorderMutex.RUnlock()
	return len(fake.reorderArgsForCall)
}

// ReorderArgsForCall returns fake arguments for Reorder function
func (fake *WorkItemRepository) ReorderArgsForCall(i int) (context.Context, string, app.WorkItem) {
	fake.reorderMutex.RLock()
	defer fake.reorderMutex.RUnlock()
	return fake.reorderArgsForCall[i].ctx, fake.reorderArgsForCall[i].before, fake.reorderArgsForCall[i].wi
}

// ReorderReturns returns fake values for Reorder function
func (fake *WorkItemRepository) ReorderReturns(result1 *app.WorkItem, result2 error) {
	fake.ReorderStub = nil
	fake.reorderReturns = struct {
		result1 *app.WorkItem
		result2 error
	}{result1, result2}
}

func (fake *WorkItemRepository) Delete(ctx context.Context, ID string) error {
	fake.deleteMutex.Lock()
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		ctx context.Context
		ID  string
	}{ctx, ID})
	fake.recordInvocation("Delete", []interface{}{ctx, ID})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(ctx, ID)
	} else {
		return fake.deleteReturns.result1
	}
}

func (fake *WorkItemRepository) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *WorkItemRepository) DeleteArgsForCall(i int) (context.Context, string) {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return fake.deleteArgsForCall[i].ctx, fake.deleteArgsForCall[i].ID
}

func (fake *WorkItemRepository) DeleteReturns(result1 error) {
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *WorkItemRepository) Create(ctx context.Context, typeID string, fields map[string]interface{}, creator string) (*app.WorkItem, error) {
	fake.createMutex.Lock()
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		ctx    context.Context
		typeID string
		fields map[string]interface{}
	}{ctx, typeID, fields})
	fake.recordInvocation("Create", []interface{}{ctx, typeID, fields})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(ctx, typeID, fields)
	} else {
		return fake.createReturns.result1, fake.createReturns.result2
	}
}

func (fake *WorkItemRepository) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *WorkItemRepository) CreateArgsForCall(i int) (context.Context, string, map[string]interface{}) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].ctx, fake.createArgsForCall[i].typeID, fake.createArgsForCall[i].fields
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
	} else {
		return fake.listReturns.result1, fake.listReturns.result2, fake.listReturns.result3
	}
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
