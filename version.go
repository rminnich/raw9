package main

// 9P2000 message types
const (
	Tversion = 100 + iota
	Rversion
	Tauth
	Rauth
	Tattach
	Rattach
	Terror
	Rerror
	Tflush
	Rflush
	Twalk
	Rwalk
	Topen
	Ropen
	Tcreate
	Rcreate
	Tread
	Rread
	Twrite
	Rwrite
	Tclunk
	Rclunk
	Tremove
	Rremove
	Tstat
	Rstat
	Twstat
	Rwstat
	Tlast
)

const (
	MSIZE   = 2*1048576 + IOHDRSZ // default message size (1048576+IOHdrSz)
	IOHDRSZ = 24                  // the non-data size of the Twrite messages
	PORT    = 564                 // default port for 9P file servers
	NumFID  = 1 << 16
	QIDLen  = 13
)

// QID types
const (
	QTDIR     = 0x80 // directories
	QTAPPEND  = 0x40 // append only files
	QTEXCL    = 0x20 // exclusive use files
	QTMOUNT   = 0x10 // mounted channel
	QTAUTH    = 0x08 // authentication file
	QTTMP     = 0x04 // non-backed-up file
	QTSYMLINK = 0x02 // symbolic link (Unix, 9P2000.u)
	QTLINK    = 0x01 // hard link (Unix, 9P2000.u)
	QTFILE    = 0x00
)

// Flags for the mode field in Topen and Tcreate messages
const (
	OREAD   = 0x0    // open read-only
	OWRITE  = 0x1    // open write-only
	ORDWR   = 0x2    // open read-write
	OEXEC   = 0x3    // execute (== read but check execute permission)
	OTRUNC  = 0x10   // or'ed in (except for exec), truncate file first
	OCEXEC  = 0x20   // or'ed in, close on exec
	ORCLOSE = 0x40   // or'ed in, remove on close
	OAPPEND = 0x80   // or'ed in, append only
	OEXCL   = 0x1000 // or'ed in, exclusive client use
)

// File modes
const (
	DMDIR    = 0x80000000 // mode bit for directories
	DMAPPEND = 0x40000000 // mode bit for append only files
	DMEXCL   = 0x20000000 // mode bit for exclusive use files
	DMMOUNT  = 0x10000000 // mode bit for mounted channel
	DMAUTH   = 0x08000000 // mode bit for authentication file
	DMTMP    = 0x04000000 // mode bit for non-backed-up file
	DMREAD   = 0x4        // mode bit for read permission
	DMWRITE  = 0x2        // mode bit for write permission
	DMEXEC   = 0x1        // mode bit for execute permission
)

const (
	NOTAG = 0xFFFF     // no tag specified
	NOFID = 0xFFFFFFFF // no fid specified
	// We reserve tag NOTAG and tag 0. 0 is a troublesome value to pass
	// around, since it is also a default value and using it can hide errors
	// in the code.
	NumTags = 1<<16 - 2
)

// Error values
const (
	EPERM   = 1
	ENOENT  = 2
	EIO     = 5
	EACCES  = 13
	EEXIST  = 17
	ENOTDIR = 20
	EINVAL  = 22
)

// Types contained in 9p messages.
type (
	MType      byte
	Mode       uint8
	NumEntries uint16
	Tag        uint16
	FID        uint32
	MaxSize    uint32
	Count      int32
	Perm       uint32
	Offset     uint64
	Data       []byte
	// Some []byte fields are encoded with a 16-bit length, e.g. stat data.
	// We use this type to tag such fields. The parameters are still []byte,
	// this was just the only way I could think of to make the stub generator do the right
	// thing.
	DataCnt16 byte // []byte with a 16-bit count.
)

// Error represents a 9P2000 error
type Error struct {
	Err string
}

// File identifier
type QID struct {
	Type    uint8  // type of the file (high 8 bits of the mode)
	Version uint32 // version number for the path
	Path    uint64 // server's unique identification of the file
}

// Dir describes a file
type Dir struct {
	Type    uint16
	Dev     uint32
	QID     QID    // file's QID
	Mode    uint32 // permissions and flags
	Atime   uint32 // last access time in seconds
	Mtime   uint32 // last modified time in seconds
	Length  uint64 // file length in bytes
	Name    string // file name
	User    string // owner name
	Group   string // group name
	ModUser string // name of the last user that modified the file
}

//
//	size[4] Tversion tag[2] msize[4] version[s]
//      size[4] Rversion tag[2] msize[4] version[s]

// assumption: we have a byte slice big enough to hold what
// we will produce.
func xversion(b []byte, version string, msize uint32) int {
	s := 4 + 1 + 2 + 4 + 2 + len(version)
	b[0], b[1], b[2], b[3] = uint8(s), uint8(s>>8), uint8(s>>16), uint8(s>>24)
	b[4] = byte(Tversion)
	b[5], b[6] = uint8(0xff), uint8(0xff)
	b[7], b[8], b[9], b[10] = uint8(msize), uint8(msize>>8), uint8(msize>>16), uint8(msize>>24)
	b[11], b[12] = uint8(len(version)), uint8(len(version)>>8)
	copy(b[13:], version)
	return s
}

// Assumption: something read the first 4 bytes of a message and returned
// a byte slice with that many bytes. For consistency, the message will include
// that length field.
func rversion(b []byte) (version string, msize uint32) {
	msize = uint32(b[7]) | uint32(b[8])<<8 | uint32(b[9])<<16 | uint32(b[10])<<24
	l := uint32(b[11]) | uint32(b[12])<<8
	version = string(b[13 : 13+l])
	return version, msize
}
