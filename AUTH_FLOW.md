# 🔐 Authentication Flow - Run-Sync

## Overview
Run-Sync menggunakan sistem autentikasi berbasis OTP (One-Time Password) dengan JWT token untuk memproteksi endpoints. User harus **terverifikasi dan aktif** sebelum bisa login dan mengakses fitur-fitur protected.

## Authentication Flow Diagram

```
┌─────────────┐
│   REGISTER  │ → User creates account
└──────┬──────┘
       │
       ↓ (OTP sent via SMS/Email)
┌─────────────┐
│ VERIFY OTP  │ → Account activated & verified
└──────┬──────┘
       │
       ↓ (JWT Token returned)
┌─────────────┐
│   LOGIN     │ → Future logins (verified users only)
└──────┬──────┘
       │
       ↓ (JWT Token)
┌─────────────┐
│  USE APIs   │ → Access protected endpoints
└─────────────┘
```

---

## Step-by-Step Guide

### 1️⃣ **Register (Sign Up)**

**Endpoint:** `POST /auth/register`

User membuat akun baru dengan data:
- Name
- Phone number (unique)
- Email (optional, unique)
- Gender
- Password (min 6 characters)

**Request:**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone_number": "081234567890",
    "gender": "male",
    "password": "SecurePass123!"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "User berhasil dibuat. Silakan verifikasi dengan kode OTP.",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "phone_number": "081234567890",
      "is_verified": false,  ← Not verified yet!
      "is_active": false     ← Not active yet!
    },
    "otp": "123456",
    "message": "Kode OTP telah dikirim ke nomor telepon Anda"
  }
}
```

**What happens:**
- User account dibuat dengan status `is_verified = false` dan `is_active = false`
- OTP 6-digit di-generate dan disimpan di Redis (expire 15 menit)
- OTP dikirim via SMS/Email (dalam development, OTP dikembalikan di response)
- Rate limit: Max 5 OTP requests per 15 menit per phone number

---

### 2️⃣ **Verify OTP**

**Endpoint:** `POST /auth/verify`

User memasukkan kode OTP yang diterima untuk mengaktifkan akun.

**Request:**
```bash
curl -X POST http://localhost:8080/auth/verify \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "081234567890",
    "otp_code": "123456"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Akun berhasil diverifikasi dan diaktifkan",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "phone_number": "081234567890",
      "is_verified": true,   ← Now verified!
      "is_active": true      ← Now active!
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5AZXhhbXBsZS5jb20iLCJleHAiOjE3NDA0MzY4MDAsImlzX2FjdGl2ZSI6dHJ1ZSwiaXNfdmVyaWZpZWQiOnRydWUsImlzcyI6InJ1bi1zeW5jIiwicGhvbmVfbnVtYmVyIjoiMDgxMjM0NTY3ODkwIiwidXNlcl9pZCI6IjU1MGU4NDAwLWUyOWItNDFkNC1hNzE2LTQ0NjY1NTQ0MDAwMCJ9.xyz..."
  }
}
```

**What happens:**
- OTP divalidasi dengan yang tersimpan di Redis
- User status diupdate: `is_verified = true`, `is_active = true`
- OTP dihapus dari Redis
- JWT token di-generate dan dikembalikan
- Token berlaku 24 jam

**Error Cases:**
- OTP salah: `"Kode OTP salah"`
- OTP expired: `"Kode OTP tidak valid atau sudah kadaluarsa"`
- User not found: `"User tidak ditemukan"`

---

### 3️⃣ **Login** (For Returning Users)

**Endpoint:** `POST /auth/login`

Untuk user yang sudah pernah verify, bisa langsung login dengan phone + password.

**Request:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "081234567890",
    "password": "SecurePass123!"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "phone_number": "081234567890",
      "is_verified": true,
      "is_active": true
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Requirements:**
- ✅ User must be **verified** (`is_verified = true`)
- ✅ User must be **active** (`is_active = true`)
- ✅ Password must be correct

**Error Cases:**
- Wrong credentials: `"nomor telepon atau password salah"`
- Not verified: `"Akun belum diverifikasi"` (HTTP 403)
- Inactive account: `"Akun tidak aktif"` (HTTP 403)

---

### 4️⃣ **Resend OTP** (Optional)

**Endpoint:** `POST /auth/resend-otp`

Jika OTP tidak diterima atau expired, user bisa request OTP baru.

**Request:**
```bash
curl -X POST http://localhost:8080/auth/resend-otp \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "081234567890"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Kode OTP baru telah dikirim",
  "data": {
    "otp": "654321"
  }
}
```

**Requirements:**
- User must exist and **not yet verified**
- Rate limit: Max 5 requests per 15 menit

---

### 5️⃣ **Using JWT Token in Protected Endpoints**

Setelah mendapatkan JWT token (dari verify atau login), include token di header `Authorization`:

**Format:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Example:**
```bash
curl -X POST http://localhost:8080/runs/groups \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "name": "Weekend Morning Run",
    "avg_pace": 5.5,
    "preferred_distance": 10,
    "latitude": -6.2088,
    "longitude": 106.8456,
    "scheduled_at": "2026-02-22T06:00:00Z",
    "max_member": 15,
    "is_women_only": false
  }'
