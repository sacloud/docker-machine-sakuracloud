package sakuracloud

import "strconv"

func ToSakuraID(id string) (int64, bool) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, false
	}
	return i, true
}

func ToSakuraIDAll(ids []string) ([]int64, bool) {
	res := []int64{}
	for _, strId := range ids {
		id, r := ToSakuraID(strId)
		if !r {
			return []int64{}, false
		}
		res = append(res, id)
	}
	return res, true
}
