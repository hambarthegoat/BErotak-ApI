package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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

func pinjamSlot(mulai, selesai string) string {
	return strings.Replace(bodyPinjam, `"slot_mulai":3,"slot_selesai":4`, `"slot_mulai":`+mulai+`,"slot_selesai":`+selesai, 1)
}

func handlerDenganRuang(t *testing.T, stok string) http.Handler {
	t.Helper()
	h := NewHandler()
	kode, _, _ := kirim(t, h, "POST", "/sumber-daya", `{"nama":"Lab A","jenis":"ruang","stok":`+stok+`}`)
	harapKode(t, kode, 201, "buat sumber daya")
	return h
}

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

	kasus := []string{
		`{"nama":"","jenis":"ruang","stok":1}`,
		`{"nama":"x","jenis":"gedung","stok":1}`,
		`{"nama":"x","jenis":"alat","stok":0}`,
		`{"nama":"x","jenis":"alat","stok":1.5}`,
		`{"nama":"x","jenis":"alat"}`,
		`{"nama":"` + strings.Repeat("é", 51) + `","jenis":"alat","stok":1}`,
		`bukan json`,
	}
	for _, c := range kasus {
		kode, obj, _ := kirim(t, h, "POST", "/sumber-daya", c)
		harapKode(t, kode, 400, c)
		if obj["error"] == nil {
			t.Fatalf("error harus ada untuk %s", c)
		}
	}
	kode, _, _ = kirim(t, h, "POST", "/sumber-daya", `{"nama":"`+strings.Repeat("é", 50)+`","jenis":"alat","stok":1}`)
	harapKode(t, kode, 201, "50 rune")

	kode, _, _ = kirim(t, h, "GET", "/sumber-daya/abc", "")
	harapKode(t, kode, 400, "id bukan angka")
	kode, _, _ = kirim(t, h, "GET", "/sumber-daya/99", "")
	harapKode(t, kode, 404, "id tidak ada")
	kode, _, _ = kirim(t, h, "GET", "/sumber-daya/1", "")
	harapKode(t, kode, 200, "detail")
}

