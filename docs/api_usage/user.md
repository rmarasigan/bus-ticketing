# User
The User API Schema contains the user account information. In this module, it will let you:
* [Create an account](#create-an-account)
* [Log in to the user account](#login)
* [Get user account information](#get-user-account)
* [Update user account properties](#update-user-account)

## Data Structure
<table>
  <tr>
    <th>Field</th>
    <th>Type</th>
    <th>Description</th>
  </tr>
  <tr>
    <td>
      <code>id</code>
    </td>
    <td>string</td>
    <td>The unique user ID and the sort key.</td>
  </tr>
  <tr>
    <td>
      <code>user_type</code>
    </td>
    <td>string</td>
    <td>
      The type of this account.<br />
      There are 2 different user types: <br />
      1 = "ADMIN" <br />
      2 = "CUSTOMER"
    </td>
  </tr>
  <tr>
    <td>
      <code>first_name</code>
    </td>
    <td>string</td>
    <td>
      The first name of the user.
    </td>
  </tr>
  <tr>
    <td>
      <code>last_name</code>
    </td>
    <td>string</td>
    <td>The last name of the user.</td>
  </tr>
  <tr>
    <td>
      <code>username</code>
    </td>
    <td>string</td>
    <td>The username of the user account and the primary key.</td>
  </tr>
  <tr>
    <td>
      <code>password</code>
    </td>
    <td>string</td>
    <td>
      The user security password for the account.
    </td>
  </tr>
  <tr>
    <td>
      <code>address</code>
    </td>
    <td>string</td>
    <td>The user address.</td>
  </tr>
  <tr>
    <td>
      <code>email</code>
    </td>
    <td>string</td>
    <td>The user e-mail address.</td>
  </tr>
  <tr>
    <td>
      <code>mobile_number</code>
    </td>
    <td>string</td>
    <td>The user phone number.</td>
  </tr>
  <tr>
    <td>
      <code>date_created</code>
    </td>
    <td>string</td>
    <td>The date that this account was created.</td>
  </tr>
  <tr>
    <td>
      <code>last_login</code>
    </td>
    <td>string</td>
    <td>The last login session of the user.</td>
  </tr>
</table>

## API Usage and Specification
#### Headers
<table>
  <tr>
    <th>Key</th>
    <th>Value</th>
  </tr>
  <tr>
    <td>
      <code>Content-Type</code>
    </td>
    <td>
      <code>application/json</code>
    </td>
  </tr>
</table>

Setting to `application/json` is recommended.

#### HTTP Response Status Codes
<table>
  <tr>
    <th>Status Code</th>
    <th>Description</th>
  </tr>
  <tr>
    <td>200</td>
    <td>OK</td>
  </tr>
  <tr>
    <td>400</td>
    <td>Bad Request</td>
  </tr>
  <tr>
    <td>500</td>
    <td>Internal Server Error</td>
  </tr>
</table>

### Create an account
To create a new user account instance, you need to instantiate an object that represents a user account property. The user account instance holds the information related to the user.

**Method**: `POST`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/user/create

#### Payload
<table>
  <tr>
    <th>Field</th>
    <th>Type</th>
    <th>Description</th>
    <th>Required</th>
  </tr>
  <tr>
    <td>
      <code>user_type</code>
    </td>
    <td>string</td>
    <td>
      The type of this account.<br />
      There are 2 different user types: <br />
      1 = "ADMIN" <br />
      2 = "CUSTOMER"
    </td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>first_name</code>
    </td>
    <td>string</td>
    <td>
      The first name of the user.
    </td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>last_name</code>
    </td>
    <td>string</td>
    <td>The last name of the user.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>username</code>
    </td>
    <td>string</td>
    <td>The username of the user account.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>password</code>
    </td>
    <td>string</td>
    <td>
      Minimum length: 8 <br/>
      The user security password for the account.
    </td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>address</code>
    </td>
    <td>string</td>
    <td>The user address.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>email</code>
    </td>
    <td>string</td>
    <td>The user e-mail address.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>mobile_number</code>
    </td>
    <td>string</td>
    <td>The user phone number.</td>
    <td>✅</td>
  </tr>
</table>

#### Sample Payload
```json
{
  "user_type": "1",
  "first_name": "Emily",
  "last_name": "Davis",
  "username": "emilydavis",
  "password": "passwordabc",
  "address": "321 Cedar Road",
  "email": "emilydavis@example.com",
  "mobile_number": "4449876543"
}
```

### Login
By providing the `username` and `password` in the payload, the authentication process will be initiated, allowing the user to access their account. Upon successful login, it will return a representation of an account that is related to the user as a response. 

**Method**: `POST`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/user/login

#### Payload
<table>
  <tr>
    <th>Field</th>
    <th>Type</th>
    <th>Description</th>
    <th>Required</th>
  </tr>
  <tr>
    <td>
      <code>username</code>
    </td>
    <td>string</td>
    <td>The username of the user account.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>password</code>
    </td>
    <td>string</td>
    <td>
      Minimum length: 8 <br/>
      The user security password for the account.
    </td>
    <td>✅</td>
  </tr>
</table>

#### Sample Request
Payload:
```json
{
  "username": "emilydavis",
  "password": "passwordabc"
}
```

Response:
```json
{
  "id": "ADMN-878495",
  "user_type": "ADMIN",
  "first_name": "Emily",
  "last_name": "Davis",
  "username": "emilydavis",
  "address": "321 Cedar Road",
  "email": "emilydavis@example.com",
  "mobile_number": "4449876543"
}
```

### Get user account
When retrieving the specific user account information, the `id` and `username` query parameters must be present in the URL. These parameters identify which user account should be returned. It will either return a representation of a specific an account that is related to the user or a list of user account.

**Method**: `GET`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/user/account/get

#### Specific User Account
**Query Parameters**
<table>
  <tr>
    <th>Parameter</th>
    <th>Type</th>
    <th>Description</th>
    <th>Required</th>
  </tr>
  <tr>
    <td>
      <code>id</code>
    </td>
    <td>string</td>
    <td>The user account unique ID.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>username</code>
    </td>
    <td>string</td>
    <td>The username of the user account.</td>
    <td>✅</td>
  </tr>
</table>

#### Sample Response
```json
[
  {
    "id": "ADMN-878495",
    "user_type": "ADMIN",
    "first_name": "Emily",
    "last_name": "Davis",
    "username": "emilydavis",
    "address": "321 Cedar Road",
    "email": "emilydavis@example.com",
    "mobile_number": "(407) 435-6841"
  }
]
```

### Update user account
When modifying the user profile, the `id` and `username` query parameters must be present in the URL. These parameters identify which user account should be modified. After the update is performed, it will return a representation of the updated account information.

**Method**: `POST`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/user/account/update?id=xxxxx&username=xxxxx

#### Query Parameters
<table>
  <tr>
    <th>Parameter</th>
    <th>Type</th>
    <th>Description</th>
    <th>Required</th>
  </tr>
  <tr>
    <td>
      <code>id</code>
    </td>
    <td>string</td>
    <td>The user account unique ID.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>username</code>
    </td>
    <td>string</td>
    <td>The username of the user account.</td>
    <td>✅</td>
  </tr>
</table>

#### Payload
<table>
  <tr>
    <th>Field</th>
    <th>Type</th>
    <th>Description</th>
    <th>Required</th>
  </tr>
  <tr>
    <td>
      <code>first_name</code>
    </td>
    <td>string</td>
    <td>The first name of the user.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>last_name</code>
    </td>
    <td>string</td>
    <td>The last name of the user.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>address</code>
    </td>
    <td>string</td>
    <td>The user address.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>email</code>
    </td>
    <td>string</td>
    <td>The user e-mail address.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>mobile_number</code>
    </td>
    <td>string</td>
    <td>The user phone number.</td>
    <td>❌</td>
  </tr>
</table>

#### Sample Request
Payload:
```json
{
  "mobile_number": "(407) 435-6841"
}
```

Response:
```json
{
  "id": "ADMN-878495",
  "user_type": "ADMIN",
  "first_name": "Emily",
  "last_name": "Davis",
  "username": "emilydavis",
  "password": "passwordabc",
  "address": "321 Cedar Road",
  "email": "emilydavis@example.com",
  "mobile_number": "(407) 435-6841",
  "date_created": "1687849585"
}
```