// Code generated by mockery v2.33.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	bun "github.com/uptrace/bun"
)

// Transaction is an autogenerated mock type for the Transaction type
type Transaction struct {
	mock.Mock
}

// Commit provides a mock function with given fields:
func (_m *Transaction) Commit() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetTransaction provides a mock function with given fields:
func (_m *Transaction) GetTransaction() *bun.Tx {
	ret := _m.Called()

	var r0 *bun.Tx
	if rf, ok := ret.Get(0).(func() *bun.Tx); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bun.Tx)
		}
	}

	return r0
}

// Rollback provides a mock function with given fields:
func (_m *Transaction) Rollback() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTransaction creates a new instance of Transaction. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransaction(t interface {
	mock.TestingT
	Cleanup(func())
}) *Transaction {
	mock := &Transaction{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
