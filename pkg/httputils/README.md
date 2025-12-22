# httputils

High level easy to use HTTP client and testserver.

## Examples

* [BasicAuth protection for http endpoints](basicauth/README.md)
    * [Get request using BasicAuth](Example_PerformGetRequestBasicAuth_test.go)
* [Download as file](Example_DownloadAsFile_test.go)
* [GET request](Example_PerformGetRequest_test.go)
    * [Get request using BasicAuth](Example_PerformGetRequestBasicAuth_test.go)
    * [GET request of a nonexisting page: 404 not found](Example_PerformGetRequest404_test.go)
    * [GET JSON data and run jq](Example_GetJsonDataAndRunJq_test.go)
    * [GET YAML data and run yq](Example_GetYamlDataAndRunYq_test.go)
* HTTPClient:
    * [Set base URL on client](Example_SetBaseUrlOnClient_test.go): This is useful if many requests are send to the same webserver using the same client.
* [POST request](Example_PostRequest_test.go)