```

**JWT Middleware Validation:**
1. Check if token format is valid (Bearer scheme)
2. Validate token signature and expiry
3. Extract user_id, phone_number, email from claims
4. **Check `is_verified = true`** (else 403 Forbidden)
5. **Check `is_active = true`** (else 403 Forbidden)
6. Set user context in request

**Protected Endpoints:**
All endpoints with `middleware.AuthorizeJWT()` require valid token:
- `POST /runs/groups` - Create run group
- `PUT /runs/groups/:id` - Update run group
- `DELETE /runs/groups/:id` - Delete run group
- `POST /runs/groups/:id/join` - Join group
- `POST /runs/activities` - Create activity
- `POST /match` - Create match
- `PATCH /match/:id` - Update match
- `GET /match/me` - Get my matches
- `POST /chats/direct` - Send direct message
- `POST /chats/group` - Send group message
- `POST /media/photos` - Upload photo
- `POST /media/safety` - Create safety log

---

## JWT Token Structure

**Claims:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "phone_number": "081234567890",
  "email": "john@example.com",
  "is_verified": true,
  "is_active": true,
  "exp": 1740436800,
  "iss": "run-sync"
}
```

**Token Expiry:** 24 hours from generation

---

## Security Features

### 1. OTP Rate Limiting
- **Max 5 OTP requests** per phone number dalam 15 menit
- Prevents brute force attacks
- Redis-based tracking

### 2. Password Hashing
- Passwords hashed dengan bcrypt
- Salt automatically generated per password
- Never stored in plain text

### 3. JWT Token Validation
- HMAC-SHA256 signature
- Expiry time checked on every request
- Secret key from environment variable

### 4. Account Status Checks
- `is_verified`: Must complete OTP verification
- `is_active`: Admin can deactivate accounts
- Both checked on login and protected endpoints

### 5. Global Rate Limiting
- Auth endpoints: 10 requests per minute
- User endpoints: 20 requests per minute
- Prevents API abuse

---

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "message": "Permintaan tidak valid",
  "error_code": "INVALID_REQUEST",
  "error_field": "body",
  "error_detail": "Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "message": "Unauthorized",
  "error_code": "INVALID_TOKEN",
  "error_field": "Authorization",
  "error_detail": "token is expired"
}
```

### 403 Forbidden
```json
{
  "success": false,
  "message": "Forbidden",
  "error_code": "NOT_VERIFIED",
  "error_field": "user",
  "error_detail": "Akun belum diverifikasi"
}
```

### 429 Too Many Requests
```json
{
  "success": false,
  "message": "Rate limit exceeded",
  "error_code": "RATE_LIMIT",
  "error_field": "body",
  "error_detail": "Terlalu banyak permintaan OTP, silakan coba lagi nanti"
}
```

---

## Testing Flow (Complete Example)

### Scenario: New User Registration & First API Call

```bash
# 1. Register new user
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Runner",
    "email": "jane@runner.com",
    "phone_number": "082112345678",
    "gender": "female",
    "password": "MySecurePass123"
  }')

