# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

Catatan umum: kode di repo ini mostly digenerate AI (Claude Code) berdasarkan dokumen Requirement dan Kartu Persona (persona 5 = E). Alurnya: saya bertanya tentang requirement secara sederhana, lalu saya bertanya bagaimana implementasi, lalu saya bertanya tentang code nya. Semua dikerjakan dalam satu percakapan. Awalnya AI sempat menulis sampai fitur persona, lalu saya minta dipangkas ke B1-B3 (checkpoint 1) dan baru ditambah B4-B7 + blokir di checkpoint 2. Prompt di bawah adalah ringkasan permintaan per bagian, urut sesuai pengerjaan.

## 1. Rencana implementasi
- Prompt: "Plan-kan dulu implementasi tugas ini sesuai Requirement dan kartu persona 5. Jelaskan struktur file dan kenapa dipilih begitu."
- Dipakai: pembagian file (main, models, store, helpers, validation, conflict, handlers) dan urutan pemeriksaan 400 → 404 → 409. Belum ada kode.
- Alasan dipilih: satu `Store` berisi slice dan counter id per jenis, karena daftar harus terurut menurut id naik dan data kecil di memori. Input memakai pointer supaya field yang tidak dikirim bisa dibedakan dari `0` atau `""`.
- Diubah: saya menambah syarat agar AI-LOG ditulis sesuai kenyataan, tidak seolah-olah AI dipakai sedikit.

## 2. Struktur data dan setup (B1)
- Prompt: "Coba generate-kan `go mod init`, struct sumber daya dan peminjaman, dan Store-nya. Jelaskan kenapa dibuat begini."
- Dipakai: `go.mod`, `models.go` (`SumberDaya`, `Peminjaman`, struct input berpointer), `store.go` (`Store`, `NewStore`, `cariSumberDaya`, `indeksPeminjaman`).
- Alasan dipilih: tag JSON persis seperti kontrak supaya frontend tidak rusak; `NewStore()` dipanggil dari `NewHandler()` supaya tiap handler punya penyimpanan kosong sendiri.
- Diubah: directive `go` di `go.mod` diturunkan ke `go 1.22` sesuai syarat minimal.

## 3. Kelola sumber daya, buat dan lihat peminjaman (B1-B3)
- Prompt: "Coba generate-kan endpoint `POST/GET /sumber-daya`, `GET /sumber-daya/{id}`, `POST/GET /peminjaman` dengan filter, dan `GET /peminjaman/{id}`. Jelaskan bagaimana alurnya dan kenapa memakai `http.NewServeMux`."
- Dipakai: `main.go` (`NewHandler`, `main`), `handlers_sumberdaya.go`, `handlers_peminjaman.go` (`buatPeminjaman`, `daftarPeminjaman`, `detailPeminjaman`), `helpers.go` (`writeJSON`, `writeError`, `parseID`, `decodeJSON`).
- Alasan dipilih: `NewServeMux` dengan pola `"POST /peminjaman"` dan `r.PathValue("id")` tersedia sejak Go 1.22 dan itu yang diminta spesifikasi; `NewHandler()` wajib ada supaya tes dosen bisa memanggilnya langsung. Daftar dimulai dari `[]Peminjaman{}` supaya kosong menjadi `[]`, bukan `null`.
- Diubah: bagian ini awalnya sudah memuat bentrok, PUT/DELETE, dan blokir; saya minta dibuang dulu karena di luar checkpoint 1. Saya juga minta menghapus yang tidak terpakai, termasuk mutex dan validasi aturan isi, yang ditambahkan lagi di entri 4 dan 6.

## 4. Validasi dan error JSON (B4)
- Prompt: "Coba generate-kan validasi semua input sesuai kontrak (trim, maksimal 50 karakter, tanggal harus ada, slot 1-8, jumlah 1..stok) dan error selalu JSON. Jelaskan kenapa cara cek ini yang dipilih."
- Dipakai: `validation.go` (`validasiTeks`, `validasiTanggal`, `validasiSlot`, `validasiJadwal`, `validasiSumberDaya`, `validasiPeminjaman`), `notFoundHandler` di `helpers.go`.
- Alasan dipilih: panjang dihitung dengan `utf8.RuneCountInString` karena spesifikasi menghitung per karakter, bukan byte; tanggal diperiksa dengan `time.Parse("2006-01-02")` supaya 30 Februari ditolak dan 29 Februari 2028 diterima; `*int` membuat `1.5` atau `"1"` otomatis gagal decode (400).
- Diubah: `validasiPeminjaman` sekarang juga mengembalikan `SumberDaya` supaya stok bisa dipakai cek bentrok tanpa mencari dua kali.

## 5. Deteksi bentrok (B5)
- Prompt: "Coba generate-kan fungsi cek bentrok kapasitas. Jelaskan bagaimana kerjanya dan kenapa dicek per slot, bukan per rentang."
- Dipakai: `unitTerpakai` dan `kapasitasTerlampaui` di `conflict.go`, dipanggil dari `buatPeminjaman`.
- Alasan dipilih: dijumlahkan per slot karena alat bisa berstok lebih dari 1; membandingkan slot_mulai/slot_selesai saja tidak cukup untuk rentang yang beririsan sebagian. Hanya peminjaman `aktif` yang dihitung.
- Diubah: belum ada.

