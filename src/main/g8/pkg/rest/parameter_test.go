package rest

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Gympass/gcore/v3/gtest"
	"github.com/gofrs/uuid"
)

type BaseParamRequestTest struct {
	Name,
	URL,
	Param,
	ExpectedError string
}
type TestStringParamRequest struct {
	BaseParamRequestTest
	ExpectedValue string
}

type TestIntParamRequest struct {
	BaseParamRequestTest
	ExpectedValue int
}

type TestUUIDParamRequest struct {
	BaseParamRequestTest
	ExpectedValue uuid.UUID
}

func TestGetString(t *testing.T) {

	tt := []TestStringParamRequest{
		{
			BaseParamRequestTest: BaseParamRequestTest{
				Name:  "test get string from url with success",
				URL:   "cia?teste=FILTER",
				Param: "teste",
			},
			ExpectedValue: "FILTER",
		},
		{
			BaseParamRequestTest: BaseParamRequestTest{
				Name:          "test error getting string from url parameter does not exist ",
				URL:           "cia?inexistente=abcd",
				Param:         "teste",
				ExpectedError: "param not found",
			},
			ExpectedValue: "",
		},
	}

	for _, testCase := range tt {
		t.Run(testCase.Name, func(t *testing.T) {
			url, err := url.Parse(testCase.URL)
			gtest.AssertNil(t, err)

			val, err := GetString(&http.Request{
				Method: http.MethodGet,
				URL:    url,
			}, testCase.Param)

			validateExpectations(err, testCase.BaseParamRequestTest, t)

			if val != testCase.ExpectedValue {
				t.Fatalf("Expected %s and got %s", testCase.ExpectedValue, val)
			}
		})
	}
}

func TestGetInt(t *testing.T) {

	tt := []TestIntParamRequest{
		{
			BaseParamRequestTest: BaseParamRequestTest{
				Name:  "test get int from url success",
				URL:   "cia?teste=123",
				Param: "teste",
			},
			ExpectedValue: 123,
		},
		{
			BaseParamRequestTest: BaseParamRequestTest{
				Name:          "test error getting int from url parameter does not exist ",
				URL:           "cia?inexistente=abcd",
				Param:         "teste",
				ExpectedError: "param not found",
			},
			ExpectedValue: 0,
		},
	}

	for _, testCase := range tt {
		t.Run(testCase.Name, func(t *testing.T) {
			url, err := url.Parse(testCase.URL)
			gtest.AssertNil(t, err)

			val, err := GetInt(&http.Request{
				Method: http.MethodGet,
				URL:    url,
			}, testCase.Param)

			validateExpectations(err, testCase.BaseParamRequestTest, t)

			if val != testCase.ExpectedValue {
				t.Fatalf("Expected %d and got %d", testCase.ExpectedValue, val)
			}
		})
	}
}

func TestGetUUID(t *testing.T) {

	tt := []TestUUIDParamRequest{
		{
			BaseParamRequestTest: BaseParamRequestTest{
				Name:  "test get uuid from url with success",
				URL:   "cia?id=df7d4231-0c0e-4d80-a72d-98cb276ea754",
				Param: "id",
			},
			ExpectedValue: uuid.Must(uuid.FromString("df7d4231-0c0e-4d80-a72d-98cb276ea754")),
		},
		{
			BaseParamRequestTest: BaseParamRequestTest{
				Name:          "test get uuid error parameter does not exist ",
				URL:           "cia?inexistente=abcd",
				Param:         "teste",
				ExpectedError: "param not found",
			},
			ExpectedValue: uuid.Nil,
		},
	}

	for _, testCase := range tt {
		t.Run(testCase.Name, func(t *testing.T) {
			url, err := url.Parse(testCase.URL)
			gtest.AssertNil(t, err)

			val, err := GetUUID(&http.Request{
				Method: http.MethodGet,
				URL:    url,
			}, testCase.Param)

			validateExpectations(err, testCase.BaseParamRequestTest, t)

			if val != testCase.ExpectedValue {
				t.Fatalf("Expected %v and got %v", testCase.ExpectedValue, val)
			}
		})
	}
}

func validateExpectations(err error, testCase BaseParamRequestTest, t *testing.T) {
	if err != nil && testCase.ExpectedError == "" {
		t.Fatalf("Expected success and got error %s", err.Error())
	}

	if err == nil && testCase.ExpectedError != "" {
		t.Fatalf("Expected error %s and got %s", testCase.ExpectedError, err.Error())
	}
}
