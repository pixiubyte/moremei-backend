package cipher

type Category string

const (
    CategoryName   Category = "name"
    CategoryIdNum  Category = "id_num"
    CategoryMobile Category = "mobile"
    CategoryEmail  Category = "email"
    CategoryText   Category = "text"

    MaskSymbol     = "*"
    EmailSeparator = "@"
)