echo $REGISTER_RESPONSE
# Extract OTP from response (in production, get from SMS/Email)
OTP=$(echo $REGISTER_RESPONSE | jq -r '.data.otp')
echo "OTP: $OTP"

# 2. Verify OTP
VERIFY_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/verify \
  -H "Content-Type: application/json" \
  -d "{
    \"phone_number\": \"082112345678\",
    \"otp_code\": \"$OTP\"
  }")

echo $VERIFY_RESPONSE
# Extract JWT token
TOKEN=$(echo $VERIFY_RESPONSE | jq -r '.data.token')
echo "Token: $TOKEN"

# 3. Use token to create run group (protected endpoint)
curl -X POST http://localhost:8080/runs/groups \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Ladies Morning Jog",
    "avg_pace": 6.0,
    "preferred_distance": 5,
    "latitude": -6.2088,
    "longitude": 106.8456,
    "scheduled_at": "2026-02-23T06:00:00Z",
    "max_member": 10,
    "is_women_only": true
  }'

# 4. Login again later (skip register & verify)
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "082112345678",
    "password": "MySecurePass123"
  }')

NEW_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
echo "New Token: $NEW_TOKEN"
```

---

## Best Practices

### For Frontend Developers
1. **Store JWT token securely** (localStorage/sessionStorage)
2. **Auto-attach token** to all protected API calls
3. **Handle 401 errors** → redirect to login
4. **Handle 403 errors** → show verification/activation message
5. **Implement token refresh** before 24h expiry
6. **Clear token on logout**

### For Backend Developers
1. **Never log passwords or tokens**
2. **Use environment variables** for JWT secret
3. **Implement OTP delivery** (WhatsApp/SMS/Email)
4. **Monitor rate limiting** effectiveness
5. **Add audit logging** for auth events
6. **Consider refresh tokens** for better UX

### For Testing
1. **Use Postman environment variables** for tokens
2. **Test all error scenarios**
3. **Verify rate limiting behavior**
4. **Test token expiry handling**
5. **Check case sensitivity** (phone numbers, emails)

---

## Environment Variables

Required for production:

```env
# JWT Secret (change in production!)
JWT_SECRET=your-super-secret-key-change-in-production

# Redis (for OTP storage)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=run_sync

# OTP Delivery (configure one)
WHATSAPP_API_KEY=your-whatsapp-api-key
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
```

---

## FAQ

**Q: Can I skip OTP verification in development?**  
A: OTP is currently returned in the response for development. In production, remove the `"otp"` field from register/resend-otp responses.

**Q: What happens if I forget my password?**  
A: Password reset feature not yet implemented. Plan: Use OTP to verify identity, then allow password change.

**Q: Can admin bypass is_verified check?**  
A: Not currently. All users must verify. Future: Add admin roles.

**Q: How to invalidate a token (logout)?**  
A: Client-side: Delete token from storage. Server-side: Implement token blacklist (Redis) for revocation.

**Q: Can I use email instead of phone for login?**  
A: Not currently. Phone number is the primary identifier. Future: Support email login.

**Q: OTP not received, what to do?**  
A: Use `/auth/resend-otp` endpoint. Check rate limits (max 5 per 15 min).

---

## Next Steps

- [ ] Implement email/SMS OTP delivery (production)
- [ ] Add password reset flow
- [ ] Implement refresh tokens
- [ ] Add OAuth2 (Google, Facebook login)
- [ ] Implement role-based access control (RBAC)
- [ ] Add 2FA for sensitive operations
- [ ] Implement token blacklist for logout
- [ ] Add account lockout after failed login attempts

---

**Last Updated:** February 21, 2026  
**Version:** 1.0.0  
**Author:** Run-Sync Team
