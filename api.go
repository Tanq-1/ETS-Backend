package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
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
	mu           sync.RWMutex
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
	mux.HandleFunc("GET /sumber-daya/{id}/ketersediaan", s.ketersediaan)

	mux.HandleFunc("POST /peminjaman", s.tambahPeminjaman)
	mux.HandleFunc("GET /peminjaman", s.daftarPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", s.detailPeminjaman)
	mux.HandleFunc("PUT /peminjaman/{id}", s.ubahPeminjaman)
	mux.HandleFunc("DELETE /peminjaman/{id}", s.hapusPeminjaman)

	// Rute yang tidak dikenal serta method yang tidak didukung juga memakai error JSON.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pattern := mux.Handler(r)
		if pattern != "" {
			mux.ServeHTTP(w, r)
			return
		}

		allowed := make([]string, 0, 4)
		for _, method := range []string{"GET", "POST", "PUT", "DELETE"} {
			alternative := r.Clone(r.Context())
			alternative.Method = method
			if _, alternativePattern := mux.Handler(alternative); alternativePattern != "" {
				allowed = append(allowed, method)
			}
		}
		if len(allowed) > 0 {
			w.Header().Set("Allow", strings.Join(allowed, ", "))
			writeError(w, http.StatusMethodNotAllowed, "method tidak diizinkan")
			return
		}
		writeError(w, http.StatusNotFound, "endpoint tidak ditemukan")
	})
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

	item := SumberDaya{Nama: nama, Jenis: jenis, Stok: *input.Stok}
	s.mu.Lock()
	item.ID = s.nextSumber
	s.sumberDaya[item.ID] = item
	s.nextSumber++
	s.mu.Unlock()
	writeJSON(w, http.StatusCreated, item)
}

