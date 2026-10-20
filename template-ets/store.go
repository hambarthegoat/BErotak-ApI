package main

import "sync"

type Store struct {
	mu         sync.Mutex
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

func (s *Store) indeksPeminjaman(id int) int {
	for i, p := range s.peminjaman {
		if p.ID == id {
			return i
		}
	}
	return -1
}
