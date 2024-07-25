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

	f.SetCellValue(sheetName, "A2", "เดือน")
	f.SetCellValue(sheetName, "B2", "วันที่")

	totalDay := timeutil.DaysInMonth(excelObj.Year, time.Month(*excelObj.Month))
	endCell := 2 + totalDay
	err = f.MergeCell(sheetName, "A3", fmt.Sprintf("A%d", endCell))
	if err != nil {
		return nil, err
	}

	f.SetCellValue(sheetName, "A3", timeutil.FullThaiMonths[time.Month(*excelObj.Month)])
	f.SetCellStyle(sheetName, "A3", fmt.Sprintf("A%d", endCell), styleMidVert)

	// Set dates in column B from row 3
	for day := 1; day <= totalDay; day++ {
		cell := fmt.Sprintf("B%d", day+2) // Starting from B3
		date := time.Date(year, time.Month(*excelObj.Month), day, 0, 0, 0, 0, time.UTC)

		// Format the date as "1 ม.ค. 2564" (in B.E.)
		dayString := fmt.Sprintf("%d %s %d", day, timeutil.ThaiMonths[date.Month()], date.Year()+543)

		f.SetCellValue(sheetName, cell, dayString)

		startCell := fmt.Sprintf("C%d", day+2)
		endCell := fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+1), day+2)
		if day%2 == 0 {
			f.SetCellStyle(sheetName, startCell, endCell, styleFillData_B)
		} else {
			f.SetCellStyle(sheetName, startCell, endCell, styleFillData_A)
		}

		f.SetCellFormula(sheetName, fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), day+2), fmt.Sprintf("SUM(C%d:%s%d)", day+2, excelutil.ColName(len(excelObj.PondNames)+1), day+2))
	}

	// Set pond names header
	for i, pondName := range excelObj.PondNames {
		cell := fmt.Sprintf("%s2", excelutil.ColName(i+2))
		trimPond := strings.TrimLeft(pondName, "บ่อ")
		trimPond = strings.TrimSpace(trimPond)
		f.SetCellValue(sheetName, cell, trimPond)

		f.SetCellFormula(sheetName, fmt.Sprintf("%s%d", excelutil.ColName(i+2), endCell+1), fmt.Sprintf("SUM(%s%d:%s%d)", excelutil.ColName(i+2), 3, excelutil.ColName(i+2), endCell))
	}

	bottomRowIdx := endCell + 1
	columnTotalRight := fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), 2)
	err = f.MergeCell(sheetName, fmt.Sprintf("A%d", bottomRowIdx), fmt.Sprintf("B%d", bottomRowIdx))
	if err != nil {
		return nil, err
	}
	columnTotalBottom := fmt.Sprintf("A%d", bottomRowIdx)
	f.SetCellValue(sheetName, columnTotalRight, fmt.Sprintf("รวม%s", excelObj.FarmName))
	f.SetCellValue(sheetName, columnTotalBottom, "รวม")

	f.SetCellFormula(sheetName, fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), bottomRowIdx), fmt.Sprintf("SUM(%s%d:%s%d)", excelutil.ColName(len(excelObj.PondNames)+2), 3, excelutil.ColName(len(excelObj.PondNames)+2), endCell))

	// Set row 2 center
	f.SetCellStyle(sheetName, "A2", fmt.Sprintf("%s2", excelutil.ColName(len(excelObj.PondNames)+2)), styleCenterAlign)
	f.SetColWidth(sheetName, "B", "B", 12)
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", bottomRowIdx), fmt.Sprintf("B%d", bottomRowIdx), styleRightAlign)

	// Set background color
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", bottomRowIdx), fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), bottomRowIdx), styleBackground)
	f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), 2), fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), bottomRowIdx), styleBackground)

	// // Set borders
	// f.SetCellStyle(sheetName, "A2", fmt.Sprintf("%s%d", excelutil.ColName(len(excelObj.PondNames)+2), bottomRowIdx), styleBorders)

	// Convert the Excel file to bytes
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	// Send the response
	return buf.Bytes(), nil
}
