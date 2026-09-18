package services

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"grip/dto"
	"grip/models"
	"grip/repository"

	"github.com/google/uuid"
)

// ImportContactsFromCSV parses a CSV stream with header name,type,value
// and bulk-creates contacts for userID. Invalid rows are skipped and
// reported; the user must exist (else gorm.ErrRecordNotFound).
func ImportContactsFromCSV(userID string, r io.Reader) (*dto.ImportSummary, error) {
	if _, err := repository.GetUserByID(userID); err != nil {
		return nil, err
	}

	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errors.New("empty CSV file")
		}
		return nil, fmt.Errorf("read CSV header: %w", err)
	}
	if len(header) != 3 ||
		!strings.EqualFold(strings.TrimSpace(header[0]), "name") ||
		!strings.EqualFold(strings.TrimSpace(header[1]), "type") ||
		!strings.EqualFold(strings.TrimSpace(header[2]), "value") {
		return nil, errors.New("invalid CSV header: expected name,type,value")
	}

	var rows []dto.ContactImportRow
	var preErrors []dto.RowError
	line := 1 // header is row 1
	for {
		record, err := reader.Read()
		line++
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read CSV row %d: %w", line, err)
		}
		if len(record) != 3 {
			preErrors = append(preErrors, dto.RowError{
				Row:     line,
				Message: fmt.Sprintf("expected 3 columns (name,type,value), got %d", len(record)),
			})
			continue
		}
		rows = append(rows, dto.ContactImportRow{
			Name:  strings.TrimSpace(record[0]),
			Type:  strings.TrimSpace(strings.ToLower(record[1])),
			Value: strings.TrimSpace(record[2]),
		})
	}

	return importRows(userID, rows, 2, preErrors)
}

// ImportContactsFromJSON parses a JSON array of {name,type,value} objects
// and bulk-creates contacts for userID.
func ImportContactsFromJSON(userID string, r io.Reader) (*dto.ImportSummary, error) {
	if _, err := repository.GetUserByID(userID); err != nil {
		return nil, err
	}

	var rows []dto.ContactImportRow
	if err := json.NewDecoder(r).Decode(&rows); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}
	if len(rows) == 0 {
		return nil, errors.New("empty JSON array")
	}

	return importRows(userID, rows, 1, nil)
}

// importRows creates every valid row, skipping invalid ones with per-row
// errors. firstRow is the 1-based number of rows[0] in the source file;
// pre holds file-level row errors (e.g. CSV column mismatches).
func importRows(userID string, rows []dto.ContactImportRow, firstRow int, pre []dto.RowError) (*dto.ImportSummary, error) {
	summary := &dto.ImportSummary{Message: "Import finished"}
	summary.Errors = append(summary.Errors, pre...)
	summary.Failed += len(pre)

	for i, row := range rows {
		rowNum := firstRow + i
		row.Type = strings.TrimSpace(strings.ToLower(row.Type))
		if err := dto.ValidateContactValue(row.Type, row.Value); err != nil {
			summary.Failed++
			summary.Errors = append(summary.Errors, dto.RowError{Row: rowNum, Message: err.Error()})
			continue
		}

		contact := &models.Contact{
			ID:     uuid.New().String(),
			UserID: userID,
			Name:   row.Name,
			Type:   row.Type,
			Value:  row.Value,
		}
		if err := repository.CreateContact(contact); err != nil {
			summary.Failed++
			summary.Errors = append(summary.Errors, dto.RowError{Row: rowNum, Message: err.Error()})
			continue
		}
		summary.Imported++
	}

	return summary, nil
}
