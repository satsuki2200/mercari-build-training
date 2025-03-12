package app

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseAddItemRequest(t *testing.T) {
	t.Parallel()

	type wants struct {
		req *AddItemRequest
		err bool
	}

	// STEP 6-1: define test cases
	cases := map[string]struct {
		// args map[string]string
		name string
		category string
		image []byte
		wants
	}{
		"ok: valid request": {
			name: "Alice",
			category: "people",
			image: []byte("dummy image data"),
			wants: wants{
				req: &AddItemRequest{
					Name: "Alice", // fill here
					Category: "people", // fill here
					Image: []byte("dummy image data"),
				},
				err: false,
			},
		},
		"ng: empty request": {
			name: "",
			category: "",
			image: []byte(""),
			wants: wants{
				req: nil,
				err: true,
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var body bytes.Buffer
			writer := multipart.NewWriter(&body)

			_ = writer.WriteField("name", tt.name)
			_ = writer.WriteField("category", tt.category)
			if len(tt.image) > 0 {
				part, err := writer.CreateFormFile("image", "dummy.jpg")
				if err != nil {
					t.Fatalf("failed to create form file: %v", err)
				}
				_, err = part.Write(tt.image)
				if err != nil {
					t.Fatalf("failed to write form file: %v", err)
				}
			}
			writer.Close()

			// prepare HTTP request
			req, err := http.NewRequest("POST", "http://localhost:9000/items", &body)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())

			// execute test target
			got, err := parseAddItemRequest(req)

			// confirm the result
			if err != nil {
				if !tt.err {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}
			if diff := cmp.Diff(tt.wants.req, got); diff != "" {
				t.Errorf("unexpected request (-want +got):\n%s", diff)
			}		
		})
	}
}

// func TestHelloHandler(t *testing.T) {
// 	t.Parallel()

// 	// Please comment out for STEP 6-2
// 	// predefine what we want
// 	type wants struct {
// 		code int               // desired HTTP status code
// 		body map[string]string // desired body
// 	}
// 	want := wants{
// 		code: http.StatusOK,
// 		body: map[string]string{"message": "Hello, world!"},
// 	}

// 	// set up test
// 	req := httptest.NewRequest("GET", "/hello", nil)
// 	res := httptest.NewRecorder()

// 	h := &Handlers{}
// 	h.Hello(res, req)

// 	// STEP 6-2: confirm the status code
// 	if diff := cmp.Diff(want.code, res.Code); diff != "" {
// 		t.Errorf("unexpected request (-want +got):\n%s", diff)
// 	}

// 	// STEP 6-2: confirm response body
// 	if diff := cmp.Diff(want.body, res.Body.String()); diff != "" {
// 		t.Errorf("unexpected request (-want +got):\n%s", diff)
// 	}
// }

// func TestAddItem(t *testing.T) {
// 	t.Parallel()

// 	type wants struct {
// 		code int
// 	}
// 	cases := map[string]struct {
// 		args     map[string]string
// 		injector func(m *MockItemRepository)
// 		wants
// 	}{
// 		"ok: correctly inserted": {
// 			args: map[string]string{
// 				"name":     "used iPhone 16e",
// 				"category": "phone",
// 			},
// 			injector: func(m *MockItemRepository) {
// 				// STEP 6-3: define mock expectation
// 				// succeeded to insert
// 			},
// 			wants: wants{
// 				code: http.StatusOK,
// 			},
// 		},
// 		"ng: failed to insert": {
// 			args: map[string]string{
// 				"name":     "used iPhone 16e",
// 				"category": "phone",
// 			},
// 			injector: func(m *MockItemRepository) {
// 				// STEP 6-3: define mock expectation
// 				// failed to insert
// 			},
// 			wants: wants{
// 				code: http.StatusInternalServerError,
// 			},
// 		},
// 	}

// 	for name, tt := range cases {
// 		t.Run(name, func(t *testing.T) {
// 			t.Parallel()

