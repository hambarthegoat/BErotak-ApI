# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

Catatan umum: kode di repo ini mostly digenerate AI (Claude Code) berdasarkan dokumen Requirement dan Kartu Persona (persona 5 = E). Alurnya: saya bertanya tentang requirement secara sederhana, lalu saya bertanya bagaimana implementasi, lalu saya bertanya tentang code nya. 

## Setup dan struktur data (persiapan B1)
- Prompt: minta rencana implementasi dari dokumen Requirement + Kartu Persona, lalu `go mod init pinjamkampus` dan struct untuk sumber daya dan peminjaman.
- Dipakai: `go.mod`, `models.go` (struct `SumberDaya`, `Peminjaman`, struct input berpointer), `store.go` (`Store`, `NewStore`, `cariSumberDaya`, `cariPeminjaman`, counter id).
- Diubah: directive `go` di go.mod diturunkan dari versi yang terpasang ke `go 1.22` supaya sesuai syarat minimal. Selain itu belum ada.

## Kelola sumber daya (B1)
- Prompt: minta `POST /sumber-daya`, `GET /sumber-daya`, `GET /sumber-daya/{id}` dengan `http.NewServeMux` pola Go 1.22 dan fungsi `NewHandler()`.
- Dipakai: `main.go` (`NewHandler`, `main`), `handlers_sumberdaya.go`, `helpers.go` (`writeJSON`, `writeError`, `parseID`, `decodeJSON`), `validasiSumberDaya` di `validation.go` (baru cek field wajib ada).
- Diubah: belum ada.

## Buat peminjaman (B2)
- Prompt: minta `POST /peminjaman` dengan enam field dan status awal `aktif`.
- Dipakai: `buatPeminjaman` di `handlers_peminjaman.go`, `validasiPeminjaman` di `validation.go` (baru cek field wajib ada).
- Diubah: belum ada. Aturan isi (B4) dan cek bentrok (B5) belum ada.

## Lihat peminjaman (B3)
- Prompt: minta `GET /peminjaman` dengan filter `sumber_daya_id`, `tanggal`, `status` (boleh digabung, kosong → `[]`) dan `GET /peminjaman/{id}`.
- Dipakai: `daftarPeminjaman` dan `detailPeminjaman` di `handlers_peminjaman.go`.
- Diubah: belum ada.

## Test
- Prompt: minta test `httptest` untuk B1-B3. (test untuk memastikan kebenaran function, karena api masih belum bisa dicoba)
- Dipakai: `main_test.go`.
- Diubah: belum dijalankan, jadi belum ada perubahan.
