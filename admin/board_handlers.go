package admin

import (
	"city-route-game/domain"
	"city-route-game/util"
	"fmt"
	"net/http"

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
			Action: "/boards",
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
		ParseAndExecuteAdminTemplate(w, "boards/new", &form, "board/_form")
		return
	}

	board := domain.Board{Name: form.Name}

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

	err = ParseAndExecuteAdminTemplate(w, "boards/show", &board)
	if err != nil {
		panic(err)
	}
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

	form := BoardForm{
		Form: Form{
			Action: "/boards/" + key,
			Method: "PATCH",
		},
		ID:   &board.ID,
		Name: board.Name,
	}

	err = ParseAndExecuteAdminTemplate(w, "boards/edit", &form, "boards/_form")
	if err != nil {
		panic(err)
	}
}

func UpdateBoardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	var form BoardForm

	err := formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		internalServerError(err, w, r)
		return
	}

	form.NormalizeInputs()

	if !form.IsValid(db) {
		w.WriteHeader(http.StatusBadRequest)
		ParseAndExecuteAdminTemplate(w, "boards/edit", &form, "boards/_form")
		return
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		var board domain.Board

		err := tx.First(&board, key).Error
		if err != nil {
			return err
		}

		board.Name = form.Name

		err = tx.Save(&board).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		handleDBErr(w, r, err)
		return
	}

	// Hide the form
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	fmt.Fprintf(w, ";updateFormSucceeded();")
}

func DeleteBoardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	err := db.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&domain.Board{}, key)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})

	if err != nil {
		handleDBErr(w, r, err)
		return
	}

	util.TurbolinksVisit("/boards", true, w, r)
}
