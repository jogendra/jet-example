package domain

type ContentBlocksRequest struct {
	Page struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	} `json:"page"`
	// directly using struct here as for now not using this field
	Query struct {
		LeftOperand struct {
			Property       string      `json:"property"`
			SimpleOperator string      `json:"simpleOperator"`
			Value          interface{} `json:"value"`
		} `json:"leftOperand"`
		LogicalOperator string `json:"logicalOperator"`
		RightOperand    struct {
			Property       string      `json:"property"`
			SimpleOperator string      `json:"simpleOperator"`
			Value          interface{} `json:"value"`
		} `json:"rightOperand"`
	} `json:"query"`
	Sort []struct {
		Property  string `json:"property"`
		Direction string `json:"direction"`
	} `json:"sort"`
	Fields []string `json:"fields"`
}

type ContentBlock struct {
	Content string `json:"content"`
}
