# mockserver-client-go

The mockserver client is a golang client for use with the fantastic http://www.mock-server.com/ HTTP mock server.

Usage:
* Import the library `import "github.com/sandvikcode/mockserver-client-go/pkg/mockclient"`

Create an expectation example:
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

Create a verification example:
```
mockServer := mockclient.Client{
    T: t,
    BaseURL: os.Getenv("MOCKSERVER_HOST"),
}

mockServer.AddVerification(
    mockclient.CreateVerification(
        mockclient.WhenRequestPath("/v1/jobs/(.*)"),
        mockclient.ThenAtLeastCalls(2),
        mockclient.ThenAtMostCalls(4),
    ))
```

Create a verification sequence example:
```
mockServer := mockclient.Client{
    T: t,
    BaseURL: os.Getenv("MOCKSERVER_HOST"),
}

mockServer.AddVerificationSequence(
    mockclient.CreateVerification(
        mockclient.WhenRequestPath("/a"),
    ),
    mockclient.CreateVerification(
        mockclient.WhenRequestPath("/b(.*)"),
    ),
    mockclient.CreateVerification(
        mockclient.WhenRequestPath("/c"),
        mockclient.WhenRequestMethod("POST"),
    ),
)
    
```

Expectation defaults:
* unlimited calls will respond to a match
* calls are not delayed
* status of matched calls is 200 OK
* body of matched calls is empty

Verification defaults:
* matched request occurs once i.e. at 1 least call and at most 1 call

Verification sequence notes:
* only the request part is used for matching the sequence i.e. request count is not applicable

Links:
* Expectations - http://www.mock-server.com/mock_server/creating_expectations.html
* Verifications - http://www.mock-server.com/mock_server/verification.html
* Clearing & Resetting http://www.mock-server.com/mock_server/clearing_and_resetting.html

