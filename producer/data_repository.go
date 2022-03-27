package producer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/5g-core/mongodblibrary"
	"github.com/5g-core/http_wrapper"
	. "github.com/5g-core/openapi/models"
	udr_context "github.com/5g-core/udr/context"
	"github.com/5g-core/udr/logger"
	"github.com/5g-core/udr/util"
	jsonpatch "github.com/evanphx/json-patch"
)

const (
	APPDATA_INFLUDATA_DB_COLLECTION_NAME       = "applicationData.influenceData"
	APPDATA_INFLUDATA_SUBSC_DB_COLLECTION_NAME = "applicationData.influenceData.subsToNotify"
	APPDATA_PFD_DB_COLLECTION_NAME             = "applicationData.pfds"
)

var CurrentResourceUri string

func getDataFromDB(collName string, filter bson.M) (map[string]interface{}, *ProblemDetails) {
	logger.DataRepoLog.Infof("getDataFromDB")
	data := mongodblibrary.RestfulAPIGetOne(collName, filter)
	if data == nil {
		return nil, util.ProblemDetailsNotFound("DATA_NOT_FOUND")
	}

	// Delete "_id" entry which is auto-inserted by MongoDB
	delete(data, "_id")
	return data, nil
}
func deleteDataFromDB(collName string, filter bson.M) {
	mongodblibrary.RestfulAPIDeleteOne(collName, filter)
}

func HandleAmfContext3gpp(ueId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {

	logger.DataRepoLog.Infof("Handle AmfContext3gpp")
	collName := "subscriptionData.contextData.amf3gppAccess"

	problemDetails := AmfContext3gppProcedure(collName, ueId, patchItem)
	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

}
func AmfContext3gppProcedure(collName string, ueId string, patchItem []PatchItem) *ProblemDetails {
	filter := bson.M{"ueId": ueId}
	origValue := mongodblibrary.RestfulAPIGetOne(collName, filter)

	patchJSON, err := json.Marshal(patchItem)
	if err != nil {
		logger.DataRepoLog.Error(err)
	}
	success := mongodblibrary.RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		newValue := mongodblibrary.RestfulAPIGetOne(collName, filter)
		PreHandleOnDataChangeNotify(ueId, CurrentResourceUri, patchItem, origValue, newValue)
		return nil
	} else {
		return util.ProblemDetailsModifyNotAllowed("")
	}
}

func HandleCreateAmfContext3gpp(ueId string, amf3GppAccessRegistration Amf3GppAccessRegistration, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateAmfContext3gpp")
	collName := "subscriptionData.contextData.amf3gppAccess"
	CreateAmfContext3gppProcedure(collName, ueId, amf3GppAccessRegistration)
	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}
func CreateAmfContext3gppProcedure(collName string, ueId string,
	amf3GppAccessRegistration Amf3GppAccessRegistration) {
	filter := bson.M{"ueId": ueId}
	putData := util.ToBsonM(amf3GppAccessRegistration)
	putData["ueId"] = ueId
	mongodblibrary.RestfulAPIPutOne(collName, filter, putData)
}
func HandleQueryAmfContext3gpp(ueId string, fields []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryAmfContext3gpp")
	collName := "subscriptionData.contextData.amf3gppAccess"
	response, problemDetails := QueryAmfContext3gppProcedure(collName, ueId)
	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}
