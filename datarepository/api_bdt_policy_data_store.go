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

// ReadBdtPolicyData - Retrieve applied BDT Policy Data
func HttpReadBdtPolicyData(c *gin.Context) {
	logger.HandlerLog.Infof("ReadBdtPolicyData is called")
    /// <param name="bdtPolicyIds type: []string">Each element identifies a service. (optional)</param>
    /// <param name="internalGroupIds type: []string">Each element identifies a group of users. (optional)</param>
    /// <param name="supis type: []string">Each element identifies the user. (optional)</param>
	// Getting query params 
	queryParams := c.Request.URL.Query()
	// query parameter : basename bdt-policy-ids paramName bdtPolicyIds
	bdtPolicyIds := queryParams[""]
	// query parameter : basename internal-group-ids paramName internalGroupIds
	internalGroupIds := queryParams[""]
	// query parameter : basename supis paramName supis
	supis := queryParams[""]
	rsp := HandleReadBdtPolicyData( bdtPolicyIds,internalGroupIds,supis, c)

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
