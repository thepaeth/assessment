package expenses

// func serverStart(t *testing.T)

// func seedExpense(t *testing.T) Expenses {
// 	var exp Expenses
// 	body := bytes.NewBufferString(`{
// 		"title": "strawberry smoothie",
// 		"amount": 79,
// 		"note": "night market promotion discount 10 bath",
// 		"tags": ["food", "beverage"]
// 	}`)

// 	err := request(http.MethodPost, uri("expenses"), body).Decode(&exp)
// 	if err != nil {
// 		t.Fatal("can't create expense:", err)
// 	}
// 	return exp
// }

// func uri(paths ...string) string {
// 	host := fmt.Sprint("http://localhost", os.Getenv("PORT"))
// 	if paths == nil {
// 		return host
// 	}

// 	url := append([]string{host}, paths...)
// 	return strings.Join(url, "/")
// }

// type Response struct {
// 	*http.Response
// 	err error
// }

// func (r *Response) Decode(v interface{}) error {
// 	if r.err != nil {
// 		return r.err
// 	}

// 	return json.NewDecoder(r.Body).Decode(v)
// }

// func request(method, url string, body io.Reader) *Response {
// 	req, _ := http.NewRequest(method, url, body)
// 	req.Header.Add("Authorization", "November 10, 2009")
// 	req.Header.Add("Content-Type", "application/json")
// 	client := http.Client{}
// 	res, err := client.Do(req)
// 	return &Response{res, err}
// }
