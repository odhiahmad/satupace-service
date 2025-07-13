# kasirku-service

Starter service untuk aplikasi Kasirku berbasis Golang menggunakan framework Gin.

## Fitur Utama

- **Manajemen Bisnis & User**: Registrasi bisnis, manajemen user bisnis, otentikasi, dan verifikasi email/nomor telepon.
- **Manajemen Produk**: CRUD produk, kategori, bundle, varian, dan unit.
- **Manajemen Transaksi**: Proses transaksi kasir, item, dan pengelolaan pembayaran.
- **Manajemen Lokasi**: Mendukung pencarian provinsi, kota, kecamatan, dan desa.
- **Membership & Keanggotaan**: Pengelolaan membership pelanggan.
- **Metode Pembayaran**: Pengelolaan metode pembayaran.
- **Pagination & Filtering** pada seluruh resource.

## Teknologi

- **Go (Golang)**
- **Gin Web Framework**
- **GORM (ORM)**
- **PostgreSQL** (database utama)
- **Redis** (caching)
- **Validator** (validasi input)
- **JWT** (otentikasi)
- **Postman Collection** (testing API)

## Struktur Folder

```
.
├── .env                   # Konfigurasi environment
├── config/                # Konfigurasi aplikasi & database
├── controller/            # Handler & controller HTTP API
├── data/                  # DTO/request/response
├── entity/                # Model/entitas database
├── helper/                # Helper/utilitas
├── middleware/            # Middleware Gin (auth, logging, dsb)
├── repository/            # Layer akses data (CRUD)
├── routes/                # Routing API
├── service/               # Bisnis logic
├── main.go                # Entry point aplikasi
└── readme.md              # Dokumentasi singkat
```

## Setup & Menjalankan

1. **Buat database** `kasirku` di PostgreSQL.
2. **Impor Struktur Database** sesuai file DB yang tersedia.
3. **Konfigurasi** file `.env` (lihat template `.env`).
4. **Jalankan aplikasi:**

   ```bash
   go run main.go
   ```

5. **Testing API:** Gunakan koleksi Postman yang disediakan untuk mencoba endpoint API.
