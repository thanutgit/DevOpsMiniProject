package service_healthz

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

// func TestHandleHealthCheck(t *testing.T) {
// 	// สร้าง mock ของ HealthzService โดยไม่ต้องเชื่อมต่อฐานข้อมูลจริง
// 	mockService := &healthzService{
// 		db: nil, // ไม่ต้องใช้ฐานข้อมูลจริงใน test นี้
// 	}

// 	got := mockService.HandleHealthCheck(nil) // ส่ง nil เป็น context เพราะเราไม่ใช้จริงใน test นี้
// 	if !strings.Contains(got.Error(), "unhealthy") && !strings.Contains(got.Error(), "ok") {
// 		t.Errorf("HandleHealthCheck() = %q, ต้องมี 'unhealthy' หรือ 'ok' อยู่ในผลลัพธ์", got.Error())
// 	}
// }

func TestHandleLiveCheck(t *testing.T) {
	app := fiber.New()
	svc := ProvideHealthzService(nil) // ส่ง nil เป็น db เพราะเราไม่ใช้จริงใน test นี้
	app.Get("/livez", svc.HandleLiveCheck)

	req := httptest.NewRequest("GET", "/livez", nil)
	reps, err := app.Test(req)
	if err != nil {
		t.Fatalf("เกิดข้อผิดพลาดในการทดสอบ: %v", err)
	}
	defer reps.Body.Close()

	body, err := io.ReadAll(reps.Body)
	if err != nil {
		t.Fatalf("เกิดข้อผิดพลาดในการทดสอบ: %v", err)
	}

	if reps.StatusCode != 200 {
		t.Errorf("HandleLiveCheck() = %d, ต้องเป็น 200", reps.StatusCode)
	}

	if !strings.Contains(string(body), "alive") {
		t.Errorf("HandleLiveCheck() = %q, ต้องมี 'alive' อยู่ในผลลัพธ์", string(body))
	}
}
