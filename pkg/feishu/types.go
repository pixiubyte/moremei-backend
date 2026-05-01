package feishu

type TokenResp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

type Common struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type SheetCreateReq struct {
	Title string `json:"title"`
}

type SheetCreateResp struct {
	Common
	Data struct {
		Spreadsheet Spreadsheet `json:"spreadsheet"`
	} `json:"data"`
}

type Spreadsheet struct {
	Title            string `json:"title"`
	FolderToken      string `json:"folder_token"`
	Url              string `json:"url"`
	SpreadsheetToken string `json:"spreadsheet_token"`
}

type QuerySheetResp struct {
	Common
	Data struct {
		Sheets []struct {
			SheetId        string `json:"sheet_id"`
			Title          string `json:"title"`
			Index          int    `json:"index"`
			Hidden         bool   `json:"hidden"`
			GridProperties struct {
				FrozenRowCount    int `json:"frozen_row_count"`
				FrozenColumnCount int `json:"frozen_column_count"`
				RowCount          int `json:"row_count"`
				ColumnCount       int `json:"column_count"`
			} `json:"grid_properties"`
			ResourceType string `json:"resource_type"`
			Merges       []struct {
				StartRowIndex    int `json:"start_row_index"`
				EndRowIndex      int `json:"end_row_index"`
				StartColumnIndex int `json:"start_column_index"`
				EndColumnIndex   int `json:"end_column_index"`
			} `json:"merges"`
		} `json:"sheets"`
	} `json:"data"`
}

type WriteSheetReq struct {
	ValueRange ValueRange `json:"valueRange"`
}

type ValueRange struct {
	Range  string          `json:"range"`
	Values [][]interface{} `json:"values"`
}

type WriteSheetResp struct {
	Common
	Data WriteSheetRespData `json:"data"`
}

type WriteSheetRespData struct {
	Revision         int    `json:"revision"`
	SpreadsheetToken string `json:"spreadsheetToken"`
	TableRange       string `json:"tableRange"`
	Updates          struct {
		Revision         int    `json:"revision"`
		SpreadsheetToken string `json:"spreadsheetToken"`
		UpdatedCells     int    `json:"updatedCells"`
		UpdatedColumns   int    `json:"updatedColumns"`
		UpdatedRange     string `json:"updatedRange"`
		UpdatedRows      int    `json:"updatedRows"`
	} `json:"updates"`
}
