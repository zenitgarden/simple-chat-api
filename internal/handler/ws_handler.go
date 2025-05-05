package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"github.com/zenitgarden/simple-chat-api/internal/dto"
	"github.com/zenitgarden/simple-chat-api/internal/service"
)

type WSMessage struct {
	Type           string    `json:"type"`
	ConversationID uuid.UUID `json:"conversationId"`
	Content        string    `json:"content"`
}

type WSConnection struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	ConversationID uuid.UUID
	Conn           *websocket.Conn
	LastActive     time.Time
	WriteMu        sync.Mutex
}

var connections = make(map[uuid.UUID]*WSConnection)
var connMu sync.RWMutex

type WsHandler struct {
	messageService     *service.MessageService
	participantService *service.ParticipantService
	userService        *service.UserService
}

func NewWsHandler(messageService *service.MessageService, participantService *service.ParticipantService, userService *service.UserService) *WsHandler {
	return &WsHandler{messageService: messageService, participantService: participantService, userService: userService}
}

func (h *WsHandler) ConnectToWs(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)

	userID, err := uuid.Parse(userId)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	upgrader := websocket.FastHTTPUpgrader{
		CheckOrigin: func(r *fasthttp.RequestCtx) bool { return true },
	}

	return upgrader.Upgrade(ctx.Context(), func(conn *websocket.Conn) {
		defer cleanupConnection(userID)
		registerConn(userID, conn)

		conn.SetPongHandler(func(appData string) error {
			connMu.Lock()
			if c, ok := connections[userID]; ok {
				c.LastActive = time.Now()
			}
			connMu.Unlock()
			return nil
		})

		go sendPingLoop(userID)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("WebSocket read error:", err)
				break
			}

			fmt.Println("Received:", string(msg))

			connMu.Lock()
			connections[userID].LastActive = time.Now()
			connMu.Unlock()

			var parsed WSMessage
			if err := json.Unmarshal(msg, &parsed); err != nil {
				log.Println("Invalid message format:", err)
				continue
			}

			message := dto.CreateMessageRequest{
				Type:           parsed.Type,
				ConversationID: parsed.ConversationID,
				Content:        parsed.Content,
				SenderID:       userID,
			}

			switch message.Type {
			case "typing", "stop_typing":
				go func(msg dto.CreateMessageRequest) {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()

					if msg.ConversationID == uuid.Nil {
						return
					}

					participantIDs, err := h.participantService.FindByConversationId(ctx, msg.ConversationID)
					if err != nil {
						log.Println("Failed to fetch participants:", err)
						return
					}

					user, err := h.userService.FindUserByID(ctx, msg.SenderID)

					userName := "Unknown"
					if err == nil {
						userName = user.Name
					} else {
						log.Printf("Failed to fetch user name for %v: %v", msg.SenderID, err)
					}

					data, _ := json.Marshal(map[string]interface{}{
						"type":           msg.Type,
						"userId":         msg.SenderID,
						"userName":       userName,
						"conversationId": msg.ConversationID,
					})

					broadcastToParticipants(participantIDs, data, msg.SenderID, msg.Type)
				}(message)
			case "message":
				go func(msg dto.CreateMessageRequest) {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()

					if msg.ConversationID == uuid.Nil {
						return
					}

					savedMsg, err := h.messageService.CreateMessage(ctx, &msg)
					if err != nil {
						log.Println("DB store error:", err)
						return
					}

					participantIDs, err := h.participantService.FindByConversationId(ctx, msg.ConversationID)
					if err != nil {
						log.Println("Failed to fetch participants:", err)
						return
					}

					user, err := h.userService.FindUserByID(ctx, msg.SenderID)

					userName := "Unknown"
					if err == nil {
						userName = user.Name
					} else {
						log.Printf("Failed to fetch user name for %v: %v", msg.SenderID, err)
					}

					data, _ := json.Marshal(map[string]any{
						"type":           "message",
						"userId":         msg.SenderID,
						"userName":       userName,
						"conversationId": msg.ConversationID,
						"content":        savedMsg.Content,
						"time":           savedMsg.SentAt,
					})

					broadcastToParticipants(participantIDs, data, msg.SenderID, "message")
				}(message)
			}
		}
	})
}

func ConnectionCount() int {
	connMu.RLock()
	defer connMu.RUnlock()
	return len(connections)
}

func registerConn(userID uuid.UUID, conn *websocket.Conn) {
	connMu.Lock()
	defer connMu.Unlock()

	if existing, ok := connections[userID]; ok {
		_ = existing.Conn.Close()
		delete(connections, userID)
	}

	wsConn := &WSConnection{
		ID:             uuid.New(),
		Conn:           conn,
		ConversationID: uuid.Nil,
		UserID:         userID,
		LastActive:     time.Now(),
	}

	connections[userID] = wsConn
}

func cleanupConnection(userID uuid.UUID) {
	connMu.Lock()
	defer connMu.Unlock()

	if existing, ok := connections[userID]; ok {
		if err := existing.Conn.Close(); err != nil {
			log.Printf("WebSocket close error for %v: %v", userID, err)
		}
		delete(connections, userID)
	} else {
		log.Printf("No connection found for participant: %v", userID)
	}
}

func broadcastToParticipants(participantIDs []*uuid.UUID, data []byte, senderID uuid.UUID, messageType string) {
	connMu.RLock()
	defer connMu.RUnlock()

	for _, pid := range participantIDs {
		if pid == nil || (*pid == senderID && messageType != "message") {
			continue
		}

		if conn, ok := connections[*pid]; ok {
			conn.WriteMu.Lock()
			err := conn.Conn.WriteMessage(websocket.TextMessage, data)
			conn.WriteMu.Unlock()

			if err != nil {
				log.Printf("Broadcast error to %v: %v", *pid, err)
			}
		} else {
			log.Printf("No connection found for participant: %v", *pid)
		}
	}
}

func sendPingLoop(userID uuid.UUID) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		connMu.RLock()
		conn, ok := connections[userID]
		connMu.RUnlock()

		if !ok {
			return
		}

		if time.Since(conn.LastActive) > 60*time.Second {
			log.Printf("Connection timed out for %v", userID)
			cleanupConnection(userID)
			return
		}
		conn.WriteMu.Lock()
		err := conn.Conn.WriteMessage(websocket.PingMessage, []byte("ping")); 
		conn.WriteMu.Unlock()
		if err != nil {
			log.Printf("Ping failed for %v: %v", userID, err)
			cleanupConnection(userID)
			return
		}
	}
}
