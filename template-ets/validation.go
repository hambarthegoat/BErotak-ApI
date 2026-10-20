package main

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

func validasiTeks(nilai *string, label string, maks int) (string, string) {
	if nilai == nil {
		return "", label + " wajib diisi"
	}
	teks := strings.TrimSpace(*nilai)
	if teks == "" {
		return "", label + " tidak boleh kosong"
	}
	if utf8.RuneCountInString(teks) > maks {
		return "", fmt.Sprintf("%s maksimal %d karakter", label, maks)
	}
	return teks, ""
}

func validasiTanggal(nilai *string) (string, string) {
	if nilai == nil {
		return "", "tanggal wajib diisi"
	}
	teks := strings.TrimSpace(*nilai)
	if len(teks) != len("2006-01-02") {
		return "", "tanggal harus berformat YYYY-MM-DD"
	}
	if _, err := time.Parse("2006-01-02", teks); err != nil {
		return "", "tanggal tidak valid, gunakan format YYYY-MM-DD dengan tanggal yang ada"
	}
	return teks, ""
}

func validasiSlot(mulai, selesai *int) string {
	if mulai == nil || selesai == nil {
		return "slot_mulai dan slot_selesai wajib diisi"
	}
	if *mulai < slotMin || *mulai > slotMax || *selesai < slotMin || *selesai > slotMax {
		return "slot harus antara 1 sampai 8"
	}
	if *mulai > *selesai {
		return "slot_mulai tidak boleh lebih besar dari slot_selesai"
	}
	return ""
}

func validasiJadwal(tanggal *string, mulai, selesai *int) (string, string) {
	tgl, pesan := validasiTanggal(tanggal)
	if pesan != "" {
		return "", pesan
	}
	return tgl, validasiSlot(mulai, selesai)
}

func validasiSumberDaya(in sumberDayaInput) (SumberDaya, string) {
	nama, pesan := validasiTeks(in.Nama, "nama", maxPanjangNama)
	if pesan != "" {
		return SumberDaya{}, pesan
	}
	if in.Jenis == nil {
		return SumberDaya{}, "jenis wajib diisi"
	}
	jenis := strings.TrimSpace(*in.Jenis)
	if jenis != jenisRuang && jenis != jenisAlat {
		return SumberDaya{}, `jenis harus "ruang" atau "alat"`
	}
	if in.Stok == nil {
		return SumberDaya{}, "stok wajib diisi"
	}
	if *in.Stok < 1 {
		return SumberDaya{}, "stok minimal 1"
	}
	return SumberDaya{Nama: nama, Jenis: jenis, Stok: *in.Stok}, ""
}

func (s *Store) validasiPeminjaman(in peminjamanInput) (Peminjaman, SumberDaya, string) {
	if in.SumberDayaID == nil {
		return Peminjaman{}, SumberDaya{}, "sumber_daya_id wajib diisi"
	}
	peminjam, pesan := validasiTeks(in.Peminjam, "peminjam", maxPanjangNama)
	if pesan != "" {
		return Peminjaman{}, SumberDaya{}, pesan
	}
	tanggal, pesan := validasiJadwal(in.Tanggal, in.SlotMulai, in.SlotSelesai)
	if pesan != "" {
		return Peminjaman{}, SumberDaya{}, pesan
	}
	if in.Jumlah == nil {
		return Peminjaman{}, SumberDaya{}, "jumlah wajib diisi"
	}
	sd, ada := s.cariSumberDaya(*in.SumberDayaID)
	if !ada {
		return Peminjaman{}, SumberDaya{}, "sumber daya tidak ditemukan"
	}
	if *in.Jumlah < 1 || *in.Jumlah > sd.Stok {
		return Peminjaman{}, SumberDaya{}, fmt.Sprintf("jumlah harus antara 1 dan %d", sd.Stok)
	}
	return Peminjaman{
		SumberDayaID: sd.ID,
		Peminjam:     peminjam,
		Tanggal:      tanggal,
		SlotMulai:    *in.SlotMulai,
		SlotSelesai:  *in.SlotSelesai,
		Jumlah:       *in.Jumlah,
	}, sd, ""
}

func (s *Store) validasiBlokir(in blokirInput) (Blokir, string) {
	if in.SumberDayaID == nil {
		return Blokir{}, "sumber_daya_id wajib diisi"
	}
	tanggal, pesan := validasiJadwal(in.Tanggal, in.SlotMulai, in.SlotSelesai)
	if pesan != "" {
		return Blokir{}, pesan
	}
	keterangan, pesan := validasiTeks(in.Keterangan, "keterangan", maxPanjangKet)
	if pesan != "" {
		return Blokir{}, pesan
	}
	if _, ada := s.cariSumberDaya(*in.SumberDayaID); !ada {
		return Blokir{}, "sumber daya tidak ditemukan"
	}
	return Blokir{
		SumberDayaID: *in.SumberDayaID,
		Tanggal:      tanggal,
		SlotMulai:    *in.SlotMulai,
		SlotSelesai:  *in.SlotSelesai,
		Keterangan:   keterangan,
	}, ""
}
