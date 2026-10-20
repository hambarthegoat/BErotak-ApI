package main

import (
	"net/http"
	"strconv"
)

const (
	pesanBentrokKapasitas = "kapasitas sumber daya terlampaui pada salah satu slot"
	pesanBentrokBlokir    = "slot yang diminta diblokir"
)

func (s *Store) buatPeminjaman(w http.ResponseWriter, r *http.Request) {
	var in peminjamanInput
	if !decodeJSON(w, r, &in) {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	p, sd, pesan := s.validasiPeminjaman(in)
	if pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}
	if s.bentrokDenganBlokir(p.SumberDayaID, p.Tanggal, p.SlotMulai, p.SlotSelesai) {
		writeError(w, http.StatusConflict, pesanBentrokBlokir)
		return
	}
	if s.kapasitasTerlampaui(p, sd.Stok) {
		writeError(w, http.StatusConflict, pesanBentrokKapasitas)
		return
	}

	p.ID = s.nextPeminjamanID
	s.nextPeminjamanID++
	p.Status = statusAktif
	s.peminjaman = append(s.peminjaman, p)
	writeJSON(w, http.StatusCreated, p)
}

func (s *Store) daftarPeminjaman(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filterSD := q.Get("sumber_daya_id")
	filterTanggal := q.Get("tanggal")
	filterStatus := q.Get("status")

	sdID := 0
	if filterSD != "" {
		n, err := strconv.Atoi(filterSD)
		if err != nil {
			writeError(w, http.StatusBadRequest, "sumber_daya_id harus berupa angka")
			return
		}
		sdID = n
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	hasil := []Peminjaman{}
	for _, p := range s.peminjaman {
		if filterSD != "" && p.SumberDayaID != sdID {
			continue
		}
		if filterTanggal != "" && p.Tanggal != filterTanggal {
			continue
		}
		if filterStatus != "" && p.Status != filterStatus {
			continue
		}
		hasil = append(hasil, p)
	}
	writeJSON(w, http.StatusOK, hasil)
}

func (s *Store) detailPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.indeksPeminjaman(id)
	if i < 0 {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, s.peminjaman[i])
}

func (s *Store) ubahPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var in peminjamanInput
	// error body ditunda: 404 harus dicek lebih dulu daripada 400
	errBody := parseJSON(r, &in)

	s.mu.Lock()
	defer s.mu.Unlock()

	i := s.indeksPeminjaman(id)
	if i < 0 {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	if errBody != nil {
		writeError(w, http.StatusBadRequest, errBodyTidakValid)
		return
	}
	p, sd, pesan := s.validasiPeminjaman(in)
	if pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}

	p.ID = id
	p.Status = s.peminjaman[i].Status
	if s.bentrokDenganBlokir(p.SumberDayaID, p.Tanggal, p.SlotMulai, p.SlotSelesai) {
		writeError(w, http.StatusConflict, pesanBentrokBlokir)
		return
	}
	if p.Status == statusAktif && s.kapasitasTerlampaui(p, sd.Stok) {
		writeError(w, http.StatusConflict, pesanBentrokKapasitas)
		return
	}
	s.peminjaman[i] = p
	writeJSON(w, http.StatusOK, p)
}

func (s *Store) hapusPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.indeksPeminjaman(id)
	if i < 0 {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	s.peminjaman = append(s.peminjaman[:i], s.peminjaman[i+1:]...)
	w.WriteHeader(http.StatusNoContent)
}
