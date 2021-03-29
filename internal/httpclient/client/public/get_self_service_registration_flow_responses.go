// Code generated by go-swagger; DO NOT EDIT.

package public

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"kratos/internal/httpclient/models"
)

// GetSelfServiceRegistrationFlowReader is a Reader for the GetSelfServiceRegistrationFlow structure.
type GetSelfServiceRegistrationFlowReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetSelfServiceRegistrationFlowReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetSelfServiceRegistrationFlowOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 403:
		result := NewGetSelfServiceRegistrationFlowForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetSelfServiceRegistrationFlowNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 410:
		result := NewGetSelfServiceRegistrationFlowGone()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetSelfServiceRegistrationFlowInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetSelfServiceRegistrationFlowOK creates a GetSelfServiceRegistrationFlowOK with default headers values
func NewGetSelfServiceRegistrationFlowOK() *GetSelfServiceRegistrationFlowOK {
	return &GetSelfServiceRegistrationFlowOK{}
}

/*GetSelfServiceRegistrationFlowOK handles this case with default header values.

registrationFlow
*/
type GetSelfServiceRegistrationFlowOK struct {
	Payload *models.RegistrationFlow
}

func (o *GetSelfServiceRegistrationFlowOK) Error() string {
	return fmt.Sprintf("[GET /self-service/registration/flows][%d] getSelfServiceRegistrationFlowOK  %+v", 200, o.Payload)
}

func (o *GetSelfServiceRegistrationFlowOK) GetPayload() *models.RegistrationFlow {
	return o.Payload
}

func (o *GetSelfServiceRegistrationFlowOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RegistrationFlow)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetSelfServiceRegistrationFlowForbidden creates a GetSelfServiceRegistrationFlowForbidden with default headers values
func NewGetSelfServiceRegistrationFlowForbidden() *GetSelfServiceRegistrationFlowForbidden {
	return &GetSelfServiceRegistrationFlowForbidden{}
}

/*GetSelfServiceRegistrationFlowForbidden handles this case with default header values.

genericError
*/
type GetSelfServiceRegistrationFlowForbidden struct {
	Payload *models.GenericError
}

func (o *GetSelfServiceRegistrationFlowForbidden) Error() string {
	return fmt.Sprintf("[GET /self-service/registration/flows][%d] getSelfServiceRegistrationFlowForbidden  %+v", 403, o.Payload)
}

func (o *GetSelfServiceRegistrationFlowForbidden) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *GetSelfServiceRegistrationFlowForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetSelfServiceRegistrationFlowNotFound creates a GetSelfServiceRegistrationFlowNotFound with default headers values
func NewGetSelfServiceRegistrationFlowNotFound() *GetSelfServiceRegistrationFlowNotFound {
	return &GetSelfServiceRegistrationFlowNotFound{}
}

/*GetSelfServiceRegistrationFlowNotFound handles this case with default header values.

genericError
*/
type GetSelfServiceRegistrationFlowNotFound struct {
	Payload *models.GenericError
}

func (o *GetSelfServiceRegistrationFlowNotFound) Error() string {
	return fmt.Sprintf("[GET /self-service/registration/flows][%d] getSelfServiceRegistrationFlowNotFound  %+v", 404, o.Payload)
}

func (o *GetSelfServiceRegistrationFlowNotFound) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *GetSelfServiceRegistrationFlowNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetSelfServiceRegistrationFlowGone creates a GetSelfServiceRegistrationFlowGone with default headers values
func NewGetSelfServiceRegistrationFlowGone() *GetSelfServiceRegistrationFlowGone {
	return &GetSelfServiceRegistrationFlowGone{}
}

/*GetSelfServiceRegistrationFlowGone handles this case with default header values.

genericError
*/
type GetSelfServiceRegistrationFlowGone struct {
	Payload *models.GenericError
}

func (o *GetSelfServiceRegistrationFlowGone) Error() string {
	return fmt.Sprintf("[GET /self-service/registration/flows][%d] getSelfServiceRegistrationFlowGone  %+v", 410, o.Payload)
}

func (o *GetSelfServiceRegistrationFlowGone) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *GetSelfServiceRegistrationFlowGone) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetSelfServiceRegistrationFlowInternalServerError creates a GetSelfServiceRegistrationFlowInternalServerError with default headers values
func NewGetSelfServiceRegistrationFlowInternalServerError() *GetSelfServiceRegistrationFlowInternalServerError {
	return &GetSelfServiceRegistrationFlowInternalServerError{}
}

/*GetSelfServiceRegistrationFlowInternalServerError handles this case with default header values.

genericError
*/
type GetSelfServiceRegistrationFlowInternalServerError struct {
	Payload *models.GenericError
}

func (o *GetSelfServiceRegistrationFlowInternalServerError) Error() string {
	return fmt.Sprintf("[GET /self-service/registration/flows][%d] getSelfServiceRegistrationFlowInternalServerError  %+v", 500, o.Payload)
}

func (o *GetSelfServiceRegistrationFlowInternalServerError) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *GetSelfServiceRegistrationFlowInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
