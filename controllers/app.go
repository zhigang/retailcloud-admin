package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/zhigang/retailcloud-admin/factory"
)

// AppController is retail cloud app struct
type AppController struct {
}

// GetAppList is return retail cloud applications list
func (a *AppController) GetAppList(c echo.Context) error {
	request := retailcloud.CreateListAppRequest()
	request.Scheme = "https"

	pageNumber, err := strconv.Atoi(c.QueryParam("pageNumber"))
	if err == nil {
		request.PageNumber = requests.NewInteger(pageNumber)
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err == nil {
		request.PageSize = requests.NewInteger(pageSize)
	}

	response, err := factory.GetRetailCloudClient().ListApp(request)
	if err != nil {
		log.Errorf("Get app list failed, error: %+v", err)
		return err
	}

	if response.IsSuccess() {
		log.Infof("Get app list succeeded, requestID: %s", response.RequestId)
	} else {
		log.Warnf("Get app list failed, response: %+v", response)
	}

	return c.JSON(http.StatusOK, response)
}

// GetEnvList is return retail cloud application's environment list
func (a *AppController) GetEnvList(c echo.Context) error {

	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "param 'id' is required")
	}

	pageNumber, err := strconv.Atoi(c.QueryParam("pageNumber"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "query param 'pageNumber' is required")
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "query param 'pageSize' is required")
	}

	envType := -1
	if c.QueryParam("envType") != "" {
		envType, err = strconv.Atoi(c.QueryParam("envType"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, "query param 'envType' type is integer")
		}
	}

	envName := c.QueryParam("envName")
	envList, err := a.getEnvList(appID, pageNumber, pageSize, envType, envName)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, envList)
}

// DeployApp is deploy application to retail cloud
func (a *AppController) DeployApp(c echo.Context) error {

	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "param 'id' is required")
	}

	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, "query param 'name' is required")
	}

	envName := c.QueryParam("envName")

	envType, err := strconv.Atoi(c.QueryParam("envType"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "query param 'envType' is required")
	}

	image := c.QueryParam("image")
	if image == "" {
		return c.JSON(http.StatusBadRequest, "query param 'image' is required")
	}

	envList, err := a.getEnvList(appID, 1, 100, envType, envName)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	result := []string{}
	if envList != nil && len(envList) > 0 {
		for _, v := range envList {
			success, err := a.deployApp(int(v.EnvId), 1, name, image)
			deployResult := fmt.Sprintf("appID: %d, envID: %d, envName: %s, deployed: %t", v.AppId, v.EnvId, v.EnvName, success)
			if err != nil {
				deployResult += fmt.Sprintf(", err: %s", err.Error())
			}
			result = append(result, deployResult)
		}
		return c.JSON(http.StatusOK, result)
	}
	result = append(result, "env not found, nothing deployed")
	return c.JSON(http.StatusNotFound, result)
}

// getEnvID is return a application's environment list
func (a *AppController) getEnvList(appID, pageNumber, pageSize, envType int, envName string) ([]retailcloud.AppEnvironmentResponse, error) {
	request := retailcloud.CreateListAppEnvironmentRequest()
	request.Scheme = "https"

	request.AppId = requests.NewInteger(appID)
	request.PageNumber = requests.NewInteger(pageNumber)
	request.PageSize = requests.NewInteger(pageSize)
	if envType >= 0 {
		request.EnvType = requests.NewInteger(envType)
	}
	if envName != "" {
		request.EnvName = envName
	}

	response, err := factory.GetRetailCloudClient().ListAppEnvironment(request)
	if err != nil {
		log.Errorf("Get env failed, error: %+v", err)
		return nil, err
	}

	if response.IsSuccess() {
		log.Infof("Get env succeeded, requestID: %s", response.RequestId)
		return response.Data, nil
	}

	log.Warnf("Get env failed, response: %+v", response)
	return nil, errors.New("env not found")
}

func (a *AppController) deployApp(envID, totalPartitions int, name, image string) (bool, error) {
	request := retailcloud.CreateDeployAppRequest()
	request.Scheme = "https"
	request.Name = name
	request.EnvId = requests.NewInteger(envID)
	request.TotalPartitions = requests.NewInteger(totalPartitions)

	request.ContainerImageList = &[]string{image}

	response, err := factory.GetRetailCloudClient().DeployApp(request)
	if err != nil {
		log.Errorf("Deploy app failed, error: %+v", err)
		return false, err
	}

	if response.IsSuccess() {
		log.Infof("Deploy app succeeded, requestID: %s", response.RequestId)
	} else {
		log.Warnf("Deploy app failed, response: %+v", response)
	}

	return response.IsSuccess(), nil
}
