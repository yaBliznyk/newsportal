package portal

//go:generate colgen -imports=github.com/yaBliznyk/newsportal/internal/db
//colgen:News,Category,Tag
//colgen:News:UniqueTagIDs,Fill(Tags,TagIDs),MapP(db)
//colgen:Category:MapP(db)
//colgen:Tag:MapP(db)

// MapP converts slice of type T to slice of type M with given converter with pointers.
func MapP[T, M any](a []T, f func(*T) *M) []M {
	n := make([]M, len(a))
	for i := range a {
		n[i] = *f(&a[i])
	}
	return n
}

// Map converts slice of type T to slice of type M with given converter.
func Map[T, M any](a []T, f func(T) M) []M {
	n := make([]M, len(a))
	for i := range a {
		n[i] = f(a[i])
	}
	return n
}
