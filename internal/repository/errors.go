package repository

import "errors"

var ErrReadingCharts = errors.New("error reading charts")
var ErrVersionAlreadyExists = errors.New("error version already exists")
var ErrSavingChart = errors.New("error saving chart")
var ErrConfiguringCLI = errors.New("error configuring cli")
var ErrWritingFile = errors.New("error while writing to file")
var ErrReadingFile = errors.New("error while reading file")
var ErrPackagingChart = errors.New("error packaging chart")
var ErrInvalidChartDir = errors.New("error invalid chart directory")
var ErrLoadingChart = errors.New("error loading chart")
var ErrGettingHomeDir = errors.New("error getting home directory for current user")
var ErrInvalidRepoURL = errors.New("error invalid repository url")
var ErrCloningRepo = errors.New("error cloning the repository")
var ErrInstallingChart = errors.New("error installing chart")
var ErrReadingValues = errors.New("error reading values.yaml")
var ErrReadingContainerImages = errors.New("error reading container images")
