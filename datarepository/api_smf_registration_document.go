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

// CreateSmfContextNon3gpp - To create an individual SMF context data of a UE in the UDR
func HttpCreateSmfContextNon3gpp(c *gin.Context) {
	logger.HandlerLog.Infof("CreateSmfContextNon3gpp is called")
    /// <param name="ueId type: string">UE id</param>
    /// <param name="pduSessionId type: int32">PDU session id</param>
    /// <param name="smfRegistration type: SmfRegistration"></param>
    // Getting the path params
	ueId :=c.Params.ByName("ueId")
	pduSessionId :=c.Params.ByName("pduSessionId")
	// get and parse body
	var smfRegistration SmfRegistration
	if err := getDataFromRequestBody(c, &smfRegistration); err != nil {
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
	rsp := HandleCreateSmfContextNon3gpp( ueId,pduSessionId,smfRegistration, c)

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

// DeleteSmfContext - To remove an individual SMF context data of a UE the UDR
func HttpDeleteSmfContext(c *gin.Context) {
	logger.HandlerLog.Infof("DeleteSmfContext is called")
    /// <param name="ueId type: string">UE id</param>
    /// <param name="pduSessionId type: int32">PDU session id</param>
    // Getting the path params
	ueId :=c.Params.ByName("ueId")
	pduSessionId :=c.Params.ByName("pduSessionId")
	rsp := HandleDeleteSmfContext( ueId,pduSessionId, c)

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

// QuerySmfRegistration - Retrieves the individual SMF registration of a UE
func HttpQuerySmfRegistration(c *gin.Context) {
	logger.HandlerLog.Infof("QuerySmfRegistration is called")
    /// <param name="ueId type: string">UE id</param>
    /// <param name="pduSessionId type: int32">PDU session id</param>
    /// <param name="fields type: []string">attributes to be retrieved (optional)</param>
    /// <param name="supportedFeatures type: string">Supported Features (optional)</param>
    // Getting the path params
	ueId :=c.Params.ByName("ueId")
	pduSessionId :=c.Params.ByName("pduSessionId")
	// Getting query params 
	queryParams := c.Request.URL.Query()
	// query parameter : basename fields paramName fields
	fields := queryParams[""]
	// query parameter : basename supported-features paramName supportedFeatures
	supportedFeatures := queryParams[""]
	rsp := HandleQuerySmfRegistration( ueId,pduSessionId,fields,supportedFeatures, c)

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
