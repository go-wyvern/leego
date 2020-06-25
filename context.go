package leego

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-wyvern/leego/engine"
	"github.com/go-wyvern/logger"

	"golang.org/x/net/context"
)

type (
	// Context represents the context of the current HTTP request. It holds request and
	// response objects, path, path parameters, data and registered handler.
	Context interface {
		// Context returns `net/context.Context`.
		Context() context.Context

		// SetContext sets `net/context.Context`.
		SetContext(context.Context)

		// Deadline returns the time when work done on behalf of this context
		// should be canceled.  Deadline returns ok==false when no deadline is
		// set.  Successive calls to Deadline return the same results.
		Deadline() (deadline time.Time, ok bool)

		// Done returns a channel that's closed when work done on behalf of this
		// context should be canceled.  Done may return nil if this context can
		// never be canceled.  Successive calls to Done return the same value.
		Done() <-chan struct{}

		// Err returns a non-nil error value after Done is closed.  Err returns
		// Canceled if the context was canceled or DeadlineExceeded if the
		// context's deadline passed.  No other values for Err are defined.
		// After Done is closed, successive calls to Err return the same value.
		Err() error

		// Value returns the value associated with this context for key, or nil
		// if no value is associated with key.  Successive calls to Value with
		// the same key returns the same result.
		Value(key interface{}) interface{}

		// Request returns `engine.Request` interface.
		Request() engine.Request

		// Request returns `engine.Response` interface.
		Response() engine.Response

		// Path returns the registered path for the handler.
		Path() string

		// SetPath sets the registered path for the handler.
		SetPath(string)

		// P returns path parameter by index.
		P(int) string

		// Param returns path parameter by name.
		Param(string) string

		// ParamNames returns path parameter names.
		ParamNames() []string

		// SetParamNames sets path parameter names.
		SetParamNames(...string)

		// ParamValues returns path parameter values.
		ParamValues() []string

		// SetParamValues sets path parameter values.
		SetParamValues(...string)

		// QueryParam returns the query param for the provided name. It is an alias
		// for `engine.URL#QueryParam()`.
		QueryParam(string) string

		// QueryParams returns the query parameters as map.
		// It is an alias for `engine.URL#QueryParams()`.
		QueryParams() map[string][]string

		// FormValue returns the form field value for the provided name. It is an
		// alias for `engine.Request#FormValue()`.
		FormValue(string) string

		// FormParams returns the form parameters as map.
		// It is an alias for `engine.Request#FormParams()`.
		FormParams() map[string][]string

		// FormFile returns the multipart form file for the provided name. It is an
		// alias for `engine.Request#FormFile()`.
		FormFile(string) (*multipart.FileHeader, error)

		// MultipartForm returns the multipart form.
		// It is an alias for `engine.Request#MultipartForm()`.
		MultipartForm() (*multipart.Form, error)

		// Cookie returns the named cookie provided in the request.
		// It is an alias for `engine.Request#Cookie()`.
		Cookie(string) (engine.Cookie, error)

		// SetCookie adds a `Set-Cookie` header in HTTP response.
		// It is an alias for `engine.Response#SetCookie()`.
		SetCookie(engine.Cookie)

		// Cookies returns the HTTP cookies sent with the request.
		// It is an alias for `engine.Request#Cookies()`.
		Cookies() []engine.Cookie

		// Get retrieves data from the context.
		Get(string) interface{}

		// Set saves data in the context.
		Set(string, interface{})

		// Bind binds the request body into provided type `i`. The default binder
		// does it based on Content-Type header.
		Bind(interface{}) error

		// Render renders a template with data and sends a text/html response with status
		// code. Templates can be registered using `leego.SetRenderer()`.
		//Render(int, string, interface{}) error

		// HTML sends an HTTP response with status code.
		HTML(int, string) error

		// String sends a string response with status code.
		String(int, string) error

		// JSON sends a JSON response with status code.
		JSON(int, interface{}) error

		// JSONBlob sends a JSON blob response with status code.
		JSONBlob(int, []byte) error

		// JSONP sends a JSONP response with status code. It uses `callback` to construct
		// the JSONP payload.
		JSONP(int, string, interface{}) error

		// XML sends an XML response with status code.
		XML(int, interface{}) error

		// XMLBlob sends a XML blob response with status code.
		XMLBlob(int, []byte) error

		// File sends a response with the content of the file.
		File(string) error

		// Attachment sends a response from `io.ReaderSeeker` as attachment, prompting
		// client to save the file.
		Attachment(io.ReadSeeker, string) error

		// NoContent sends a response with no body and a status code.
		NoContent(int) error

		// Redirect redirects the request with status code.
		Redirect(int, string) error

		// Error invokes the registered HTTP error handler. Generally used by middleware.
		Error(err error)

		// Handler returns the matched handler by router.
		Handler() HandlerFunc

		// SetHandler sets the matched handler by router.
		SetHandler(HandlerFunc)

		SetParamsMap(m map[string]string)

		GetParamsMap() map[string]string

		// Logger returns the `Logger` instance.
		Logger() *logger.Logger

		// leego returns the `leego` instance.
		Leego() *Leego

		SetLogger(*logger.Logger)

		// ServeContent sends static content from `io.Reader` and handles caching
		// via `If-Modified-Since` request header. It automatically sets `Content-Type`
		// and `Last-Modified` response headers.
		ServeContent(io.ReadSeeker, string, time.Time) error

		// Reset resets the context after request completes. It must be called along
		// with `leego#AcquireContext()` and `leego#ReleaseContext()`.
		// See `leego#ServeHTTP()`
		Reset(engine.Request, engine.Response)

		SetData(string, interface{})

		GetData(string) interface{}

		Language() string

		SetLang(string)
	}

	leegoContext struct {
		context   context.Context
		request   engine.Request
		response  engine.Response
		logger    *logger.Logger
		path      string
		pnames    []string
		pvalues   []string
		paramsMap map[string]string
		handler   HandlerFunc
		leego     *Leego
		lang      string
		data      map[string]interface{}
	}
)

