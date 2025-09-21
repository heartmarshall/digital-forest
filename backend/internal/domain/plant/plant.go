package plant

import "time"

// Plant представляет цифровое растение в лесу
type Plant struct {
	ID        int
	Author    string
	ImageData string
	CreatedAt time.Time
}
