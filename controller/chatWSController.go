package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/helper"
	"run-sync/repository"
	"run-sync/service"
	ws "run-sync/websocket"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins; tighten in production
	},
}

// ────────────────────────────────────────────────
// WebSocket message types (JSON over WS)
// ────────────────────────────────────────────────

// WSIncomingMessage is what the client sends.
type WSIncomingMessage struct {
	Type    string `json:"type"`    // "message", "typing", "read"
	Message string `json:"message"` // for type=message
}

// WSOutgoingMessage is what the server broadcasts.
type WSOutgoingMessage struct {
	Type       string                 `json:"type"`        // "message", "typing", "system"
	Id         string                 `json:"id"`          // message ID
	RoomID     string                 `json:"room_id"`     // match_id or group_id
	SenderId   string                 `json:"sender_id"`   // who sent it
	SenderName string                 `json:"sender_name"` // display name
	Sender     *response.UserResponse `json:"sender,omitempty"`
	Message    string                 `json:"message"`
	CreatedAt  time.Time              `json:"created_at"`
}

// ────────────────────────────────────────────────
// Chat WS Controller
// ────────────────────────────────────────────────

type ChatWSController interface {
	// WebSocket endpoints
	HandleDirectChat(ctx *gin.Context)
	HandleGroupChat(ctx *gin.Context)

	// REST endpoints (history)
	GetDirectHistory(ctx *gin.Context)
	GetGroupHistory(ctx *gin.Context)
	DeleteMessage(ctx *gin.Context)
}

type chatWSController struct {
	hub            *ws.Hub
	directChatRepo repository.DirectChatMessageRepository
	groupChatRepo  repository.GroupChatMessageRepository
	userRepo       repository.UserRepository
	memberRepo     repository.RunGroupMemberRepository
	jwtService     service.JWTService
}

func NewChatWSController(
	hub *ws.Hub,
	directChatRepo repository.DirectChatMessageRepository,
	groupChatRepo repository.GroupChatMessageRepository,
	userRepo repository.UserRepository,
	memberRepo repository.RunGroupMemberRepository,
	jwtService service.JWTService,
) ChatWSController {
	return &chatWSController{
		hub:            hub,
		directChatRepo: directChatRepo,
		groupChatRepo:  groupChatRepo,
		userRepo:       userRepo,
		memberRepo:     memberRepo,
		jwtService:     jwtService,
	}
}

// ────────────────────────────────────────────────
// HandleDirectChat — ws://host/ws/direct/:matchId?token=JWT
// ────────────────────────────────────────────────

func (c *chatWSController) HandleDirectChat(ctx *gin.Context) {
	userID, err := c.authenticateWS(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, helper.BuildErrorResponse(
			"Unauthorized", "UNAUTHORIZED", "token", err.Error(), nil,
		))
		return
	}

	matchID := ctx.Param("matchId")
	if matchID == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Match ID required", "INVALID_REQUEST", "matchId", "matchId is required", nil,
		))
		return
	}

	roomID := "direct:" + matchID

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WS upgrade error: %v", err)
		return
	}

	// Look up sender name for system messages
	senderName := c.getUserName(userID)

	client := &ws.Client{
		Hub:      c.hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		UserName: senderName,
		RoomID:   roomID,
	}

	c.hub.Register(client)

	// Write pump in goroutine
	go client.WritePump()

	// Read pump blocks — on each message, save to DB & broadcast
	client.ReadPump(func(cl *ws.Client, raw []byte) {
		c.handleDirectMessage(cl, raw, matchID)
	})
}

// ────────────────────────────────────────────────
// HandleGroupChat — ws://host/ws/group/:groupId?token=JWT
// ────────────────────────────────────────────────

func (c *chatWSController) HandleGroupChat(ctx *gin.Context) {
	userID, err := c.authenticateWS(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, helper.BuildErrorResponse(
			"Unauthorized", "UNAUTHORIZED", "token", err.Error(), nil,
		))
		return
	}

	groupID := ctx.Param("groupId")
	if groupID == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Group ID required", "INVALID_REQUEST", "groupId", "groupId is required", nil,
		))
		return
	}

	// Verify user is a member of this group
	userUUID, _ := uuid.Parse(userID)
	groupUUID, _ := uuid.Parse(groupID)
	member, err := c.memberRepo.FindByGroupAndUser(groupUUID, userUUID)
	if err != nil || member == nil || member.Status != "joined" {
		ctx.JSON(http.StatusForbidden, helper.BuildErrorResponse(
			"Anda bukan anggota grup ini", "FORBIDDEN", "groupId", "harus bergabung terlebih dahulu", nil,
		))
		return
	}

	roomID := "group:" + groupID

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WS upgrade error: %v", err)
		return
	}

	// Look up sender name for system messages
	groupSenderName := c.getUserName(userID)

	client := &ws.Client{
		Hub:      c.hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		UserName: groupSenderName,
		RoomID:   roomID,
	}

	c.hub.Register(client)

	go client.WritePump()

	client.ReadPump(func(cl *ws.Client, raw []byte) {
		c.handleGroupMessage(cl, raw, groupID)
	})
}

