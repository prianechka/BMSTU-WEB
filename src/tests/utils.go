package tests

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
)

func AssertErrors(t *testing.T, realError, expectedError error) {
	if realError != expectedError {
		t.Errorf("unexpected err: %v", realError)
	}
}

func AssertMocks(t *testing.T, mock sqlmock.Sqlmock) {
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
	}
}

func AssertResult(t *testing.T, realResult any, expectedResult any) {
	if !reflect.DeepEqual(realResult, expectedResult) {
		t.Errorf("results not match, want %v, have %v", expectedResult, realResult)
	}
}
