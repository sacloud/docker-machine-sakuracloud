package sakuracloud

import "strconv"

// ToSakuraID returns ID
func ToSakuraID(id string) (int64, bool) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, false
	}
	return i, true
}

// ToSakuraIDAll return IDs
func ToSakuraIDAll(ids []string) ([]int64, bool) {
	res := []int64{}
	for _, strID := range ids {
		id, r := ToSakuraID(strID)
		if !r {
			return []int64{}, false
		}
		res = append(res, id)
	}
	return res, true
}
