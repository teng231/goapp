package live

import (
	"fmt"
	"log"
	"sync"

	"github.com/steampoweredtaco/gotiktoklive"
)

// LiveHub quản lý nhiều Room
type LiveHub struct {
	tiktok *gotiktoklive.TikTok
	mu     sync.Mutex
	rooms  map[string]*Room // key = userID
}

// NewLiveHub tạo hub mới
func NewLiveHub() (*LiveHub, error) {
	tiktok, err := gotiktoklive.NewTikTok()
	if err != nil {
		return nil, err
	}
	return &LiveHub{
		tiktok: tiktok,
		rooms:  make(map[string]*Room),
	}, nil
}

// Add tạo room và tự kết nối
func (h *LiveHub) Add(userID string) (*Room, error) {
	h.mu.Lock()
	if _, exists := h.rooms[userID]; exists {
		h.mu.Unlock()
		return nil, fmt.Errorf("room %s already exists", userID)
	}
	room := &Room{
		TikTok: h.tiktok,
		UserID: userID,
		Status: "disconnected",
	}
	h.rooms[userID] = room
	h.mu.Unlock()

	if err := room.Connect(); err != nil {
		h.Remove(userID) // rollback nếu connect fail
		return nil, err
	}
	log.Printf("Room %s connected", userID)
	return room, nil
}

// Remove ngắt kết nối và xóa room
func (h *LiveHub) Remove(userID string) error {
	h.mu.Lock()
	room, exists := h.rooms[userID]
	if !exists {
		h.mu.Unlock()
		return fmt.Errorf("room %s not found", userID)
	}
	delete(h.rooms, userID)
	h.mu.Unlock()

	if err := room.Disconnect(); err != nil {
		return err
	}
	log.Printf("Room %s removed", userID)
	return nil
}

// DisconnectAll tắt tất cả rooms
func (h *LiveHub) DisconnectAll() {
	h.mu.Lock()
	roomsCopy := make(map[string]*Room, len(h.rooms))
	for k, v := range h.rooms {
		roomsCopy[k] = v
	}
	h.rooms = make(map[string]*Room)
	h.mu.Unlock()

	for id, r := range roomsCopy {
		if err := r.Disconnect(); err != nil {
			log.Printf("Room %s disconnect error: %v", id, err)
		}
	}
}
