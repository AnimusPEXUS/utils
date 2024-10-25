package fsnotify

// #include <sys/inotify.h>
import "C"

const (
	// for Add
	// IN_ATTRI         = C.IN_ATTRI
	IN_ACCESS        uint32 = C.IN_ACCESS
	IN_CLOSE_WRITE   uint32 = C.IN_CLOSE_WRITE
	IN_CLOSE_NOWRITE uint32 = C.IN_CLOSE_NOWRITE
	IN_CREATE        uint32 = C.IN_CREATE
	IN_DELETE        uint32 = C.IN_DELETE
	IN_DELETE_SELF   uint32 = C.IN_DELETE_SELF
	IN_MODIFY        uint32 = C.IN_MODIFY
	IN_MOVE_SELF     uint32 = C.IN_MOVE_SELF
	IN_MOVED_FROM    uint32 = C.IN_MOVED_FROM
	IN_MOVED_TO      uint32 = C.IN_MOVED_TO
	IN_OPEN          uint32 = C.IN_OPEN
	IN_ALL_EVENTS    uint32 = C.IN_ALL_EVENTS
	IN_MOVE          uint32 = C.IN_MOVE
	IN_CLOSE         uint32 = C.IN_CLOSE
	IN_DONT_FOLLOW   uint32 = C.IN_DONT_FOLLOW
	IN_EXCL_UNLINK   uint32 = C.IN_EXCL_UNLINK
	IN_MASK_ADD      uint32 = C.IN_MASK_ADD
	IN_ONESHOT       uint32 = C.IN_ONESHOT
	IN_ONLYDIR       uint32 = C.IN_ONLYDIR
	IN_MASK_CREATE   uint32 = C.IN_MASK_CREATE
	IN_IGNORED       uint32 = C.IN_IGNORED
	IN_ISDIR         uint32 = C.IN_ISDIR
	IN_Q_OVERFLOW    uint32 = C.IN_Q_OVERFLOW
	IN_UNMOUNT       uint32 = C.IN_UNMOUNT

	// for Init1
	IN_CLOEXEC  int = C.IN_CLOEXEC
	IN_NONBLOCK int = C.IN_NONBLOCK
)
