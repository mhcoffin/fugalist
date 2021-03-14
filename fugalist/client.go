package fugalist

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
)

const projectID = "fugalist"

func client() (*firestore.Client, context.Context, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, nil, err
	}
	return client, ctx, nil
}

type Client struct {
	client *firestore.Client
	uid    string
}

func NewClient(uid string) (Client, context.Context, error) {
	client, ctx, err := client()
	if err != nil {
		return Client{}, nil, err
	}
	return Client{client, uid}, ctx, nil
}

func (c *Client) ReadProject(ctx context.Context, pid ProjectId) (*Project, error) {
	path := fmt.Sprintf("Users/%s/Projects/%s", c.uid, pid)
	// fmt.Print(path)
	snap, err := c.client.Doc(path).Get(ctx)
	if err != nil {
		return nil, err
	}
	var p Project
	err = snap.DataTo(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Client) ReadUserInfo(ctx context.Context, uid string) (*UserInfo, error) {
	path := fmt.Sprintf("Users/%s", uid)
	snap, err := c.client.Doc(path).Get(ctx)
	if err != nil {
		return nil, err
	}
	var ui UserInfo
	err = snap.DataTo(&ui)
	if err != nil {
		return nil, err
	}
	return &ui, nil
}
