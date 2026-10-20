package main

import "net/http"

func (s *Store) buatSumberDaya(w http.ResponseWriter, r *http.Request) {
	var in sumberDayaInput
	if !decodeJSON(w, r, &in) {
		return
	}
	sd, pesan := validasiSumberDaya(in)
	if pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	sd.ID = s.nextSumberDayaID
	s.nextSumberDayaID++
	s.sumberDaya = append(s.sumberDaya, sd)
	writeJSON(w, http.StatusCreated, sd)
}

func (s *Store) daftarSumberDaya(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	writeJSON(w, http.StatusOK, s.sumberDaya)
}

func (s *Store) detailSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	sd, ada := s.cariSumberDaya(id)
	if !ada {
		writeError(w, http.StatusNotFound, "sumber daya tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, sd)
}
