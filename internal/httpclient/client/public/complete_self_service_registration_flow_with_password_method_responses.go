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

// CompleteSelfServiceRegistrationFlowWithPasswordMethodReader is a Reader for the CompleteSelfServiceRegistrationFlowWithPasswordMethod structure.
type CompleteSelfServiceRegistrationFlowWithPasswordMethodReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CompleteSelfServiceRegistrationFlowWithPasswordMethodReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCompleteSelfServiceRegistrationFlowWithPasswordMethodOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 302:
		result := NewCompleteSelfServiceRegistrationFlowWithPasswordMethodFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 400:
		result := NewCompleteSelfServiceRegistrationFlowWithPasswordMethodBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewCompleteSelfServiceRegistrationFlowWithPasswordMethodInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewCompleteSelfServiceRegistrationFlowWithPasswordMethodOK creates a CompleteSelfServiceRegistrationFlowWithPasswordMethodOK with default headers values
func NewCompleteSelfServiceRegistrationFlowWithPasswordMethodOK() *CompleteSelfServiceRegistrationFlowWithPasswordMethodOK {
	return &CompleteSelfServiceRegistrationFlowWithPasswordMethodOK{}
}

/*CompleteSelfServiceRegistrationFlowWithPasswordMethodOK handles this case with default header values.

registrationViaApiResponse
*/
type CompleteSelfServiceRegistrationFlowWithPasswordMethodOK struct {
	Payload *models.RegistrationViaAPIResponse
}

func (o *CompleteSelfServiceRegistrationFlowWithPasswordMethodOK) Error() string {
	return fmt.Sprintf("[POST /self-service/registration/methods/password][%d] completeSelfServiceRegistrationFlowWithPasswordMethodOK  %+v", 200, o.Payload)
}

func (o *CompleteSelfServiceRegistrationFlowWithPasswordMethodOK) GetPayload() *models.RegistrationViaAPIResponse {
	return o.Payload
}

func (o *CompleteSelfServiceRegistrationFlowWithPasswordMethodOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RegistrationViaAPIResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCompleteSelfServiceRegistrationFlowWithPasswordMethodFound creates a CompleteSelfServiceRegistrationFlowWithPasswordMethodFound with default headers values
func NewCompleteSelfServiceRegistrationFlowWithPasswordMethodFound() *CompleteSelfServiceRegistrationFlowWithPasswordMethodFound {
	return &CompleteSelfServiceRegistrationFlowWithPasswordMethodFound{}
}

/*CompleteSelfServiceRegistrationFlowWithPasswordMethodFound handles this case with default header values.

Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
typically 201.
*/
type CompleteSelfServiceRegistrationFlowWithPasswordMethodFound struct {
}

func (o *CompleteSelfServiceRegistrationFlowWithPasswordMethodFound) Error() string {
	return fmt.Sprintf("[POST /self-service/registration/methods/password][%d] completeSelfServiceRegistrationFlowWithPasswordMethodFound ", 302)
}

func (o *CompleteSelfServiceRegistrationFlowWithPasswordMethodFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCompleteSelfServiceRegistrationFlowWithPasswordMethodBadRequest creates a CompleteSelfServiceRegistrationFlowWithPasswordMethodBadRequest with default headers values
func NewCompleteSelfServiceRegistrationFlowWithPasswordMethodBadRequest() *CompleteSelfServiceRegistrationFlowWithPasswordMethodBadRequest {
	return &CompleteSelfServiceRegistrationFlowWithPasswordMethodBadRequest{}
}

/*CompleteSelfServiceRegistrationFlowWithPasswordMethodBadRequest handles this case with default header values.

registrationFlow
*/
type CompleteSelfServiceRegistrationFlowWithPasswordMethodBadRequest struct {
	Payload *models.RegistrationFlow
}

func (o *CompleteSelfServiceRegistrationFlowWithPasswordMethodBadRequest) Error() string {
	return fmt.Sprintf("[POST /self-service/registration/methods/password][%d] completeSelfServiceRegistrationFlowWithPasswordMethodBadRequest  %+v", 400, o.Payload)
}

func (o *CompleteSelfServiceRegistrationFlowWithPasswordMethodBadRequest) GetPayload() *models.RegistrationFlow {
	return o.Payload
}

func (o *CompleteSelfServiceRegistrationFlowWithPasswordMethodBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RegistrationFlow)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCompleteSelfServiceRegistrationFlowWithPasswordMethodInternalServerError creates a CompleteSelfServiceRegistrationFlowWithPasswordMethodInternalServerError with default headers values
func NewCompleteSelfServiceRegistrationFlowWithPasswordMethodInternalServerError() *CompleteSelfServiceRegistrationFlowWithPasswordMethodInternalServerError {
	return &CompleteSelfServiceRegistrationFlowWithPasswordMethodInternalServerError{}
}

/*CompleteSelfServiceRegistrationFlowWithPasswordMethodInternalServerError handles this case with default header values.

genericError
*/
type CompleteSelfServiceRegistrationFlowWithPasswordMethodInternalServerError struct {
	Payload *models.GenericError
}

func (o *CompleteSelfServiceRegistrationFlowWithPasswordMethodInternalServerError) Error() string {
	return fmt.Sprintf("[POST /self-service/registration/methods/password][%d] completeSelfServiceRegistrationFlowWithPasswordMethodInternalServerError  %+v", 500, o.Payload)
}

func (o *CompleteSelfServiceRegistrationFlowWithPasswordMethodInternalServerError) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *CompleteSelfServiceRegistrationFlowWithPasswordMethodInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
