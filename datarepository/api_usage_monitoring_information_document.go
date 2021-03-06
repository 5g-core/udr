/*
 * Nudr_DataRepository API OpenAPI file
 *
 * Unified Data Repository Service. © 2021, 3GPP Organizational Partners (ARIB, ATIS, CCSA, ETSI, TSDSI, TTA, TTC). All rights reserved.
 *
 * API version: 2.1.5
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package datarepository

import (
	"net/http"

	"github.com/5g-core/openapi"
	. "github.com/5g-core/openapi/models"
	"github.com/5g-core/udr/logger"
	. "github.com/5g-core/udr/producer"
	"github.com/gin-gonic/gin"
)

// CreateUsageMonitoringResource - Create a usage monitoring resource
func HttpCreateUsageMonitoringResource(c *gin.Context) {
	logger.HandlerLog.Infof("CreateUsageMonitoringResource is called")
    /// <param name="ueId type: string"></param>
    /// <param name="usageMonId type: string"></param>
    /// <param name="usageMonData type: UsageMonData"></param>
    // Getting the path params
	ueId :=c.Params.ByName("ueId")
	usageMonId :=c.Params.ByName("usageMonId")
	// get and parse body
	var usageMonData UsageMonData
	if err := getDataFromRequestBody(c, &usageMonData); err != nil {
		problemDetail := "[Request Body] " + err.Error()
		var status int32
		status =int32(http.StatusBadRequest)
		title:="Malformed request syntax"
		rsp := ProblemDetails{
			Title:  &title,
			Status: &status,
			Detail: &problemDetail,
		}
		logger.DataRepoLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
	}
	rsp := HandleCreateUsageMonitoringResource( ueId,usageMonId,usageMonData, c)

	// send response
	for k, v := range rsp.Header {
		// TODO: concatenate all values
		c.Header(k, v[0])
	}
	serializedBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.DataRepoLog.Errorf("Serialize Response Body error: %+v", err)
		problemDetail := "Serialize Response Body error: " + err.Error()
		var status int32
		status =int32(http.StatusBadRequest)
		title:="Malformed request syntax"
		rsp := ProblemDetails{
			Title:  &title,
			Status: &status,
			Detail: &problemDetail,
		}
		logger.DataRepoLog.Errorln(problemDetail)
		c.JSON(http.StatusInternalServerError, rsp)
	} else {
		c.Data(rsp.Status, "application/json", serializedBody)
	}
}

// DeleteUsageMonitoringInformation - Delete a usage monitoring resource
func HttpDeleteUsageMonitoringInformation(c *gin.Context) {
	logger.HandlerLog.Infof("DeleteUsageMonitoringInformation is called")
    /// <param name="ueId type: string"></param>
    /// <param name="usageMonId type: string"></param>
    // Getting the path params
	ueId :=c.Params.ByName("ueId")
	usageMonId :=c.Params.ByName("usageMonId")
	rsp := HandleDeleteUsageMonitoringInformation( ueId,usageMonId, c)

	// send response
	for k, v := range rsp.Header {
		// TODO: concatenate all values
		c.Header(k, v[0])
	}
	serializedBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.DataRepoLog.Errorf("Serialize Response Body error: %+v", err)
		problemDetail := "Serialize Response Body error: " + err.Error()
		var status int32
		status =int32(http.StatusBadRequest)
		title:="Malformed request syntax"
		rsp := ProblemDetails{
			Title:  &title,
			Status: &status,
			Detail: &problemDetail,
		}
		logger.DataRepoLog.Errorln(problemDetail)
		c.JSON(http.StatusInternalServerError, rsp)
	} else {
		c.Data(rsp.Status, "application/json", serializedBody)
	}
}

// ReadUsageMonitoringInformation - Retrieve a usage monitoring resource
func HttpReadUsageMonitoringInformation(c *gin.Context) {
	logger.HandlerLog.Infof("ReadUsageMonitoringInformation is called")
    /// <param name="ueId type: string"></param>
    /// <param name="usageMonId type: string"></param>
    /// <param name="suppFeat type: string">Supported Features (optional)</param>
    // Getting the path params
	ueId :=c.Params.ByName("ueId")
	usageMonId :=c.Params.ByName("usageMonId")
	// Getting query params 
	queryParams := c.Request.URL.Query()
	// query parameter : basename supp-feat paramName suppFeat
	suppFeat := queryParams[""]
	rsp := HandleReadUsageMonitoringInformation( ueId,usageMonId,suppFeat, c)

	// send response
	for k, v := range rsp.Header {
		// TODO: concatenate all values
		c.Header(k, v[0])
	}
	serializedBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.DataRepoLog.Errorf("Serialize Response Body error: %+v", err)
		problemDetail := "Serialize Response Body error: " + err.Error()
		var status int32
		status =int32(http.StatusBadRequest)
		title:="Malformed request syntax"
		rsp := ProblemDetails{
			Title:  &title,
			Status: &status,
			Detail: &problemDetail,
		}
		logger.DataRepoLog.Errorln(problemDetail)
		c.JSON(http.StatusInternalServerError, rsp)
	} else {
		c.Data(rsp.Status, "application/json", serializedBody)
	}
}
