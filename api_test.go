package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
