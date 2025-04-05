# 📋 Activity API Service

Sebuah RESTful API sederhana untuk manajemen **activity** menggunakan Go dengan Fiber framework, PostgreSQL sebagai database, dan validasi input menggunakan `go-playground/validator`.

## 🚀 Fitur

- ✅ Get semua aktivitas
- ➕ Tambah aktivitas baru
- ✏️ Edit/update aktivitas
- ❌ Hapus aktivitas berdasarkan ID

## 🧱 Struktur Data

```go
type Activity struct {
    ID           uuid.UUID `json:"id"`
    Title        string    `json:"title" validate:"required"`
    Category     string    `json:"category" validate:"required,oneof=TASK EVENT"`
    Description  string    `json:"description" validate:"required"`
    ActivityDate time.Time `json:"activityDate" validate:"required"`
    Status       string    `json:"status" validate:"required,oneof=NEW 'ON PROGRESS' EXPIRED"`
    CreatedAt    time.Time `json:"createdAt"`
}
```

## 💠 Cara Setup

1. **Clone repo**

   ```bash
   git clone https://github.com/username/activity-api.git
   cd activity-api
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Atur koneksi database**

   Di fungsi `initDB()` ganti DNS dengan konfigurasi database PostgreSQL kamu:

   ```go
   dns := "postgresql://<username>:<password>@<host>:<port>/<database>"
   ```

4. **Buat tabel `activities` di PostgreSQL**

   ```sql
   CREATE EXTENSION IF NOT EXISTS "pgcrypto";

   CREATE TABLE activities (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     title TEXT NOT NULL,
     category TEXT NOT NULL CHECK (category IN ('TASK', 'EVENT')),
     description TEXT NOT NULL,
     activity_date DATE NOT NULL,
     status TEXT NOT NULL CHECK (status IN ('NEW', 'ON PROGRESS', 'EXPIRED'))
   );
   ```

5. **Jalankan aplikasi**

   ```bash
   go run main.go
   ```

   Server akan berjalan di: `http://localhost:3000`

## 📬 API Endpoints

### 🔹 GET `/activities`

Ambil semua data aktivitas.

**Response:**

```json
[
  {
    "id": "uuid",
    "title": "Meeting",
    "category": "EVENT",
    "description": "Weekly sync",
    "activityDate": "2025-04-05T00:00:00Z",
    "status": "NEW",
    "createdAt": "2025-04-01T12:00:00Z"
  }
]
```

---

### 🔹 POST `/activities`

Tambah aktivitas baru.

**Request Body:**

```json
{
  "title": "Submit Report",
  "category": "TASK",
  "description": "Submit monthly report",
  "activityDate": "2025-04-10T00:00:00Z",
  "status": "ON PROGRESS"
}
```

**Response:**

```json
{
  "message": "Success",
  "id": "uuid"
}
```

---

### 🔹 PUT `/activities/:id`

Update aktivitas berdasarkan ID.

**Request Body:** (sama seperti POST)

**Response:**

```json
{
  "message": "Success"
}
```

---

### 🔹 DELETE `/activities/:id`

Hapus aktivitas berdasarkan ID.

**Response:**

```json
{
  "message": "Success"
}
```

## 📦 Dependencies

- [Fiber](https://github.com/gofiber/fiber) – Web framework untuk Go
- [PostgreSQL](https://www.postgresql.org/)
- [uuid](https://github.com/gofrs/uuid)
- [go-playground/validator](https://github.com/go-playground/validator)

---

🙌 Selamat ngoding!

