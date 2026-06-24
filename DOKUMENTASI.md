# Dokumentasi Proyek: Todo List Gamification

**Kelompok 4 — Mata Kuliah Pemrograman Web**

---

## Daftar Isi

1. [Gambaran Umum](#gambaran-umum)
2. [Struktur Folder](#struktur-folder)
3. [Alur Kerja Aplikasi](#alur-kerja-aplikasi)
4. [Cara Setup & Menjalankan](#cara-setup--menjalankan)
5. [Referensi API Endpoint](#referensi-api-endpoint)
6. [Penjelasan Per File](#penjelasan-per-file)
7. [Sistem Gamifikasi](#sistem-gamifikasi)
8. [Panduan Git untuk Anggota Kelompok](#panduan-git-untuk-anggota-kelompok)

---

## Gambaran Umum

Aplikasi ini adalah **web app Todo List** yang dilengkapi sistem gamifikasi. Ketika user menyelesaikan tugas, mereka dapat XP, naik level, dan membuka badge secara otomatis.

**Stack yang dipakai:**
- **Go (Golang)** — bahasa pemrograman utama
- **Fiber v2** — web framework, mirip Express.js tapi lebih cepat
- **GORM** — library untuk operasi database (ORM)
- **MySQL** — database
- **JWT** — autentikasi via token untuk route API
- **Cookie-based auth** — autentikasi untuk halaman web HTML

---

## Struktur Folder

```
Projek_Pemweb_Kel4/
│
├── main.go                  # Entry point, setup Fiber + template engine
│
├── routes/
│   └── routes.go            # Semua URL endpoint didaftarkan di sini
│
├── controllers/
│   ├── auth_controller.go   # Register & Login (response JSON untuk API)
│   ├── auth_view.go         # Login, Register, Logout (response halaman HTML)
│   ├── task_controller.go   # CRUD tugas + proses gamifikasi saat complete
│   ├── task_view.go         # Render halaman daftar tugas
│   ├── badge_controller.go  # CRUD badge
│   ├── badge_view.go        # Render halaman koleksi badge
│   ├── category_controller.go # CRUD kategori
│   ├── dashboard_controller.go # Ambil data statistik user (API)
│   └── dashboard_view.go    # Render halaman dashboard
│
├── models/
│   ├── user.go              # Struct tabel users
│   ├── task.go              # Struct tabel tasks
│   ├── badge.go             # Struct tabel badges
│   └── category.go          # Struct tabel categories
│
├── middlewares/
│   ├── auth_middleware.go   # Cek JWT token untuk route /api
│   └── cookie_middleware.go # Cek cookie untuk route halaman web
│
├── helpers/
│   └── gamification_helper.go # Kalkulasi XP dan level
│
├── database/
│   └── connect.go           # Koneksi ke MySQL + auto-migrate tabel
│
├── views/                   # Template HTML (Fiber render engine)
│   ├── layouts/
│   │   ├── base.html        # Layout untuk halaman auth (login/register)
│   │   └── app.html         # Layout untuk halaman setelah login
│   ├── auth/
│   │   ├── login.html
│   │   └── register.html
│   ├── dashboard/
│   │   └── index.html
│   ├── tasks/
│   │   └── index.html
│   └── badges/
│       └── index.html
│
├── static/
│   ├── css/main.css
│   └── js/app.js
│
├── .env                     # Konfigurasi lokal (jangan di-push ke Git!)
├── go.mod                   # Daftar dependensi
└── go.sum                   # Checksum dependensi
```

---

## Alur Kerja Aplikasi

```
Browser buka /
    └── redirect ke /login

/login (GET)  → tampil form login (auth_view.go → LoginPage)
/login (POST) → proses login (auth_view.go → LoginSubmit)
                  ├── cek email di DB
                  ├── bandingkan password dengan bcrypt
                  ├── buat JWT token
                  ├── simpan token ke Cookie
                  └── redirect ke /dashboard

/dashboard, /tasks, /badges
    └── dicek dulu oleh RequireAuth (cookie_middleware.go)
        ├── tidak ada cookie → redirect ke /login
        └── ada cookie → render halaman HTML

/api/* (semua route API)
    └── dicek dulu oleh Protected() (auth_middleware.go)
        ├── tidak ada header Authorization → 401
        └── ada token valid → lanjut ke controller
```

---

## Cara Setup & Menjalankan

### Prasyarat
- Go minimal versi 1.20 → [https://go.dev/dl/](https://go.dev/dl/)
- MySQL (bisa pakai XAMPP atau Laragon)
- Postman atau Thunder Client untuk test API

### Langkah-langkah

**1. Clone repo**
```bash
git clone https://github.com/justdotzy69/Projek_Pemweb_Kel4.git
cd Projek_Pemweb_Kel4
```

**2. Install semua dependensi**
```bash
go mod tidy
```

**3. Buat database di MySQL**
```sql
CREATE DATABASE todo_gamification;
```

**4. Buat file `.env`**

Salin dari contoh yang sudah ada:
```bash
cp .env.example .env
```

Lalu isi sesuai konfigurasi lokal kamu:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=todo_gamification

JWT_SECRET=bebas_isi_string_apapun_yang_susah_ditebak
PORT=3000
```

> Kalau pakai XAMPP/Laragon biasanya `DB_PASSWORD` dikosongkan.

**5. Jalankan server**
```bash
go run main.go
```

Kalau berhasil, terminal akan menampilkan:
```
Koneksi ke database MySQL berhasil!
Migrasi database selesai!
🚀 Server berjalan di http://localhost:3000
```

Tabel di database akan dibuat otomatis oleh GORM.

---

## Referensi API Endpoint

Base URL: `http://localhost:3000`

### Halaman Web

| Method | URL         | Keterangan                          |
|--------|-------------|-------------------------------------|
| GET    | `/`         | Redirect ke `/login`                |
| GET    | `/login`    | Halaman form login                  |
| POST   | `/login`    | Proses login, set cookie            |
| GET    | `/register` | Halaman form registrasi             |
| POST   | `/register` | Proses registrasi                   |
| GET    | `/logout`   | Hapus cookie, redirect ke `/login`  |
| GET    | `/dashboard`| Halaman dashboard *(perlu login)*   |
| GET    | `/tasks`    | Halaman daftar tugas *(perlu login)*|
| GET    | `/badges`   | Halaman badge *(perlu login)*       |

### API JSON

Semua route `/api/*` selain auth membutuhkan header:
```
Authorization: Bearer <token_dari_login>
```

**Auth**

| Method | URL                  | Body (JSON)                        | Keterangan          |
|--------|----------------------|------------------------------------|---------------------|
| POST   | `/api/auth/register` | `{"email":"...", "password":"..."}` | Daftar akun baru   |
| POST   | `/api/auth/login`    | `{"email":"...", "password":"..."}` | Login, dapat token |

**Dashboard**

| Method | URL             | Keterangan                                       |
|--------|-----------------|--------------------------------------------------|
| GET    | `/api/dashboard`| Ambil data user (XP, level, badge yang dimiliki) |

**Tugas**

| Method | URL                       | Body (JSON)                                                                 | Keterangan              |
|--------|---------------------------|-----------------------------------------------------------------------------|-------------------------|
| GET    | `/api/tasks`              | —                                                                           | Ambil semua tugas milik user |
| POST   | `/api/tasks`              | `{"title":"...", "difficulty":"easy/medium/hard", "category_id":1, ...}` | Buat tugas baru        |
| PUT    | `/api/tasks/:id/complete` | —                                                                           | Selesaikan tugas, dapat XP |

**Kategori**

| Method | URL               | Body (JSON)        | Keterangan          |
|--------|-------------------|--------------------|---------------------|
| GET    | `/api/categories` | —                  | Ambil semua kategori|
| POST   | `/api/categories` | `{"name":"Kuliah"}`| Buat kategori baru  |

**Badge**

| Method | URL           | Body (JSON)                                                   | Keterangan        |
|--------|---------------|---------------------------------------------------------------|-------------------|
| GET    | `/api/badges` | —                                                             | Ambil semua badge |
| POST   | `/api/badges` | `{"name":"...", "image_url":"...", "required_level":5}` | Buat badge baru   |

---

## Penjelasan Per File

### `main.go`
Entry point aplikasi. Tugasnya:
1. Load file `.env`
2. Konek ke database
3. Setup Fiber dengan template engine HTML
4. Daftarkan semua route via `routes.SetupRoutes(app)`
5. Jalankan server

Dua fungsi template custom didaftarkan di sini:
- `add` — penjumlahan dua angka, dipakai di template HTML
- `xpOf` — mengembalikan nilai XP berdasarkan difficulty, dipakai di template HTML

---

### `routes/routes.go`
Satu-satunya tempat semua URL didaftarkan. Terbagi dua bagian:
- **Web routes** — URL yang return halaman HTML, pakai `RequireAuth` (cookie)
- **API routes** — URL yang return JSON, pakai `Protected()` (JWT Bearer)

Kalau mau tambah endpoint baru, cukup tambahkan di file ini.

---

### `controllers/auth_controller.go`
Menangani Register dan Login untuk kebutuhan API (return JSON).

- `Register` — validasi input → hash password dengan bcrypt → simpan ke DB
- `Login` — cari user → bandingkan password → buat JWT token → return token

---

### `controllers/auth_view.go`
Menangani halaman web untuk auth (return HTML).

- `LoginPage` / `RegisterPage` — render form HTML
- `LoginSubmit` — sama seperti Login API tapi setelah sukses menyimpan token ke **cookie** lalu redirect ke `/dashboard`
- `RegisterSubmit` — simpan user baru lalu redirect ke `/login`
- `Logout` — hapus cookie lalu redirect ke `/login`

---

### `controllers/task_controller.go`
Yang paling kompleks. Fungsi utama:

- `CreateTask` — buat tugas baru, status awal selalu `pending`
- `GetTasks` — ambil semua tugas user + data kategorinya (pakai `Preload`)
- `CompleteTask` — proses selesaikan tugas menggunakan **database transaction**:
  1. Cek tugas ada dan milik user yang benar
  2. Cek belum pernah diselesaikan sebelumnya
  3. Ubah status jadi `completed`
  4. Tambah XP ke user
  5. Hitung ulang level
  6. Kalau naik level, cek badge yang syaratnya sudah terpenuhi
  7. Berikan badge yang belum dimiliki user
  8. Simpan semua perubahan

Transaction dipakai agar kalau ada error di tengah proses, semua perubahan dibatalkan (tidak ada data setengah tersimpan).

---

### `controllers/badge_controller.go`
- `CreateBadge` — buat data master badge baru (nama, url gambar, syarat level)
- `GetBadges` — ambil semua badge yang tersedia di sistem

---

### `controllers/category_controller.go`
- `CreateCategory` — buat kategori baru
- `GetCategories` — ambil semua kategori

---

### `controllers/dashboard_controller.go`
- `GetDashboard` — ambil data user (XP, level, badge) tanpa menyertakan password

---

### `models/`

| File          | Tabel di DB    | Relasi                                              |
|---------------|----------------|-----------------------------------------------------|
| `user.go`     | `users`        | HasMany Tasks, ManyToMany Badges (via `user_badges`)|
| `task.go`     | `tasks`        | BelongsTo User, BelongsTo Category (opsional)       |
| `badge.go`    | `badges`       | ManyToMany Users                                    |
| `category.go` | `categories`   | HasMany Tasks                                       |

Field `json:"-"` pada password di `user.go` memastikan password tidak pernah ikut keluar di response API.

---

### `middlewares/auth_middleware.go`
`Protected()` — dipakai untuk route `/api/*`.

Alurnya: baca header `Authorization` → pastikan format `Bearer <token>` → parse dan validasi JWT dengan `JWT_SECRET` → kalau valid, simpan `user_id` ke `c.Locals` → lanjut ke controller.

---

### `middlewares/cookie_middleware.go`
Dipakai untuk route halaman web.

- `RequireAuth` — cek apakah ada cookie `token`, kalau tidak ada redirect ke `/login`
- `RedirectIfAuth` — kalau sudah login dan buka `/login` atau `/register`, langsung redirect ke `/dashboard`
- `SetTokenCookie` / `ClearTokenCookie` — helper untuk set dan hapus cookie

---

### `helpers/gamification_helper.go`
Dua fungsi kalkulasi yang dipisah agar mudah diubah:

- `CalculateXP(difficulty string) int`
  - `easy` → 10 XP
  - `medium` → 20 XP
  - `hard` → 30 XP

- `CalculateLevel(totalXP int) int`
  - Formula: `floor(totalXP / 100) + 1`
  - Contoh: 0–99 XP = Level 1, 100–199 XP = Level 2, dst.

Kalau mau ubah sistem XP atau rumus level, cukup edit file ini saja.

---

### `database/connect.go`
- Baca konfigurasi DB dari `.env`
- Buka koneksi ke MySQL via GORM
- Jalankan `AutoMigrate` — GORM akan otomatis buat/update tabel sesuai struct di `models/`

---

## Sistem Gamifikasi

```
User selesaikan tugas (PUT /api/tasks/:id/complete)
        │
        ▼
Dapat XP berdasarkan difficulty:
  easy   → +10 XP
  medium → +20 XP
  hard   → +30 XP
        │
        ▼
Hitung level baru: floor(totalXP / 100) + 1
        │
        ├── Level tidak naik → selesai
        │
        └── Level naik → cek badge
                │
                ▼
        Cari badge dengan required_level <= level baru
                │
                ▼
        Berikan badge yang belum dimiliki user
```

**Contoh progres:**

| Total XP | Level |
|----------|-------|
| 0 – 99   | 1     |
| 100 – 199| 2     |
| 200 – 299| 3     |
| ...      | ...   |

---

## Panduan Git untuk Anggota Kelompok

**Sebelum mulai coding, selalu pull dulu:**
```bash
git pull origin main
```

**Setelah selesai coding:**
```bash
# Cek file apa saja yang berubah
git status

# Masukkan semua perubahan
git add .

# Commit dengan pesan yang jelas
git commit -m "feat: tambah endpoint edit tugas"

# Push ke GitHub
git push origin main
```

**Format pesan commit yang disarankan:**
- `feat:` — fitur baru
- `fix:` — perbaikan bug
- `docs:` — perubahan dokumentasi
- `refactor:` — perubahan kode tanpa mengubah fungsi

**File yang tidak boleh di-push (sudah ada di .gitignore):**
- `.env` — berisi password dan secret key
