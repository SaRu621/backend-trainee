package structs

type Void struct{}

type User struct {
    ID string `json:"id"`
}

type Slug_ser struct {
    Slug    string  `json:"slug"`
    Percent float32 `json:"percent"`
}

type UserToSlug struct {
    ToDelete []string `json:"del"`
    ToAdd    []string `json:"add"`
    Exp      []Exp    `json:"exp"`
    ID         string `json:"id"`    
}

type UserEvents struct {
    ID    string `json:"id"`
    Year  string `json:"year"`
    Month string `json:"month"`
}

type CSV_link struct {
    URL string `json:"url"`
}

type Exp struct {
    Y   string `json:"y"`
    Mo  string `json:"mo"`
    D   string `json:"d"`
    H   string `json:"h"`
    Mi  string `json:"mi"`
}
