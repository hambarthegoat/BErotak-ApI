package main

import "net/http"

func (s *Store) buatBlokir(w http.ResponseWriter, r *http.Request) {
	var in blokirInput
	if !decodeJSON(w, r, &in) {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	b, pesan := s.validasiBlokir(in)
	if pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}
	if s.bentrokDenganBlokir(b.SumberDayaID, b.Tanggal, b.SlotMulai, b.SlotSelesai) {
		writeError(w, http.StatusConflict, "slot beririsan dengan blokir lain")
		return
	}

	b.ID = s.nextBlokirID
	s.nextBlokirID++
	s.blokir = append(s.blokir, b)
	dibatalkan := s.batalkanTerdampak(b)
	writeJSON(w, http.StatusCreated, blokirResponse{Blokir: b, Dibatalkan: dibatalkan})
}

func (s *Store) daftarBlokir(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	writeJSON(w, http.StatusOK, s.blokir)
}

func (s *Store) hapusBlokir(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.indeksBlokir(id)
	if i < 0 {
		writeError(w, http.StatusNotFound, "blokir tidak ditemukan")
		return
	}
	s.blokir = append(s.blokir[:i], s.blokir[i+1:]...)
	w.WriteHeader(http.StatusNoContent)
}
