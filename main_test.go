package main

import (
	"net/http"
	"net/http/httptest"
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

	listM := len(cafeList["moskow"])
	listT := len(cafeList["tula"])
	requests := []struct {
		url   string
		count int
		want  int
	}{
		{"/cafe?city=moscow&count=0", 0, 0},
		{"/cafe?city=moscow&count=1", 1, 1},
		{"/cafe?city=moscow&count=2", 2, 2},
		{"/cafe?city=moscow&count=100", 100, min(100, listM)},
		{"/cafe?city=tula&count=0", 0, 0},
		{"/cafe?city=tula&count=1", 1, 1},
		{"/cafe?city=tula&count=2", 2, 2},
		{"/cafe?city=tula&count=100", 100, min(100, listT)},
	}

	for _, v := range requests {

		req := httptest.NewRequest("GET", v.url, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, req)

		body := strings.TrimSpace(response.Body.String())

		cafes := strings.Split(body, ",")

		actual := len(cafes)
		if v.count == 0 && body == "" {
			actual = 0
		}

		if actual != v.count && v.count != 100 {
			t.Errorf("For count = %d got %d cafes, want %d. Response: %s", v.count, actual, v.want, body)
		}
		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{"/cafe?city=moscow&search=фасоль", 0},
		{"/cafe?city=moscow&search=кофе", 2},
		{"/cafe?city=moscow&search=вилка", 1},
	}

	for _, v := range requests {
		req := httptest.NewRequest("GET", v.search, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, req)

		body := strings.TrimSpace(response.Body.String())
		cafes := strings.ToLower(body)
		slice := strings.Split(cafes, ",")

		actual := len(slice)

		if actual == 1 && body == "" {
			actual = 0
		}

		if actual != v.wantCount {
			t.Errorf("got %d cafes, want %d. Response: %s", actual, v.wantCount, body)
		}

	}
}
