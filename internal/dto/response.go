package dto

type SuccessResponse[Data any] struct {
	Message  string    `json:"message,omitempty"`
	Data     Data      `json:"data"`
	Metadata *Metadata `json:"metadata,omitempty"`
}
