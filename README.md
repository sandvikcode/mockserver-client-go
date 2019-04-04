# flames-library-mockserver-client

The mockserver client is a golang client for use with the fantastic http://www.mock-server.com/ HTTP mock server.

Usage:
* Import the library `import "github.com/sandvikcode/flames-library-mockserver-client/pkg/mock"`

Example:
```
mockServer := mock.Client{
    T: t, 
    BaseURL: os.Getenv("MOCKSERVER_HOST")
}

mockServer.AddExpectation(
    mock.CreateExpectation(
        mock.WhenRequestPath("/(.*)"),
        mock.ThenResponseStatus(http.StatusOK),
    ))

defer mockServer.Clear("/(.*)")
```

Links:
* Expectations - http://www.mock-server.com/mock_server/creating_expectations.html
* Verifications - http://www.mock-server.com/mock_server/verification.html
* Clearing & Resetting http://www.mock-server.com/mock_server/clearing_and_resetting.html
