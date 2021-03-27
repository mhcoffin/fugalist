package fugalist

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
)

const projectID = "fugalist"

func client(ctx context.Context) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type Client struct {
	client *firestore.Client
	uid    string
}

func NewClient(ctx context.Context, uid string) (Client, error) {
	client, err := client(ctx)
	if err != nil {
		return Client{}, err
	}
	return Client{client, uid}, nil
}

func (c *Client) ReadProject(ctx context.Context, pid ProjectId) (*Project, error) {
	snap, err := c.client.Collection("Users").Doc(c.uid).Collection("Projects").Doc(pid).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read document: %w", err)
	}
	var p Project
	err = snap.DataTo(&p)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}
	return &p, nil
}

func (c *Client) ReadUserInfo(ctx context.Context) (*UserInfo, error) {
	snap, err := c.client.Collection("Users").Doc(c.uid).Get(ctx)
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

// Set the URL and the URLTimestamp of a project
func (c *Client) SetUrl(ctx context.Context, pid string, url string) error {
	project := c.client.Collection("Users").Doc(c.uid).Collection("Projects").Doc(pid)
	_, err := project.Update(ctx, []firestore.Update{
		{
			Path:  "URL",
			Value: url,
		},
		{
			Path:  "URLTimestamp",
			Value: firestore.ServerTimestamp,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to set url: %w", err)
	}
	return nil
}

func (c *Client) WriteShare(ctx context.Context, share Share) error {
	shareDoc := c.client.Collection("Shared").Doc(fmt.Sprintf("%s.%d", share.PID, share.Summary.Version))
	userDoc := c.client.Collection("Users").Doc(share.UID)
	batch := c.client.Batch()
	batch.Set(shareDoc, share)
	batch.Update(userDoc, []firestore.Update{
		{
			Path:  fmt.Sprintf("Projects.%s.SharedTime", share.PID),
			Value: firestore.ServerTimestamp,
		},
	})
	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to share project %s.%s.%d: %w", share.UID, share.PID, share.Summary.Version, err)
	}
	return nil
}
