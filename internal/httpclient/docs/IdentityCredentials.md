# IdentityCredentials

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**Identifiers** | Pointer to **[]string** | Identifiers represents a list of unique identifiers this credential type matches. | [optional] 
**Type** | Pointer to **string** | and so on. | [optional] 

## Methods

### NewIdentityCredentials

`func NewIdentityCredentials() *IdentityCredentials`

NewIdentityCredentials instantiates a new IdentityCredentials object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIdentityCredentialsWithDefaults

`func NewIdentityCredentialsWithDefaults() *IdentityCredentials`

NewIdentityCredentialsWithDefaults instantiates a new IdentityCredentials object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConfig

`func (o *IdentityCredentials) GetConfig() map[string]interface{}`

GetConfig returns the Config field if non-nil, zero value otherwise.

### GetConfigOk

`func (o *IdentityCredentials) GetConfigOk() (*map[string]interface{}, bool)`

GetConfigOk returns a tuple with the Config field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfig

`func (o *IdentityCredentials) SetConfig(v map[string]interface{})`

SetConfig sets Config field to given value.

### HasConfig

`func (o *IdentityCredentials) HasConfig() bool`

HasConfig returns a boolean if a field has been set.

### GetIdentifiers

`func (o *IdentityCredentials) GetIdentifiers() []string`

GetIdentifiers returns the Identifiers field if non-nil, zero value otherwise.

### GetIdentifiersOk

`func (o *IdentityCredentials) GetIdentifiersOk() (*[]string, bool)`

GetIdentifiersOk returns a tuple with the Identifiers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdentifiers

`func (o *IdentityCredentials) SetIdentifiers(v []string)`

SetIdentifiers sets Identifiers field to given value.

### HasIdentifiers

`func (o *IdentityCredentials) HasIdentifiers() bool`

HasIdentifiers returns a boolean if a field has been set.

### GetType

`func (o *IdentityCredentials) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *IdentityCredentials) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *IdentityCredentials) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *IdentityCredentials) HasType() bool`

HasType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


