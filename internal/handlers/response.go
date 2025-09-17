package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, r *http.Request, data any) {
	w.Header().Set("Content-Type", "application/json")
	// data_byte, err := json.Marshal(data)
	// if err != nil {
	// 	fmt.Println("error while marshal data - ", err)
	// 	return
	// }
	// w.Write(data_byte)

	enc_json := json.NewEncoder(w)
	enc_json.SetEscapeHTML(false) // Don’t change &, <, or > to /u0026,/u003c,/u003e ,just keep them as they are in my JSON
	err := enc_json.Encode(data)
	if err != nil {
		fmt.Println("error while encoding data - ", err)
		return
	}

}
