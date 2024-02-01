package httputil

//type HttpUtil struct {
//}
//
//type MyRequest struct {
//	method  string
//	url     string
//	payLoad interface{}
//}

//func NewRequest(method, url string, payload interface{}) *MyRequest {
//	return &MyRequest{
//		method:  method,
//		url:     url,
//		payLoad: payload,
//	}
//}
//
//func (h HttpUtil) SendRequest(req *MyRequest, resp interface{}) error {
//	var (
//		body []byte
//		err  error
//	)
//
//	switch req.method {
//	case http.MethodGet:
//		body, err = h.handleGet(req)
//	case http.MethodPost:
//		body, err = h.handlePost(req)
//	default:
//		return errors.New("unsupported method: " + req.method)
//	}
//
//	if err != nil {
//		return err
//	}
//
//	return json.Unmarshal(body, resp)
//}
//
//func (h HttpUtil) handleGet(req *MyRequest) ([]byte, error) {
//	return get(h.genRequestUrl(req))
//}
//
//func (h HttpUtil) handlePost(req *MyRequest) ([]byte, error) {
//	return post(req.url, req.payLoad)
//}
//
//func (h HttpUtil) genRequestUrl(req *MyRequest) string {
//	rt := reflect.TypeOf(req.urlParam)
//
//	if rt.NumField() == 0 {
//		return ""
//	}
//
//	rv := reflect.ValueOf(req.urlParam)
//
//	keyValuePattern := make([]string, rt.NumField())
//
//	for i := 0; i < rt.NumField(); i++ {
//		field := rt.Field(i)
//		key := field.Tag.Get("json")
//		value := rv.Field(i).Interface()
//		keyValuePattern[i] = fmt.Sprintf("%s=%v", key, value)
//	}
//
//	log.Info("url: %s", fmt.Sprintf("%s?%s", req.url, strings.Join(keyValuePattern, "&")))
//
//	return fmt.Sprintf("%s?%s", req.url, strings.Join(keyValuePattern, "&"))
//}
