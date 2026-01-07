package main

import (
	"net/http"
	"strconv"

	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// Get user wallet
func getWalletHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		wallet, err := store.GetUserWallet(db, uint(user.ID))
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch wallet")
			return
		}

		writeJSON(w, http.StatusOK, wallet)
	}
}

// Get wallet transactions
func getWalletTransactionsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		limit := 50
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		transactions, err := store.GetWalletTransactions(db, uint(user.ID), limit)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch transactions")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"transactions": transactions,
		})
	}
}

// Transfer wallet funds to another user
func transferWalletHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		var req struct {
			ToUserID uint    `json:"to_user_id"`
			Amount   float64 `json:"amount"`
			Currency string  `json:"currency"` // coins or cash
		}

		if err := readJSON(w, r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Amount <= 0 {
			writeJSONError(w, http.StatusBadRequest, "amount must be positive")
			return
		}

		// Get sender wallet
		senderWallet, err := store.GetUserWallet(db, uint(user.ID))
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch wallet")
			return
		}

		// Check balance
		if req.Currency == "coins" {
			if senderWallet.Coins < int(req.Amount) {
				writeJSONError(w, http.StatusBadRequest, "insufficient coins")
				return
			}
		} else if req.Currency == "cash" {
			if senderWallet.Cash < req.Amount {
				writeJSONError(w, http.StatusBadRequest, "insufficient cash")
				return
			}
		} else {
			writeJSONError(w, http.StatusBadRequest, "invalid currency")
			return
		}

		// Get receiver wallet
		receiverWallet, err := store.GetUserWallet(db, req.ToUserID)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, "receiver not found")
			return
		}

		// Perform transfer in transaction
		err = db.Transaction(func(tx *gorm.DB) error {
			// Deduct from sender
			if req.Currency == "coins" {
				senderWallet.Coins -= int(req.Amount)
			} else {
				senderWallet.Cash -= req.Amount
			}
			if err := tx.Save(senderWallet).Error; err != nil {
				return err
			}

			// Add to receiver
			if req.Currency == "coins" {
				receiverWallet.Coins += int(req.Amount)
			} else {
				receiverWallet.Cash += req.Amount
			}
			if err := tx.Save(receiverWallet).Error; err != nil {
				return err
			}

			// Create transactions
			toUserIDUint := uint(req.ToUserID)
			userIDUint := uint(user.ID)
			refType := "transfer"

			senderTx := &store.WalletTransaction{
				UserID:        uint(user.ID),
				Type:          "debit",
				Amount:        req.Amount,
				Currency:      req.Currency,
				Description:   "Transfer to user",
				ReferenceID:    &toUserIDUint,
				ReferenceType: &refType,
			}
			if err := tx.Create(senderTx).Error; err != nil {
				return err
			}

			receiverTx := &store.WalletTransaction{
				UserID:        req.ToUserID,
				Type:          "credit",
				Amount:        req.Amount,
				Currency:      req.Currency,
				Description:   "Transfer from user",
				ReferenceID:    &userIDUint,
				ReferenceType: &refType,
			}
			return tx.Create(receiverTx).Error
		})

		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "transfer failed")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "transfer completed successfully",
		})
	}
}


