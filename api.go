package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

type Peminjaman struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
	Status       string `json:"status"`
}

// Pointer dipakai agar field yang tidak dikirim dapat dibedakan dari angka 0.
type sumberDayaInput struct {
	Nama  *string `json:"nama"`
	Jenis *string `json:"jenis"`
	Stok  *int    `json:"stok"`
}

type peminjamanInput struct {
	SumberDayaID *int    `json:"sumber_daya_id"`
	Peminjam     *string `json:"peminjam"`
	Tanggal      *string `json:"tanggal"`
	SlotMulai    *int    `json:"slot_mulai"`
	SlotSelesai  *int    `json:"slot_selesai"`
	Jumlah       *int    `json:"jumlah"`
}

type store struct {
	sumberDaya  map[int]SumberDaya
	peminjaman  map[int]Peminjaman
	nextSumber  int
	nextPinjam  int
}

// NewHandler selalu membuat penyimpanan baru, terpisah dari handler lain.
func NewHandler() http.Handler {
	s := &store{
		sumberDaya: make(map[int]SumberDaya),
		peminjaman: make(map[int]Peminjaman),
		nextSumber: 1,
		nextPinjam: 1,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /sumber-daya", s.tambahSumberDaya)
	mux.HandleFunc("GET /sumber-daya", s.daftarSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", s.detailSumberDaya)

	mux.HandleFunc("POST /peminjaman", s.tambahPeminjaman)
	mux.HandleFunc("GET /peminjaman", s.daftarPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", s.detailPeminjaman)

	return mux
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("body harus berisi tepat satu objek JSON")
	}
	return nil
}

func parseID(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func validTanggal(value string) bool {
	parsed, err := time.Parse("2006-01-02", value)
	return err == nil && parsed.Format("2006-01-02") == value
}

func validTeks(value string) bool {
	return value != "" && utf8.RuneCountInString(value) <= 50
}

func (s *store) tambahSumberDaya(w http.ResponseWriter, r *http.Request) {
	var input sumberDayaInput
	if err := decodeJSON(r, &input); err != nil || input.Nama == nil || input.Jenis == nil || input.Stok == nil {
		writeError(w, http.StatusBadRequest, "body sumber daya tidak valid")
		return
	}

	nama := strings.TrimSpace(*input.Nama)
	jenis := strings.TrimSpace(*input.Jenis)
	if !validTeks(nama) || (jenis != "ruang" && jenis != "alat") || *input.Stok < 1 {
		writeError(w, http.StatusBadRequest, "nama, jenis, atau stok tidak valid")
		return
	}

	item := SumberDaya{ID: s.nextSumber, Nama: nama, Jenis: jenis, Stok: *input.Stok}
	s.sumberDaya[item.ID] = item
	s.nextSumber++
	writeJSON(w, http.StatusCreated, item)
}

func (s *store) daftarSumberDaya(w http.ResponseWriter, _ *http.Request) {
	items := make([]SumberDaya, 0, len(s.sumberDaya))
	for _, item := range s.sumberDaya {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	writeJSON(w, http.StatusOK, items)
}

func (s *store) detailSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}
	item, ok := s.sumberDaya[id]
	if !ok {
		writeError(w, http.StatusNotFound, "sumber daya tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *store) tambahPeminjaman(w http.ResponseWriter, r *http.Request) {
	var input peminjamanInput
	if err := decodeJSON(r, &input); err != nil || input.SumberDayaID == nil ||
		input.Peminjam == nil || input.Tanggal == nil || input.SlotMulai == nil ||
		input.SlotSelesai == nil || input.Jumlah == nil {
		writeError(w, http.StatusBadRequest, "body peminjaman tidak valid")
		return
	}

	peminjam := strings.TrimSpace(*input.Peminjam)
	tanggal := strings.TrimSpace(*input.Tanggal)
	sumber, ada := s.sumberDaya[*input.SumberDayaID]
	if !ada || !validTeks(peminjam) || !validTanggal(tanggal) ||
		*input.SlotMulai < 1 || *input.SlotSelesai > 8 ||
		*input.SlotMulai > *input.SlotSelesai ||
		*input.Jumlah < 1 || *input.Jumlah > sumber.Stok {
		writeError(w, http.StatusBadRequest, "data peminjaman tidak valid")
		return
	}

	// Pemeriksaan bentrok kapasitas (B5) ditambahkan pada Checkpoint 2.
	item := Peminjaman{
		ID: s.nextPinjam, SumberDayaID: *input.SumberDayaID,
		Peminjam: peminjam, Tanggal: tanggal,
		SlotMulai: *input.SlotMulai, SlotSelesai: *input.SlotSelesai,
		Jumlah: *input.Jumlah, Status: "aktif",
	}
	s.peminjaman[item.ID] = item
	s.nextPinjam++
	writeJSON(w, http.StatusCreated, item)
}

func (s *store) daftarPeminjaman(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filterSumber := query.Get("sumber_daya_id")
	var idSumber int
	if filterSumber != "" {
		parsed, err := strconv.Atoi(filterSumber)
		if err != nil {
			writeError(w, http.StatusBadRequest, "filter sumber_daya_id harus angka")
			return
		}
		idSumber = parsed
	}

	items := make([]Peminjaman, 0)
	for _, item := range s.peminjaman {
		if filterSumber != "" && item.SumberDayaID != idSumber {
			continue
		}
		if tanggal := query.Get("tanggal"); tanggal != "" && item.Tanggal != tanggal {
			continue
		}
		if status := query.Get("status"); status != "" && item.Status != status {
			continue
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	writeJSON(w, http.StatusOK, items)
}

func (s *store) detailPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}
	item, ok := s.peminjaman[id]
	if !ok {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, item)
}
