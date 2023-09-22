// Package cps provides access to the Akamai CPS APIs
package cps

import (
	"errors"

	"github.com/findmyname666/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
)

var (
	// ErrStructValidation is returned returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// CPS is the cps api interface
	CPS interface {
		ChangeManagementInfo
		ChangeOperations
		Deployments
		DeploymentSchedules
		DVChallenges
		Enrollments
		History
		PostVerification
		PreVerification
		ThirdPartyCSR
	}

	cps struct {
		session.Session
	}

	// Option defines a CPS option
	Option func(*cps)

	// ClientFunc is a cps client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) CPS
)

// Client returns a new cps Client instance with the specified controller
func Client(sess session.Session, opts ...Option) CPS {
	c := &cps{
		Session: sess,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}
