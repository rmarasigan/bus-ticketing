# Bus Unit
The Bus Unit API Schema contains the Bus Company's active bus unit and the unit's capacity. In this module, it will let you:
* [Create a new Bus Unit Records](#create-bus-unit-records)
* [Get specific Bus Unit Record](#get-bus-unit-record)
* [Filter Bus Unit Records](#filter-bus-unit-record)
* [Update Bus Unit Record](#update-bus-unit-record)

## Data Structure
<table>
  <tr>
    <th>Field</th>
    <th>Type</th>
    <th>Description</th>
  </tr>
  <tr>
    <td>
      <code>bus_id</code>
    </td>
    <td>string</td>
    <td>The unique bus ID and the sort key.</td>
  </tr>
  <tr>
    <td>
      <code>code</code>
    </td>
    <td>string</td>
    <td>The code is a unique identification of a Bus Unit and is the primary key.</td>
  </tr>
  <tr>
    <td>
      <code>active</code>
    </td>
    <td>boolean</td>
    <td>Defines if the Bus Unit is on "trip".</td>
  </tr>
  <tr>
    <td>
      <code>min_capacity</code>
    </td>
    <td>number</td>
    <td>The minimum number of passenger.</td>
  </tr>
  <tr>
    <td>
      <code>max_capacity</code>
    </td>
    <td>number</td>
    <td>The maximum number of passenger.</td>
  </tr>
  <tr>
    <td>
      <code>date_created</code>
    </td>
    <td>string</td>
    <td>The date that this bus unit record was created.</td>
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

### Create Bus Unit Records
To create a new bus unit instance, you must initialize an array of objects representing bus units. It should contain at least one item in the array and each item represents specific bus unit properties.

**Method**: `POST`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-unit/create

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
      <code>code</code>
    </td>
    <td>string</td>
    <td>The code is a unique identification of a Bus Unit.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>bus_id</code>
    </td>
    <td>string</td>
    <td>The unique bus ID.</td>
    <td>✅</td>
  </tr>
   <tr>
    <td>
      <code>active</code>
    </td>
    <td>boolean</td>
    <td>Defines if the Bus Unit is on "trip".</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>min_capacity</code>
    </td>
    <td>number</td>
    <td>The minimum number of passenger.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>max_capacity</code>
    </td>
    <td>number</td>
    <td>The maximum number of passenger.</td>
    <td>✅</td>
  </tr>
</table>

#### Sample Payload
```json
[
  {
    "bus_id": "BCBSCMPN-875011",
    "code": "BCBSCMPNBUS001",
    "active": true,
    "min_capacity": 30,
    "max_capacity": 60
  },
  {
    "bus_id": "BCBSCMPN-875011",
    "code": "BCBSCMPNBUS002",
    "active": true,
    "min_capacity": 30,
    "max_capacity": 60
  }
]
```

### Get Bus Unit Record
When retrieving the specific bus unit record, the `code` and `bus_id` query parameters must be present in the URL. These parameters identify which bus unit record should be returned. It will either return a representation of a specific bus unit record or a list of bus unit records.

**Method**: `GET`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-unit/get

#### Specific Bus Unit
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
      <code>code</code>
    </td>
    <td>string</td>
    <td>The code is a unique identification of a Bus Unit.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>bus_id</code>
    </td>
    <td>string</td>
    <td>The unique bus ID.</td>
    <td>✅</td>
  </tr>
</table>

#### Sample Response
```json
[
  {
    "bus_id": "BCBSCMPN-875011",
    "code": "BCBSCMPNBUS002",
    "active": true,
    "min_capacity": 30,
    "max_capacity": 60
  }
]
```

### Filter Bus Unit Record
When retrieving a list of bus unit records, the `bus_id` query parameter must be present in the URL, and either `code` or `active` is optional in the query parameter. These parameters will identify which bus unit record(s) should be returned.

**Method**: `GET`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-unit/search?bus_id=xxxxxx

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
      <code>bus_id</code>
    </td>
    <td>string</td>
    <td>The unique bus ID.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>code</code>
    </td>
    <td>string</td>
    <td>The code is a unique identification of a Bus Unit.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>active</code>
    </td>
    <td>boolean</td>
    <td>Defines if the Bus Unit is on "trip".</td>
    <td>❌</td>
  </tr>
</table>

#### Sample Response
```json
[
  {
    "bus_id": "BCBSCMPN-875011",
    "code": "BCBSCMPNBUS001",
    "active": true,
    "min_capacity": 30,
    "max_capacity": 60,
    "date_created": "1687501761"
  },
  {
    "bus_id": "BCBSCMPN-875011",
    "code": "BCBSCMPNBUS003",
    "active": true,
    "min_capacity": 45,
    "max_capacity": 70,
    "date_created": "1687501761"
  },
  {
    "bus_id": "BCBSCMPN-875011",
    "code": "BCBSCMPNBUS002",
    "active": true,
    "min_capacity": 30,
    "max_capacity": 60,
    "date_created": "1687501761"
  }
]
```

### Update Bus Unit Record
When modifying the bus unit record, the `code` and `bus_id` query parameters must be present in the URL. These parameters identify which bus unit record should be modified. After the update is performed, it will return a representation of the updated bus unit record.

**Method**: `POST`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-unit/update?code=xxxxx&bus_id=xxxxx

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
      <code>code</code>
    </td>
    <td>string</td>
    <td>The code is a unique identification of a Bus Unit.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>bus_id</code>
    </td>
    <td>string</td>
    <td>The unique bus ID.</td>
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
      <code>active</code>
    </td>
    <td>boolean</td>
    <td>Defines if the Bus Unit is on "trip".</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>min_capacity</code>
    </td>
    <td>number</td>
    <td>The minimum number of passenger.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>max_capacity</code>
    </td>
    <td>number</td>
    <td>The maximum number of passenger.</td>
    <td>❌</td>
  </tr>
</table>

#### Sample Request
Payload:
```json
{
  "active": false
}
```

Response:
```json
{
  "bus_id": "BCBSCMPN-875011",
  "code": "BCBSCMPNBUS002",
  "active": false,
  "min_capacity": 30,
  "max_capacity": 60,
  "date_created": "1687501761"
}
```