package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"stockx-backend/reserr"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/logrusorgru/aurora"
	"github.com/rs/xid"
)

type Value struct {
	Verbose         bool
	RequestID       string
	RequestIDPrefix string
}

func NewValue(verbose bool, id string) Value {
	prefix := ""
	if id != "" {
		prefix = " " + aurora.Italic(aurora.BrightBlack(fmt.Sprintf("[%s]", id))).String()
	}

	return Value{Verbose: verbose, RequestID: id, RequestIDPrefix: prefix}
}

const (
	requestDataKey contextKey = iota
)

type contextKey int

func WithValue(ctx context.Context, v Value) context.Context {
	return context.WithValue(ctx, requestDataKey, v)
}

func GetValue(ctx context.Context) Value {
	if v, ok := ctx.Value(requestDataKey).(Value); ok {
		return v
	}

	return NewValue(false, xid.New().String())
}

func (v Value) Info(a ...interface{}) {
	log.Printf("%s%s %s\n", aurora.Blue("INFO"), v.RequestIDPrefix, fmt.Sprint(a...))
}

func (v Value) Infof(format string, a ...interface{}) {
	log.Printf("%s%s %s\n", aurora.Blue("INFO"), v.RequestIDPrefix, fmt.Sprintf(format, a...))
}

func (v Value) Debug(a ...interface{}) {
	log.Printf("%s%s %s\n", aurora.Cyan("DEBUG"), v.RequestIDPrefix, fmt.Sprint(a...))
}

func (v Value) Debugf(format string, a ...interface{}) {
	log.Printf("%s%s %s\n", aurora.Cyan("DEBUG"), v.RequestIDPrefix, fmt.Sprintf(format, a...))
}

func (v Value) Warning(a ...interface{}) {
	log.Printf("%s%s %s\n", aurora.BrightRed("WARNING"), v.RequestIDPrefix, fmt.Sprint(a...))
}

func (v Value) Warningf(format string, a ...interface{}) {
	log.Printf("%s %s %s\n", aurora.BrightRed("WARNING"), v.RequestIDPrefix, fmt.Sprintf(format, a...))
}

func (v Value) Error(a ...interface{}) {
	log.Printf("%s%s %s\n", aurora.Red("ERROR"), v.RequestIDPrefix, fmt.Sprint(a...))
}

func (v Value) Errorf(format string, a ...interface{}) {
	log.Printf("%s%s %s\n", aurora.Red("ERROR"), v.RequestIDPrefix, fmt.Sprintf(format, a...))
}

func (v Value) DebugRequestData(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if len(r.Form) > 0 {
		parameters := make([]string, 0)

		for k, vs := range r.Form {
			values := make([]string, 0)

			values = append(values, vs...)

			parameters = append(parameters, fmt.Sprintf("%s:[%s]", k, strings.Join(values, " ")))
		}

		if len(parameters) > 0 {
			v.Debugf("Request params: [%s]", strings.Join(parameters, " "))
		}
	}

	if r.Method != http.MethodGet {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		if string(body) != "" {
			v.Debugf("Request data: %s", string(body))
		}
	}

	return nil
}

type OnRespondDelegate func(w http.ResponseWriter, r *http.Request, code int, err error)

var OnRespondEvent OnRespondDelegate

func RespondWithJSON(w http.ResponseWriter, r *http.Request, code int, output interface{}, err error) {
	if OnRespondEvent != nil {
		OnRespondEvent(w, r, code, err)
	}

	if err != nil {
		SendJSONError(r.Context(), w, err)
		return
	}

	SendJSON(r.Context(), w, code, output)
}

func SendJSON(ctx context.Context, w http.ResponseWriter, status int, output interface{}, headers ...string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	b, err := json.Marshal(output)
	v := GetValue(ctx)

	if err != nil {
		v.Error("send json: ", err)
		w.WriteHeader(http.StatusInternalServerError)

		if _, err = w.Write([]byte(`{error:"internal",description:"error while marshaling JSON output"}`)); err != nil {
			v.Error("send json: ", err)
		}

		return
	}

	for i := 0; i < (len(headers) - 1); i += 2 {
		w.Header().Set(headers[i], headers[i+1])
	}
	w.WriteHeader(status)

	if _, err = w.Write(b); err != nil {
		v.Error("send json: ", err)
	}

	if v.Verbose {
		v.Debugf("output: %s", b)
	}
}

type StatusFailed struct {
	Status string `json:"status"`
}

