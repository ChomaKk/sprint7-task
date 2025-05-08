package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, min(100, len(cafeList["tula"]))},
	}

	// Tula
	for _, v := range requests {
		response := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/cafe?city=tula&count="+strconv.Itoa(v.count), nil)
		handler.ServeHTTP(response, req)
		if v.count == 0 {
			assert.Equal(t, "", response.Body.String())
			continue
		}
		resp := strings.Split(response.Body.String(), ",")
		assert.Equal(t, v.want, len(resp))
	}

	// Moscow
	rq := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, min(100, len(cafeList["moscow"]))},
	}
	for _, v := range rq {
		response := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/cafe?city=moscow&count="+strconv.Itoa(v.count), nil)
		handler.ServeHTTP(response, req)
		if v.count == 0 {
			assert.Equal(t, "", response.Body.String())
			continue
		}
		resp := strings.Split(response.Body.String(), ",")
		assert.Equal(t, v.want, len(resp))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		search string
		want   int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	for _, v := range requests {
		count := 0
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&search="+v.search, nil)
		handler.ServeHTTP(response, req)

		eachCafe := strings.Split(response.Body.String(), ",")
		for _, e := range eachCafe {
			if strings.Contains(strings.ToLower(e), v.search) {
				count++
				continue
			}
		}
		assert.Equal(t, v.want, count)
	}
}
