package restagentdemo

import procedures "TD3/comsoc"

type Request struct {
	Operator string                   `json:"op"`
	Prefs    []procedures.Alternative `json:"prefs"`
}

type Response struct {
	Result int `json:"res"`
}
