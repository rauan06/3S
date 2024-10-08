package test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
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
				"  <Message>Incorrect bucket name</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:       "PUT request without session 2",
			requestURL: "/asdasdas/asdasdsad/asd?session_id=123",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>404</Code>\n" +
				"  <Message>404 Not Found</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusNotFound,
		},
		{
			name:       "PUT request without session 3",
			requestURL: "/:/:/:/:/",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>400</Code>\n" +
				"  <Message>Incorrect bucket name</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:       "PUT request without session 4",
			requestURL: "/123?session_id=123",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>200</Code>\n" +
				"  <Message>Bucket session id: [a-z0-9]+</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusOK,
		},
		{
			name:       "PUT request without session 5",
			requestURL: "/123",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>409</Code>\n" +
				"  <Message>Bucket name already exists</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusConflict,
		},
		{
			name:       "PUT request without session 6",
			requestURL: "/12345?session_id=1234567",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>200</Code>\n" +
				"  <Message>Bucket session id: [a-z0-9]+</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusOK,
		},
		{
			name:       "PUT request without session 7",
			requestURL: "/192.168.0.1",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>400</Code>\n" +
				"  <Message>Incorrect bucket name</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:       "PUT request without session 8",
			requestURL: "/192..rauan",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>400</Code>\n" +
				"  <Message>Incorrect bucket name</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:       "PUT request without session 9",
			requestURL: "/707560180ca618f34faacf97f346bae6a59eef815e47177dfa6c8d2233696c40a",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>400</Code>\n" +
				"  <Message>Incorrect bucket name</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:       "PUT request without session 10",
			requestURL: "/77",
			expectedBody: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				"<Response>\n" +
				"  <Code>400</Code>\n" +
				"  <Message>Incorrect bucket name</Message>\n" +
				"</Response>\n",
			expectedCode: http.StatusBadRequest,
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

			if tt.expectedCode == http.StatusOK {
				re := regexp.MustCompile(`^<\?xml version="1\.0" encoding="UTF-8"\?>\n<Response>\n  <Code>200</Code>\n  <Message>Bucket session id: [a-zA-Z0-9]+</Message>\n</Response>\n$`)
				if !re.MatchString(rr.Body.String()) {
					t.Errorf("handler returned unexpected body: got\n %v\n want to match regex\n %v\n", rr.Body.String(), tt.expectedBody)
				}
			} else {
				if rr.Body.String() != tt.expectedBody {
					t.Errorf("handler returned unexpected body: got\n %v\n want\n %v\n", rr.Body.String(), tt.expectedBody)
				}
			}
		})
	}
}
