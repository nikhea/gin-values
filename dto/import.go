package dto

// ContactImportRow is one record of a contacts import file (CSV row or
// JSON array element). Types must be one of email/phone/address/other.
type ContactImportRow struct {
	Name string `json:"name" csv:"name" example:"Work email"`

	Type string `json:"type" binding:"required,oneof=email phone address other" csv:"type" example:"email"`

	Value string `json:"value" binding:"required" csv:"value" example:"kaige@work.com"`
}

// RowError describes one failed import row.
type RowError struct {
	Row     int    `json:"row" example:"3"`
	Message string `json:"message" example:"invalid email value"`
}

// ImportSummary reports a bulk import: valid rows are created, invalid
// rows are skipped and listed in Errors.
type ImportSummary struct {
	Message  string     `json:"message" example:"Import finished"`
	Imported int        `json:"imported" example:"2"`
	Failed   int        `json:"failed" example:"1"`
	Errors   []RowError `json:"errors,omitempty"`
}

// AvatarResponse returns the public URL of an uploaded avatar.
type AvatarResponse struct {
	Message   string `json:"message" example:"Avatar uploaded"`
	AvatarURL string `json:"avatar_url" example:"/uploads/avatars/users/458622d8-daba-4252-8ce1-846277353139.png"`
}
