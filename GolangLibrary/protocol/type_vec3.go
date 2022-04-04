package protocol

type McVec3 struct {
	X float32 `json:"X"`
	Y float32 `json:"Y"`
	Z float32 `json:"Z"`
}

func Vec3FromSlice(slice []float32) McVec3 {
	return McVec3{slice[0], slice[1], slice[2]}
}
