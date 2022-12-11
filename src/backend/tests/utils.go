package tests

import (
	"errors"
	"github.com/bloomberg/go-testgroup"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"log"
	"reflect"
	"time"
)

var TestErr = errors.New("error")

func AssertErrors(t *testgroup.T, realError, expectedError error) {
	if realError != expectedError {
		t.Errorf(TestErr, "unexpected err: %v", realError)
	}
}

func AssertMocks(t *testgroup.T, mock sqlmock.Sqlmock) {
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf(TestErr, "there were unfulfilled expectations: %s", expectationErr)
	}
}

func AssertResult(t *testgroup.T, realResult any, expectedResult any) {
	if !reflect.DeepEqual(realResult, expectedResult) {
		t.Errorf(TestErr, "results not match, want %v, have %v", expectedResult, realResult)
	}
}

func TimeTrack(start time.Time) {
	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)
}
