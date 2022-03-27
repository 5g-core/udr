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

// ModifyEeGroupSubscription - Modify an individual ee subscription for a group of a UEs
func HttpModifyEeGroupSubscription(c *gin.Context) {
	logger.HandlerLog.Infof("ModifyEeGroupSubscription is called")
    /// <param name="ueGroupId type: string"></param>
    /// <param name="subsId type: string"></param>
    /// <param name="patchItem type: []PatchItem"></param>
    /// <param name="supportedFeatures type: string">Features required to be supported by the target NF (optional)</param>
    // Getting the path params
	ueGroupId :=c.Params.ByName("ueGroupId")
	subsId :=c.Params.ByName("subsId")
	// Getting query params 
	queryParams := c.Request.URL.Query()
	// query parameter : basename supported-features paramName supportedFeatures
	supportedFeatures := queryParams[""]
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
	rsp := HandleModifyEeGroupSubscription( ueGroupId,subsId,patchItem,supportedFeatures, c)

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

// QueryEeGroupSubscription - Retrieve a individual eeSubscription for a group of UEs or any UE
func HttpQueryEeGroupSubscription(c *gin.Context) {
	logger.HandlerLog.Infof("QueryEeGroupSubscription is called")
    /// <param name="ueGroupId type: string"></param>
    /// <param name="subsId type: string">Unique ID of the subscription to remove</param>
    // Getting the path params
	ueGroupId :=c.Params.ByName("ueGroupId")
	subsId :=c.Params.ByName("subsId")
	rsp := HandleQueryEeGroupSubscription( ueGroupId,subsId, c)

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

// RemoveEeGroupSubscriptions - Deletes a eeSubscription for a group of UEs or any UE
func HttpRemoveEeGroupSubscriptions(c *gin.Context) {
	logger.HandlerLog.Infof("RemoveEeGroupSubscriptions is called")
    /// <param name="ueGroupId type: string"></param>
    /// <param name="subsId type: string">Unique ID of the subscription to remove</param>
    // Getting the path params
	ueGroupId :=c.Params.ByName("ueGroupId")
	subsId :=c.Params.ByName("subsId")
	rsp := HandleRemoveEeGroupSubscriptions( ueGroupId,subsId, c)

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

// UpdateEeGroupSubscriptions - Update an individual ee subscription of a group of UEs or any UE
func HttpUpdateEeGroupSubscriptions(c *gin.Context) {
	logger.HandlerLog.Infof("UpdateEeGroupSubscriptions is called")
    /// <param name="ueGroupId type: string"></param>
    /// <param name="subsId type: string"></param>
    /// <param name="eeSubscription type: EeSubscription"></param>
    // Getting the path params
	ueGroupId :=c.Params.ByName("ueGroupId")
	subsId :=c.Params.ByName("subsId")
	// get and parse body
	var eeSubscription EeSubscription
	if err := getDataFromRequestBody(c, &eeSubscription); err != nil {
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
	rsp := HandleUpdateEeGroupSubscriptions( ueGroupId,subsId,eeSubscription, c)

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
