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

// peminjaman baru berid 0, jadi tidak ada yang dikecualikan
func (s *Store) kapasitasTerlampaui(p Peminjaman, stok int) bool {
	for slot := p.SlotMulai; slot <= p.SlotSelesai; slot++ {
		if s.unitTerpakai(p.SumberDayaID, p.Tanggal, slot, p.ID)+p.Jumlah > stok {
			return true
		}
	}
	return false
}
