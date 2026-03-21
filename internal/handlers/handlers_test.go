package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStorageHandler_Handler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		wantStatus     int
		wantResponse   string // часть ответа, которую проверяем
		wantErr        bool
		wantErrMessage string // часть сообщения об ошибке
	}{
		// ========== POST - CREATE DEPARTMENT (path = "") ==========
		// {
		// 	name:       "POST /departments/ - valid request",
		// 	method:     http.MethodPost,
		// 	url:        "/departments/",
		// 	wantStatus: http.StatusOK,
		// 	wantErr:    false,
		// },
		// {
		// 	name:       "POST /departments - without trailing slash",
		// 	method:     http.MethodPost,
		// 	url:        "/departments",
		// 	wantStatus: http.StatusOK,
		// 	wantErr:    false,
		// },

		// ========== POST - CREATE EMPLOYEE (path = {id}/employees) ==========
		{
			name:       "POST /departments/1/employees - valid request",
			method:     http.MethodPost,
			url:        "/departments/1/employees",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "POST /departments/123/employees - valid with large id",
			method:     http.MethodPost,
			url:        "/departments/123/employees",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "POST /departments/1/employees/ - with trailing slash",
			method:     http.MethodPost,
			url:        "/departments/1/employees/",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "POST /departments/abc/employees - invalid department id",
			method:     http.MethodPost,
			url:        "/departments/abc/employees",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "POST /departments/-1/employees - negative id",
			method:     http.MethodPost,
			url:        "/departments/-1/employees",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "POST /departments/1/employees/extra - too many segments",
			method:     http.MethodPost,
			url:        "/departments/1/employees/extra",
			wantStatus: http.StatusOK, // в твоем коде не обрабатывается, падает в default?
			wantErr:    false,         // нужно проверить
		},
		{
			name:       "POST /departments/1/wrong - wrong second segment",
			method:     http.MethodPost,
			url:        "/departments/1/wrong",
			wantStatus: http.StatusBadRequest, // сейчас вернет 200, хотя должен 404
			wantErr:    true,
		},

		// ========== GET - DEPARTMENT (path = {id}) ==========
		{
			name:       "GET /departments/1 - valid request",
			method:     http.MethodGet,
			url:        "/departments/1",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "GET /departments/123 - valid with large id",
			method:     http.MethodGet,
			url:        "/departments/123",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "GET /departments/1/ - with trailing slash",
			method:     http.MethodGet,
			url:        "/departments/1/",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "GET /departments/abc - invalid id (string)",
			method:     http.MethodGet,
			url:        "/departments/abc",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "GET /departments/-1 - negative id",
			method:     http.MethodGet,
			url:        "/departments/-1",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "GET /departments/ - empty id",
			method:     http.MethodGet,
			url:        "/departments/",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "GET /departments - no id",
			method:     http.MethodGet,
			url:        "/departments",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},

		// ========== PATCH - UPDATE DEPARTMENT ==========
		{
			name:       "PATCH /departments/1 - valid request",
			method:     http.MethodPatch,
			url:        "/departments/1",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "PATCH /departments/123 - valid with large id",
			method:     http.MethodPatch,
			url:        "/departments/123",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "PATCH /departments/abc - invalid id",
			method:     http.MethodPatch,
			url:        "/departments/abc",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "PATCH /departments/ - empty id",
			method:     http.MethodPatch,
			url:        "/departments/",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},

		// ========== DELETE - DEPARTMENT ==========
		{
			name:       "DELETE /departments/1 - valid request",
			method:     http.MethodDelete,
			url:        "/departments/1",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "DELETE /departments/123 - valid with large id",
			method:     http.MethodDelete,
			url:        "/departments/123",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "DELETE /departments/abc - invalid id",
			method:     http.MethodDelete,
			url:        "/departments/abc",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "DELETE /departments/ - empty id",
			method:     http.MethodDelete,
			url:        "/departments/",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},

		// ========== METHOD NOT ALLOWED ==========
		{
			name:       "PUT not allowed",
			method:     http.MethodPut,
			url:        "/departments/1",
			wantStatus: http.StatusMethodNotAllowed,
			wantErr:    true,
		},
		{
			name:       "OPTIONS not allowed",
			method:     http.MethodOptions,
			url:        "/departments/1",
			wantStatus: http.StatusMethodNotAllowed,
			wantErr:    true,
		},
		{
			name:       "HEAD not allowed",
			method:     http.MethodHead,
			url:        "/departments/1",
			wantStatus: http.StatusMethodNotAllowed,
			wantErr:    true,
		},

		// ========== EDGE CASES ==========
		{
			name:       "POST with invalid path - should return 404",
			method:     http.MethodPost,
			url:        "/departments/extra/segments",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "POST with very long id",
			method:     http.MethodPost,
			url:        "/departments/999999999999/employees",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "GET with very long id",
			method:     http.MethodGet,
			url:        "/departments/999999999999",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем запрос
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			// Создаем хендлер (Storage = nil, т.к. тестируем только парсинг URL)
			handler := &StorageHandler{
				Storage: nil, // для тестов URL он не нужен
			}

			// Вызываем хендлер
			got, err := handler.Handler(w, req)

			// Проверяем статус код
			if w.Code != tt.wantStatus {
				t.Errorf("Status code: expected %d, got %d", tt.wantStatus, w.Code)
				t.Logf("Response body: %s", w.Body.String())
			}

			// Проверяем наличие ошибки
			if (err != nil) != tt.wantErr {
				t.Errorf("Error: expected error = %v, got error = %v", tt.wantErr, err)
			}

			// Логируем результат для отладки
			t.Logf("Method: %s, URL: %s", tt.method, tt.url)
			t.Logf("Got: %s", got)
			t.Logf("Response: %s", w.Body.String())
			if err != nil {
				t.Logf("Error: %v", err)
			}
		})
	}
}

// Отдельный тест для проверки конкретных кейсов с разбором URL
func TestStorageHandler_URLParsing(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		expectedPath   string
		expectedID     string
		expectedAction string
	}{
		{
			name:         "POST create department root",
			method:       http.MethodPost,
			url:          "/departments/",
			expectedPath: "",
		},
		{
			name:         "POST create employee",
			method:       http.MethodPost,
			url:          "/departments/5/employees",
			expectedPath: "5/employees",
			expectedID:   "5",
		},
		{
			name:         "GET department by id",
			method:       http.MethodGet,
			url:          "/departments/42",
			expectedPath: "42",
			expectedID:   "42",
		},
		{
			name:         "PATCH department",
			method:       http.MethodPatch,
			url:          "/departments/10",
			expectedPath: "10",
			expectedID:   "10",
		},
		{
			name:         "DELETE department",
			method:       http.MethodDelete,
			url:          "/departments/7",
			expectedPath: "7",
			expectedID:   "7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)

			// Эмулируем логику парсинга из хендлера
			path := strings.TrimPrefix(req.URL.Path, "/departments/")
			path = strings.TrimSuffix(path, "/")

			if path != tt.expectedPath {
				t.Errorf("Path: expected %q, got %q", tt.expectedPath, path)
			}

			// Для POST с employees проверяем разбор
			if tt.method == http.MethodPost && tt.expectedID != "" {
				sl := strings.Split(path, "/")
				if len(sl) == 3 {
					id := sl[0]
					if id != tt.expectedID {
						t.Errorf("ID: expected %q, got %q", tt.expectedID, id)
					}
					if sl[1] != "employees" {
						t.Errorf("Action: expected 'employees', got %q", sl[1])
					}
				}
			}

			t.Logf("URL: %s -> Path: %q", tt.url, path)
		})
	}
}
