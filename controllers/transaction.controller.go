package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	A "perpus_api/db"
	M "perpus_api/models"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var TransactionLock = sync.Mutex{}

func GetTransactionInOutDataAll(c echo.Context) error {

	var obj M.TransactionInventoryInOut
	var arrobj []M.TransactionInventoryInOut
	var res M.ResponseMultiple

	con := A.GetDB()

	trans_type := c.FormValue("type")

	sql := "select t1.id as trans_id, t1.date, t1.detail, t1.created_at, t2.* from public.transactions t1 inner join public.users t2 on t1.approver_id = t2.id where transaction_type = '" + trans_type + "';"

	rows, err := con.Query(sql)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return c.JSON(http.StatusInternalServerError, res)
	}
	defer rows.Close()

	for rows.Next() {

		var user M.User

		err := rows.Scan(
			&obj.Id,
			&obj.Date,
			&obj.Detail,
			&obj.CreatedAt,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.LastActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.CreatedBy,
			&user.UpdatedBy,
			&user.DeletedAt,
			&user.DeletedBy,
			&user.Role,
		)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return c.JSON(http.StatusInternalServerError, res)
		}

		*user.Password = ""

		obj.Approver = user

		arrobj = append(arrobj, obj)

	}

	res.Status = http.StatusOK
	res.Msg = "Founded Transaction!"
	res.Success = true
	res.Datas = arrobj

	return c.JSON(http.StatusOK, res)

}

func GetTransactionInOutDataAllObj(c echo.Context, trans_type string) ([]M.TransactionInventoryInOut, error) {

	var obj M.TransactionInventoryInOut
	var arrobj []M.TransactionInventoryInOut

	con := A.GetDB()

	sql := "select t1.id as trans_id, t1.date, t1.detail, t1.created_at, t2.* from public.transactions t1 inner join public.users t2 on t1.approver_id = t2.id where transaction_type = '" + trans_type + "';"

	rows, err := con.Query(sql)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var user M.User

		err := rows.Scan(
			&obj.Id,
			&obj.Date,
			&obj.Detail,
			&obj.CreatedAt,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.LastActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.CreatedBy,
			&user.UpdatedBy,
			&user.DeletedAt,
			&user.DeletedBy,
			&user.Role,
		)

		if err != nil {
			return nil, err
		}

		*user.Password = ""

		obj.Approver = user

		arrobj = append(arrobj, obj)

	}

	return arrobj, nil

}


func GetLoanDataAll(c echo.Context) error {

	var obj M.TransactionLoan
	var arrobj []M.TransactionLoan
	var res M.ResponseMultiple

	con := A.GetDB()

	claims, errm := M.GetUserDataByJWT(c)

	if errm != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"msg": *errm})
	}

	var sql string

	if claims.Role == 0 {
		sql = `
	select 
		t1.id as trans_id, 
		t1.date, 
		t1.detail, 
		t1.created_at, 
		t1.expected_return_date,
		t1.approval_status,
		t1.is_returned,
		t3.*, t2.* 
	from public.transactions t1 
	left join public.users t2 on t1.approver_id = t2.id 
	left join public.members t3 on t1.member_id = t3.id 
	where transaction_type = 'LOAN' and member_id = `+strconv.Itoa(claims.ID)+`;
	`
	} else {
		sql = `
	select 
		t1.id as trans_id, 
		t1.date, 
		t1.detail, 
		t1.created_at, 
		t1.expected_return_date,
		t1.approval_status,
		t1.is_returned,
		t3.*, t2.* 
	from public.transactions t1 
	left join public.users t2 on t1.approver_id = t2.id 
	left join public.members t3 on t1.member_id = t3.id 
	where transaction_type = 'LOAN';
	`
	}

	rows, err := con.Query(sql)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return c.JSON(http.StatusInternalServerError, res)
	}
	defer rows.Close()

	for rows.Next() {

		var user M.User
		var member M.Member

		err := rows.Scan(
			&obj.Id,
			&obj.Date,
			&obj.Detail,
			&obj.CreatedAt,
			&obj.ExpectedReturnDate,
			&obj.ApprovalStatus,
			&obj.IsReturned,
			&member.ID,
			&member.Username,
			&member.FullName,
			&member.Email,
			&member.Password,
			&member.PhoneNumber,
			&member.Address,
			&member.CreatedAt,
			&member.UpdatedAt,
			&member.DeletedAt,
			&member.LastActive,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.LastActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.CreatedBy,
			&user.UpdatedBy,
			&user.DeletedAt,
			&user.DeletedBy,
			&user.Role,
		)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return c.JSON(http.StatusInternalServerError, res)
		}

		if user.Password != nil {
			*user.Password = ""
			*user.Role = 404
		}

		member.Password = ""

		obj.Approver = &user
		obj.Member = &member

		arrobj = append(arrobj, obj)

	}

	res.Status = http.StatusOK
	res.Msg = "Founded Transaction!"
	res.Success = true
	res.Datas = arrobj

	return c.JSON(http.StatusOK, res)
}

