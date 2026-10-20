package main

import "testing"

func TestChangeRequestBentrokDengan(t *testing.T) {
	h := handlerDenganRuang(t, "1")
	kirim(t, h, "POST", "/peminjaman", pinjamSlot("2", "3")) // id 1
	kirim(t, h, "POST", "/peminjaman", pinjamSlot("5", "6")) // id 2
	// tanggal lain, tidak boleh ikut
	kirim(t, h, "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"Budi","tanggal":"2026-11-03","slot_mulai":3,"slot_selesai":5,"jumlah":1}`) // id 3

	cek := func(slotMulai, slotSelesai string, harap []any) {
		t.Helper()
		kode, obj, _ := kirim(t, h, "POST", "/peminjaman", pinjamSlot(slotMulai, slotSelesai))
		harapKode(t, kode, 409, "slot "+slotMulai+"-"+slotSelesai)
		if obj["error"] == nil {
			t.Fatalf("field error harus tetap ada: %v", obj)
		}
		got, ada := obj["bentrok_dengan"].([]any)
		if !ada || len(got) != len(harap) {
			t.Fatalf("slot %s-%s: bentrok_dengan %v, harusnya %v", slotMulai, slotSelesai, obj["bentrok_dengan"], harap)
		}
		for i := range harap {
			if got[i] != harap[i] {
				t.Fatalf("slot %s-%s: bentrok_dengan %v, harusnya %v", slotMulai, slotSelesai, got, harap)
			}
		}
	}
	cek("3", "5", []any{float64(1), float64(2)})
	cek("1", "2", []any{float64(1)})

	// 201 dan 400 tidak memuat bentrok_dengan
	kode, obj, _ := kirim(t, h, "POST", "/peminjaman", pinjamSlot("7", "8"))
	harapKode(t, kode, 201, "sukses")
	if _, ada := obj["bentrok_dengan"]; ada {
		t.Fatalf("201 tidak boleh memuat bentrok_dengan: %v", obj)
	}
	kode, obj, _ = kirim(t, h, "POST", "/peminjaman", pinjamSlot("0", "1"))
	harapKode(t, kode, 400, "validasi")
	if _, ada := obj["bentrok_dengan"]; ada {
		t.Fatalf("400 tidak boleh memuat bentrok_dengan: %v", obj)
	}

	// PUT tidak berubah
	kode, obj, _ = kirim(t, h, "PUT", "/peminjaman/1", pinjamSlot("5", "6"))
	harapKode(t, kode, 409, "put bentrok")
	if _, ada := obj["bentrok_dengan"]; ada {
		t.Fatalf("PUT tidak boleh memuat bentrok_dengan: %v", obj)
	}

	// blokir: boleh berisi []
	kirim(t, h, "POST", "/blokir", `{"sumber_daya_id":1,"tanggal":"2026-12-01","slot_mulai":1,"slot_selesai":2,"keterangan":"Servis"}`)
	kode, obj, _ = kirim(t, h, "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-12-01","slot_mulai":2,"slot_selesai":3,"jumlah":1}`)
	harapKode(t, kode, 409, "kena blokir")
	if got, ada := obj["bentrok_dengan"].([]any); !ada || len(got) != 0 {
		t.Fatalf("blokir: bentrok_dengan harus [], dapat %v", obj["bentrok_dengan"])
	}
}
