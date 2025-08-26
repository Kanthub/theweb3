package routes

import (
	"fmt"

	"net/http"
)

func (h Routes) GetSupportSymbolList(w http.ResponseWriter, r *http.Request) {
	supRet, err := h.srv.GetSupportSymbolList()
	if err != nil {
		return
	}
	err = jsonResponse(w, supRet, http.StatusOK)
	if err != nil {
		fmt.Println("Error writing response", "err", err.Error())
	}
}

func (h Routes) GetMarketPrice(w http.ResponseWriter, r *http.Request) {
	addrRet, err := h.srv.GetMarketPrice()
	if err != nil {
		return
	}
	err = jsonResponse(w, addrRet, http.StatusOK)
	if err != nil {
		fmt.Println("Error writing response", "err", err.Error())
	}
}
