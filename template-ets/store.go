package main

import "sync"

type Store struct {
	mu         sync.Mutex
	sumberDaya []SumberDaya
	peminjaman []Peminjaman
	blokir     []Blokir

	nextSumberDayaID int
	nextPeminjamanID int
	nextBlokirID     int
}

func NewStore() *Store {
	return &Store{
		sumberDaya:       []SumberDaya{},
		peminjaman:       []Peminjaman{},
		blokir:           []Blokir{},
		nextSumberDayaID: 1,
		nextPeminjamanID: 1,
		nextBlokirID:     1,
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

func (s *Store) indeksBlokir(id int) int {
	for i, b := range s.blokir {
		if b.ID == id {
			return i
		}
	}
	return -1
}
