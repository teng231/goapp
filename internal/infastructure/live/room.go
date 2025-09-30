package live

import (
	"context"
	"fmt"
	"sync"

	"github.com/steampoweredtaco/gotiktoklive"
)

// IRoom định nghĩa các hành vi của một room
type IRoom interface {
	Connect() error
	Disconnect() error
	Event() <-chan gotiktoklive.Event
}

// Room đại diện cho một livestream
type Room struct {
	TikTok *gotiktoklive.TikTok
	Live   *gotiktoklive.Live
	UserID string // username TikTok
	Status string // connected / disconnected
	mu     sync.Mutex
}

func (r *Room) Connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Live != nil {
		return nil
	}
	live, err := r.TikTok.TrackUser(r.UserID)
	if err != nil {
		return err
	}
	r.Live = live
	r.Status = "connected"
	return nil
}

func (r *Room) Disconnect() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Live == nil {
		return nil
	}
	r.Live.Close()
	r.Live = nil
	r.Status = "disconnected"
	return nil
}
func (r *Room) Event(ctx context.Context, newEv chan<- gotiktoklive.Event) {
	// Không nên giữ lock lâu khi loop; chỉ lock khi cần truy cập Live
	r.mu.Lock()
	if r.Live == nil {
		r.mu.Unlock()
		return
	}
	events := r.Live.Events
	user := r.UserID
	r.mu.Unlock()

	for {
		select {
		case <-ctx.Done():
			// context đã bị hủy hoặc hết hạn
			fmt.Printf("[User %s] Event listener canceled: %v\n", user, ctx.Err())
			return
		case ev, ok := <-events:
			if !ok {
				fmt.Printf("[User %s] Events channel closed\n", user)
				return
			}
			switch e := ev.(type) {
			case gotiktoklive.ChatEvent:
				fmt.Printf("[User %s] Chat: %s => %s\n", user, e.User.Username, e.Comment)
				newEv <- e
			}
		}
	}
}
