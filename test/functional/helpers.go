package functional

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// waitHTTPServer waits for the HTTP server to be ready
func waitHTTPServer(listenAddress string, sleepTime time.Duration, retries int) error {
	var err error
	var conn net.Conn

	for i := 0; i < retries; i++ {
		conn, err = net.DialTimeout("tcp", listenAddress, sleepTime)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(sleepTime)
	}

	return errors.New("HTTP server is not ready")
}

// actAndAssert performs the request and validates the response
func actAndAssert(t *testing.T, input *InputFunctionalTest) error {

	var err error
	var httpReq *http.Request
	var httpResp *http.Response
	var body []byte

	httpReq, err = http.NewRequest(input.method, input.url, input.parameters)
	if err != nil {
		return fmt.Errorf("%s. error creating HTTP request: %s", t.Name(), err)
	}

	for key, value := range input.headers {
		httpReq.Header.Set(key, value)
	}

	passed := t.Run(fmt.Sprintf("Functional %s", input.desc), func(t *testing.T) {

		client := &http.Client{}
		httpResp, err = client.Do(httpReq)
		if err != nil {
			t.Errorf("%s. error performing HTTP request: %s", t.Name(), err)
			return
		}
		defer httpResp.Body.Close()

		body, err = io.ReadAll(httpResp.Body)
		if err != nil {
			t.Errorf("%s. Error reading response body: %s", t.Name(), err)
			return
		}

		fmt.Println(">>>>", httpResp.StatusCode, string(body))

		assert.NoError(t, err)
		assert.Equal(t, input.expectedStatusCode, httpResp.StatusCode)
	})

	if !passed {
		return fmt.Errorf("test failed")
	}

	return nil
}
