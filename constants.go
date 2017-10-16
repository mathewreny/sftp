package sftp

// SFTP initialization request and response.
const (
	FXP_INIT byte = 1 + iota
	FXP_VERSION
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

// For server extensions.
const (
	FXP_EXTENDED byte = 200 + iota
	FXP_EXTENDEDREPLY
)

// Pflags bit masks used when opening files. Documentation copied from
//
// https://tools.ietf.org/html/draft-ietf-secsh-filexfer-02#section-6.3
const (
	// Open the file for reading.
	FXF_READ uint32 = 0x00000001
	// Open the file for writing.  If both this and SSH_FXF_READ are
	// specified, the file is opened for both reading and writing.
	FXF_WRITE uint32 = 0x00000002
	// Force all writes to append data at the end of the file.
	FXF_APPEND uint32 = 0x00000004
	// If this flag is specified, then a new file will be created if one
	// does not already exist (if O_TRUNC is specified, the new file will
	// be truncated to zero length if it previously exists).
	FXF_CREAT uint32 = 0x00000008
	// Same thing as FXF_CREAT, but with an 'E' for Ken Thompson.
	FXF_CREATE uint32 = FXF_CREAT
	// Forces an existing file with the same name to be truncated to zero
	// length when creating a file by specifying SSH_FXF_CREAT.
	// SSH_FXF_CREAT MUST also be specified if this flag is used.
	FXF_TRUNC uint32 = 0x00000010
	// Causes the request to fail if the named file already exists.
	// SSH_FXF_CREAT MUST also be specified if this flag is used.
	FXF_EXCL uint32 = 0x00000020
)

// Bit masks for Attrs. See the ietf documentation for more information.
//
// https://tools.ietf.org/html/draft-ietf-secsh-filexfer-02#section-5
const (
	FILEXFER_ATTR_SIZE        uint32 = 0x00000001
	FILEXFER_ATTR_UIDGID      uint32 = 0x00000002
	FILEXFER_ATTR_PERMISSIONS uint32 = 0x00000004
	FILEXFER_ATTR_ACMODTIME   uint32 = 0x00000008
	FILEXFER_ATTR_EXTENDED    uint32 = 0x80000000
)

// Status codes found within the Status packet. Documentation copied from
//
// https://tools.ietf.org/html/draft-ietf-secsh-filexfer-02#page-20
//
// Note: Status responses with STATUS_OK will always be returned as
// a nil error by this library's Client type.
const (
	// Indicates successful completion of the operation.
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
