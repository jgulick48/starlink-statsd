// Package starlink talks to the dish over gRPC. Only read-only requests
// (get_status, get_history) are ever issued.
package starlink

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device"
)

const handleMethod = "/SpaceX.API.Device.Device/Handle"

// Client is a read-only dish client.
type Client struct {
	conn *grpc.ClientConn
}

// New creates a client. The connection is established lazily.
func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("create grpc client: %w", err)
	}
	return &Client{conn: conn}, nil
}

// Close closes the underlying connection.
func (c *Client) Close() error { return c.conn.Close() }

func (c *Client) handle(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	resp := &pb.Response{}
	if err := c.conn.Invoke(ctx, handleMethod, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Status performs a get_status request.
func (c *Client) Status(ctx context.Context) (*pb.DishGetStatusResponse, error) {
	resp, err := c.handle(ctx, &pb.Request{Request: &pb.Request_GetStatus{GetStatus: &pb.GetStatusRequest{}}})
	if err != nil {
		return nil, err
	}
	st := resp.GetDishGetStatus()
	if st == nil {
		return nil, fmt.Errorf("get_status: response had no dish status")
	}
	return st, nil
}

// History performs a get_history request.
func (c *Client) History(ctx context.Context) (*pb.DishGetHistoryResponse, error) {
	resp, err := c.handle(ctx, &pb.Request{Request: &pb.Request_GetHistory{GetHistory: &pb.GetHistoryRequest{}}})
	if err != nil {
		return nil, err
	}
	h := resp.GetDishGetHistory()
	if h == nil {
		return nil, fmt.Errorf("get_history: response had no dish history")
	}
	return h, nil
}
