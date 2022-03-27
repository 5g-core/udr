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

// ReadOperatorSpecificData - Retrieve the operator specific policy data of an UE
func HttpReadOperatorSpecificData(c *gin.Context) {
	logger.HandlerLog.Infof("ReadOperatorSpecificData is called")
    /// <param name="ueId type: string">UE Id</param>
    /// <param name="fields type: []string">attributes to be retrieved (optional)</param>
    /// <param name="suppFeat type: string">Supported Features (optional)</param>
    // Getting the path params
	ueId :=c.Params.ByName("ueId")
	// Getting query params 
	queryParams := c.Request.URL.Query()
	// query parameter : basename fields paramName fields
	fields := queryParams[""]
	// query parameter : basename supp-feat paramName suppFeat
	suppFeat := queryParams[""]
	rsp := HandleReadOperatorSpecificData( ueId,fields,suppFeat, c)

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

// ReplaceOperatorSpecificData - Modify the operator specific policy data of an UE
func HttpReplaceOperatorSpecificData(c *gin.Context) {
	logger.HandlerLog.Infof("ReplaceOperatorSpecificData is called")
    /// <param name="ueId type: string">UE Id</param>
    /// <param name="requestBody type: map[string]OperatorSpecificDataContainer"></param>
    // Getting the path params
	ueId :=c.Params.ByName("ueId")
	// get and parse body
	var requestBody map[string]OperatorSpecificDataContainer
	if err := getDataFromRequestBody(c, &requestBody); err != nil {
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
	rsp := HandleReplaceOperatorSpecificData( ueId,requestBody, c)

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

// UpdateOperatorSpecificData - Modify the operator specific policy data of an UE
func HttpUpdateOperatorSpecificData(c *gin.Context) {
	logger.HandlerLog.Infof("UpdateOperatorSpecificData is called")
    /// <param name="ueId type: string">UE Id</param>
    /// <param name="patchItem type: []PatchItem"></param>
    // Getting the path params
	ueId :=c.Params.ByName("ueId")
	// get and parse body
	var patchItem []PatchItem
	if err := getDataFromRequestBody(c, &patchItem); err != nil {
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
	rsp := HandleUpdateOperatorSpecificData( ueId,patchItem, c)

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
