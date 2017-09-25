package sftp

type FxpAttrs struct {
	Flags        uint32
	Size         uint64
	Uid, Gid     uint32
	Permissions  uint32
	Atime, Mtime uint32
	Extended     [][2]string
}

func (a FxpAttrs) Len() uint32 {
	var length uint32 = 4
	if 0 != a.Flags&FILEXFER_ATTR_SIZE {
		length += 8
	}
	if 0 != a.Flags&FILEXFER_ATTR_UIDGID {
		length += 8
	}
	if 0 != a.Flags&FILEXFER_ATTR_PERMISSIONS {
		length += 4
	}
	if 0 != a.Flags&FILEXFER_ATTR_ACMODTIME {
		length += 8
	}
	if 0 != a.Flags&FILEXFER_ATTR_EXTENDED {
		// extended needs a count.
		length += 4
		for _, ex := range a.Extended {
			length += 8 + uint32(len(ex[0])+len(ex[1]))
		}
	}
	return length
}

type FxpName struct {
	Path, Long string
	Attrs      FxpAttrs
}

func (n FxpName) Len() uint32 {
	return 8 + uint32(len(n.Path)+len(n.Long)) + n.Attrs.Len()
}
