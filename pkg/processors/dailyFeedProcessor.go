package processors

import (
	"boonmafarm/api/pkg/services"
	"boonmafarm/api/utils/excelutil"
	"boonmafarm/api/utils/timeutil"
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type IDailyFeedProcessor interface {
	DownloadExcelForm(clientId int, formType string, feedId, farmId int, date string) ([]byte, error)
}

type dailyFeedProcessorImp struct {
	DailyFeedService      services.IDailyFeedService
	FeedCollectionService services.IFeedCollectionService
	FarmService           services.IFarmService
	PondService           services.IPondService
}

func NewDailyFeedProcessor(dailyFeedService services.IDailyFeedService, feedCollectionService services.IFeedCollectionService, farmService services.IFarmService, pondService services.IPondService) IDailyFeedProcessor {
	return &dailyFeedProcessorImp{
		DailyFeedService:      dailyFeedService,
		FeedCollectionService: feedCollectionService,
		FarmService:           farmService,
		PondService:           pondService,
	}
}

func (p dailyFeedProcessorImp) DownloadExcelForm(clientId int, formType string, feedId, farmId int, date string) ([]byte, error) {
	type excelForm struct {
		FeedId    int
		FeedType  string
		FarmName  string
		Year      int
		Month     *int
		PondNames []string
	}

	// get feed collection
	feedCollection, err := p.FeedCollectionService.Get(feedId)
	if err != nil {
		return nil, err
	}

	// get farm
	farm, err := p.FarmService.Get(farmId, clientId)
	if err != nil {
		return nil, err
	}

	// get pondNames
	pondNames, err := p.PondService.GetPondNameList(farmId)
	if err != nil {
		return nil, err
	}

	dateParsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	year := dateParsed.Year()
	var month *int
	if formType == "month" {
		monthInt := int(dateParsed.Month())
		month = &monthInt
	}

	// create struct
	excelObj := excelForm{
		FeedId:    feedId,
		FeedType:  feedCollection.Name,
		FarmName:  farm.Name,
		Year:      year,
		Month:     month,
		PondNames: pondNames,
	}

	// create excel form
	f := excelize.NewFile()

	// Rename the default sheet to "รายเดือน" or "รายปี"
	sheetName := "รายปี"
	if formType == "month" {
		sheetName = "รายเดือน"
	}

	// Rename the default sheet (Sheet1) to the desired name
	defaultSheetName := "Sheet1"
	if err := f.SetSheetName(defaultSheetName, sheetName); err != nil {
		return nil, err
	}
	err = f.SetDefaultFont("TH SarabunPSK")

	// Set active sheet of the workbook
	f.SetActiveSheet(0)
	if err != nil {
		return nil, err
	}

	// Align text to the right in A1
	styleRightAlign, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	styleCenterAlign, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	// Set vertical alignment (middle) and rotated text
	styleMidVert, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical:     "center",
			Horizontal:   "center",
			TextRotation: 90,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	styleTotalRow, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#A52A2A"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	styleBackground, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#01FF00"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	styleFillData_A, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#ffff04"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	styleFillData_B, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#f1c233"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	// Write data to the sheet
	f.SetCellValue(sheetName, "A1", "ประเภท:")
	f.SetCellStyle(sheetName, "A1", "A1", styleRightAlign)

	err = f.MergeCell(sheetName, "B1", "C1")
	if err != nil {
		return nil, err
	}

	err = f.MergeCell(sheetName, "D1", "E1")
	if err != nil {
		return nil, err
	}

	f.SetCellValue(sheetName, "B1", excelObj.FeedType)
	f.SetCellStyle(sheetName, "B1", "C1", styleCenterAlign)

	f.SetCellValue(sheetName, "D1", "ไอดีอาหาร:")
	f.SetCellValue(sheetName, "F1", excelObj.FeedId)
	f.SetCellStyle(sheetName, "D1", "F1", styleRightAlign)

	f.SetCellValue(sheetName, "G1", "ฟาร์ม:")
	f.SetCellValue(sheetName, "H1", excelObj.FarmName)
	f.SetCellStyle(sheetName, "G1", "H1", styleRightAlign)

	f.SetCellValue(sheetName, "I1", "ปี:")
	f.SetCellValue(sheetName, "J1", excelObj.Year+543)
	f.SetCellStyle(sheetName, "I1", "J1", styleRightAlign)

	f.SetCellValue(sheetName, "A2", "เดือน")
	f.SetCellValue(sheetName, "B2", "วันที่")

	// Set pond names header
	for i, pondName := range excelObj.PondNames {
		cell := fmt.Sprintf("%s2", excelutil.ColName(i+2))
		trimPond := strings.TrimLeft(pondName, "บ่อ")
		trimPond = strings.TrimSpace(trimPond)
		f.SetCellValue(sheetName, cell, trimPond)
	}

	columnTotalRight := fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), 2)
	f.SetCellValue(sheetName, columnTotalRight, fmt.Sprintf("รวม%s", excelObj.FarmName))

	//use for loop to set month or year

	// variable before loop
	latestRow := 2

	for i := 1; i <= 12; i++ {
		if excelObj.Month != nil && i != *excelObj.Month {
			continue
		}

		totalDay := timeutil.DaysInMonth(excelObj.Year, time.Month(i))
		startRow := latestRow + 1
		endRow := latestRow + totalDay
		totalMonthRow := endRow + 1

		startCell := fmt.Sprintf("A%d", startRow)
		err = f.MergeCell(sheetName, startCell, fmt.Sprintf("A%d", endRow))
		if err != nil {
			return nil, err
		}

		f.SetCellValue(sheetName, startCell, timeutil.FullThaiMonths[i])
		f.SetCellStyle(sheetName, startCell, fmt.Sprintf("A%d", endRow), styleMidVert)

		// Set dates in column B from row 3
		for day := 1; day <= totalDay; day++ {
			cell := fmt.Sprintf("B%d", day+latestRow)
			date := time.Date(year, time.Month(i), day, 0, 0, 0, 0, time.UTC)

			// Format the date as "1 ม.ค. 2564" (in B.E.)
			dayString := fmt.Sprintf("%d %s %d", day, timeutil.ThaiMonths[i], date.Year()+543)

			f.SetCellValue(sheetName, cell, dayString)

			startStylingrow := fmt.Sprintf("C%d", day+latestRow)
			endStylingrow := fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+1), day+latestRow)
			if day%2 == 0 {
				f.SetCellStyle(sheetName, startStylingrow, endStylingrow, styleFillData_B)
			} else {
				f.SetCellStyle(sheetName, startStylingrow, endStylingrow, styleFillData_A)
			}

			f.SetCellFormula(sheetName, fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), day+latestRow), fmt.Sprintf("SUM(C%d:%s%d)", day+latestRow, excelutil.ColName(len(excelObj.PondNames)+1), day+latestRow))
		}

		err = f.MergeCell(sheetName, fmt.Sprintf("A%d", totalMonthRow), fmt.Sprintf("B%d", totalMonthRow))
		if err != nil {
			return nil, err
		}
		columnTotalBottom := fmt.Sprintf("A%d", totalMonthRow)
		f.SetCellValue(sheetName, columnTotalBottom, fmt.Sprintf("รวมเดือน: %s", timeutil.FullThaiMonths[time.Month(i)]))
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", totalMonthRow), fmt.Sprintf("B%d", totalMonthRow), styleRightAlign)
		// Set background color of the total row
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", totalMonthRow), fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), totalMonthRow), styleBackground)

		// set formula for total row
		for j := range excelObj.PondNames {
			cell := fmt.Sprintf("%s%d", excelutil.ColName(j+2), totalMonthRow)
			f.SetCellFormula(sheetName, cell, fmt.Sprintf("SUM(%s%d:%s%d)", excelutil.ColName(j+2), startRow+1, excelutil.ColName(j+2), endRow))
		}

		f.SetCellFormula(sheetName, fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), totalMonthRow), fmt.Sprintf("SUM(%s%d:%s%d)", excelutil.ColName(len(excelObj.PondNames)+2), startRow, excelutil.ColName(len(excelObj.PondNames)+2), endRow))
		f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), totalMonthRow), fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), totalMonthRow), styleTotalRow)

		// set variable
		latestRow = totalMonthRow
	}

	// Set row 2 center
	f.SetCellStyle(sheetName, "A2", fmt.Sprintf("%s2", excelutil.ColName(len(excelObj.PondNames)+2)), styleCenterAlign)
	f.SetColWidth(sheetName, "B", "B", 12)

	f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), 2), fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), latestRow), styleBackground)

	// Convert the Excel file to bytes
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	// Send the response
	return buf.Bytes(), nil
}