# VMA APIs
VMA Swagger APIs.

## Version: 1.0

### Terms of service
http://swagger.io/terms/

**Contact information:**  
API Support  
http://www.swagger.io/support  
rafael.xavier.lima@gmail.com  

### Security
**JWT**  

|apiKey|*API Key*|
|---|---|
|In|header|
|Name|token|

### /secured/rating

#### GET
##### Summary:

Endpoint to get the rating list in use

##### Description:

Get the rating list

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [helper.Ratings](#helper.Ratings) ] |

### /secured/rating/byage/{age}

#### GET
##### Summary:

Endpoint to get the rating

##### Description:

Get an rating classification depending of the required age

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| age | path | An Age | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [helper.Ratings](#helper.Ratings) |
| 204 | No Content | string |
| 400 | Bad request | string |

### /secured/user

#### GET
##### Summary:

Endpoint to load a list of users limited to 100

##### Description:

Get a list of users

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [models.User](#models.User) ] |
| 400 | Bad request | string |

### /secured/user/{email}

#### GET
##### Summary:

Endpoint to load an user by it's email

##### Description:

Get an user

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| email | path | User Email | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.User](#models.User) |
| 400 | Bad request | string |

### /secured/user/{id}

#### DELETE
##### Summary:

Endpoint to exclude an user

##### Description:

Delete an user

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | User Id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK |  |
| 204 | No Content | string |

### /secured/user/update

#### PUT
##### Summary:

Endpoint to update common user information

##### Description:

Update user information

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| age | formData |  | No | integer |
| createdAt | formData |  | No | string |
| email | formData |  | Yes | string |
| firstName | formData |  | Yes | string |
| id | formData |  | No | string |
| lastName | formData |  | No | string |
| updatedAt | formData |  | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.User](#models.User) |
| 204 | No Content | string |

### /secured/user/update/email

#### PUT
##### Summary:

Endpoint to register a new user

##### Description:

Add a new user

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| email | query | User email | Yes | string |
| newemail | query | User newemail | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.User](#models.User) |
| 204 | No Content | string |
| 400 | Bad request | string |

### /secured/user/update/password

#### PUT
##### Summary:

Endpoint to update the user password

##### Description:

Update the user password

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| email | query | User email | Yes | string |
| password | query | User password | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK |  |
| 204 | No Content | string |
| 400 | Bad Request | string |

### /user/register

#### DELETE
##### Summary:

Endpoint to register a new user

##### Description:

Add a new user

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| email | path | User Data | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK |  |
| 204 | No Content | string |

#### POST
##### Summary:

Endpoint to register a new user

##### Description:

Add a new user

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| age | formData |  | No | integer |
| createdAt | formData |  | No | string |
| email | formData |  | Yes | string |
| firstName | formData |  | Yes | string |
| id | formData |  | No | string |
| lastName | formData |  | No | string |
| updatedAt | formData |  | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.User](#models.User) |
| 204 | No Content | string |
| 400 | Bad request | string |
| 406 | Not Acceptable | string |
| 500 | Server Error | string |

### Models


#### helper.Ratings

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string |  | No |
| minAge | integer |  | No |
| rating | string |  | No |

#### models.Address

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| city | string |  | No |
| country | string |  | No |
| state | string |  | No |
| street | string |  | No |

#### models.Password

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| createdAt | string |  | No |
| encryptedPass | string |  | No |
| expire | string |  | No |

#### models.User

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | [models.Address](#models.Address) |  | No |
| age | integer |  | No |
| createdAt | string |  | No |
| email | string |  | Yes |
| firstName | string |  | Yes |
| id | string |  | No |
| lastName | string |  | No |
| password | [models.Password](#models.Password) |  | No |
| updatedAt | string |  | No |