// ────────────────────────────────────────────────
// REST: Chat history
// ────────────────────────────────────────────────

func (c *chatWSController) GetDirectHistory(ctx *gin.Context) {
	matchID, err := uuid.Parse(ctx.Param("matchId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Invalid match ID", "INVALID_REQUEST", "matchId", err.Error(), nil,
		))
		return
	}

	messages, err := c.directChatRepo.FindByMatchId(matchID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil pesan", "FETCH_FAILED", "body", err.Error(), nil,
		))
		return
	}

	var result []response.DirectChatMessageDetailResponse
	for _, msg := range messages {
		sender := c.getUserResponse(msg.SenderId)
		senderName := ""
		if sender != nil && sender.Name != nil {
			senderName = *sender.Name
		}
		result = append(result, response.DirectChatMessageDetailResponse{
			Id:         msg.Id.String(),
			MatchId:    msg.MatchId.String(),
			SenderId:   msg.SenderId.String(),
			SenderName: senderName,
			Sender:     sender,
			Message:    msg.Message,
			CreatedAt:  msg.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil pesan", result))
}

func (c *chatWSController) GetGroupHistory(ctx *gin.Context) {
	groupID, err := uuid.Parse(ctx.Param("groupId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Invalid group ID", "INVALID_REQUEST", "groupId", err.Error(), nil,
		))
		return
	}

	// Verify user is a member of this group
	userId := ctx.MustGet("user_id").(uuid.UUID)
	member, mErr := c.memberRepo.FindByGroupAndUser(groupID, userId)
	if mErr != nil || member == nil || member.Status != "joined" {
		ctx.JSON(http.StatusForbidden, helper.BuildErrorResponse(
			"Anda bukan anggota grup ini", "FORBIDDEN", "groupId", "harus bergabung terlebih dahulu", nil,
		))
		return
	}

	messages, err := c.groupChatRepo.FindByGroupId(groupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil pesan", "FETCH_FAILED", "body", err.Error(), nil,
		))
		return
	}

	var result []response.GroupChatMessageDetailResponse
	for _, msg := range messages {
		sender := c.getUserResponse(msg.SenderId)
		senderName := ""
		if sender != nil && sender.Name != nil {
			senderName = *sender.Name
		}
		result = append(result, response.GroupChatMessageDetailResponse{
			Id:         msg.Id.String(),
			GroupId:    msg.GroupId.String(),
			SenderId:   msg.SenderId.String(),
			SenderName: senderName,
			Sender:     sender,
			Message:    msg.Message,
			CreatedAt:  msg.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil pesan grup", result))
}

func (c *chatWSController) DeleteMessage(ctx *gin.Context) {
	msgID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Invalid message ID", "INVALID_REQUEST", "id", err.Error(), nil,
		))
		return
	}

	chatType := ctx.Query("type") // "direct" or "group"
	if chatType == "group" {
		if err := c.groupChatRepo.Delete(msgID); err != nil {
			ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
				"Gagal menghapus pesan", "DELETE_FAILED", "body", err.Error(), nil,
			))
			return
		}
	} else {
		if err := c.directChatRepo.Delete(msgID); err != nil {
			ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
				"Gagal menghapus pesan", "DELETE_FAILED", "body", err.Error(), nil,
			))
			return
		}
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Pesan berhasil dihapus", nil))
}

// ────────────────────────────────────────────────
// Internal helpers
// ────────────────────────────────────────────────

// authenticateWS gets JWT from query parameter (?token=xxx) for WebSocket.
func (c *chatWSController) authenticateWS(ctx *gin.Context) (string, error) {
	tokenStr := ctx.Query("token")
	if tokenStr == "" {
		// Also check Authorization header
		auth := ctx.GetHeader("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			tokenStr = strings.TrimPrefix(auth, "Bearer ")
		}
	}
	if tokenStr == "" {
		return "", fmt.Errorf("token is required")
	}

	token, err := c.jwtService.ValidateToken(tokenStr)
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	uid, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in token")
	}

	return uid, nil
}

