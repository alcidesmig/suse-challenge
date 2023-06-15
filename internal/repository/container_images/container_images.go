package container_images

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"suse-cli-challenge/internal/repository"
)

// ImplContainerImagesRepository implements repository.ContainerImagesRepository.
type ImplContainerImagesRepository struct {
	repository.ContainerImagesRepository
}

// NewContainerImagesRepository creates a new instance of ContainerImagesRepository.
func NewContainerImagesRepository() repository.ContainerImagesRepository {
	return &ImplContainerImagesRepository{}
}

// GetReferencedContainerImagesFromChartDir retrieves referenced container images from a chart directory.
// It lists all the container images that will be used in the installation.
// This method has the following constraints:
// - On Linux PCs, it requires the presence of "bash", "helm", "grep", "sort", and "sed" installations.
// - On Windows PCs, it requires the presence of "cmd", "helm", "findstr", "sort", and "sed" installations.
func (icir *ImplContainerImagesRepository) GetReferencedContainerImagesFromChartDir(path string) ([]string, error) {
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
