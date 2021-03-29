package fugalist

import (
	"context"
	"github.com/mhcoffin/go-doricolib/doricolib"
	"log"
)

/**
 * Creates a ScoreLib by reading a project from fugalist and putting it all together.
 */
func CreateDoricoLib(uid string, pid string) (*doricolib.ScoreLib, error) {
	ctx := context.Background()
	db, err := NewClient(ctx, uid)
	if err != nil {
		log.Fatalf("failed to create firestore client: %s", err)
	}
	projectSummary, err := db.ReadProjectSummary(ctx, pid)
	if err != nil {
		log.Fatalf("failed to read user project summary: %v", err)
	}
	project, err := db.ReadProject(ctx, pid)
	if err != nil {
		log.Fatalf("failed: %s", err)
	}
	xmap, err := CreateExpressionMap(*projectSummary, project)
	if err != nil {
		log.Fatal(err)
	}
	lib := doricolib.CreateDoricoLib([]doricolib.ExpressionMap{*xmap})
	return lib, nil
}