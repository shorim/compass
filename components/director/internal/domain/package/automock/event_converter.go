// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
)
import mock "github.com/stretchr/testify/mock"
import model "github.com/kyma-incubator/compass/components/director/internal/model"

// EventConverter is an autogenerated mock type for the EventConverter type
type EventConverter struct {
	mock.Mock
}

// MultipleInputFromGraphQL provides a mock function with given fields: in
func (_m *EventConverter) MultipleInputFromGraphQL(in []*externalschema.EventDefinitionInput) ([]*model.EventDefinitionInput, error) {
	ret := _m.Called(in)

	var r0 []*model.EventDefinitionInput
	if rf, ok := ret.Get(0).(func([]*externalschema.EventDefinitionInput) []*model.EventDefinitionInput); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.EventDefinitionInput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*externalschema.EventDefinitionInput) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MultipleToGraphQL provides a mock function with given fields: in
func (_m *EventConverter) MultipleToGraphQL(in []*model.EventDefinition) []*externalschema.EventDefinition {
	ret := _m.Called(in)

	var r0 []*externalschema.EventDefinition
	if rf, ok := ret.Get(0).(func([]*model.EventDefinition) []*externalschema.EventDefinition); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*externalschema.EventDefinition)
		}
	}

	return r0
}

// ToGraphQL provides a mock function with given fields: in
func (_m *EventConverter) ToGraphQL(in *model.EventDefinition) *externalschema.EventDefinition {
	ret := _m.Called(in)

	var r0 *externalschema.EventDefinition
	if rf, ok := ret.Get(0).(func(*model.EventDefinition) *externalschema.EventDefinition); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*externalschema.EventDefinition)
		}
	}

	return r0
}
