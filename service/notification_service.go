package service

import (
	"context"
	"log"

	"run-sync/config"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/repository"

	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
)

type NotificationService interface {
	// Send menyimpan notifikasi ke DB lalu kirim FCM push ke semua device aktif user.
	// actorId: user yang memicu event (opsional, bisa nil)
	// refId & refType: referensi ke entity terkait, misal match_id + "match" (opsional)
	Send(userId uuid.UUID, notifType, title, body string, actorId *uuid.UUID, refId, refType *string) error

	// GetByUser mengambil daftar notifikasi milik user beserta jumlah unread.
	GetByUser(userId uuid.UUID, page, limit int) (response.NotificationListResponse, error)

	// MarkAsRead menandai notifikasi berdasarkan list id sebagai sudah dibaca.
	MarkAsRead(userId uuid.UUID, ids []string) error

	// MarkAllAsRead menandai semua notifikasi user sebagai sudah dibaca.
	MarkAllAsRead(userId uuid.UUID) error

	// RegisterDeviceToken menyimpan/memperbarui FCM token device user.
	RegisterDeviceToken(userId uuid.UUID, fcmToken, platform string) error

	// RemoveDeviceToken menghapus FCM token (saat user logout dari device).
	RemoveDeviceToken(fcmToken string) error
}

type notificationService struct {
	notifRepo  repository.NotificationRepository
	deviceRepo repository.UserDeviceTokenRepository
}

func NewNotificationService(
	notifRepo repository.NotificationRepository,
	deviceRepo repository.UserDeviceTokenRepository,
) NotificationService {
	return &notificationService{
		notifRepo:  notifRepo,
		deviceRepo: deviceRepo,
	}
}

func (s *notificationService) Send(
	userId uuid.UUID,
	notifType, title, body string,
	actorId *uuid.UUID,
	refId, refType *string,
) error {
	// 1. Simpan notifikasi ke database
	notif := &entity.Notification{
		Id:      uuid.New(),
		UserId:  userId,
		Type:    notifType,
		Title:   title,
		Body:    body,
		IsRead:  false,
		ActorId: actorId,
		RefId:   refId,
		RefType: refType,
	}
	if err := s.notifRepo.Create(notif); err != nil {
		return err
	}

	// 2. Kirim FCM push ke semua device aktif user (fire-and-forget, tidak memblokir response)
	go s.sendFCMToUser(userId, title, body, notifType, refId, refType)

	return nil
}

// sendFCMToUser mengambil semua token aktif user lalu mengirim FCM multicast.
func (s *notificationService) sendFCMToUser(
	userId uuid.UUID,
	title, body, notifType string,
	refId, refType *string,
) {
	tokens, err := s.deviceRepo.FindByUserId(userId)
	if err != nil || len(tokens) == 0 {
		return
	}

	// Kumpulkan FCM token strings
	fcmTokens := make([]string, 0, len(tokens))
	for _, t := range tokens {
		fcmTokens = append(fcmTokens, t.FCMToken)
	}

	// Build data payload agar Flutter bisa route ke halaman yang benar
	data := map[string]string{
		"type": notifType,
	}
	if refId != nil {
		data["ref_id"] = *refId
	}
	if refType != nil {
		data["ref_type"] = *refType
	}

	message := &messaging.MulticastMessage{
		Tokens: fcmTokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				Sound:       "default",
				ClickAction: "FLUTTER_NOTIFICATION_CLICK",
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
				},
			},
		},
	}

	resp, err := config.FCMClient.SendEachForMulticast(context.Background(), message)
	if err != nil {
		log.Printf("FCM multicast error untuk user %s: %v", userId, err)
		return
	}

	// Hapus token yang sudah tidak valid (unregistered)
	if resp.FailureCount > 0 {
		for i, r := range resp.Responses {
			if !r.Success && messaging.IsRegistrationTokenNotRegistered(r.Error) {
				_ = s.deviceRepo.DeleteByToken(fcmTokens[i])
			}
		}
	}
}

func (s *notificationService) GetByUser(userId uuid.UUID, page, limit int) (response.NotificationListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	notifs, err := s.notifRepo.FindByUserId(userId, limit, offset)
	if err != nil {
		return response.NotificationListResponse{}, err
	}

	unread, err := s.notifRepo.FindUnreadCount(userId)
	if err != nil {
		return response.NotificationListResponse{}, err
	}

	items := make([]response.NotificationResponse, 0, len(notifs))
	for _, n := range notifs {
		item := response.NotificationResponse{
			Id:        n.Id.String(),
			Type:      n.Type,
			Title:     n.Title,
			Body:      n.Body,
			IsRead:    n.IsRead,
			ReadAt:    n.ReadAt,
			CreatedAt: n.CreatedAt,
		}
		if n.ActorId != nil {
			s := n.ActorId.String()
			item.ActorId = &s
		}
		item.RefId = n.RefId
		item.RefType = n.RefType
		items = append(items, item)
	}

	return response.NotificationListResponse{
		Notifications: items,
		UnreadCount:   unread,
	}, nil
}

func (s *notificationService) MarkAsRead(userId uuid.UUID, ids []string) error {
	uuids := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		parsed, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		uuids = append(uuids, parsed)
	}
	if len(uuids) == 0 {
		return nil
	}
	return s.notifRepo.MarkAsRead(uuids, userId)
}

func (s *notificationService) MarkAllAsRead(userId uuid.UUID) error {
	return s.notifRepo.MarkAllAsRead(userId)
}

func (s *notificationService) RegisterDeviceToken(userId uuid.UUID, fcmToken, platform string) error {
	token := &entity.UserDeviceToken{
		Id:       uuid.New(),
		UserId:   userId,
		FCMToken: fcmToken,
		Platform: platform,
		IsActive: true,
	}
	return s.deviceRepo.Upsert(token)
}

func (s *notificationService) RemoveDeviceToken(fcmToken string) error {
	return s.deviceRepo.DeleteByToken(fcmToken)
}
