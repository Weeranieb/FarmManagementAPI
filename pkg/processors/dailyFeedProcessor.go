package processors

import (
	dbContext "boonmafarm/api/pkg/dbcontext"
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/services"
	"boonmafarm/api/utils/excelutil"
	"boonmafarm/api/utils/timeutil"
	"bytes"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type IDailyFeedProcessor interface {
	DownloadExcelForm(clientId int, formType string, feedId, farmId int, date string) ([]byte, error)
	UploadExcelForm(file *multipart.FileHeader, username string, clientId int) error
	BulkCreateAndUpdate(dailyFeeds []*models.DailyFeed, username string) error
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

func (p dailyFeedProcessorImp) UploadExcelForm(file *multipart.FileHeader, username string, clientId int) error {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return err
	}

	// Read the file
	excelFile, err := excelize.OpenReader(src)
	if err != nil {
		return err
	}

	// Get the sheet name
	sheetName := excelFile.GetSheetName(0)
	fmt.Printf("Sheet Name: %s", sheetName)

	if sheetName != "รายเดือน" && sheetName != "รายปี" {
		return fmt.Errorf("invalid sheet name")
	}

	// Get detail from the sheet
	sFeedId, _ := excelFile.GetCellValue(sheetName, "F1")

	feedId, err := strconv.Atoi(sFeedId)
	if err != nil {
		return err
	}

	sYear, _ := excelFile.GetCellValue(sheetName, "J1")

	year, err := strconv.Atoi(sYear)
	if err != nil {
		return err
	}

	year -= 543

	farmName, _ := excelFile.GetCellValue(sheetName, "H1")

	sMonth, _ := excelFile.GetCellValue(sheetName, "A3")
	month := timeutil.FullThaiMonthToTime[sMonth]
	monthInt := int(month)

	// get farmId by name
	farmId, err := p.FarmService.GetFarmIdByName(farmName, clientId)
	if err != nil {
		return err
	}

	if sheetName == "รายเดือน" {
		isAvailable, err := p.DailyFeedService.IsFeedOnDateAvailable(feedId, farmId, year, &monthInt)
		if err != nil {
			return err
		}

		if !isAvailable {
			return fmt.Errorf("feed already exist")
		}
	} else if sheetName == "รายปี" {
		isAvailable, err := p.DailyFeedService.IsFeedOnDateAvailable(feedId, farmId, year, nil)
		if err != nil {
			return err
		}

		if !isAvailable {
			return fmt.Errorf("feed already exist")
		}
	}

	// get pond List
	pondList, err := p.PondService.GetPondList(farmId)
	if err != nil {
		return err
	}

	// check if pond Name in excel is the same as the pond name in the database
	expectedTotalCell := fmt.Sprintf("%s%d", excelutil.ColName(len(pondList)+2), 2)
	if val, _ := excelFile.GetCellValue(sheetName, expectedTotalCell); !strings.Contains(val, "รวม") {
		return fmt.Errorf("excel format is out of date")
	}

	// isLeapYear := timeutil.IsLeapYear(year)
	rows, _ := excelFile.GetRows(sheetName)
	rows = rows[2:]

	var payload []*models.AddDailyFeed

	day := 1
	totalDay := timeutil.DaysInMonth(year, timeutil.FullThaiMonthToTime[sMonth])
	totalCellRow := totalDay
	for rowIdx, row := range rows {
		column := row[2 : len(row)-1]
		if rowIdx == totalCellRow {
			if sheetName == "รายปี" {
				day = 1

				month++
				totalCellRow += timeutil.DaysInMonth(year, month) + 1
			}

			continue
			// set new totalDay
		}
		for i, cell := range column {
			if cell == "" {
				continue
			}

			temp := &models.AddDailyFeed{}
			amount, err := strconv.ParseFloat(cell, 64)
			if err != nil {
				return err
			}

			temp.FeedCollectionId = feedId
			temp.Amount = amount
			temp.FeedDate = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
			temp.PondId = pondList[i].Id

			payload = append(payload, temp)
		}
		day++
	}

	// bulk insert daily feed
	if err := p.DailyFeedService.BulkCreate(payload, username); err != nil {
		return err
	}

	return nil
}

func (p dailyFeedProcessorImp) BulkCreateAndUpdate(dailyFeeds []*models.DailyFeed, username string) error {
	// start transaction
	db := dbContext.Context.Postgresql
	tx := db.Begin()

	// separate daily feed to create and update by checking the id
	var createPayload []*models.AddDailyFeed
	var updatePayload []*models.DailyFeed

	for _, dailyFeed := range dailyFeeds {
		if dailyFeed.Id == -1 {
			temp := &models.AddDailyFeed{}
			temp.Amount = dailyFeed.Amount
			temp.FeedCollectionId = dailyFeed.FeedCollectionId
			temp.FeedDate = dailyFeed.FeedDate
			temp.PondId = dailyFeed.PondId

			createPayload = append(createPayload, temp)
		} else {
			updatePayload = append(updatePayload, dailyFeed)
		}
	}

	// bulk insert daily feed
	if len(createPayload) > 0 {
		if err := p.DailyFeedService.WithTrx(tx).BulkCreate(createPayload, username); err != nil {
			return err
		}
	}

	// bulk update daily feed
	if len(updatePayload) > 0 {
		if err := p.DailyFeedService.WithTrx(tx).BulkUpdate(updatePayload, username); err != nil {
			return err
		}
	}

	// commit transaction
	tx.Commit()

	return nil
}
