package healthz

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"DevOpsMiniProject/entity"

	"github.com/gofiber/fiber/v3"
)

type fakeRepo struct {
	pingErr error
	hang    bool
}

func (f fakeRepo) PingDB(ctx context.Context) error {
	if f.hang {
		<-ctx.Done()
		return ctx.Err()
	}
	return f.pingErr
}
func (fakeRepo) CreateUser(u entity.User) (entity.User, error) { return u, nil }
func (fakeRepo) DeleteUser(string) error                       { return nil }
func (fakeRepo) GetAllUser() ([]entity.User, error)            { return nil, nil }

func doRequest(t *testing.T, repo fakeRepo, path string) (int, string, time.Duration) {
	t.Helper()

	app := fiber.New()
	svc := ProvideHealthzService(repo)
	app.Get("/livez", svc.HandleLiveCheck)
	app.Get("/readyz", svc.HandleHealthCheck)

	req := httptest.NewRequest("GET", path, nil)

	start := time.Now()
	resp, err := app.Test(req, fiber.TestConfig{Timeout: 3 * time.Second})
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b), elapsed
}

func TestHandleLiveCheck(t *testing.T) {
	status, body, _ := doRequest(t, fakeRepo{pingErr: errors.New("db down")}, "/livez")

	if status != fiber.StatusOK {
		t.Errorf("status = %d, want 200", status)
	}
	if !strings.Contains(body, "alive") {
		t.Errorf("body = %q, want contains %q", body, "alive")
	}
}

func TestHandleHealthCheck(t *testing.T) {
	tests := []struct {
		name       string
		repo       fakeRepo
		wantStatus int
		wantBody   string
		maxTime    time.Duration
	}{
		{"DB ปกติ", fakeRepo{}, fiber.StatusOK, `"database":"ok"`, 0},
		{"DB ต่อไม่ได้", fakeRepo{pingErr: errors.New("connection refused")}, fiber.StatusServiceUnavailable, `"database":"fail"`, 0},
		{"DB ค้าง ต้องตัดที่ ~800ms", fakeRepo{hang: true}, fiber.StatusServiceUnavailable, `"database":"fail"`, 950 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body, elapsed := doRequest(t, tt.repo, "/readyz")

			if status != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", status, tt.wantStatus, body)
			}
			if !strings.Contains(body, tt.wantBody) {
				t.Errorf("body = %q, want contains %q", body, tt.wantBody)
			}
			if tt.maxTime > 0 && elapsed > tt.maxTime {
				t.Errorf("took %v, want <= %v (timeout ไม่ทำงาน?)", elapsed, tt.maxTime)
			}
		})
	}
}
