package repositoryimpl

type Config struct {
	Table Table `json:"table" required:"true"`
}

type Table struct {
	ReviewComment      string `json:"review_comment"        required:"true"`
	TranslationComment string `json:"translation_comment"   required:"true"`
}
