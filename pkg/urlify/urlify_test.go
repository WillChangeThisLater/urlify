package urlify

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestUrlify(t *testing.T) {

	bucket := "urlify"
	prefix := "urlify"
	region := "us-east-2"
	buffer := []byte("i contain some random data")
	url, err := Urlify(bucket, prefix, region, buffer)
	if err != nil {
		t.Errorf("urlify(%s, %s, %s, %v) failed: %v\n", bucket, prefix, region, buffer, err)
	}

	response, err := http.Get(url)
	if err != nil {
		t.Errorf("Invalid response from presigned URL %s: %v\n", url, err)
	}

	if response.StatusCode != 200 {
		t.Errorf("Unexpected status code %d (expected %d)", response.StatusCode, 200)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Could not read response body: %v", err)
	}
	if !bytes.Equal(buffer, body) {
		t.Errorf("Buffer value provided to urlify != response from presigned URL")
	}
}
