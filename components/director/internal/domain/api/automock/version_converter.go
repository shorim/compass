// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	version "github.com/kyma-incubator/compass/components/director/internal/domain/version"
)

// VersionConverter is an autogenerated mock type for the VersionConverter type
type VersionConverter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: _a0
func (_m *VersionConverter) FromEntity(_a0 version.Version) *model.Version {
	ret := _m.Called(_a0)

	var r0 *model.Version
	if rf, ok := ret.Get(0).(func(version.Version) *model.Version); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Version)
		}
	}

	return r0
}

// InputFromGraphQL provides a mock function with given fields: in
func (_m *VersionConverter) InputFromGraphQL(in *externalschema.VersionInput) *model.VersionInput {
	ret := _m.Called(in)

	var r0 *model.VersionInput
	if rf, ok := ret.Get(0).(func(*externalschema.VersionInput) *model.VersionInput); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.VersionInput)
		}
	}

	return r0
}

// ToEntity provides a mock function with given fields: _a0
func (_m *VersionConverter) ToEntity(_a0 model.Version) version.Version {
	ret := _m.Called(_a0)

	var r0 version.Version
	if rf, ok := ret.Get(0).(func(model.Version) version.Version); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(version.Version)
	}

	return r0
}

// ToGraphQL provides a mock function with given fields: in
func (_m *VersionConverter) ToGraphQL(in *model.Version) *externalschema.Version {
	ret := _m.Called(in)

	var r0 *externalschema.Version
	if rf, ok := ret.Get(0).(func(*model.Version) *externalschema.Version); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*externalschema.Version)
		}
	}

	return r0
}
