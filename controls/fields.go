package controls

type TextField struct {
	Value string `form:"value" binding:"required"`
}

type ChoiceField struct {
	Value   string `form:"value" binding:"required"`
	Choices []string
}

type MultiChoiceField struct {
	Values  []string `form:"values" binding:"required,min=1"`
	Choices []string
}
