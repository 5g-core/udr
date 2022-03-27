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

// QuerySmData - Retrieves the Session Management subscription data of a UE
func HttpQuerySmData(c *gin.Context) {
	logger.HandlerLog.Infof("QuerySmData is called")
    /// <param name="ueId type: string">UE id</param>
    /// <param name="servingPlmnId type: string">PLMN ID</param>
    /// <param name="singleNssai type: VarSnssai">single NSSAI (optional)</param>
    /// <param name="dnn type: string">DNN (optional)</param>
    /// <param name="fields type: []string">attributes to be retrieved (optional)</param>
    /// <param name="supportedFeatures type: string">Supported Features (optional)</param>
    /// <param name="ifNoneMatch type: string">Validator for conditional requests, as described in RFC 7232, 3.2 (optional)</param>
    /// <param name="ifModifiedSince type: string">Validator for conditional requests, as described in RFC 7232, 3.3 (optional)</param>
    // Getting the path params
	ueId :=c.Params.ByName("ueId")
	servingPlmnId :=c.Params.ByName("servingPlmnId")
	// Getting query params 
	queryParams := c.Request.URL.Query()
	// query parameter : basename single-nssai paramName singleNssai
	singleNssai := queryParams[""]
	// query parameter : basename dnn paramName dnn
	dnn := queryParams[""]
	// query parameter : basename fields paramName fields
	fields := queryParams[""]
	// query parameter : basename supported-features paramName supportedFeatures
	supportedFeatures := queryParams[""]
	// Getting the header params
	ifNoneMatch :=  c.Request.Header["ifNoneMatch"]	
	ifModifiedSince :=  c.Request.Header["ifModifiedSince"]	
	rsp := HandleQuerySmData( ueId,servingPlmnId,singleNssai,dnn,fields,supportedFeatures,ifNoneMatch,ifModifiedSince, c)

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
