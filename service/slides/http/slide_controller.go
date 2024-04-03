package http

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"slide-share/model"
	"slide-share/service/slides/usecase"
	"slide-share/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ISlideController interface {
	GetNewestSlideGroup(c echo.Context) error
	GetSlideGroupByPage(c echo.Context) error
	GetSlideGroup(c echo.Context) error
	GetSlideGroups(c echo.Context) error
	CreateSlideGroup(c echo.Context) error
	UpdateSlideGroup(c echo.Context) error
	GetSlide(c echo.Context) error
	UpdateSlide(c echo.Context) error
	UploadSlideBySlidesURL(c echo.Context) error
	UploadSlideByPDF(c echo.Context) error
}

type SlideController struct {
	su usecase.ISlideUsecase
}

func NewSlideController(su usecase.ISlideUsecase) ISlideController {
	return &SlideController{su: su}
}

func (sc *SlideController) GetNewestSlideGroup(c echo.Context) error {
	slideGroup, err := sc.su.GetNewestSlideGroup()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, slideGroup)
}

func (sc *SlideController) GetSlideGroupByPage(c echo.Context) error {
	page := c.QueryParam("page")
	if page == "" {
		return c.JSON(http.StatusBadRequest, "page query parameter is required")
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid page query parameter")
	}

	slideGroups, err := sc.su.GetSlideGroupByPage(pageInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, slideGroups)
}

func (sc *SlideController) GetSlideGroup(c echo.Context) error {
	slideGroupID := c.Param("slide_group_id")
	slideGroup, err := sc.su.GetSlideGroup(slideGroupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, slideGroup)
}

func (sc *SlideController) GetSlideGroups(c echo.Context) error {
	slideGroups, err := sc.su.GetSlideGroups()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, slideGroups)
}

func (sc *SlideController) CreateSlideGroup(c echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")
	secret := os.Getenv("AUTH_SECRET")
	payload, err := utils.VerifyAndGetUserClaims(authToken, secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	if payload.Role == "user" || payload.Role == "" {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	slideGroup := model.SlideGroup{}
	if err := c.Bind(&slideGroup); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	slideGroupID := c.Param("slide_group_id")
	slideGroup.ID = slideGroupID

	DriveID, err := sc.su.CreateSlideGroup(&slideGroup)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, DriveID)
}

func (sc *SlideController) UpdateSlideGroup(c echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")
	secret := os.Getenv("AUTH_SECRET")
	payload, err := utils.VerifyAndGetUserClaims(authToken, secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	if payload.Role == "user" || payload.Role == "" {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	slideGroup := model.SlideGroup{}
	if err := c.Bind(&slideGroup); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	slideGroupID := c.Param("slide_group_id")

	err = sc.su.UpdateSlideGroup(slideGroupID, &slideGroup)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (sc *SlideController) GetSlide(c echo.Context) error {
	slideGroupID := c.Param("slide_group_id")
	slideID := c.Param("slide_id")
	slide, err := sc.su.GetSlide(slideGroupID, slideID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, slide)
}

func (sc *SlideController) UpdateSlide(c echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")
	secret := os.Getenv("AUTH_SECRET")
	payload, err := utils.VerifyAndGetUserClaims(authToken, secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	if payload.Role == "user" || payload.Role == "" {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	slide := model.Slide{}
	if err := c.Bind(&slide); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	slideGroupID := c.Param("slide_group_id")
	slideID := c.Param("slide_id")

	err = sc.su.UpdateSlide(slideGroupID, slideID, &slide)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (sc *SlideController) UploadSlideBySlidesURL(c echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")
	secret := os.Getenv("AUTH_SECRET")
	payload, err := utils.VerifyAndGetUserClaims(authToken, secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	if payload.Role == "user" || payload.Role == "" {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	slideUploadBySlidesURL := model.SlideUploadBySlidesURL{}
	if err = c.Bind(&slideUploadBySlidesURL); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = sc.su.UploadSlideBySlidesURL(&slideUploadBySlidesURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, nil)
}

func (sc *SlideController) UploadSlideByPDF(c echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")
	secret := os.Getenv("AUTH_SECRET")
	payload, err := utils.VerifyAndGetUserClaims(authToken, secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	if payload.Role == "user" || payload.Role == "" {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	// リクエストボディのデータを取得
	data := c.FormValue("data")
	slideUploadByPDF := model.SlideUploadByPDF{}
	if err := json.Unmarshal([]byte(data), &slideUploadByPDF); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error()) // JSONの解析エラー
	}

	// ファイルの取得
	file, err := c.FormFile("pdf")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "PDF file is required")
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	pdfBytes, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// サムネイル画像の取得
	thumbnailFile, err := c.FormFile("thumbnail")

	// サムネイル画像がない場合はnilで渡す
	if err != nil {
		if err != http.ErrMissingFile {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		err = sc.su.UploadSlideByPDF(pdfBytes, nil, &slideUploadByPDF)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusCreated, nil)
	}

	thumbnailSrc, err := thumbnailFile.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer thumbnailSrc.Close()

	// サムネイル画像をバイト配列に読み込む（必要に応じて）
	thumbnailBytes, err := io.ReadAll(thumbnailSrc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = sc.su.UploadSlideByPDF(pdfBytes, thumbnailBytes, &slideUploadByPDF)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, nil)
}
