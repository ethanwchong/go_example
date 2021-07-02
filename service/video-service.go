package service

import (
	"videoproject/entity"
	"videoproject/repository"
)

type VideoService interface {
	Save(entity.Video) error
	Update(entity.Video) error
	Delete(entity.Video) error
	FindAll() []entity.Video
}

type videoService struct {
	//videos []entity.Video
	videoRepository repository.VideoRepository
}

func NewVideoService(repo repository.VideoRepository) VideoService {
	return &videoService{
		videoRepository: repo,
	}
}

func (videoService *videoService) Save(video entity.Video) error {
	//videoService.videos = append(videoService.videos, video)
	videoService.videoRepository.Save(video)
	return nil
}

func (videoService *videoService) Update(video entity.Video) error {
	videoService.videoRepository.Update(video)
	return nil
}

func (videoService *videoService) Delete(video entity.Video) error {
	videoService.videoRepository.Delete(video)
	return nil
}

func (videoService *videoService) FindAll() []entity.Video {
	return videoService.videoRepository.FindAll()
}
