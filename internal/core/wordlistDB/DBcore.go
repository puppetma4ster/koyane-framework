package wordlistDB

import (
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type WordlistRepository struct {
	db *gorm.DB
}

// NewWordlistRepository
// constructor
func NewWordlistRepository(dsn string) (*WordlistRepository, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &WordlistRepository{db: db}, nil
}

type WordlistFilter struct {
	name                                                     *[]string
	EntitiesMin, EntitiesMax                                 *uint64
	SmallestEntityMin, SmallestEntityMax                     *uint64
	BiggestEntityMin, BiggestEntityMax                       *uint64
	AverageLengthMin, AverageLengthMax                       *float64
	AverageEntropyMin, AverageEntropyMax                     *float64
	DigitsPercentMin, DigitsPercentMax                       *float64
	UpperCasePercentMin, UpperCasePercentMax                 *float64
	SpecialCharPercentMin, SpecialCharPercentMax             *float64
	DigitAndUpperCaseMin, DigitAndUpperCaseMax               *float64
	DigitAndSpecialCharMin, DigitAndSpecialCharMax           *float64
	UpperCaseAndSpecialCharMin, UpperCaseAndSpecialCharMax   *float64
	DigitUpperCaseAndSpecialMin, DigitUpperCaseAndSpecialMax *float64
	Encoding                                                 *[]string
	Language                                                 *[]string
	Category                                                 *[]string
	Author                                                   *[]string
	SizeMin, SizeMax                                         *uint64
	Tags                                                     *[]string
	Link                                                     *[]string
}

// NewWordlistFilter builds a WordlistFilter from optional string slices and optional numeric ranges.
// All inputs are optional (nil-safe). If a range is nil, the corresponding Min/Max pointers remain nil.
// Note: Category is not provided in the arguments and will be left nil intentionally.
func NewWordlistFilter(
	name, encoding, language, author, link *[]string,
	entitiesRange, SmallestEntitiesRange, BiggestEntitiesRange, sizeRange *utils.Uint64Range,
	AverageLengthRange, AverageEntropyRange, DigitsPercentRange,
	UpperCasePercentRange, SpecialCharPercentRange, DigitAndUpperCaseRange,
	DigitAndSpecialCharRange, UpperCaseAndSpecialCharRange, DigitUpperCaseAndSpecialRange *utils.Float64Range,
	tags *[]string,
) *WordlistFilter {

	f := &WordlistFilter{
		// string-slice filters (as pointers to slices)
		name:     name,
		Encoding: encoding,
		Language: language,
		Author:   author,
		Link:     link,
		Tags:     tags,

		// Category is intentionally left nil because the constructor doesn't receive it.
		Category: nil,
	}

	// Helper to map uint64 range -> (Min, Max) pointers on the filter
	setUintRange := func(r *utils.Uint64Range, dstMin **uint64, dstMax **uint64) {
		if r == nil {
			return
		}
		if r.Min != nil {
			*dstMin = r.Min
		}
		if r.Max != nil {
			*dstMax = r.Max
		}
	}

	// Helper to map float64 range -> (Min, Max) pointers on the filter
	setFloatRange := func(r *utils.Float64Range, dstMin **float64, dstMax **float64) {
		if r == nil {
			return
		}
		if r.Min != nil {
			*dstMin = r.Min
		}
		if r.Max != nil {
			*dstMax = r.Max
		}
	}

	// Map uint64 ranges
	setUintRange(entitiesRange, &f.EntitiesMin, &f.EntitiesMax)
	setUintRange(SmallestEntitiesRange, &f.SmallestEntityMin, &f.SmallestEntityMax)
	setUintRange(BiggestEntitiesRange, &f.BiggestEntityMin, &f.BiggestEntityMax)
	setUintRange(sizeRange, &f.SizeMin, &f.SizeMax)

	// Map float64 ranges
	setFloatRange(AverageLengthRange, &f.AverageLengthMin, &f.AverageLengthMax)
	setFloatRange(AverageEntropyRange, &f.AverageEntropyMin, &f.AverageEntropyMax)
	setFloatRange(DigitsPercentRange, &f.DigitsPercentMin, &f.DigitsPercentMax)
	setFloatRange(UpperCasePercentRange, &f.UpperCasePercentMin, &f.UpperCasePercentMax)
	setFloatRange(SpecialCharPercentRange, &f.SpecialCharPercentMin, &f.SpecialCharPercentMax)
	setFloatRange(DigitAndUpperCaseRange, &f.DigitAndUpperCaseMin, &f.DigitAndUpperCaseMax)
	setFloatRange(DigitAndSpecialCharRange, &f.DigitAndSpecialCharMin, &f.DigitAndSpecialCharMax)
	setFloatRange(UpperCaseAndSpecialCharRange, &f.UpperCaseAndSpecialCharMin, &f.UpperCaseAndSpecialCharMax)
	setFloatRange(DigitUpperCaseAndSpecialRange, &f.DigitUpperCaseAndSpecialMin, &f.DigitUpperCaseAndSpecialMax)

	return f
}

func (database *WordlistRepository) Filter(f WordlistFilter) ([]utils.Wordlist, error) {
	q := database.db.Model(&utils.Wordlist{})
	if f.name != nil { // wordlist name search (insensitive)
		q = q.Where("name LIKE ?", "%"+*f.name+"%")
	}
	if f.EntitiesMin != nil && f.EntitiesMax != nil { // searches entities range
		q = q.Where("entities BETWEEN ? AND ?", *f.EntitiesMin, *f.EntitiesMax)
	} else if f.EntitiesMin != nil && f.EntitiesMax == nil { // searches min entities
		q = q.Where("entities >= ?", *f.EntitiesMin)

	} else if f.EntitiesMin == nil && f.EntitiesMax != nil { // searches max entities
		q = q.Where("entities <= ?", *f.EntitiesMax)
	}
	if f.SizeMin != nil && f.SizeMax != nil {
		q = q.Where("size BETWEEN ? AND ?", *f.SizeMin, *f.SizeMax)
	}
	if f.Language != nil {
		q = q.Where("language = ?", *f.Language)
	}
	if f.Category != nil {
		q = q.Where("category = ?", *f.Category)
	}
	var out []utils.Wordlist
	err := q.Find(&out).Error
	return out, err
}