func (s *StatusFailed) Error() string {
	return fmt.Sprintf("status %q", s.Status)
}

func SendJSONError(ctx context.Context, w http.ResponseWriter, err error, headers ...string) {
	v := GetValue(ctx)

	if err, ok := err.(*StatusFailed); ok {
		SendJSON(ctx, w, http.StatusOK, err)
		return
	}

	var output reserr.ErrorOutput
	if err, ok := err.(reserr.HTTPError); ok {
		output.Error = err.Name()
		output.Description = err.Error()
		output.Message = err.Message()
		v.Errorf("%s: %s", output.Error, output.Description)
		SendJSON(ctx, w, err.StatusCode(), &output, headers...)

		return
	}

	output.Error = "internal"
	output.Description = err.Error()
	v.Errorf("internal: %s", output.Description)
	SendJSON(ctx, w, http.StatusInternalServerError, &output, headers...)
}

func DecodeBodyAsJSON(w http.ResponseWriter, r *http.Request, input interface{}) bool {
	if r.Header.Get("Content-Type") != "application/json" {
		SendJSONError(r.Context(), w, reserr.NewHTTPError(http.StatusUnsupportedMediaType, "input", fmt.Errorf("only application/json content type is allowed"), `Dane wejściowe powinny być w formacie json!`))
		return false
	}

	return DecodeJSON(w, r, input)
}

func DecodeJSON(w http.ResponseWriter, req *http.Request, input interface{}) bool {
	v := GetValue(req.Context())
	b, err := ioutil.ReadAll(req.Body)

	if err != nil {
		v.Error("parsing input: ", err)
		RespondWithJSON(w, req, http.StatusBadRequest, struct{}{}, reserr.NewHTTPError(http.StatusBadRequest, "input", err, `Incorrect input data`))

		return false
	}

	defer func() { req.Body = ioutil.NopCloser(bytes.NewBuffer(b)) }()

	tmpBody := ioutil.NopCloser(bytes.NewBuffer(b))

	if v.Verbose {
		var b []byte
		b, err = ioutil.ReadAll(tmpBody)

		if err != nil {
			v.Infof("partial input: %s", b)
		} else {
			v.Infof("input: %s", b)
			err = json.Unmarshal(b, input)
		}
	} else {
		err = json.NewDecoder(tmpBody).Decode(input)
	}

	if err != nil {
		v.Error("parsing input: ", err)
		RespondWithJSON(w, req, http.StatusBadRequest, struct{}{}, reserr.NewHTTPError(http.StatusBadRequest, "input", err, `Incorrect input data`))

		return false
	}

	return true
}

const (
	ContentTypeJSON ContentType = "application/json"
	ContentTypeXML  ContentType = "application/xml"
)

type ContentType string

func SendRequest(method, address, endPoint string, params map[string]string, headers map[string]string, contentType *ContentType, body io.Reader, output interface{}) (int, http.Header, error) {
	u, err := url.Parse(address)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	u.Path = path.Join(u.Path, endPoint)

	if len(params) > 0 {
		uParams := url.Values{}
		for p, v := range params {
			uParams.Add(p, v)
		}

		u.RawQuery = uParams.Encode()
	}

	fmt.Printf("sending request %s %s", method, u.String())

	request, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if len(headers) > 0 {
		for h, v := range headers {
			request.Header.Set(h, v)
		}
	}

	if contentType != nil {
		request.Header.Set("Content-Type", string(*contentType))
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer response.Body.Close()

	fmt.Printf("received response status code %d", response.StatusCode)

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		errorMessage, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		return response.StatusCode, nil, fmt.Errorf("%s", string(errorMessage))
	}

	if output != nil {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		response.Body.Close()

		body = bytes.Replace(body, []byte(":NaN"), []byte(":null"), -1)
		body = bytes.Replace(body, []byte(": NaN"), []byte(": null"), -1)
		response.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.NewDecoder(response.Body).Decode(output)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
	}

	return response.StatusCode, response.Header, nil
}

func DecodeVarAsString(w http.ResponseWriter, r *http.Request, field string, data *string) bool {
	if !validValue(w, r, field, data) {
		return false
	}

	var err error

	vars := mux.Vars(r)
	*data = vars[field]

	if *data == "" {
		v := GetValue(r.Context())
		err = fmt.Errorf("bad %s value", field)
		v.Error("parsing input: ", err)
		RespondWithJSON(w, r, http.StatusBadRequest, struct{}{}, reserr.NewHTTPError(http.StatusBadRequest, "input", err, `Incorrect input data`))

		return false
	}

	return true
}

