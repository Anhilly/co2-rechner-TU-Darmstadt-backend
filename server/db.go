package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
)

func RouteDB() chi.Router {
	r := chi.NewRouter()

	r.Post("/addFaktor", PostAddFaktor)
	r.Post("/addZaehlerdaten", PostAddZaehlerdaten)
	r.Post("/insertZaehler", PostInsertZaehler)
	r.Post("/insertGebaeude", PostInsertGebaeude)

	return r
}

func PostAddFaktor(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	data := structs.AddCO2Faktor{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	//Datenerarbeitung
	ordner, err := database.CreateDump("PostAddFaktor")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	err = database.EnergieversorgungAddFaktor(data)
	if err != nil {
		err := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt

		res.WriteHeader(http.StatusInternalServerError)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	//Response
	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   nil,
		Error:  nil,
	})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
}

func PostAddZaehlerdaten(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	data := structs.AddZaehlerdaten{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	//Datenerarbeitung
	ordner, err := database.CreateDump("PostAddFaktor")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	err = database.ZaehlerAddZaehlerdaten(data)
	if err != nil {
		err := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt

		res.WriteHeader(http.StatusInternalServerError)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	//Response
	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   nil,
		Error:  nil,
	})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
}

func PostInsertZaehler(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	data := structs.InsertZaehler{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	//Datenerarbeitung
	ordner, err := database.CreateDump("PostAddFaktor")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	err = database.ZaehlerInsert(data)
	if err != nil {
		err := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt

		res.WriteHeader(http.StatusInternalServerError)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	//Response
	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   nil,
		Error:  nil,
	})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
}

func PostInsertGebaeude(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	data := structs.InsertGebaeude{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	//Datenerarbeitung
	ordner, err := database.CreateDump("PostAddFaktor")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	err = database.GebaeudeInsert(data)
	if err != nil {
		err := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt

		res.WriteHeader(http.StatusInternalServerError)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	//Response
	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   nil,
		Error:  nil,
	})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		response, _ := json.Marshal(structs.Response{
			Status: structs.ResponseError,
			Data:   nil,
			Error: structs.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		_, _ = res.Write(response)

		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
}
