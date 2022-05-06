# vmware-download-sdk
This SDK builds a layer of abstraction above customerconnect.vmware.com to hide the complexity for the client. This allows for downloads to be requested using the minimum of information.

## Usage
See test `TestFetchDownloadLinkVersionGlob` in `download_test.go` for an end to end example of how to use the SDK to download a file based on a version glob, meaning the latest version matching the pattern is downloaded.

By setting product to "DownloadGroup", subroduct to an actual download group and version to a productId you can download files that are not part of the official product lists.