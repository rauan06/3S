package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"triples/http_utils"
)

type nullWriter struct{}

func (nw *nullWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func TestGET(t *testing.T) {
	log.SetOutput(&nullWriter{})

	tests := []struct {
		name         string
		requestURL   string
		expectedBody string
		expectedCode int
	}{
		{
			name:       "Get request without session 1",
			requestURL: "/",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Message>You dont have enough rights to access this bucket or object</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusForbidden,
		},
		{
			name:       "Get request without session 2",
			requestURL: "/////////////////////////////////",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Message>You dont have enough rights to access this bucket or object</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusForbidden,
		},
		{
			name:       "Get request without session 3",
			requestURL: "/:/:/:/:/",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Message>You dont have enough rights to access this bucket or object</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusForbidden,
		},
		{
			name:       "Get request without session 4",
			requestURL: "/123",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Message>You dont have enough rights to access this bucket or object</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.requestURL, nil)
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

func TestGET_withSession(t *testing.T) {
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
			name:         "Get request 1",
			requestURL:   "/",
			expectedBody: "<ListAllMyAllBucketsResult>",
			expectedCode: http.StatusOK,
		},
		{
			name:       "Get request 2",
			requestURL: "/////////////////////////////////",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Message>You dont have enough rights to access this bucket or object</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusOK,
		},
		{
			name:       "Get request 3",
			requestURL: "/:/:/:/:/",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>400</Code>\n" +
				"  <Message>400 Bad request</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:       "Get request 4",
			requestURL: "/123?session_id=rauan",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>403</Code>\n" +
				"  <Message>You dont have enough rights to access this bucket or object</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardownSuite := SetupWithSession(t)
			defer teardownSuite(t)

			req, err := http.NewRequest("GET", tt.requestURL, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(http_utils.Handler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			if rr.Body.String() == tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
