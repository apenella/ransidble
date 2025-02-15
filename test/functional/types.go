package functional

import "io"

// InputFunctionalTest is the input for the functional test
type InputFunctionalTest struct {
	// desc is a description of the test
	desc string
	// method is the HTTP method to use in the test
	method string
	// url is the URL to test
	url string
	// parameters is the request body
	parameters io.ReadCloser
	// expectedStatusCode is the expected status code of the response
	expectedStatusCode int
	// headers is the expected headers of the response
	headers map[string]string
	// expectedBody is the expected body of the response
	expectedBody string
}
