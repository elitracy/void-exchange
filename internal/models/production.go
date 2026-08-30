package models

type ProductType string

const (
	ProductGenerator ProductType = "product_generator"
	ProductDrone     ProductType = "product_drone"
)

var ProductRecipes = map[ProductType]ProductRecipe{
	ProductDrone: {
		ProductType: ProductDrone,
		RequiredResources: map[ResourceType]int{
			ResourceMineral: 5,
		},
	},
	ProductGenerator: {
		ProductType: ProductGenerator,
		RequiredResources: map[ResourceType]int{
			ResourceMineral: 10,
		},
	},
}

type ProductRecipe struct {
	ProductType       ProductType
	RequiredResources map[ResourceType]int
}

type Factory struct {
	*CoreEntity
	TargetProduct ProductType // NOTE: should this be synced up to entities?
	PowerSource   EntityId    // TODO: runtime validation for generator
	Rate          int         // per tick

	CurrentResources map[ResourceType]int
	CurrentUnits     int // just a count, they will be created on the fly once needed
}

type Generator struct {
	*CoreEntity
	OutputTarget EntityId
}

type CreditMiner struct {
	*CoreEntity
	PowerSource    EntityId
	Rate           int // per tick
	CurrentCredits int
}

func NewFactory(targetProduct ProductType) *Factory {
	f := &Factory{
		CoreEntity:       NewCoreEntity(),
		TargetProduct:    targetProduct,
		PowerSource:      -1,
		Rate:             1,
		CurrentResources: make(map[ResourceType]int),
	}

	for r := range ProductRecipes[f.TargetProduct].RequiredResources {
		f.CurrentResources[r] = 0
	}

	return f
}

func NewGenerator() *Generator {
	return &Generator{
		CoreEntity: NewCoreEntity(),
	}
}

func NewCreditMiner() *CreditMiner {
	return &CreditMiner{
		CoreEntity: NewCoreEntity(),
		Rate:       100,
	}
}

func (g *Generator) SetOutputTarget(id EntityId) { g.OutputTarget = id }

func (f *Factory) SetGenerator(id EntityId) { f.PowerSource = id }

// TODO: refactor flags to tickChecks() for separation of concerns
func (f *Factory) Tick() error {

	// check power
	var powerFlag bool = f.PowerSource != -1

	// check resources
	var resourceFlag bool = true

	recipe := ProductRecipes[f.TargetProduct]
	for resource, resourceCount := range f.CurrentResources {
		if resourceCount < recipe.RequiredResources[resource] {
			resourceFlag = false
		}
	}

	if powerFlag && resourceFlag {
		f.CurrentUnits += f.Rate
		for resource, n := range recipe.RequiredResources {
			f.CurrentResources[resource] -= (n * f.Rate)
		}
	}

	return nil
}

func (cm *CreditMiner) SetGenerator(id EntityId) { cm.PowerSource = id }

// TODO: refactor flags to tickChecks() for separation of concerns
func (cm *CreditMiner) Tick() error {

	// check power
	var powerFlag bool = cm.PowerSource != -1

	if powerFlag {
		cm.CurrentCredits += cm.Rate
	}

	return nil
}
