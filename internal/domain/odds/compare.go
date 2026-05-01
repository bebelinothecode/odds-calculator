package odds

import "sort"

// BestByOutcome returns the best (highest) decimal price per outcome from a
// slice of quotes.
func BestByOutcome(quotes []Quote) map[Outcome]BestPrice {
	best := map[Outcome]BestPrice{}
	for _, q := range quotes {
		if b, ok := best[q.Outcome]; !ok || q.Price.Decimal > b.Decimal {
			best[q.Outcome] = BestPrice{Book: q.Price.Book, Decimal: q.Price.Decimal}
		}
	}
	return best
}

// SortedBooks returns the unique bookmaker keys present in the quotes, sorted
// alphabetically.
func SortedBooks(quotes []Quote) []string {
	seen := map[string]struct{}{}
	for _, q := range quotes {
		seen[q.Price.Book] = struct{}{}
	}
	books := make([]string, 0, len(seen))
	for b := range seen {
		books = append(books, b)
	}
	sort.Strings(books)
	return books
}
