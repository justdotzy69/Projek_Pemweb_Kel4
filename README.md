# To-Do List Gamification API (Backend)

Proyek ini adalah *backend* RESTful API untuk aplikasi **To-Do List Gamification**. Dibangun menggunakan bahasa pemrograman Go, aplikasi ini mengubah pengalaman mencatat tugas harian menjadi seperti bermain *game*. Pengguna bisa mendapatkan *Experience Points* (XP), naik level, dan mendapatkan penghargaan (*Badges*) secara otomatis ketika menyelesaikan tugas.

Proyek ini dikembangkan oleh **Kelompok 4** untuk memenuhi tugas mata kuliah Pemrograman Web.

---

## Teknologi yang Digunakan
- **Bahasa Pemrograman:** [Go (Golang)](https://go.dev/) (Minimal versi 1.20)
- **Web Framework:** [Fiber v2](https://gofiber.io/) (Framework Go yang sangat cepat, terinspirasi dari Express.js)
- **Database ORM:** [GORM](https://gorm.io/)
- **Database System:** MySQL (Bisa menggunakan XAMPP, Laragon, dll)
- **Keamanan:** `crypto/bcrypt` (Hashing Password) & `golang-jwt/jwt/v5` (Autentikasi Token)

---

## Persiapan (Prerequisites)
Sebelum menjalankan proyek ini, pastikan komputer Anda sudah terinstal perangkat lunak berikut:
1. **Go:** [Download & Install Go](https://go.dev/dl/)
2. **Database MySQL:** Gunakan [XAMPP](https://www.apachefriends.org/), [Laragon](https://laragon.org/), atau MySQL Server bawaan.
3. **API Tester:** [Postman](https://www.postman.com/downloads/) atau Thunder Client (Ekstensi VS Code).

---

## Cara Instalasi & Inisialisasi Project

Ikuti langkah-langkah di bawah ini untuk menjalankan *server* di komputer lokal Anda:

### 1. Kloning Repositori
Buka terminal/CMD, lalu kloning repositori ini dan masuk ke dalam foldernya:
```bash
git clone <url-repository-github-anda>
cd Projek_Pemweb_Kel4
```

### 2. Inisialisasi Dependensi (Modul Go)
Karena proyek ini menggunakan *framework* Fiber dan modul lainnya, kita perlu mengunduh semua dependensi yang tercatat. Jalankan perintah:
```bash
go mod tidy
```
*(Perintah ini akan secara otomatis membaca file `go.mod` dan mengunduh seluruh framework/library yang dibutuhkan ke dalam komputer Anda).*

### 3. Persiapan Database MySQL
1. Nyalakan server MySQL Anda (Start MySQL di XAMPP/Laragon).
2. Buka *client* database (phpMyAdmin atau DBeaver).
3. Buat database baru (kosong) bernama `todo_gamification`:
   ```sql
   CREATE DATABASE todo_gamification;
   ```

### 4. Konfigurasi Environment Variables (.env)
Aplikasi membutuhkan file `.env` untuk menyimpan kredensial rahasia.
1. Salin template yang sudah disediakan:
   ```bash
   cp .env.example .env
   ```
2. Buka file `.env` yang baru dibuat di *code editor* Anda.
3. Sesuaikan `DB_USER` dan `DB_PASSWORD` dengan konfigurasi MySQL Anda (Jika pakai XAMPP/Laragon, biasanya `DB_USER=root` dan `DB_PASSWORD=` dikosongkan).

### 5. Jalankan Server!
Eksekusi program utama. GORM akan otomatis melakukan *Auto-Migration* (membuatkan tabel di MySQL) jika belum ada.
```bash
go run main.go
```
*Jika berhasil, akan muncul logo Fiber di terminal dan keterangan bahwa server berjalan di `http://127.0.0.1:3000`.*

---

## Cara Penggunaan & Pengetesan API (Postman)

Karena aplikasi ini dilindungi oleh JWT (JSON Web Token), Anda harus melakukan *Login* terlebih dahulu untuk mengakses fitur-fitur utama.

### Langkah 1: Registrasi & Login (Mendapatkan Kunci)
1. **Daftar Akun Baru:** Buka Postman, buat *request* `POST` ke `http://localhost:3000/api/auth/register` dengan format body JSON (isi `email` dan `password`).
2. **Login:** Buat *request* `POST` ke `http://localhost:3000/api/auth/login` menggunakan akun tadi.
3. **Copy Token:** Jika *login* berhasil, server akan merespons dengan JSON berisi teks panjang bernama `"token"`. Salin teks token tersebut.

### Langkah 2: Mengakses Rute Privat (Gunakan Token)
Untuk mencoba fitur seperti melihat *Dashboard* atau membuat tugas, Anda **wajib** menyisipkan token tadi ke dalam Postman:
1. Di Postman, klik tab **Headers**.
2. Tambahkan kolom baru:
   - **Key:** `Authorization`
   - **Value:** `Bearer <paste_token_anda_disini>` *(Perhatikan ada spasi setelah kata Bearer)*

### Langkah 3: Alur Pengetesan Gamifikasi
1. **Buat Master Data:** Lakukan `POST` ke `/api/categories` dan `/api/badges` untuk menyiapkan data awal.
2. **Buat Tugas Baru:** Lakukan `POST` ke `/api/tasks` dengan tingkat kesulitan (`easy`, `medium`, atau `hard`).
3. **Selesaikan Tugas:** Lakukan `PUT` ke `/api/tasks/{id}/complete` untuk menyelesaikan tugas. Lihat responsnya untuk melihat XP yang Anda dapatkan!
4. **Cek Profil:** Lakukan `GET` ke `/api/dashboard` untuk melihat peningkatan level dan *badge* yang berhasil Anda buka.

---

## Panduan Pengembangan & Menyimpan Perubahan (Git Workflow)

Bagi anggota kelompok atau pengembang lain yang ingin melanjutkan kodingan di proyek ini, ikuti alur standar berikut untuk menyimpan perubahan ke GitHub agar tetap rapi:

### 1. Periksa Status File
Pastikan file rahasia seperti `.env` tidak ikut terbaca (harus sudah masuk ke dalam `.gitignore`).
```bash
git status
```

### 2. Tambahkan Perubahan
Masukkan semua file yang baru saja diedit atau ditambahkan:
```bash
git add .
```

### 3. Buat Commit Message yang Rapi
Berikan deskripsi singkat tentang fitur apa yang baru saja diselesaikan. Contoh:
```bash
git commit -m "feat: menambahkan endpoint fitur edit dan hapus tugas"
```

### 4. Push ke GitHub
Dorong perubahan tersebut ke repositori utama:
```bash
git push origin main
```
*(Catatan: Jika branch default menggunakan nama master, gunakan `git push origin master`)*.

---

## Struktur Folder Utama
```text
Projek_Pemweb_Kel4/
├── controllers/     # Menangani logika bisnis dan respons HTTP (Auth, Task, dll)
├── database/        # Konfigurasi koneksi MySQL dan Auto-Migration GORM
├── helpers/         # Fungsi bantuan untuk kalkulasi XP dan Level
├── middlewares/     # Penjaga rute (JWT Auth Middleware)
├── models/          # Struktur data (Struct) untuk tabel database
├── routes/          # Definisi seluruh endpoint API URL
├── .env.example     # Template konfigurasi environment variables
├── .gitignore       # Daftar file yang dikecualikan dari Git (seperti .env)
├── go.mod           # Daftar dependensi modul Go
└── main.go          # Titik masuk utama (Entry Point) aplikasi
```