package admin

import (
	"city-route-game/internal/app"
	"city-route-game/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type BoardController struct {
	Controller
	boardEditorService app.BoardEditorService
}

func NewBoardController(config ControllerConfig, service app.BoardEditorService) BoardController {
	return BoardController{
		Controller: Controller{
			FormDecoder:  config.FormDecoder,
			TemplateRoot: config.TemplateRoot,
			AssetHost:    config.AssetHost,
		},
		boardEditorService: service,
	}
}

func (c BoardController)Index(w http.ResponseWriter, r *http.Request) {
	boards, err := c.boardEditorService.FindAll(r.Context())
	if err != nil {
		panic(err)
	}

	page := NewPageWithData(c.AssetHost, boards)

	if err = c.ParseAndExecuteAdminTemplate(w, "boards/index", &page); err != nil {
		panic(err)
	}
}

func (c BoardController)New(w http.ResponseWriter, r *http.Request) {
	data := app.NewCreateBoardForm()
	page := NewPageWithData(c.AssetHost, &data)

	err := c.ParseAndExecuteAdminTemplate(w, "boards/new", &page, "boards/_form")
	if err != nil {
		panic(err)
	}
}

func (c BoardController)Create(w http.ResponseWriter, r *http.Request) {
	var form app.CreateBoardForm
	var err error

	if err = c.FormDecoder.Decode(&form, r.PostForm); err != nil {
		c.InternalServerError(err, w, r)
		return
	}

	_, err = c.boardEditorService.CreateBoard(r.Context(), &form)

	if err != nil {
		if errors.Is(err, app.ErrInvalidForm) {
			page := NewPageWithData(c.AssetHost, &form)
			w.WriteHeader(http.StatusBadRequest)
			if err = c.ParseAndExecuteAdminTemplate(w, "boards/new", &page, "boards/_form"); err != nil {
				panic(err)
			}
			return
		} else {
			c.InternalServerError(err, w, r)
			return
		}
	}

	util.TurbolinksVisit("/boards", true, w, r)
}

type idParam struct {
	ID app.ID
}

func (c BoardController)GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	board, err := c.boardEditorService.FindByID(r.Context(), id)
	if err != nil {
		c.HandleServiceError(err, w, r)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		util.MustReturnJson(w, board)
	} else {
		page := NewPageWithData(c.AssetHost, &board)
		err = c.ParseAndExecuteAdminTemplate(w, "boards/show", &page)
		if err != nil {
			panic(err)
		}
	}
}

type EditBoardPage struct {
	BoardForm *app.UpdateBoardForm
	BoardJSON string
}

func (c BoardController)Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	board, err := c.boardEditorService.FindByID(r.Context(), id)
	if err != nil {
		c.HandleServiceError(err, w, r)
		return
	}

	boardForm := app.NewUpdateBoardForm(board)

	boardJson, err := json.Marshal(board)
	if err != nil {
		panic(err)
	}

	page := NewPageWithData(
		c.AssetHost,
			EditBoardPage{
			BoardForm: &boardForm,
			BoardJSON: string(boardJson),
		})

	err = c.ParseAndExecuteAdminTemplate(w, "boards/edit", &page, "boards/_form")
	if err != nil {
		panic(err)
	}
}

func (c BoardController)Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	accept := r.Header.Get("Accept")
	gotJson := strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")
	respondWithJson := strings.HasPrefix(accept, "application/json")

	var form app.UpdateBoardForm
	var gotName bool
	var gotDimensions bool

	if gotJson {
		err := json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			panic(err)
		}
		// assume JSON contains everything
		gotName = true
		gotDimensions = true
	} else {
		_, gotName = r.PostForm["name"]
		_, gotDimensions = r.PostForm["width"]
		err := c.FormDecoder.Decode(&form, r.PostForm)
		if err != nil {
			panic(err)
		}
	}

	var err error
	var board *app.Board

	if gotName && gotDimensions {
		board, err = c.boardEditorService.Update(r.Context(), id, &form)
	} else if !gotName && gotDimensions {
		board, err = c.boardEditorService.UpdateDimensions(r.Context(), id, &form)
	} else {
		board, err = c.boardEditorService.UpdateName(r.Context(), id, &form)
	}

	if err != nil {
		if errors.Is(err, app.ErrInvalidForm) {
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

				page := NewPageWithData(c.AssetHost, EditBoardPage{
					BoardForm: &form,
					BoardJSON: string(boardJson),
				})

				util.SetHTMLContentType(w)
				w.WriteHeader(http.StatusBadRequest)
				err = c.ParseAndExecuteAdminTemplate(w, "boards/edit", &page, "boards/_form")
				if err != nil {
					panic(err)
				}
			}
		} else {
			c.HandleServiceError(err, w, r)
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

func (c BoardController)Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := c.boardEditorService.DeleteByID(r.Context(), id); err != nil {
		c.HandleServiceError(err, w, r)
		return
	}

	util.TurbolinksVisit("/boards", true, w, r)
}