// 			ctrl := gomock.NewController(t)

// 			mockIR := NewMockItemRepository(ctrl)
// 			tt.injector(mockIR)
// 			h := &Handlers{itemRepo: mockIR}

// 			values := url.Values{}
// 			for k, v := range tt.args {
// 				values.Set(k, v)
// 			}
// 			req := httptest.NewRequest("POST", "/items", strings.NewReader(values.Encode()))
// 			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 			rr := httptest.NewRecorder()
// 			h.AddItem(rr, req)

// 			if tt.wants.code != rr.Code {
// 				t.Errorf("expected status code %d, got %d", tt.wants.code, rr.Code)
// 			}
// 			if tt.wants.code >= 400 {
// 				return
// 			}

// 			for _, v := range tt.args {
// 				if !strings.Contains(rr.Body.String(), v) {
// 					t.Errorf("response body does not contain %s, got: %s", v, rr.Body.String())
// 				}
// 			}
// 		})
// 	}
// }

// STEP 6-4: uncomment this test
// func TestAddItemE2e(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping e2e test")
// 	}

// 	db, closers, err := setupDB(t)
// 	if err != nil {
// 		t.Fatalf("failed to set up database: %v", err)
// 	}
// 	t.Cleanup(func() {
// 		for _, c := range closers {
// 			c()
// 		}
// 	})

// 	type wants struct {
// 		code int
// 	}
// 	cases := map[string]struct {
// 		args map[string]string
// 		wants
// 	}{
// 		"ok: correctly inserted": {
// 			args: map[string]string{
// 				"name":     "used iPhone 16e",
// 				"category": "phone",
// 			},
// 			wants: wants{
// 				code: http.StatusOK,
// 			},
// 		},
// 		"ng: failed to insert": {
// 			args: map[string]string{
// 				"name":     "",
// 				"category": "phone",
// 			},
// 			wants: wants{
// 				code: http.StatusBadRequest,
// 			},
// 		},
// 	}

// 	for name, tt := range cases {
// 		t.Run(name, func(t *testing.T) {
// 			h := &Handlers{itemRepo: &itemRepository{db: db}}

// 			values := url.Values{}
// 			for k, v := range tt.args {
// 				values.Set(k, v)
// 			}
// 			req := httptest.NewRequest("POST", "/items", strings.NewReader(values.Encode()))
// 			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 			rr := httptest.NewRecorder()
// 			h.AddItem(rr, req)

// 			// check response
// 			if tt.wants.code != rr.Code {
// 				t.Errorf("expected status code %d, got %d", tt.wants.code, rr.Code)
// 			}
// 			if tt.wants.code >= 400 {
// 				return
// 			}
// 			for _, v := range tt.args {
// 				if !strings.Contains(rr.Body.String(), v) {
// 					t.Errorf("response body does not contain %s, got: %s", v, rr.Body.String())
// 				}
// 			}

// 			// STEP 6-4: check inserted data
// 		})
// 	}
// }

// func setupDB(t *testing.T) (db *sql.DB, closers []func(), e error) {
// 	t.Helper()

// 	defer func() {
// 		if e != nil {
// 			for _, c := range closers {
// 				c()
// 			}
// 		}
// 	}()

// 	// create a temporary file for e2e testing
// 	f, err := os.CreateTemp(".", "*.sqlite3")
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	closers = append(closers, func() {
// 		f.Close()
// 		os.Remove(f.Name())
// 	})

// 	// set up tables
// 	db, err = sql.Open("sqlite3", f.Name())
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	closers = append(closers, func() {
// 		db.Close()
// 	})

// 	// TODO: replace it with real SQL statements.
// 	cmd := `CREATE TABLE IF NOT EXISTS items (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		name VARCHAR(255),
// 		category VARCHAR(255)
// 	)`
// 	_, err = db.Exec(cmd)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return db, closers, nil
// }
