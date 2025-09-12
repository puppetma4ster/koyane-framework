package wordlistDB

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// WordlistRepository wraps a gorm.DB and provides
// methods for interacting with the wordlists database.
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

// WordlistFilter for compact filtering
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
	name, encoding, language, author, link, category *[]string,
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
		Category: category,
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

// Filter applies a WordlistFilter to the query and returns matching rows.
func (database *WordlistRepository) Filter(f *WordlistFilter) ([]utils.Wordlist, error) {
	q := database.db.Model(&utils.Wordlist{})

	// fuzzy name search (LIKE over multiple terms, case-insensitive on SQLite)
	q = applyStringListLike(q, "name", f.name, true)

	// exact matches with IN
	q = applyStringListExact(q, "encoding", f.Encoding)
	q = applyStringListExact(q, "language", f.Language)
	q = applyStringListExact(q, "category", f.Category)
	q = applyStringListExact(q, "author", f.Author)

	// optionally fuzzy on link
	q = applyStringListLike(q, "link", f.Link, true)

	// numeric ranges (uint64-backed)
	q = applyUintRange(q, "entities", f.EntitiesMin, f.EntitiesMax)
	q = applyUintRange(q, "smallestEntity", f.SmallestEntityMin, f.SmallestEntityMax)
	q = applyUintRange(q, "biggestEntity", f.BiggestEntityMin, f.BiggestEntityMax)
	q = applyUintRange(q, "size", f.SizeMin, f.SizeMax)

	// float ranges
	q = applyFloatRange(q, "averageLength", f.AverageLengthMin, f.AverageLengthMax)
	q = applyFloatRange(q, "averageEntropy", f.AverageEntropyMin, f.AverageEntropyMax)
	q = applyFloatRange(q, "digitsPercent", f.DigitsPercentMin, f.DigitsPercentMax)
	q = applyFloatRange(q, "upperCasePercent", f.UpperCasePercentMin, f.UpperCasePercentMax)
	q = applyFloatRange(q, "specialCharPercent", f.SpecialCharPercentMin, f.SpecialCharPercentMax)
	q = applyFloatRange(q, "digitAndUpperCase", f.DigitAndUpperCaseMin, f.DigitAndUpperCaseMax)
	q = applyFloatRange(q, "digitAndSpecialChar", f.DigitAndSpecialCharMin, f.DigitAndSpecialCharMax)
	q = applyFloatRange(q, "upperCaseAndSpecialChar", f.UpperCaseAndSpecialCharMin, f.UpperCaseAndSpecialCharMax)
	q = applyFloatRange(q, "digitUpperCaseAndSpecialChar", f.DigitUpperCaseAndSpecialMin, f.DigitUpperCaseAndSpecialMax)

	// tags: choose "any" or "all"
	q = applyTagsFilter(q, f.Tags, "any")

	var out []utils.Wordlist
	if err := q.Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

// ----------------------
// helper functions
// ----------------------

// applyStringListExact applies an exact-match IN filter: column IN (values...).
func applyStringListExact(q *gorm.DB, column string, listPtr *[]string) *gorm.DB {
	if listPtr == nil || len(*listPtr) == 0 {
		return q
	}
	return q.Where(column+" IN ?", *listPtr)
}

// applyStringListLike builds an OR chain of LIKE clauses for one column.
func applyStringListLike(q *gorm.DB, column string, listPtr *[]string, noCase bool) *gorm.DB {
	if listPtr == nil || len(*listPtr) == 0 {
		return q
	}
	terms := *listPtr

	// Filter out empty terms
	tmp := make([]string, 0, len(terms))
	for _, t := range terms {
		if t != "" {
			tmp = append(tmp, t)
		}
	}
	if len(tmp) == 0 {
		return q
	}

	col := column
	if noCase {
		col += " COLLATE NOCASE"
	}

	placeholders := make([]string, 0, len(tmp))
	args := make([]any, 0, len(tmp))
	for _, t := range tmp {
		placeholders = append(placeholders, col+" LIKE ?")
		args = append(args, "%"+t+"%")
	}
	expr := "(" + strings.Join(placeholders, " OR ") + ")"
	return q.Where(expr, args...)
}

// applyUintRange applies >= and/or <= for a uint64 range (nil = open bound).
func applyUintRange(q *gorm.DB, column string, min, max *uint64) *gorm.DB {
	if min != nil {
		q = q.Where(column+" >= ?", *min)
	}
	if max != nil {
		q = q.Where(column+" <= ?", *max)
	}
	return q
}

// applyFloatRange applies >= and/or <= for a float64 range (nil = open bound).
func applyFloatRange(q *gorm.DB, column string, min, max *float64) *gorm.DB {
	if min != nil {
		q = q.Where(column+" >= ?", *min)
	}
	if max != nil {
		q = q.Where(column+" <= ?", *max)
	}
	return q
}

// applyTagsFilter filters by tags using EXISTS subquery.
func applyTagsFilter(q *gorm.DB, tagsPtr *[]string, mode string) *gorm.DB {
	if tagsPtr == nil || len(*tagsPtr) == 0 {
		return q
	}
	tags := *tagsPtr

	switch mode {
	case "all":
		sub := q.Session(&gorm.Session{}).
			Model(&utils.WordlistTag{}).
			Select("COUNT(DISTINCT tag)").
			Where("wordlist_id = wordlists.id").
			Where("tag IN ?", tags)
		return q.Where("(?) = ?", sub, len(tags))
	default:
		return q.Where(`
			EXISTS (
				SELECT 1 FROM wordlist_tags wt
				WHERE wt.wordlist_id = wordlists.id
				  AND wt.tag IN ?
			)`, tags)
	}
}

// FindByID searches for a Wordlist by its ID.
// Returns the Wordlist pointer if found, or an error if not found or if the query fails.
func (r *WordlistRepository) FindByID(id uint64) (*utils.Wordlist, error) {
	var wl utils.Wordlist
	result := r.db.Preload("Tags").First(&wl, id) // Preload tags relation

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("ID not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &wl, nil
}

// DownloadWordlist downloads a file from the given URL into the specified folder.
// It shows a progress bar with percentage, downloaded size, and speed.
//
// Parameters:
//   - url: the file URL to download
//   - folder: the target folder where the file will be saved
//
// Returns:
//   - error: any error that occurs during the request, file creation, or download
func DownloadWordlist(url, folder, usrAgent string) error {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	// Set a custom User-Agent to avoid blocked requests
	req.Header.Set("User-Agent", usrAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if server responded with a success status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Get the content length (total size in bytes)
	size := resp.ContentLength
	if size <= 0 {
		fmt.Println("Warning: can't get content length, progress bar won't show correctly")
	}

	// Extract filename from URL and create destination file
	filename := path.Base(url)
	outFile, err := os.Create(path.Join(folder, filename))
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Track time and number of downloaded bytes
	start := time.Now()
	var downloaded int64

	// Buffer for reading chunks of data
	buf := make([]byte, 32*1024) // 32 KB buffer
	for {
		// Read a chunk from the response body
		n, err := resp.Body.Read(buf)
		if n > 0 {
			// Write chunk to file
			_, writeErr := outFile.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			// Update downloaded size
			downloaded += int64(n)

			// Print progress if total size is known
			if size > 0 {
				percent := float64(downloaded) / float64(size) * 100
				elapsed := time.Since(start).Seconds()
				speed := float64(downloaded) / 1024 / 1024 / elapsed // MB/s
				fmt.Printf("[*] \r%.2f%% [%.2f MB / %.2f MB] at %.2f MB/s",
					percent,
					float64(downloaded)/1024/1024,
					float64(size)/1024/1024,
					speed)
			}
		}
		// If end of file reached, stop reading
		if err == io.EOF {
			break
		}
		// Handle other errors
		if err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}
