package utils

func MakeBucket(id int64) int {
	BUCKET_SIZE := 1000 * 60 * 60 * 24 * 10 // 10 days
	t := id >> 22

	return int(t / int64(BUCKET_SIZE))
}
