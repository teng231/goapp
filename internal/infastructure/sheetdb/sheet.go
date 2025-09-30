package sheetdb

import (
	"context"
	"fmt"
	"log"
	"teng231/goapp/internal/domain"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// SheetDB định nghĩa interface thao tác
type SheetDB interface {
	AppendComment(c *domain.Comment) error
}

// sheetService cài đặt SheetDB
type sheetService struct {
	srv           *sheets.Service
	spreadsheetID string
	sheetName     string
}

// NewSheetDB tạo instance sheetService
func NewSheetDB(ctx context.Context, credentials, spreadsheetID, sheetName string) SheetDB {
	conf, err := google.JWTConfigFromJSON([]byte(credentials), sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse credentials: %v", err)
	}
	client := conf.Client(ctx)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Sheets client: %v", err)
	}
	return &sheetService{
		srv:           srv,
		spreadsheetID: spreadsheetID,
		sheetName:     sheetName,
	}
}

func (s *sheetService) AppendComment(c *domain.Comment) error {
	// Lấy dữ liệu hiện tại
	readRange := fmt.Sprintf("%s!A:Z", s.sheetName)
	resp, err := s.srv.Spreadsheets.Values.Get(s.spreadsheetID, readRange).Do()
	if err != nil {
		return fmt.Errorf("unable to read sheet: %v", err)
	}

	rows := resp.Values
	userRow := -1
	for i, row := range rows {
		if len(row) > 0 && row[0] == c.Username {
			userRow = i + 1 // Sheet index bắt đầu từ 1
			break
		}
	}

	if userRow == -1 {
		// Username chưa có -> thêm dòng mới
		newRow := []interface{}{c.Username, c.Message}
		_, err = s.srv.Spreadsheets.Values.Append(s.spreadsheetID, s.sheetName+"!A:B", &sheets.ValueRange{
			Values: [][]interface{}{newRow},
		}).ValueInputOption("RAW").Do()
		if err != nil {
			return fmt.Errorf("unable to append new row: %v", err)
		}
		log.Printf("Added new row for %s", c.Username)
		return nil
	}

	// Username có rồi -> tìm cột trống
	row := rows[userRow-1]
	col := len(row) + 1 // cột tiếp theo
	cell := fmt.Sprintf("%s!%s%d", s.sheetName, colToLetter(col), userRow)
	_, err = s.srv.Spreadsheets.Values.Update(s.spreadsheetID, cell, &sheets.ValueRange{
		Values: [][]interface{}{{c.Message}},
	}).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("unable to update cell: %v", err)
	}
	log.Printf("Appended comment for %s at %s", c.Username, cell)
	return nil
}

// Helper: convert số cột -> chữ (1=A, 2=B, ...)
func colToLetter(n int) string {
	res := ""
	for n > 0 {
		rem := (n - 1) % 26
		res = string(rune('A'+rem)) + res
		n = (n - rem) / 26
	}
	return res
}
