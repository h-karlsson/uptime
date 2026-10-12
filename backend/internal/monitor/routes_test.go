package monitor

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

const validBody = `{"name":"Example","url":"https://example.com","method":"GET","intervalSeconds":60,"timeoutMs":10000,"expectedStatus":200,"enabled":true}`

func TestMonitorLifecycle(t *testing.T) {
	router := newTestRouter(&fakeStore{}, 1)

	create := send(router, http.MethodPost, "/api/monitors", validBody)
	if create.Code != http.StatusCreated {
		t.Fatalf("create: status = %d, want 201; body = %s", create.Code, create.Body.String())
	}
	var created Monitor
	if err := json.Unmarshal(create.Body.Bytes(), &created); err != nil || created.ID == 0 {
		t.Fatalf("create: body = %s, want the created monitor", create.Body.String())
	}

	list := send(router, http.MethodGet, "/api/monitors", "")
	var listed []Monitor
	if err := json.Unmarshal(list.Body.Bytes(), &listed); err != nil || len(listed) != 1 {
		t.Fatalf("list: body = %s, want one monitor", list.Body.String())
	}

	update := send(router, http.MethodPut, "/api/monitors/1", `{"name":"Renamed","url":"https://example.com","method":"HEAD","intervalSeconds":300,"timeoutMs":5000,"expectedStatus":204,"enabled":false}`)
	var updated Monitor
	if err := json.Unmarshal(update.Body.Bytes(), &updated); err != nil || update.Code != 200 || updated.Name != "Renamed" || updated.Enabled {
		t.Fatalf("update: status = %d, body = %s", update.Code, update.Body.String())
	}

	if response := send(router, http.MethodDelete, "/api/monitors/1", ""); response.Code != http.StatusNoContent {
		t.Fatalf("delete: status = %d, want 204", response.Code)
	}
	if response := send(router, http.MethodGet, "/api/monitors/1", ""); response.Code != http.StatusNotFound {
		t.Fatalf("get after delete: status = %d, want 404", response.Code)
	}
}

func TestEmptyListIsJSONArray(t *testing.T) {
	response := send(newTestRouter(&fakeStore{}, 1), http.MethodGet, "/api/monitors", "")
	if body := response.Body.String(); body != "[]\n" {
		t.Fatalf("body = %q, want an empty JSON array", body)
	}
}

func TestOtherUsersMonitorsAreNotFound(t *testing.T) {
	store := &fakeStore{}
	if _, err := store.Create(context.Background(), 2, validInput()); err != nil {
		t.Fatal(err)
	}
	router := newTestRouter(store, 1)

	tests := []struct {
		method string
		body   string
	}{
		{http.MethodGet, ""},
		{http.MethodPut, validBody},
		{http.MethodDelete, ""},
	}
	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {
			if response := send(router, test.method, "/api/monitors/1", test.body); response.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want 404", response.Code)
			}
		})
	}
	if monitors, _ := store.List(context.Background(), 2); len(monitors) != 1 {
		t.Fatal("another user's monitor was changed or deleted")
	}
}

func TestRequestErrors(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		body   string
		status int
	}{
		{"malformed json", http.MethodPost, "/api/monitors", `{`, 400},
		{"unknown field", http.MethodPost, "/api/monitors", `{"name":"x","extra":true}`, 400},
		{"invalid input", http.MethodPost, "/api/monitors", `{"name":"x","url":"nope","method":"GET","intervalSeconds":60,"timeoutMs":10000,"expectedStatus":200}`, 400},
		{"non-numeric id", http.MethodGet, "/api/monitors/abc", "", 404},
		{"unknown id", http.MethodDelete, "/api/monitors/99", "", 404},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := send(newTestRouter(&fakeStore{}, 1), test.method, test.path, test.body)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, test.status, response.Body.String())
			}
		})
	}
}
