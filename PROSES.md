# PROSES

## 1. Persona dan masalahnya
Pengguna persona saya adalah petugas gedung. Dengan fitur khasnya, petugas gedung dapat menjawab ketersediaan ruangan pada suatu waktu tertentu dengan lebih cepat.

## 2. Rancangan struktur data
Pada Checkpoint 1, saya menggunakan dua struct utama:

- `SumberDaya`: menyimpan `ID`, `Nama`, `Jenis`, dan `Stok` untuk data ruang atau alat.
- `Peminjaman`: menyimpan `ID`, `SumberDayaID`, `Peminjam`, `Tanggal`, `SlotMulai`, `SlotSelesai`, `Jumlah`, dan `Status`.

Kedua jenis data disimpan dalam `store` menggunakan `map[int]SumberDaya` dan `map[int]Peminjaman`, dengan ID sebagai kunci. Saya juga memakai `nextSumber` dan `nextPinjam` untuk menentukan ID baru mulai dari 1. Ketika endpoint daftar dipanggil, isi map dimasukkan ke slice lalu diurutkan berdasarkan ID supaya hasilnya selalu berurutan.

## 3. Satu keputusan desain
Saya memilih `map` untuk penyimpanan utama. Pilihan lainnya adalah `slice`, tetapi jika memakai slice saya perlu mencari data satu per satu ketika endpoint detail menerima ID tertentu. Dengan map, saya bisa langsung mengambil data berdasarkan ID. Kekurangannya, urutan isi map tidak tetap, sehingga pada endpoint daftar saya mengubahnya menjadi slice dan mengurutkannya terlebih dahulu.

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
