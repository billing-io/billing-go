package billingio

// pageFunc fetches a single page of results. It receives the cursor for the
// page to fetch (nil for the first page) and returns the items, whether more
// pages exist, the next cursor, and any error.
type pageFunc[T any] func(cursor *string) (items []T, hasMore bool, nextCursor *string, err error)

// Iter is a generic auto-pagination iterator. It lazily fetches pages from the
// API as you advance through results.
//
// Usage:
//
//	iter := client.Checkouts.ListAutoPaginate(ctx, nil)
//	for iter.Next() {
//	    checkout := iter.Current()
//	    fmt.Println(checkout.CheckoutID)
//	}
//	if err := iter.Err(); err != nil {
//	    log.Fatal(err)
//	}
type Iter[T any] struct {
	fetch   pageFunc[T]
	items   []T
	index   int
	hasMore bool
	cursor  *string
	err     error
	started bool
}

// newIter creates a new Iter using the given page-fetching function.
func newIter[T any](fetch pageFunc[T]) *Iter[T] {
	return &Iter[T]{
		fetch:   fetch,
		hasMore: true, // assume there is at least one page
	}
}

// Next advances the iterator to the next item. It returns false when there are
// no more items or an error occurred. Call Current to get the item and Err to
// check for errors.
func (it *Iter[T]) Next() bool {
	if it.err != nil {
		return false
	}

	// If we haven't fetched the first page yet, or we exhausted the current
	// page and there are more pages to fetch, load the next page.
	if !it.started || (it.index >= len(it.items) && it.hasMore) {
		it.started = true
		items, hasMore, nextCursor, err := it.fetch(it.cursor)
		if err != nil {
			it.err = err
			return false
		}
		it.items = items
		it.hasMore = hasMore
		it.cursor = nextCursor
		it.index = 0

		if len(items) == 0 {
			return false
		}
		return true
	}

	// Advance within the current page.
	it.index++
	if it.index >= len(it.items) {
		if it.hasMore {
			// Recursively load the next page.
			return it.Next()
		}
		return false
	}
	return true
}

// Current returns the item at the current iterator position.
// It is only valid to call Current after a successful call to Next.
func (it *Iter[T]) Current() T {
	return it.items[it.index]
}

// Err returns the first error encountered during iteration, if any.
func (it *Iter[T]) Err() error {
	return it.err
}