## 6. Ubah dan hapus (B6)
- Prompt: "Coba generate-kan `PUT` dan `DELETE /peminjaman/{id}`. Jelaskan bagaimana PUT menjaga aturan bentrok dan kenapa data lama aman kalau 409."
- Dipakai: `ubahPeminjaman` dan `hapusPeminjaman` di `handlers_peminjaman.go`; `idPengecualian` di `unitTerpakai`.
- Alasan dipilih: peminjaman yang diubah dikecualikan lewat id-nya supaya tidak bentrok dengan dirinya sendiri; perubahan dibuat di salinan `p` dan baru menimpa `s.peminjaman[i]` setelah semua cek lolos. Body dibaca sebelum kunci tapi errornya baru dilaporkan setelah cek 404, supaya urutan 400 → 404 → body → 409 terpenuhi.
- Diubah: belum ada.

## 7. Aman untuk permintaan bersamaan (B7)
- Prompt: "Coba generate-kan perlindungan data bersama untuk request bersamaan. Jelaskan kenapa mutex dan bagaimana letaknya."
- Dipakai: `Store.mu`; setiap handler `Lock()` lalu `defer Unlock()`.
- Alasan dipilih: satu mutex untuk seluruh `Store`, sehingga cek bentrok dan penyimpanan terjadi dalam satu kunci; dua request yang berebut slot yang sama menghasilkan satu 201 dan satu 409. Satu kunci lebih sederhana daripada kunci per sumber daya, dan data kecil.
- Diubah: belum ada.

## 8. Fitur khas: blokir jadwal (persona E)
- Prompt: "Coba generate-kan `POST/GET/DELETE /blokir` untuk persona E. Jelaskan bagaimana peminjaman yang terdampak dibatalkan dan kenapa statusnya `dibatalkan`, bukan dihapus."
- Dipakai: `handlers_blokir.go`, `validasiBlokir`, `bentrokDenganBlokir` dan `batalkanTerdampak` di `conflict.go`.
- Alasan dipilih: spesifikasi meminta data peminjaman tetap ada dengan status `dibatalkan` yang tidak menahan slot; slice peminjaman terurut menurut id, jadi `dibatalkan` otomatis naik. Blokir yang dihapus tidak memulihkan peminjaman.
- Diubah: keputusan AI yang belum saya periksa: PUT pada peminjaman `dibatalkan` tetap dicek ke blokir, tapi cek kapasitas dilewati.

## 9. Tanya penjelasan kode (tidak menghasilkan kode)
- Prompt: "Kenapa ada mux di main? Apakah ini untuk implementasi B1-B3?"
- Dipakai: penjelasan bahwa `http.NewServeMux` dibutuhkan untuk memetakan method dan path ke handler, untuk `r.PathValue("id")` di `parseID`, dan karena `NewHandler()` wajib ada di kontrak. Tidak ada perubahan kode.
- Diubah: tidak ada.

## 10. Test
- Prompt: "Coba generate-kan test `httptest` untuk tiap fitur, termasuk 40 request bersamaan ke slot yang sama. Jelaskan apa yang dites tiap fungsi."
- Dipakai: `main_test.go`.
- Alasan dipilih: API belum bisa dicoba, jadi `httptest` memanggil `NewHandler()` tanpa membuka port.
- Diubah: test B2/B3 diubah supaya membuat sumber daya dulu, karena `sumber_daya_id` sekarang harus menunjuk sumber daya yang ada. Test dijalankan dan lolos, lihat entri 12.

## 11. Change request B: bentrok_dengan
- Prompt: "Ini lembar change request-nya (Lembar CR-B Daftar Bentrok), kerjakan sesuai dokumen. Jelaskan bagaimana dan kenapa dibuat begitu."
- Dipakai: `idPenghalang` di `conflict.go`, `writeBentrok` di `helpers.go`, perubahan pemanggilan di `buatPeminjaman`, dan `cr_test.go`.
- Alasan dipilih: daftar id mengambil semua peminjaman aktif yang beririsan, sesuai contoh di lembar (stok 1, slot 3-5 menghasilkan `[1, 2]`). Fungsinya terpisah dari `kapasitasTerlampaui` supaya aturan penolakan tidak berubah. PUT sengaja tidak diubah karena lembar hanya meminta POST. Untuk blokir saya memakai `[]` karena lembar membolehkannya.
- Diubah: belum ada. Test dijalankan dan lolos, lihat entri 12.

## 12. Jalankan dan uji sendiri
- Prompt: "Jalankan sendiri, cek apakah berjalan sesuai permintaan. Jalankan server sendiri, buat request sendiri."
- Dipakai: `go vet`, `go test -race`, lalu server dijalankan dan dicoba dengan curl: sumber daya, contoh di lembar CR, PUT/DELETE, blokir, dan 50 request serentak ke slot yang sama. Hasilnya sesuai: `[1, 2]` dan `[1]` untuk contoh CR, `dibatalkan: [1, 2]` untuk blokir, dan tepat satu 201 dari 50 request serentak.
- Diubah: binary `pinjamkampus` hasil `go build ./...` dihapus (lihat Bug 1 di PROSES.md). Tidak ada perubahan di kode.

## 13. Rapikan kode dan cek ke requirement
- Prompt: "Kurangi komentarnya, perbaiki kodenya dan pastikan kode-kodenya sudah sesuai dengan requirement."
- Dipakai: komentar dikurangi dan tersisa yang menjelaskan alasan non-obvious (`idPengecualian`, kenapa error body PUT ditunda); `validasiJadwal` di `validation.go` menggabungkan pengecekan tanggal+slot yang sebelumnya duplikat di `validasiPeminjaman` dan `validasiBlokir`; `errBodyTidakValid` dipindah ke atas `helpers.go`.
- Diubah: dicek ulang ke dokumen Requirement, kartu persona E, dan lembar CR (urutan 400 → 404 → 409, `[]` bukan `null`, id naik, 204 tanpa body, `bentrok_dengan` hanya pada 409 kapasitas/blokir POST). Tidak ada perbedaan perilaku yang ditemukan.
