# PinjamKampus API

- Nama dan NRP: Muhammad Ilham Akbar / 5053241009
- Persona: 5 (E, Asisten laboratorium: blokir jadwal)
- Cara menjalankan: `go run .` lalu buka http://localhost:8080

## Fitur
- [x] B1 Sumber daya
- [x] B2 Buat peminjaman
- [x] B3 Lihat peminjaman (daftar, filter, detail)
- [x] B4 Validasi dan error JSON
- [x] B5 Deteksi bentrok (409)
- [x] B6 Ubah dan hapus
- [x] B7 Aman untuk permintaan bersamaan
- [x] Fitur khas persona (blokir jadwal: `POST/GET/DELETE /blokir`)
- [x] Change request B (`bentrok_dengan` pada 409 `POST /peminjaman`)

## Pengujian
- `go vet ./...` bersih dan `go test -race ./...` lolos (test otomatis di `main_test.go` dan `cr_test.go`).
- Server juga dicoba langsung dengan curl: contoh di lembar CR (`[1, 2]` dan `[1]`), blokir, PUT/DELETE, dan 50 request serentak ke slot yang sama (tepat satu 201, sisanya 409).

## Yang belum selesai
- Tidak ada fitur yang sengaja ditinggalkan.
