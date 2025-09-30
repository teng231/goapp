package sheet

import (
	"fmt"

	"google.golang.org/api/sheets/v4"
)

// Sheet mô tả 1 bảng tính
type Sheet struct {
	SheetID string // Spreadsheet ID
	Url     string // Link tới sheet (nếu cần)
	Config  string // Tên sheet/tab, ví dụ "Sheet1"
}

// SheetDB định nghĩa interface thao tác
type SheetDB interface {
	Append(sheet *Sheet, data map[string]interface{}) error
	Edit(sheet *Sheet, data map[string]interface{}) error
}

// sheetService cài đặt SheetDB
type sheetService struct {
	srv *sheets.Service
}

// Append thêm 1 dòng cuối sheet
func (s *sheetService) Append(sheet *Sheet, data map[string]interface{}) error {
	var row []interface{}
	for k, v := range data {
		row = append(row, fmt.Sprintf("%s:%v", k, v))
	}
	vr := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}
	_, err := s.srv.Spreadsheets.Values.Append(sheet.SheetID, sheet.Config, vr).
		ValueInputOption("USER_ENTERED").Do()
	return err
}

// Edit: ví dụ đơn giản — tìm & ghi đè toàn bộ sheet bằng data mới
// (thực tế bạn có thể dùng tìm theo range hoặc cell để update)
func (s *sheetService) Edit(sheet *Sheet, data map[string]interface{}) error {
	var row []interface{}
	for k, v := range data {
		row = append(row, fmt.Sprintf("%s:%v", k, v))
	}
	vr := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}
	_, err := s.srv.Spreadsheets.Values.Update(sheet.SheetID, sheet.Config+"!A2",
		vr).ValueInputOption("USER_ENTERED").Do()
	return err
}
