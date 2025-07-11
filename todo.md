# Daftar Tugas Proyek In-Memory Storage

## Fitur Inti

- [x] Implementasi struktur data dasar:
  - [x] String
  - [x] List
  - [x] Hash
  - [ ] Set
  - [ ] Sorted Set
- [x] Implementasi perintah (command) untuk setiap struktur data:
  - [x] Perintah String (SET, GET, DEL)
  - [x] Perintah String (INCR, DECR)
  - [x] Perintah List (LPUSH, RPUSH, LPOP, RPOP)
  - [x] Perintah List (LLEN)
  - [x] Perintah Hash (HSET, HGET, HDEL, HGETALL)
  - [x] Perintah Set (SADD, SREM, SMEMBERS, SISMEMBER)
  - [x] Perintah Sorted Set (ZADD, ZRANGE, ZREM)
- [x] Penanganan koneksi klien secara konkuren.
- [x] Parser untuk protokol komunikasi (misalnya, RESP - REdis Serialization Protocol).

## Persistensi Data

- [ ] Snapshotting (seperti RDB di Redis) untuk menyimpan state ke disk.
- [ ] Append-Only File (AOF) untuk mencatat setiap operasi tulis.

## Manajemen Memori

- [ ] Kebijakan penggusuran (Eviction Policy) saat memori penuh:
  - [ ] LRU (Least Recently Used)
  - [ ] LFU (Least Frequently Used)
- [x] Dukungan TTL (Time To Live) untuk kunci (key) agar bisa kedaluwarsa.

## Replikasi

- [ ] Implementasi replikasi master-slave.

## Pengujian

- [ ] Unit test untuk semua perintah dan struktur data.
- [ ] Integration test untuk simulasi interaksi klien dan server.
- [ ] Benchmark test untuk mengukur performa.

## Dokumentasi

- [ ] Dokumentasi API untuk setiap perintah.
- [ ] Memperbarui `README.md` dengan instruksi cara menjalankan dan menggunakan proyek.
- [ ] Contoh penggunaan.

## Tambahan

- [ ] Implementasi Pub/Sub (Publish/Subscribe).
- [ ] Transaksi (MULTI/EXEC).
