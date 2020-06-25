package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-wyvern/leego"
	"github.com/go-wyvern/leego/engine/standard"
	"github.com/stretchr/testify/assert"
)

func TestAddTrailingSlash(t *testing.T) {
	lee := leego.New()

	req := standard.NewRequest(httptest.NewRequest(leego.GET, "/add-slash", nil))
	rec := standard.NewResponse(httptest.NewRecorder())
	c := lee.NewContext(req, rec)
	h := AddTrailingSlash()(func(c leego.Context) leego.LeeError {
		return nil
	})
	h(c)
	assert.Equal(t, "/add-slash/", req.URL().Path())
	assert.Equal(t, "/add-slash/", req.URI())

	// With config
	req = standard.NewRequest(httptest.NewRequest(leego.GET, "/add-slash?key=value", nil))
	rec = standard.NewResponse(httptest.NewRecorder())
	c = lee.NewContext(req, rec)
	h = AddTrailingSlashWithConfig(TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	})(func(c leego.Context) leego.LeeError {
		return nil
	})
	h(c)
	assert.Equal(t, http.StatusMovedPermanently, rec.Status())
	assert.Equal(t, "/add-slash/?key=value", rec.Header().Get(leego.HeaderLocation))
}

func TestRemoveTrailingSlash(t *testing.T) {
	lee := leego.New()
	req := standard.NewRequest(httptest.NewRequest(leego.GET, "/remove-slash/", nil))
	rec := standard.NewResponse(httptest.NewRecorder())
	c := lee.NewContext(req, rec)
	h := RemoveTrailingSlash()(func(c leego.Context) leego.LeeError {
		return nil
	})
	h(c)
	assert.Equal(t, "/remove-slash", req.URL().Path())
	assert.Equal(t, "/remove-slash", req.URI())

	// With config
	req = standard.NewRequest(httptest.NewRequest(leego.GET, "/remove-slash/?key=value", nil))
	rec = standard.NewResponse(httptest.NewRecorder())
	c = lee.NewContext(req, rec)
	h = RemoveTrailingSlashWithConfig(TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	})(func(c leego.Context) leego.LeeError {
		return nil
	})
	h(c)
	assert.Equal(t, http.StatusMovedPermanently, rec.Status())
	assert.Equal(t, "/remove-slash?key=value", rec.Header().Get(leego.HeaderLocation))

	// With bare URL
	req = standard.NewRequest(httptest.NewRequest(leego.GET, "http://localhost", nil))
	rec = standard.NewResponse(httptest.NewRecorder())
	c = lee.NewContext(req, rec)
	h = RemoveTrailingSlash()(func(c leego.Context) leego.LeeError {
		return nil
	})
	h(c)
	assert.Equal(t, "", req.URL().Path())
	assert.Equal(t, "http://localhost", req.URI())
}
