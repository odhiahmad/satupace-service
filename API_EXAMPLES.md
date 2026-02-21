# Run-Sync API - cURL Examples

Base URL: `http://localhost:8080`

## 📝 Table of Contents
- [🔐 Authentication (New!)](#authentication)
- [Users](#users)
- [Run Groups](#run-groups)
- [Run Activities](#run-activities)
- [Direct Matching](#direct-matching)
- [Chat Messages](#chat-messages)
- [Media & Photos](#media--photos)
- [Safety Logs](#safety-logs)
- [WhatsApp Verification](#whatsapp-verification)

---

## 🔐 Authentication

**⚠️ IMPORTANT: Authentication Flow**
1. **Register** - Create account and receive OTP
2. **Verify OTP** - Confirm registration and get JWT token
3. **Login** - Get JWT token for existing verified users
4. **Use Token** - Include JWT in Authorization header for protected endpoints

### Register (Sign Up)
Creates a new user account and sends OTP for verification.

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
      "email": "john@example.com",
      "phone_number": "081234567890",
      "is_verified": false,
      "is_active": false
    },
    "otp": "123456",
    "message": "Kode OTP telah dikirim ke nomor telepon Anda"
  }
}
```

### Verify OTP
Verifies OTP code and activates user account. Returns JWT token.

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
      "is_verified": true,
      "is_active": true
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Login
Login with phone number and password. Only works for verified and active users.

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

### Resend OTP
Resends OTP code if expired or not received.

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

### Using JWT Token in Protected Endpoints
After login or verification, include the JWT token in the Authorization header:

```bash
curl -X GET http://localhost:8080/runs/groups \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## Users

**Note:** User registration is now handled via `/auth/register` endpoint. See [Authentication](#authentication) section.

### Get All Users
```bash
curl -X GET http://localhost:8080/users
```

### Get User by ID
```bash
curl -X GET http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000
```

### Update User
```bash
curl -X PUT http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated",
    "gender": "male"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000
```

---

## Run Groups

### Create Run Group (Auth Required)
```bash
curl -X POST http://localhost:8080/runs/groups \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
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

### Get All Run Groups
```bash
curl -X GET http://localhost:8080/runs/groups
```

### Get Run Group by ID
```bash
curl -X GET http://localhost:8080/runs/groups/GROUP_UUID
```

### Update Run Group (Auth Required)
```bash
curl -X PUT http://localhost:8080/runs/groups/GROUP_UUID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Updated Group Name",
    "max_member": 20,
    "status": "open"
  }'
```

### Delete Run Group (Auth Required)
```bash
curl -X DELETE http://localhost:8080/runs/groups/GROUP_UUID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Join Run Group (Auth Required)
```bash
curl -X POST http://localhost:8080/runs/groups/GROUP_UUID/join \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{}'
```

### Get Group Members
```bash
curl -X GET http://localhost:8080/runs/groups/GROUP_UUID/members
```

### Update Member Status (Auth Required)
```bash
curl -X PUT http://localhost:8080/runs/members/MEMBER_UUID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "status": "confirmed"
  }'
```

### Remove Member (Auth Required)
```bash
curl -X DELETE http://localhost:8080/runs/members/MEMBER_UUID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Run Activities

### Create Run Activity (Auth Required)
```bash
curl -X POST http://localhost:8080/runs/activities \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "distance": 10.5,
    "duration": 3600,
    "avg_pace": 5.7,
    "calories": 650,
    "route_data": "{\"coordinates\": [...]}",
    "started_at": "2026-02-21T06:00:00Z",
    "ended_at": "2026-02-21T07:00:00Z"
  }'
```

### Get Activity by ID
```bash
curl -X GET http://localhost:8080/runs/activities/ACTIVITY_UUID
```

### Get User Activities
```bash
curl -X GET http://localhost:8080/runs/users/USER_UUID/activities
```

---

## Direct Matching

### Create Match Request (Auth Required)
```bash
curl -X POST http://localhost:8080/match \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "user_id_2": "TARGET_USER_UUID",
    "message": "Hey! Would you like to run together?"
  }'
```

### Update Match Status (Auth Required)
```bash
curl -X PATCH http://localhost:8080/match/MATCH_UUID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "status": "accepted"
  }'
```

### Get Match Details (Auth Required)
```bash
curl -X GET http://localhost:8080/match/MATCH_UUID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get My Matches (Auth Required)
```bash
curl -X GET http://localhost:8080/match/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Chat Messages

### Send Direct Message (Auth Required)
```bash
curl -X POST http://localhost:8080/chats/direct \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "match_id": "MATCH_UUID",
    "content": "Hello! Ready for tomorrow run?"
  }'
```

### Get Direct Chat History (Auth Required)
```bash
curl -X GET http://localhost:8080/chats/direct/MATCH_UUID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Send Group Message (Auth Required)
```bash
curl -X POST http://localhost:8080/chats/group \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "group_id": "GROUP_UUID",
    "content": "See you all at 6 AM!"
  }'
```

### Get Group Chat History (Auth Required)
```bash
curl -X GET http://localhost:8080/chats/group/GROUP_UUID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Media & Photos

### Upload Photo (Auth Required)
```bash
curl -X POST http://localhost:8080/media/photos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "url": "https://res.cloudinary.com/demo/image/upload/sample.jpg",
    "type": "profile",
    "is_primary": true
  }'
```

### Get Photo by ID
```bash
curl -X GET http://localhost:8080/media/photos/PHOTO_UUID
```

### Get User Photos
```bash
curl -X GET http://localhost:8080/media/users/USER_UUID/photos
```

---

## Safety Logs

### Create Safety Log (Auth Required)
```bash
curl -X POST http://localhost:8080/media/safety \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "run_group_id": "GROUP_UUID",
    "latitude": -6.2088,
    "longitude": 106.8456,
    "notes": "Starting the run, all good!",
    "emergency_contact": "081234567890"
  }'
