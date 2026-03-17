package functional

import "io"

// InputFunctionalTest is the input for the functional test
type InputFunctionalTest struct {
	// desc is a description of the test
	desc string
	// expectedBody is the expected body of the response
	expectedBody string
	// expectedHttpHeaders is the expected headers of the response
	expectedHttpHeaders map[string]string
	// expectedStatusCode is the expected status code of the response
	expectedStatusCode int
	// headers is the expected headers of the request
	headers map[string]string
	// method is the HTTP method to use in the test
	method string
	// parameters is the request body
	parameters io.ReadCloser
	// url is the URL to test
	url string
}
