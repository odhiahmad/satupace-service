# Run-Sync API Service

A high-performance backend service for a running social app built with Go, Gin, and PostgreSQL. Connects runners for group activities, direct matching, and safety tracking.

## 🚀 Features

- **User Management**: Registration, authentication, profile management
- **Runner Profiles**: Running experience, preferred pace, and distance tracking  
- **Run Groups**: Create, join, and manage running groups with location-based discovery
- **Direct Matching**: Connect with other runners for 1-on-1 runs
- **Real-time Messaging**: Direct and group chat functionality
- **Activity Tracking**: Record and monitor running activities
- **Safety Logs**: Emergency contact and safety check-ins
- **Photo Upload**: Profile and activity photos via Cloudinary
- **WhatsApp Integration**: OTP verification via WhatsApp
- **Redis Caching**: High-performance data caching and rate limiting

## 🛠 Tech Stack

- **Language**: Go 1.25
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with GORM
- **Caching**: Redis
- **Authentication**: JWT tokens
- **Media Storage**: Cloudinary
- **Messaging**: WhatsApp Business API (via whatsmeow)
- **Validation**: go-playground/validator

## 📁 Project Structure

```
run-sync/
├── config/              # Database and Redis configuration
├── controller/          # HTTP request handlers
├── data/               
│   ├── request/        # Request DTOs
│   └── response/       # Response DTOs
├── entity/             # Database models
├── helper/             # Utility functions
│   └── mapper/         # Entity to DTO mappers
├── middleware/         # Auth, CORS, rate limiting
├── repository/         # Database access layer
├── routes/             # API route definitions
├── seeder/             # Database seeders
├── service/            # Business logic layer
├── main.go             # Application entry point
└── create_indexes.sql  # Database performance indexes
```

## 🔧 Setup & Installation

### Prerequisites

- Go 1.25 or higher
- PostgreSQL 14+
- Redis 6+
- Cloudinary account (for image uploads)
- WhatsApp Business API credentials (optional, for OTP)

### Environment Variables

Create a `.env` file in the root directory:

```env
# Server Configuration
GIN_MODE=debug
PORT=8080

# Database Configuration
DB_USER=postgres
DB_PASS=your_password
DB_HOST=localhost
DB_NAME=run_sync
DB_PORT=5432

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT Secret
JWT_SECRET=your_secret_key_here

# Cloudinary Configuration
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret

# Email Configuration (for notifications)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password

# WhatsApp Configuration (optional)
DB_NAME_WHATSAPP=whatsapp_sessions
```

### Installation Steps

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd run-sync
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup database**
   ```bash
   # Create database
   createdb run_sync
   
   # Run performance indexes
   psql -d run_sync -f create_indexes.sql
   ```

4. **Run the application**
   ```bash
   # Development mode
   GIN_MODE=debug go run main.go
   
   # Production mode
   GIN_MODE=release go run main.go
   ```

5. **Build for production**
   ```bash
   go build -o run-sync .
   ./run-sync
   ```

## 📡 API Endpoints

### Users
- `POST /users` - Create new user
- `GET /users` - Get all users
- `GET /users/:id` - Get user by ID
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

### Run Groups
- `POST /runs/groups` - Create run group (auth required)
- `GET /runs/groups` - List all run groups
- `GET /runs/groups/:id` - Get group details
- `PUT /runs/groups/:id` - Update group (auth required)
- `DELETE /runs/groups/:id` - Delete group (auth required)
- `POST /runs/groups/:id/join` - Join group (auth required)
- `GET /runs/groups/:groupId/members` - List group members

### Run Activities
- `POST /runs/activities` - Create activity (auth required)
- `GET /runs/activities/:id` - Get activity details
- `GET /runs/users/:userId/activities` - Get user activities

### Direct Matching
- `POST /match` - Create match request (auth required)
- `PATCH /match/:id` - Update match status (auth required)
- `GET /match/:id` - Get match details (auth required)
- `GET /match/me` - Get user matches (auth required)

### Chat
- `POST /chats/direct` - Send direct message (auth required)
- `GET /chats/direct/:matchId` - Get direct chat history (auth required)
- `POST /chats/group` - Send group message (auth required)
- `GET /chats/group/:groupId` - Get group chat history (auth required)

### Media
- `POST /media/photos` - Upload photo (auth required)
- `GET /media/photos/:id` - Get photo details
- `GET /media/users/:userId/photos` - Get user photos

### Safety
- `POST /media/safety` - Create safety log (auth required)
- `GET /media/safety/:id` - Get safety log details

### WhatsApp Verification
- `POST /whatsapp/register` - Register phone number
- `POST /whatsapp/verify` - Verify OTP code

## 🔐 Authentication

The API uses JWT tokens for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## ⚡ Performance Optimizations

### Database Indexes
All critical database queries are optimized with proper indexes. Run `create_indexes.sql` to apply them.

### Caching Strategy
- Redis caching for frequently accessed data
- Rate limiting to prevent abuse
- Connection pooling for database efficiency

### Code Optimizations
- Removed unused functions and files
- Optimized import statements
- Eliminated duplicate code
- Efficient query patterns using GORM
- Proper error handling throughout

## 🧹 Code Quality Improvements

### What Was Fixed
1. **Module Name**: Changed from `loka-kasir` to `run-sync` throughout the codebase
2. **Import Cleanup**: Fixed all import paths and removed unused imports
3. **Database Models**: Removed business/POS related entities, kept only running app models
4. **Helper Functions**: Removed unused utility functions, kept only essential ones
5. **Type Consistency**: Fixed type mismatches in Run Group entity
6. **JWT Service**: Simplified and optimized for the running app use case
7. **CORS Middleware**: Moved to dedicated middleware package
8. **Redis Helper**: Simplified to only include used functions
9. **Error Handling**: Improved error handling throughout services

### Removed Files/Code
- Removed: `home_service.go` (unused)
- Removed: Business-related entity migrations
- Removed: Unused product autocomplete functions
- Removed: Duplicate `pointer.go` helper
- Cleaned: Unused dependencies from go.mod

## 📊 Database Schema

Key entities:
- **users**: User accounts and authentication
- **runner_profiles**: Runner experience and preferences
- **run_groups**: Group running sessions
- **run_group_members**: Group membership
- **run_activities**: Individual running activities
- **direct_matches**: Runner matching
- **direct_chat_messages**: Direct messaging
- **group_chat_messages**: Group messaging
- **user_photos**: User profile photos
- **safety_logs**: Safety check-ins

## 🔄 Migration & Seeding

The application auto-migrates database schemas in development mode. Sample data is seeded automatically via seeders in `seeder/` directory.

## 🐛 Debugging

Enable debug mode for detailed logging:
```env
GIN_MODE=debug
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📝 License

This project is licensed under the MIT License.

## 👨‍💻 Author

Built with ❤️ for the running community

## 🆘 Support

For issues and questions, please open an issue in the repository.
