// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
)
import mock "github.com/stretchr/testify/mock"
import model "github.com/kyma-incubator/compass/components/director/internal/model"

// APIConverter is an autogenerated mock type for the APIConverter type
type APIConverter struct {
	mock.Mock
}

// InputFromGraphQL provides a mock function with given fields: in
func (_m *APIConverter) InputFromGraphQL(in *externalschema.APIDefinitionInput) (*model.APIDefinitionInput, error) {
	ret := _m.Called(in)

	var r0 *model.APIDefinitionInput
	if rf, ok := ret.Get(0).(func(*externalschema.APIDefinitionInput) *model.APIDefinitionInput); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.APIDefinitionInput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*externalschema.APIDefinitionInput) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MultipleInputFromGraphQL provides a mock function with given fields: in
func (_m *APIConverter) MultipleInputFromGraphQL(in []*externalschema.APIDefinitionInput) ([]*model.APIDefinitionInput, error) {
	ret := _m.Called(in)

	var r0 []*model.APIDefinitionInput
	if rf, ok := ret.Get(0).(func([]*externalschema.APIDefinitionInput) []*model.APIDefinitionInput); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.APIDefinitionInput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*externalschema.APIDefinitionInput) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MultipleToGraphQL provides a mock function with given fields: in
func (_m *APIConverter) MultipleToGraphQL(in []*model.APIDefinition) []*externalschema.APIDefinition {
	ret := _m.Called(in)

	var r0 []*externalschema.APIDefinition
	if rf, ok := ret.Get(0).(func([]*model.APIDefinition) []*externalschema.APIDefinition); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*externalschema.APIDefinition)
		}
	}

	return r0
}

// SpecToGraphQL provides a mock function with given fields: definitionID, in
func (_m *APIConverter) SpecToGraphQL(definitionID string, in *model.APISpec) *externalschema.APISpec {
	ret := _m.Called(definitionID, in)

	var r0 *externalschema.APISpec
	if rf, ok := ret.Get(0).(func(string, *model.APISpec) *externalschema.APISpec); ok {
		r0 = rf(definitionID, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*externalschema.APISpec)
		}
	}

	return r0
}

// ToGraphQL provides a mock function with given fields: in
func (_m *APIConverter) ToGraphQL(in *model.APIDefinition) *externalschema.APIDefinition {
	ret := _m.Called(in)

	var r0 *externalschema.APIDefinition
	if rf, ok := ret.Get(0).(func(*model.APIDefinition) *externalschema.APIDefinition); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*externalschema.APIDefinition)
		}
	}

	return r0
}
