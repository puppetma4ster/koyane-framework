package utils

type Wordlist struct {
	ID                       uint          `gorm:"primaryKey;column:id"`
	Name                     string        `gorm:"column:name;not null"`
	Entities                 int64         `gorm:"column:entities;not null"`
	SmallestEntity           int64         `gorm:"column:smallestEntity;not null"`
	BiggestEntity            int64         `gorm:"column:biggestEntity;not null"`
	AverageLength            float64       `gorm:"column:averageLength;not null"`
	AverageEntropy           float64       `gorm:"column:averageEntropy;not null"`
	DigitsPercent            float64       `gorm:"column:digitsPercent;not null"`
	UpperCasePercent         float64       `gorm:"column:upperCasePercent;not null"`
	SpecialCharPercent       float64       `gorm:"column:specialCharPercent;not null"`
	DigitAndUpperCase        float64       `gorm:"column:digitAndUpperCase;not null"`
	DigitAndSpecialChar      float64       `gorm:"column:digitAndSpecialChar;not null"`
	UpperCaseAndSpecialChar  float64       `gorm:"column:upperCaseAndSpecialChar;not null"`
	DigitUpperCaseAndSpecial float64       `gorm:"column:digitUpperCaseAndSpecialChar;not null"`
	Encoding                 string        `gorm:"column:encoding;not null"`
	Language                 string        `gorm:"column:language;not null"`
	Category                 string        `gorm:"column:category;not null"`
	Author                   string        `gorm:"column:author;not null"`
	Size                     int64         `gorm:"column:size;not null"`
	Link                     string        `gorm:"column:link;not null"`
	Tags                     []WordlistTag `gorm:"foreignKey:WordlistID;references:ID"`
}

type WordlistTag struct {
	WordlistID uint   `gorm:"primaryKey;column:wordlist_id"`
	Tag        string `gorm:"primaryKey;column:tag"`
}

// TableName überschreibt den Tabellennamen, falls nötig
func (Wordlist) TableName() string {
	return "wordlists"
}
