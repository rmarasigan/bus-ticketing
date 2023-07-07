# Bookings
The Bookings API Schema contains the details of reserving seats for a particular bus. In this module, it will let you:
* [Create a new bus booking](#create-a-booking)
* [Get Booking Records](#get-booking-records)
* [Get Cancelled Booking Records]()
* [Filter Booking Records](#filter-booking-records)
* [Update Booking Status Record](#update-booking-status-record)

## Data Structure
### Bookings
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
    <td>The unique booking ID and the primary key.</td>
  </tr>
  <tr>
    <td>
      <code>bus_id</code>
    </td>
    <td>string</td>
    <td>The unique bus ID.</td>
  </tr>
  <tr>
    <td>
      <code>bus_route_id</code>
    </td>
    <td>string</td>
    <td>The bus route ID and the sort key.</td>
  </tr>
  <tr>
    <td>
      <code>user_id</code>
    </td>
    <td>string</td>
    <td>The unique user ID.</td>
  </tr>
  <tr>
    <td>
      <code>status</code>
    </td>
    <td>string</td>
    <td>
      The status of the particular booking.
      There are 3 different status types: <br />
      - PENDING <br />
      - CONFIRMED <br />
      - CANCELLED
    </td>
  </tr>
  <tr>
    <td>
      <code>seat_number</code>
    </td>
    <td>string</td>
    <td>The specific seat number(s) for the particular booking.</td>
  </tr>
  <tr>
    <td>
      <code>travel_date</code>
    </td>
    <td>string</td>
    <td>The date when to travel.</td>
  </tr>
  <tr>
    <td>
      <code>date_confirmed</code>
    </td>
    <td>string</td>
    <td>The date the booking was confirmed.</td>
  </tr>
  <tr>
    <td>
      <code>is_cancelled</code>
    </td>
    <td>boolean</td>
    <td>Indicates if the booking is cancelled or not.</td>
  </tr>
  <tr>
    <td>
      <code>cancelled</code>
    </td>
    <td>object</td>
    <td>Contains the cancelled booking information.</td>
  </tr>
  <tr>
    <td>
      <code>timestamp</code>
    </td>
    <td>string</td>
    <td>The timestamp when the request was made.</td>
  </tr>
  <tr>
    <td>
      <code>date_created</code>
    </td>
    <td>string</td>
    <td>The date that this booking record was created.</td>
  </tr>
</table>

### Cancelled Bookings
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
    <td>The unique booking cancellation ID.</td>
  </tr>
  <tr>
    <td>
      <code>booking_id</code>
    </td>
    <td>string</td>
    <td>The unique booking ID as the primary key.</td>
  </tr>
  <tr>
    <td>
      <code>reason</code>
    </td>
    <td>string</td>
    <td>The reason for booking cancellation.</td>
  </tr>
  <tr>
    <td>
      <code>cancelled_by</code>
    </td>
    <td>string</td>
    <td>Indicates who cancelled the booking.</td>
  </tr>
  <tr>
    <td>
      <code>date_cancelled</code>
    </td>
    <td>string</td>
    <td>The date that this booking record was cancelled.</td>
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

### Create a Booking
To create a new booking instance, you need to instantiate an object that represents the booking property. The booking instance holds the information related to the details of reserving seats for a particular bus.

**Method**: `POST`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bookings/create

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
      <code>user_id</code>
    </td>
    <td>string</td>
    <td>The unique user ID.</td>
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
      <code>bus_route_id</code>
    </td>
    <td>string</td>
    <td>The bus route ID and the sort key.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>seat_number</code>
    </td>
    <td>string</td>
    <td>The specific seat number(s) for the particular booking.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>travel_date</code>
    </td>
    <td>string</td>
    <td>The date when to travel.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>status</code>
    </td>
    <td>string</td>
    <td>The status of the particular booking. Should be set to "PENDING".</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>timestamp</code>
    </td>
    <td>string</td>
    <td>The timestamp when the request was made.</td>
    <td>✅</td>
  </tr>
</table>

#### Sample Payload
```json
{
  "user_id": "ADMN-878495",
  "bus_id": "BCBSCMPN-884690",
  "bus_route_id": "RTBRTC15001900884691",
  "seat_number": "23,24,25,26",
  "status": "PENDING",
  "timestamp": "2023-07-01 10:30",
  "travel_date": "2023-07-06 19:30"
}
```

### Get Booking Records
When retrieving the specific booking record, the `booking_id` query parameter must be present in the URL. This parameter identify which information should be returned. It will either return a representation of a specific cancelled booking record or a list of cancelled booking record.

**Method**: `GET`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bookings/get

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
    <td>The unique booking ID.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>bus_route_id</code>
    </td>
    <td>string</td>
    <td>The unique bus route ID.</td>
    <td>✅</td>
  </tr>
</table>

#### Sample Response
```json
[
  {
    "id": "bd866a7e-34cd-4ea1-8411-5351a6b76ffd",
    "user_id": "ADMN-878495",
    "bus_id": "BCBSCMPN-884690",
    "bus_route_id": "RTBRTC15001900884691",
    "status": "PENDING",
    "seat_number": "23,24,25,26",
    "travel_date": "2023-07-06 19:30",
    "date_created": "2023-07-05 07:48:26",
    "cancelled": {
      "id": "",
      "booking_id": "",
      "reason": "",
      "cancelled_by": "",
      "date_cancelled": ""
    },
    "timestamp": "2023-07-01 10:30"
  }
]
```

### Get Cancelled Booking Records
When retrieving the specific cancelled booking record, the `booking_id` query parameter must be present in the URL. This parameter identify which information should be returned. It will either return a representation of a specific cancelled booking record or a list of cancelled booking record.

**Method**: `GET`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bookings/cancelled/get

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
      <code>booking_id</code>
    </td>
    <td>string</td>
    <td>The unique booking ID.</td>
    <td>❌</td>
  </tr>
</table>

#### Sample Response
```json
[
  {
    "id": "053607ed-3dc6-40a3-aea4-d7e87fd015f6",
    "booking_id": "ce4e0245-b772-47f8-92fc-0d70cbd511c0",
    "reason": "sample reason",
    "cancelled_by": "ADMN-878495",
    "date_cancelled": "2023-07-05 04:16:41"
  }
]
```


### Filter Booking Records
When retrieving a list of booking records, the `status` query parameter must be present in the URL, and either of the `bus_id`, or `route_id` is optional in the query parameter. These parameters will identify which booking record(s) should be returned.

**Method**: `GET`

**Endpoint**: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bookings/search?status=xxxxxx

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
      <code>status</code>
    </td>
    <td>string</td>
    <td>The status of the particular booking. There are 3 different status types: <br />
      - PENDING <br />
      - CONFIRMED <br />
      - CANCELLED <br />
      Should be defaulted to "ALL" if fetching all records with different statuses.
    </td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>bus_id</code>
    </td>
    <td>string</td>
    <td>The unique bus ID.</td>
    <td>❌</td>
  </tr>
  <tr>
    <td>
      <code>route_id</code>
    </td>
    <td>string</td>
    <td>The unique bus route ID.</td>
    <td>❌</td>
  </tr>
</table>

#### Sample Response
```json
[
  {
    "id": "bd866a7e-34cd-4ea1-8411-5351a6b76ffd",
    "user_id": "ADMN-878495",
    "bus_id": "BCBSCMPN-884690",
    "bus_route_id": "RTBRTC15001900884691",
    "status": "PENDING",
    "seat_number": "23,24,25,26",
    "travel_date": "2023-07-06 19:30",
    "date_created": "2023-07-05 07:48:26",
    "cancelled": {
      "id": "",
      "booking_id": "",
      "reason": "",
      "cancelled_by": "",
      "date_cancelled": ""
    },
    "timestamp": "2023-07-01 10:30"
  }
]
```

### Update Booking Status Record
When modifying the booking record, the `id` and `bus_route_id` query parameters must be present in the URL. These parameters identify which booking record should be modified. After the update is performed, it will return a representation of the updated booking record.

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
    <td>The unique booking ID.</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>route_id</code>
    </td>
    <td>string</td>
    <td>The unique bus route ID.</td>
    <td>✅</td>
  </tr>
</table>

#### Status: `CONFIRMED`
**Payload**
<table>
  <tr>
    <th>Parameter</th>
    <th>Type</th>
    <th>Description</th>
    <th>Required</th>
  </tr>
  <tr>
    <td>
      <code>status</code>
    </td>
    <td>string</td>
    <td>The status of the particular booking. Should be set as "CONFIRMED".</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>seat_number</code>
    </td>
    <td>string</td>
    <td>The specific seat number(s) for the particular booking.</td>
    <td>❌</td>
  </tr>
</table>

**Sample Request**
Payload:
```json
{
  "status": "CONFIRMED",
  "seat_number": "23,24,25"
}
```

#### Status: `CANCELLED`
**Payload**
<table>
  <tr>
    <th>Parameter</th>
    <th>Type</th>
    <th>Description</th>
    <th>Required</th>
  </tr>
  <tr>
    <td>
      <code>status</code>
    </td>
    <td>string</td>
    <td>The status of the particular booking. Should be set as "CANCELLED".</td>
    <td>✅</td>
  </tr>
  <tr>
    <td>
      <code>cancelled</code>
    </td>
    <td>object</td>
    <td>Contains the cancelled booking information (reason and cancelled_by fields).</td>
    <td>✅</td>
  </tr>
</table>

**Sample Request**
Payload:
```json
{
  "status": "CANCELLED",
  "cancelled": {
    "reason": "sample reason",
    "cancelled_by": "ADMN-878495"
  }
}
```