# PinjamKampus API

- Nama dan NRP: Titan Avery Hamonangan Manurung - 5053241005
- Persona: A - Petugas gedung
- Cara menjalankan: `go run .` lalu buka http://localhost:8080

## Fitur
- [x] B1 Sumber daya
- [x] B2 Buat peminjaman
- [x] B3 Lihat peminjaman (daftar, filter, detail)
- [ ] B4 Validasi dan error JSON
- [ ] B5 Deteksi bentrok (409)
- [ ] B6 Ubah dan hapus
- [ ] B7 Aman untuk permintaan bersamaan
- [ ] Fitur khas persona
- [ ] Change request

## Yang belum selesai
- B4 belum lengkap; baru ada validasi dasar untuk endpoint B1-B3.
- B5 belum memeriksa bentrok atau kapasitas slot.
- B6 belum menyediakan PUT dan DELETE peminjaman.
- B7 belum memakai mutex untuk mengamankan data bersama.
- Fitur Persona A `GET /sumber-daya/{id}/ketersediaan` belum dibuat.
- Change request belum diterima atau dikerjakan.
- B1-B3 sudah ditulis, tetapi masih perlu diverifikasi dengan `go test ./...`, `go build ./...`, dan `go vet ./...` pada komputer yang memiliki Go.
