# mockserver-client-go

The mockserver client is a golang client for use with the fantastic http://www.mock-server.com/ HTTP mock server.

Usage:
* Import the library `import "github.com/sandvikcode/mockserver-client-go/pkg/mockclient"`

Example:
```
mockServer := mockclient.Client{
    T: t, 
    BaseURL: os.Getenv("MOCKSERVER_HOST"),
}

mockServer.AddExpectation(
    mockclient.CreateExpectation(
        mockclient.WhenRequestPath("/(.*)"),
        mockclient.ThenResponseStatus(http.StatusOK),
    ))

defer mockServer.Clear("/(.*)")
```

Links:
* Expectations - http://www.mock-server.com/mock_server/creating_expectations.html
* Verifications - http://www.mock-server.com/mock_server/verification.html
* Clearing & Resetting http://www.mock-server.com/mock_server/clearing_and_resetting.html
