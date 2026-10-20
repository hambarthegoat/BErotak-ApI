package main

import "strings"

func validasiSumberDaya(in sumberDayaInput) (SumberDaya, string) {
	if in.Nama == nil || in.Jenis == nil || in.Stok == nil {
		return SumberDaya{}, "nama, jenis, dan stok wajib diisi"
	}
	return SumberDaya{
		Nama:  strings.TrimSpace(*in.Nama),
		Jenis: strings.TrimSpace(*in.Jenis),
		Stok:  *in.Stok,
	}, ""
}

func validasiPeminjaman(in peminjamanInput) (Peminjaman, string) {
	if in.SumberDayaID == nil || in.Peminjam == nil || in.Tanggal == nil ||
		in.SlotMulai == nil || in.SlotSelesai == nil || in.Jumlah == nil {
		return Peminjaman{}, "sumber_daya_id, peminjam, tanggal, slot_mulai, slot_selesai, dan jumlah wajib diisi"
	}
	return Peminjaman{
		SumberDayaID: *in.SumberDayaID,
		Peminjam:     strings.TrimSpace(*in.Peminjam),
		Tanggal:      strings.TrimSpace(*in.Tanggal),
		SlotMulai:    *in.SlotMulai,
		SlotSelesai:  *in.SlotSelesai,
		Jumlah:       *in.Jumlah,
	}, ""
}
