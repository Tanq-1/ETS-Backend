# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Pembuatan API dasar (Checkpoint 1: B1-B3)
- Prompt: "Persona saya adalah Persona A. Tolong kerjakan ETS-nya dan berikan jawaban sesuai urutan checkpoint. Berikan chekpoint 1 dan setelah konfirmasi baru checkpoint 2. Sebelum memulai, jelaskan apa isi AI-LOG.md, PROSES.md, dan README-mahasiswa.md."
- Dipakai: `go.mod`; `main.go` untuk menjalankan server pada port 8080; `api.go` yang berisi struct `SumberDaya`, `Peminjaman`, `store`, fungsi `NewHandler()`, helper JSON dan validasi awal, serta handler B1-B3; `api_test.go` untuk pengujian awal B1-B3 dan penyimpanan terpisah setiap `NewHandler()`.
- Diubah: Belum ada perubahan kode secara manual yang dapat dicatat. Kode masih berupa hasil awal AI dan perlu diuji, dipahami, serta diperbaiki jika ada kekeliruan. Fitur B4-B7 dan endpoint ketersediaan Persona A belum dibuat pada Checkpoint 1.

## Perbaikan dokumentasi penggunaan AI
- Prompt: "Berikan perbaikan untuk AI-LOG.md untuk mengikuti format yang sudah saya lampirkan. Kembalikan dulu AI-LOG.md yang benar. Lain kali, jangan isi proses.md agar saya isi sendiri catat juga perubahan ini di AI-LOG.md."
- Dipakai: Penyusunan ulang catatan `AI-LOG.md` menjadi entri dengan tiga bagian: Prompt, Dipakai, dan Diubah, mengikuti contoh pada dokumen Requirement ETS.
- Diubah: Format log yang sebelumnya memuat bagian terpisah seperti Verifikasi diganti menjadi format tiga butir tersebut. `PROSES.md` sebelumnya sempat dibuatkan draf oleh AI pada Checkpoint 1; mulai permintaan ini, isi `PROSES.md` diserahkan sepenuhnya kepada mahasiswa. Pada revisi ini `PROSES.md` tidak diubah, dan tidak ada perubahan pada kode Go.
