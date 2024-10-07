package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"triples/http_utils"
)

func TestPUTWithoutSession(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name         string
		requestURL   string
		expectedBody string
		expectedCode int
	}{
		{
			name:       "PUT request without session 1",
			requestURL: "/",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Messege>403 Forbidden</Messege>\n" +
				"</Response>\n",
			expectedCode: http.StatusForbidden,
		},
		{
			name:       "PUT request without session 2",
			requestURL: "/////////////////////////////////",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Messege>403 Forbidden</Messege>\n" +
				"</Response>\n",
			expectedCode: http.StatusForbidden,
		},
		{
			name:       "PUT request without session 3",
			requestURL: "/:/:/:/:/",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Messege>403 Forbidden</Messege>\n" +
				"</Response>\n",
			expectedCode: http.StatusForbidden,
		},
		{
			name:       "PUT request without session 4",
			requestURL: "/123",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Messege>403 Forbidden</Messege>\n" +
				"</Response>\n",
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardownSuite := SetupSuite(t)
			defer teardownSuite(t)
			req, err := http.NewRequest("PUT", tt.requestURL, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(http_utils.Handler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
