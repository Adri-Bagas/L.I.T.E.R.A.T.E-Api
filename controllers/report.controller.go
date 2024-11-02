package controllers

import (
	"fmt"
	"log"
	"net/http"
	A "perpus_api/db"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo/v4"
)

func GetTransactionExcel(c echo.Context) error {
	xlsx := excelize.NewFile()

	sheet1Name := "Loan"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)
	sheet2Name := "Inventory In"
	sheetIndex := xlsx.NewSheet(sheet2Name)
	xlsx.SetActiveSheet(sheetIndex)
	sheet3Name := "Inventory Out"
	sheetIndex = xlsx.NewSheet(sheet3Name)
	xlsx.SetActiveSheet(sheetIndex)
	sheet4Name := "Return"
	sheetIndex = xlsx.NewSheet(sheet4Name)
	xlsx.SetActiveSheet(sheetIndex)

	SheetOneData, err := GetLoanDataAllObj(c)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	SheetTwoData, err := GetTransactionInOutDataAllObj(c, "INVENTORY_IN")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	SheetThreeData, err := GetTransactionInOutDataAllObj(c, "INVENTORY_OUT")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	SheetFourData, err := GetReturnDataAllObj(c)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	style, err := xlsx.NewStyle(`{
		"font": {
			"bold": true,
			"size": 12
		},
		"fill": {
			"type": "pattern",
			"color": ["#E0EBF5"],
			"pattern": 1
		}
	}`)

	if err != nil {
		fmt.Println(err)
	}

	xlsx.SetCellStyle(sheet1Name, "A1", "G1", style)
	xlsx.SetCellStyle(sheet4Name, "A1", "F1", style)
	xlsx.SetCellStyle(sheet2Name, "A1", "D1", style)
	xlsx.SetCellStyle(sheet3Name, "A1", "D1", style)

	xlsx.SetCellValue(sheet1Name, "A1", "Loaner")
	xlsx.SetCellValue(sheet1Name, "B1", "Date")
	xlsx.SetCellValue(sheet1Name, "C1", "Expected Return")
	xlsx.SetCellValue(sheet1Name, "D1", "Is Returned")
	xlsx.SetCellValue(sheet1Name, "E1", "Detail")
	xlsx.SetCellValue(sheet1Name, "F1", "Approval")
	xlsx.SetCellValue(sheet1Name, "G1", "Approved by")

	xlsx.SetCellValue(sheet4Name, "A1", "Loaner")
	xlsx.SetCellValue(sheet4Name, "B1", "Date")
	xlsx.SetCellValue(sheet4Name, "C1", "Penalty")
	xlsx.SetCellValue(sheet4Name, "D1", "Detail")
	xlsx.SetCellValue(sheet4Name, "E1", "Approval")
	xlsx.SetCellValue(sheet4Name, "F1", "Approved by")

	xlsx.SetCellValue(sheet2Name, "A1", "Date")
	xlsx.SetCellValue(sheet2Name, "B1", "Detail")
	xlsx.SetCellValue(sheet2Name, "C1", "Created at")
	xlsx.SetCellValue(sheet2Name, "D1", "Approved by")

	xlsx.SetCellValue(sheet3Name, "A1", "Date")
	xlsx.SetCellValue(sheet3Name, "B1", "Detail")
	xlsx.SetCellValue(sheet3Name, "C1", "Created at")
	xlsx.SetCellValue(sheet3Name, "D1", "Approved by")

	err = xlsx.AutoFilter(sheet1Name, "A1", "G1", "")
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	err = xlsx.AutoFilter(sheet4Name, "A1", "F1", "")
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	err = xlsx.AutoFilter(sheet2Name, "A1", "D1", "")
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	err = xlsx.AutoFilter(sheet3Name, "A1", "D1", "")
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	for i, each := range SheetFourData {
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+2), each.Member.FullName)
		if each.Date != nil {
			xlsx.SetCellValue(sheet4Name, fmt.Sprintf("B%d", i+2), *each.Date)
		}
		if each.PenaltyFee != nil {
			xlsx.SetCellValue(sheet4Name, fmt.Sprintf("C%d", i+2), *each.PenaltyFee)
		}

		if each.Detail != nil {
			xlsx.SetCellValue(sheet4Name, fmt.Sprintf("D%d", i+2), *each.Detail)
		}

		if each.ApprovalStatus != nil {
			xlsx.SetCellValue(sheet4Name, fmt.Sprintf("E%d", i+2), *each.ApprovalStatus)
		}

		if each.Approver.Username != nil {
			xlsx.SetCellValue(sheet4Name, fmt.Sprintf("F%d", i+2), *each.Approver.Username)
		}
	}

	for i, each := range SheetOneData {
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+2), each.Member.FullName)
		if each.Date != nil {
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+2), *each.Date)
		}
		if each.ExpectedReturnDate != nil {
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", i+2), *each.ExpectedReturnDate)
		}
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", i+2), each.IsReturned)

		if each.Detail != nil {
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", i+2), *each.Detail)
		}

		if each.ApprovalStatus != nil {
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", i+2), *each.ApprovalStatus)
		}

		if each.Approver.Username != nil {
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", i+2), *each.Approver.Username)
		}
	}

	for i, each := range SheetTwoData {
		if each.Date != nil {
			xlsx.SetCellValue(sheet2Name, fmt.Sprintf("A%d", i+2), *each.Date)
		}
		if each.Detail != nil {
			xlsx.SetCellValue(sheet2Name, fmt.Sprintf("B%d", i+2), *each.Detail)
		}
		if each.CreatedAt != nil {
			xlsx.SetCellValue(sheet2Name, fmt.Sprintf("C%d", i+2), *each.CreatedAt)
		}
		if each.Approver.Username != nil {
			xlsx.SetCellValue(sheet2Name, fmt.Sprintf("D%d", i+2), *each.Approver.Username)
		}
	}

	for i, each := range SheetThreeData {
		if each.Date != nil {
			xlsx.SetCellValue(sheet3Name, fmt.Sprintf("A%d", i+2), *each.Date)
		}
		if each.Detail != nil {
			xlsx.SetCellValue(sheet3Name, fmt.Sprintf("B%d", i+2), *each.Detail)
		}
		if each.CreatedAt != nil {
			xlsx.SetCellValue(sheet3Name, fmt.Sprintf("C%d", i+2), *each.CreatedAt)
		}
		if each.Approver.Username != nil {
			xlsx.SetCellValue(sheet3Name, fmt.Sprintf("D%d", i+2), *each.Approver.Username)
		}
	}

	err = xlsx.SaveAs("./file1.xlsx")
	if err != nil {
		fmt.Println(err)
	}

	return c.File("./file1.xlsx")

}

