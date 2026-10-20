package main

import (
	"net/http"
	"strconv"
)

func (s *Store) buatPeminjaman(w http.ResponseWriter, r *http.Request) {
	var in peminjamanInput
	if !decodeJSON(w, r, &in) {
		return
	}
	p, pesan := validasiPeminjaman(in)
	if pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
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
	p, ada := s.cariPeminjaman(id)
	if !ada {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, p)
}