func (c *chatWSController) handleDirectMessage(client *ws.Client, raw []byte, matchID string) {
	var incoming WSIncomingMessage
	if err := json.Unmarshal(raw, &incoming); err != nil {
		log.Printf("WS: Invalid JSON from %s: %v", client.UserID, err)
		return
	}

	switch incoming.Type {
	case "message":
		if incoming.Message == "" {
			return
		}

		senderUUID, _ := uuid.Parse(client.UserID)
		matchUUID, _ := uuid.Parse(matchID)

		msg := entity.DirectChatMessage{
			Id:        uuid.New(),
			MatchId:   matchUUID,
			SenderId:  senderUUID,
			Message:   incoming.Message,
			CreatedAt: time.Now(),
		}

		// Save to database asynchronously
		go func() {
			if err := c.directChatRepo.Create(&msg); err != nil {
				log.Printf("WS: Failed to save direct message: %v", err)
			}
		}()

		sender := c.getUserResponse(senderUUID)
		outgoing := WSOutgoingMessage{
			Type:       "message",
			Id:         msg.Id.String(),
			RoomID:     matchID,
			SenderId:   client.UserID,
			SenderName: client.UserName,
			Sender:     sender,
			Message:    msg.Message,
			CreatedAt:  msg.CreatedAt,
		}

		data, _ := json.Marshal(outgoing)
		c.hub.BroadcastToRoom(client.RoomID, data)

	case "typing":
		out, _ := json.Marshal(map[string]string{
			"type":        "typing",
			"sender_id":   client.UserID,
			"sender_name": client.UserName,
			"room_id":     matchID,
		})
		c.hub.BroadcastToRoom(client.RoomID, out)

	case "read":
		out, _ := json.Marshal(map[string]string{
			"type":        "read",
			"sender_id":   client.UserID,
			"sender_name": client.UserName,
			"room_id":     matchID,
		})
		c.hub.BroadcastToRoom(client.RoomID, out)
	}
}

func (c *chatWSController) handleGroupMessage(client *ws.Client, raw []byte, groupID string) {
	var incoming WSIncomingMessage
	if err := json.Unmarshal(raw, &incoming); err != nil {
		log.Printf("WS: Invalid JSON from %s: %v", client.UserID, err)
		return
	}

	switch incoming.Type {
	case "message":
		if incoming.Message == "" {
			return
		}

		senderUUID, _ := uuid.Parse(client.UserID)
		groupUUID, _ := uuid.Parse(groupID)

		msg := entity.GroupChatMessage{
			Id:        uuid.New(),
			GroupId:   groupUUID,
			SenderId:  senderUUID,
			Message:   incoming.Message,
			CreatedAt: time.Now(),
		}

		// Save to database asynchronously
		go func() {
			if err := c.groupChatRepo.Create(&msg); err != nil {
				log.Printf("WS: Failed to save group message: %v", err)
			}
		}()

		sender := c.getUserResponse(senderUUID)
		outgoing := WSOutgoingMessage{
			Type:       "message",
			Id:         msg.Id.String(),
			RoomID:     groupID,
			SenderId:   client.UserID,
			SenderName: client.UserName,
			Sender:     sender,
			Message:    msg.Message,
			CreatedAt:  msg.CreatedAt,
		}

		data, _ := json.Marshal(outgoing)
		c.hub.BroadcastToRoom(client.RoomID, data)

	case "typing":
		out, _ := json.Marshal(map[string]string{
			"type":        "typing",
			"sender_id":   client.UserID,
			"sender_name": client.UserName,
			"room_id":     groupID,
		})
		c.hub.BroadcastToRoom(client.RoomID, out)
	}
}

func (c *chatWSController) getUserResponse(userID uuid.UUID) *response.UserResponse {
	user, err := c.userRepo.FindById(userID)
	if err != nil {
		return nil
	}
	return &response.UserResponse{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		HasProfile:  user.HasProfile,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

// getUserName looks up a user's display name by their ID string.
func (c *chatWSController) getUserName(userIDStr string) string {
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return ""
	}
	user, err := c.userRepo.FindById(uid)
	if err != nil {
		return ""
	}
	return helper.DerefOrEmpty(user.Name)
}
