package entity

import "encoding/json"

type Order struct {
	UID  string
	Data json.RawMessage
}
