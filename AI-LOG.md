# AI-LOG

## Pembuatan API dasar (Checkpoint 1: B1-B3)

- Prompt: "Persona saya adalah Persona A. Tolong kerjakan ETS-nya dan berikan jawaban sesuai urutan checkpoint. Berikan chekpoint 1 dan setelah konfirmasi baru checkpoint 2. Sebelum memulai, jelaskan apa isi AI-LOG.md, PROSES.md, dan README-mahasiswa.md."
- Dipakai: `go.mod`; `main.go` untuk menjalankan server pada port 8080; `api.go` yang berisi struct `SumberDaya`, `Peminjaman`, `store`, fungsi `NewHandler()`, helper JSON dan validasi awal, serta handler B1-B3; `api_test.go` untuk pengujian awal B1-B3 dan penyimpanan terpisah setiap `NewHandler()`.
- Diubah: Belum ada perubahan kode secara manual yang dapat dicatat. Kode masih berupa hasil awal AI dan perlu diuji, dipahami, serta diperbaiki jika ada kekeliruan. Fitur B4-B7 dan endpoint ketersediaan Persona A belum dibuat pada Checkpoint 1.

## Perbaikan dokumentasi penggunaan AI

- Prompt: "Berikan perbaikan untuk AI-LOG.md untuk mengikuti format yang sudah saya lampirkan. Kembalikan dulu AI-LOG.md yang benar. Lain kali, jangan isi proses.md agar saya isi sendiri catat juga perubahan ini di AI-LOG.md."
- Dipakai: Penyusunan ulang catatan `AI-LOG.md` menjadi entri dengan tiga bagian: Prompt, Dipakai, dan Diubah, mengikuti contoh pada dokumen Requirement ETS.
- Diubah: Format log yang sebelumnya memuat bagian terpisah seperti Verifikasi diganti menjadi format tiga butir tersebut. `PROSES.md` sebelumnya sempat dibuatkan draf oleh AI pada Checkpoint 1; mulai permintaan ini, isi `PROSES.md` diserahkan sepenuhnya kepada mahasiswa. Pada revisi ini `PROSES.md` tidak diubah, dan tidak ada perubahan pada kode Go.

## Pembaruan PROSES.md berdasarkan versi mahasiswa

- Prompt: "Maaf kamu boleh mengubah PROSES.md. hanya saja jangan ubah bagian pertama dan ikuti dari sini."
- Dipakai: Dokumen `PROSES.md` versi mahasiswa yang dilampirkan, serta struktur data dan handler yang sudah ditulis pada Checkpoint 1.
- Diubah: Instruksi sebelumnya diperbarui: AI boleh membantu mengisi `PROSES.md`, tetapi bagian 1 (Persona dan masalahnya) harus tetap persis seperti tulisan mahasiswa. Bagian 2 (Rancangan struktur data) dan bagian 3 (Satu keputusan desain) diisi berdasarkan kode Checkpoint 1. Bagian 4 (Bug atau kendala) dan bagian 5 (Change request) dibiarkan untuk diisi berdasarkan kejadian nyata. Tidak ada kode Go yang diubah.

## Implementasi lanjutan (Checkpoint 2: B4-B7 dan Persona A)

- Prompt: "Sekarang lanjutkan ke checkpoint 2 dan berikan AI-LOG.md, README.md, dan PROSES.md terbaru jika memang berubah."
- Dipakai: `api.go` ditambahkan penanganan 404/405 dalam bentuk JSON, validasi dan filter lebih ketat, perhitungan bentrok berdasarkan total `jumlah` di setiap slot, endpoint `PUT /peminjaman/{id}` dan `DELETE /peminjaman/{id}`, `sync.RWMutex` untuk operasi bersamaan, serta `GET /sumber-daya/{id}/ketersediaan?tanggal=YYYY-MM-DD` khusus Persona A. `checkpoint2_test.go` dibuat untuk menguji kasus validasi, bentrok 409, pembaruan yang gagal tidak mengubah data, penghapusan, ketersediaan delapan slot, dan permintaan peminjaman paralel. `README.md` diperbarui menurut template mahasiswa, serta bagian 2 dan 4 `PROSES.md` dilanjutkan tanpa mengubah bagian 1.
- Diubah: -
