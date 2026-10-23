package user

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"DevOpsMiniProject/entity"
	"DevOpsMiniProject/repository"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type fakeRepo struct {
	createErr error
	deleteErr error
	getAllErr error
	users     []entity.User
}

func (f fakeRepo) CreateUser(u entity.User) (entity.User, error) {
	if f.createErr != nil {
		return entity.User{}, f.createErr
	}
	u.ID = 1
	return u, nil
}
func (f fakeRepo) DeleteUser(string) error            { return f.deleteErr }
func (f fakeRepo) GetAllUser() ([]entity.User, error) { return f.users, f.getAllErr }
func (f fakeRepo) PingDB(context.Context) error       { return nil }

func doRequest(t *testing.T, repo fakeRepo, method, body string) (int, string) {
	t.Helper() // ถ้า fail จะชี้บรรทัดใน test ไม่ใช่บรรทัดใน helper

	app := fiber.New()
	svc := ProvideUserService(repo)
	app.Post("/user", svc.HandleCreateUser)
	app.Delete("/user", svc.HandleDeleteUser)
	app.Get("/user", svc.HandleGetAllUser)

	req := httptest.NewRequest(method, "/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b)
}

func TestHandleCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		repo       fakeRepo
		wantStatus int
		wantBody   string // ข้อความที่ต้องมีใน body ("" = ไม่เช็ก)
	}{
		{"สร้างสำเร็จ", `{"username":"thanut","name":"Thanut","age":30}`, fakeRepo{}, fiber.StatusCreated, `"id":1`},
		{"JSON พัง", `{"username":`, fakeRepo{}, fiber.StatusBadRequest, ""},
		{"ไม่มี name", `{"username":"thanut"}`, fakeRepo{}, fiber.StatusBadRequest, "Name"},
		{"username สั้นเกิน", `{"username":"ab","name":"T"}`, fakeRepo{}, fiber.StatusBadRequest, "Username"},
		{"age เป็น string", `{"username":"thanut","name":"T","age":"five"}`, fakeRepo{}, fiber.StatusBadRequest, ""},
		{"username ซ้ำ", `{"username":"thanut","name":"T"}`, fakeRepo{createErr: gorm.ErrDuplicatedKey}, fiber.StatusConflict, "already exists"},
		{"DB พัง", `{"username":"thanut","name":"T"}`, fakeRepo{createErr: errors.New("db down")}, fiber.StatusInternalServerError, "internal error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body := doRequest(t, tt.repo, "POST", tt.body)

			if status != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", status, tt.wantStatus, body)
			}
			if tt.wantBody != "" && !strings.Contains(body, tt.wantBody) {
				t.Errorf("body = %q, want contains %q", body, tt.wantBody)
			}
		})
	}
}

func TestHandleDeleteUser(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		repo       fakeRepo
		wantStatus int
		wantBody   string
	}{
		{"ลบสำเร็จ", `{"username":"thanut"}`, fakeRepo{}, fiber.StatusOK, "thanut"},
		{"JSON พัง", `{"username":`, fakeRepo{}, fiber.StatusBadRequest, ""},
		{"username สั้นเกิน", `{"username":"ab"}`, fakeRepo{}, fiber.StatusBadRequest, "Username"},
		{"ไม่มี user นี้", `{"username":"nobody"}`, fakeRepo{deleteErr: repository.ErrUserNotFound}, fiber.StatusNotFound, "not found"},
		{"DB พัง", `{"username":"thanut"}`, fakeRepo{deleteErr: errors.New("db down")}, fiber.StatusInternalServerError, "internal error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body := doRequest(t, tt.repo, "DELETE", tt.body)

			if status != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", status, tt.wantStatus, body)
			}
			if tt.wantBody != "" && !strings.Contains(body, tt.wantBody) {
				t.Errorf("body = %q, want contains %q", body, tt.wantBody)
			}
		})
	}
}

func TestHandleGetAllUser(t *testing.T) {
	tests := []struct {
		name       string
		repo       fakeRepo
		wantStatus int
		wantBody   string
	}{
		{"มี user", fakeRepo{users: []entity.User{{Username: "thanut", Name: "Thanut", Age: 30}}}, fiber.StatusOK, "thanut"},
		{"ไม่มี user", fakeRepo{}, fiber.StatusOK, "No user found"},
		{"DB พัง", fakeRepo{getAllErr: errors.New("db down")}, fiber.StatusInternalServerError, "internal error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body := doRequest(t, tt.repo, "GET", "")

			if status != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", status, tt.wantStatus, body)
			}
			if !strings.Contains(body, tt.wantBody) {
				t.Errorf("body = %q, want contains %q", body, tt.wantBody)
			}
		})
	}
}
