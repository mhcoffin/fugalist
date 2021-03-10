package fugalistdb

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
)

const uid = "VRf7soDS0BQ6praLnktgJfD5CVa2"
const pid = "8I3dGF1qFu"

func TestClient(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "fugalist")
	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf("users/%s/projects/%s", uid, pid)
	snap, err := client.Doc(path).Get(ctx)
	assert.Nilf(t, err, "%v", err)

	assignments, err := snap.DataAt("assignments")
	assert.Nil(t, err)
	assert.NotNil(t, assignments)

	errcnt := 0
	for k := 0; k < 4; k++ {
		var p = Project{}
		err = snap.DataTo(&p)
		assert.Nil(t, err)
		if p.Assignments == nil {
			errcnt++
		}
	}
	assert.Equal(t, 0, errcnt)
}
