package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gta2024/pkg/models"
	"gta2024/tests/utils"
)

func getUrl(bannerId int64) string {
	return fmt.Sprintf("/banner/%d", bannerId)
}

func TestSuccessDelete(t *testing.T) {
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
		tagIDs     = []int64{1, 2, 3}
		contentStr = "{\"title\": \"example\"}"
		isAcive    = true

		token = utils.GenerateToken(t, "admin")
	)

	// init db data
	utils.CreateOrUpdateBanner(db, models.DBBanner{
		BannerID: bannerId,
		Content:  contentStr,
		IsActive: isAcive,
	})
	for _, tagId := range tagIDs {
		utils.CreateFeatureTagBanner(db, models.DBFeatureTagBanner{
			BannerID:  bannerId,
			FeatureID: featureId,
			TagID:     tagId,
		})
	}

	// create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", getUrl(bannerId), bytes.NewBuffer([]byte{}))
	utils.SetAuthToken(req, token)
	router.ServeHTTP(w, req)

	// check status
	if w.Code != http.StatusNoContent {
		t.Fatalf("Status is not NoContent: %d", w.Code)
	}

	// check db data
	dbBanners, _ := utils.GetBanners(db)
	if len(dbBanners) != 0 {
		t.Fatalf("banner table must be empty, current - %+v", dbBanners)
	}
	dbFTBs, _ := utils.GetFeatureTagBanners(db)
	if len(dbFTBs) != 0 {
		t.Fatalf("feature_tag_banner table must empty, current - %+v", dbFTBs)
	}
}

func TestNotFoundDelete(t *testing.T) {
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
		bannerId = int64(10)

		token = utils.GenerateToken(t, "admin")
	)

	// create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", getUrl(bannerId), bytes.NewBuffer([]byte{}))
	utils.SetAuthToken(req, token)
	router.ServeHTTP(w, req)

	// check status
	if w.Code != http.StatusNotFound {
		t.Fatalf("Status is not NotFound: %d", w.Code)
	}
}

func TestUnauthorizedDeleteBanner(t *testing.T) {
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
		bannerId = int64(10)
	)

	// create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", getUrl(bannerId), bytes.NewBuffer([]byte{}))
	router.ServeHTTP(w, req)

	// check status
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Status is not Unauthorized: %d", w.Code)
	}
}

func TestForbiddenDeleteBanner(t *testing.T) {
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
		bannerId = int64(10)

		token = utils.GenerateToken(t, "user")
	)

	// create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", getUrl(bannerId), bytes.NewBuffer([]byte{}))
	utils.SetAuthToken(req, token)
	router.ServeHTTP(w, req)

	// check status
	if w.Code != http.StatusForbidden {
		t.Fatalf("Status is not Forbidden: %d", w.Code)
	}
}
