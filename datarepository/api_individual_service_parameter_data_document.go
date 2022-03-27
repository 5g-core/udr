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

// CreateOrReplaceServiceParameterData - Create or update an individual Service Parameter Data resource
func HttpCreateOrReplaceServiceParameterData(c *gin.Context) {
	logger.HandlerLog.Infof("CreateOrReplaceServiceParameterData is called")
    /// <param name="serviceParamId type: string">The Identifier of an Individual Service Parameter Data to be created or updated. It shall apply the format of Data type string.</param>
    /// <param name="serviceParameterData type: ServiceParameterData"></param>
    // Getting the path params
	serviceParamId :=c.Params.ByName("serviceParamId")
	// get and parse body
	var serviceParameterData ServiceParameterData
	if err := getDataFromRequestBody(c, &serviceParameterData); err != nil {
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
	rsp := HandleCreateOrReplaceServiceParameterData( serviceParamId,serviceParameterData, c)

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

// DeleteIndividualServiceParameterData - Delete an individual Service Parameter Data resource
func HttpDeleteIndividualServiceParameterData(c *gin.Context) {
	logger.HandlerLog.Infof("DeleteIndividualServiceParameterData is called")
    /// <param name="serviceParamId type: string">The Identifier of an Individual Service Parameter Data to be updated. It shall apply the format of Data type string.</param>
    // Getting the path params
	serviceParamId :=c.Params.ByName("serviceParamId")
	rsp := HandleDeleteIndividualServiceParameterData( serviceParamId, c)

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

// UpdateIndividualServiceParameterData - Modify part of the properties of an individual Service Parameter Data resource
func HttpUpdateIndividualServiceParameterData(c *gin.Context) {
	logger.HandlerLog.Infof("UpdateIndividualServiceParameterData is called")
    /// <param name="serviceParamId type: string">The Identifier of an Individual Service Parameter Data to be updated. It shall apply the format of Data type string.</param>
    /// <param name="serviceParameterDataPatch type: ServiceParameterDataPatch"></param>
    // Getting the path params
	serviceParamId :=c.Params.ByName("serviceParamId")
	// get and parse body
	var serviceParameterDataPatch ServiceParameterDataPatch
	if err := getDataFromRequestBody(c, &serviceParameterDataPatch); err != nil {
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
	rsp := HandleUpdateIndividualServiceParameterData( serviceParamId,serviceParameterDataPatch, c)

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