func GetLoanDataAllObj(c echo.Context) ([]M.TransactionLoan, error) {

	var obj M.TransactionLoan
	var arrobj []M.TransactionLoan

	con := A.GetDB()

	claims, errm := M.GetUserDataByJWT(c)

	if errm != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"msg": *errm})
	}

	var sql string

	if claims.Role == 0 {
		sql = `
	select 
		t1.id as trans_id, 
		t1.date, 
		t1.detail, 
		t1.created_at, 
		t1.expected_return_date,
		t1.approval_status,
		t1.is_returned,
		t3.*, t2.* 
	from public.transactions t1 
	left join public.users t2 on t1.approver_id = t2.id 
	left join public.members t3 on t1.member_id = t3.id 
	where transaction_type = 'LOAN' and member_id = `+strconv.Itoa(claims.ID)+`;
	`
	} else {
		sql = `
	select 
		t1.id as trans_id, 
		t1.date, 
		t1.detail, 
		t1.created_at, 
		t1.expected_return_date,
		t1.approval_status,
		t1.is_returned,
		t3.*, t2.* 
	from public.transactions t1 
	left join public.users t2 on t1.approver_id = t2.id 
	left join public.members t3 on t1.member_id = t3.id 
	where transaction_type = 'LOAN';
	`
	}

	rows, err := con.Query(sql)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var user M.User
		var member M.Member

		err := rows.Scan(
			&obj.Id,
			&obj.Date,
			&obj.Detail,
			&obj.CreatedAt,
			&obj.ExpectedReturnDate,
			&obj.ApprovalStatus,
			&obj.IsReturned,
			&member.ID,
			&member.Username,
			&member.FullName,
			&member.Email,
			&member.Password,
			&member.PhoneNumber,
			&member.Address,
			&member.CreatedAt,
			&member.UpdatedAt,
			&member.DeletedAt,
			&member.LastActive,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.LastActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.CreatedBy,
			&user.UpdatedBy,
			&user.DeletedAt,
			&user.DeletedBy,
			&user.Role,
		)

		if err != nil {
			return nil, err
		}

		if user.Password != nil {
			*user.Password = ""
			*user.Role = 404
		}

		member.Password = ""

		obj.Approver = &user
		obj.Member = &member

		arrobj = append(arrobj, obj)

	}

	return arrobj, nil
}

func GetReturnDataAll(c echo.Context) error {

	var obj M.TransactionLoan
	var arrobj []M.TransactionLoan
	var res M.ResponseMultiple

	con := A.GetDB()

	sql := `
	select 
		t1.id as trans_id, 
		t1.date, 
		t1.detail, 
		t1.created_at, 
		t1.penalty_fee,
		t3.*, t2.* 
	from public.transactions t1 
		left join public.users t2 on t1.approver_id = t2.id 
		left join public.members t3 on t1.member_id = t3.id 
	where transaction_type = 'RETURN';
	`

	rows, err := con.Query(sql)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return c.JSON(http.StatusInternalServerError, res)
	}
	defer rows.Close()

	for rows.Next() {

		var user M.User
		var member M.Member

		err := rows.Scan(
			&obj.Id,
			&obj.Date,
			&obj.Detail,
			&obj.CreatedAt,
			&obj.PenaltyFee,
			&member.ID,
			&member.Username,
			&member.FullName,
			&member.Email,
			&member.Password,
			&member.PhoneNumber,
			&member.Address,
			&member.CreatedAt,
			&member.UpdatedAt,
			&member.DeletedAt,
			&member.LastActive,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.LastActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.CreatedBy,
			&user.UpdatedBy,
			&user.DeletedAt,
			&user.DeletedBy,
			&user.Role,
		)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return c.JSON(http.StatusInternalServerError, res)
		}

		if user.Password != nil {
			*user.Password = ""
			*user.Role = 404
		}
		member.Password = ""

		obj.Approver = &user
		obj.Member = &member

		arrobj = append(arrobj, obj)

	}

	res.Status = http.StatusOK
	res.Msg = "Founded Transaction!"
	res.Success = true
	res.Datas = arrobj

	return c.JSON(http.StatusOK, res)
}

