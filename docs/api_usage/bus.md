# Bus
The Bus API Schema contains the bus information. In this module, it will let you:
* [Create a new bus information](#create-bus-information)
* [Get specific Bus Information](#get-bus-information)
* [Filter Bus Record](#filter-bus-record)
* [Update Bus Record](#update-bus-record)

## Table Structure
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
    <td>The unique bus ID and the sort key.</td>
  </tr>
  <tr>
    <td>
      <code>name</code>
    </td>
    <td>string</td>
    <td>
     The name of the bus line and is the primary key.
    </td>
  </tr>
  <tr>
    <td>
      <code>owner</code>
    </td>
    <td>string</td>
    <td>
      The name of bus company owner.
    </td>
  </tr>
  <tr>
    <td>
      <code>email</code>
    </td>
    <td>string</td>
    <td>The bus company e-mail address.</td>
  </tr>
  <tr>
    <td>
      <code>address</code>
    </td>
    <td>string</td>
    <td>The bus company address.</td>
  </tr>
  <tr>
    <td>
      <code>company</code>
    </td>
    <td>string</td>
    <td>
      The name of the bus company and is the sort key.
    </td>
  </tr>
  <tr>
    <td>
      <code>mobile_number</code>
    </td>
    <td>string</td>
    <td>The bus company phone number.</td>
  </tr>
  <tr>
    <td>
      <code>date_created</code>
    </td>
    <td>string</td>
    <td>The date that this bus information was created.</td>
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
      <code>
        Content-Type
      </code>
    </td>
    <td>
      <code>
        application/json
      </code>
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

### Create Bus Information
To create a new bus instance, you must initialize an array of objects representing buses. It should contain at least one item in the array and each item represents specific bus properties.

**Method**: `POST`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus/create

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
      <code>name</code>
    </td>
    <td>string</td>
    <td>The name of the bus line.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>owner</code>
    </td>
    <td>string</td>
    <td>The name of bus company owner.</td>
    <td>✅</td>
  </tr>
   <tr>
    <td>
      <code>email</code>
    </td>
    <td>string</td>
    <td>The bus company e-mail address.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>address</code>
    </td>
    <td>string</td>
    <td>The bus company address.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>company</code>
    </td>
    <td>string</td>
    <td>The name of the bus company.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>mobile_number</code>
    </td>
    <td>string</td>
    <td>The bus company phone number.</td>
    <td>✅</td>
  </tr>
</table>

#### Sample Payload
```json
[
  {
    "name": "Blue Horizon",
    "owner": "John Doe",
    "email": "john.doe@example.com",
    "address": "123 Main Street, City",
    "company": "ABC Bus Company",
    "mobile_number": "123-456-7890"
  },
  {
    "name": "Green Wave",
    "owner": "Jane Smith",
    "email": "jane.smith@example.com",
    "address": "456 Elm Avenue, Town",
    "company": "XYZ Bus Services",
    "mobile_number": "987-654-3210"
  }
]
```

### Get Bus Information
When retrieving the bus information, the `id` and `name` query parameters must be present in the URL. These parameters identify which bus information should be returned. It will return a representation of a specific bus information.

**Method**: `GET`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus/get?id=xxxxxx&name=xxxxxx

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
    <td>The unique bus ID.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>name</code>
    </td>
    <td>string</td>
    <td>The name of the bus line.</td>
    <td>✅</td>
  </tr>
</table>

#### Sample Response
```json
{
  "id": "BCBSCMPN-875011",
  "name": "Blue Horizon",
  "owner": "John Doe",
  "email": "john.doe@example.com",
  "address": "123 Main Street, City",
  "company": "ABC Bus Company",
  "mobile_number": "123-456-7890"
}
```

### Filter Bus Record
When retrieving a list of bus records, either `name` or `company` query parameter must be present in the URL. This parameter will identify which bus records should be returned.

**Method**: `GET`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus/search?name=xxxxxx&company=xxxxxx

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
      <code>name</code>
    </td>
    <td>string</td>
    <td>The name of the bus line.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>company</code>
    </td>
    <td>string</td>
    <td>The name of the bus company.</td>
    <td>❌</td>
  </tr>
</table>

#### Sample Response
```json
[
  {
    "id": "TRNSTBSC-875011",
    "name": "Yellow Sunshine",
    "owner": "Melissa Anderson",
    "email": "melissa.anderson@example.com",
    "address": "741 Oak Avenue, Suburb",
    "company": "Transit Bus Co",
    "mobile_number": "999-333-7777",
    "date_created": "1687501112"
  },
  {
    "id": "BCBSCMPN-875011",
    "name": "Blue Horizon",
    "owner": "John Doe",
    "email": "john.doe@example.com",
    "address": "123 Main Street, City",
    "company": "ABC Bus Company",
    "mobile_number": "123-456-7890",
    "date_created": "1687501112"
  }
]
```

### Update Bus Record
When modifiying the bus record, the `id` and `name` query parameters must be present in the URL. These parameters identify which bus record should be modified. After the update is performed, it will return a representation of the updated bus record.

**Method**: `POST`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus/update?id=xxxxx&name=xxxxx

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
    <td>The unique bus ID.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>name</code>
    </td>
    <td>string</td>
    <td>The name of the bus line.</td>
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
      <code>owner</code>
    </td>
    <td>string</td>
    <td>The name of bus company owner.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>email</code>
    </td>
    <td>string</td>
    <td>The bus company e-mail address.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>address</code>
    </td>
    <td>string</td>
    <td>The bus company address.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>mobile_number</code>
    </td>
    <td>string</td>
    <td>The bus company phone number.</td>
    <td>❌</td>
  </tr>
</table>

#### Sample Request
Payload:
```json
{
  "owner": "Daniel Martinez",
  "email": "daniel.martinez@example.com"
}
```

Response:
```json
{
  "id": "BCBSCMPN-875011",
  "name": "Blue Horizon",
  "owner": "Daniel Martinez",
  "email": "daniel.martinez@example.com",
  "address": "123 Main Street, City",
  "company": "ABC Bus Company",
  "mobile_number": "123-456-7890",
  "date_created": "1687501112"
}
```