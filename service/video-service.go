package service

import "videoproject/entity"

type VideoService interface {
	Save(entity.Video) entity.Video
	FindAll() []entity.Video
}

type videoService struct {
	videos []entity.Video
}

func New() VideoService {
	return &videoService{}
}

func (videoService *videoService) Save(video entity.Video) entity.Video {
	videoService.videos = append(videoService.videos, video)
	return video
}

func (videoService *videoService) FindAll() []entity.Video {
	return videoService.videos
}
