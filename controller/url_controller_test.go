package controller

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/vvikash157/url_shortener/mocks"
	"github.com/vvikash157/url_shortener/models"
)

func TestShortUrlHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//create a mock service
	mockService := mocks.NewMockURLService(ctrl)
	controler := NewUrlController(mockService)

	//prepare input and expectations
	input := `{"long_url":"www.abc.com/xyz/pqrs"}`
	expectedResponse := models.ShortenResponse{ShortUrl: "www.abc.com/2"}

	mockService.EXPECT().UrlShortener(models.ShortenRequest{LongUrl: "www.abc.com/xyz/pqrs"}).Return(expectedResponse, nil)

	//prepares request/response recoder
	req := httptest.NewRequest("POST", "/shortner", bytes.NewBufferString(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	//call the handler
	controler.ShortUrlHandler(w, req)

	//assert
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", w.Code)
	}

	if body := w.Body.String(); !bytes.Contains([]byte(body), []byte("www.abc.com/2")) {
		t.Errorf("body doesnot contain shortened code; got %s", body)
	}

}

func TestShortUrlHandler_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//create a mock service
	mockService := mocks.NewMockURLService(ctrl)
	controler := NewUrlController(mockService)

	//prepare input and expectations
	// input:=models.ShortenRequest{LongUrl: ""}
	// expectations:=models.ShortenResponse{ShortUrl: }

	//prepare request response recoder
	req := httptest.NewRequest("POST", "/shortner", nil)
	w := httptest.NewRecorder()

	//call handler
	controler.ShortUrlHandler(w, req)

	//assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid input; we got %d", w.Code)
	}

}

func TestShortUrlHandler_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockURLService(ctrl)
	controler := NewUrlController(mockService)

	input := `{"long_url":"www.abc.com/vvviikkkaasshh"}`

	mockService.EXPECT().UrlShortener(models.ShortenRequest{LongUrl:"www.abc.com/vvviikkkaasshh"}).Return(models.ShortenResponse{}, errors.New("simulating service error"))

	req := httptest.NewRequest("POST", "/shortner", bytes.NewBufferString(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	controler.ShortUrlHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code 500 because error from service; got %d", w.Code)
	}

}
