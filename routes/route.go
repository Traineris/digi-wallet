package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Shercosta/digi-wallet/middleware"
	"github.com/Shercosta/digi-wallet/models"
	"github.com/Shercosta/digi-wallet/request"
	"github.com/Shercosta/digi-wallet/response"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type AddBalanceRequest struct {
    Amount int `json:"amount"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	response.JSONSuccess(w, "hello from routes", nil, nil)
}

func GetAllUsers(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        sort := r.URL.Query().Get("sort_amount")
        users := []models.User{}

        query := db
        if sort == "asc" {
            query = query.Order("balance asc")
        } else if sort == "desc" {
            query = query.Order("balance desc")
        }

        if err := query.Find(&users).Error; err != nil {
            response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
            return
        }

        response.JSONSuccess(w, users, nil, nil)
    }
}

func GetBalance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var balance models.Balance

		userID := middleware.GetUserID(r.Context())

		fmt.Println("pass")

		err := db.Where("user_id = ?", userID).First(&balance).Error
		if err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, balance, nil, nil)
	}
}

func InitializeBalance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var balance models.Balance

		err := db.First(&balance).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				balance = models.Balance{Amount: 100000}
				if err := db.Create(&balance).Error; err != nil {
					response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
					return
				}
				response.JSONSuccess(w, "Balance initialized", nil, nil)
				return
			}

			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		balance.Amount = 100000
		if err := db.Save(&balance).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, "Balance reset", nil, nil)
	}
}

func PostTakeBalance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var balance models.Balance

		var req request.TakeRequest
		req.AssignFormValues(r)

		userID := middleware.GetUserID(r.Context())

		err := db.Where("user_id = ?", userID).First(&balance).Error
		if err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		balance.Amount -= *req.Amount
		if err := db.Save(&balance).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, balance, nil, nil)
	}
}

func AddBalance(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var body AddBalanceRequest
        if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
            response.JSONError(w, http.StatusBadRequest, err.Error(), nil)
            return
        }

        userID := middleware.GetUserID(r.Context())

        var balance models.Balance
        if err := db.Where("user_id = ?", userID).First(&balance).Error; err != nil {
            response.JSONError(w, http.StatusNotFound, "balance not found", nil)
            return
        }

 
        balance.Amount += float64(body.Amount)
        
        if err := db.Save(&balance).Error; err != nil {
            response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
            return
        }

        var user models.User
        db.First(&user, userID)
        user.Level = UpdateLevel(int(balance.Amount))
        db.Save(&user)

        response.JSONSuccess(w, balance, nil, nil)
    }
}

func UpdateLevel(balance int) int {
	switch {
	case balance > 3000000:
		return 4
	case balance > 2000000:
		return 3
	case balance > 1000000:
		return 2
	default:
		return 1
	}
}

func DeleteUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetID, _ := strconv.Atoi(chi.URLParam(r, "id"))
		currentUserID := middleware.GetUserID(r.Context())

		var currentUser, targetUser models.User
		if err := db.First(&currentUser, currentUserID).Error; err != nil {
			response.JSONError(w, http.StatusUnauthorized, "user not found", nil)
			return
		}
		if err := db.First(&targetUser, targetID).Error; err != nil {
			response.JSONError(w, http.StatusNotFound, "target user not found", nil)
			return
		}

		if currentUser.Level <= targetUser.Level {
			response.JSONError(w, http.StatusForbidden, "cannot delete user with equal/higher level", nil)
			return
		}

		db.Delete(&targetUser)
		response.JSONSuccess(w, map[string]string{"message": "user deleted"}, nil, nil)
	}
}


