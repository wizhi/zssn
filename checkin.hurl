POST {{host}}/survivors/{{survivor}}/checkins
{
  "latitude": {{lat}},
  "longitude": {{long}}
}

HTTP 204

GET {{host}}/survivors/{{survivor}}

HTTP 200
[Asserts]
jsonpath "$.Location.Latitude" == {{lat}}
jsonpath "$.Location.Longitude" == {{long}}
