package binarylane

import (
	"context"
	"fmt"
	"net/http"
)

const floatingBasePath = "v2/floating_ips"

// FloatingIPsService is an interface for interfacing with the floating IPs
// endpoints of the BinaryLane API.
// See: https://api.binarylane.com.au/reference#floating-ips
type FloatingIPsService interface {
	List(context.Context, *ListOptions) ([]FloatingIP, *Response, error)
	Get(context.Context, string) (*FloatingIP, *Response, error)
	Create(context.Context, *FloatingIPCreateRequest) (*FloatingIP, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

// FloatingIPsServiceOp handles communication with the floating IPs related methods of the
// BinaryLane API.
type FloatingIPsServiceOp struct {
	client *Client
}

var _ FloatingIPsService = &FloatingIPsServiceOp{}

// FloatingIP represents a BinaryLane floating IP.
type FloatingIP struct {
	Region *Region `json:"region"`
	Server *Server `json:"server"`
	IP     string  `json:"ip"`
}

func (f FloatingIP) String() string {
	return Stringify(f)
}

// URN returns the floating IP in a valid BL API URN form.
func (f FloatingIP) URN() string {
	return ToURN("FloatingIP", f.IP)
}

type floatingIPsRoot struct {
	FloatingIPs []FloatingIP `json:"floating_ips"`
	Links       *Links       `json:"links"`
	Meta        *Meta        `json:"meta"`
}

type floatingIPRoot struct {
	FloatingIP *FloatingIP `json:"floating_ip"`
	Links      *Links      `json:"links,omitempty"`
}

// FloatingIPCreateRequest represents a request to create a floating IP.
// If ServerId is not empty, the floating IP will be assigned to the server.
type FloatingIPCreateRequest struct {
	Region   string `json:"region"`
	ServerID int    `json:"server_id,omitempty"`
}

// List all floating IPs.
func (f *FloatingIPsServiceOp) List(ctx context.Context, opt *ListOptions) ([]FloatingIP, *Response, error) {
	path := floatingBasePath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := f.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(floatingIPsRoot)
	resp, err := f.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.FloatingIPs, resp, err
}

// Get an individual floating IP.
func (f *FloatingIPsServiceOp) Get(ctx context.Context, ip string) (*FloatingIP, *Response, error) {
	path := fmt.Sprintf("%s/%s", floatingBasePath, ip)

	req, err := f.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(floatingIPRoot)
	resp, err := f.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.FloatingIP, resp, err
}

// Create a floating IP. If the ServerID field of the request is not empty,
// the floating IP will also be assigned to the server.
func (f *FloatingIPsServiceOp) Create(ctx context.Context, createRequest *FloatingIPCreateRequest) (*FloatingIP, *Response, error) {
	path := floatingBasePath

	req, err := f.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(floatingIPRoot)
	resp, err := f.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.FloatingIP, resp, err
}

// Delete a floating IP.
func (f *FloatingIPsServiceOp) Delete(ctx context.Context, ip string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", floatingBasePath, ip)

	req, err := f.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.client.Do(ctx, req, nil)

	return resp, err
}
