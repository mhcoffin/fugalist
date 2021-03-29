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

func (c *Client) ReadProjectSummary(ctx context.Context, pid string) (*ProjectSummary, error) {
	snap, err := c.client.Collection("Users").Doc(c.uid).Collection("Summaries").Doc(pid).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read project summary for %s.%s: %w", c.uid, pid, err)
	}
	result := &ProjectSummary{}
	err = snap.DataTo(result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall project summary for %s.%s: %w", c.uid, pid, err)
	}
	return result, nil
}

// After generating an expression map, set the URL and timestamp in the summary.
func (c *Client) SetUrl(ctx context.Context, pid string, url string) error {
	summary := c.client.Collection("Users").Doc(c.uid).Collection("Summaries").Doc(pid)
	_, err := summary.Update(ctx, []firestore.Update{
		{
			Path:  "ExpressionMapURL",
			Value: url,
		},
		{
			Path:  "ExpressionMapTime",
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
	summaryDoc := c.client.Collection("Users").Doc(share.UID).Collection("Summaries").Doc(share.PID)
	batch := c.client.Batch()
	batch.Set(shareDoc, share)
	batch.Update(summaryDoc, []firestore.Update{
		{
			Path:  "ShareTime",
			Value: firestore.ServerTimestamp,
		},
		{
			Path: "Version",
			Value: firestore.Increment(1),
		},
	})
	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to share project %s.%s.%d: %w", share.UID, share.PID, share.Summary.Version, err)
	}
	return nil
}
