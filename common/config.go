package common

// Config contains shared config parameters, common to the source and
// destination. If you don't need shared parameters you can entirely remove this
// file.
type Config struct {
	// Topic is the topic to publish to when receiving a record to write
	Topic string `json:"topic" validate:"required"`
}
