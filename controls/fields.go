package controls

type TextField struct {
	Value string `form:"value" binding:"required"`
}

type ChoiceField struct {
	Value   string `form:"value" binding:"required"`
	Choices []string
}