func (s *store) daftarSumberDaya(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
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
	s.mu.RLock()
	item, ok := s.sumberDaya[id]
	s.mu.RUnlock()
	if !ok {
		writeError(w, http.StatusNotFound, "sumber daya tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *store) tambahPeminjaman(w http.ResponseWriter, r *http.Request) {
	var input peminjamanInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "body peminjaman tidak valid")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.validasiPeminjaman(input)
	if !ok {
		writeError(w, http.StatusBadRequest, "data peminjaman tidak valid")
		return
	}
	if s.adaBentrok(item, 0) {
		writeJSON(w, http.StatusConflict, struct {
			Error         string `json:"error"`
			BentrokDengan []int  `json:"bentrok_dengan"`
		}{
			Error:         "kapasitas slot tidak mencukupi",
			BentrokDengan: s.idPeminjamanBentrok(item),
		})
		return
	}
	item.ID = s.nextPinjam
	item.Status = "aktif"
	s.peminjaman[item.ID] = item
	s.nextPinjam++
	writeJSON(w, http.StatusCreated, item)
}

func (s *store) daftarPeminjaman(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filterSumber := strings.TrimSpace(query.Get("sumber_daya_id"))
	var idSumber int
	_, adaFilterSumber := query["sumber_daya_id"]
	if adaFilterSumber {
		parsed, err := strconv.Atoi(filterSumber)
		if err != nil {
			writeError(w, http.StatusBadRequest, "filter sumber_daya_id harus angka")
			return
		}
		idSumber = parsed
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]Peminjaman, 0)
	for _, item := range s.peminjaman {
		if adaFilterSumber && item.SumberDayaID != idSumber {
			continue
		}
		if tanggal := strings.TrimSpace(query.Get("tanggal")); tanggal != "" && item.Tanggal != tanggal {
			continue
		}
		if status := strings.TrimSpace(query.Get("status")); status != "" && item.Status != status {
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
	s.mu.RLock()
	item, ok := s.peminjaman[id]
	s.mu.RUnlock()
	if !ok {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// validasiPeminjaman dipanggil saat mutex store sedang dikunci.
func (s *store) validasiPeminjaman(input peminjamanInput) (Peminjaman, bool) {
	if input.SumberDayaID == nil || input.Peminjam == nil || input.Tanggal == nil ||
		input.SlotMulai == nil || input.SlotSelesai == nil || input.Jumlah == nil {
		return Peminjaman{}, false
	}
	peminjam := strings.TrimSpace(*input.Peminjam)
	tanggal := strings.TrimSpace(*input.Tanggal)
	sumber, ada := s.sumberDaya[*input.SumberDayaID]
	if !ada || !validTeks(peminjam) || !validTanggal(tanggal) ||
		*input.SlotMulai < 1 || *input.SlotMulai > 8 ||
		*input.SlotSelesai < 1 || *input.SlotSelesai > 8 ||
		*input.SlotMulai > *input.SlotSelesai ||
		*input.Jumlah < 1 || *input.Jumlah > sumber.Stok {
		return Peminjaman{}, false
	}
	return Peminjaman{
		SumberDayaID: *input.SumberDayaID, Peminjam: peminjam,
		Tanggal: tanggal, SlotMulai: *input.SlotMulai,
		SlotSelesai: *input.SlotSelesai, Jumlah: *input.Jumlah,
	}, true
}

// adaBentrok menghitung penggunaan unit pada SETIAP slot, bukan hanya irisan rentang.
// exceptID membuat PUT tidak menghitung peminjaman yang sedang diubah.
func (s *store) adaBentrok(calon Peminjaman, exceptID int) bool {
	stok := s.sumberDaya[calon.SumberDayaID].Stok
	for slot := calon.SlotMulai; slot <= calon.SlotSelesai; slot++ {
		terpakai := 0
		for _, lama := range s.peminjaman {
			if lama.ID == exceptID || lama.Status != "aktif" ||
				lama.SumberDayaID != calon.SumberDayaID || lama.Tanggal != calon.Tanggal {
				continue
			}
			if lama.SlotMulai <= slot && slot <= lama.SlotSelesai {
				terpakai += lama.Jumlah
			}
		}
		if terpakai+calon.Jumlah > stok {
			return true
		}
	}
	return false
}

// idPeminjamanBentrok hanya dipanggil setelah POST terbukti gagal karena kapasitas.
// Pemanggil sudah memegang mutex agar daftar ID berasal dari data yang sama.
func (s *store) idPeminjamanBentrok(calon Peminjaman) []int {
	ids := make([]int, 0)
	for _, lama := range s.peminjaman {
		if lama.Status != "aktif" || lama.SumberDayaID != calon.SumberDayaID ||
			lama.Tanggal != calon.Tanggal {
			continue
		}
		if lama.SlotMulai <= calon.SlotSelesai && calon.SlotMulai <= lama.SlotSelesai {
			ids = append(ids, lama.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

func (s *store) ubahPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}
	s.mu.RLock()
	_, ada := s.peminjaman[id]
	s.mu.RUnlock()
	if !ada {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	var input peminjamanInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "body peminjaman tidak valid")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	lama, ada := s.peminjaman[id]
	if !ada {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	item, valid := s.validasiPeminjaman(input)
	if !valid {
		writeError(w, http.StatusBadRequest, "data peminjaman tidak valid")
		return
	}
	item.ID, item.Status = id, lama.Status
	if s.adaBentrok(item, id) {
		writeError(w, http.StatusConflict, "kapasitas slot tidak mencukupi")
		return
	}
	s.peminjaman[id] = item
	writeJSON(w, http.StatusOK, item)
}

func (s *store) hapusPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ada := s.peminjaman[id]; !ada {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	delete(s.peminjaman, id)
	w.WriteHeader(http.StatusNoContent)
}

type slotTersedia struct {
	Slot     int  `json:"slot"`
	Sisa     int  `json:"sisa"`
	Tersedia bool `json:"tersedia"`
}

func (s *store) ketersediaan(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	sumber, ada := s.sumberDaya[id]
	if !ada {
		writeError(w, http.StatusNotFound, "sumber daya tidak ditemukan")
		return
	}
	tanggal := strings.TrimSpace(r.URL.Query().Get("tanggal"))
	if !validTanggal(tanggal) {
		writeError(w, http.StatusBadRequest, "tanggal tidak valid")
		return
	}
	slots := make([]slotTersedia, 0, 8)
	for slot := 1; slot <= 8; slot++ {
		sisa := sumber.Stok
		for _, pinjam := range s.peminjaman {
			if pinjam.Status == "aktif" && pinjam.SumberDayaID == id &&
				pinjam.Tanggal == tanggal && pinjam.SlotMulai <= slot && slot <= pinjam.SlotSelesai {
				sisa -= pinjam.Jumlah
			}
		}
		slots = append(slots, slotTersedia{Slot: slot, Sisa: sisa, Tersedia: sisa > 0})
	}
	writeJSON(w, http.StatusOK, struct {
		SumberDayaID int            `json:"sumber_daya_id"`
		Tanggal      string         `json:"tanggal"`
		Slot         []slotTersedia `json:"slot"`
	}{SumberDayaID: id, Tanggal: tanggal, Slot: slots})
}
