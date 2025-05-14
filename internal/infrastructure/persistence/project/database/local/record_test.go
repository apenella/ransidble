package local

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateDataHash(t *testing.T) {
	tests := []struct {
		desc     string
		data     json.RawMessage
		expected string
		err      error
	}{
		{
			desc:     "Testing hash generation with valid data",
			data:     json.RawMessage(`{"key": "value"}`),
			expected: "e43abcf3375244839c012f9633f95862d232a95b00d5bc7348b3098b9fed7f32",
			err:      nil,
		},
		{
			desc:     "Testing an error when generating hash with invalid data",
			data:     json.RawMessage(`{"key": "value"`),
			expected: "",
			err:      fmt.Errorf("error generating hash for data: json: error calling MarshalJSON for type json.RawMessage: unexpected end of JSON input"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			hash, err := generateDataHash(test.data)
			if err != nil && test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err, "occured an unexpected error")
				assert.Nil(t, test.err, "an expected error not occurred")
				assert.Equal(t, test.expected, hash)
			}
		})
	}
}

func TestCreateRecord(t *testing.T) {
	tests := []struct {
		desc     string
		id       string
		data     json.RawMessage
		expected *Record
		err      error
	}{
		{
			desc: "Testing record creation with valid data",
			id:   "test-id",
			data: json.RawMessage(`{"key": "value"}`),
			expected: &Record{
				Hash: "e43abcf3375244839c012f9633f95862d232a95b00d5bc7348b3098b9fed7f32",
				Data: json.RawMessage(`{"key": "value"}`),
			},
			err: nil,
		},
		{
			desc:     "Testing error when creating record with invalid data",
			id:       "test-id",
			data:     json.RawMessage(`{"key": "value"`), // Invalid JSON
			expected: nil,
			err:      fmt.Errorf("error creating record: error generating hash for data: json: error calling MarshalJSON for type json.RawMessage: unexpected end of JSON input"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			record, err := CreateRecord(test.id, test.data)
			if err != nil && test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err, "occured an unexpected error")
				assert.Nil(t, test.err, "an expected error not occurred")
				assert.Equal(t, test.expected.Data, record.Data)
				assert.Equal(t, test.expected.Hash, record.Hash)
			}
		})
	}
}

func TestGenerateRecordReference(t *testing.T) {
	tests := []struct {
		desc     string
		id       string
		expected string
		err      error
	}{
		{
			desc:     "Testing generate record reference with valid ID and version",
			id:       "test-id",
			expected: "6cc41d5ec590ab78cccecf81ef167d418c309a4598e8e45fef78039f7d9aa9fe",
		},
		{
			desc:     "Testing error when generating record reference with empty ID",
			id:       "",
			expected: "",
			err:      fmt.Errorf("error generating record reference: invalid ID. id is empty"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			reference, err := generateRecordReference(test.id)

			if err != nil && test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err, "occured an unexpected error")
				assert.Nil(t, test.err, "an expected error not occurred")
				assert.Equal(t, test.expected, reference)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	tests := []struct {
		desc     string
		record   *Record
		expected bool
		err      error
	}{
		{
			desc: "Testing record verification with valid data",
			record: &Record{
				Hash: "9724c1e20e6e3e4d7f57ed25f9d4efb006e508590d528c90da597f6a775c13e5",
				Data: json.RawMessage(`{"key": "value"}`),
			},
			expected: true,
			err:      nil,
		},

		{
			desc: "Testing error verifing record when data is nil",
			record: &Record{
				Hash: "",
				Data: nil,
			},
			expected: false,
			err:      fmt.Errorf(ErrRecordDataIsEmpty),
		},
		{
			desc: "Testing error verifing record when hash is empty",
			record: &Record{
				Hash: "",
				Data: json.RawMessage(`{"key": "value"}`),
			},
			expected: false,
			err:      fmt.Errorf(ErrRcordHashIsEmpty),
		},

		{
			desc: "Testing error verifing record when hash is not the expected",
			record: &Record{
				Hash: "invalid-hash",
				Data: json.RawMessage(`{"key": "value"}`),
			},
			expected: false,
			err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			result, err := test.record.Verify()

			if err != nil && test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err, "occured an unexpected error")
				assert.Nil(t, test.err, "an expected error not occurred")
				assert.Equal(t, test.expected, result)
			}
		})
	}
}
