package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/vvikash157/url_shortener/models"
	"github.com/vvikash157/url_shortener/repository"
	"github.com/vvikash157/url_shortener/utils"

	"github.com/pkg/errors"
)

type URLService interface {
	UrlShortener(request models.ShortenRequest) (models.ShortenResponse, error)
	ResolveUrl(code string) (string, error)
}

type urlService struct {
	urlRepo    repository.UrlRepository
	cacheRepo  repository.CacheRepository
	counter    int
	counterMux sync.Mutex
	baseURL    string
}

func NewURLService(urlRepo repository.UrlRepository, cacheRepo repository.CacheRepository, baseUrl string) URLService {
	return &urlService{
		urlRepo:   urlRepo,
		cacheRepo: cacheRepo,
		baseURL:   baseUrl,
		counter:   1,
	}
}

func (u *urlService) UrlShortener(req models.ShortenRequest) (models.ShortenResponse, error) {
	if req.LongUrl == "" {
		fmt.Println("longUrl is empty")
		return models.ShortenResponse{}, errors.New("long url not available")
	}

	//check cache for existing code
	if code, err := u.cacheRepo.Get(req.LongUrl); err == nil {
		return models.ShortenResponse{ShortUrl: fmt.Sprintf("%s/%s", u.baseURL, code)}, err
	}

	if code, err := u.urlRepo.GetCodeByLongUrl(req.LongUrl); err == nil {
		u.cacheRepo.Set(req.LongUrl, code, 7*24*time.Hour)
		return models.ShortenResponse{ShortUrl: fmt.Sprintf("%s/%s", u.baseURL, code)}, err
	}

	u.counterMux.Lock()
	defer u.counterMux.Unlock()
	u.counter++

	code := utils.EncodeBase62(u.counter)

	if err := u.urlRepo.InsertUrl(req.LongUrl, code); err != nil {
		return models.ShortenResponse{}, err
	}

	u.cacheRepo.Set(req.LongUrl, code, 7*24*time.Hour)
	u.cacheRepo.Set(code, req.LongUrl, 7*24*time.Hour)

	fmt.Printf("short url created %s", code)
	return models.ShortenResponse{ShortUrl: fmt.Sprintf("%s/%s", u.baseURL, code)}, nil

}

func (u *urlService) ResolveUrl(code string) (string, error) {
	if longUrl, err := u.cacheRepo.Get(code); err != nil {
		return longUrl, err
	}
	return u.urlRepo.GetLongUrlByCode(code)
}
