package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	A "perpus_api/db"
	M "perpus_api/models"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type TransactionForm struct {
	TransactionType     string                  `json:"transaction_type"`
	MemberId            *int                    `json:"member_id"`
	Date                *string                 `json:"date"`
	TransactionBeforeId *int                    `json:"transaction_before_id"`
	Detail              *string                 `json:"detail"`
	CreatedBy           *int                    `json:"created_by"`
	ApprovalStatus      *string                 `json:"approval_status"`
	ApproverId          *int                    `json:"approver_id"`
	BooksId             *[]int                  `json:"books_id"`
	BooksQty            *map[int]map[string]int `json:"books_qty"`
	PenaltyFee          *int                    `json:"penalty"`
}

var TransactionLock = sync.Mutex{}

func CreateTransaction(c echo.Context) error {

	TransactionLock.Lock()
	defer TransactionLock.Unlock()

	requestBody := new(TransactionForm)

	if err := c.Bind(requestBody); err != nil {
		return err
	}

	if err := c.Validate(requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	claims, errm := M.GetUserDataByJWT(c)

	if errm != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"msg": *errm})
	}

	switch requestBody.TransactionType {
	case "LOAN":
		res, err := CreateLoan(*requestBody, *claims)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, res)
		}

		return c.JSON(http.StatusOK, res)

	case "INVENTORY_IN":
	case "INVENTORY_OUT":
	case "RETURN":
	case "LOST":
	}

	return nil
}

func CreateLoan(data TransactionForm, claims M.JwtCustomClaims) (M.ResponseNoData, error) {

	TransactionLock.Lock()
	defer TransactionLock.Unlock()

	var res M.ResponseNoData

	if data.BooksId == nil {
		res.Status = http.StatusInternalServerError
		res.Msg = "Books cannot be nil!"
		res.Success = false
		return res, nil
	}

	con := A.GetDB()

	sql := `
	INSERT INTO public.transactions(
		transaction_type, member_id, date, detail, created_at, created_by, approval_status, approver_id, expected_return_date)
		VALUES ($1, $2, $3, $4, NOW(), $5, $6, $7, $8) RETURNING id
	`

	t1, err := time.Parse("2006-01-02", *data.Date)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	t2 := t1.AddDate(0, 0, 14) //make this dynamic prob

	var transId *int

	if claims.Role == 0 {
		err := con.QueryRow(sql, data.TransactionType, data.MemberId, t1, data.Detail, claims.ID, "WAITING", nil, t2).Scan(&transId)

		if transId == nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, nil
		}
	} else {
		err := con.QueryRow(sql, data.TransactionType, data.MemberId, t1, data.Detail, claims.ID, "APPROVE", claims.ID, t2).Scan(&transId)

		if transId == nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, nil
		}
	}

	sql = `
	INSERT INTO public.transaction_detail(
		transaction_id, book_id, detail_book_id, created_at, book_condition)
		VALUES ($1, $2, $3, NOW(), $4)
	`

	sqlFind := `SELECT book_id, status, condition FROM book_details WHERE id = $1`

	sqlUpdate := `
		UPDATE public.book_details
		SET status = $1
		WHERE id = $2
	`

	//update book statuses
	for _, elem := range *data.BooksId {
		var book_id *int
		var status *string
		var condition *string

		err := con.QueryRow(sqlFind, elem).Scan(&book_id, &status, &condition)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

		if *status == "REMOVED" || *status == "MISSING" {
			res.Status = http.StatusInternalServerError
			res.Msg = "Book is missing or already removed!"
			res.Success = false
			return res, err
		}

		_, err = con.Exec(sql, transId, book_id, elem, condition)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

		_, err = con.Exec(sqlUpdate, "BORROWED", elem)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}
	}

	res.Status = http.StatusOK
	res.Msg = "Success to create transaction!"
	res.Success = true
	return res, nil

}

func InventoryIn(data TransactionForm, claims M.JwtCustomClaims) (M.ResponseNoData, error) {
	TransactionLock.Lock()
	defer TransactionLock.Unlock()

	var res M.ResponseNoData

	if data.BooksId == nil {
		res.Status = http.StatusInternalServerError
		res.Msg = "Books cannot be nil!"
		res.Success = false
		return res, nil
	}

	con := A.GetDB()

	sql := `
	INSERT INTO public.transactions(
		transaction_type, date, detail, created_at, created_by, approval_status, approver_id)
		VALUES ($1, $2, $3, NOW(), $4, $5, $6) RETURNING id;
	`

	t1, err := time.Parse("2006-01-02", *data.Date)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	var transId *int

	if claims.Role == 0 {
		res.Status = http.StatusUnauthorized
		res.Msg = "Unouthorized!"
		res.Success = false
		return res, nil
	} else {
		err := con.QueryRow(sql, data.TransactionType, t1, data.Detail, claims.ID, "APPROVE", claims.ID).Scan(&transId)

		if transId == nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, nil
		}
	}

	sqlInsertBook := `
	INSERT INTO public.book_details(
		book_id, serial_number, status, condition)
		VALUES ($1, $2, $3, $4);
	`

	sqlInsertTransactionDetail := `
	INSERT INTO public.transaction_detail(
		transaction_id, book_id, book_condition, qty, created_at)
		VALUES ($1, $2, $3, $4, NOW());
	`

	for key, value := range *data.BooksQty {
		currentDate := time.Now().Format("20060102")

		formattedBookID := fmt.Sprintf("%07d", key)

		var condition string

		switch value["condition"] {
		case 0:
			condition = "MINT"
		case 1:
			condition = "FINE"
		case 2:
			condition = "GOOD"
		case 3:
			condition = "FAIR"
		case 4:
			condition = "POOR"
		default:
			condition = "MINT"
		}

		for i := 0; i < value["qty"]; i++ {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			randomNumber := r.Intn(900000) + 100000
			serialNumber := fmt.Sprintf("%s%s%d", currentDate, formattedBookID, randomNumber)

			_, err := con.Exec(sqlInsertBook, key, serialNumber, "STORED", condition)

			if err != nil {
				res.Status = http.StatusInternalServerError
				res.Msg = err.Error()
				res.Success = false
				return res, nil
			}
		}

		_, err := con.Exec(sqlInsertTransactionDetail, transId, key, condition, value["qty"])

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, nil
		}
	}

	res.Status = http.StatusOK
	res.Msg = "Success to create transaction!"
	res.Success = true
	return res, nil
}