func DecodeVarAsInt64(w http.ResponseWriter, r *http.Request, field string, data *int64) bool {
	if !validValue(w, r, field, data) {
		return false
	}

	var err error

	vars := mux.Vars(r)
	*data, err = strconv.ParseInt(vars[field], 10, 64)

	if err != nil {
		v := GetValue(r.Context())
		err = fmt.Errorf("bad %s value", field)
		v.Error("parsing input: ", err)
		RespondWithJSON(w, r, http.StatusBadRequest, struct{}{}, reserr.NewHTTPError(http.StatusBadRequest, "input", err, `Incorrect input data`))

		return false
	}

	return true
}

func DecodeFormValueAsDate(w http.ResponseWriter, r *http.Request, format, fieldName string, out *time.Time) bool {
	var err error

	if !validValue(w, r, fieldName, out) {
		return false
	}

	value := r.FormValue(fieldName)
	if value != "" {
		*out, err = time.Parse(format, value)
		if err != nil {
			v := GetValue(r.Context())
			err = fmt.Errorf("bad %s value", fieldName)
			v.Error("parsing input: ", err)
			RespondWithJSON(w, r, http.StatusBadRequest, struct{}{}, reserr.NewHTTPError(http.StatusBadRequest, "input", err, `Incorrect input data`))

			return false
		}
	}

	return true
}

type Filters interface {
	BuildQuery(*gorm.DB, *Sort) *gorm.DB
}

func DecodeFormValueAsFilters(w http.ResponseWriter, r *http.Request, filters Filters) bool {
	if filters == nil {
		return true
	}

	v := GetValue(r.Context())

	if err := r.ParseForm(); err != nil {
		err = fmt.Errorf("bad params")
		v.Error("parsing input: ", err)
		RespondWithJSON(w, r, http.StatusBadRequest, struct{}{}, reserr.NewHTTPError(http.StatusBadRequest, "input", err, `Incorrect input params!`))

		return false
	}

	if len(r.Form) > 0 {
		err := decodeFilters(filters, r.Form)
		if err != nil {
			v.Error("parsing input: ", err)
			RespondWithJSON(w, r, http.StatusBadRequest, struct{}{}, reserr.NewHTTPError(http.StatusBadRequest, "input", err, `Incorrect filter params!`))

			return false
		}
	}

	return true
}

type Sort struct {
	Field string    `param:"field"`
	Order SortOrder `param:"order"`
}

const (
	SortOrderAscending  SortOrder = "asc"
	SortOrderDescending SortOrder = "desc"
)

type SortOrder string

func DecodeFormValueAsSort(w http.ResponseWriter, r *http.Request, sort **Sort) bool {
	v := GetValue(r.Context())

	if err := r.ParseForm(); err != nil {
		err = fmt.Errorf("bad params")
		v.Error("parsing input: ", err)
		RespondWithJSON(w, r, http.StatusBadRequest, struct{}{}, reserr.NewHTTPError(http.StatusBadRequest, "input", err, `Incorrect input params!`))

		return false
	}

	if len(r.Form) > 0 {
		formSort, err := decodeSort(r.Form)
		if err != nil {
			v.Error("parsing input: ", err)
			RespondWithJSON(w, r, http.StatusBadRequest, struct{}{}, reserr.NewHTTPError(http.StatusBadRequest, "input", err, `Nieprawidłowe parametry sortowania!`))

			return false
		}

		*sort = formSort
	}

	return true
}

func DecodeFormValueAsString(w http.ResponseWriter, r *http.Request, fieldName string, out *string) bool {
	if !validValue(w, r, fieldName, out) {
		return false
	}

	value := r.FormValue(fieldName)
	if value != "" {
		*out = value
	}

	return true
}

func DecodeFormValueAsURL(w http.ResponseWriter, r *http.Request, fieldName string, out **url.URL) bool {
	if !validValue(w, r, fieldName, out) {
		return false
	}

	var err error

	value := r.FormValue(fieldName)
	if value != "" {
		*out, err = url.Parse(value)
		if err != nil {
			v := GetValue(r.Context())
			err = fmt.Errorf("bad %s value", fieldName)
			v.Error("parsing input: ", err)
			RespondWithJSON(w, r, http.StatusBadRequest, struct{}{}, reserr.NewHTTPError(http.StatusBadRequest, "input", err, `Incorrect input data`))

			return false
		}
	}

	return true
}

