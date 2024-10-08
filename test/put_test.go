package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"triples/http_utils"
)

func TestPUT(t *testing.T) {
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
				"  <Code>400</Code>\n" +
				"  <Message>400 Bad Request</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:       "PUT request without session 2",
			requestURL: "/asdasdas/asdasdsad/asd?session_id=123",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>400</Code>\n" +
				"  <Message>400 Bad Request</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:       "PUT request without session 3",
			requestURL: "/:/:/:/:/",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>400</Code>\n" +
				"  <Message>400 Bad Request</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:       "PUT request without session 4",
			requestURL: "/123?session_id=123",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>200</Code>\n" +
				"  <Message>Bucket session id: 123</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusOK,
		},
		{
			name:       "PUT request without session 5",
			requestURL: "/123",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>409</Code>\n" +
				"  <Message>409 Conflict</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusConflict,
		},
		{
			name:       "PUT request without session 6",
			requestURL: "/12345?session_id=1234567",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>200</Code>\n" +
				"  <Message>Bucket session id: 1234567</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
				t.Errorf("handler returned unexpected body: got\n %v\n want\n %v\n", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
