package sqltests

import (
	"net/http"
	"net/url"
	"strings"
)

var (
	crudTests = map[string][]testCase{
		"GetSingle":     getSingleTestCases,
		"GetCollection": getCollectionTestCases,
		"CreateSingle":  createSingleTestCases,
		"UpdateSingle":  updateSingleTestCases,
		"DeleteSingle":  deleteSingleTestCases,
	}
	getSingleTestCases = []testCase{
		{
			Kind:                "success",
			Name:                "it succeeds when equipment table contains more than one column and returns associated recipes normalized",
			ExpectedStatusCodes: []int{http.StatusOK},
			ExpectedResults: []string{`{
  "acquisitionCost": {
    "amount": 8477.85,
    "currency": "EUR"
  },
  "id": "2",
  "name": "Sanremo Café Racer",
  "purchaseDate": "2017-12-12T12:00:00.123Z",
  "recipes": [
    "1"
  ]
}`},
			Request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/equipment/2", nil)
				return req
			},
		},
		{
			Kind:                "success",
			Name:                "it succeeds when equipment table contains more than one column and returns associated recipes denormalized when provided with the query option",
			ExpectedStatusCodes: []int{http.StatusOK},
			ExpectedResults: []string{`{
  "acquisitionCost": {
    "amount": 8477.85,
    "currency": "EUR"
  },
  "id": "2",
  "name": "Sanremo Café Racer",
  "purchaseDate": "2017-12-12T12:00:00.123Z",
  "recipes": [
    {
      "creationDate": "2017-12-13T00:00:00.000Z",
      "equipmentId": "2",
      "id": "1",
      "instructions": "do this",
      "lastAccessed": "%sT00:00:01.000Z",
      "lastModified": "2017-12-14T00:00:00.123Z",
      "name": "Espresso single shot"
    }
  ]
}`, `.*`},
			Request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/equipment/2?denormalize=true", nil)
				return req
			},
		},
		{

			Kind:                "failure",
			Name:                "it fails and returns 404 NOT FOUND when querying a non existent equipment id",
			ExpectedStatusCodes: []int{http.StatusNotFound},
			ExpectedResults: []string{`{
  "status": {
    "code": 404,
    "description": "Resource with uniqueID '42' not found in equipment table"
  }
}`},
			Request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/equipment/42", nil)
				return req
			},
		},
		{
			Kind:                "success",
			Name:                "it succeeds when recipes table contains more than one column and returns associated equipment normalized",
			ExpectedStatusCodes: []int{http.StatusOK},
			ExpectedResults: []string{`{
  "creationDate": "2017-12-13T00:00:00.000Z",
  "equipment": "2",
  "equipmentId": "2",
  "id": "1",
  "instructions": "do this",
  "lastAccessed": "%sT00:00:01.000Z",
  "lastModified": "2017-12-14T00:00:00.123Z",
  "name": "Espresso single shot"
}`, `.*`},
			Request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/recipes/1", nil)
				return req
			},
		},
		{
			Kind:                "success",
			Name:                "it succeeds when recipes table contains more than one column and returns associated equipment denormalized when provided with the query option",
			ExpectedStatusCodes: []int{http.StatusOK},
			ExpectedResults: []string{`{
  "creationDate": "2017-12-13T00:00:00.000Z",
  "equipment": {
    "acquisitionCost": {
      "amount": 8477.85,
      "currency": "EUR"
    },
    "id": "2",
    "name": "Sanremo Café Racer",
    "purchaseDate": "2017-12-12T12:00:00.123Z"
  },
  "equipmentId": "2",
  "id": "1",
  "instructions": "do this",
  "lastAccessed": "%sT00:00:01.000Z",
  "lastModified": "2017-12-14T00:00:00.123Z",
  "name": "Espresso single shot"
}`, `.*`},
			Request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/recipes/1?denormalize=true", nil)
				return req
			},
		},
	}
	getCollectionTestCases = []testCase{
		{
			Kind:                "success",
			Name:                "it returns 200 OK with empty array when querying an empty table",
			ExpectedStatusCodes: []int{http.StatusOK},
			ExpectedResults:     []string{`[]`},
			Request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/zeroRows", nil)
				return req
			},
		},
		{
			Kind:                "success",
			Name:                "it returns 200 OK with a single result in an array when querying a table containing one row",
			ExpectedStatusCodes: []int{http.StatusOK},
			ExpectedResults: []string{`[
  {
    "id": "1",
    "name": "TESTNAME"
  }
]`},
			Request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/oneRows", nil)
				return req
			},
		},
		{
			Kind:                "success",
			Name:                "it succeeds when equipment table contains more than one column",
			ExpectedStatusCodes: []int{http.StatusOK},
			ExpectedResults: []string{`[
  {
    "acquisitionCost": {
      "amount": 25.95,
      "currency": "EUR"
    },
    "id": "1",
    "name": "Bialetti Moka Express 6 cup",
    "purchaseDate": "2017-12-11T12:00:00.123Z"
  },
  {
    "acquisitionCost": {
      "amount": 8477.85,
      "currency": "EUR"
    },
    "id": "2",
    "name": "Sanremo Café Racer",
    "purchaseDate": "2017-12-12T12:00:00.123Z"
  },
  {
    "acquisitionCost": {
      "amount": 39.95,
      "currency": "EUR"
    },
    "id": "3",
    "name": "Buntfink SteelKettle",
    "purchaseDate": "2017-12-12T12:00:00.000Z"
  },
  {
    "acquisitionCost": {
      "amount": 49.95,
      "currency": "EUR"
    },
    "id": "4",
    "name": "Copper Coffee Pot Cezve",
    "purchaseDate": "2017-12-12T12:00:00.000Z"
  }
]`},
			Request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/equipment", nil)
				return req
			},
		},
	}
	createSingleTestCases = []testCase{
		{
			Kind: "success",
			Name: "it returns a 200 OK with the newly created resource or a 204 No Content when provided with valid URL parameters on POST",
			ExpectedResults: []string{`%s{
  "acquisitionCost": {
    "amount": 35.99,
    "currency": "EUR"
  },
  "id": "5",
  "name": "French Press",
  "purchaseDate": "2017-04-02T00:00:00.000Z",
  "recipes": []
}%s`, `(`, `|^$)`},
			ExpectedStatusCodes: []int{http.StatusCreated, http.StatusNoContent},
			ExpectedHeader: http.Header(map[string][]string{
				"Location": []string{"/equipment/5"},
			}),
			Request: func() *http.Request {
				postData := url.Values{}
				postData.Set("id", "5")
				postData.Set("name", "French Press")
				postData.Set("acquisitionCost", "35.99")
				postData.Set("purchaseDate", "2017-04-02T00:00:00.000Z")
				req, _ := http.NewRequest("POST", "/equipment", strings.NewReader(postData.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			},
		},
		{
			Kind: "success",
			Name: "it returns a 200 OK with the newly created resource or a 204 No Content when creating a new resource in an empty table",
			ExpectedResults: []string{`%s{
  "id": "1",
  "name": "Graef CM800 Coffee Burr Grinder"
}%s`, `(`, `|^$)`},
			ExpectedStatusCodes: []int{http.StatusCreated, http.StatusNoContent},
			ExpectedHeader: http.Header(map[string][]string{
				"Location": []string{"/zeroRows/1"},
			}),
			Request: func() *http.Request {
				postData := url.Values{}
				postData.Set("id", "1")
				postData.Set("name", "Graef CM800 Coffee Burr Grinder")
				req, _ := http.NewRequest("POST", "/zeroRows", strings.NewReader(postData.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			},
		},
	}
	updateSingleTestCases = []testCase{
		{
			Kind:                "success",
			Name:                "it succeeds when provided with valid parameters as URL parameters",
			ExpectedStatusCodes: []int{http.StatusOK},
			ExpectedResults: []string{`{
  "acquisitionCost": {
    "amount": 9283.99,
    "currency": "EUR"
  },
  "id": "2",
  "name": "Sanremo Café Racer",
  "purchaseDate": "2017-12-01T12:34:56.789Z",
  "recipes": [
    "1"
  ]
}`},
			Request: func() *http.Request {
				postData := url.Values{}
				postData.Set("name", "Sanremo Café Racer")
				postData.Set("acquisitionCost", "9283.99")
				postData.Set("purchaseDate", "2017-12-01T12:34:56.789Z")
				req, _ := http.NewRequest("PATCH", "/equipment/2?"+postData.Encode(), nil)
				return req
			},
		},
		{
			Kind:                "success",
			Name:                "it succeeds when user explicitely wants to insert a null value",
			ExpectedStatusCodes: []int{http.StatusOK},
			ExpectedResults: []string{`{
  "acquisitionCost": {
    "amount": 8477.85,
    "currency": "EUR"
  },
  "id": "2",
  "name": "Sanremo Café Racer",
  "purchaseDate": %s,
  "recipes": [
    "1"
  ]
}`, `(null|"0001-01-01T00:00:00.000Z")`},
			Request: func() *http.Request {
				body := strings.NewReader(
					`{"name": "Sanremo Café Racer", "acquisitionCost": 8477.85, "purchaseDate": null}`,
				)
				req, _ := http.NewRequest(
					"PATCH",
					"/equipment/2",
					body,
				)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
		},
		{
			Kind:                "success",
			Name:                "it succeeds when provided with valid parameters as json in the request body",
			ExpectedStatusCodes: []int{http.StatusOK},
			ExpectedResults: []string{`{
  "acquisitionCost": {
    "amount": 8477.85,
    "currency": "EUR"
  },
  "id": "2",
  "name": "Sanremo Café Racer",
  "purchaseDate": "2017-12-12T12:00:00.123Z",
  "recipes": [
    "1"
  ]
}`},
			Request: func() *http.Request {
				body := strings.NewReader(
					`{"name": "Sanremo Café Racer", "acquisitionCost": 8477.85, "purchaseDate": "2017-12-12T12:00:00.123Z"}`,
				)
				req, _ := http.NewRequest(
					"PATCH",
					"/equipment/2",
					body,
				)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
		},
		{

			Kind:                "failure",
			Name:                "it fails and returns 404 NOT FOUND when trying to update a non existent id",
			ExpectedStatusCodes: []int{http.StatusNotFound},
			ExpectedResults: []string{`{
  "status": {
    "code": 404,
    "description": "Resource with uniqueID '42' not found in equipment table"
  }
}
`},
			Request: func() *http.Request {
				postData := url.Values{}
				postData.Set("name", "Sanremo Café Racer")
				postData.Set("acquisitionCost", "512.23")
				req, _ := http.NewRequest("PATCH", "/equipment/42?"+postData.Encode(), nil)
				return req
			},
		},
	}
	deleteSingleTestCases = []testCase{
		{
			Kind:                "success",
			Name:                "it succeeds in deleting an existing resource",
			ExpectedStatusCodes: []int{http.StatusOK},
			ExpectedResults: []string{`{
  "status": {
    "code": 200,
    "description": "Resource with uniqueID '5' successfully deleted from equipment table"
  }
}`},
			Request: func() *http.Request {
				req, _ := http.NewRequest("DELETE", "/equipment/5", nil)
				return req
			},
		},
		{
			Kind:                "success",
			Name:                "it succeeds in deleting an existing resource from the zeroRows table",
			ExpectedStatusCodes: []int{http.StatusOK},
			ExpectedResults: []string{`{
  "status": {
    "code": 200,
    "description": "Resource with uniqueID '1' successfully deleted from zero_rows table"
  }
}`},
			Request: func() *http.Request {
				req, _ := http.NewRequest("DELETE", "/zeroRows/1", nil)
				return req
			},
		},
		{

			Kind:                "failure",
			Name:                "it fails and returns 404 NOT FOUND when trying to delete a non existent id",
			ExpectedStatusCodes: []int{http.StatusNotFound},
			ExpectedResults: []string{`{
  "status": {
    "code": 404,
    "description": "Resource with uniqueID '42' not found in equipment table"
  }
}`},
			Request: func() *http.Request {
				req, _ := http.NewRequest("DELETE", "/equipment/42", nil)
				return req
			},
		},
	}
)
