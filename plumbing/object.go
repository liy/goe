package plumbing

type ObjectType int8

const (
	OBJ_INVALID ObjectType = 0
	OBJ_COMMIT  ObjectType = 1
	OBJ_TREE    ObjectType = 2
	OBJ_BLOB    ObjectType = 3
	OBJ_TAG     ObjectType = 4
	// 5 is reserved for future expansion
	OBJ_OFS_DELTA ObjectType = 6
	OBJ_REF_DELTA ObjectType = 7
)

type Hash [20]byte
