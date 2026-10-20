# PROSES

## 1. Persona dan masalahnya
Persona saya E, asisten laboratorium. Lab sering dipakai praktikum atau diservis di jam tertentu, tapi tetap bisa dipinjam orang lain lewat aplikasi, jadi bentrok dengan kegiatan lab. Fitur blokir jadwal membuat slot tertentu tidak bisa dipinjam, dan peminjaman yang sudah terlanjur ada di slot itu otomatis dibatalkan.

(Draf, tulis ulang dengan kata-katamu sendiri.)

## 2. Rancangan struktur data
Semua data ada di satu struct `Store` (store.go):
- `sumberDaya []SumberDaya` dan `peminjaman []Peminjaman`: slice biasa. Daftar harus terurut menurut id naik dan id hanya bertambah, jadi urutan append sudah terurut.
- `nextSumberDayaID` dan `nextPeminjamanID`: counter id terpisah per jenis, mulai dari 1. Counter hanya naik setelah input lolos, jadi request yang gagal tidak menghabiskan id.
- `Peminjaman.Status` bernilai `aktif` saat dibuat, dan `dibatalkan` kalau terkena blokir. Hanya `aktif` yang menahan slot.
- `blokir []Blokir` dan `nextBlokirID`: blokir disimpan terpisah dari peminjaman dengan counter id sendiri.
- `mu sync.Mutex` di `Store`: satu kunci untuk semua data, supaya cek bentrok dan penyimpanan terjadi dalam satu kunci yang sama.

Pencarian dilakukan linear karena data kecil dan di memori.

## 3. Satu keputusan desain
Cek kapasitas dilakukan per slot dengan menjumlahkan unit (`unitTerpakai` di conflict.go), bukan sekadar membandingkan slot_mulai/slot_selesai antar peminjaman. Alternatif membandingkan rentang tidak bisa menangani alat dengan stok lebih dari 1 dan mudah salah pada rentang yang beririsan sebagian. Per slot juga memudahkan PUT: peminjaman yang sedang diubah cukup dilewati lewat id-nya.

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala: setelah `go build ./...`, muncul file binary `pinjamkampus` di folder `template-ets` yang ikut terbaca sebagai file baru di `git status`.
- Cara menemukan: `git status --short` setelah menjalankan build dan server.
- Penyebab dan solusi: karena hanya ada satu package `main`, `go build ./...` menulis binary ke folder kerja. Binary saya hapus dan untuk pengecekan berikutnya dipakai `go build -o /dev/null .`, sehingga tidak ada file yang tertinggal dan ikut ter-commit.

### Bug 2
- Gejala: test B2/B3 dari checkpoint 1 akan gagal dengan 400 setelah validasi (B4) ditambahkan, padahal fungsinya benar.
- Cara menemukan: saat menambah validasi, saya membaca ulang test lama dan melihat `bodyPinjam` memakai `sumber_daya_id: 1` tanpa membuat sumber daya terlebih dulu. Ini ketahuan dari membaca kode, sebelum test dijalankan.
- Penyebab dan solusi: sebelum B4, POST tidak memeriksa apakah sumber daya ada, jadi test lama lolos tanpa membuatnya. Test diubah dengan helper `handlerDenganRuang` yang membuat sumber daya lebih dulu.

## 5. Change request
Change Request B: pada `POST /peminjaman`, 409 karena slot atau kapasitas harus memuat `bentrok_dengan`, yaitu id peminjaman yang menghalangi (naik), selain `error`.

Yang diubah: `idPenghalang` baru di conflict.go, `writeBentrok` di helpers.go, dan dua pemanggilan di `buatPeminjaman`. PUT dan endpoint lain tidak disentuh.

Keputusan desain: daftar diambil dari semua peminjaman aktif yang beririsan (sumber daya, tanggal, dan rentang slot sama), bukan hanya yang membuat kapasitas penuh, sesuai contoh di lembar CR. Pencarian ini terpisah dari `kapasitasTerlampaui` supaya aturan penolakan tidak berubah, dan hanya dipanggil setelah penolakan diputuskan. Untuk 409 karena blokir, field diisi `[]`, yang diizinkan lembar CR.
