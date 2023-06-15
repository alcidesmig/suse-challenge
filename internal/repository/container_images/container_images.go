package container_images

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"suse-cli-challenge/internal/repository"
)

type ContainerImagesRepository struct {
	repository.ContainerImagesRepository
}

func NewContainerImagesRepository() repository.ContainerImagesRepository {
	return &ContainerImagesRepository{}
}

func (ir *ContainerImagesRepository) GetReferencedContainerImagesFromChartDir(path string) ([]string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", fmt.Sprintf(`helm template %s | findstr image: | sed -e "s/[ ]*image:[ ]*//" -e "s/\"//g" | sort /unique`, path))
	} else {
		cmd = exec.Command("bash", "-c", fmt.Sprintf(`helm template %s | grep image: | sed -e 's/[ ]*image:[ ]*//' -e 's/"//g' | sort -u`, path))
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%w: error simulating `helm template` for analyzing the chart container images, details: %s", repository.ErrReadingContainerImages, err.Error())
	}

	imageList := strings.Split(string(output), "\n")
	return imageList, nil
}
