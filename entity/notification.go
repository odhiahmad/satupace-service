package entity

import (
	"time"

	"github.com/google/uuid"
)

// Notification type constants
const (
	// --- Chat ---
	NotifDirectMessage = "direct_message" // pesan baru di direct chat
	NotifGroupMessage  = "group_message"  // pesan baru di group chat

	// --- Match ---
	NotifMatchRequest  = "match_request"  // seseorang mengirim request match
	NotifMatchAccepted = "match_accepted" // request match diterima
	NotifMatchRejected = "match_rejected" // request match ditolak

	// --- Group ---
	NotifGroupInvite        = "group_invite"         // diundang masuk ke group
	NotifGroupJoinRequest   = "group_join_request"   // ada member baru yang minta join (untuk owner/admin)
	NotifGroupJoinApproved  = "group_join_approved"  // permintaan join disetujui
	NotifGroupJoinRejected  = "group_join_rejected"  // permintaan join ditolak
	NotifGroupRoleChanged   = "group_role_changed"   // role member diubah (misal jadi admin)
	NotifGroupMemberKicked  = "group_member_kicked"  // member dikeluarkan dari group
	NotifGroupMemberLeft    = "group_member_left"    // member meninggalkan group (untuk owner/admin)
	NotifGroupFull          = "group_full"           // group sudah penuh
	NotifGroupCompleted     = "group_completed"      // group run selesai
	NotifGroupCancelled     = "group_cancelled"      // group run dibatalkan
	NotifGroupScheduleStart = "group_schedule_start" // grup run mau dimulai (reminder H-1 jam)

	// --- Activity ---
	NotifActivityLogged = "activity_logged" // aktivitas lari berhasil dicatat

	// --- Safety ---
	NotifUserReported   = "user_reported"   // akun kamu dilaporkan oleh pengguna lain
	NotifUserBlocked    = "user_blocked"    // kamu memblokir / diblokir pengguna
	NotifAutoSuspended  = "auto_suspended"  // akun disuspend otomatis karena banyak laporan

	// --- Account / System ---
	NotifAccountVerified    = "account_verified"    // email berhasil diverifikasi
	NotifProfileIncomplete  = "profile_incomplete"  // profil runner belum dilengkapi
	NotifPasswordChanged    = "password_changed"    // password berhasil diubah
	NotifEmailChangeRequest = "email_change_request" // ada permintaan ganti email
)

// Notification is the persisted notification record.
type Notification struct {
	Id         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId     uuid.UUID  `gorm:"type:uuid;not null;index"` // penerima notifikasi
	Type       string     `gorm:"type:varchar(100);not null"`
	Title      string     `gorm:"type:varchar(255);not null"`
	Body       string     `gorm:"type:text;not null"`
	IsRead     bool       `gorm:"default:false"`
	ReadAt     *time.Time
	ActorId    *uuid.UUID `gorm:"type:uuid"` // user yang memicu notifikasi (opsional)
	RefId      *string    `gorm:"type:varchar(255)"` // id referensi (match_id, group_id, message_id, dll)
	RefType    *string    `gorm:"type:varchar(100)"` // tipe referensi: "match", "group", "message", dll
	CreatedAt  time.Time
}
