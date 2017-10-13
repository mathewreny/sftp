package sftp

type Handle struct {
	client *Client
	h      string
}
