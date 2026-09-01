package util

import (
	"strings"
	"testing"
)

// TestBuildtime ทดสอบว่า Buildtime() แปลงค่า buildTime ที่ฝังตอน build
// ให้เป็นเวลาไทยที่อ่านง่าย และรับมือกับค่าที่ผิดรูปแบบได้
func TestBuildtime(t *testing.T) {
	// เก็บค่าเดิมไว้ แล้วคืนหลัง test จบ กัน test อื่นได้รับผลกระทบ
	original := buildTime
	defer func() { buildTime = original }()

	tests := []struct {
		name  string
		given string
		want  string
	}{
		{
			name:  "แปลงเวลา RFC3339 เป็นเวลาไทย (UTC+7)",
			given: "2026-06-29T10:00:00Z",
			want:  "2026/06/29 | 17:00:00",
		},
		{
			name:  "ค่าผิดรูปแบบต้องไม่ทำให้ panic",
			given: "not-a-timestamp",
			want:  "Invalid Time!",
		},
		{
			name:  "ค่าว่าง (กรณี build โดยไม่ส่ง ldflags)",
			given: "",
			want:  "Invalid Time!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildTime = tt.given

			got := Buildtime()

			if got != tt.want {
				t.Errorf("Buildtime() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestUptime ทดสอบว่า Uptime() คืน format ที่ถูกต้อง
// ไม่เช็คค่าเป๊ะ เพราะเวลาเดินตลอด
func TestUptime(t *testing.T) {
	got := Uptime()

	for _, unit := range []string{"h", "m", "s"} {
		if !strings.Contains(got, unit) {
			t.Errorf("Uptime() = %q, ต้องมีหน่วย %q อยู่ในผลลัพธ์", got, unit)
		}
	}
}

// TestIncrementRequest ทดสอบว่าตัวนับ request เพิ่มขึ้นจริง
// เทียบเป็นส่วนต่าง ไม่เทียบค่าสัมบูรณ์ เพราะ test อื่นอาจนับไปแล้ว
func TestIncrementRequest(t *testing.T) {
	before := GetTotalRequests()

	const times = 3
	for i := 0; i < times; i++ {
		IncrementRequest()
	}

	after := GetTotalRequests()

	if diff := after - before; diff != times {
		t.Errorf("นับได้ %d ครั้ง, ต้องการ %d ครั้ง", diff, times)
	}
}
