package response

type Success struct {
	Data string `json:"data"`
	Msg  string `json:"msg"`
}
type Err struct {
	Err string `json:"error"`
}
