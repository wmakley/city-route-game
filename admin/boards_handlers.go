package admin

import (
	"city-route-game/domain"
	"city-route-game/util"
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
	data := NewBoardForm{
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
	var form NewBoardForm

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
	err := db.Find(&board, key).Error
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
	err := db.Find(&board, key).Error
	if err != nil {
		handleDBErr(w, r, err)
		return
	}

	form := EditBoardForm{
		Form: Form{
			Action: "/boards/" + key,
			Method: "PATCH",
		},
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

	var form NewBoardForm

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

		err := tx.Find(&board, key).Error
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

	util.TurbolinksVisit("/boards/"+key, true, w, r)
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
