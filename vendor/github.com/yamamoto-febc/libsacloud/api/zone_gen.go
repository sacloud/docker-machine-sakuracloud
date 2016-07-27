package api

/************************************************
  generated by IDE. for [ZoneAPI]
************************************************/

import (
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

/************************************************
   To support influent interface for Find()
************************************************/

func (api *ZoneAPI) Reset() *ZoneAPI {
	api.reset()
	return api
}

func (api *ZoneAPI) Offset(offset int) *ZoneAPI {
	api.offset(offset)
	return api
}

func (api *ZoneAPI) Limit(limit int) *ZoneAPI {
	api.limit(limit)
	return api
}

func (api *ZoneAPI) Include(key string) *ZoneAPI {
	api.include(key)
	return api
}

func (api *ZoneAPI) Exclude(key string) *ZoneAPI {
	api.exclude(key)
	return api
}

func (api *ZoneAPI) FilterBy(key string, value interface{}) *ZoneAPI {
	api.filterBy(key, value, false)
	return api
}

// func (api *ZoneAPI) FilterMultiBy(key string, value interface{}) *ZoneAPI {
// 	api.filterBy(key, value, true)
// 	return api
// }

func (api *ZoneAPI) WithNameLike(name string) *ZoneAPI {
	return api.FilterBy("Name", name)
}

func (api *ZoneAPI) WithTag(tag string) *ZoneAPI {
	return api.FilterBy("Tags.Name", tag)
}
func (api *ZoneAPI) WithTags(tags []string) *ZoneAPI {
	return api.FilterBy("Tags.Name", []interface{}{tags})
}

// func (api *ZoneAPI) WithSizeGib(size int) *ZoneAPI {
// 	api.FilterBy("SizeMB", size*1024)
// 	return api
// }

// func (api *ZoneAPI) WithSharedScope() *ZoneAPI {
// 	api.FilterBy("Scope", "shared")
// 	return api
// }

// func (api *ZoneAPI) WithUserScope() *ZoneAPI {
// 	api.FilterBy("Scope", "user")
// 	return api
// }

func (api *ZoneAPI) SortBy(key string, reverse bool) *ZoneAPI {
	api.sortBy(key, reverse)
	return api
}

func (api *ZoneAPI) SortByName(reverse bool) *ZoneAPI {
	api.sortByName(reverse)
	return api
}

// func (api *ZoneAPI) SortBySize(reverse bool) *ZoneAPI {
// 	api.sortBy("SizeMB", reverse)
// 	return api
// }

/************************************************
  To support CRUD(Create/Read/Update/Delete)
************************************************/

// func (api *ZoneAPI) New() *sacloud.Zone {
// 	return &sacloud.Zone{}
// }

// func (api *ZoneAPI) Create(value *sacloud.Zone) (*sacloud.Zone, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.create(api.createRequest(value), res)
// 	})
// }

func (api *ZoneAPI) Read(id string) (*sacloud.Zone, error) {
	return api.request(func(res *sacloud.Response) error {
		return api.read(id, nil, res)
	})
}

// func (api *ZoneAPI) Update(id string, value *sacloud.Zone) (*sacloud.Zone, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.update(id, api.createRequest(value), res)
// 	})
// }

// func (api *ZoneAPI) Delete(id string) (*sacloud.Zone, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.delete(id, nil, res)
// 	})
// }

/************************************************
  Inner functions
************************************************/

func (api *ZoneAPI) setStateValue(setFunc func(*sacloud.Request)) *ZoneAPI {
	api.baseAPI.setStateValue(setFunc)
	return api
}

func (api *ZoneAPI) request(f func(*sacloud.Response) error) (*sacloud.Zone, error) {
	res := &sacloud.Response{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.Zone, nil
}

func (api *ZoneAPI) createRequest(value *sacloud.Zone) *sacloud.Request {
	req := &sacloud.Request{}
	req.Zone = value
	return req
}
