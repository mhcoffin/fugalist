package fugalist

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	firestoreClient *firestore.Client
	ctx             context.Context
	pid             = Uniq()
	createTime      = time.Now().Add(-time.Hour)
	modifyTime      = createTime.Add(time.Hour)
	axis1           = Axis{
		Id:   Uniq(),
		Name: "Length",
		Techniques: []Technique{
			{Id: Uniq(), Name: "Staccato"},
			{Id: Uniq(), Name: "Tenuto"},
		},
		SortOrder: 100,
	}
	axis2 = Axis{
		Id:   Uniq(),
		Name: "Legato",
		Techniques: []Technique{
			{Id: Uniq(), Name: "Normal"},
			{Id: Uniq(), Name: "Legato"},
		},
		SortOrder: 200,
	}
	pigment1 = &Pigment{
		PigmentId: Uniq(),
		Name:      "sus",
		Midi:      "c#1,d2",
		Stop:      "c1",
		Dynamics:  "cc3",
	}
	pigment2 = &Pigment{
		PigmentId: Uniq(),
		Name:      "short",
		Midi:      "c#1,d3",
		Stop:      "c1",
		Dynamics:  "cc3",
	}
	summary1 = ProjectSummary{
		CreateTime:  createTime,
		ProjectID:   pid,
		Version:     0,
		Name:        "Test Project",
		Description: "Test description",
		Plugins:     "Test,plug,ins",
	}
	project1 = Project{
		ProjectId:  pid,
		Public:     false,
		CreateTime: createTime,
		ModifyTime: modifyTime,
		Axes: map[string]Axis{
			axis1.Id: axis1,
			axis2.Id: axis2,
		},
		Pigments: map[PigmentId]*Pigment{
			pigment1.PigmentId: pigment1,
			pigment2.PigmentId: pigment2,
		},
		Palette:      nil,
		Tints:        nil,
		Assignments:  nil,
		URL:          nil,
		URLTimestamp: nil,
		MiddleC:      "",
	}
)

func init() {
	err := os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")
	if err != nil {
		panic(fmt.Errorf("failed to set FIRESTORE_EMULATOR_HOST: %w", err))
	}
	ctx = context.Background()
	firestoreClient, err = client(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to create firestore client: %w", err))
	}
}

func SetUp() string {
	uid := uuid.New().String()
	err := CreateUser(uid, summary1)
	if err != nil {
		panic(err)
	}
	err = CreateProjects(uid, project1)
	if err != nil {
		panic(err)
	}
	return uid
}

func CreateUser(uid string, summaries ...ProjectSummary) error {
	proj := make(map[string]ProjectSummary)
	for _, summary := range summaries {
		proj[summary.ProjectID] = summary
	}
	user := UserInfo{
		Projects: proj,
		Preferences: map[string]string{
			"theme": "dark",
		},
	}
	_, err := firestoreClient.Collection("Users").Doc(uid).Create(ctx, user)
	return err
}

func CreateProjects(uid string, projects ...Project) error {
	for _, p := range projects {
		_, err := firestoreClient.Collection("Users").Doc(uid).Collection("Projects").Doc(p.ProjectId).Create(ctx, project1)
		if err != nil {
			return fmt.Errorf("failed to save project: %w", err)
		}
	}
	return nil
}

func TestClient_ReadUserInfo(t *testing.T) {
	uid := SetUp()
	ctx := context.Background()
	cl, err := NewClient(ctx, uid)
	assert.Nil(t, err)
	x, err := cl.ReadUserInfo(ctx)
	assert.Nil(t, err)
	assert.WithinDuration(t, summary1.CreateTime, x.Projects[summary1.ProjectID].CreateTime, time.Nanosecond)
	assert.Equal(t, summary1.ProjectID, x.Projects[summary1.ProjectID].ProjectID)
	assert.Equal(t, summary1.Name, x.Projects[summary1.ProjectID].Name)
	assert.Equal(t, summary1.Description, x.Projects[summary1.ProjectID].Description)
	assert.Equal(t, summary1.Plugins, x.Projects[summary1.ProjectID].Plugins)
}

func TestClient_ReadProject(t *testing.T) {
	uid := SetUp()
	ctx := context.Background()
	cl, err := NewClient(ctx, uid)
	assert.Nil(t, err)
	p, err := cl.ReadProject(ctx, pid)
	assert.Nil(t, err)
	assert.WithinDuration(t, project1.CreateTime, p.CreateTime, time.Nanosecond)
	assert.WithinDuration(t, project1.ModifyTime, p.ModifyTime, time.Nanosecond)
	assert.Equal(t, project1.Axes, p.Axes)
	assert.Equal(t, project1.Assignments, p.Assignments)
	assert.Equal(t, project1.MiddleC, p.MiddleC)
	assert.Equal(t, project1.Palette, p.Palette)
	assert.Equal(t, project1.Pigments, p.Pigments)
	assert.Equal(t, project1.URL, p.URL)
}

func TestClient_SetUrl(t *testing.T) {
	uid := SetUp()
	ctx := context.Background()
	cl, err := NewClient(ctx, uid)
	assert.Nil(t, err)
	url := "http://foo/bar"
	err = cl.SetUrl(ctx, pid, url)
	assert.Nil(t, err)
	p, err := cl.ReadProject(ctx, pid)
	assert.Nil(t, err)
	assert.Equal(t, &url, p.URL)
	assert.WithinDuration(t, time.Now(), *p.URLTimestamp, 200 * time.Millisecond)

}
