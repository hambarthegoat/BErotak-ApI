# PROSES

## 1. Persona dan masalahnya
Persona saya E, asisten laboratorium. Lab sering dipakai praktikum atau diservis di jam tertentu, tapi tetap bisa dipinjam orang lain lewat aplikasi, jadi bentrok dengan kegiatan lab. Fitur blokir jadwal membuat slot tertentu tidak bisa dipinjam, dan peminjaman yang sudah terlanjur ada di slot itu otomatis dibatalkan.

(Draf, tulis ulang dengan kata-katamu sendiri.)

## 2. Rancangan struktur data
Semua data ada di satu struct `Store` (store.go):
- `sumberDaya []SumberDaya` dan `peminjaman []Peminjaman`: slice biasa. Daftar harus terurut menurut id naik dan id hanya bertambah, jadi urutan append sudah terurut.
- `nextSumberDayaID` dan `nextPeminjamanID`: counter id terpisah per jenis, mulai dari 1. Counter hanya naik setelah input lolos, jadi request yang gagal tidak menghabiskan id.
- `Peminjaman.Status` bernilai `aktif` saat dibuat.

Pencarian dilakukan linear karena data kecil dan di memori. Struktur untuk blokir (persona E) dan mutex belum ditambahkan; catat di sini kalau sudah.

## 3. Satu keputusan desain
(Satu keputusan yang kamu ambil, pilihan lain yang kamu pertimbangkan, dan alasanmu memilih yang ini. Contoh: kenapa cek bentrok per slot, bukan per rentang.)

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
