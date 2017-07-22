package conn

type Client interface {
	Dial() error
}