func TestBuatPeminjaman(t *testing.T) {
	h := handlerDenganRuang(t, "1")
	kode, obj, _ := kirim(t, h, "POST", "/peminjaman", bodyPinjam)
	harapKode(t, kode, 201, "buat")
	if obj["status"] != "aktif" || obj["id"] != float64(1) {
		t.Fatalf("objek salah: %v", obj)
	}

	kasus := []string{
		`{"sumber_daya_id":9,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":1,"slot_selesai":1,"jumlah":1}`,
		`{"sumber_daya_id":1,"peminjam":" ","tanggal":"2026-11-02","slot_mulai":1,"slot_selesai":1,"jumlah":1}`,
		`{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-02-30","slot_mulai":1,"slot_selesai":1,"jumlah":1}`,
		`{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2027-02-29","slot_mulai":1,"slot_selesai":1,"jumlah":1}`,
		`{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":0,"slot_selesai":1,"jumlah":1}`,
		`{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":1,"slot_selesai":9,"jumlah":1}`,
		`{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":5,"slot_selesai":4,"jumlah":1}`,
		`{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":1,"slot_selesai":1,"jumlah":2}`,
		`{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":1,"slot_selesai":1}`,
		`{"sumber_daya_id":"1","peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":1,"slot_selesai":1,"jumlah":1}`,
	}
	for _, c := range kasus {
		kode, obj, _ := kirim(t, h, "POST", "/peminjaman", c)
		harapKode(t, kode, 400, c)
		if obj["error"] == nil {
			t.Fatalf("error harus ada untuk %s", c)
		}
	}
	kode, _, _ = kirim(t, h, "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2028-02-29","slot_mulai":1,"slot_selesai":1,"jumlah":1}`)
	harapKode(t, kode, 201, "29 Feb 2028")

	kode, obj, _ = kirim(t, h, "GET", "/tidak-ada", "")
	harapKode(t, kode, 404, "endpoint tidak ada")
	if obj["error"] == nil {
		t.Fatalf("404 harus berbentuk JSON error")
	}
}

func TestLihatPeminjaman(t *testing.T) {
	h := handlerDenganRuang(t, "1")
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

func TestBentrok(t *testing.T) {
	h := handlerDenganRuang(t, "1")
	kode, _, _ := kirim(t, h, "POST", "/peminjaman", bodyPinjam)
	harapKode(t, kode, 201, "pertama")

	kode, obj, _ := kirim(t, h, "POST", "/peminjaman", pinjamSlot("4", "5"))
	harapKode(t, kode, 409, "irisan")
	if obj["error"] == nil {
		t.Fatalf("409 harus punya field error")
	}
	kode, _, _ = kirim(t, h, "POST", "/peminjaman", pinjamSlot("1", "2"))
	harapKode(t, kode, 201, "bersebelahan")
	kode, _, _ = kirim(t, h, "POST", "/peminjaman", strings.Replace(bodyPinjam, "2026-11-02", "2026-11-03", 1))
	harapKode(t, kode, 201, "tanggal lain")

	// alat dengan stok 3: kapasitas dihitung per unit
	kirim(t, h, "POST", "/sumber-daya", `{"nama":"Proyektor","jenis":"alat","stok":3}`)
	alat := func(jumlah string) string {
		return strings.Replace(strings.Replace(bodyPinjam, `"sumber_daya_id":1`, `"sumber_daya_id":2`, 1), `"jumlah":1`, `"jumlah":`+jumlah, 1)
	}
	kode, _, _ = kirim(t, h, "POST", "/peminjaman", alat("2"))
	harapKode(t, kode, 201, "alat 2 dari 3")
	kode, _, _ = kirim(t, h, "POST", "/peminjaman", alat("2"))
	harapKode(t, kode, 409, "alat melebihi stok")
	kode, _, _ = kirim(t, h, "POST", "/peminjaman", alat("1"))
	harapKode(t, kode, 201, "alat sisa 1")
}

func TestPutDelete(t *testing.T) {
	h := handlerDenganRuang(t, "1")
	kirim(t, h, "POST", "/peminjaman", bodyPinjam)
	kirim(t, h, "POST", "/peminjaman", pinjamSlot("6", "7"))

	kode, obj, _ := kirim(t, h, "PUT", "/peminjaman/1", pinjamSlot("3", "5"))
	harapKode(t, kode, 200, "put tidak bentrok dengan diri sendiri")
	if obj["slot_selesai"] != float64(5) || obj["status"] != "aktif" || obj["id"] != float64(1) {
		t.Fatalf("hasil put salah: %v", obj)
	}
	kode, _, _ = kirim(t, h, "PUT", "/peminjaman/1", pinjamSlot("5", "6"))
	harapKode(t, kode, 409, "put bentrok")
	_, obj, _ = kirim(t, h, "GET", "/peminjaman/1", "")
	if obj["slot_mulai"] != float64(3) || obj["slot_selesai"] != float64(5) {
		t.Fatalf("data lama berubah: %v", obj)
	}

	kode, _, _ = kirim(t, h, "PUT", "/peminjaman/abc", bodyPinjam)
	harapKode(t, kode, 400, "put id bukan angka")
	kode, _, _ = kirim(t, h, "PUT", "/peminjaman/99", "bukan json")
	harapKode(t, kode, 404, "404 sebelum body")
	kode, _, _ = kirim(t, h, "PUT", "/peminjaman/1", "bukan json")
	harapKode(t, kode, 400, "body rusak")

	kode, _, raw := kirim(t, h, "DELETE", "/peminjaman/1", "")
	harapKode(t, kode, 204, "delete")
	if raw != "" {
		t.Fatalf("204 tidak boleh punya body: %q", raw)
	}
	kode, _, _ = kirim(t, h, "DELETE", "/peminjaman/1", "")
	harapKode(t, kode, 404, "delete ulang")
	kode, _, _ = kirim(t, h, "DELETE", "/peminjaman/abc", "")
	harapKode(t, kode, 400, "delete id bukan angka")
	kode, _, _ = kirim(t, h, "POST", "/peminjaman", bodyPinjam)
	harapKode(t, kode, 201, "slot lepas setelah delete")
}

func TestBlokir(t *testing.T) {
	h := handlerDenganRuang(t, "3")
	kirim(t, h, "POST", "/peminjaman", bodyPinjam)           // id 1, slot 3-4
	kirim(t, h, "POST", "/peminjaman", pinjamSlot("4", "5")) // id 2
	kirim(t, h, "POST", "/peminjaman", pinjamSlot("7", "8")) // id 3

	kode, _, _ := kirim(t, h, "POST", "/blokir", `{"sumber_daya_id":1,"tanggal":"2026-11-02","slot_mulai":2,"slot_selesai":3,"keterangan":" "}`)
	harapKode(t, kode, 400, "keterangan kosong")
	kode, _, _ = kirim(t, h, "POST", "/blokir", `{"sumber_daya_id":9,"tanggal":"2026-11-02","slot_mulai":2,"slot_selesai":3,"keterangan":"x"}`)
	harapKode(t, kode, 400, "sumber daya tidak ada")

	kode, _, raw := kirim(t, h, "POST", "/blokir", `{"sumber_daya_id":1,"tanggal":"2026-11-02","slot_mulai":4,"slot_selesai":4,"keterangan":"Servis"}`)
	harapKode(t, kode, 201, "buat blokir")
	var resp blokirResponse
	json.Unmarshal([]byte(raw), &resp)
	if len(resp.Dibatalkan) != 2 || resp.Dibatalkan[0] != 1 || resp.Dibatalkan[1] != 2 {
		t.Fatalf("dibatalkan salah: %s", raw)
	}
	_, obj, _ := kirim(t, h, "GET", "/peminjaman/1", "")
	if obj["status"] != "dibatalkan" {
		t.Fatalf("status harus dibatalkan: %v", obj)
	}
	_, obj, _ = kirim(t, h, "GET", "/peminjaman/3", "")
	if obj["status"] != "aktif" {
		t.Fatalf("peminjaman di luar blokir harus tetap aktif: %v", obj)
	}

	kode, _, _ = kirim(t, h, "POST", "/blokir", `{"sumber_daya_id":1,"tanggal":"2026-11-02","slot_mulai":4,"slot_selesai":6,"keterangan":"Lain"}`)
	harapKode(t, kode, 409, "blokir beririsan")
	kode, _, _ = kirim(t, h, "POST", "/peminjaman", bodyPinjam)
	harapKode(t, kode, 409, "peminjaman kena blokir")
	kode, _, _ = kirim(t, h, "PUT", "/peminjaman/3", pinjamSlot("4", "4"))
	harapKode(t, kode, 409, "put kena blokir")

	kode, _, raw = kirim(t, h, "POST", "/blokir", `{"sumber_daya_id":1,"tanggal":"2026-12-01","slot_mulai":1,"slot_selesai":1,"keterangan":"Praktikum"}`)
	harapKode(t, kode, 201, "blokir tanpa terdampak")
	if !strings.Contains(raw, `"dibatalkan":[]`) {
		t.Fatalf("dibatalkan harus []: %s", raw)
	}

	kode, _, _ = kirim(t, h, "DELETE", "/blokir/abc", "")
	harapKode(t, kode, 400, "hapus id bukan angka")
	kode, _, _ = kirim(t, h, "DELETE", "/blokir/99", "")
	harapKode(t, kode, 404, "hapus tidak ada")
	kode, _, _ = kirim(t, h, "DELETE", "/blokir/1", "")
	harapKode(t, kode, 204, "hapus blokir")
	kode, _, _ = kirim(t, h, "POST", "/peminjaman", bodyPinjam)
	harapKode(t, kode, 201, "bisa dipinjam lagi")
	_, obj, _ = kirim(t, h, "GET", "/peminjaman/1", "")
	if obj["status"] != "dibatalkan" {
		t.Fatalf("peminjaman tidak boleh pulih otomatis: %v", obj)
	}
	_, _, raw = kirim(t, h, "GET", "/blokir", "")
	var daftar []Blokir
	json.Unmarshal([]byte(raw), &daftar)
	if len(daftar) != 1 || daftar[0].ID != 2 {
		t.Fatalf("daftar blokir salah: %s", raw)
	}
}

func TestBersamaan(t *testing.T) {
	h := handlerDenganRuang(t, "1")
	srv := httptest.NewServer(h)
	defer srv.Close()

	const n = 40
	var wg sync.WaitGroup
	var mu sync.Mutex
	hitung := map[int]int{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Post(srv.URL+"/peminjaman", "application/json", strings.NewReader(bodyPinjam))
			if err != nil {
				t.Error(err)
				return
			}
			resp.Body.Close()
			mu.Lock()
			hitung[resp.StatusCode]++
			mu.Unlock()
		}()
	}
	wg.Wait()
	if hitung[201] != 1 || hitung[409] != n-1 {
		t.Fatalf("harusnya 1x201 dan %dx409, dapat %v", n-1, hitung)
	}
}
