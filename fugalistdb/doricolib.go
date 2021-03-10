package fugalistdb

import (
	"github.com/mhcoffin/go-doricolib/doricolib"
	"log"
)

/**
 * Creates a ScoreLib by reading a project from fugalist and putting it all together.
 */
func CreateDoricoLib(uid string, pid string, version int) (*doricolib.ScoreLib, error) {
	db, ctx, err := NewClient(uid)
	if err != nil {
		log.Fatalf("failed to create firestore client: %s", err)
	}
	userInfo, err := db.ReadUserInfo(ctx, uid)
	if err != nil {
		log.Fatalf("failed to read user info: %v", err)
	}
	projectSummary, ok := userInfo.Projects[pid]
	if !ok {
		log.Fatalf("summary does not contain project %v", pid)
	}
	project, err := db.ReadProject(ctx, pid, version)
	if err != nil {
		log.Fatalf("failed: %s", err)
	}
	xmap, err := CreateExpressionMap(projectSummary, project)
	if err != nil {
		log.Fatal(err)
	}
	lib := doricolib.CreateDoricoLib([]doricolib.ExpressionMap{*xmap})
	return lib, nil
}