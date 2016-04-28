package docker

type DockerRegistry interface {
	Close()
	Ping() error
	ImageExists(repoName string, tag string) (bool, error)
	GetImageInfo(repoName string, tag string) (digest string, 
		layerAr []map[string]interface{}, err error)
	GetImage(repoName string, tag string, filepath string) error
	DeleteImage(repoName, tag string) error
	PushImage(imageName, imageFilePath, digestString string) error
}
