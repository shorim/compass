// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
)
import mock "github.com/stretchr/testify/mock"
import model "github.com/kyma-incubator/compass/components/director/internal/model"

// TokenConverter is an autogenerated mock type for the TokenConverter type
type TokenConverter struct {
	mock.Mock
}

// ToGraphQLForApplication provides a mock function with given fields: _a0
func (_m *TokenConverter) ToGraphQLForApplication(_a0 model.OneTimeToken) (externalschema.OneTimeTokenForApplication, error) {
	ret := _m.Called(_a0)

	var r0 externalschema.OneTimeTokenForApplication
	if rf, ok := ret.Get(0).(func(model.OneTimeToken) externalschema.OneTimeTokenForApplication); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(externalschema.OneTimeTokenForApplication)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.OneTimeToken) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToGraphQLForRuntime provides a mock function with given fields: _a0
func (_m *TokenConverter) ToGraphQLForRuntime(_a0 model.OneTimeToken) externalschema.OneTimeTokenForRuntime {
	ret := _m.Called(_a0)

	var r0 externalschema.OneTimeTokenForRuntime
	if rf, ok := ret.Get(0).(func(model.OneTimeToken) externalschema.OneTimeTokenForRuntime); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(externalschema.OneTimeTokenForRuntime)
	}

	return r0
}
