package controller

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net/http"

	"github.com/pancakem/rides/v1/src/pkg/common"
	"github.com/pancakem/rides/v1/src/pkg/model"
)

// GetRiders a json encoded list of drivers
func GetRiders(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	users := model.GetAllRider()
	dat := &bytes.Buffer{}
	err := binary.Write(dat, binary.BigEndian, users)
	common.CheckError(err)
	en := json.NewEncoder(w)
	en.Encode(dat)
}