type TransactionData struct {
	Date             string `json:"date"`
	TransactionCount int    `json:"transaction_count"`
}

func GetDataDashboardChart(c echo.Context) error {
	sql := `SELECT DATE(created_at) AS date, COUNT(*) AS transaction_count FROM transactions GROUP BY DATE(created_at) ORDER BY transaction_count`

	con := A.GetDB()

	rows, err := con.Query(sql)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var transactionData []TransactionData

	for rows.Next() {
		var date string
		var transactionCount int
		if err := rows.Scan(&date, &transactionCount); err != nil {
			log.Fatal(err)
		}
		transactionData = append(transactionData, TransactionData{Date: date, TransactionCount: transactionCount})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, transactionData)

}

type DataCount struct {
	PenaltySum       int `json:"penalty_sum"`
	MemberCount      int    `json:"member_count"`
	TransactionCount int    `json:"transaction_count"`
	BookCount        int    `json:"book_count"`
}

func GetDataCountAndSumDashboard(c echo.Context) error {

	var datas DataCount

	con := A.GetDB()

	sql := `SELECT * from summary_dashboard();`

	err := con.QueryRow(sql).Scan(
		&datas.PenaltySum,
		&datas.MemberCount,
		&datas.TransactionCount,
		&datas.BookCount,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, datas)
}
