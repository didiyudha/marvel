package character

import (
	"context"
	"github.com/bxcodec/faker/v3"
	"github.com/didiyudha/marvel/business/data/tests"
	"testing"
)

var dbc = tests.DBContainer{
	Image: "postgres:13-alpine",
	Port:  "5432",
	Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
}

func TestCharacter(t *testing.T) {
	_, db, teardown := tests.NewDBContainer(t, dbc)
	t.Cleanup(teardown)

	ctx := context.TODO()

	store := NewStore(db)

	t.Log("Save character")
	{
		testID := 0

		characters := []Character{
			{
				ID:          1,
				Name:        faker.Name(),
				Description: faker.Paragraph(),
			},
			{
				ID:          2,
				Name:        faker.Name(),
				Description: faker.Paragraph(),
			},
			{
				ID:          3,
				Name:        faker.Name(),
				Description: faker.Paragraph(),
			},
			{
				ID:          4,
				Name:        faker.Name(),
				Description: faker.Paragraph(),
			},
			{
				ID:          5,
				Name:        faker.Name(),
				Description: faker.Paragraph(),
			},
		}

		if err := store.Save(ctx, characters...); err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to store character(s) : %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to store character(s).", tests.Success, testID)
	}



}