func GetReturnDataAllObj(c echo.Context) ([]M.TransactionLoan, error) {

	var obj M.TransactionLoan
	var arrobj []M.TransactionLoan

	con := A.GetDB()

	sql := `
	select 
		t1.id as trans_id, 
		t1.date, 
		t1.detail, 
		t1.created_at, 
		t1.penalty_fee,
		t3.*, t2.* 
	from public.transactions t1 
		left join public.users t2 on t1.approver_id = t2.id 
		left join public.members t3 on t1.member_id = t3.id 
	where transaction_type = 'RETURN';
	`

	rows, err := con.Query(sql)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var user M.User
		var member M.Member

		err := rows.Scan(
			&obj.Id,
			&obj.Date,
			&obj.Detail,
			&obj.CreatedAt,
			&obj.PenaltyFee,
			&member.ID,
			&member.Username,
			&member.FullName,
			&member.Email,
			&member.Password,
			&member.PhoneNumber,
			&member.Address,
			&member.CreatedAt,
			&member.UpdatedAt,
			&member.DeletedAt,
			&member.LastActive,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.LastActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.CreatedBy,
			&user.UpdatedBy,
			&user.DeletedAt,
			&user.DeletedBy,
			&user.Role,
		)

		if err != nil {
			return nil, err
		}

		if user.Password != nil {
			*user.Password = ""
			*user.Role = 404
		}
		member.Password = ""

		obj.Approver = &user
		obj.Member = &member

		arrobj = append(arrobj, obj)

	}

	return arrobj, nil
}

func GetLoanDataAllIded(c echo.Context) error {

	var obj M.TransactionLoan
	var arrobj []M.TransactionLoan
	var res M.ResponseMultiple

	con := A.GetDB()

	sql := `
	select 
	t1.id as trans_id, 
	t1.date, 
	t1.detail, 
	t1.created_at, 
	t1.expected_return_date,
	t1.approval_status,
	t1.is_returned,
	t2.*
	from 
	public.transactions t1 
	inner join public.members t2 on t1.member_id = t2.id
	where transaction_type = 'LOAN' AND approval_status = 'APPROVE' and t1.is_returned = false;
	`

	rows, err := con.Query(sql)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return c.JSON(http.StatusInternalServerError, res)
	}
	defer rows.Close()

	for rows.Next() {

		var member M.Member

		err := rows.Scan(
			&obj.Id,
			&obj.Date,
			&obj.Detail,
			&obj.CreatedAt,
			&obj.ExpectedReturnDate,
			&obj.ApprovalStatus,
			&obj.IsReturned,
			&member.ID,
			&member.Username,
			&member.FullName,
			&member.Email,
			&member.Password,
			&member.PhoneNumber,
			&member.Address,
			&member.CreatedAt,
			&member.UpdatedAt,
			&member.DeletedAt,
			&member.LastActive,
		)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return c.JSON(http.StatusInternalServerError, res)
		}

		obj.Member = &member

		arrobj = append(arrobj, obj)

	}

	datas := make(map[int]string)

	for _, val := range arrobj {

		time, err := time.Parse("2006-01-02T00:00:00Z", *val.Date)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return c.JSON(http.StatusInternalServerError, res)
		}

		datas[int(val.Id)] = time.Format("2006-01-02") + " | " + val.Member.FullName
	}

	return c.JSON(http.StatusOK, datas)
}

