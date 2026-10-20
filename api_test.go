package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func kirim(t *testing.T, handler http.Handler, method, path, body string, wantCode int) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != wantCode {
		t.Fatalf("%s %s: status %d, ingin %d; body=%s", method, path, response.Code, wantCode, response.Body.String())
	}
	return response
}

func TestCheckpoint1(t *testing.T) {
	h := NewHandler()

	// B1: daftar kosong, buat sumber daya, dan lihat detail.
	kosong := kirim(t, h, "GET", "/sumber-daya", "", 200)
	if strings.TrimSpace(kosong.Body.String()) != "[]" {
		t.Fatalf("daftar awal harus [], dapat %s", kosong.Body.String())
	}
	created := kirim(t, h, "POST", "/sumber-daya", `{"nama":"  Aula A  ","jenis":"ruang","stok":1}`, 201)
	var sumber SumberDaya
	if err := json.Unmarshal(created.Body.Bytes(), &sumber); err != nil || sumber.ID != 1 || sumber.Nama != "Aula A" {
		t.Fatalf("sumber daya tidak sesuai: %+v, %v", sumber, err)
	}
	kirim(t, h, "GET", "/sumber-daya/1", "", 200)
	kirim(t, h, "GET", "/sumber-daya/99", "", 404)
	kirim(t, h, "GET", "/sumber-daya/bukan-angka", "", 400)

	// B2: buat peminjaman dengan status awal aktif.
	baru := kirim(t, h, "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":" Rina ","tanggal":"2026-11-02","slot_mulai":3,"slot_selesai":4,"jumlah":1}`, 201)
	var pinjam Peminjaman
	if err := json.Unmarshal(baru.Body.Bytes(), &pinjam); err != nil || pinjam.ID != 1 || pinjam.Peminjam != "Rina" || pinjam.Status != "aktif" {
		t.Fatalf("peminjaman tidak sesuai: %+v, %v", pinjam, err)
	}

	// B3: daftar, detail, dan filter yang dapat digabung.
	filtered := kirim(t, h, "GET", "/peminjaman?sumber_daya_id=1&tanggal=2026-11-02&status=aktif", "", 200)
	var daftar []Peminjaman
	if err := json.Unmarshal(filtered.Body.Bytes(), &daftar); err != nil || len(daftar) != 1 || daftar[0].ID != 1 {
		t.Fatalf("hasil filter tidak sesuai: %+v, %v", daftar, err)
	}
	kirim(t, h, "GET", "/peminjaman/1", "", 200)
	missing := kirim(t, h, "GET", "/peminjaman?tanggal=2026-11-03", "", 200)
	if strings.TrimSpace(missing.Body.String()) != "[]" {
		t.Fatalf("filter kosong harus [], dapat %s", missing.Body.String())
	}
	kirim(t, h, "GET", "/peminjaman?sumber_daya_id=abc", "", 400)
}

func TestNewHandlerPenyimpananTerpisah(t *testing.T) {
	h1 := NewHandler()
	h2 := NewHandler()
	kirim(t, h1, "POST", "/sumber-daya", `{"nama":"Proyektor","jenis":"alat","stok":3}`, 201)
	response := kirim(t, h2, "GET", "/sumber-daya", "", 200)
	if strings.TrimSpace(response.Body.String()) != "[]" {
		t.Fatalf("handler kedua seharusnya kosong, dapat %s", response.Body.String())
	}
}

const tanggalTes = "2026-11-02"

func buatSumberTes(t *testing.T, h http.Handler, stok int) {
	t.Helper()
	kirim(t, h, "POST", "/sumber-daya", fmt.Sprintf(`{"nama":"Aula","jenis":"ruang","stok":%d}`, stok), 201)
}

func bodyPeminjamanTes(peminjam string, mulai, selesai, jumlah int) string {
	return fmt.Sprintf(`{"sumber_daya_id":1,"peminjam":%q,"tanggal":%q,"slot_mulai":%d,"slot_selesai":%d,"jumlah":%d}`,
		peminjam, tanggalTes, mulai, selesai, jumlah)
}

func TestCheckpoint2B4ValidasiDanErrorJSON(t *testing.T) {
	h := NewHandler()
	buatSumberTes(t, h, 2)

	tests := []struct {
		method, path, body string
		status             int
	}{
		{"GET", "/sumber-daya/tidak-angka", "", 400},
		{"GET", "/sumber-daya/99", "", 404},
		{"GET", "/peminjaman?sumber_daya_id=", "", 400},
		{"POST", "/sumber-daya", `{"nama":"X","jenis":"ruang","stok":0}`, 400},
		{"POST", "/sumber-daya", `{"nama":"X","jenis":"ruang","stok":1,"tambahan":true}`, 400},
		{"POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-02-29","slot_mulai":1,"slot_selesai":2,"jumlah":1}`, 400},
		{"POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":5,"slot_selesai":4,"jumlah":1}`, 400},
		{"GET", "/tidak-ada", "", 404},
		{"PATCH", "/peminjaman/1", "", 405},
	}

	for _, tc := range tests {
		r := kirim(t, h, tc.method, tc.path, tc.body, tc.status)
		if got := r.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
			t.Errorf("%s %s: Content-Type=%q, ingin JSON", tc.method, tc.path, got)
		}
		var message map[string]string
		if err := json.Unmarshal(r.Body.Bytes(), &message); err != nil || message["error"] == "" {
			t.Errorf("%s %s: error tidak berupa JSON: %s", tc.method, tc.path, r.Body.String())
		}
	}
}

func TestCheckpoint2B5BentrokPerSlotDanJumlah(t *testing.T) {
	h := NewHandler()
	buatSumberTes(t, h, 3)

	// Penggunaan slot 2-4 adalah 2 unit dan slot 3-5 ditambah 1 unit.
	kirim(t, h, "POST", "/peminjaman", bodyPeminjamanTes("Rina", 2, 4, 2), 201)
	kirim(t, h, "POST", "/peminjaman", bodyPeminjamanTes("Budi", 3, 5, 1), 201)
	kirim(t, h, "POST", "/peminjaman", bodyPeminjamanTes("Terlalu banyak", 4, 6, 1), 409)

	// Slot 6-7 tidak beririsan dan peminjaman tanggal lain terpisah.
	kirim(t, h, "POST", "/peminjaman", bodyPeminjamanTes("Siti", 6, 7, 3), 201)
	lainTanggal := strings.Replace(bodyPeminjamanTes("Tanggal lain", 2, 4, 3), tanggalTes, "2026-11-03", 1)
	kirim(t, h, "POST", "/peminjaman", lainTanggal, 201)

	list := kirim(t, h, "GET", "/peminjaman", "", 200)
	var entries []Peminjaman
	if err := json.Unmarshal(list.Body.Bytes(), &entries); err != nil || len(entries) != 4 {
		t.Fatalf("peminjaman gagal tidak boleh tersimpan: %+v, err=%v", entries, err)
	}
}

func TestCheckpoint2B6PUTDELETEAtomicDanIDTidakDidaurUlang(t *testing.T) {
	h := NewHandler()
	buatSumberTes(t, h, 1)
	kirim(t, h, "POST", "/peminjaman", bodyPeminjamanTes("Rina", 2, 3, 1), 201)
	kirim(t, h, "POST", "/peminjaman", bodyPeminjamanTes("Budi", 5, 5, 1), 201)

	// Jika PUT bentrok, data lama harus tetap ada.
	kirim(t, h, "PUT", "/peminjaman/1", bodyPeminjamanTes("Rina", 5, 5, 1), 409)
	old := kirim(t, h, "GET", "/peminjaman/1", "", 200)
	var sebelum Peminjaman
	if err := json.Unmarshal(old.Body.Bytes(), &sebelum); err != nil || sebelum.SlotMulai != 2 || sebelum.SlotSelesai != 3 {
		t.Fatalf("PUT gagal mengubah data lama: %+v, err=%v", sebelum, err)
	}

	// Slot 3 yang sebelumnya digunakan diri sendiri tidak dihitung sebagai bentrok.
	updated := kirim(t, h, "PUT", "/peminjaman/1", bodyPeminjamanTes(" Rina ", 3, 4, 1), 200)
	var sesudah Peminjaman
	if err := json.Unmarshal(updated.Body.Bytes(), &sesudah); err != nil || sesudah.ID != 1 || sesudah.SlotMulai != 3 || sesudah.Peminjam != "Rina" || sesudah.Status != "aktif" {
		t.Fatalf("PUT berhasil tidak sesuai: %+v, err=%v", sesudah, err)
	}
	kirim(t, h, "PUT", "/peminjaman/99", "{}", 404)
	kirim(t, h, "PUT", "/peminjaman/buruk", "{}", 400)

	deleted := kirim(t, h, "DELETE", "/peminjaman/1", "", 204)
	if deleted.Body.Len() != 0 {
		t.Fatalf("DELETE 204 harus tanpa body, mendapat %s", deleted.Body.String())
	}
	kirim(t, h, "DELETE", "/peminjaman/1", "", 404)
	kirim(t, h, "GET", "/peminjaman/1", "", 404)
	kirim(t, h, "POST", "/peminjaman", bodyPeminjamanTes("Baru", 3, 4, 1), 201)
	newBooking := kirim(t, h, "GET", "/peminjaman/3", "", 200)
	var newItem Peminjaman
	if err := json.Unmarshal(newBooking.Body.Bytes(), &newItem); err != nil || newItem.ID != 3 {
		t.Fatalf("ID baru tidak boleh mendaur ulang ID 1: %+v, err=%v", newItem, err)
	}
}

type hasilKetersediaanTes struct {
	SumberDayaID int `json:"sumber_daya_id"`
	Tanggal      string `json:"tanggal"`
	Slot         []struct {
		Slot     int  `json:"slot"`
		Sisa     int  `json:"sisa"`
		Tersedia bool `json:"tersedia"`
	} `json:"slot"`
}

func cekKetersediaanTes(t *testing.T, h http.Handler, expected []int) {
	t.Helper()
	r := kirim(t, h, "GET", "/sumber-daya/1/ketersediaan?tanggal="+tanggalTes, "", 200)
	var hasil hasilKetersediaanTes
	if err := json.Unmarshal(r.Body.Bytes(), &hasil); err != nil || hasil.SumberDayaID != 1 || hasil.Tanggal != tanggalTes || len(hasil.Slot) != 8 {
		t.Fatalf("format ketersediaan salah: %+v, err=%v", hasil, err)
	}
	for i, slot := range hasil.Slot {
		if slot.Slot != i+1 || slot.Sisa != expected[i] || slot.Tersedia != (expected[i] > 0) {
			t.Errorf("slot %d = %+v, sisa yang diinginkan %d", i+1, slot, expected[i])
		}
	}
}

func TestCheckpoint2PersonaAKetersediaanSelaluTerkini(t *testing.T) {
	h := NewHandler()
	buatSumberTes(t, h, 3)
	cekKetersediaanTes(t, h, []int{3, 3, 3, 3, 3, 3, 3, 3})
	kirim(t, h, "POST", "/peminjaman", bodyPeminjamanTes("Rina", 2, 4, 2), 201)
	kirim(t, h, "POST", "/peminjaman", bodyPeminjamanTes("Budi", 3, 3, 1), 201)
	cekKetersediaanTes(t, h, []int{3, 1, 0, 1, 3, 3, 3, 3})
	kirim(t, h, "PUT", "/peminjaman/2", bodyPeminjamanTes("Budi", 5, 5, 1), 200)
	cekKetersediaanTes(t, h, []int{3, 1, 1, 1, 2, 3, 3, 3})
	kirim(t, h, "DELETE", "/peminjaman/1", "", 204)
	cekKetersediaanTes(t, h, []int{3, 3, 3, 3, 2, 3, 3, 3})

	kirim(t, h, "GET", "/sumber-daya/xx/ketersediaan?tanggal="+tanggalTes, "", 400)
	kirim(t, h, "GET", "/sumber-daya/99/ketersediaan?tanggal="+tanggalTes, "", 404)
	kirim(t, h, "GET", "/sumber-daya/1/ketersediaan", "", 400)
	kirim(t, h, "GET", "/sumber-daya/1/ketersediaan?tanggal=2026-02-29", "", 400)
}

func TestCheckpoint2B7PermintaanBersamaanTidakMelebihiStok(t *testing.T) {
	h := NewHandler()
	buatSumberTes(t, h, 1)
	const jumlahRequest = 24
	statuses := make(chan int, jumlahRequest)
	var wg sync.WaitGroup

	for i := 0; i < jumlahRequest; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			body := bodyPeminjamanTes(fmt.Sprintf("Mahasiswa %d", index), 2, 3, 1)
			req := httptest.NewRequest("POST", "/peminjaman", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			h.ServeHTTP(res, req)
			statuses <- res.Code
		}(i)
	}
	wg.Wait()
	close(statuses)
	created, conflicts := 0, 0
	for status := range statuses {
		switch status {
		case 201:
			created++
		case 409:
			conflicts++
		default:
			t.Errorf("status paralel tidak diharapkan: %d", status)
		}
	}
	if created != 1 || conflicts != jumlahRequest-1 {
		t.Fatalf("harus ada 1 sukses dan %d bentrok, dapat %d sukses dan %d bentrok", jumlahRequest-1, created, conflicts)
	}
	cekKetersediaanTes(t, h, []int{1, 0, 0, 1, 1, 1, 1, 1})
}
