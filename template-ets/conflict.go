package main

func beririsan(a1, a2, b1, b2 int) bool {
	return a1 <= b2 && b1 <= a2
}

// idPengecualian dilewati supaya PUT tidak bentrok dengan dirinya sendiri
func (s *Store) unitTerpakai(sumberDayaID int, tanggal string, slot, idPengecualian int) int {
	total := 0
	for _, p := range s.peminjaman {
		if p.ID == idPengecualian || p.Status != statusAktif {
			continue
		}
		if p.SumberDayaID == sumberDayaID && p.Tanggal == tanggal &&
			slot >= p.SlotMulai && slot <= p.SlotSelesai {
			total += p.Jumlah
		}
	}
	return total
}

func (s *Store) idPenghalang(p Peminjaman) []int {
	ids := []int{}
	for _, lain := range s.peminjaman {
		if lain.ID == p.ID || lain.Status != statusAktif {
			continue
		}
		if lain.SumberDayaID == p.SumberDayaID && lain.Tanggal == p.Tanggal &&
			beririsan(p.SlotMulai, p.SlotSelesai, lain.SlotMulai, lain.SlotSelesai) {
			ids = append(ids, lain.ID)
		}
	}
	return ids
}

// peminjaman baru berid 0, jadi tidak ada yang dikecualikan
func (s *Store) kapasitasTerlampaui(p Peminjaman, stok int) bool {
	for slot := p.SlotMulai; slot <= p.SlotSelesai; slot++ {
		if s.unitTerpakai(p.SumberDayaID, p.Tanggal, slot, p.ID)+p.Jumlah > stok {
			return true
		}
	}
	return false
}

func (s *Store) bentrokDenganBlokir(sumberDayaID int, tanggal string, mulai, selesai int) bool {
	for _, b := range s.blokir {
		if b.SumberDayaID == sumberDayaID && b.Tanggal == tanggal &&
			beririsan(mulai, selesai, b.SlotMulai, b.SlotSelesai) {
			return true
		}
	}
	return false
}

func (s *Store) batalkanTerdampak(b Blokir) []int {
	ids := []int{}
	for i := range s.peminjaman {
		p := &s.peminjaman[i]
		if p.Status != statusAktif {
			continue
		}
		if p.SumberDayaID == b.SumberDayaID && p.Tanggal == b.Tanggal &&
			beririsan(p.SlotMulai, p.SlotSelesai, b.SlotMulai, b.SlotSelesai) {
			p.Status = statusDibatalkan
			ids = append(ids, p.ID)
		}
	}
	return ids
}