func GetFindLoanDataWithBooks(c echo.Context) error {
	var obj M.TransactionLoan
	var res M.Response

	con := A.GetDB()

	transId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return c.JSON(http.StatusInternalServerError, res)
	}

	sql := `
	select 
		t1.id as trans_id, 
		t1.date, 
		t1.detail, 
		t1.created_at, 
		t1.expected_return_date,
		t1.approval_status,
		t1.is_returned,
		t1.member_id,
		t3.*, t2.* 
	from public.transactions t1 
		inner join public.users t2 on t1.approver_id = t2.id 
		inner join public.members t3 on t1.member_id = t3.id 
	where transaction_type = 'LOAN' AND t1.id = $1;
	`
	var user M.User
	var member M.Member

	err = con.QueryRow(sql, transId).Scan(
		&obj.Id,
		&obj.Date,
		&obj.Detail,
		&obj.CreatedAt,
		&obj.ExpectedReturnDate,
		&obj.ApprovalStatus,
		&obj.IsReturned,
		&obj.MemberId,
		&member.ID,
		&member.Username,
		&member.FullName,
		&member.Email,
		&member.Password,
		&member.PhoneNumber,
		&member.Address,
		&member.CreatedAt,
		&member.UpdatedAt,
		&member.DeletedAt,
		&member.LastActive,
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.LastActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
		&user.DeletedAt,
		&user.DeletedBy,
		&user.Role,
	)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return c.JSON(http.StatusInternalServerError, res)
	}

	*user.Password = ""
	*user.Role = 404
	member.Password = ""

	obj.Approver = &user
	obj.Member = &member

	sql = `
	select 
		t3.id,
		t2.title,
		t3.serial_number,
		t3.condition,
		t3.status,
		t2.price
	from public.transaction_detail t1 
		inner join books t2 on t2.id = t1.book_id
		inner join book_details t3 on t3.id = t1.detail_book_id
	where t1.transaction_id = $1;
	`

	var book M.BookDetails
	var books []M.BookDetails

	rows, err := con.Query(sql, transId)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return c.JSON(http.StatusInternalServerError, res)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&book.Id,
			&book.Title,
			&book.SerialNumber,
			&book.Condition,
			&book.Status,
			&book.Price,
		)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return c.JSON(http.StatusInternalServerError, res)
		}

		books = append(books, book)
	}

	obj.Books = &books

	res.Status = http.StatusOK
	res.Msg = "Founded Transaction!"
	res.Success = true
	res.Data = obj

	return c.JSON(http.StatusOK, res)
}

func CreateTransaction(c echo.Context) error {

	requestBody := new(M.TransactionForm)

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
	case "LOAN_ONLINE":
		res, err := CreateOnlineLoan(*requestBody, *claims)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, res)
		}

		return c.JSON(http.StatusOK, res)

	case "RETURN_ONLINE":
		res, err := CreateOnlineReturn(*requestBody, *claims)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, res)
		}

		return c.JSON(http.StatusOK, res)

	case "INVENTORY_IN":
		res, err := InventoryIn(*requestBody, *claims)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, res)
		}

		return c.JSON(http.StatusOK, res)
	case "INVENTORY_OUT":
		res, err := InventoryOut(*requestBody, *claims)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, res)
		}

		return c.JSON(http.StatusOK, res)
	case "RETURN":
		res, err := CreateReturn(*requestBody, *claims)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, res)
		}

		return c.JSON(http.StatusOK, res)
	case "LOST":
	}

	return nil
}

func CreateLoan(data M.TransactionForm, claims M.JwtCustomClaims) (M.ResponseNoData, error) {

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

	// tx, err := con.Begin()

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

		// err = M.UpdateBookStock(*book_id)

		// if err != nil {
		// 	res.Status = http.StatusInternalServerError
		// 	res.Msg = err.Error()
		// 	res.Success = false
		// 	return res, err
		// }
	}

	res.Status = http.StatusOK
	res.Msg = "Success to create transaction!"
	res.Success = true
	return res, nil

}

func InventoryIn(data M.TransactionForm, claims M.JwtCustomClaims) (M.ResponseNoData, error) {
	TransactionLock.Lock()
	defer TransactionLock.Unlock()

	var res M.ResponseNoData

	if data.BooksQty == nil {
		res.Status = http.StatusInternalServerError
		res.Msg = "Books cannot be nil!"
		res.Success = false
		return res, nil
	}

	con := A.GetDB()

	// tx, err := con.Begin()

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

		// err = M.UpdateBookStock(key)

		// if err != nil {
		// 	res.Status = http.StatusInternalServerError
		// 	res.Msg = err.Error()
		// 	res.Success = false
		// 	return res, err
		// }
	}

	res.Status = http.StatusOK
	res.Msg = "Success to create transaction!"
	res.Success = true
	return res, nil
}

