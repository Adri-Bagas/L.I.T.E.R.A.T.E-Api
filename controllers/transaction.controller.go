package controllers

import (
	"net/http"
	A "perpus_api/db"
	M "perpus_api/models"
	"time"

	"github.com/labstack/echo/v4"
)

type TransactionForm struct {
	TransactionType     string  `json:"transaction_type"`
	MemberId            *int    `json:"member_id"`
	Date                *string `json:"date"`
	TransactionBeforeId *int    `json:"transaction_before_id"`
	Detail              *string `json:"detail"`
	CreatedBy           *int    `json:"created_by"`
	ApprovalStatus      *string `json:"approval_status"`
	ApproverId          *int    `json:"approver_id"`
	BooksId             *[]int  `json:"books_id"`
}

func CreateTransaction(c echo.Context) error {
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
		res, err := CreateLoan(*requestBody, claims)

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

func CreateLoan(data TransactionForm, claims *M.JwtCustomClaims) (M.ResponseNoData, error) {
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

