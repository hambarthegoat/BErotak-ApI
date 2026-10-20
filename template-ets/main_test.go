package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func kirim(t *testing.T, h http.Handler, metode, path, body string) (int, map[string]any, string) {
	t.Helper()
	req := httptest.NewRequest(metode, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	var obj map[string]any
	json.Unmarshal(rec.Body.Bytes(), &obj)
	return rec.Code, obj, rec.Body.String()
}

func harapKode(t *testing.T, got, want int, ket string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: kode %d, harusnya %d", ket, got, want)
	}
}

const bodyPinjam = `{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":3,"slot_selesai":4,"jumlah":1}`

func TestSumberDaya(t *testing.T) {
	h := NewHandler()
	kode, _, raw := kirim(t, h, "GET", "/sumber-daya", "")
	harapKode(t, kode, 200, "daftar kosong")
	if strings.TrimSpace(raw) != "[]" {
		t.Fatalf("daftar kosong harus [], dapat %s", raw)
	}

	kode, obj, _ := kirim(t, h, "POST", "/sumber-daya", `{"nama":"  Aula A ","jenis":"ruang","stok":1}`)
	harapKode(t, kode, 201, "buat")
	if obj["id"] != float64(1) || obj["nama"] != "Aula A" {
		t.Fatalf("objek salah: %v", obj)
	}

	for _, c := range []string{`{"nama":"x","jenis":"alat"}`, `bukan json`} {
		kode, obj, _ := kirim(t, h, "POST", "/sumber-daya", c)
		harapKode(t, kode, 400, c)
		if obj["error"] == nil {
			t.Fatalf("error harus ada untuk %s", c)
		}
	}

	kode, _, _ = kirim(t, h, "GET", "/sumber-daya/abc", "")
	harapKode(t, kode, 400, "id bukan angka")
	kode, _, _ = kirim(t, h, "GET", "/sumber-daya/99", "")
	harapKode(t, kode, 404, "id tidak ada")
	kode, _, _ = kirim(t, h, "GET", "/sumber-daya/1", "")
	harapKode(t, kode, 200, "detail")
}

func TestBuatPeminjaman(t *testing.T) {
	h := NewHandler()
	kode, obj, _ := kirim(t, h, "POST", "/peminjaman", bodyPinjam)
	harapKode(t, kode, 201, "buat")
	if obj["status"] != "aktif" || obj["id"] != float64(1) {
		t.Fatalf("objek salah: %v", obj)
	}
	kode, _, _ = kirim(t, h, "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"Rina"}`)
	harapKode(t, kode, 400, "field kurang")
}

func TestLihatPeminjaman(t *testing.T) {
	h := NewHandler()
	_, _, raw := kirim(t, h, "GET", "/peminjaman", "")
	if strings.TrimSpace(raw) != "[]" {
		t.Fatalf("daftar kosong harus [], dapat %s", raw)
	}
	kirim(t, h, "POST", "/peminjaman", bodyPinjam)
	kirim(t, h, "POST", "/peminjaman", strings.Replace(bodyPinjam, "2026-11-02", "2026-11-03", 1))

	_, _, raw = kirim(t, h, "GET", "/peminjaman?tanggal=2026-11-03", "")
	var daftar []Peminjaman
	json.Unmarshal([]byte(raw), &daftar)
	if len(daftar) != 1 || daftar[0].ID != 2 {
		t.Fatalf("filter tanggal salah: %s", raw)
	}
	_, _, raw = kirim(t, h, "GET", "/peminjaman?status=selesai", "")
	if strings.TrimSpace(raw) != "[]" {
		t.Fatalf("filter kosong harus [], dapat %s", raw)
	}
	kode, _, _ := kirim(t, h, "GET", "/peminjaman?sumber_daya_id=abc", "")
	harapKode(t, kode, 400, "filter bukan angka")

	kode, _, _ = kirim(t, h, "GET", "/peminjaman/abc", "")
	harapKode(t, kode, 400, "id bukan angka")
	kode, _, _ = kirim(t, h, "GET", "/peminjaman/99", "")
	harapKode(t, kode, 404, "tidak ada")
	kode, _, _ = kirim(t, h, "GET", "/peminjaman/2", "")
	harapKode(t, kode, 200, "detail")
}
