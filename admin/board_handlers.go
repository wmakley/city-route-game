package admin

import (
	"city-route-game/domain"
	"city-route-game/util"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func BoardsIndexHandler(w http.ResponseWriter, r *http.Request) {
	var boards []domain.Board

	if err := db.Find(&boards).Error; err != nil {
		internalServerError(err, w, r)
		return
	}

	err := ParseAndExecuteAdminTemplate(w, "boards/index", &boards)
	if err != nil {
		panic(err)
	}
}

func NewBoardHandler(w http.ResponseWriter, r *http.Request) {
	data := BoardForm{
		Form: Form{
			Action: "/boards/",
			Method: "POST",
		},
		Name: "",
	}

	err := ParseAndExecuteAdminTemplate(w, "boards/new", &data, "boards/_form")
	if err != nil {
		panic(err)
	}
}

func CreateBoardHandler(w http.ResponseWriter, r *http.Request) {
	var form BoardForm

	err := formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		internalServerError(err, w, r)
		return
	}

	form.NormalizeInputs()

	if !form.IsValid(db) {
		w.WriteHeader(http.StatusBadRequest)
		ParseAndExecuteAdminTemplate(w, "boards/new", &form, "boards/_form")
		return
	}

	board := domain.Board{
		Name:   form.Name,
		Width:  800,
		Height: 500,
	}

	if err := db.Save(&board).Error; err != nil {
		internalServerError(err, w, r)
		return
	}

	util.TurbolinksVisit("/boards", true, w, r)
}

func GetBoardByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	var board domain.Board
	err := db.First(&board, key).Error
	if err != nil {
		handleDBErr(w, r, err)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		util.MustReturnJson(w, board)
	} else {
		err = ParseAndExecuteAdminTemplate(w, "boards/show", &board)
		if err != nil {
			panic(err)
		}
	}
}

type EditBoardPage struct {
	BoardForm BoardForm
	BoardJSON string
}

func EditBoardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	var board domain.Board
	err := db.First(&board, key).Error
	if err != nil {
		handleDBErr(w, r, err)
		return
	}

	boardForm := BoardForm{
		Form: Form{
			Action: "/boards/" + key,
			Method: "PATCH",
		},
		ID:   board.ID,
		Name: board.Name,
	}

	boardJson, err := json.Marshal(board)
	if err != nil {
		panic(err)
	}

	editBoardPage := EditBoardPage{
		BoardForm: boardForm,
		BoardJSON: string(boardJson),
	}

	err = ParseAndExecuteAdminTemplate(w, "boards/edit", &editBoardPage, "boards/_form")
	if err != nil {
		panic(err)
	}
}

func UpdateBoardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	accept := r.Header.Get("Accept")
	gotJson := strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")
	respondWithJson := strings.HasPrefix(accept, "application/json")

	var form BoardForm

	if gotJson {
		err := json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			panic(err)
		}
	} else {
		err := formDecoder.Decode(&form, r.PostForm)
		if err != nil {
			internalServerError(err, w, r)
			return
		}
	}

	form.NormalizeInputs()

	var board domain.Board
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.First(&board, key).Error
		if err != nil {
			return err
		}

		// Set missing fields to original values
		if gotJson {
			// Json is mostly used to update dimensions, so ignore name
			if form.Name == "" {
				form.Name = board.Name
			}
		} else {
			// Html is mostly used to update the name
			if r.FormValue("Width") == "" {
				form.Width = board.Width
			}
			if r.FormValue("Height") == "" {
				form.Height = board.Height
			}
		}

		if !form.IsValid(db) {
			return ErrInvalidForm
		}

		board.Name = form.Name
		board.Width = form.Width
		board.Height = form.Height

		err = tx.Save(&board).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, ErrInvalidForm) {
			if respondWithJson {
				json := make(map[string]interface{})
				json["board"] = board
				json["errors"] = form.Errors()

				util.SetJSONContentType(w)
				w.WriteHeader(http.StatusBadRequest)
				util.MustEncode(w, json)
			} else {
				boardJson, err := json.Marshal(board)
				if err != nil {
					panic(err)
				}

				invalidEditBoardPage := EditBoardPage{
					BoardForm: form,
					BoardJSON: string(boardJson),
				}

				util.SetHTMLContentType(w)
				w.WriteHeader(http.StatusBadRequest)
				err = ParseAndExecuteAdminTemplate(w, "boards/edit", &invalidEditBoardPage, "boards/_form")
				if err != nil {
					panic(err)
				}
			}
		} else {
			handleDBErr(w, r, err)
		}
		return
	}

	if respondWithJson {
		util.SetJSONContentType(w)
		util.MustEncode(w, &board)
	} else {
		// Call a global function in the admin js directly
		util.SetJavaScriptContentType(w)
		fmt.Fprintf(w, ";updateFormSucceeded();")
	}
}

func DeleteBoardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		log.Println("Bad ID:", vars["id"])
		http.NotFound(w, r)
		return
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return domain.DeleteBoard(tx, uint(id))
	})

	if err != nil {
		handleDBErr(w, r, err)
		return
	}

	util.TurbolinksVisit("/boards", true, w, r)
}
