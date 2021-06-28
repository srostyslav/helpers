package helpers

import "testing"

func TestHelpers(t *testing.T) {
	InitLogger()

	req := &Request{Url: "https://www.google.com"}
	if err := req.Get(); err != nil {
		t.Error(err)
	} else {
		InfoLogger.Println(req.ResponseCode)
	}
}