func InventoryOut(data M.TransactionForm, claims M.JwtCustomClaims) (M.ResponseNoData, error) {
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

	// tx, err := con.Begin()

	sql := `
	INSERT INTO public.transactions(
		transaction_type, date, detail, created_at, created_by, approval_status, approver_id)
		VALUES ($1, $2, $3, NOW(), $4, $5, $6) RETURNING id;
	`

	sqlInsertTransactionDetail := `
	INSERT INTO public.transaction_detail(
		transaction_id, detail_book_id, created_at)
		VALUES ($1, $2, NOW());
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
			return res, err
		}
	}

	for _, value := range *data.BooksId {
		_, err := con.Exec(sqlDelete, "REMOVED", value)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

		_, err = con.Exec(sqlInsertTransactionDetail, transId, value)

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

func CreateReturn(data M.TransactionForm, claims M.JwtCustomClaims) (M.ResponseNoData, error) {
	TransactionLock.Lock()
	defer TransactionLock.Unlock()

	var res M.ResponseNoData

	if data.BooksQty == nil {
		res.Status = http.StatusInternalServerError
		res.Msg = "Books cannot be nil!"
		res.Success = false
		return res, nil
	}

	con := A.GetDB()

	tx, err := con.Begin()

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	var transOld M.TransactionLoan
	var transId *int

	var oldMemberId *int

	sqlFind := `
		select 
			t1.id as trans_id, 
			t1.date, 
			t1.detail, 
			t1.created_at, 
			t1.expected_return_date,
			t1.approval_status,
			t1.is_returned,
			t1.member_id
		from public.transactions t1 where id = $1
	`

	err = con.QueryRow(sqlFind, data.TransactionBeforeId).Scan(
		&transOld.Id,
		&transOld.Date,
		&transOld.Detail,
		&transOld.CreatedAt,
		&transOld.ExpectedReturnDate,
		&transOld.ApprovalStatus,
		&transOld.IsReturned,
		&oldMemberId,
	)

	if err != nil {
		tx.Rollback()
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	sql := `
	INSERT INTO public.transactions(
		transaction_type, member_id, date, transaction_before_id, detail, created_at, created_by, approval_status, approver_id, penalty_fee)
	VALUES ($1, $2, $3, $4, $5, NOW(), $6, $7, $8, $9) RETURNING id;
	`

	t1, err := time.Parse("2006-01-02", *data.Date)

	if err != nil {
		tx.Rollback()
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	if claims.Role == 0 {
		tx.Rollback()
		res.Status = http.StatusUnauthorized
		res.Msg = "Unouthorized!"
		res.Success = false
		return res, nil
	} else {
		err := con.QueryRow(sql, data.TransactionType, oldMemberId, t1, data.TransactionBeforeId, data.Detail, claims.ID, "APPROVE", claims.ID, data.PenaltyFee).Scan(&transId)

		if transId == nil {
			tx.Rollback()
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
		SET status = $1, condition = $2
		WHERE id = $3
	`

	//update book statuses
	for key, elem := range *data.BooksQty {
		var book_id *int
		var status string

		var condition string

		switch elem["condition"] {
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

		err := con.QueryRow(sqlFind, key).Scan(&book_id, &status, &condition)

		if err != nil {
			tx.Rollback()
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

		if status == "REMOVED" || status == "MISSING" {
			tx.Rollback()
			res.Status = http.StatusInternalServerError
			res.Msg = "Book is missing or already removed!"
			res.Success = false
			return res, err
		}

		switch elem["status"] {
		case 0:
			status = "STORED"
		case 1:
			status = "MISSING"
		default:
			status = "STORED"
		}

		_, err = con.Exec(sql, transId, book_id, key, condition)

		if err != nil {
			tx.Rollback()
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

		_, err = con.Exec(sqlUpdate, status, condition, key)

		if err != nil {
			tx.Rollback()
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

		// err = M.UpdateBookStock(*book_id)

		// if err != nil {
		// 	tx.Rollback()
		// 	res.Status = http.StatusInternalServerError
		// 	res.Msg = err.Error()
		// 	res.Success = false
		// 	return res, err
		// }
	}

	updateSql := `update transactions set is_returned = true where id = $1;`

	_, err = con.Exec(updateSql, transOld.Id)

	if err != nil {
		tx.Rollback()
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	if err := tx.Commit(); err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Success to create transaction!"
	res.Success = true
	return res, nil
}

func CreateOnlineLoan(data M.TransactionForm, claims M.JwtCustomClaims) (M.Response, error) {

	TransactionLock.Lock()
	defer TransactionLock.Unlock()

	var res M.Response

	if data.BooksId == nil {
		res.Status = http.StatusInternalServerError
		res.Msg = "Books cannot be nil!"
		res.Success = false
		return res, nil
	}

	con := A.GetDB()

	// tx, err := con.Begin()

	sql := `
	INSERT INTO public.transactions(
		transaction_type, member_id, date, detail, created_at, created_by, approval_status, approver_id, expected_return_date)
		VALUES ($1, $2, $3, $4, NOW(), $5, $6, $7, $8) RETURNING id
	`

	t1 := time.Now()

	t2 := t1.AddDate(0, 0, 2) //make this dynamic prob

	var transId *int

	if claims.Role == 0 {
		err := con.QueryRow(sql, "LOAN", claims.ID, t1, "Online Loan", claims.ID, "APPROVE", nil, t2).Scan(&transId)

		if transId == nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}
	} else {

		res.Status = http.StatusInternalServerError
		res.Msg = "Unauthorized!"
		res.Success = false
		return res, nil

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

	sqlInsertAccess := `
		INSERT INTO public.member_access(
			member_id, book_id, book_details_id, transaction_id)
		VALUES ($1, $2, $3, $4);
	`

	sqlInsertToCollection := `
	INSERT INTO public.member_collections (member_id, book_id)
	SELECT $1, $2
	WHERE NOT EXISTS (
		SELECT 1 FROM member_collections 
		WHERE member_id = $1 AND book_id = $2
	);
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

		_, err = con.Exec(sqlInsertAccess, claims.ID, book_id, elem, transId)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

		// err = M.UpdateBookStock(*book_id)

		// if err != nil {
		// 	res.Status = http.StatusInternalServerError
		// 	res.Msg = err.Error()
		// 	res.Success = false
		// 	return res, err
		// }

		_, err = con.Exec(sqlInsertToCollection, claims.ID, book_id)

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
	res.Data = (*data.BooksId)[0]
	return res, nil

}

func CreateOnlineReturn(data M.TransactionForm, claims M.JwtCustomClaims) (M.Response, error) {
	TransactionLock.Lock()
	defer TransactionLock.Unlock()

	var res M.Response

	if data.BooksId == nil {
		res.Status = http.StatusInternalServerError
		res.Msg = "Books cannot be nil!"
		res.Success = false
		return res, nil
	}

	con := A.GetDB()

	// tx, err := con.Begin()

	var oldTransId *int

	sql := `
	SELECT 
	t1.transaction_id 
	from member_access t1 where 
	t1.member_id = $1 and t1.book_details_id = $2
	`

	err := con.QueryRow(sql, claims.ID, (*data.BooksId)[0]).Scan(&oldTransId)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	sql = `
	INSERT INTO public.transactions(
		transaction_type, member_id, date, detail, created_at, created_by, approval_status, approver_id, transaction_before_id)
		VALUES ($1, $2, $3, $4, NOW(), $5, $6, $7, $8) RETURNING id
	`

	t1 := time.Now()

	var transId *int

	if claims.Role == 0 {
		err := con.QueryRow(sql, "RETURN", claims.ID, t1, "Online Return", claims.ID, "APPROVE", nil, &oldTransId).Scan(&transId)

		if transId == nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, nil
		}
	} else {

		res.Status = http.StatusInternalServerError
		res.Msg = "Unauthorized!"
		res.Success = false
		return res, nil

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

	sqlDeleteAccess := `
		DELETE FROM public.member_access
		WHERE member_id = $1 AND transaction_id = $2 AND book_details_id = $3 AND book_id = $4;
	`

	var book_id *int
	var status *string
	var condition *string

	//update book statuses
	for _, elem := range *data.BooksId {

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

		_, err = con.Exec(sqlDeleteAccess, claims.ID, oldTransId, elem, book_id)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

		// err = M.UpdateBookStock(*book_id)

		// if err != nil {
		// 	res.Status = http.StatusInternalServerError
		// 	res.Msg = err.Error()
		// 	res.Success = false
		// 	return res, err
		// }
	}

	updateSql := `update transactions set is_returned = true where id = $1;`

	_, err = con.Exec(updateSql, data.TransactionBeforeId)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Success to create transaction!"
	res.Success = true
	res.Data = book_id
	return res, nil
}
