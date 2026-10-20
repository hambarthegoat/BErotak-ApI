# PinjamKampus API

- Nama dan NRP: Muhammad Ilham Akbar / 5053241009
- Persona: 5 (E, Asisten laboratorium: blokir jadwal)
- Cara menjalankan: `go run .` lalu buka http://localhost:8080

## Fitur
- [x] B1 Sumber daya
- [x] B2 Buat peminjaman
- [x] B3 Lihat peminjaman (daftar, filter, detail)
- [ ] B4 Validasi dan error JSON
- [ ] B5 Deteksi bentrok (409)
- [ ] B6 Ubah dan hapus
- [ ] B7 Aman untuk permintaan bersamaan
- [ ] Fitur khas persona
- [ ] Change request

## Yang belum selesai
- B4-B7, fitur khas persona (blokir jadwal), dan change request belum dikerjakan.
- Untuk sementara POST hanya memeriksa field wajib ada; aturan isi (panjang, format tanggal, rentang slot, stok) belum diperiksa.
- Belum aman untuk request bersamaan (tanpa mutex).
- Kode baru ditulis dan belum dijalankan/diuji (`go build`, `go vet`, `go test`).