func GetBodyData(w http.ResponseWriter, r *http.Request, data *[]byte) bool {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithJSON(w, r, http.StatusBadRequest, struct{}{}, reserr.NewHTTPError(http.StatusBadRequest, "input", err, `Incorrect input data`))
		return false
	}

	if data != nil {
		*data = body
	}

	return true
}
func DecodeQueryAsBool(w http.ResponseWriter, r *http.Request, field string, data *bool, defaultValue *bool) bool {
	fieldData := r.URL.Query().Get(field)
	if defaultValue != nil && fieldData == "" {
		*data = *defaultValue
		return true
	}

	var err error
	*data, err = strconv.ParseBool(r.URL.Query().Get(field))

	if err != nil {
		v := GetValue(r.Context())
		err = fmt.Errorf("bad %s value", field)
		v.Error("parsing input: ", err)
		RespondWithJSON(w, r, http.StatusBadRequest, struct{}{}, reserr.NewHTTPError(http.StatusBadRequest, "input", err, `Incorrect input data`))

		return false
	}

	return true
}

func validValue(w http.ResponseWriter, r *http.Request, field string, data interface{}) bool {
	if reflect.ValueOf(data).IsNil() {
		v := GetValue(r.Context())
		err := fmt.Errorf("nil %s value", field)
		v.Error("internal input error: ", err)
		RespondWithJSON(w, r, http.StatusBadRequest, struct{}{}, reserr.NewHTTPError(http.StatusInternalServerError, "input", err, `Błąd parsowania danych wejściowych!`))

		return false
	}

	return true
}

func decodeFilters(filters Filters, form url.Values) error {
	if filters == nil {
		return nil
	}

	if reflect.ValueOf(filters).Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer to filters")
	}

	m := make(map[string]interface{})

	for key, values := range form {
		if len(values) > 0 && strings.Contains(key, "filter.") {
			key = strings.Replace(key, "filter.", "", 1)

			t := reflect.TypeOf(filters).Elem()
			for i := 0; i < t.NumField(); i++ {
				fieldTag := t.Field(i).Tag.Get("param")
				if key == fieldTag {
					if reflect.ValueOf(filters).Elem().Field(i).Kind() == reflect.Slice {
						arrayValues := strings.Split(values[0], ",")

						if t.Field(i).Type.Elem().Kind() == reflect.Int {
							intValues := make([]int, 0)

							for y := 0; y < len(arrayValues); y++ {
								if arrayValues[y] == "" {
									continue
								}

								intValue, err := strconv.Atoi(arrayValues[y])
								if err != nil {
									return err
								}

								intValues = append(intValues, intValue)
							}

							m[key] = intValues

							break
						}

						m[key] = arrayValues

						break
					}

					switch reflect.ValueOf(filters).Elem().Field(i).Type() {
					case reflect.PtrTo(reflect.TypeOf(int64(0))):
						intValue, err := strconv.Atoi(values[0])
						if err != nil {
							return err
						}

						m[key] = intValue
					case reflect.PtrTo(reflect.TypeOf(false)):
						boolValue, err := strconv.ParseBool(values[0])
						if err != nil {
							return err
						}

						m[key] = boolValue
					default:
						m[key] = values[0]
					}
				}
			}
		}
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, filters); err != nil {
		return err
	}

	return nil
}

func decodeSort(form url.Values) (*Sort, error) {
	var sort Sort

	for key, values := range form {
		if len(values) > 0 && strings.Contains(key, "sort.") {
			key = strings.Replace(key, "sort.", "", 1)
			t := reflect.TypeOf(sort)

			for i := 0; i < t.NumField(); i++ {
				fieldTag := t.Field(i).Tag.Get("param")
				if key == fieldTag {
					v := reflect.ValueOf(&sort).Elem().FieldByName(t.Field(i).Name)
					if v.IsValid() && v.CanSet() {
						v.SetString(values[0])
						break
					}
				}
			}
		}
	}

	if sort.Field != "" || sort.Order != "" {
		if sort.Order != SortOrderAscending && sort.Order != SortOrderDescending {
			return nil, fmt.Errorf("invalid sort order")
		}

		return &sort, nil
	}

	return nil, nil
}

func ConvertJSONToRequestBody(data interface{}) (io.Reader, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
