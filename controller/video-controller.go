package controller

import (
	"net/http"
	"strconv"
	"videoproject/entity"
	"videoproject/service"
	"videoproject/validators"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type VideoController interface {
	FindAll() []entity.Video
	Save(ctx *gin.Context) error
	Update(ctx *gin.Context) error
	Delete(ctx *gin.Context) error
	ShowAll(ctx *gin.Context)
}

type videoController struct {
	videoService service.VideoService
}

var validate *validator.Validate

func New(videoService service.VideoService) VideoController {
	validate = validator.New()
	validate.RegisterValidation("is-cool", validators.ValidateCoolTitle)
	return &videoController{
		videoService: videoService,
	}
}

func (c *videoController) FindAll() []entity.Video {
	return c.videoService.FindAll()
}

func (c *videoController) Save(ctx *gin.Context) error {
	var video entity.Video
	err := ctx.ShouldBindJSON(&video)
	if err != nil {
		return err
	}
	err = validate.Struct(video)
	if err != nil {
		return err
	}
	c.videoService.Save(video)
	return nil
}

func (c *videoController) Update(ctx *gin.Context) error {
	var video entity.Video
	err := ctx.ShouldBindJSON(&video)
	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		return err
	}
	video.ID = id

	err = validate.Struct(video)
	if err != nil {
		return err
	}
	c.videoService.Update(video)
	return nil

}
func (c *videoController) Delete(ctx *gin.Context) error {
	var video entity.Video
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		return err
	}
	video.ID = id
	c.videoService.Delete(video)
	return nil
}

func (c *videoController) ShowAll(ctx *gin.Context) {
	videos := c.videoService.FindAll()
	data := gin.H{
		"title":  "Video Page",
		"videos": videos,
	}
	ctx.HTML(http.StatusOK, "index.html", data)
}
