package conn

type Client interface {
	Dial() error
	Codec(*Encoder)
}

// TODO: Add handler's for Message types. But need to implement messages first