func QueryAmfContext3gppProcedure(collName string, ueId string) (*map[string]interface{}, *ProblemDetails) {
	filter := bson.M{"ueId": ueId}
	amf3GppAccessRegistration := mongodblibrary.RestfulAPIGetOne(collName, filter)
	delete(amf3GppAccessRegistration, "_id")
	if amf3GppAccessRegistration != nil {
		return &amf3GppAccessRegistration, nil
	} else {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
}

func HandleAmfContextNon3gpp(ueId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle AmfContextNon3gpp")

	collName := "subscriptionData.contextData.amfNon3gppAccess"

	filter := bson.M{"ueId": ueId}

	problemDetails := AmfContextNon3gppProcedure(ueId, collName, patchItem, filter)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
}

func AmfContextNon3gppProcedure(ueId string, collName string, patchItem []PatchItem,
	filter bson.M) *ProblemDetails {
	origValue := mongodblibrary.RestfulAPIGetOne(collName, filter)

	patchJSON, err := json.Marshal(patchItem)
	if err != nil {
		logger.DataRepoLog.Error(err)
	}
	success := mongodblibrary.RestfulAPIJSONPatch(collName, filter, patchJSON)
	if success {
		newValue := mongodblibrary.RestfulAPIGetOne(collName, filter)
		PreHandleOnDataChangeNotify(ueId, CurrentResourceUri, patchItem, origValue, newValue)
		return nil
	} else {
		return util.ProblemDetailsModifyNotAllowed("")
	}
}
func HandleCreateAmfContextNon3gpp(ueId string, amfNon3GppAccessRegistration AmfNon3GppAccessRegistration, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateAmfContextNon3gpp")
	collName := "subscriptionData.contextData.amfNon3gppAccess"
	CreateAmfContextNon3gppProcedure(amfNon3GppAccessRegistration, collName, ueId)
	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func CreateAmfContextNon3gppProcedure(AmfNon3GppAccessRegistration AmfNon3GppAccessRegistration,
	collName string, ueId string) {
	putData := util.ToBsonM(AmfNon3GppAccessRegistration)
	putData["ueId"] = ueId
	filter := bson.M{"ueId": ueId}

	mongodblibrary.RestfulAPIPutOne(collName, filter, putData)
}

func HandleQueryAmfContextNon3gpp(ueId string, fields []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryAmfContextNon3gpp")

	collName := "subscriptionData.contextData.amfNon3gppAccess"

	response, problemDetails := QueryAmfContextNon3gppProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func QueryAmfContextNon3gppProcedure(collName string, ueId string) (*map[string]interface{}, *ProblemDetails) {
	filter := bson.M{"ueId": ueId}
	response := mongodblibrary.RestfulAPIGetOne(collName, filter)

	if response != nil {
		return &response, nil
	} else {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
}

func HandleCreateAMFSubscriptions(ueId string, subsId string, amfSubscriptionInfo []AmfSubscriptionInfo, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateAMFSubscriptions")

	problemDetails := CreateAMFSubscriptionsProcedure(subsId, ueId, amfSubscriptionInfo)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
}

func CreateAMFSubscriptionsProcedure(subsId string, ueId string,
	AmfSubscriptionInfo []AmfSubscriptionInfo) *ProblemDetails {
	udrSelf := udr_context.UDR_Self()
	value, ok := udrSelf.UESubsCollection.Load(ueId)
	if !ok {
		return util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
	UESubsData := value.(*udr_context.UESubsData)

	_, ok = UESubsData.EeSubscriptionCollection[subsId]
	if !ok {
		return util.ProblemDetailsNotFound("SUBSCRIPTION_NOT_FOUND")
	}

	UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos = AmfSubscriptionInfo
	return nil
}
func HandleCreateOrReplaceAccessAndMobilityData(ueId string, accessAndMobilityData AccessAndMobilityData, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleDeleteAccessAndMobilityData(ueId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQueryAccessAndMobilityData(ueId string, suppFeat []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleUpdateAccessAndMobilityData(ueId string, accessAndMobilityData AccessAndMobilityData, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleReadAccessAndMobilityPolicyData(ueId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQueryAmData(ueId string, servingPlmnId string, fields []string, supportedFeatures []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryAmData")

	collName := "subscriptionData.provisionedData.amData"

	response, problemDetails := QueryAmDataProcedure(collName, ueId, servingPlmnId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
}

func QueryAmDataProcedure(collName string, ueId string, servingPlmnId string) (*map[string]interface{},
	*ProblemDetails) {
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	accessAndMobilitySubscriptionData := mongodblibrary.RestfulAPIGetOne(collName, filter)
	if accessAndMobilitySubscriptionData != nil {
		return &accessAndMobilitySubscriptionData, nil
	} else {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
}

func HandleModifyAmfSubscriptionInfo(ueId string, subsId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle ModifyAmfSubscriptionInfo")

	problemDetails := ModifyAmfSubscriptionInfoProcedure(ueId, subsId, patchItem)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
}

func ModifyAmfSubscriptionInfoProcedure(ueId string, subsId string,
	patchItem []PatchItem) *ProblemDetails {
	udrSelf := udr_context.UDR_Self()
	value, ok := udrSelf.UESubsCollection.Load(ueId)
	if !ok {
		return util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
	UESubsData := value.(*udr_context.UESubsData)

	_, ok = UESubsData.EeSubscriptionCollection[subsId]

	if !ok {
		return util.ProblemDetailsNotFound("SUBSCRIPTION_NOT_FOUND")
	}

	if UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos == nil {
		return util.ProblemDetailsNotFound("AMFSUBSCRIPTION_NOT_FOUND")
	}
	var patchJSON []byte
	if patchJSONtemp, err := json.Marshal(patchItem); err != nil {
		logger.DataRepoLog.Errorln(err)
	} else {
		patchJSON = patchJSONtemp
	}
	var patch jsonpatch.Patch
	if patchtemp, err := jsonpatch.DecodePatch(patchJSON); err != nil {
		logger.DataRepoLog.Errorln(err)
		return util.ProblemDetailsModifyNotAllowed("PatchItem attributes are invalid")
	} else {
		patch = patchtemp
	}
	original, err := json.Marshal((UESubsData.EeSubscriptionCollection[subsId]).AmfSubscriptionInfos)
	if err != nil {
		logger.DataRepoLog.Warnln(err)
	}

	modified, err := patch.Apply(original)
	if err != nil {
		return util.ProblemDetailsModifyNotAllowed("Occur error when applying PatchItem")
	}
	var modifiedData []AmfSubscriptionInfo
	err = json.Unmarshal(modified, &modifiedData)
	if err != nil {
		logger.DataRepoLog.Error(err)
	}

	UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos = modifiedData
	return nil
}

func HandleCreateIndividualApplicationDataSubscription(applicationDataSubs ApplicationDataSubs, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleReadApplicationDataChangeSubscriptions(dataFilter []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleDeleteAuthenticationStatus(ueId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQueryAuthenticationStatus(ueId string, fields []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryAuthenticationStatus")

	collName := "subscriptionData.authenticationData.authenticationStatus"

	response, problemDetails := QueryAuthenticationStatusProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func QueryAuthenticationStatusProcedure(collName string, ueId string) (*map[string]interface{},
	*ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	authEvent := mongodblibrary.RestfulAPIGetOne(collName, filter)

	if authEvent != nil {
		return &authEvent, nil
	} else {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
}

func HandleQueryAuthSubsData(ueId string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryAuthSubsData")

	collName := "subscriptionData.authenticationData.authenticationSubscription"

	response, problemDetails := QueryAuthSubsDataProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func QueryAuthSubsDataProcedure(collName string, ueId string) (map[string]interface{}, *ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	authenticationSubscription := mongodblibrary.RestfulAPIGetOne(collName, filter)

	if authenticationSubscription != nil {
		return authenticationSubscription, nil
	} else {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
}

func HandleCreateAuthenticationSoR(ueId string, sorData SorData, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateAuthenticationSoR")

	putData := util.ToBsonM(sorData)

	collName := "subscriptionData.ueUpdateConfirmationData.sorData"

	CreateAuthenticationSoRProcedure(collName, ueId, putData)

	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func CreateAuthenticationSoRProcedure(collName string, ueId string, putData bson.M) {
	filter := bson.M{"ueId": ueId}
	putData["ueId"] = ueId

	mongodblibrary.RestfulAPIPutOne(collName, filter, putData)
}

func HandleQueryAuthSoR(ueId string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryAuthSoR")

	collName := "subscriptionData.ueUpdateConfirmationData.sorData"

	response, problemDetails := QueryAuthSoRProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func QueryAuthSoRProcedure(collName string, ueId string) (map[string]interface{}, *ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	sorData := mongodblibrary.RestfulAPIGetOne(collName, filter)

	if sorData != nil {
		return sorData, nil
	} else {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
}

func HandleCreateAuthenticationStatus(ueId string, authEvent AuthEvent, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateAuthenticationStatus")

	putData := util.ToBsonM(authEvent)

	collName := "subscriptionData.authenticationData.authenticationStatus"

	CreateAuthenticationStatusProcedure(collName, ueId, putData)

	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func CreateAuthenticationStatusProcedure(collName string, ueId string, putData bson.M) {
	filter := bson.M{"ueId": ueId}
	putData["ueId"] = ueId

	mongodblibrary.RestfulAPIPutOne(collName, filter, putData)
}

func HandleModifyAuthenticationSubscription(ueId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleCreateSharedData(sharedDataId string, sharedData SharedData, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateSharedData")
	collName := "subscriptionData.sharedData"
	CreateSharedDataProcedure(collName, sharedDataId, sharedData)
	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func CreateSharedDataProcedure(collName string, sharedDataId string, sharedData SharedData) {
	filter := bson.M{"sharedDataId": sharedDataId}
	putData := util.ToBsonM(sharedData)
	putData["sharedDataId"] = sharedDataId
	mongodblibrary.RestfulAPIPutOne(collName, filter, putData)
}

func HandleCreateAuthenticationUPU(ueId string, supportedFeatures []string, upuData UpuData, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQueryAuthUPU(ueId string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleReadBdtData(bdtRefIds []string, suppFeat []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleReadBdtPolicyData(bdtPolicyIds []string, internalGroupIds []string, supis []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQueryCagAck(ueId string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateCagUpdateAck(ueId string, cagAckData CagAckData, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQuery5GVnGroupInternal(internalGroupIds []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQuery5GVnGroup(gpsis []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreate5GVnGroup(externalGroupId string, model5GvnGroupConfiguration Model5GvnGroupConfiguration, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQueryContextData(ueId string, contextDatasetNames []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleDelete5GVnGroup(externalGroupId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQueryCoverageRestrictionData(ueId string, supportedFeatures []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleRemoveAmfSubscriptionsInfo(ueId string, subsId string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle RemoveAmfSubscriptionsInfo")

	problemDetails := RemoveAmfSubscriptionsInfoProcedure(subsId, ueId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
}

func RemoveAmfSubscriptionsInfoProcedure(subsId string, ueId string) *ProblemDetails {
	udrSelf := udr_context.UDR_Self()
	value, ok := udrSelf.UESubsCollection.Load(ueId)
	if !ok {
		return util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}

	UESubsData := value.(*udr_context.UESubsData)
	_, ok = UESubsData.EeSubscriptionCollection[subsId]

	if !ok {
		return util.ProblemDetailsNotFound("SUBSCRIPTION_NOT_FOUND")
	}

	if UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos == nil {
		return util.ProblemDetailsNotFound("AMFSUBSCRIPTION_NOT_FOUND")
	}

	UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos = nil

	return nil
}

func HandleQueryEEData(ueId string, fields []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryEEData")

	collName := "subscriptionData.eeProfileData"

	response, problemDetails := QueryEEDataProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func QueryEEDataProcedure(collName string, ueId string) (*map[string]interface{}, *ProblemDetails) {
	filter := bson.M{"ueId": ueId}
	eeProfileData := mongodblibrary.RestfulAPIGetOne(collName, filter)

	if eeProfileData != nil {
		return &eeProfileData, nil
	} else {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
}

func HandleQueryGroupEEData(ueGroupId string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleModifyEeGroupSubscription(ueGroupId string, subsId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQueryEeGroupSubscription(ueGroupId string, subsId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleRemoveEeGroupSubscriptions(ueGroupId string, subsId string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle RemoveEeGroupSubscriptions")

	problemDetails := RemoveEeGroupSubscriptionsProcedure(ueGroupId, subsId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
}

func RemoveEeGroupSubscriptionsProcedure(ueGroupId string, subsId string) *ProblemDetails {
	udrSelf := udr_context.UDR_Self()
	value, ok := udrSelf.UEGroupCollection.Load(ueGroupId)
	if !ok {
		return util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}

	UEGroupSubsData := value.(*udr_context.UEGroupSubsData)
	_, ok = UEGroupSubsData.EeSubscriptions[subsId]

	if !ok {
		return util.ProblemDetailsNotFound("SUBSCRIPTION_NOT_FOUND")
	}
	delete(UEGroupSubsData.EeSubscriptions, subsId)

	return nil
}
func HandleUpdateEeGroupSubscriptions(ueGroupId string, subsId string, eeSubscription EeSubscription, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle UpdateEeGroupSubscriptions")

	problemDetails := UpdateEeGroupSubscriptionsProcedure(ueGroupId, subsId, eeSubscription)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
}

func UpdateEeGroupSubscriptionsProcedure(ueGroupId string, subsId string,
	EeSubscription EeSubscription) *ProblemDetails {
	udrSelf := udr_context.UDR_Self()
	value, ok := udrSelf.UEGroupCollection.Load(ueGroupId)
	if !ok {
		return util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}

	UEGroupSubsData := value.(*udr_context.UEGroupSubsData)
	_, ok = UEGroupSubsData.EeSubscriptions[subsId]

	if !ok {
		return util.ProblemDetailsNotFound("SUBSCRIPTION_NOT_FOUND")
	}
	UEGroupSubsData.EeSubscriptions[subsId] = &EeSubscription

	return nil
}

func HandleCreateEeGroupSubscriptions(ueGroupId string, eeSubscription EeSubscription, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateEeGroupSubscriptions")

	locationHeader := CreateEeGroupSubscriptionsProcedure(ueGroupId, eeSubscription)

	headers := http.Header{}
	headers.Set("Location", locationHeader)
	return http_wrapper.NewResponse(http.StatusCreated, headers, eeSubscription)
}

func CreateEeGroupSubscriptionsProcedure(ueGroupId string, eeSubscription EeSubscription) string {
	udrSelf := udr_context.UDR_Self()

	value, ok := udrSelf.UEGroupCollection.Load(ueGroupId)
	if !ok {
		udrSelf.UEGroupCollection.Store(ueGroupId, new(udr_context.UEGroupSubsData))
		value, _ = udrSelf.UEGroupCollection.Load(ueGroupId)
	}
	UEGroupSubsData := value.(*udr_context.UEGroupSubsData)
	if UEGroupSubsData.EeSubscriptions == nil {
		UEGroupSubsData.EeSubscriptions = make(map[string]*EeSubscription)
	}

	newSubscriptionID := strconv.Itoa(udrSelf.EeSubscriptionIDGenerator)
	UEGroupSubsData.EeSubscriptions[newSubscriptionID] = &eeSubscription
	udrSelf.EeSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/nudr-dr/v1/subscription-data/group-data/{ueGroupId}/ee-subscriptions */
	locationHeader := fmt.Sprintf("%s/nudr-dr/v1/subscription-data/group-data/%s/ee-subscriptions/%s",
		udrSelf.GetIPv4GroupUri(udr_context.NUDR_DR), ueGroupId, newSubscriptionID)

	return locationHeader
}
func HandleQueryEeGroupSubscriptions(ueGroupId string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryEeGroupSubscriptions")

	response, problemDetails := QueryEeGroupSubscriptionsProcedure(ueGroupId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func QueryEeGroupSubscriptionsProcedure(ueGroupId string) ([]EeSubscription, *ProblemDetails) {
	udrSelf := udr_context.UDR_Self()

	value, ok := udrSelf.UEGroupCollection.Load(ueGroupId)
	if !ok {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}

	UEGroupSubsData := value.(*udr_context.UEGroupSubsData)
	var eeSubscriptionSlice []EeSubscription

	for _, v := range UEGroupSubsData.EeSubscriptions {
		eeSubscriptionSlice = append(eeSubscriptionSlice, *v)
	}
	return eeSubscriptionSlice, nil
}

func HandleModifyEesubscription(ueId string, subsId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQueryeeSubscription(ueId string, subsId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleRemoveeeSubscriptions(ueId string, subsId string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle RemoveeeSubscriptions")

	problemDetails := RemoveeeSubscriptionsProcedure(ueId, subsId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
}

func RemoveeeSubscriptionsProcedure(ueId string, subsId string) *ProblemDetails {
	udrSelf := udr_context.UDR_Self()
	value, ok := udrSelf.UESubsCollection.Load(ueId)
	if !ok {
		return util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}

	UESubsData := value.(*udr_context.UESubsData)
	_, ok = UESubsData.EeSubscriptionCollection[subsId]

	if !ok {
		return util.ProblemDetailsNotFound("SUBSCRIPTION_NOT_FOUND")
	}
	delete(UESubsData.EeSubscriptionCollection, subsId)
	return nil
}
func HandleUpdateEesubscriptions(ueId string, subsId string, eeSubscription EeSubscription, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle UpdateEesubscriptions")

	problemDetails := UpdateEesubscriptionsProcedure(ueId, subsId, eeSubscription)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
}

func UpdateEesubscriptionsProcedure(ueId string, subsId string,
	EeSubscription EeSubscription) *ProblemDetails {
	udrSelf := udr_context.UDR_Self()
	value, ok := udrSelf.UESubsCollection.Load(ueId)
	if !ok {
		return util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}

	UESubsData := value.(*udr_context.UESubsData)
	_, ok = UESubsData.EeSubscriptionCollection[subsId]

	if !ok {
		return util.ProblemDetailsNotFound("SUBSCRIPTION_NOT_FOUND")
	}
	UESubsData.EeSubscriptionCollection[subsId].EeSubscriptions = &EeSubscription

	return nil
}

func HandleCreateEeSubscriptions(ueId string, eeSubscription EeSubscription, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateEeSubscriptions")

	locationHeader := CreateEeSubscriptionsProcedure(ueId, eeSubscription)

	headers := http.Header{}
	headers.Set("Location", locationHeader)
	return http_wrapper.NewResponse(http.StatusCreated, headers, eeSubscription)
}

func CreateEeSubscriptionsProcedure(ueId string, EeSubscription EeSubscription) string {
	udrSelf := udr_context.UDR_Self()

	value, ok := udrSelf.UESubsCollection.Load(ueId)
	if !ok {
		udrSelf.UESubsCollection.Store(ueId, new(udr_context.UESubsData))
		value, _ = udrSelf.UESubsCollection.Load(ueId)
	}
	UESubsData := value.(*udr_context.UESubsData)
	if UESubsData.EeSubscriptionCollection == nil {
		UESubsData.EeSubscriptionCollection = make(map[string]*udr_context.EeSubscriptionCollection)
	}

	newSubscriptionID := strconv.Itoa(udrSelf.EeSubscriptionIDGenerator)
	UESubsData.EeSubscriptionCollection[newSubscriptionID] = new(udr_context.EeSubscriptionCollection)
	UESubsData.EeSubscriptionCollection[newSubscriptionID].EeSubscriptions = &EeSubscription
	udrSelf.EeSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/subscription-data/{ueId}/context-data/ee-subscriptions/{subsId} */
	locationHeader := fmt.Sprintf("%s/subscription-data/%s/context-data/ee-subscriptions/%s",
		udrSelf.GetIPv4GroupUri(udr_context.NUDR_DR), ueId, newSubscriptionID)

	return locationHeader
}

func HandleQueryeesubscriptions(ueId string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle Queryeesubscriptions")

	response, problemDetails := QueryeesubscriptionsProcedure(ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func QueryeesubscriptionsProcedure(ueId string) ([]EeSubscription, *ProblemDetails) {
	udrSelf := udr_context.UDR_Self()

	value, ok := udrSelf.UESubsCollection.Load(ueId)
	if !ok {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}

	UESubsData := value.(*udr_context.UESubsData)
	var eeSubscriptionSlice []EeSubscription

	for _, v := range UESubsData.EeSubscriptionCollection {
		eeSubscriptionSlice = append(eeSubscriptionSlice, *v.EeSubscriptions)
	}
	return eeSubscriptionSlice, nil
}

func HandleCreateIndividualExposureDataSubscription(exposureDataSubscription ExposureDataSubscription, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleGetGroupIdentifiers(extGroupId []string, intGroupId []string, ueIdInd []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateHSSSubscriptions(ueId string, subsId string, hssSubscriptionInfo HssSubscriptionInfo, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleGetHssSubscriptionInfo(ueId string, subsId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleModifyHssSubscriptionInfo(ueId string, subsId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleRemoveHssSubscriptionsInfo(ueId string, subsId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateHSSSDMSubscriptions(ueId string, subsId string, hssSubscriptionInfo HssSubscriptionInfo, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleGetHssSDMSubscriptionInfo(ueId string, subsId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleModifyHssSDMSubscriptionInfo(ueId string, subsId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleRemoveHssSDMSubscriptionsInfo(ueId string, subsId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateIpSmGwContext(ueId string, ipSmGwRegistration IpSmGwRegistration, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleDeleteIpSmGwContext(ueId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleModifyIpSmGwContext(ueId string, patchItem []PatchItem, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQueryIpSmGwContext(ueId string, fields []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleReadIPTVCongifurationData(configIds []string, dnns []string, snssais []string, supis []string, interGroupIds []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleDeleteIndividualApplicationDataSubscription(subsId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleReadIndividualApplicationDataSubscription(subsId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleReplaceIndividualApplicationDataSubscription(subsId string, applicationDataSubs ApplicationDataSubs, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateIndividualAppliedBdtPolicyData(bdtPolicyId string, bdtPolicyData BdtPolicyData, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleDeleteIndividualAppliedBdtPolicyData(bdtPolicyId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleUpdateIndividualAppliedBdtPolicyData(bdtPolicyId string, bdtPolicyDataPatch BdtPolicyDataPatch, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleDeleteIndividualAuthenticationStatus(ueId string, servingNetworkName string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQueryIndividualAuthenticationStatus(ueId string, servingNetworkName string, fields []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateIndividualAuthenticationStatus(ueId string, servingNetworkName string, authEvent AuthEvent, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateIndividualBdtData(bdtReferenceId string, bdtData BdtData, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleDeleteIndividualBdtData(bdtReferenceId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleReadIndividualBdtData(bdtReferenceId string, suppFeat []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleUpdateIndividualBdtData(bdtReferenceId string, bdtDataPatch BdtDataPatch, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleDeleteIndividualExposureDataSubscription(subId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleReplaceIndividualExposureDataSubscription(subId string, exposureDataSubscription ExposureDataSubscription, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandlePartialReplaceIndividualIPTVConfigurationData(configurationId string, iptvConfigDataPatch IptvConfigDataPatch, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateOrReplaceIndividualIPTVConfigurationData(configurationId string, iptvConfigData IptvConfigData, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleDeleteIndividualIPTVConfigurationData(configurationId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateOrReplaceIndividualInfluenceData(influenceId string, trafficInfluData TrafficInfluData, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleDeleteIndividualInfluenceData(influenceId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleUpdateIndividualInfluenceData(influenceId string, trafficInfluDataPatch TrafficInfluDataPatch, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleDeleteIndividualInfluenceDataSubscription(subscriptionId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleReadIndividualInfluenceDataSubscription(subscriptionId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleReplaceIndividualInfluenceDataSubscription(subscriptionId string, trafficInfluSub TrafficInfluSub, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateOrReplaceIndividualPFDData(appId string, pfdDataForAppExt PfdDataForAppExt, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleDeleteIndividualPFDData(appId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleReadIndividualPFDData(appId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleDeleteIndividualPolicyDataSubscription(subsId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleReplaceIndividualPolicyDataSubscription(subsId string, policyDataSubscription PolicyDataSubscription, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateOrReplaceServiceParameterData(serviceParamId string, serviceParameterData ServiceParameterData, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleDeleteIndividualServiceParameterData(serviceParamId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleUpdateIndividualServiceParameterData(serviceParamId string, serviceParameterDataPatch ServiceParameterDataPatch, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleReadInfluenceData(influenceIds []string, dnns []string, snssais []string, internalGroupIds []string, supis []string, suppFeat []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateIndividualInfluenceDataSubscription(trafficInfluSub TrafficInfluSub, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleReadInfluenceDataSubscriptions(dnn []string, snssai []string, internalGroupId []string, supi []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQueryLcsBcaData(ueId string, servingPlmnId string, supportedFeatures []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQueryLcsMoData(ueId string, fields []string, supportedFeatures []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQueryLcsPrivacyData(ueId string, fields []string, supportedFeatures []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateMessageWaitingData(ueId string, messageWaitingData MessageWaitingData, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleDeleteMessageWaitingData(ueId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleModifyMessageWaitingData(ueId string, patchItem []PatchItem, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQueryMessageWaitingData(ueId string, fields []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleModify5GVnGroup(externalGroupId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQueryNssaiAck(ueId string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateNSSAIUpdateAck(ueId string, nssaiAckData NssaiAckData, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleModifyOperSpecData(ueId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQueryOperSpecData(ueId string, fields []string, supportedFeatures []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleReadOperatorSpecificData(ueId string, fields []string, suppFeat []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleReplaceOperatorSpecificData(ueId string, requestBody map[string]OperatorSpecificDataContainer, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleUpdateOperatorSpecificData(ueId string, patchItem []PatchItem, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleReadPFDData(appId []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleGetppData(ueId string, supportedFeatures []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle GetppData")

	collName := "subscriptionData.ppData"

	response, problemDetails := GetppDataProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func GetppDataProcedure(collName string, ueId string) (*map[string]interface{}, *ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	ppData := mongodblibrary.RestfulAPIGetOne(collName, filter)

	if ppData != nil {
		return &ppData, nil
	} else {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
}

func HandleQueryPPData(ueId string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQuery5GVNGroupPPData(extGroupIds []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateOrReplaceSessionManagementData(ueId string, pduSessionId string, pduSessionManagementData PduSessionManagementData, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleDeleteSessionManagementData(ueId string, pduSessionId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQuerySessionManagementData(ueId string, pduSessionId string, ipv4Addr []string, ipv6Prefix []string, dnn []string, fields []string, suppFeat []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleReadPlmnUePolicySet(plmnId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateIndividualPolicyDataSubscription(policyDataSubscription PolicyDataSubscription, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQueryProvisionedData(ueId string, servingPlmnId string, datasetNames []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryProvisionedData")

	var provisionedDataSets ProvisionedDataSets

	response, problemDetails := QueryProvisionedDataProcedure(ueId, servingPlmnId, provisionedDataSets)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func QueryProvisionedDataProcedure(ueId string, servingPlmnId string,
	provisionedDataSets ProvisionedDataSets) (*ProvisionedDataSets, *ProblemDetails) {
	{
		collName := "subscriptionData.provisionedData.amData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		accessAndMobilitySubscriptionData := mongodblibrary.RestfulAPIGetOne(collName, filter)
		if accessAndMobilitySubscriptionData != nil {
			var tmp AccessAndMobilitySubscriptionData
			err := mapstructure.Decode(accessAndMobilitySubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.AmData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		smfSelectionSubscriptionData := mongodblibrary.RestfulAPIGetOne(collName, filter)
		if smfSelectionSubscriptionData != nil {
			var tmp SmfSelectionSubscriptionData
			err := mapstructure.Decode(smfSelectionSubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmfSelData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smsData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		smsSubscriptionData := mongodblibrary.RestfulAPIGetOne(collName, filter)
		if smsSubscriptionData != nil {
			var tmp SmsSubscriptionData
			err := mapstructure.Decode(smsSubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmsSubsData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		sessionManagementSubscriptionDatas := mongodblibrary.RestfulAPIGetMany(collName, filter)
		if sessionManagementSubscriptionDatas != nil {
			var tmp []SessionManagementSubscriptionData
			err := mapstructure.Decode(sessionManagementSubscriptionDatas, &tmp)
			if err != nil {
				panic(err)
			}
			for _, smData := range tmp {
				dnnConfigurations := smData.DnnConfigurations
				tmpDnnConfigurations := make(map[string]DnnConfiguration)
				for escapedDnn, dnnConf := range *dnnConfigurations {
					dnn := util.UnescapeDnn(escapedDnn)
					tmpDnnConfigurations[dnn] = dnnConf
				}
				smData.DnnConfigurations = &tmpDnnConfigurations
			}
			provisionedDataSets.SmData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.traceData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		traceData := mongodblibrary.RestfulAPIGetOne(collName, filter)
		if traceData != nil {
			var tmp TraceData
			err := mapstructure.Decode(traceData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.TraceData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smsMngData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		smsManagementSubscriptionData := mongodblibrary.RestfulAPIGetOne(collName, filter)
		if smsManagementSubscriptionData != nil {
			var tmp SmsManagementSubscriptionData
			err := mapstructure.Decode(smsManagementSubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmsMngData = &tmp
		}
	}

	if !reflect.DeepEqual(provisionedDataSets, ProvisionedDataSets{}) {
		return &provisionedDataSets, nil
	} else {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
}

func HandleModifyPpData(ueId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle ModifyPpData")

	collName := "subscriptionData.ppData"

	problemDetails := ModifyPpDataProcedure(collName, ueId, patchItem)
	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
}

func ModifyPpDataProcedure(collName string, ueId string, patchItem []PatchItem) *ProblemDetails {
	filter := bson.M{"ueId": ueId}

	origValue := mongodblibrary.RestfulAPIGetOne(collName, filter)

	patchJSON, err := json.Marshal(patchItem)
	if err != nil {
		logger.DataRepoLog.Errorln(err)
	}

	success := mongodblibrary.RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		newValue := mongodblibrary.RestfulAPIGetOne(collName, filter)
		PreHandleOnDataChangeNotify(ueId, CurrentResourceUri, patchItem, origValue, newValue)
		return nil
	} else {
		return util.ProblemDetailsModifyNotAllowed("")
	}
}

func HandleGet5GVnGroupConfiguration(externalGroupId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleGetAmfSubscriptionInfo(ueId string, subsId string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle GetAmfSubscriptionInfo")

	response, problemDetails := GetAmfSubscriptionInfoProcedure(subsId, ueId)
	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func GetAmfSubscriptionInfoProcedure(subsId string, ueId string) (*[]AmfSubscriptionInfo,
	*ProblemDetails) {
	udrSelf := udr_context.UDR_Self()

	value, ok := udrSelf.UESubsCollection.Load(ueId)
	if !ok {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}

	UESubsData := value.(*udr_context.UESubsData)
	_, ok = UESubsData.EeSubscriptionCollection[subsId]

	if !ok {
		return nil, util.ProblemDetailsNotFound("SUBSCRIPTION_NOT_FOUND")
	}

	if UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos == nil {
		return nil, util.ProblemDetailsNotFound("AMFSUBSCRIPTION_NOT_FOUND")
	}
	return &UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos, nil
}

func HandleGetIdentityData(ueId string, appPortId []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle GetIdentityData")

	collName := "subscriptionData.identityData"

	response, problemDetails := GetIdentityDataProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func GetIdentityDataProcedure(collName string, ueId string) (*map[string]interface{}, *ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	identityData := mongodblibrary.RestfulAPIGetOne(collName, filter)

	if identityData != nil {
		return &identityData, nil
	} else {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
}

func HandleGetNiddAuData(ueId string, singleNssai []string, dnn []string, mtcProviderInformation []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleGetOdbData(ueId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleGetIndividualSharedData(sharedDataId string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {

	logger.DataRepoLog.Infof("HandleGetIndividualSharedData")

	collName := "subscriptionData.sharedData"

	response, problemDetails := GetIndividualSharedDataProcedure(collName, sharedDataId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}
func GetIndividualSharedDataProcedure(collName string, sharedDataId string) (*map[string]interface{},
	*ProblemDetails) {

	filter := bson.M{"sharedDataId": sharedDataId}
	sharedData := mongodblibrary.RestfulAPIGetOne(collName, filter)

	delete(sharedData, "_id")

	if sharedData != nil {
		return &sharedData, nil
	} else {
		return nil, util.ProblemDetailsNotFound("DATA_NOT_FOUND")
	}
}

func HandleGetSharedData(sharedDataIds []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {

	logger.DataRepoLog.Infof("Handle GetSharedData")

	collName := "subscriptionData.sharedData"

	response, problemDetails := GetSharedDataProcedure(collName, sharedDataIds)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)

}

func GetSharedDataProcedure(collName string, sharedDataIds []string) (*[]map[string]interface{},
	*ProblemDetails) {
	var sharedDataArray []map[string]interface{}
	for _, sharedDataId := range sharedDataIds {
		filter := bson.M{"sharedDataId": sharedDataId}
		sharedData := mongodblibrary.RestfulAPIGetOne(collName, filter)
		if sharedData != nil {
			sharedDataArray = append(sharedDataArray, sharedData)
		}
	}

	if sharedDataArray != nil {
		return &sharedDataArray, nil
	} else {
		return nil, util.ProblemDetailsNotFound("DATA_NOT_FOUND")
	}
}

func HandleModifysdmSubscription(ueId string, subsId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQuerysdmSubscription(ueId string, subsId string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle Querysdmsubscriptions")

	response, problemDetails := QuerysdmsubscriptionsProcedure(ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func QuerysdmsubscriptionsProcedure(ueId string) (*[]SdmSubscription, *ProblemDetails) {
	udrSelf := udr_context.UDR_Self()

	value, ok := udrSelf.UESubsCollection.Load(ueId)
	if !ok {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}

	UESubsData := value.(*udr_context.UESubsData)
	var sdmSubscriptionSlice []SdmSubscription

	for _, v := range UESubsData.SdmSubscriptions {
		sdmSubscriptionSlice = append(sdmSubscriptionSlice, *v)
	}
	return &sdmSubscriptionSlice, nil
}

func HandleRemovesdmSubscriptions(ueId string, subsId string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle RemovesdmSubscriptions")

	problemDetails := RemovesdmSubscriptionsProcedure(ueId, subsId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
}

func RemovesdmSubscriptionsProcedure(ueId string, subsId string) *ProblemDetails {
	udrSelf := udr_context.UDR_Self()
	value, ok := udrSelf.UESubsCollection.Load(ueId)
	if !ok {
		return util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}

	UESubsData := value.(*udr_context.UESubsData)
	_, ok = UESubsData.SdmSubscriptions[subsId]

	if !ok {
		return util.ProblemDetailsNotFound("SUBSCRIPTION_NOT_FOUND")
	}
	delete(UESubsData.SdmSubscriptions, subsId)

	return nil
}
func HandleUpdatesdmsubscriptions(ueId string, subsId string, sdmSubscription SdmSubscription, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle Updatesdmsubscriptions")

	problemDetails := UpdatesdmsubscriptionsProcedure(ueId, subsId, sdmSubscription)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
}

func UpdatesdmsubscriptionsProcedure(ueId string, subsId string,
	SdmSubscription SdmSubscription) *ProblemDetails {
	udrSelf := udr_context.UDR_Self()
	value, ok := udrSelf.UESubsCollection.Load(ueId)
	if !ok {
		return util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}

	UESubsData := value.(*udr_context.UESubsData)
	_, ok = UESubsData.SdmSubscriptions[subsId]

	if !ok {
		return util.ProblemDetailsNotFound("SUBSCRIPTION_NOT_FOUND")
	}
	SdmSubscription.SubscriptionId = &subsId
	UESubsData.SdmSubscriptions[subsId] = &SdmSubscription

	return nil
}

func HandleCreateSdmSubscriptions(ueId string, sdmSubscription SdmSubscription, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateSdmSubscriptions")

	collName := "subscriptionData.contextData.amfNon3gppAccess"

	locationHeader, SdmSubscription := CreateSdmSubscriptionsProcedure(sdmSubscription, collName, ueId)

	headers := http.Header{}
	headers.Set("Location", locationHeader)
	return http_wrapper.NewResponse(http.StatusCreated, headers, SdmSubscription)
}

func CreateSdmSubscriptionsProcedure(sdmSubscription SdmSubscription,
	collName string, ueId string) (string, SdmSubscription) {
	udrSelf := udr_context.UDR_Self()

	value, ok := udrSelf.UESubsCollection.Load(ueId)
	if !ok {
		udrSelf.UESubsCollection.Store(ueId, new(udr_context.UESubsData))
		value, _ = udrSelf.UESubsCollection.Load(ueId)
	}
	UESubsData := value.(*udr_context.UESubsData)
	if UESubsData.SdmSubscriptions == nil {
		UESubsData.SdmSubscriptions = make(map[string]*SdmSubscription)
	}

	newSubscriptionID := strconv.Itoa(udrSelf.SdmSubscriptionIDGenerator)
	sdmSubscription.SubscriptionId = &newSubscriptionID
	UESubsData.SdmSubscriptions[newSubscriptionID] = &sdmSubscription
	udrSelf.SdmSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/subscription-data/{ueId}/context-data/sdm-subscriptions/{subsId}' */
	locationHeader := fmt.Sprintf("%s/subscription-data/%s/context-data/sdm-subscriptions/%s",
		udrSelf.GetIPv4GroupUri(udr_context.NUDR_DR), ueId, newSubscriptionID)

	return locationHeader, sdmSubscription
}
func HandleQuerysdmsubscriptions(ueId string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle Querysdmsubscriptions")

	response, problemDetails := QuerysdmsubscriptionsProcedure(ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func HandleCreateSMFSubscriptions(ueId string, subsId string, smfSubscriptionInfo SmfSubscriptionInfo, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleGetSmfSubscriptionInfo(ueId string, subsId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleModifySmfSubscriptionInfo(ueId string, subsId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleRemoveSmfSubscriptionsInfo(ueId string, subsId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateSmfContextNon3gpp(ueId string, pduSessionIdIn string, smfRegistration SmfRegistration, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateSmfContextNon3gpp")

	collName := "subscriptionData.contextData.smfRegistrations"

	pduSessionId, err := strconv.ParseInt(pduSessionIdIn, 10, 64)
	if err != nil {
		logger.DataRepoLog.Warnln(err)
	}

	response, status := CreateSmfContextNon3gppProcedure(smfRegistration, collName, ueId, pduSessionId)

	if status == http.StatusCreated {
		return http_wrapper.NewResponse(http.StatusCreated, nil, response)
	} else if status == http.StatusOK {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func CreateSmfContextNon3gppProcedure(SmfRegistration SmfRegistration,
	collName string, ueId string, pduSessionIdInt int64) (bson.M, int) {
	putData := util.ToBsonM(SmfRegistration)
	putData["ueId"] = ueId
	putData["pduSessionId"] = int32(pduSessionIdInt)

	filter := bson.M{"ueId": ueId, "pduSessionId": pduSessionIdInt}
	isExisted := mongodblibrary.RestfulAPIPutOne(collName, filter, putData)

	if !isExisted {
		return putData, http.StatusCreated
	} else {
		return putData, http.StatusOK
	}
}
func HandleDeleteSmfContext(ueId string, pduSessionId string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle DeleteSmfContext")

	collName := "subscriptionData.contextData.smfRegistrations"

	DeleteSmfContextProcedure(collName, ueId, pduSessionId)
	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func DeleteSmfContextProcedure(collName string, ueId string, pduSessionId string) {
	pduSessionIdInt, err := strconv.ParseInt(pduSessionId, 10, 32)
	if err != nil {
		logger.DataRepoLog.Error(err)
	}
	filter := bson.M{"ueId": ueId, "pduSessionId": pduSessionIdInt}

	mongodblibrary.RestfulAPIDeleteOne(collName, filter)
}

func HandleQuerySmfRegistration(ueId string, pduSessionId string, fields []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QuerySmfRegistration")

	collName := "subscriptionData.contextData.smfRegistrations"

	response, problemDetails := QuerySmfRegistrationProcedure(collName, ueId, pduSessionId)
	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}

	pd := util.ProblemDetailsUpspecified("")
	return http_wrapper.NewResponse(int(*pd.Status), nil, pd)
}

func QuerySmfRegistrationProcedure(collName string, ueId string,
	pduSessionId string) (*map[string]interface{}, *ProblemDetails) {
	pduSessionIdInt, err := strconv.ParseInt(pduSessionId, 10, 32)
	if err != nil {
		logger.DataRepoLog.Error(err)
	}

	filter := bson.M{"ueId": ueId, "pduSessionId": pduSessionIdInt}

	smfRegistration := mongodblibrary.RestfulAPIGetOne(collName, filter)

	if smfRegistration != nil {
		return &smfRegistration, nil
	} else {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
}

func HandleQuerySmfRegList(ueId string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QuerySmfRegList")

	collName := "subscriptionData.contextData.smfRegistrations"

	response := QuerySmfRegListProcedure(collName, ueId)

	if response == nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, []map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	}
}

func QuerySmfRegListProcedure(collName string, ueId string) *[]map[string]interface{} {
	filter := bson.M{"ueId": ueId}
	smfRegList := mongodblibrary.RestfulAPIGetMany(collName, filter)

	if smfRegList != nil {
		return &smfRegList
	} else {
		// Return empty array instead
		return nil
	}
}

func HandleQuerySmfSelectData(ueId string, servingPlmnId string, fields []string, supportedFeatures []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QuerySmfSelectData")

	collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"

	response, problemDetails := QuerySmfSelectDataProcedure(collName, ueId, servingPlmnId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return http_wrapper.NewResponse(int(*problemDetails.Status), nil, problemDetails)
	}
}

func QuerySmfSelectDataProcedure(collName string, ueId string,
	servingPlmnId string) (*map[string]interface{}, *ProblemDetails) {
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	smfSelectionSubscriptionData := mongodblibrary.RestfulAPIGetOne(collName, filter)

	if smfSelectionSubscriptionData != nil {
		return &smfSelectionSubscriptionData, nil
	} else {
		return nil, util.ProblemDetailsNotFound("USER_NOT_FOUND")
	}
}

func HandleCreateSmsfContext3gpp(ueId string, smsfRegistration SmsfRegistration, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleDeleteSmsfContext3gpp(ueId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQuerySmsfContext3gpp(ueId string, fields []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateSmsfContextNon3gpp(ueId string, smsfRegistration SmsfRegistration, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleDeleteSmsfContextNon3gpp(ueId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQuerySmsfContextNon3gpp(ueId string, fields []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQuerySmsMngData(ueId string, servingPlmnId string, supportedFeatures []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQuerySmsData(ueId string, servingPlmnId string, supportedFeatures []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleReadServiceParameterData(serviceParamIds []string, dnns []string, snssais []string, internalGroupIds []string, supis []string, ueIpv4s []string, ueIpv6s []string, ueMacs []string, suppFeat []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleReadSessionManagementPolicyData(ueId string, snssai []string, dnn []string, fields []string, suppFeat []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleUpdateSessionManagementPolicyData(ueId string, smPolicyDataPatch SmPolicyDataPatch, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQuerySmData(ueId string, servingPlmnId string, singleNssai []string, dnn []string, fields []string, supportedFeatures []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleReadSponsorConnectivityData(sponsorId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQuerySubsToNotify(ueId []string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleRemoveMultipleSubscriptionDataSubscriptions(ueId []string, nfInstanceId []string, deleteAllNfs []string, implicitUnsubscribeIndication []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleSubscriptionDataSubscriptions(subscriptionDataSubscriptions SubscriptionDataSubscriptions, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleModifysubscriptionDataSubscription(subsId string, patchItem []PatchItem, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleQuerySubscriptionDataSubscriptions(subsId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleRemovesubscriptionDataSubscriptions(subsId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQueryTraceData(ueId string, servingPlmnId string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateOrReplaceUEPolicySet(ueId string, uePolicySet UePolicySet, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleReadUEPolicySet(ueId string, suppFeat []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleUpdateUEPolicySet(ueId string, uePolicySetPatch UePolicySetPatch, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQueryUeLocation(ueId string, supportedFeatures []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleCreateUsageMonitoringResource(ueId string, usageMonId string, usageMonData UsageMonData, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleDeleteUsageMonitoringInformation(ueId string, usageMonId string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
func HandleReadUsageMonitoringInformation(ueId string, usageMonId string, suppFeat []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}

func HandleQueryV2xData(ueId string, supportedFeatures []string, ifNoneMatch []string, ifModifiedSince []string, c *gin.Context) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, map[string]interface{}{})
}
