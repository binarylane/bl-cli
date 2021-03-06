package binarylane

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const vpcsBasePath = "/v2/vpcs"

// VPCsService is an interface for managing Virtual Private Cloud configurations with the
// BinaryLane API.
// See: https://api.binarylane.com.au/reference#vpcs
type VPCsService interface {
	Create(context.Context, *VPCCreateRequest) (*VPC, *Response, error)
	Get(context.Context, int) (*VPC, *Response, error)
	List(context.Context, *ListOptions) ([]*VPC, *Response, error)
	Update(context.Context, int, *VPCUpdateRequest) (*VPC, *Response, error)
	Set(context.Context, int, ...VPCSetField) (*VPC, *Response, error)
	Delete(context.Context, int) (*Response, error)
}

var _ VPCsService = &VPCsServiceOp{}

// VPCsServiceOp interfaces with VPC endpoints in the BinaryLane API.
type VPCsServiceOp struct {
	client *Client
}

// VPCCreateRequest represents a request to create a Virtual Private Cloud.
type VPCCreateRequest struct {
	Name        string `json:"name,omitempty"`
	RegionSlug  string `json:"region,omitempty"`
	Description string `json:"description,omitempty"`
	IPRange     string `json:"ip_range,omitempty"`
}

// VPCUpdateRequest represents a request to update a Virtual Private Cloud.
type VPCUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Default     *bool  `json:"default,omitempty"`
}

// VPCSetField allows one to set individual fields within a VPC configuration.
type VPCSetField interface {
	vpcSetField(map[string]interface{})
}

// VPCSetName is used when one want to set the `name` field of a VPC.
// Ex.: VPCs.Set(..., VPCSetName("new-name"))
type VPCSetName string

// VPCSetDescription is used when one want to set the `description` field of a VPC.
// Ex.: VPCs.Set(..., VPCSetDescription("vpc description"))
type VPCSetDescription string

// VPCSetDefault is used when one wants to enable the `default` field of a VPC, to
// set a VPC as the default one in the region
// Ex.: VPCs.Set(..., VPCSetDefault())
func VPCSetDefault() VPCSetField {
	return &vpcSetDefault{}
}

// vpcSetDefault satisfies the VPCSetField interface
type vpcSetDefault struct{}

// VPC represents a BinaryLane Virtual Private Cloud configuration.
type VPC struct {
	ID          int       `json:"id,,float64,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	IPRange     string    `json:"ip_range,omitempty"`
	RegionSlug  string    `json:"region,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	Default     bool      `json:"default,omitempty"`
}

type vpcRoot struct {
	VPC *VPC `json:"vpc"`
}

type vpcsRoot struct {
	VPCs  []*VPC `json:"vpcs"`
	Links *Links `json:"links"`
	Meta  *Meta  `json:"meta"`
}

func (v VPC) URN() string {
	return ToURN("VPC", v.ID)
}

// Get returns the details of a Virtual Private Cloud.
func (v *VPCsServiceOp) Get(ctx context.Context, id int) (*VPC, *Response, error) {
	path := fmt.Sprintf("%s/%d", vpcsBasePath, id)
	req, err := v.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(vpcRoot)
	resp, err := v.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.VPC, resp, nil
}

// Create creates a new Virtual Private Cloud.
func (v *VPCsServiceOp) Create(ctx context.Context, create *VPCCreateRequest) (*VPC, *Response, error) {
	path := vpcsBasePath
	req, err := v.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, nil, err
	}

	root := new(vpcRoot)
	resp, err := v.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.VPC, resp, nil
}

// List returns a list of the caller's VPCs, with optional pagination.
func (v *VPCsServiceOp) List(ctx context.Context, opt *ListOptions) ([]*VPC, *Response, error) {
	path, err := addOptions(vpcsBasePath, opt)
	if err != nil {
		return nil, nil, err
	}
	req, err := v.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(vpcsRoot)
	resp, err := v.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.VPCs, resp, nil
}

// Update updates a Virtual Private Cloud's properties.
func (v *VPCsServiceOp) Update(ctx context.Context, id int, update *VPCUpdateRequest) (*VPC, *Response, error) {
	path := fmt.Sprintf("%s/%d", vpcsBasePath, id)
	req, err := v.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, nil, err
	}

	root := new(vpcRoot)
	resp, err := v.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.VPC, resp, nil
}

func (n VPCSetName) vpcSetField(in map[string]interface{}) {
	in["name"] = n
}

func (n VPCSetDescription) vpcSetField(in map[string]interface{}) {
	in["description"] = n
}

func (*vpcSetDefault) vpcSetField(in map[string]interface{}) {
	in["default"] = true
}

// Set updates specific properties of a Virtual Private Cloud.
func (v *VPCsServiceOp) Set(ctx context.Context, id int, fields ...VPCSetField) (*VPC, *Response, error) {
	path := fmt.Sprintf("%s/%d", vpcsBasePath, id)
	update := make(map[string]interface{}, len(fields))
	for _, field := range fields {
		field.vpcSetField(update)
	}

	req, err := v.client.NewRequest(ctx, http.MethodPatch, path, update)
	if err != nil {
		return nil, nil, err
	}

	root := new(vpcRoot)
	resp, err := v.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.VPC, resp, nil
}

// Delete deletes a Virtual Private Cloud. There is no way to recover a VPC once it has been
// destroyed.
func (v *VPCsServiceOp) Delete(ctx context.Context, id int) (*Response, error) {
	path := fmt.Sprintf("%s/%d", vpcsBasePath, id)
	req, err := v.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := v.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