func InventoryOut(data TransactionForm, claims M.JwtCustomClaims) (M.ResponseNoData, error) {
	TransactionLock.Lock()
	defer TransactionLock.Unlock()

	var res M.ResponseNoData

	if data.BooksId == nil {
		res.Status = http.StatusInternalServerError
		res.Msg = "Books cannot be nil!"
		res.Success = false
		return res, nil
	}

	con := A.GetDB()

	sql := `
	INSERT INTO public.transactions(
		transaction_type, date, detail, created_at, created_by, approval_status, approver_id)
		VALUES ($1, $2, $3, NOW(), $4, $5, $6) RETURNING id;
	`

	sqlDelete := `
	UPDATE public.book_details
		SET status= $1
		WHERE id = $2
	`

	t1, err := time.Parse("2006-01-02", *data.Date)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	var transId *int

	if claims.Role == 0 {
		res.Status = http.StatusUnauthorized
		res.Msg = "Unouthorized!"
		res.Success = false
		return res, nil
	} else {
		err := con.QueryRow(sql, data.TransactionType, t1, data.Detail, claims.ID, "APPROVE", claims.ID).Scan(&transId)

		if transId == nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, nil
		}
	}

	for _, value := range *data.BooksId {
		_, err := con.Exec(sqlDelete, "REMOVED", value)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, nil
		}
	}

	res.Status = http.StatusOK
	res.Msg = "Success to create transaction!"
	res.Success = true
	return res, nil
}

func CreateReturn(data TransactionForm, claims M.JwtCustomClaims) (M.ResponseNoData, error) {
	TransactionLock.Lock()
	defer TransactionLock.Unlock()

	var res M.ResponseNoData

	if data.BooksId == nil {
		res.Status = http.StatusInternalServerError
		res.Msg = "Books cannot be nil!"
		res.Success = false
		return res, nil
	}

	con := A.GetDB()

	var transId *int

	sqlFind := "select id from public.transactions where id = $1"

	err := con.QueryRow(sqlFind).Scan(&transId)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	sql := `
	INSERT INTO public.transactions(
		transaction_type, member_id, date, transaction_before_id, detail, created_at, created_by, approval_status, approver_id)
		VALUES ($1, $2, $3, $4, $5, NOW(), $6, $7, $8) RETURNING id
	`

	t1, err := time.Parse("2006-01-02", *data.Date)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	if claims.Role == 0 {
		res.Status = http.StatusUnauthorized
		res.Msg = "Unouthorized!"
		res.Success = false
		return res, nil
	} else {
		err := con.QueryRow(sql, data.TransactionType, data.MemberId, t1, data.TransactionBeforeId, data.Detail, claims.ID, "APPROVE", claims.ID).Scan(&transId)

		if transId == nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, nil
		}
	}

	sql = `
	INSERT INTO public.transaction_detail(
		transaction_id, book_id, detail_book_id, created_at, book_condition)
		VALUES ($1, $2, $3, NOW(), $4)
	`

	sqlFind = `SELECT book_id, status, condition FROM book_details WHERE id = $1`

	sqlUpdate := `
		UPDATE public.book_details
		SET status = $1
		WHERE id = $2
	`

		//update book statuses
		for _, elem := range *data.BooksId {
			var book_id *int
			var status *string
			var condition *string
	
			err := con.QueryRow(sqlFind, elem).Scan(&book_id, &status, &condition)
	
			if err != nil {
				res.Status = http.StatusInternalServerError
				res.Msg = err.Error()
				res.Success = false
				return res, err
			}
	
			if *status == "REMOVED" || *status == "MISSING" {
				res.Status = http.StatusInternalServerError
				res.Msg = "Book is missing or already removed!"
				res.Success = false
				return res, err
			}
	
			_, err = con.Exec(sql, transId, book_id, elem, condition)
	
			if err != nil {
				res.Status = http.StatusInternalServerError
				res.Msg = err.Error()
				res.Success = false
				return res, err
			}
	
			_, err = con.Exec(sqlUpdate, "STORED", elem)
	
			if err != nil {
				res.Status = http.StatusInternalServerError
				res.Msg = err.Error()
				res.Success = false
				return res, err
			}
		}

	res.Status = http.StatusOK
	res.Msg = "Success to create transaction!"
	res.Success = true
	return res, nil
}