```

### Get Safety Log by ID
```bash
curl -X GET http://localhost:8080/media/safety/SAFETY_LOG_UUID
```

---

## WhatsApp Verification

### Register Phone Number
```bash
curl -X POST http://localhost:8080/whatsapp/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "phone_number": "6281234567890"
  }'
```

### Verify OTP Code
```bash
curl -X POST http://localhost:8080/whatsapp/verify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "code": "123456"
  }'
```

---

## 🔑 Getting JWT Token

After creating a user, you would typically have a login endpoint that returns a JWT token. For testing, you can:

1. **Manual Token Generation**: Use the JWT service to generate a token for your user ID
2. **Login Endpoint**: If you have a login endpoint, use it to get the token

Example response with token:
```json
{
  "status": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {...}
  }
}
```

Then use the token in subsequent requests:
```bash
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 📦 Import to Postman

### Method 1: Import as Collection
1. Open Postman
2. Click **Import** button
3. Select **Raw text** tab
4. Copy-paste the curl commands
5. Postman will auto-convert them to requests

### Method 2: Direct Paste
1. In Postman, click **Import**
2. Paste any individual curl command
3. Click **Import**
4. The request will be added to your collection

### Method 3: Create Collection JSON (see below)

---

## 🔧 Environment Variables for Postman

Create environment variables in Postman:

```
BASE_URL = http://localhost:8080
JWT_TOKEN = your_jwt_token_here
USER_ID = user_uuid_here
GROUP_ID = group_uuid_here
MATCH_ID = match_uuid_here
```

Then replace hardcoded values with `{{BASE_URL}}`, `{{JWT_TOKEN}}`, etc.

---

## 🧪 Testing Flow

### 1. Register & Setup
```bash
# 1. Create user
curl -X POST http://localhost:8080/users ...

# 2. Register WhatsApp (optional)
curl -X POST http://localhost:8080/whatsapp/register ...

# 3. Verify OTP
curl -X POST http://localhost:8080/whatsapp/verify ...
```

### 2. Create Run Group
```bash
# Create a group (need JWT token)
curl -X POST http://localhost:8080/runs/groups ...
```

### 3. Join & Interact
```bash
# Join the group
curl -X POST http://localhost:8080/runs/groups/GROUP_UUID/join ...

# Send group message
curl -X POST http://localhost:8080/chats/group ...
```

### 4. Track Activity
```bash
# Create run activity
curl -X POST http://localhost:8080/runs/activities ...

# Create safety log
curl -X POST http://localhost:8080/media/safety ...
```

---

## 💡 Tips

1. **Save Response IDs**: Save UUIDs from responses to use in subsequent requests
2. **Use Variables**: Replace UUIDs with Postman variables for easier testing
3. **Environment Setup**: Create separate environments for dev/staging/prod
4. **Tests**: Add Postman tests to auto-extract tokens and IDs
5. **Rate Limiting**: Some endpoints have rate limits (e.g., WhatsApp: 10 req/min)

---

## 📄 Response Examples

### Success Response
```json
{
  "status": true,
  "message": "Success message",
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "status": false,
  "message": "Error message",
  "errors": [
    {
      "field": "email",
      "message": "Invalid email format"
    }
  ]
}
```

---

## 🐛 Common Issues

### 401 Unauthorized
- Missing or invalid JWT token
- Token expired
- Solution: Get a fresh token

### 400 Bad Request
- Invalid request body
- Missing required fields
- Solution: Check request payload

### 429 Too Many Requests
- Rate limit exceeded
- Solution: Wait before retrying

### 404 Not Found
- Invalid UUID
- Resource doesn't exist
- Solution: Verify the ID exists
