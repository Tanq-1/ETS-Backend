# PROSES

## 1. Persona dan masalahnya
Pengguna persona saya adalah petugas gedung. Dengan fitur khasnya, petugas gedung dapat menjawab ketersediaan ruangan pada suatu waktu tertentu dengan lebih cepat.

## 2. Rancangan struktur data
Pada Checkpoint 1, saya menggunakan dua struct utama:

- `SumberDaya`: menyimpan `ID`, `Nama`, `Jenis`, dan `Stok` untuk data ruang atau alat.
- `Peminjaman`: menyimpan `ID`, `SumberDayaID`, `Peminjam`, `Tanggal`, `SlotMulai`, `SlotSelesai`, `Jumlah`, dan `Status`.

Kedua jenis data disimpan dalam `store` menggunakan `map[int]SumberDaya` dan `map[int]Peminjaman`, dengan ID sebagai kunci. Saya juga memakai `nextSumber` dan `nextPinjam` untuk menentukan ID baru mulai dari 1. Ketika endpoint daftar dipanggil, isi map dimasukkan ke slice lalu diurutkan berdasarkan ID supaya hasilnya selalu berurutan.

Pada Checkpoint 2, saya menambahkan `sync.RWMutex` di dalam `store`. Operasi yang mengubah map atau nomor ID menggunakan `Lock`, sedangkan operasi membaca menggunakan `RLock`. Bentrok peminjaman dihitung dari jumlah unit yang terpakai pada masing-masing slot, untuk sumber daya dan tanggal yang sama. Endpoint ketersediaan Persona A menghitung sisa stok untuk setiap slot 1 sampai 8 dari data peminjaman yang masih aktif, sehingga hasilnya ikut berubah setelah POST, PUT, dan DELETE.

## 3. Satu keputusan desain
Saya memilih `map` untuk penyimpanan utama. Pilihan lainnya adalah `slice`, tetapi jika memakai slice saya perlu mencari data satu per satu ketika endpoint detail menerima ID tertentu. Dengan map, saya bisa langsung mengambil data berdasarkan ID. Kekurangannya, urutan isi map tidak tetap, sehingga pada endpoint daftar saya mengubahnya menjadi slice dan mengurutkannya terlebih dahulu.

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala: Pada rancangan awal Checkpoint 2 terdapat pembacaan `nextSumber` saat menyusun `SumberDaya` sebelum mutex dikunci. Ini berpotensi menimbulkan *data race* ketika beberapa permintaan membuat sumber daya bersamaan.
- Cara menemukan: Peninjauan kode dengan bantuan AI, dengan memeriksa seluruh akses terhadap map dan penghitung ID setelah `sync.RWMutex` ditambahkan. Masalah ini ditemukan lewat pemeriksaan kode, bukan hasil menjalankan `go test -race`.
- Penyebab dan solusi: Penghitung ID dibaca di luar area yang dilindungi mutex. Inisialisasi objek diubah agar pembacaan dan penetapan `nextSumber` dilakukan setelah `Lock()` dan sebelum `Unlock()`.

### Bug 2
- Gejala: Pada Checkpoint 1, permintaan `GET /peminjaman?sumber_daya_id=` yang nilai filternya kosong diperlakukan seperti tidak ada filter dan menghasilkan respons 200.
- Cara menemukan: Membaca implementasi filter di Checkpoint 1 dan membandingkannya dengan kebutuhan bahwa nilai filter `sumber_daya_id` harus berupa angka.
- Penyebab dan solusi: Kode lama hanya memeriksa apakah nilai hasil `query.Get()` tidak kosong. Kode diganti untuk memeriksa keberadaan kunci query, lalu melakukan `strconv.Atoi()` terhadap nilainya. Karena string kosong tidak bisa dikonversi ke angka, responsnya menjadi 400 JSON.

## 5. Change request
CR-B dari tim frontend meminta agar respons `POST /peminjaman` yang ditolak dengan status 409 karena slot atau kapasitas penuh menambahkan `bentrok_dengan`, yaitu daftar ID peminjaman yang menghalangi. ID harus terurut naik; peminjaman pada tanggal, sumber daya, atau slot yang tidak beririsan tidak ikut. Respons sukses, error lainnya, dan `PUT /peminjaman/{id}` tidak boleh berubah.

Saya mengubah handler `tambahPeminjaman` pada `api.go`. Ketika `adaBentrok` menemukan kapasitas tidak cukup, handler tetap mengembalikan `error` lalu menambahkan `bentrok_dengan` dari helper baru `idPeminjamanBentrok`. Helper tersebut hanya mengambil peminjaman aktif pada sumber daya dan tanggal yang sama dengan irisan slot, kemudian mengurutkan ID menggunakan `sort.Ints`.

Keputusan desain saya adalah **memisahkan pemeriksaan kapasitas dari pengumpulan daftar ID**. Pemeriksaan kapasitas per slot tetap menentukan apakah permintaan harus ditolak; daftar ID baru dihitung jika POST memang menghasilkan 409. Dengan begitu, aturan bentrok B5 tidak berubah dan respons PUT juga tetap seperti sebelumnya. Kedua operasi dilakukan ketika mutex masih terkunci supaya hasil daftar ID sesuai dengan kondisi yang menyebabkan penolakan. Pengujian CR-B ditambahkan ke `api_test.go`, tetapi pelaksanaannya masih menunggu lingkungan yang memiliki Go.
