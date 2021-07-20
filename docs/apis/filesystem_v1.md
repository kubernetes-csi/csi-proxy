# CSI Proxy FileSystem v1 API
<a name="top"></a>

## Table of Contents

- [FileSystem RPCs](#v1.FileSystemRPCs)

- [FileSystem Messages](#v1.FileSystemMessages)


<a name="v1.FileSystemRPCs"></a>

## v1 FileSystem RPCs

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| PathExists | [PathExistsRequest](#v1.PathExistsRequest) | [PathExistsResponse](#v1.PathExistsResponse) | PathExists checks if the requested path exists in the host filesystem. |
| Mkdir | [MkdirRequest](#v1.MkdirRequest) | [MkdirResponse](#v1.MkdirResponse) | Mkdir creates a directory at the requested path in the host filesystem. |
| Rmdir | [RmdirRequest](#v1.RmdirRequest) | [RmdirResponse](#v1.RmdirResponse) | Rmdir removes the directory at the requested path in the host filesystem. This may be used for unlinking a symlink created through CreateSymlink. |
| CreateSymlink | [CreateSymlinkRequest](#v1.CreateSymlinkRequest) | [CreateSymlinkResponse](#v1.CreateSymlinkResponse) | CreateSymlink creates a symbolic link called target_path that points to source_path in the host filesystem (target_path is the name of the symbolic link created, source_path is the existing path). |
| IsSymlink | [IsSymlinkRequest](#v1.IsSymlinkRequest) | [IsSymlinkResponse](#v1.IsSymlinkResponse) | IsSymlink checks if a given path is a symlink. |


<a name="v1.FileSystemMessages"></a>
<p align="right"><a href="#top">Top</a></p>
## v1 FileSystem Messages

<a name="v1.CreateSymlinkRequest"></a>
### CreateSymlinkRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| source_path | string |  | The path of the existing directory to be linked. All special characters allowed by Windows in path names will be allowed except for restrictions noted below. For details, please check: https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file
| target_path | string |  | Target path is the location of the new directory entry to be created in the host's filesystem. All special characters allowed by Windows in path names will be allowed except for restrictions noted below. For details, please check: https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file


Restrictions: Only absolute path (indicated by a drive letter prefix: e.g. &#34;C:\&#34;) is accepted. The path prefix needs needs to match the paths specified as kubelet-path parameter of csi-proxy. UNC paths of the form &#34;\\server\share\path\file&#34; are not allowed. All directory separators need to be backslash character: &#34;\&#34;. Characters: .. / : | ? * in the path are not allowed. source_path cannot already exist in the host filesystem. Maximum path length will be capped to 260 characters.

<a name="v1.CreateSymlinkResponse"></a>
### CreateSymlinkResponse
Intentionally empty.

<a name="v1.IsSymlinkRequest"></a>
### IsSymlinkRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | string |  | The path whose existence as a symlink we want to check in the host&#39;s filesystem. |


<a name="v1.IsSymlinkResponse"></a>
### IsSymlinkResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| is_symlink | bool |  | Indicates whether the path in IsSymlinkRequest is a symlink. |

<a name="v1.MkdirRequest"></a>
### MkdirRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | string |  | The path to create in the host&#39;s filesystem. All special characters allowed by Windows in path names will be allowed except for restrictions noted below. For details, please check: https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file Non-existent parent directories in the path will be automatically created. Directories will be created with Read and Write privileges of the Windows User account under which csi-proxy is started (typically LocalSystem).

Restrictions: Only absolute path (indicated by a drive letter prefix: e.g. &#34;C:\&#34;) is accepted. Depending on the context parameter of this function, the path prefix needs to match the paths specified either as kubelet-csi-plugins-path or as kubelet-pod-path parameters of csi-proxy. The path parameter cannot already exist in the host&#39;s filesystem. UNC paths of the form &#34;\\server\share\path\file&#34; are not allowed. All directory separators need to be backslash character: &#34;\&#34;. Characters: .. / : | ? * in the path are not allowed. Maximum path length will be capped to 260 characters.

<a name="v1.MkdirResponse"></a>
### MkdirResponse
Intentionally empty.

<a name="v1.PathExistsRequest"></a>
### PathExistsRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | string |  | The path whose existence we want to check in the host&#39;s filesystem |

<a name="v1.PathExistsResponse"></a>
### PathExistsResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| exists | bool |  | Indicates whether the path in PathExistsRequest exists in the host&#39;s filesystem |

<a name="v1.RmdirRequest"></a>
### RmdirRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | string |  | The path to remove in the host&#39;s filesystem. All special characters allowed by Windows in path names will be allowed except for restrictions noted below. For details, please check: https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file
| force | bool |  | Force remove all contents under path (if any). |

Restrictions: Only absolute path (indicated by a drive letter prefix: e.g. &#34;C:\&#34;) is accepted. Depending on the context parameter of this function, the path prefix needs to match the paths specified either as kubelet-csi-plugins-path or as kubelet-pod-path parameters of csi-proxy. UNC paths of the form &#34;\\server\share\path\file&#34; are not allowed. All directory separators need to be backslash character: &#34;\&#34;. Characters: .. / : | ? * in the path are not allowed. Path cannot be a file of type symlink. Maximum path length will be capped to 260 characters.

<a name="v1.RmdirResponse"></a>
### RmdirResponse
Intentionally empty.
