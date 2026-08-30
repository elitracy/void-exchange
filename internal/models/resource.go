package models

type ResourceId int
type ResourceType string

var currentResourceId ResourceId = 0

const (
	ResourceMineral ResourceType = "resource_mineral"
	ResourceEnergy  ResourceType = "resource_energy"
)

type Resource struct {
	Id   ResourceId
	Type ResourceType
}

type ResourceDeposit struct {
	*CoreEntity
	Resource  *Resource
	Total     int
	Remaining int
}

func NewResource(t ResourceType) *Resource {
	currentResourceId++
	return &Resource{
		Id:   currentResourceId,
		Type: t,
	}
}

func NewResourceDeposit(t ResourceType, total int) *ResourceDeposit {
	return &ResourceDeposit{
		CoreEntity: NewCoreEntity(),
		Resource:   NewResource(t),
		Total:      total,
		Remaining:  total,
	}
}
