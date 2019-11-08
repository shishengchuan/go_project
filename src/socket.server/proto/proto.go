package proto

type ProtoBase struct{
	ProtoName string  `json:"protoName"`
}

type Fire struct{
	ProtoName string `json:"protoName"`	
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}
