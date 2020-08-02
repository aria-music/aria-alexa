package aria

type Request struct {
	Token string      `json:"token"`
	OP    string      `json:"op"`
	Data  interface{} `json:"data"`
}

type RepeatRequestData struct {
	URI string `json:"uri"`
}
