# PinjamKampus API

- Nama dan NRP: [isi nama] - [isi NRP]
- Persona: A - Petugas gedung
- Cara menjalankan: `go run .` lalu buka http://localhost:8080

## Fitur
- [x] B1 Sumber daya
- [x] B2 Buat peminjaman
- [x] B3 Lihat peminjaman (daftar, filter, detail)
- [x] B4 Validasi dan error JSON
- [x] B5 Deteksi bentrok (409)
- [x] B6 Ubah dan hapus
- [x] B7 Aman untuk permintaan bersamaan
- [x] Fitur khas persona
- [ ] Change request

## Yang belum selesai
- Change request belum diterima atau dikerjakan.
- B1-B7 dan fitur Persona A sudah diimplementasikan beserta tes otomatis, tetapi tes belum dapat dijalankan di lingkungan pembuatan karena program Go belum tersedia. Jalankan `go test ./...`, `go test -race ./...`, `go build ./...`, dan `go vet ./...` pada komputer yang memiliki Go 1.22 atau lebih baru.
