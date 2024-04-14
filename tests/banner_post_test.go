package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/sirupsen/logrus"

	"gta2024/pkg/models"
	"gta2024/tests/utils"
)

func buildRequestBody(featureId int64, tagIds []int64, content models.BannerContent, isActive bool) models.CreateBanner {
	return models.CreateBanner{
		FeatureID: &featureId,
		TagIDs:    &tagIds,
		Content:   &content,
		IsActive:  &isActive,
	}
}

func TestSuccessCreate(t *testing.T) {
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
		featureId  = int64(1)
		tagIDs     = []int64{1, 2, 3}
		content    = models.BannerContent{"title": "example"}
		contentStr = "{\"title\": \"example\"}"
		isAcive    = true

		token = utils.GenerateToken(t, "admin")
	)

	// create request
	requestBody := buildRequestBody(featureId, tagIDs, content, isAcive)
	requestBodyBytes, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/banner", bytes.NewBuffer(requestBodyBytes))
	utils.SetAuthToken(req, token)
	router.ServeHTTP(w, req)

	// check status
	if w.Code != http.StatusCreated {
		t.Fatalf("Status is not Ok: %d", w.Code)
	}

	// check response body
	responseBody, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Error read response: %s", err.Error())
	}

	var bannerPostResp models.BannerPostResp201
	err = json.Unmarshal(responseBody, &bannerPostResp)
	if err != nil {
		t.Fatalf("Error parsing response: %s", err.Error())
	}

	if bannerPostResp.BannerID <= 0 {
		t.Fatalf("Invalid bannerId in response: %d", bannerPostResp.BannerID)
	}

	// check db data
	dbBanners, _ := utils.GetBanners(db)
	if len(dbBanners) != 1 {
		t.Fatalf("No single banner in db, expected - 1, current - %+v", dbBanners)
	}
	if dbBanners[0].BannerID != bannerPostResp.BannerID ||
		dbBanners[0].Content != contentStr ||
		dbBanners[0].IsActive != isAcive {
		t.Fatalf("Incorrect banner data in db: %+v", dbBanners)
	}

	dbFTBs, _ := utils.GetFeatureTagBanners(db)
	if len(dbFTBs) != len(tagIDs) {
		t.Fatalf("Incorrect feature tag banner rows in db: %+v", dbFTBs)

	}
	for _, dbFTB := range dbFTBs {
		if !slices.Contains(tagIDs, dbFTB.TagID) ||
			dbFTB.FeatureID != featureId ||
			dbFTB.BannerID != bannerPostResp.BannerID {
			logrus.Info()
			t.Fatalf("Incorrect data in feature_tag_banner: %+v", dbFTBs)
		}
	}
}

func TestBadRequestCreate(t *testing.T) {
	// init db
	db, err := psql.SetUp()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer psql.TearDown(db)

	// init service
	router := utils.InitRouter(db)

	// test data
	token := utils.GenerateToken(t, "admin")
	cases := []map[string]interface{}{
		{
			"tag_ids":   []int64{1, 2, 3},
			"content":   models.BannerContent{"title": "example"},
			"is_active": false,
		},
		{
			"feature_id": 1,
			"content":    models.BannerContent{"title": "example"},
			"is_active":  false,
		},
		{
			"feature_id": 1,
			"tag_ids":    []int64{1, 2, 3},
			"is_active":  false,
		},
		{
			"feature_id": 1,
			"tag_ids":    []int64{1, 2, 3},
			"content":    models.BannerContent{"title": "example"},
		},
		{
			"feature_id": 1,
			"tag_ids":    []int64{},
			"content":    models.BannerContent{"title": "example"},
			"is_active":  false,
		},
	}

	for _, testCase := range cases {
		// create request
		requestBodyBytes, _ := json.Marshal(testCase)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/banner", bytes.NewBuffer(requestBodyBytes))
		utils.SetAuthToken(req, token)
		router.ServeHTTP(w, req)

		// check status
		if w.Code != http.StatusBadRequest {
			t.Fatalf("Status is not BadRequest: %d", w.Code)
		}
	}

	// check db data
	dbBanners, _ := utils.GetBanners(db)
	if len(dbBanners) != 0 {
		t.Fatalf("banner table must be empty, current - %+v", dbBanners)
	}
	dbFTBs, _ := utils.GetFeatureTagBanners(db)
	if len(dbFTBs) != 0 {
		t.Fatalf("feature_tag_banner table must be empty, current - %+v", dbFTBs)
	}
}

func TestFeatureTagAlreadyExistsCreate(t *testing.T) {
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
		dbBannerId = int64(10)
		featureId  = int64(1)
		dbTagIDs   = []int64{1}
		newTagIDs  = []int64{1, 2, 3}
		content    = models.BannerContent{"title": "example"}
		contentStr = "{\"title\": \"example\"}"
		isAcive    = true

		token = utils.GenerateToken(t, "admin")
	)

	// init db data
	utils.CreateOrUpdateBanner(db, models.DBBanner{
		BannerID: dbBannerId,
		Content:  contentStr,
		IsActive: isAcive,
	})
	for _, tagId := range dbTagIDs {
		utils.CreateFeatureTagBanner(db, models.DBFeatureTagBanner{
			BannerID:  dbBannerId,
			FeatureID: featureId,
			TagID:     tagId,
		})
	}

	// create request
	requestBody := buildRequestBody(featureId, newTagIDs, content, isAcive)
	requestBodyBytes, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/banner", bytes.NewBuffer(requestBodyBytes))
	utils.SetAuthToken(req, token)
	router.ServeHTTP(w, req)

	// check status
	if w.Code != http.StatusConflict {
		t.Fatalf("Status is not Conflict: %d", w.Code)
	}

	// check db data
	dbBanners, _ := utils.GetBanners(db)
	if len(dbBanners) != 1 {
		t.Fatalf("banner table must not have new rows, current - %+v", dbBanners)
	}
	dbFTBs, _ := utils.GetFeatureTagBanners(db)
	if len(dbFTBs) != len(dbTagIDs) {
		t.Fatalf("feature_tag_banner table must not have new rows, current - %+v", dbFTBs)
	}
}

func TestUnauthorizedCreate(t *testing.T) {
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
		tagIDs    = []int64{1, 2, 3}
		content   = models.BannerContent{"title": "example"}
		isAcive   = true
	)

	// create request
	requestBody := buildRequestBody(featureId, tagIDs, content, isAcive)
	requestBodyBytes, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/banner", bytes.NewBuffer(requestBodyBytes))
	router.ServeHTTP(w, req)

	// check status
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Status is not Unauthorized: %d", w.Code)
	}
}

func TestForbiddenCreate(t *testing.T) {
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
		tagIDs    = []int64{1, 2, 3}
		content   = models.BannerContent{"title": "example"}
		isAcive   = true

		token = utils.GenerateToken(t, "user")
	)

	// create request
	requestBody := buildRequestBody(featureId, tagIDs, content, isAcive)
	requestBodyBytes, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/banner", bytes.NewBuffer(requestBodyBytes))
	utils.SetAuthToken(req, token)
	router.ServeHTTP(w, req)

	// check status
	if w.Code != http.StatusForbidden {
		t.Fatalf("Status is not Forbidden: %d", w.Code)
	}
}
