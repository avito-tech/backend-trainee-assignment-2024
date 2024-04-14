package tests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gta2024/pkg/models"
	"gta2024/tests/utils"
)

func getQueryParamsUrl(featureId int64, tagId int64) string {
	return fmt.Sprintf("/user_banner?feature_id=%d&tag_id=%d", featureId, tagId)
}

func TestSuccessUserGetActiveBanner(t *testing.T) {
	// init db
	db, err := psql.SetUp()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer psql.TearDown(db)

	// init service
	router := utils.InitRouter(db)

	// test data
	var (
		bannerId   = int64(10)
		featureId  = int64(1)
		tagId      = int64(2)
		contentStr = "{\"title\":\"example\"}"
		isAcive    = true

		token = utils.GenerateToken(t, "user")
	)

	// init db data
	utils.CreateOrUpdateBanner(db, models.DBBanner{
		BannerID: bannerId,
		Content:  contentStr,
		IsActive: isAcive,
	})
	utils.CreateFeatureTagBanner(db, models.DBFeatureTagBanner{
		BannerID:  bannerId,
		FeatureID: featureId,
		TagID:     tagId,
	})

	// create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", getQueryParamsUrl(featureId, tagId), bytes.NewBuffer([]byte{}))
	utils.SetAuthToken(req, token)
	router.ServeHTTP(w, req)

	// check status
	if w.Code != http.StatusOK {
		t.Fatalf("Status is not Ok: %d", w.Code)
	}

	// check response body
	responseBody, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Error read response: %s", err.Error())
	}
	if string(responseBody) != contentStr {
		t.Fatalf("Invalid response body: %s", string(responseBody))
	}
}

func TestNotFoundUserGetNotActiveBanner(t *testing.T) {
	// init db
	db, err := psql.SetUp()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer psql.TearDown(db)

	// init service
	router := utils.InitRouter(db)

	// test data
	var (
		bannerId   = int64(10)
		featureId  = int64(1)
		tagId      = int64(2)
		contentStr = "{\"title\":\"example\"}"
		isAcive    = false

		token = utils.GenerateToken(t, "user")
	)

	// init db data
	utils.CreateOrUpdateBanner(db, models.DBBanner{
		BannerID: bannerId,
		Content:  contentStr,
		IsActive: isAcive,
	})
	utils.CreateFeatureTagBanner(db, models.DBFeatureTagBanner{
		BannerID:  bannerId,
		FeatureID: featureId,
		TagID:     tagId,
	})

	// create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", getQueryParamsUrl(featureId, tagId), bytes.NewBuffer([]byte{}))
	utils.SetAuthToken(req, token)
	router.ServeHTTP(w, req)

	// check status
	if w.Code != http.StatusNotFound {
		t.Fatalf("Status is not NotFound: %d", w.Code)
	}
}

func TestSuccessAdminGetActiveBanner(t *testing.T) {
	// init db
	db, err := psql.SetUp()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer psql.TearDown(db)

	// init service
	router := utils.InitRouter(db)

	// test data
	var (
		bannerId   = int64(10)
		featureId  = int64(1)
		tagId      = int64(2)
		contentStr = "{\"title\":\"example\"}"
		isAcive    = true

		token = utils.GenerateToken(t, "admin")
	)

	// init db data
	utils.CreateOrUpdateBanner(db, models.DBBanner{
		BannerID: bannerId,
		Content:  contentStr,
		IsActive: isAcive,
	})
	utils.CreateFeatureTagBanner(db, models.DBFeatureTagBanner{
		BannerID:  bannerId,
		FeatureID: featureId,
		TagID:     tagId,
	})

	// create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", getQueryParamsUrl(featureId, tagId), bytes.NewBuffer([]byte{}))
	utils.SetAuthToken(req, token)
	router.ServeHTTP(w, req)

	// check status
	if w.Code != http.StatusOK {
		t.Fatalf("Status is not Ok: %d", w.Code)
	}

	// check response body
	responseBody, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Error read response: %s", err.Error())
	}
	if string(responseBody) != contentStr {
		t.Fatalf("Invalid response body: %s", string(responseBody))
	}
}

func TestSuccessAdminGetNotActiveBanner(t *testing.T) {
	// init db
	db, err := psql.SetUp()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer psql.TearDown(db)

	// init service
	router := utils.InitRouter(db)

	// test data
	var (
		bannerId   = int64(10)
		featureId  = int64(1)
		tagId      = int64(2)
		contentStr = "{\"title\":\"example\"}"
		isAcive    = false

		token = utils.GenerateToken(t, "admin")
	)

	// init db data
	utils.CreateOrUpdateBanner(db, models.DBBanner{
		BannerID: bannerId,
		Content:  contentStr,
		IsActive: isAcive,
	})
	utils.CreateFeatureTagBanner(db, models.DBFeatureTagBanner{
		BannerID:  bannerId,
		FeatureID: featureId,
		TagID:     tagId,
	})

	// create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", getQueryParamsUrl(featureId, tagId), bytes.NewBuffer([]byte{}))
	utils.SetAuthToken(req, token)
	router.ServeHTTP(w, req)

	// check status
	if w.Code != http.StatusOK {
		t.Fatalf("Status is not Ok: %d", w.Code)
	}

	// check response body
	responseBody, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Error read response: %s", err.Error())
	}
	if string(responseBody) != contentStr {
		t.Fatalf("Invalid response body: %s", string(responseBody))
	}
}

func TestNotFoundGetBanner(t *testing.T) {
	// init db
	db, err := psql.SetUp()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer psql.TearDown(db)

	// init service
	router := utils.InitRouter(db)

	// test data
	var (
		featureId = int64(1)
		tagId     = int64(2)
	)

	for _, role := range []string{"user", "admin"} {
		token := utils.GenerateToken(t, role)

		// create request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", getQueryParamsUrl(featureId, tagId), bytes.NewBuffer([]byte{}))
		utils.SetAuthToken(req, token)
		router.ServeHTTP(w, req)

		// check status
		if w.Code != http.StatusNotFound {
			t.Fatalf("Status is not NotFound: %d", w.Code)
		}
	}
}

func TestUnauthorizedGetBanner(t *testing.T) {
	// init db
	db, err := psql.SetUp()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer psql.TearDown(db)

	// init service
	router := utils.InitRouter(db)

	// test data
	var (
		featureId = int64(1)
		tagId     = int64(2)
	)

	// create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", getQueryParamsUrl(featureId, tagId), bytes.NewBuffer([]byte{}))
	router.ServeHTTP(w, req)

	// check status
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Status is not Unauthorized: %d", w.Code)
	}
}