var _ Context = new(leegoContext)

func (c *leegoContext) Language() string {
	return c.lang
}

func (c *leegoContext) SetLang(lang string) {
	if lang != "" && len(lang) >= 5 {
		lang = lang[:5]
	} else {
		lang = "zh-CN"
	}
	c.lang = lang
}

func (c *leegoContext) SetParamsMap(m map[string]string) {
	c.paramsMap = m
}

func (c *leegoContext) Logger() *logger.Logger {
	if c.logger != nil {
		return c.logger
	}
	return c.leego.logger
}

func (c *leegoContext) SetLogger(l *logger.Logger) {
	c.logger = l
}

func (c *leegoContext) GetParamsMap() map[string]string {
	return c.paramsMap
}

func (c *leegoContext) SetData(key string, data interface{}) {
	c.data[key] = data
}

func (c *leegoContext) GetData(key string) interface{} {
	return c.data[key]
}

func (c *leegoContext) Context() context.Context {
	return c.context
}

func (c *leegoContext) SetContext(ctx context.Context) {
	c.context = ctx
}

func (c *leegoContext) Deadline() (deadline time.Time, ok bool) {
	return c.context.Deadline()
}

func (c *leegoContext) Done() <-chan struct{} {
	return c.context.Done()
}

func (c *leegoContext) Err() error {
	return c.context.Err()
}

func (c *leegoContext) Value(key interface{}) interface{} {
	return c.context.Value(key)
}

func (c *leegoContext) Request() engine.Request {
	return c.request
}

func (c *leegoContext) Response() engine.Response {
	return c.response
}

func (c *leegoContext) Path() string {
	return c.path
}

func (c *leegoContext) SetPath(p string) {
	c.path = p
}

func (c *leegoContext) P(i int) (value string) {
	l := len(c.pnames)
	if i < l {
		value = c.pvalues[i]
	}
	return
}

func (c *leegoContext) Param(name string) (value string) {
	l := len(c.pnames)
	for i, n := range c.pnames {
		if n == name && i < l {
			value = c.pvalues[i]
			break
		}
	}
	return
}

func (c *leegoContext) ParamNames() []string {
	return c.pnames
}

func (c *leegoContext) SetParamNames(names ...string) {
	c.pnames = names
}

func (c *leegoContext) ParamValues() []string {
	return c.pvalues
}

func (c *leegoContext) SetParamValues(values ...string) {
	c.pvalues = values
}

func (c *leegoContext) QueryParam(name string) string {
	return c.request.URL().QueryParam(name)
}

func (c *leegoContext) QueryParams() map[string][]string {
	return c.request.URL().QueryParams()
}

func (c *leegoContext) FormValue(name string) string {
	return c.request.FormValue(name)
}

func (c *leegoContext) FormParams() map[string][]string {
	return c.request.FormParams()
}

func (c *leegoContext) FormFile(name string) (*multipart.FileHeader, error) {
	return c.request.FormFile(name)
}

func (c *leegoContext) MultipartForm() (*multipart.Form, error) {
	return c.request.MultipartForm()
}

func (c *leegoContext) Cookie(name string) (engine.Cookie, error) {
	return c.request.Cookie(name)
}

func (c *leegoContext) SetCookie(cookie engine.Cookie) {
	c.response.SetCookie(cookie)
}

func (c *leegoContext) Cookies() []engine.Cookie {
	return c.request.Cookies()
}

func (c *leegoContext) Set(key string, val interface{}) {
	c.context = context.WithValue(c.context, key, val)
}

func (c *leegoContext) Get(key string) interface{} {
	return c.context.Value(key)
}

func (c *leegoContext) Bind(i interface{}) error {
	return c.leego.binder.Bind(i, c)
}

//func (c *leegoContext) Render(code int, name string, data interface{}) (err error) {
//	if c.leego.renderer == nil {
//		return ErrRendererNotRegistered
//	}
//	buf := new(bytes.Buffer)
//	if err = c.leego.renderer.Render(buf, name, data, c); err != nil {
//		return
//	}
//	c.response.Header().Set(HeaderContentType, MIMETextHTMLCharsetUTF8)
//	c.response.WriteHeader(code)
//	_, err = c.response.Write(buf.Bytes())
//	return
//}

