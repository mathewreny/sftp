package sftp

// Initialization packets
const (
	INIT byte = 1 + iota
	VERSION
)

// Requests to server
const (
	FXP_OPEN byte = 3 + iota
	FXP_CLOSE
	FXP_READ
	FXP_WRITE
	FXP_LSTAT
	FXP_FSTAT
	FXP_SETSTAT
	FXP_FSETSTAT
	FXP_OPENDIR
	FXP_READDIR
	FXP_REMOVE
	FXP_MKDIR
	FXP_RMDIR
	FXP_REALPATH
	FXP_STAT
	FXP_RENAME
	FXP_READLINK
	FXP_SYMLINK
)

// Responses from server
const (
	FXP_STATUS byte = 101 + iota
	FXP_HANDLE
	FXP_DATA
	FXP_NAME
	FXP_ATTRS
)

// Ignored by sshftp
const (
	FXP_EXTENDED byte = 200 + iota
	FXP_EXTENDEDREPLY
)

// Masks defined in the ietf at the following link
// https://tools.ietf.org/html/draft-ietf-secsh-filexfer-02#section-5
const (
	FILEXFER_ATTR_SIZE        uint32 = 0x00000001
	FILEXFER_ATTR_UIDGID      uint32 = 0x00000002
	FILEXFER_ATTR_PERMISSIONS uint32 = 0x00000004
	FILEXFER_ATTR_ACMODTIME   uint32 = 0x00000008
	FILEXFER_ATTR_EXTENDED    uint32 = 0x80000000
)

// Defined in https://tools.ietf.org/html/draft-ietf-secsh-filexfer-02#page-20
// The documentation for these codes was copied from the link above.
const (
	// Indicates successful completion of the operation.
	//
	// Note: Status responses with this code will never be returned as a
	// go "error" in the sshftp library. This code will not be seen by clients.
	STATUS_OK uint32 = iota
	// Indicates end-of-file condition; for SSH_FX_READ it means that no more
	// data is available in the file, and for SSH_FX_READDIR it indicates that
	// no more files are contained in the directory.
	STATUS_EOF
	// Is returned when a reference is made to a file which should exist but
	// doesn't.
	STATUS_NO_SUCH_FILE
	// Is returned when the authenticated user does not have sufficient
	// permissions to perform the operation.
	STATUS_PERMISSION_DENIED
	// Is a generic catch-all error message; it should be returned if an error
	// occurs for which there is no more specific error code defined.
	STATUS_FAILURE
	// May be returned if a badly formatted packet or protocol incompatibility
	// is detected.
	STATUS_BAD_MESSAGE
	// Is a pseudo-error which indicates that the client has no connection to
	// the server (it can only be generated locally by the client, and MUST NOT
	// be returned by servers).
	STATUS_NO_CONNECTION
	// Is a pseudo-error which indicates that the connection to the server has
	// been lost (it can only be generated locally by the client, and MUST NOT
	// be returned by servers).
	STATUS_CONNECTION_LOST
	// Indicates that an attempt was made to perform an operation which is not
	// supported for the server (it may be generated locally by the client if
	// e.g.  the version number exchange indicates that a required feature is
	// not supported by the server, or it may be returned by the server if the
	// server does not implement an operation).
	STATUS_OP_UNSUPPORTED
)
