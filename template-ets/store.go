package main

type Store struct {
	sumberDaya []SumberDaya
	peminjaman []Peminjaman

	nextSumberDayaID int
	nextPeminjamanID int
}

func NewStore() *Store {
	return &Store{
		sumberDaya:       []SumberDaya{},
		peminjaman:       []Peminjaman{},
		nextSumberDayaID: 1,
		nextPeminjamanID: 1,
	}
}

func (s *Store) cariSumberDaya(id int) (SumberDaya, bool) {
	for _, sd := range s.sumberDaya {
		if sd.ID == id {
			return sd, true
		}
	}
	return SumberDaya{}, false
}

func (s *Store) cariPeminjaman(id int) (Peminjaman, bool) {
	for _, p := range s.peminjaman {
		if p.ID == id {
			return p, true
		}
	}
	return Peminjaman{}, false
}
