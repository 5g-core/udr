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

// CreateIndividualAppliedBdtPolicyData - Create an individual applied BDT Policy Data resource
func HttpCreateIndividualAppliedBdtPolicyData(c *gin.Context) {
	logger.HandlerLog.Infof("CreateIndividualAppliedBdtPolicyData is called")
    /// <param name="bdtPolicyId type: string">The Identifier of an Individual Applied BDT Policy Data to be created or updated. It shall apply the format of Data type string.</param>
    /// <param name="bdtPolicyData type: BdtPolicyData"></param>
    // Getting the path params
	bdtPolicyId :=c.Params.ByName("bdtPolicyId")
	// get and parse body
	var bdtPolicyData BdtPolicyData
	if err := getDataFromRequestBody(c, &bdtPolicyData); err != nil {
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
	rsp := HandleCreateIndividualAppliedBdtPolicyData( bdtPolicyId,bdtPolicyData, c)

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

// DeleteIndividualAppliedBdtPolicyData - Delete an individual Applied BDT Policy Data resource
func HttpDeleteIndividualAppliedBdtPolicyData(c *gin.Context) {
	logger.HandlerLog.Infof("DeleteIndividualAppliedBdtPolicyData is called")
    /// <param name="bdtPolicyId type: string">The Identifier of an Individual Applied BDT Policy Data to be updated. It shall apply the format of Data type string.</param>
    // Getting the path params
	bdtPolicyId :=c.Params.ByName("bdtPolicyId")
	rsp := HandleDeleteIndividualAppliedBdtPolicyData( bdtPolicyId, c)

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

// UpdateIndividualAppliedBdtPolicyData - Modify part of the properties of an individual Applied BDT Policy Data resource
func HttpUpdateIndividualAppliedBdtPolicyData(c *gin.Context) {
	logger.HandlerLog.Infof("UpdateIndividualAppliedBdtPolicyData is called")
    /// <param name="bdtPolicyId type: string">The Identifier of an Individual Applied BDT Policy Data to be updated. It shall apply the format of Data type string.</param>
    /// <param name="bdtPolicyDataPatch type: BdtPolicyDataPatch"></param>
    // Getting the path params
	bdtPolicyId :=c.Params.ByName("bdtPolicyId")
	// get and parse body
	var bdtPolicyDataPatch BdtPolicyDataPatch
	if err := getDataFromRequestBody(c, &bdtPolicyDataPatch); err != nil {
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
	rsp := HandleUpdateIndividualAppliedBdtPolicyData( bdtPolicyId,bdtPolicyDataPatch, c)

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