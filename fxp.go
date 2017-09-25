package sftp

type FxpAttrs struct {
	Flags        uint32
	Size         uint64
	Uid, Gid     uint32
	Permissions  uint32
	Atime, Mtime uint32
	Extended     [][2]string
}

type FxpName struct {
	Path, Long string
	Attrs      FxpAttrs
}
