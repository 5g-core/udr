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

// CreateCagUpdateAck - To store the CAG update acknowledgement information of a UE
func HttpCreateCagUpdateAck(c *gin.Context) {
	logger.HandlerLog.Infof("CreateCagUpdateAck is called")
    /// <param name="ueId type: string">UE id</param>
    /// <param name="cagAckData type: CagAckData"></param>
    /// <param name="supportedFeatures type: string">Supported Features (optional)</param>
    // Getting the path params
	ueId :=c.Params.ByName("ueId")
	// Getting query params 
	queryParams := c.Request.URL.Query()
	// query parameter : basename supported-features paramName supportedFeatures
	supportedFeatures := queryParams[""]
	// get and parse body
	var cagAckData CagAckData
	if err := getDataFromRequestBody(c, &cagAckData); err != nil {
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
	rsp := HandleCreateCagUpdateAck( ueId,cagAckData,supportedFeatures, c)

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
