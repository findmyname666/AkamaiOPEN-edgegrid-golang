package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// PolicyProperty interface is a cloudlets API interface for policy associated properties
	PolicyProperty interface {
		// GetPolicyProperties gets all the associated properties by the policyID
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#getpolicyproperties
		GetPolicyProperties(context.Context, int64) (GetPolicyPropertiesResponse, error)

		// DeletePolicyProperty removes a property from a policy activation associated_properties list
		DeletePolicyProperty(context.Context, DeletePolicyPropertyRequest) error
	}

	// GetPolicyPropertiesResponse contains response data for GetPolicyProperties
	GetPolicyPropertiesResponse map[string]AssociateProperty

	// AssociateProperty contains the response data for a single property
	AssociateProperty struct {
		GroupID       int64         `json:"groupId"`
		ID            int64         `json:"id"`
		Name          string        `json:"name"`
		NewestVersion NetworkStatus `json:"newestVersion"`
		Production    NetworkStatus `json:"production"`
		Staging       NetworkStatus `json:"staging"`
	}

	// NetworkStatus is the type for NetworkStatus of any activation
	NetworkStatus struct {
		ActivatedBy        string                     `json:"activatedBy"`
		ActivationDate     string                     `json:"activationDate"`
		Version            int64                      `json:"version"`
		CloudletsOrigins   map[string]CloudletsOrigin `json:"cloudletsOrigins"`
		ReferencedPolicies []string                   `json:"referencedPolicies"`
	}

	// CloudletsOrigin is the type for CloudletsOrigins in NetworkStatus
	CloudletsOrigin struct {
		OriginID    string     `json:"id"`
		Hostname    string     `json:"hostname"`
		Type        OriginType `json:"type"`
		Checksum    string     `json:"checksum"`
		Description string     `json:"description"`
	}

	// DeletePolicyPropertyRequest contains the request parameters for DeletePolicyProperty
	DeletePolicyPropertyRequest struct {
		PolicyID   int64
		PropertyID int64
		Network    VersionActivationNetwork
	}
)

var (
	// ErrGetPolicyProperties is returned when GetPolicyProperties fails
	ErrGetPolicyProperties = errors.New("get policy properties")
	// ErrDeletePolicyProperty is returned when DeletePolicyProperty fails
	ErrDeletePolicyProperty = errors.New("delete policy property")
)

// Validate validates DeletePolicyPropertyRequest
func (r DeletePolicyPropertyRequest) Validate() error {
	return validation.Errors{
		"PolicyID":   validation.Validate(r.PolicyID, validation.Required),
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
	}.Filter()
}

// GetPolicyProperties gets all the associated properties by the policyID
func (c *cloudlets) GetPolicyProperties(ctx context.Context, policyID int64) (GetPolicyPropertiesResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetPolicyProperties")

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d/properties", policyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetPolicyProperties, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPolicyProperties, err)
	}

	var result GetPolicyPropertiesResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPolicyProperties, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPolicyProperties, c.Error(resp))
	}

	return result, nil
}

func (c *cloudlets) DeletePolicyProperty(ctx context.Context, params DeletePolicyPropertyRequest) error {
	c.Log(ctx).Debug("DeletePolicyProperty")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeletePolicyProperty, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d/properties/%d", params.PolicyID, params.PropertyID))
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrDeletePolicyProperty, err.Error())
	}

	if params.Network != "" {
		q := uri.Query()
		q.Set("network", string(params.Network))

		uri.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeletePolicyProperty, err)
	}

	resp, err := c.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeletePolicyProperty, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%w: %d", ErrDeletePolicyProperty, resp.StatusCode)
	}

	return nil
}
