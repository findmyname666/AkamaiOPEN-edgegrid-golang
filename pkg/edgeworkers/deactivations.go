package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Deactivations is an EdgeWorkers deactivations API interface
	Deactivations interface {
		// ListDeactivations lists all deactivations for a given EdgeWorker ID
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/deactivations#get-deactivations-1
		ListDeactivations(context.Context, EdgeWorkerListDeactivationsRequest) (*EdgeWorkerListDeactivationsResponse, error)

		// DeactivateVersion deactivates an existing EdgeWorker version on the akamai network
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/deactivations#post-deactivations-1
		DeactivateVersion(context.Context, EdgeWorkerDeactivateVersionRequest) (*Deactivation, error)

		// GetDeactivation gets details for a specific deactivation
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/deactivations#get-deactivation-1
		GetDeactivation(context.Context, EdgeWorkerGetDeactivationRequest) (*Deactivation, error)
	}

	// Deactivation is the response returned by GetDeactivation, DeactivateVersion and ListDeactivation
	Deactivation struct {
		EdgeWorkerID     int               `json:"edgeWorkerId"`
		Version          string            `json:"version"`
		DeactivationID   int               `json:"deactivationId"`
		AccountID        string            `json:"accountId"`
		Status           string            `json:"status"`
		Network          ActivationNetwork `json:"network"`
		Note             string            `json:"note,omitempty"`
		CreatedBy        string            `json:"createdBy"`
		CreatedTime      string            `json:"createdTime"`
		LastModifiedTime string            `json:"lastModifiedTime"`
	}

	// EdgeWorkerListDeactivationsRequest describes the parameters for the list deactivations request
	EdgeWorkerListDeactivationsRequest struct {
		EdgeWorkerID int
		Version      string
	}

	// EdgeWorkerDeactivateVersionRequest describes the request parameters for DeactivateVersion
	EdgeWorkerDeactivateVersionRequest struct {
		EdgeWorkerID int
		Body         EdgeWorkerDeactivateVersionPayload
	}

	// EdgeWorkerGetDeactivationRequest describes the request parameters for GetDeactivation
	EdgeWorkerGetDeactivationRequest struct {
		EdgeWorkerID   int
		DeactivationID int
	}

	// EdgeWorkerDeactivateVersionPayload is the request payload for DeactivateVersion
	EdgeWorkerDeactivateVersionPayload struct {
		Network ActivationNetwork `json:"network"`
		Note    string            `json:"note"`
		Version string            `json:"version"`
	}

	// EdgeWorkerListDeactivationsResponse describes the list deactivations response
	EdgeWorkerListDeactivationsResponse struct {
		Deactivations []Deactivation `json:"deactivations"`
	}
)

// Validate validates EdgeWorkerListDeactivationsRequest
func (r *EdgeWorkerListDeactivationsRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(r.EdgeWorkerID, validation.Required),
	}.Filter()
}

// Validate validates EdgeWorkerDeactivateVersionRequest
func (r *EdgeWorkerDeactivateVersionRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(r.EdgeWorkerID, validation.Required),
		"Body.Network": validation.Validate(r.Body.Network, validation.Required, validation.In(
			ActivationNetworkProduction, ActivationNetworkStaging,
		).Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'",
			r.Body.Network, ActivationNetworkStaging, ActivationNetworkProduction))),
		"Body.Version": validation.Validate(r.Body.Version, validation.Required),
	}.Filter()
}

// Validate validates EdgeWorkerGetDeactivationRequest
func (r *EdgeWorkerGetDeactivationRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID":   validation.Validate(r.EdgeWorkerID, validation.Required),
		"DeactivationID": validation.Validate(r.DeactivationID, validation.Required),
	}.Filter()
}

var (
	// ErrListDeactivations is returned when ListDeactivations fails
	ErrListDeactivations = errors.New("list deactivations")
	// ErrDeactivateVersion is returned when DeactivateVersion fails
	ErrDeactivateVersion = errors.New("deactivate version")
	// ErrGetDeactivation is returned when GetDeactivation fails
	ErrGetDeactivation = errors.New("get deactivation")
)

func (e *edgeworkers) ListDeactivations(ctx context.Context, params EdgeWorkerListDeactivationsRequest) (*EdgeWorkerListDeactivationsResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListDeactivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListDeactivations, ErrStructValidation, err.Error())
	}

	uri, err := url.Parse(fmt.Sprintf("/edgeworkers/v1/ids/%d/deactivations", params.EdgeWorkerID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse URL: %s", ErrListDeactivations, err.Error())
	}

	q := uri.Query()
	if params.Version != "" {
		q.Add("version", params.Version)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListDeactivations, err.Error())
	}

	var result EdgeWorkerListDeactivationsResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListDeactivations, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListDeactivations, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) DeactivateVersion(ctx context.Context, params EdgeWorkerDeactivateVersionRequest) (*Deactivation, error) {
	logger := e.Log(ctx)
	logger.Debug("DeactivateVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeactivateVersion, ErrStructValidation, err.Error())
	}

	uri, err := url.Parse(fmt.Sprintf("/edgeworkers/v1/ids/%d/deactivations", params.EdgeWorkerID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse URL: %s", ErrDeactivateVersion, err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeactivateVersion, err.Error())
	}

	var result Deactivation
	resp, err := e.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeactivateVersion, err.Error())
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrDeactivateVersion, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) GetDeactivation(ctx context.Context, params EdgeWorkerGetDeactivationRequest) (*Deactivation, error) {
	logger := e.Log(ctx)
	logger.Debug("GetDeactivation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetDeactivation, ErrStructValidation, err.Error())
	}

	uri, err := url.Parse(fmt.Sprintf("/edgeworkers/v1/ids/%d/deactivations/%d", params.EdgeWorkerID, params.DeactivationID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse URL: %s", ErrGetDeactivation, err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)

	var result Deactivation
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetDeactivation, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetDeactivation, e.Error(resp))
	}

	return &result, nil
}
