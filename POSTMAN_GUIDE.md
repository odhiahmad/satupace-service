# 🚀 Quick Start Guide - Postman Setup

## 📦 Files yang Tersedia

1. **API_EXAMPLES.md** - Dokumentasi lengkap dengan contoh curl
2. **Run-Sync-API.postman_collection.json** - Postman Collection (import ke Postman)
3. **Run-Sync-Local.postman_environment.json** - Postman Environment variables

---

## 🎯 Cara Import ke Postman

### Step 1: Import Collection

1. Buka **Postman**
2. Klik tombol **Import** (pojok kiri atas)
3. Pilih **File** tab
4. Drag & drop file **Run-Sync-API.postman_collection.json**
5. Klik **Import**

✅ Collection sudah masuk dengan 8 folders:
- Users
- Run Groups
- Run Activities
- Direct Matching
- Chat Messages
- Media & Photos
- Safety Logs
- WhatsApp Verification

### Step 2: Import Environment

1. Di Postman, klik icon **⚙️ Gear** (pojok kanan atas)
2. Pilih **Import**
3. Drag & drop file **Run-Sync-Local.postman_environment.json**
4. Klik **Import**

### Step 3: Aktifkan Environment

1. Di dropdown environment (pojok kanan atas)
2. Pilih **Run-Sync Local**

✅ Sekarang base_url sudah otomatis: `http://localhost:8080`

---

## 🧪 Testing Flow

### 1. Jalankan Server
```bash
cd C:\Users\odhia\Documents\run-sync
go run main.go
```

Server running di: `http://localhost:8080`

### 2. Test Basic Endpoint

**Get All Users:**
- Buka collection: **Users** → **Get All Users**
- Klik **Send**
- Harusnya dapat response list users

### 3. Create User (Register)

**Request:**
```json
{
  "name": "Test User",
  "email": "test@example.com",
  "phone_number": "081234567890",
  "gender": "male",
  "password": "Password123!"
}
```

**Steps:**
1. **Users** → **Create User**
2. Klik **Send**
3. Copy `user_id` dari response
4. Paste ke Environment variable `user_id`

### 4. Get JWT Token

⚠️ **PENTING**: Untuk endpoint yang butuh auth, kamu perlu JWT token.

**Manual cara:**
1. Pakai JWT service di code untuk generate token
2. Atau buat endpoint `/login` untuk dapat token
3. Copy token ke environment variable `jwt_token`

### 5. Create Run Group (Need Auth)

**Request:**
```json
{
  "name": "Weekend Morning Run",
  "avg_pace": 5.5,
  "preferred_distance": 10,
  "latitude": -6.2088,
  "longitude": 106.8456,
  "scheduled_at": "2026-02-22T06:00:00Z",
  "max_member": 15,
  "is_women_only": false
}
```

**Steps:**
1. Pastikan `jwt_token` sudah diisi di environment
2. **Run Groups** → **Create Run Group**
3. Klik **Send**
4. Copy `group_id` dari response
5. Paste ke environment variable `group_id`

### 6. Join Group

1. **Run Groups** → **Join Run Group**
2. URL otomatis pakai `{{group_id}}`
3. Klik **Send**

### 7. Send Group Message

**Request:**
```json
{
  "group_id": "{{group_id}}",
  "content": "Hello everyone!"
}
```

**Steps:**
1. **Chat Messages** → **Send Group Message**
2. Klik **Send**

---

## 🔑 Environment Variables

Setelah import, isi variables berikut sesuai hasil testing:

| Variable | Cara Dapat | Contoh |
|----------|------------|---------|
| `base_url` | Sudah diset | `http://localhost:8080` |
| `jwt_token` | Dari login/generate | `eyJhbGciOiJIUzI1NiIs...` |
| `user_id` | Dari create user response | `550e8400-e29b-41d4-...` |
| `group_id` | Dari create group response | `660f9511-f3ac-52e5-...` |
| `match_id` | Dari create match response | `770g0622-g4bd-63f6-...` |
| `activity_id` | Dari create activity response | `880h1733-h5ce-74g7-...` |

