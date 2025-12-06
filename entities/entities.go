package entities

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"name"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Role      string `json:"role"`
	// RegDate   string `json:"reg_date"`
	RegDate  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:reg_date" json:"reg_date,omitempty"`
}

type Orders struct {
	ID          string `json:"id" gorm:"column:id;primaryKey"`
	Description string `json:"description"`
	DateStart   string `json:"date_start"`
	DateEnd     string `json:"date_end"`
	Percent     string `json:"percent"`
}

type Review struct {
	ID      string `json:"id"`
	Rating  string `json:"rating"`
	Comment string `json:"comment"`
}

type Discount struct {
	ID          string  `json:"id" gorm:"column:id;primaryKey"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Percentage  float64 `json:"percentage"`
}

type Publisher struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Book struct {
	ID          string    `json:"id" gorm:"column:id;primaryKey"`
	ISBN        string    `json:"isbn"`
	PublisherID string    `json:"publisher_id"`
	Title       string    `json:"title"`
	AuthorID    string    `json:"author_id"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Language    string    `json:"language"`
	CategoryID  string    `json:"category_id"`
	Category    string    `json:"category"`
	Author      string    `json:"author"`
	Link        string    `json:"link"`
	Publisher   Publisher `json:"publisher" gorm:"foreignKey:PublisherID"`
}
