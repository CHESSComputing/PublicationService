package main

// ZenodoError represents individual zenodo error struct
type ZenodoError struct {
	Field    string
	Messages []string
}

// ZenodoReponse represents zenodo response
type ZenodoResponse struct {
	Status  int
	Message string
	Error   []ZenodoError
}

// ZenodoCreatedResponse represents output of /create API
type ZenodoCreatedResponse struct {
	Id        int64    `json:"id"`
	MetaData  MetaData `json:"metadata"`
	Created   string   `json:"created"`
	Modified  string   `json:"modified"`
	Owner     int      `json:"owner"`
	RecordId  int64    `json:"record_id"`
	State     string   `json:"state"`
	Submitted bool     `json:"submitted"`
	Title     string   `json:"title"`
	Links     Links    `json:"self"`
}

// PrereserveDoi represents PrereserveDoi struct
type PrereserveDoi struct {
	Doi   string `json:"doi"`
	RecId int64  `json:"recid"`
}

// MetaData represents meta-data struct
type MetaData struct {
	AccessRight   string        `json:"access_right"`
	PrereserveDoi PrereserveDoi `json:"prereserve_doi"`
}

// Links contains zenodo links
type Links struct {
	Badge              string `json:"badge"`
	Discard            string `json:"discard"`
	Edit               string `json:"edit"`
	Files              string `json:"files"`
	Html               string `json:"html"`
	LatestDraft        string `json:"latest_draft"`
	LasestDraftHtml    string `json:"latest_draft_html"`
	NewVersion         string `json:"newversion"`
	Publish            string `json:"publish"`
	RegisterConceptDoi string `json:"registerconceptdoi"`
	Self               string `json:"self"`
}

// ZenodoAddResponse represents output of /create API
type ZenodoAddResponse struct {
	Created  string `json:"created"`
	Modified string `json:"modified"`
	Size     int64  `json:"size"`
	Key      string `json:"key"`
	MimeType string `json:"mimetype"`
	Checksum string `json:"checksum"`
	Owner    int    `json:"owner"`
	RecordId int64  `json:"record_id"`
	Links    Links  `json:"self"`
}