**Cara update variable:**
1. Klik icon **👁️ Eye** (pojok kanan atas)
2. Klik di value yang mau diubah
3. Paste UUID/token yang baru
4. Klik di luar untuk save

---

## 💡 Pro Tips

### 1. Auto-Extract IDs dengan Tests

Tambahkan script di **Tests** tab request untuk auto-save response IDs:

```javascript
// Di Tests tab request "Create User"
if (pm.response.code === 200 || pm.response.code === 201) {
    var jsonData = pm.response.json();
    if (jsonData.data && jsonData.data.id) {
        pm.environment.set("user_id", jsonData.data.id);
    }
}
```

### 2. Pre-request Script untuk Token

Tambahkan di **Pre-request Script** tab:

```javascript
// Auto set authorization header
if (pm.environment.get("jwt_token")) {
    pm.request.headers.add({
        key: 'Authorization',
        value: 'Bearer ' + pm.environment.get("jwt_token")
    });
}
```

### 3. Organize dengan Folders

Collection sudah diorganize berdasarkan feature:
- ✅ Keep it organized
- ✅ Easy to find endpoints
- ✅ Group related requests

### 4. Test Automation

Buat **Test Suite** untuk automation:
1. Pilih Collection
2. Klik **Run** (kanan atas)
3. Select requests yang mau dirun
4. Klik **Run Run-Sync API**

---

## 🐛 Troubleshooting

### Error: "Could not send request"

**Solusi:**
- Pastikan server running: `go run main.go`
- Check port 8080 tidak dipakai app lain
- Test di browser: `http://localhost:8080/users`

### Error: 401 Unauthorized

**Solusi:**
- JWT token kosong atau expired
- Update `jwt_token` di environment
- Generate token baru

### Error: 404 Not Found

**Solusi:**
- Check URL path benar
- Pastikan UUID valid
- Resource mungkin sudah dihapus

### Error: 400 Bad Request

**Solusi:**
- Check request body format (JSON valid)
- Lihat required fields
- Check data types (string, number, boolean)

### Environment Variables tidak kerja

**Solusi:**
- Pastikan environment **Run-Sync Local** aktif (selected)
- Check variable name pakai `{{variable_name}}`
- Case sensitive!

---

## 📝 Testing Checklist

### Basic Flow
- [ ] Get all users
- [ ] Create user
- [ ] Get user by ID
- [ ] Update user

### Run Groups Flow
- [ ] Create run group (dengan auth)
- [ ] Get all groups
- [ ] Get group details
- [ ] Join group
- [ ] Get group members

### Chat Flow
- [ ] Send group message
- [ ] Get chat history
- [ ] Send direct message (setelah match)

### Activities
- [ ] Create run activity
- [ ] Get user activities

---

## 🎨 Customization

### Buat Environment Baru (Production)

1. Duplicate environment **Run-Sync Local**
2. Rename jadi **Run-Sync Production**
3. Update `base_url` ke production URL:
   ```
   https://api.run-sync.com
   ```

### Tambah Custom Variables

Di environment, klik **Add new variable**:
- `api_key` - untuk API key jika ada
- `timeout` - untuk custom timeout
- `version` - untuk API versioning

---

## 📚 Resources

- **API Documentation**: Baca `API_EXAMPLES.md`
- **Project README**: Baca `README.md`
- **Changelog**: Baca `CHANGELOG.md`

---

## ✅ Success!

Setelah import collection & environment, kamu bisa langsung:
1. ⚡ Test all endpoints dengan satu klik
2. 🔄 Auto-populate IDs dengan environment variables
3. 🎯 Organize requests per feature
4. 🚀 Ready untuk development!

**Happy Testing!** 🎉
