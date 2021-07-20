# CSI Proxy SMB v1 API
<a name="top"></a>

## Table of Contents

- [SMB RPCs](#v1.SMBRPCs)

- [SMB Messages](#v1.SMBMessages)


<a name="v1.SMBRPCs"></a>

## v1 SMB RPCs

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| NewSmbGlobalMapping | [NewSmbGlobalMappingRequest](#v1.NewSmbGlobalMappingRequest) | [NewSmbGlobalMappingResponse](#v1.NewSmbGlobalMappingResponse) | NewSmbGlobalMapping creates an SMB mapping on the SMB client to an SMB share. |
| RemoveSmbGlobalMapping | [RemoveSmbGlobalMappingRequest](#v1.RemoveSmbGlobalMappingRequest) | [RemoveSmbGlobalMappingResponse](#v1.RemoveSmbGlobalMappingResponse) | RemoveSmbGlobalMapping removes the SMB mapping to an SMB share. |


<a name="v1.SMBMessages"></a>
<p align="right"><a href="#top">Top</a></p>

## v1 SMB Messages

<a name="v1.NewSmbGlobalMappingRequest"></a>
### NewSmbGlobalMappingRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| remote_path | string |  | A remote SMB share to mount All unicode characters allowed in SMB server name specifications are permitted except for restrictions below
| local_path | string |  | Optional local path to mount the smb on |
| username | string |  | Username credential associated with the share |
| password | string |  | Password credential associated with the share |

Restrictions: SMB remote path specified in the format: \\server-name\sharename, \\server.fqdn\sharename or \\a.b.c.d\sharename If not an IP address, share name has to be a valid DNS name. UNC specifications to local paths or prefix: \\?\ is not allowed. Characters: &#43; [ ] &#34; / : ; | &lt; &gt; , ? * = $ are not allowed.

<a name="v1.NewSmbGlobalMappingResponse"></a>
### NewSmbGlobalMappingResponse
Intentionally empty.

<a name="v1.RemoveSmbGlobalMappingRequest"></a>
### RemoveSmbGlobalMappingRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| remote_path | string |  | A remote SMB share mapping to remove All unicode characters allowed in SMB server name specifications are permitted except for restrictions below

Restrictions: SMB share specified in the format: \\server-name\sharename, \\server.fqdn\sharename or \\a.b.c.d\sharename If not an IP address, share name has to be a valid DNS name. UNC specifications to local paths or prefix: \\?\ is not allowed. Characters: &#43; [ ] &#34; / : ; | &lt; &gt; , ? * = $ are not allowed.

<a name="v1.RemoveSmbGlobalMappingResponse"></a>

### RemoveSmbGlobalMappingResponse
Intentionally empty.
