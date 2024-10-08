package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"triples/http_utils"
)

func TestDELETE(t *testing.T) {
	log.SetOutput(&nullWriter{})
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name         string
		requestURL   string
		expectedBody string
		expectedCode int
	}{
		{
			name:       "DELETE request without session 1",
			requestURL: "/",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Message>403 Forbidden</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusForbidden,
		},
		{
			name:       "DELETE request without session 2",
			requestURL: "/////////////////////////////////",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Message>403 Forbidden</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusForbidden,
		},
		{
			name:       "DELETE request without session 3",
			requestURL: "/:/:/:/:/",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Message>403 Forbidden</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusForbidden,
		},
		{
			name:       "DELETE request without session 4",
			requestURL: "/123",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Message>403 Forbidden</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardownSuite := SetupSuite(t)
			defer teardownSuite(t)

			req, err := http.NewRequest("DELETE", tt.requestURL, nil)
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

func TestDELETE_withSession(t *testing.T) {
	log.SetOutput(&nullWriter{})
	teardownSuite := SetupWithSession(t)
	defer teardownSuite(t)

	tests := []struct {
		name         string
		requestURL   string
		expectedBody string
		expectedCode int
	}{
		{
			name:       "DELETE request without session 1",
			requestURL: "/",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>400</Code>\n" +
				"  <Message>400 Bad Request</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardownSuite := SetupWithSession(t)
			defer teardownSuite(t)

			req, err := http.NewRequest("DELETE", tt.requestURL, nil)
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
