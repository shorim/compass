// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
)
import mock "github.com/stretchr/testify/mock"
import model "github.com/kyma-incubator/compass/components/director/internal/model"

// PackageConverter is an autogenerated mock type for the PackageConverter type
type PackageConverter struct {
	mock.Mock
}

// MultipleCreateInputFromGraphQL provides a mock function with given fields: in
func (_m *PackageConverter) MultipleCreateInputFromGraphQL(in []*externalschema.PackageCreateInput) ([]*model.PackageCreateInput, error) {
	ret := _m.Called(in)

	var r0 []*model.PackageCreateInput
	if rf, ok := ret.Get(0).(func([]*externalschema.PackageCreateInput) []*model.PackageCreateInput); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.PackageCreateInput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*externalschema.PackageCreateInput) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MultipleToGraphQL provides a mock function with given fields: in
func (_m *PackageConverter) MultipleToGraphQL(in []*model.Package) ([]*externalschema.Package, error) {
	ret := _m.Called(in)

	var r0 []*externalschema.Package
	if rf, ok := ret.Get(0).(func([]*model.Package) []*externalschema.Package); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*externalschema.Package)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*model.Package) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToGraphQL provides a mock function with given fields: in
func (_m *PackageConverter) ToGraphQL(in *model.Package) (*externalschema.Package, error) {
	ret := _m.Called(in)

	var r0 *externalschema.Package
	if rf, ok := ret.Get(0).(func(*model.Package) *externalschema.Package); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*externalschema.Package)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Package) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
