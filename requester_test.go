package requester

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestResult is the structure used in table test cases
type TestResult struct {
	statusCode int
	body       string
}

// ClientResponse is response which will be returned by the fake server
type ClientResponse struct {
	Message string `json:"message"`
}

func TestMake(t *testing.T) {
	// Json response
	fakeServer1 := newTestServer(200, `{"message":"ok"}`)
	defer fakeServer1.Close()

	requester := New(http.DefaultClient)
	requester.SerUserAgent("test")

	var clientJSONResponse ClientResponse
	response, err := requester.Make("GET", fakeServer1.URL, map[string]string{}, &clientJSONResponse)

	expectedResponse := ClientResponse{Message: "ok"}
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"message":"ok"}`), response.Body)
	assert.Equal(t, 200, response.Status)
	assert.Equal(t, expectedResponse, clientJSONResponse)

	// String response
	fakeServer2 := newTestServer(200, `ok`)
	defer fakeServer2.Close()

	requester = New(http.DefaultClient)
	requester.SerUserAgent("test")

	var clientStringResponse string
	response, err = requester.Make("GET", fakeServer2.URL, map[string]string{}, &clientStringResponse)

	expectedStringResponse := "ok"
	assert.Nil(t, err)
	assert.Equal(t, []byte(expectedStringResponse), response.Body)
	assert.Equal(t, expectedStringResponse, clientStringResponse)
}

func TestMakeError(t *testing.T) {
	testCases := []struct {
		testName           string
		fakeServerResponse TestResult
		method             string

		expectedStatusCode  int
		expectedErr         error
		expectedRawResponse []byte
		expectedResponse    interface{}
	}{
		{
			testName:           "Error Case - Invalid Json",
			fakeServerResponse: TestResult{statusCode: 500, body: `{{dsadsa}}`},
			method:             "GET",

			expectedErr:         errors.New("invalid character '{' looking for beginning of object key string"),
			expectedRawResponse: []byte(`{{dsadsa}}`),
			expectedResponse:    ClientResponse{},
		},
		{
			testName:            "Error Case - Server response code is incorrect",
			fakeServerResponse:  TestResult{statusCode: 500, body: `{"error":"test"}`},
			method:              "GET",
			expectedErr:         nil,
			expectedRawResponse: []byte(`{"error":"test"}`),
			expectedResponse:    ClientResponse{},
		},
	}

	for _, test := range testCases {
		fakeServer := newTestServer(test.fakeServerResponse.statusCode, test.fakeServerResponse.body)
		defer fakeServer.Close()

		requester := New(http.DefaultClient)
		requester.SerUserAgent("test")

		var clientResponse ClientResponse
		response, err := requester.Make("GET", fakeServer.URL, map[string]string{}, &clientResponse)

		assert.Equal(t, test.expectedErr, err)
		assert.Equal(t, test.expectedRawResponse, response.Body)
		assert.Equal(t, test.expectedResponse, clientResponse)
	}
}

func newTestServer(code int, body string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, body)
	}))

	return server
}
