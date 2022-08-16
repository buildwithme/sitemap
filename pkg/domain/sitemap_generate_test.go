package domain

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/buildwithme/sitemap/pkg"
	"github.com/buildwithme/sitemap/pkg/models"
)

func TestEmpty(t *testing.T) {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	options := &pkg.ParameterOptions{
		Parallel:   1,
		MaxDepth:   1,
		Timeout:    3,
		RetryDetay: 1,
		MaxRetry:   3,
	}

	var expectedURLDetails []models.UrlResult
	var expectedFailedURLs []models.FailedURL

	server := httptest.NewServer(http.HandlerFunc(handleFunc))
	defer server.Close()

	urlDetails, failedURLs := StartSitemapGeneration(server.URL, options)

	if !reflect.DeepEqual(urlDetails, expectedURLDetails) {
		t.Errorf("expectedURLs doesn't match, got %v", urlDetails)
	}

	if !reflect.DeepEqual(failedURLs, expectedFailedURLs) {
		t.Errorf("failedURLs doesn't match, got %v", failedURLs)
	}
}

func TestWithDepth1(t *testing.T) {

	fmt.Println("Test with depth 1 - should not mark the URLs on level 1 as failed as no request was sent to them")
	options := &pkg.ParameterOptions{
		Parallel:   1,
		MaxDepth:   1,
		Timeout:    3,
		RetryDetay: 1,
		MaxRetry:   3,
	}
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
			<html lang="en">
			<head>
			</head>
			<body>
				<a href="/">Test base </a>
				<a href="/view">Test view</a>
				<a href="/root">Test root</a>
			</body>
			</html>
			`))
			return
		} else if r.URL.Path == "/view" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
	var expectedURLs = []models.UrlResult{
		{
			URL:   "%SERVER%/",
			Title: "Test base ",
		},
		{
			URL:   "%SERVER%/view",
			Title: "Test view",
		},
		{
			URL:   "%SERVER%/root",
			Title: "Test root",
		},
	}
	var expectedFailedURLs []models.FailedURL

	server := httptest.NewServer(http.HandlerFunc(handleFunc))
	defer server.Close()

	for i := range expectedURLs {
		expectedURLs[i].URL = strings.Replace(expectedURLs[i].URL, "%SERVER%", server.URL, 1)
	}

	for i := range expectedFailedURLs {
		expectedFailedURLs[i].URL = strings.Replace(expectedFailedURLs[i].URL, "%SERVER%", server.URL, 1)
	}

	urlDetails, failedURLs := StartSitemapGeneration(server.URL, options)

	if !reflect.DeepEqual(urlDetails, expectedURLs) {
		t.Errorf("expectedURLs %v doesn't match, got %v", expectedURLs, urlDetails)
	}

	if !reflect.DeepEqual(failedURLs, expectedFailedURLs) {
		t.Errorf("failedURLs %v doesn't match, got %v", expectedFailedURLs, failedURLs)
	}
}

func TestWithDepth2(t *testing.T) {

	fmt.Println("Test with depth 2 - should not mark the URLs on level 1 as failed as no request was sent to them")
	options := &pkg.ParameterOptions{
		Parallel:   1,
		MaxDepth:   2,
		Timeout:    3,
		RetryDetay: 1,
		MaxRetry:   3,
	}
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
			<html lang="en">
			<head>
			</head>
			<body>
				<a href="/">Test base </a>
				<a href="/view">Test view</a>
				<a href="/root">Test root</a>
			</body>
			</html>
			`))
			return
		} else if r.URL.Path == "/view" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
	var expectedURLs = []models.UrlResult{
		{
			URL:   "%SERVER%/",
			Title: "Test base ",
		},
	}
	var expectedFailedURLs []models.FailedURL = []models.FailedURL{
		{
			URL:    "%SERVER%/view",
			Reason: "500 Internal Server Error",
		},
		{
			URL:    "%SERVER%/root",
			Reason: "400 Bad Request",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(handleFunc))
	defer server.Close()

	for i := range expectedURLs {
		expectedURLs[i].URL = strings.Replace(expectedURLs[i].URL, "%SERVER%", server.URL, 1)
	}

	for i := range expectedFailedURLs {
		expectedFailedURLs[i].URL = strings.Replace(expectedFailedURLs[i].URL, "%SERVER%", server.URL, 1)
	}

	urlDetails, failedURLs := StartSitemapGeneration(server.URL, options)

	if !reflect.DeepEqual(urlDetails, expectedURLs) {
		t.Errorf("expectedURLs %v doesn't match, got %v", expectedURLs, urlDetails)
	}

	if !reflect.DeepEqual(failedURLs, expectedFailedURLs) {
		t.Errorf("failedURLs %v doesn't match, got %v", expectedFailedURLs, failedURLs)
	}
}

func TestBadURLs(t *testing.T) {

	fmt.Println("Test with depth 2 - should not mark the URLs on level 1 as failed as no request was sent to them")
	options := &pkg.ParameterOptions{
		Parallel:   1,
		MaxDepth:   3,
		Timeout:    3,
		RetryDetay: 1,
		MaxRetry:   3,
	}
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
			<html lang="en">
			<head>
			</head>
			<body>
				<a href="/">Test base </a>
				<a href="/view">Test view<a>
				<a href="/root">Test root</a>
			</body>
			</html>
			`))
			return
		} else if r.URL.Path == "/view" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
	var expectedURLs = []models.UrlResult{
		{
			URL:   "%SERVER%/",
			Title: "Test base ",
		},
	}
	var expectedFailedURLs []models.FailedURL = []models.FailedURL{
		{
			URL:    "%SERVER%/view",
			Reason: "500 Internal Server Error",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(handleFunc))
	defer server.Close()

	for i := range expectedURLs {
		expectedURLs[i].URL = strings.Replace(expectedURLs[i].URL, "%SERVER%", server.URL, 1)
	}

	for i := range expectedFailedURLs {
		expectedFailedURLs[i].URL = strings.Replace(expectedFailedURLs[i].URL, "%SERVER%", server.URL, 1)
	}

	urlDetails, failedURLs := StartSitemapGeneration(server.URL, options)

	if !reflect.DeepEqual(urlDetails, expectedURLs) {
		t.Errorf("expectedURLs %v doesn't match, got %v", expectedURLs, urlDetails)
	}

	if !reflect.DeepEqual(failedURLs, expectedFailedURLs) {
		t.Errorf("failedURLs %v doesn't match, got %v", expectedFailedURLs, failedURLs)
	}
}
