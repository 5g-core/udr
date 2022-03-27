package util

import (
	"net/http"

	"github.com/5g-core/openapi/models"
)

func ProblemDetailsSystemFailure(detail string) *models.ProblemDetails {
	title := "System failure"
	var status int32
	status = int32(http.StatusInternalServerError)
	cause := "SYSTEM_FAILURE"
	return &models.ProblemDetails{
		Title:  &title,
		Status: &status,
		Detail: &detail,
		Cause:  &cause,
	}
}

func ProblemDetailsMalformedReqSyntax(detail string) *models.ProblemDetails {
	title := "Malformed request syntax"
	var status int32
	status = int32(http.StatusBadRequest)
	return &models.ProblemDetails{
		Title:  &title,
		Status: &status,
		Detail: &detail,
	}
}

func ProblemDetailsNotFound(cause string) *models.ProblemDetails {
	title := ""
	var status int32
	status = int32(http.StatusNotFound)
	if cause == "USER_NOT_FOUND" {
		title = "User not found"
	} else if cause == "SUBSCRIPTION_NOT_FOUND" {
		title = "Subscription not found"
	} else if cause == "AMFSUBSCRIPTION_NOT_FOUND" {
		title = "AMF Subscription not found"
	} else {
		title = "Data not found"
	}
	return &models.ProblemDetails{
		Title:  &title,
		Status: &status,
		Cause:  &cause,
	}
}

func ProblemDetailsModifyNotAllowed(detail string) *models.ProblemDetails {
	var status int32
	status = int32(http.StatusForbidden)
	title := "Modify not allowed"
	cause := "MODIFY_NOT_ALLOWED"
	return &models.ProblemDetails{
		Title:  &title,
		Status: &status,
		Cause:  &cause,
		Detail: &detail,
	}
}

func ProblemDetailsUpspecified(detail string) *models.ProblemDetails {
	var status int32
	status = int32(http.StatusForbidden)
	title := "Unspecified"
	cause := "UNSPECIFIED"
	return &models.ProblemDetails{
		Title:  &title,
		Status: &status,
		Cause:  &cause,
		Detail: &detail,
	}
}
