package main

import (
	"fmt"
	"strconv"
	"strings"
)

// FilterSet is a function that applies a set of filters and returns the filtered records.
type FilterSet func([]string) []string

// Filter is a filter function applied to a single record.
type Filter func(string) bool

// FilterBulk is a bulk filter function applied to an entire slice of records.
type FilterBulk func([]string) []string

var filters = map[int]FilterSet{
	1: FilterForAnimals,
	2: FilterForIDs,
}

func main() {
	// Initialize some contrived records.
	records := []string{
		"Cat",
		"A sentence is not a valid record.",
		"Minotaur",
		"cd5169bf-3649-4091-862b-c7ec1de92fd9-cd5169bf-3649-4091-862b-c7ec1de92fd9-cd5169bf-3649-4091-862b-c7ec1de92fd9",
		"3412-3241",
		"Dragon",
		"Cat",
	}

	// Call the filter functions.
	// ApplyFilters will be applied first, in order from top to bottom.
	animals := filters[1](records)
	ids := filters[2](records)

	// The only thing that should be left is one record of "Cat".
	fmt.Println("Animals:")
	for _, animal := range animals {
		fmt.Println(animal)
	}

	// The only thing left should be the integer.
	fmt.Println("IDs:")
	for _, id := range ids {
		fmt.Println(id)
	}
}

// FilterForAnimals applies a set of filters removing any non-animals.
func FilterForAnimals(records []string) []string {
	return ApplyBulkFilters(
		ApplyFilters(records,
			FilterMagicalCreatures,
			FilterStringLength,
			FilterInts,
			FilterWords,
		),
		FilterDuplicates,
	)
}

// FilterForIDs applies a set of filters removing any non-IDs.
func FilterForIDs(records []string) []string {
	return ApplyBulkFilters(
		ApplyFilters(records,
			FilterIDs,
		),
		FilterDuplicates,
	)
}

// ApplyFilters applies a set of filters to a record list.
// Each record will be checked against each filter.
// The filters are applied in the order they are passed in.
func ApplyFilters(records []string, filters ...Filter) []string {
	// Make sure there are actually filters to be applied.
	if len(filters) == 0 {
		return records
	}

	filteredRecords := make([]string, 0, len(records))

	// Range over the records and apply all the filters to each record.
	// If the record passes all the filters, add it to the final slice.
	for _, r := range records {
		keep := true

		for _, f := range filters {
			if !f(r) {
				keep = false
				break
			}
		}

		if keep {
			filteredRecords = append(filteredRecords, r)
		}
	}

	return filteredRecords
}

// ApplyBulkFilters applies a set of filters to the entire slice of records.
// Used when each record filter requires knowledge of the other records, e.g. de-duping.
func ApplyBulkFilters(records []string, filters ...FilterBulk) []string {
	for _, f := range filters {
		records = f(records)
	}

	return records
}

// FilterDuplicates is a bulk filter to remove any duplicates from the set.
func FilterDuplicates(records []string) []string {
	recordMap := map[string]bool{}
	filteredRecords := []string{}

	for _, record := range records {
		if ok := recordMap[record]; ok {
			continue
		}
		recordMap[record] = true
		filteredRecords = append(filteredRecords, record)
	}

	return filteredRecords
}

// FilterMagicalCreatures filters out common mythical creatures.
func FilterMagicalCreatures(record string) bool {
	magicalCreatures := []string{
		"Unicorn",
		"Dragon",
		"Griffin",
		"Minotaur",
	}

	for _, c := range magicalCreatures {
		if record == c {
			return false
		}
	}

	return true
}

// FilterStringLength removes any long records.
func FilterStringLength(record string) bool {
	if len(record) > 75 {
		return false
	}

	return true
}

// FilterInts removes an integers disguised as strings.
func FilterInts(record string) bool {
	if _, err := strconv.Atoi(record); err == nil {
		return false
	}

	return true
}

// FilterWords removes any records with spaces.
func FilterWords(record string) bool {
	split := strings.Split(record, " ")
	if len(split) > 1 {
		return false
	}

	return true
}

// FilterIDs removes any IDs that don't match contrived criteria.
func FilterIDs(record string) bool {
	split := strings.Split(record, "-")
	if len(split) != 2 {
		return false
	}

	if _, err := strconv.Atoi(split[0]); err == nil {
		return false
	}
	if _, err := strconv.Atoi(split[1]); err == nil {
		return false
	}

	return true
}
