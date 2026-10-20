# AI-LOG.md

## Pembuatan API dasar (Checkpoint 1: B1-B3)

- Prompt: "Persona saya adalah Persona A. Tolong kerjakan ETS-nya dan berikan jawaban sesuai urutan checkpoint. Berikan chekpoint 1 dan setelah konfirmasi baru checkpoint 2. Sebelum memulai, jelaskan apa isi AI-LOG.md, PROSES.md, dan README-mahasiswa.md."
- Dipakai: `go.mod`; `main.go` untuk menjalankan server pada port 8080; `api.go` yang berisi struct `SumberDaya`, `Peminjaman`, `store`, fungsi `NewHandler()`, helper JSON dan validasi awal, serta handler B1-B3; `api_test.go` untuk pengujian awal B1-B3 dan penyimpanan terpisah setiap `NewHandler()`.
- Diubah: Mengubah nama file agar lebih sesuai

## Perbaikan dokumentasi penggunaan AI

- Prompt: "Berikan perbaikan untuk AI-LOG.md untuk mengikuti format yang sudah saya lampirkan. Kembalikan dulu AI-LOG.md yang benar. Lain kali, jangan isi proses.md agar saya isi sendiri catat juga perubahan ini di AI-LOG.md."
- Dipakai: Penyusunan ulang catatan `AI-LOG.md` menjadi entri dengan tiga bagian: Prompt, Dipakai, dan Diubah, mengikuti contoh pada dokumen Requirement ETS.
- Diubah: Pada prompt berikutnya, instruksi bahwa PROSES.md tidak bisa diubah akan diubah

## Pembaruan PROSES.md berdasarkan versi mahasiswa

- Prompt: "Maaf kamu boleh mengubah PROSES.md. hanya saja jangan ubah bagian pertama dan ikuti dari sini."
- Dipakai: Dokumen `PROSES.md` versi mahasiswa yang dilampirkan, serta struktur data dan handler yang sudah ditulis pada Checkpoint 1.
- Diubah: -

## Implementasi lanjutan (Checkpoint 2: B4-B7 dan Persona A)

- Prompt: "Sekarang lanjutkan ke checkpoint 2 dan berikan AI-LOG.md, README.md, dan PROSES.md terbaru jika memang berubah."
- Dipakai: `api.go` ditambahkan penanganan 404/405 dalam bentuk JSON, validasi dan filter lebih ketat, perhitungan bentrok berdasarkan total `jumlah` di setiap slot, endpoint `PUT /peminjaman/{id}` dan `DELETE /peminjaman/{id}`, `sync.RWMutex` untuk operasi bersamaan, serta `GET /sumber-daya/{id}/ketersediaan?tanggal=YYYY-MM-DD` khusus Persona A. `checkpoint2_test.go` dibuat untuk menguji kasus validasi, bentrok 409, pembaruan yang gagal tidak mengubah data, penghapusan, ketersediaan delapan slot, dan permintaan peminjaman paralel. `README.md` diperbarui menurut template mahasiswa, serta bagian 2 dan 4 `PROSES.md` dilanjutkan tanpa mengubah bagian 1.
- Diubah: Pada prompt berikutnya, akan ada permintaan untuk mempersatukan kedua file test

## Penyatuan file tes setelah Checkpoint 2

- Prompt: "berikan api_test.go yang lengkap dan berikan juga AI lognya."
- Dipakai: Seluruh pengujian awal B1-B3 dalam `api_test.go` serta pengujian B4-B7 dan fitur Persona A yang sebelumnya terdapat dalam `checkpoint2_test.go`.
- Diubah: -

## Penerapan Change Request B (Daftar Bentrok)

- Prompt: "Terapkan change request ini dan ubah PROSES.md, AI-LOG.md, dan README.md sesuai kebutuhan serta file-file lain." Lampiran: `Lembar CR-B Daftar Bentrok.docx`.
- Dipakai: Aturan CR-B untuk `POST /peminjaman` yang gagal karena kapasitas, struktur map peminjaman yang sudah ada, helper `adaBentrok`, dan satu berkas pengujian `api_test.go`. Pada respons 409 yang relevan ditambahkan `bentrok_dengan` berisi ID peminjaman aktif dengan sumber daya, tanggal, serta irisan slot yang sesuai, diurutkan naik. Pengujian mencakup dua ID bentrok, pembatasan tanggal/sumber daya/slot, dan ketiadaan field pada respons lain.
- Diubah: Bagian akhir README-mahasiswa.md
