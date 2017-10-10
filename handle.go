package sftp

type Handle struct {
	conn *Conn
	h    string
}
