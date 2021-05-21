package admin

import (
	"city-route-game/domain"
	"city-route-game/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func BoardsIndexHandler(w http.ResponseWriter, r *http.Request) {
	boards, err := domain.ListBoards(domain.InitialContext())
	if err != nil {
		panic(err)
	}

	page := NewPageWithData(boards)

	if err = ParseAndExecuteAdminTemplate(w, "boards/index", &page); err != nil {
		panic(err)
	}
}

func NewBoardHandler(w http.ResponseWriter, r *http.Request) {
	data := domain.NewBoardForm()
	page := NewPageWithData(&data)

	err := ParseAndExecuteAdminTemplate(w, "boards/new", &page, "boards/_form")
	if err != nil {
		panic(err)
	}
}

func CreateBoardHandler(w http.ResponseWriter, r *http.Request) {
	var form domain.BoardForm
	var err error

	if err = formDecoder.Decode(&form, r.PostForm); err != nil {
		internalServerError(err, w, r)
		return
	}

	_, err = domain.CreateBoard(domain.InitialContext(), &form)

	if err != nil {
		if errors.Is(err, domain.ErrInvalidForm) {
			page := NewPageWithData(&form)
			w.WriteHeader(http.StatusBadRequest)
			if err = ParseAndExecuteAdminTemplate(w, "boards/new", &page, "boards/_form"); err != nil {
				panic(err)
			}
			return
		} else {
			internalServerError(err, w, r)
			return
		}
	}

	util.TurbolinksVisit("/boards", true, w, r)
}

func GetBoardByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	board, err := domain.GetBoardByID(domain.InitialContext(), key)
	if err != nil {
		handleDBErr(w, r, err)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		util.MustReturnJson(w, board)
	} else {
		page := NewPageWithData(&board)
		err = ParseAndExecuteAdminTemplate(w, "boards/show", &page)
		if err != nil {
			panic(err)
		}
	}
}

type EditBoardPage struct {
	BoardForm *domain.BoardForm
	BoardJSON string
}

func EditBoardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	board, err := domain.GetBoardByID(domain.InitialContext(), key)
	if err != nil {
		handleDBErr(w, r, err)
		return
	}

	boardForm := domain.NewEditBoardForm(board)

	boardJson, err := json.Marshal(board)
	if err != nil {
		panic(err)
	}

	page := NewPageWithData(EditBoardPage{
		BoardForm: &boardForm,
		BoardJSON: string(boardJson),
	})

	err = ParseAndExecuteAdminTemplate(w, "boards/edit", &page, "boards/_form")
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

	var form domain.BoardForm

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

	var err error
	var board *domain.Board

	if gotJson && form.Name == "" {
		board, err = domain.UpdateBoardDimensions(domain.InitialContext(), key, &form)
	} else if form.Width == 0 && form.Height == 0 {
		board, err = domain.UpdateBoardName(domain.InitialContext(), key, &form)
	} else {
		board, err = domain.UpdateBoard(domain.InitialContext(), key, &form)
	}

	if err != nil {
		if errors.Is(err, domain.ErrInvalidForm) {
			if respondWithJson {
				body := make(map[string]interface{})
				body["board"] = board
				body["errors"] = form.Errors

				util.SetJSONContentType(w)
				w.WriteHeader(http.StatusBadRequest)
				util.MustEncode(w, body)
			} else {
				boardJson, err := json.Marshal(board)
				if err != nil {
					panic(err)
				}

				page := NewPageWithData(EditBoardPage{
					BoardForm: &form,
					BoardJSON: string(boardJson),
				})

				util.SetHTMLContentType(w)
				w.WriteHeader(http.StatusBadRequest)
				err = ParseAndExecuteAdminTemplate(w, "boards/edit", &page, "boards/_form")
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
		util.MustEncode(w, board)
	} else {
		// Call a global function in the admin js directly
		util.SetJavaScriptContentType(w)
		_, err := fmt.Fprint(w, ";updateFormSucceeded();")
		if err != nil {
			panic(err)
		}
	}
}

func DeleteBoardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := domain.DeleteBoardById(id); err != nil {
		handleDBErr(w, r, err)
		return
	}

	util.TurbolinksVisit("/boards", true, w, r)
}