func (c *leegoContext) HTML(code int, html string) (err error) {
	c.response.Header().Set(HeaderContentType, MIMETextHTMLCharsetUTF8)
	c.response.WriteHeader(code)
	_, err = c.response.Write([]byte(html))
	return
}

func (c *leegoContext) String(code int, s string) (err error) {
	c.response.Header().Set(HeaderContentType, MIMETextPlainCharsetUTF8)
	c.response.WriteHeader(code)
	_, err = c.response.Write([]byte(s))
	return
}

func (c *leegoContext) JSON(code int, i interface{}) (err error) {
	b, err := json.Marshal(i)
	c.Response().SetBody(string(b))
	//if c.leego.Debug() {
	//	b, err = json.MarshalIndent(i, "", "  ")
	//}
	if err != nil {
		return err
	}
	return c.JSONBlob(code, b)
}

func (c *leegoContext) JSONBlob(code int, b []byte) (err error) {
	c.response.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	c.response.WriteHeader(code)
	_, err = c.response.Write(b)
	return
}

func (c *leegoContext) JSONP(code int, callback string, i interface{}) (err error) {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	c.response.Header().Set(HeaderContentType, MIMEApplicationJavaScriptCharsetUTF8)
	c.response.WriteHeader(code)
	if _, err = c.response.Write([]byte(callback + "(")); err != nil {
		return
	}
	if _, err = c.response.Write(b); err != nil {
		return
	}
	_, err = c.response.Write([]byte(");"))
	return
}

func (c *leegoContext) XML(code int, i interface{}) (err error) {
	b, err := xml.Marshal(i)
	c.Response().SetBody(string(b))
	//if c.leego.Debug() {
	//	b, err = xml.MarshalIndent(i, "", "  ")
	//}
	if err != nil {
		return err
	}
	return c.XMLBlob(code, b)
}

func (c *leegoContext) XMLBlob(code int, b []byte) (err error) {
	c.response.Header().Set(HeaderContentType, MIMEApplicationXMLCharsetUTF8)
	c.response.WriteHeader(code)
	if _, err = c.response.Write([]byte(xml.Header)); err != nil {
		return
	}
	_, err = c.response.Write(b)
	return
}

func (c *leegoContext) File(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return ErrNotFound
	}
	defer f.Close()

	fi, _ := f.Stat()
	if fi.IsDir() {
		file = filepath.Join(file, "index.html")
		f, err = os.Open(file)
		if err != nil {
			return ErrNotFound
		}
		if fi, err = f.Stat(); err != nil {
			return err
		}
	}
	return c.ServeContent(f, fi.Name(), fi.ModTime())
}

func (c *leegoContext) Attachment(r io.ReadSeeker, name string) (err error) {
	c.response.Header().Set(HeaderContentType, ContentTypeByExtension(name))
	c.response.Header().Set(HeaderContentDisposition, "attachment; filename="+name)
	c.response.WriteHeader(http.StatusOK)
	_, err = io.Copy(c.response, r)
	return
}

func (c *leegoContext) NoContent(code int) error {
	c.response.WriteHeader(code)
	return nil
}

func (c *leegoContext) Redirect(code int, url string) error {
	if code < http.StatusMultiplleegoices || code > http.StatusTemporaryRedirect {
		return ErrInvalidRedirectCode
	}
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(code)
	return nil
}

func (c *leegoContext) Error(err error) {
	c.leego.httpErrorHandler(err, c)
}

func (c *leegoContext) Leego() *Leego {
	return c.leego
}

func (c *leegoContext) Handler() HandlerFunc {
	return c.handler
}

func (c *leegoContext) SetHandler(h HandlerFunc) {
	c.handler = h
}

//func (c *leegoContext) Logger() log.Logger {
//	return c.leego.logger
//}

func (c *leegoContext) ServeContent(content io.ReadSeeker, name string, modtime time.Time) error {
	req := c.Request()
	res := c.Response()

	if t, err := time.Parse(http.TimeFormat, req.Header().Get(HeaderIfModifiedSince)); err == nil && modtime.Before(t.Add(1*time.Second)) {
		res.Header().Del(HeaderContentType)
		res.Header().Del(HeaderContentLength)
		return c.NoContent(http.StatusNotModified)
	}

	res.Header().Set(HeaderContentType, ContentTypeByExtension(name))
	res.Header().Set(HeaderLastModified, modtime.UTC().Format(http.TimeFormat))
	res.WriteHeader(http.StatusOK)
	_, err := io.Copy(res, content)
	return err
}

// ContentTypeByExtension returns the MIME type associated with the file based on
// its extension. It returns `application/octet-stream` incase MIME type is not
// found.
func ContentTypeByExtension(name string) (t string) {
	if t = mime.TypeByExtension(filepath.Ext(name)); t == "" {
		t = MIMEOctetStream
	}
	return
}

func (c *leegoContext) Reset(req engine.Request, res engine.Response) {
	c.context = context.Background()
	c.request = req
	c.response = res
	c.handler = NotFoundHandler
	c.data = make(map[string]interface{})
}
