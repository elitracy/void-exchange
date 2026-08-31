package worldgen

import "github.com/elitracy/space-war-sim/internal/models"

func PopulateTerritory(
	em *models.EntityManager,
	t *models.Territory,
	types []models.ResourceType,
	min, max int,
	amount func(min, max int) int,
) error {
	for _, rt := range types {
		deposit := models.NewResourceDeposit(rt, amount(min, max))
		models.Register(em, deposit)

		if err := t.AddDeposit(deposit.Id()); err != nil {
			return err
		}
	}

	return nil
}
