package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

/* --- JSONB STRUCTURE TYPES --- */
type Item struct {
	Title       string  `json:"title"`
	Desc        string  `json:"desc,omitempty"`
	ModelUrl    string  `json:"modelUrl,omitempty"`
	ImgUrl      string  `json:"imgUrl"`
	Orientation float64 `json:"orientation,omitempty"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Z           float64 `json:"z"`
	Type        string  `json:"type"`
	Volume      float64 `json:"volume,omitempty"`
	Mass        float64 `json:"mass,omitempty"`
	Power       float64 `json:"power,omitempty"`
	Noise       float64 `json:"noise,omitempty"`
	Notes       string  `json:"notes,omitempty"`
}

type Offset struct {
	X float64 `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
	Z float64 `json:"z,omitempty"`
}

type Floor struct {
	Level             int      `json:"level"`
	Volume            float64  `json:"volume,omitempty"`
	Offset            *Offset  `json:"offset,omitempty"`
	Type              string   `json:"type"`
	ModelUrl          string   `json:"modelUrl"`
	AcceptedItemTypes []string `json:"acceptedItemTypes"`
	Items             []Item   `json:"items"`
}

type HubStructure struct {
	Title  string  `json:"title,omitempty"`
	Desc   string  `json:"desc"`
	Author string  `json:"author,omitempty"`
	Floors []Floor `json:"floors"`
}

/* --- Implement Valuer & Scanner --- */
func (s HubStructure) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *HubStructure) Scan(value interface{}) error {
	if value == nil {
		*s = HubStructure{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan HubStructure: %v", value)
	}
	return json.Unmarshal(b, s)
}

/* --- MAIN HUB MODEL --- */
type Hub struct {
	ID        int64        `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Author    string       `json:"author"`
	Title     string       `json:"title,omitempty"`
	Structure HubStructure `json:"structure" gorm:"type:jsonb"`
	CreatedAt time.Time    `json:"created_at" gorm:"autoCreateTime"`
}
