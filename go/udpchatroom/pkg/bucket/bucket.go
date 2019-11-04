package bucket

type Bucket struct {
	room map[int64]*Channel
}

func NewBucket() *Bucket {
	return &Bucket{room: make(map[int64]*Channel)}
}
