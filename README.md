# <span style="color:red">**DEPRECATED - This SDK is can now be found [here](https://github.com/vmware-labs/vmware-customer-connect-sdk)**</span>.

## vmware-download-sdk
This SDK builds a layer of abstraction above customerconnect.vmware.com to hide the complexity for the client. This allows for downloads to be requested using the minimum of information.

## Usage
See test `TestFetchDownloadLinkVersionGlob` in `download_test.go` for an end to end example of how to use the SDK to download a file based on a version glob, meaning the latest version matching the pattern is downloaded.