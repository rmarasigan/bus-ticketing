# Bus Route
The Bus Route API Schema contains the bus unit routing information. In this module, it will let you:
* [Create a new bus route](#create-bus-route)
* [Get Bus Route Records](#get-bus-route-information)
* [Filter Bus Route Records](#filter-bus-route-record)
* [Update Bus Route Record](#update-bus-route-record)

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
    <td>The unique bus route ID and the primary key.</td>
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
      <code>bus_unit_id</code>
    </td>
    <td>string</td>
    <td>The bus unit ID for the identification of specific bus unit route.</td>
  </tr>
  <tr>
    <td>
      <code>currency_code</code>
    </td>
    <td>string</td>
    <td>The medium of exchange for goods and services
        (<a href = "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies.json">Currency Codes List</a>).
    </td>
  </tr>
  <tr>
    <td>
      <code>rate</code>
    </td>
    <td>number</td>
    <td>The fare charged to the passenger.</td>
  </tr>
  <tr>
    <td>
      <code>active</code>
    </td>
    <td>boolean</td>
    <td>The availability of the bus unit for the defined route.</td>
  </tr>
  <tr>
    <td>
      <code>departure_time</code>
    </td>
    <td>string</td>
    <td>The expected departure time on the starting point and in 24-hour format (e.g. 15:00).</td>
  </tr>
  <tr>
    <td>
      <code>arrival_time</code>
    </td>
    <td>string</td>
    <td>The expected arrival time on the destination and in 24-hour format (e.g. 15:00).</td>
  </tr>
  <tr>
    <td>
      <code>from_route</code>
    </td>
    <td>string</td>
    <td>The starting point of a bus.</td>
  </tr>
  <tr>
    <td>
      <code>to_route</code>
    </td>
    <td>string</td>
    <td>The destination of a bus.</td>
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

### Create Bus Route
To create a new bus route instance, you need to instantiate an object that represents the bus route property. The bus route instance holds the information related to the specific bus unit route.

**Method**: `POST`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-route/create

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
      <code>bus_id</code>
    </td>
    <td>string</td>
    <td>The unique bus ID.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>bus_unit_id</code>
    </td>
    <td>string</td>
    <td>The bus unit ID for the identification of specific bus unit route.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>currency_code</code>
    </td>
    <td>string</td>
    <td>The medium of exchange for goods and services
        (<a href = "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies.json">Currency Codes List</a>).
    </td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>rate</code>
    </td>
    <td>number</td>
    <td>The fare charged to the passenger.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>active</code>
    </td>
    <td>boolean</td>
    <td>The availability of the bus unit for the defined route.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>departure_time</code>
    </td>
    <td>string</td>
    <td>The expected departure time on the starting point and in 24-hour format (e.g. 15:00).</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>arrival_time</code>
    </td>
    <td>string</td>
    <td>The expected arrival time on the destination and in 24-hour format (e.g. 15:00).</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>from_route</code>
    </td>
    <td>string</td>
    <td>The starting point of a bus.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>to_route</code>
    </td>
    <td>string</td>
    <td>The destination of a bus.</td>
    <td>✅</td>
  </tr>
</table>

#### Sample Payload
```json
{
  "rate": 120,
  "active": true,
  "currency_code": "PHP",
  "bus_id": "SNRSBSS-875011",
  "bus_unit_id": "SNRSBSSBUS002",
  "departure_time": "15:00",
  "arrival_time": "17:00",
  "from_route": "Route A",
  "to_route": "Route B"
}
```

### Get Bus Route Information
When retrieving the specific bus route information, the `id` and `bus_id` query parameters must be present in the URL. These parameters identify which information should be returned. It will either return a representation of a specific bus route information or a list of bus route information.

**Method**: `GET`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-route/get

#### Specific Bus Route
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
    <td>The unique bus route ID.</td>
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
    "rate": 120,
    "active": true,
    "currency_code": "PHP",
    "id": "RTRTC15001900877753",
    "bus_id": "SNRSBSS-875011",
    "bus_unit_id": "SNRSBSSBUS002",
    "departure_time": "15:00",
    "arrival_time": "17:00",
    "from_route": "Route A",
    "to_route": "Route B"
  }
]
```


### Filter Bus Route Record
When retrieving a list of bus route records, the `bus_id` query parameter must be present in the URL, and either of the `active`, `bus_unit_id`, `departure_time`, `arrival_time`, `from_route`, or `to_route` is optional in the query parameter. These parameters will identify which bus route record(s) should be returned.

**Method**: `GET`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-route/search?bus_id=xxxxxx

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
      <code>bus_unit_id</code>
    </td>
    <td>string</td>
    <td>The bus unit ID for the identification of specific bus unit route.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>active</code>
    </td>
    <td>boolean</td>
    <td>The availability of the bus unit for the defined route.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>departure_time</code>
    </td>
    <td>string</td>
    <td>The expected departure time on the starting point and in 24-hour format (e.g. 15:00).</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>arrival_time</code>
    </td>
    <td>string</td>
    <td>The expected arrival time on the destination and in 24-hour format (e.g. 15:00).</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>from_route</code>
    </td>
    <td>string</td>
    <td>The starting point of a bus.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>to_route</code>
    </td>
    <td>string</td>
    <td>The destination of a bus.</td>
    <td>❌</td>
  </tr>
</table>

#### Sample Response
```json
[
  {
    "id": "RTBRTC15001900880102",
    "bus_id": "SNRSBSS-875011",
    "bus_unit_id": "SNRSBSSBUS002",
    "currency_code": "PHP",
    "rate": 90,
    "active": true,
    "departure_time": "15:00",
    "arrival_time": "19:00",
    "from_route": "Route B",
    "to_route": "Route C",
    "date_created": "1688010233"
  },
  {
    "id": "RTRTB15001900880101",
    "bus_id": "SNRSBSS-875011",
    "bus_unit_id": "SNRSBSSBUS002",
    "currency_code": "PHP",
    "rate": 90,
    "active": true,
    "departure_time": "15:00",
    "arrival_time": "19:00",
    "from_route": "Route A",
    "to_route": "Route B",
    "date_created": "1688010114"
  }
]
```

### Update Bus Route Record
When modifying the bus route record, the `id` and `bus_id` query parameters must be present in the URL. These parameters identify which bus route record should be modified. After the update is performed, it will return a representation of the updated bus route record.

**Method**: `POST`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-route/update?id=xxxxx&bus_id=xxxxx

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
    <td>The unique bus route ID.</td>
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
    <th>Parameter</th>
    <th>Type</th>
    <th>Description</th>
    <th>Required</th>
  </tr>
  <tr>
    <td>
      <code>currency_code</code>
    </td>
    <td>string</td>
    <td>The medium of exchange for goods and services
        (<a href = "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies.json">Currency Codes List</a>).
    </td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>rate</code>
    </td>
    <td>number</td>
    <td>The fare charged to the passenger.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>active</code>
    </td>
    <td>boolean</td>
    <td>The availability of the bus unit for the defined route.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>departure_time</code>
    </td>
    <td>string</td>
    <td>The expected departure time on the starting point and in 24-hour format (e.g. 15:00).</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>arrival_time</code>
    </td>
    <td>string</td>
    <td>The expected arrival time on the destination and in 24-hour format (e.g. 15:00).</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>from_route</code>
    </td>
    <td>string</td>
    <td>The starting point of a bus.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>to_route</code>
    </td>
    <td>string</td>
    <td>The destination of a bus.</td>
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
  "id": "RTRTB15001900880101",
  "bus_id": "SNRSBSS-875011",
  "bus_unit_id": "SNRSBSSBUS002",
  "currency_code": "PHP",
  "rate": 90,
  "active": false,
  "departure_time": "15:00",
  "arrival_time": "19:00",
  "from_route": "Route A",
  "to_route": "Route B",
  "date_created": "1688010114"
}
```