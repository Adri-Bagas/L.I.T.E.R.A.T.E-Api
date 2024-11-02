package models

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

type TransactionInventoryInOut struct {
	Id        int64   `json:"id"`
	Date      *string `json:"date"`
	Detail    *string `json:"detail"`
	CreatedAt *string `json:"created_at"`
	Approver  User  `json:"approver"`
}

type TransactionLoan struct {
	Id                 int64            `json:"id"`
	Date               *string          `json:"date"`
	ExpectedReturnDate *string          `json:"expected_return_date,omitempty"`
	Detail             *string          `json:"detail"`
	CreatedAt          *string          `json:"created_at"`
	ApprovalStatus     *string          `json:"approval_status,omitempty"`
	IsReturned         bool             `json:"is_returned"`
	PenaltyFee         *int             `json:"penalty,omitempty"`
	MemberId           *int64           `json:"member_id,omitempty"`
	Member             *Member        `json:"member,omitempty"`
	Approver           *User          `json:"approver,omitempty"`
	Books              *[]BookDetails `json:"books,omitempty"`
}