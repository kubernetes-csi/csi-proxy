// Code generated by protoc-gen-go. DO NOT EDIT.
// source: volume/v1beta1/api.proto

package v1beta1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ListVolumesOnDiskRequest struct {
	// Disk device ID of the disk to query for volumes
	DiskId               string   `protobuf:"bytes,1,opt,name=disk_id,json=diskId,proto3" json:"disk_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListVolumesOnDiskRequest) Reset()         { *m = ListVolumesOnDiskRequest{} }
func (m *ListVolumesOnDiskRequest) String() string { return proto.CompactTextString(m) }
func (*ListVolumesOnDiskRequest) ProtoMessage()    {}
func (*ListVolumesOnDiskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{0}
}

func (m *ListVolumesOnDiskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListVolumesOnDiskRequest.Unmarshal(m, b)
}
func (m *ListVolumesOnDiskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListVolumesOnDiskRequest.Marshal(b, m, deterministic)
}
func (m *ListVolumesOnDiskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListVolumesOnDiskRequest.Merge(m, src)
}
func (m *ListVolumesOnDiskRequest) XXX_Size() int {
	return xxx_messageInfo_ListVolumesOnDiskRequest.Size(m)
}
func (m *ListVolumesOnDiskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListVolumesOnDiskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListVolumesOnDiskRequest proto.InternalMessageInfo

func (m *ListVolumesOnDiskRequest) GetDiskId() string {
	if m != nil {
		return m.DiskId
	}
	return ""
}

type ListVolumesOnDiskResponse struct {
	// Volume device IDs of volumes on the specified disk
	VolumeIds            []string `protobuf:"bytes,1,rep,name=volume_ids,json=volumeIds,proto3" json:"volume_ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListVolumesOnDiskResponse) Reset()         { *m = ListVolumesOnDiskResponse{} }
func (m *ListVolumesOnDiskResponse) String() string { return proto.CompactTextString(m) }
func (*ListVolumesOnDiskResponse) ProtoMessage()    {}
func (*ListVolumesOnDiskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{1}
}

func (m *ListVolumesOnDiskResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListVolumesOnDiskResponse.Unmarshal(m, b)
}
func (m *ListVolumesOnDiskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListVolumesOnDiskResponse.Marshal(b, m, deterministic)
}
func (m *ListVolumesOnDiskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListVolumesOnDiskResponse.Merge(m, src)
}
func (m *ListVolumesOnDiskResponse) XXX_Size() int {
	return xxx_messageInfo_ListVolumesOnDiskResponse.Size(m)
}
func (m *ListVolumesOnDiskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListVolumesOnDiskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListVolumesOnDiskResponse proto.InternalMessageInfo

func (m *ListVolumesOnDiskResponse) GetVolumeIds() []string {
	if m != nil {
		return m.VolumeIds
	}
	return nil
}

type MountVolumeRequest struct {
	// Volume device ID of the volume to mount
	VolumeId string `protobuf:"bytes,1,opt,name=volume_id,json=volumeId,proto3" json:"volume_id,omitempty"`
	// Path in the host's file system where the volume needs to be mounted
	Path                 string   `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MountVolumeRequest) Reset()         { *m = MountVolumeRequest{} }
func (m *MountVolumeRequest) String() string { return proto.CompactTextString(m) }
func (*MountVolumeRequest) ProtoMessage()    {}
func (*MountVolumeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{2}
}

func (m *MountVolumeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MountVolumeRequest.Unmarshal(m, b)
}
func (m *MountVolumeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MountVolumeRequest.Marshal(b, m, deterministic)
}
func (m *MountVolumeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MountVolumeRequest.Merge(m, src)
}
func (m *MountVolumeRequest) XXX_Size() int {
	return xxx_messageInfo_MountVolumeRequest.Size(m)
}
func (m *MountVolumeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MountVolumeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MountVolumeRequest proto.InternalMessageInfo

func (m *MountVolumeRequest) GetVolumeId() string {
	if m != nil {
		return m.VolumeId
	}
	return ""
}

func (m *MountVolumeRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

type MountVolumeResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MountVolumeResponse) Reset()         { *m = MountVolumeResponse{} }
func (m *MountVolumeResponse) String() string { return proto.CompactTextString(m) }
func (*MountVolumeResponse) ProtoMessage()    {}
func (*MountVolumeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{3}
}

func (m *MountVolumeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MountVolumeResponse.Unmarshal(m, b)
}
func (m *MountVolumeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MountVolumeResponse.Marshal(b, m, deterministic)
}
func (m *MountVolumeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MountVolumeResponse.Merge(m, src)
}
func (m *MountVolumeResponse) XXX_Size() int {
	return xxx_messageInfo_MountVolumeResponse.Size(m)
}
func (m *MountVolumeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MountVolumeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MountVolumeResponse proto.InternalMessageInfo

type DismountVolumeRequest struct {
	// Volume device ID of the volume to dismount
	VolumeId string `protobuf:"bytes,1,opt,name=volume_id,json=volumeId,proto3" json:"volume_id,omitempty"`
	// Path where the volume has been mounted.
	Path                 string   `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DismountVolumeRequest) Reset()         { *m = DismountVolumeRequest{} }
func (m *DismountVolumeRequest) String() string { return proto.CompactTextString(m) }
func (*DismountVolumeRequest) ProtoMessage()    {}
func (*DismountVolumeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{4}
}

func (m *DismountVolumeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DismountVolumeRequest.Unmarshal(m, b)
}
func (m *DismountVolumeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DismountVolumeRequest.Marshal(b, m, deterministic)
}
func (m *DismountVolumeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DismountVolumeRequest.Merge(m, src)
}
func (m *DismountVolumeRequest) XXX_Size() int {
	return xxx_messageInfo_DismountVolumeRequest.Size(m)
}
func (m *DismountVolumeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DismountVolumeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DismountVolumeRequest proto.InternalMessageInfo

func (m *DismountVolumeRequest) GetVolumeId() string {
	if m != nil {
		return m.VolumeId
	}
	return ""
}

func (m *DismountVolumeRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

type DismountVolumeResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DismountVolumeResponse) Reset()         { *m = DismountVolumeResponse{} }
func (m *DismountVolumeResponse) String() string { return proto.CompactTextString(m) }
func (*DismountVolumeResponse) ProtoMessage()    {}
func (*DismountVolumeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{5}
}

func (m *DismountVolumeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DismountVolumeResponse.Unmarshal(m, b)
}
func (m *DismountVolumeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DismountVolumeResponse.Marshal(b, m, deterministic)
}
func (m *DismountVolumeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DismountVolumeResponse.Merge(m, src)
}
func (m *DismountVolumeResponse) XXX_Size() int {
	return xxx_messageInfo_DismountVolumeResponse.Size(m)
}
func (m *DismountVolumeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DismountVolumeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DismountVolumeResponse proto.InternalMessageInfo

type IsVolumeFormattedRequest struct {
	// Volume device ID of the volume to check
	VolumeId             string   `protobuf:"bytes,1,opt,name=volume_id,json=volumeId,proto3" json:"volume_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IsVolumeFormattedRequest) Reset()         { *m = IsVolumeFormattedRequest{} }
func (m *IsVolumeFormattedRequest) String() string { return proto.CompactTextString(m) }
func (*IsVolumeFormattedRequest) ProtoMessage()    {}
func (*IsVolumeFormattedRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{6}
}

func (m *IsVolumeFormattedRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IsVolumeFormattedRequest.Unmarshal(m, b)
}
func (m *IsVolumeFormattedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IsVolumeFormattedRequest.Marshal(b, m, deterministic)
}
func (m *IsVolumeFormattedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IsVolumeFormattedRequest.Merge(m, src)
}
func (m *IsVolumeFormattedRequest) XXX_Size() int {
	return xxx_messageInfo_IsVolumeFormattedRequest.Size(m)
}
func (m *IsVolumeFormattedRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IsVolumeFormattedRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IsVolumeFormattedRequest proto.InternalMessageInfo

func (m *IsVolumeFormattedRequest) GetVolumeId() string {
	if m != nil {
		return m.VolumeId
	}
	return ""
}

type IsVolumeFormattedResponse struct {
	// Is the volume formatted with NTFS
	Formatted            bool     `protobuf:"varint,1,opt,name=formatted,proto3" json:"formatted,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IsVolumeFormattedResponse) Reset()         { *m = IsVolumeFormattedResponse{} }
func (m *IsVolumeFormattedResponse) String() string { return proto.CompactTextString(m) }
func (*IsVolumeFormattedResponse) ProtoMessage()    {}
func (*IsVolumeFormattedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{7}
}

func (m *IsVolumeFormattedResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IsVolumeFormattedResponse.Unmarshal(m, b)
}
func (m *IsVolumeFormattedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IsVolumeFormattedResponse.Marshal(b, m, deterministic)
}
func (m *IsVolumeFormattedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IsVolumeFormattedResponse.Merge(m, src)
}
func (m *IsVolumeFormattedResponse) XXX_Size() int {
	return xxx_messageInfo_IsVolumeFormattedResponse.Size(m)
}
func (m *IsVolumeFormattedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IsVolumeFormattedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IsVolumeFormattedResponse proto.InternalMessageInfo

func (m *IsVolumeFormattedResponse) GetFormatted() bool {
	if m != nil {
		return m.Formatted
	}
	return false
}

type FormatVolumeRequest struct {
	// Volume device ID of the volume to format
	VolumeId             string   `protobuf:"bytes,1,opt,name=volume_id,json=volumeId,proto3" json:"volume_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FormatVolumeRequest) Reset()         { *m = FormatVolumeRequest{} }
func (m *FormatVolumeRequest) String() string { return proto.CompactTextString(m) }
func (*FormatVolumeRequest) ProtoMessage()    {}
func (*FormatVolumeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{8}
}

func (m *FormatVolumeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FormatVolumeRequest.Unmarshal(m, b)
}
func (m *FormatVolumeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FormatVolumeRequest.Marshal(b, m, deterministic)
}
func (m *FormatVolumeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FormatVolumeRequest.Merge(m, src)
}
func (m *FormatVolumeRequest) XXX_Size() int {
	return xxx_messageInfo_FormatVolumeRequest.Size(m)
}
func (m *FormatVolumeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FormatVolumeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FormatVolumeRequest proto.InternalMessageInfo

func (m *FormatVolumeRequest) GetVolumeId() string {
	if m != nil {
		return m.VolumeId
	}
	return ""
}

type FormatVolumeResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FormatVolumeResponse) Reset()         { *m = FormatVolumeResponse{} }
func (m *FormatVolumeResponse) String() string { return proto.CompactTextString(m) }
func (*FormatVolumeResponse) ProtoMessage()    {}
func (*FormatVolumeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{9}
}

func (m *FormatVolumeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FormatVolumeResponse.Unmarshal(m, b)
}
func (m *FormatVolumeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FormatVolumeResponse.Marshal(b, m, deterministic)
}
func (m *FormatVolumeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FormatVolumeResponse.Merge(m, src)
}
func (m *FormatVolumeResponse) XXX_Size() int {
	return xxx_messageInfo_FormatVolumeResponse.Size(m)
}
func (m *FormatVolumeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FormatVolumeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FormatVolumeResponse proto.InternalMessageInfo

type ResizeVolumeRequest struct {
	// Volume device ID of the volume to dismount
	VolumeId string `protobuf:"bytes,1,opt,name=volume_id,json=volumeId,proto3" json:"volume_id,omitempty"`
	// New size of the volume
	Size                 int64    `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResizeVolumeRequest) Reset()         { *m = ResizeVolumeRequest{} }
func (m *ResizeVolumeRequest) String() string { return proto.CompactTextString(m) }
func (*ResizeVolumeRequest) ProtoMessage()    {}
func (*ResizeVolumeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{10}
}

func (m *ResizeVolumeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResizeVolumeRequest.Unmarshal(m, b)
}
func (m *ResizeVolumeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResizeVolumeRequest.Marshal(b, m, deterministic)
}
func (m *ResizeVolumeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResizeVolumeRequest.Merge(m, src)
}
func (m *ResizeVolumeRequest) XXX_Size() int {
	return xxx_messageInfo_ResizeVolumeRequest.Size(m)
}
func (m *ResizeVolumeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ResizeVolumeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ResizeVolumeRequest proto.InternalMessageInfo

func (m *ResizeVolumeRequest) GetVolumeId() string {
	if m != nil {
		return m.VolumeId
	}
	return ""
}

func (m *ResizeVolumeRequest) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

type ResizeVolumeResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResizeVolumeResponse) Reset()         { *m = ResizeVolumeResponse{} }
func (m *ResizeVolumeResponse) String() string { return proto.CompactTextString(m) }
func (*ResizeVolumeResponse) ProtoMessage()    {}
func (*ResizeVolumeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{11}
}

func (m *ResizeVolumeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResizeVolumeResponse.Unmarshal(m, b)
}
func (m *ResizeVolumeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResizeVolumeResponse.Marshal(b, m, deterministic)
}
func (m *ResizeVolumeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResizeVolumeResponse.Merge(m, src)
}
func (m *ResizeVolumeResponse) XXX_Size() int {
	return xxx_messageInfo_ResizeVolumeResponse.Size(m)
}
func (m *ResizeVolumeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ResizeVolumeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ResizeVolumeResponse proto.InternalMessageInfo

type VolumeStatsRequest struct {
	// Volume device Id of the volume to dismount
	VolumeId             string   `protobuf:"bytes,1,opt,name=volume_id,json=volumeId,proto3" json:"volume_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VolumeStatsRequest) Reset()         { *m = VolumeStatsRequest{} }
func (m *VolumeStatsRequest) String() string { return proto.CompactTextString(m) }
func (*VolumeStatsRequest) ProtoMessage()    {}
func (*VolumeStatsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{12}
}

func (m *VolumeStatsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VolumeStatsRequest.Unmarshal(m, b)
}
func (m *VolumeStatsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VolumeStatsRequest.Marshal(b, m, deterministic)
}
func (m *VolumeStatsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VolumeStatsRequest.Merge(m, src)
}
func (m *VolumeStatsRequest) XXX_Size() int {
	return xxx_messageInfo_VolumeStatsRequest.Size(m)
}
func (m *VolumeStatsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_VolumeStatsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_VolumeStatsRequest proto.InternalMessageInfo

func (m *VolumeStatsRequest) GetVolumeId() string {
	if m != nil {
		return m.VolumeId
	}
	return ""
}

type VolumeStatsResponse struct {
	//Total size of the volume
	DiskSize int64 `protobuf:"varint,1,opt,name=diskSize,proto3" json:"diskSize,omitempty"`
	//Capacity of the volume
	VolumeSize int64 `protobuf:"varint,2,opt,name=volumeSize,proto3" json:"volumeSize,omitempty"`
	//Used bytes
	VolumeUsedSize       int64    `protobuf:"varint,3,opt,name=volumeUsedSize,proto3" json:"volumeUsedSize,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VolumeStatsResponse) Reset()         { *m = VolumeStatsResponse{} }
func (m *VolumeStatsResponse) String() string { return proto.CompactTextString(m) }
func (*VolumeStatsResponse) ProtoMessage()    {}
func (*VolumeStatsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e25953f9b6119981, []int{13}
}

func (m *VolumeStatsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VolumeStatsResponse.Unmarshal(m, b)
}
func (m *VolumeStatsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VolumeStatsResponse.Marshal(b, m, deterministic)
}
func (m *VolumeStatsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VolumeStatsResponse.Merge(m, src)
}
func (m *VolumeStatsResponse) XXX_Size() int {
	return xxx_messageInfo_VolumeStatsResponse.Size(m)
}
func (m *VolumeStatsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VolumeStatsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VolumeStatsResponse proto.InternalMessageInfo

func (m *VolumeStatsResponse) GetDiskSize() int64 {
	if m != nil {
		return m.DiskSize
	}
	return 0
}

func (m *VolumeStatsResponse) GetVolumeSize() int64 {
	if m != nil {
		return m.VolumeSize
	}
	return 0
}

func (m *VolumeStatsResponse) GetVolumeUsedSize() int64 {
	if m != nil {
		return m.VolumeUsedSize
	}
	return 0
}

func init() {
	proto.RegisterType((*ListVolumesOnDiskRequest)(nil), "v1beta1.ListVolumesOnDiskRequest")
	proto.RegisterType((*ListVolumesOnDiskResponse)(nil), "v1beta1.ListVolumesOnDiskResponse")
	proto.RegisterType((*MountVolumeRequest)(nil), "v1beta1.MountVolumeRequest")
	proto.RegisterType((*MountVolumeResponse)(nil), "v1beta1.MountVolumeResponse")
	proto.RegisterType((*DismountVolumeRequest)(nil), "v1beta1.DismountVolumeRequest")
	proto.RegisterType((*DismountVolumeResponse)(nil), "v1beta1.DismountVolumeResponse")
	proto.RegisterType((*IsVolumeFormattedRequest)(nil), "v1beta1.IsVolumeFormattedRequest")
	proto.RegisterType((*IsVolumeFormattedResponse)(nil), "v1beta1.IsVolumeFormattedResponse")
	proto.RegisterType((*FormatVolumeRequest)(nil), "v1beta1.FormatVolumeRequest")
	proto.RegisterType((*FormatVolumeResponse)(nil), "v1beta1.FormatVolumeResponse")
	proto.RegisterType((*ResizeVolumeRequest)(nil), "v1beta1.ResizeVolumeRequest")
	proto.RegisterType((*ResizeVolumeResponse)(nil), "v1beta1.ResizeVolumeResponse")
	proto.RegisterType((*VolumeStatsRequest)(nil), "v1beta1.VolumeStatsRequest")
	proto.RegisterType((*VolumeStatsResponse)(nil), "v1beta1.VolumeStatsResponse")
}

func init() { proto.RegisterFile("volume/v1beta1/api.proto", fileDescriptor_e25953f9b6119981) }

var fileDescriptor_e25953f9b6119981 = []byte{
	// 508 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x5d, 0x6f, 0xd3, 0x30,
	0x14, 0x6d, 0xe9, 0xd4, 0x35, 0x17, 0x34, 0x69, 0xb7, 0x6c, 0x64, 0xd9, 0x07, 0xc3, 0x0f, 0x68,
	0x2f, 0x6b, 0xd4, 0xed, 0x01, 0x81, 0x84, 0x90, 0xd0, 0x98, 0x28, 0x62, 0x42, 0x4a, 0x05, 0x0f,
	0x68, 0xd2, 0x94, 0x36, 0x86, 0x59, 0x5d, 0x93, 0xd0, 0xeb, 0x4c, 0x8c, 0x3f, 0xc9, 0x5f, 0x42,
	0x71, 0xdc, 0xd4, 0x59, 0x1d, 0x54, 0x24, 0xde, 0xec, 0x63, 0x9f, 0x73, 0x7c, 0xaf, 0x8f, 0x0d,
	0xee, 0x6d, 0x72, 0x93, 0x4d, 0xb9, 0x7f, 0xdb, 0x1f, 0x71, 0x19, 0xf6, 0xfd, 0x30, 0x15, 0xbd,
	0x74, 0x96, 0xc8, 0x04, 0xd7, 0x35, 0xc4, 0x4e, 0xc1, 0xfd, 0x28, 0x48, 0x7e, 0x51, 0x1b, 0xe9,
	0x53, 0x7c, 0x26, 0x68, 0x12, 0xf0, 0x1f, 0x19, 0x27, 0x89, 0x4f, 0x60, 0x3d, 0x12, 0x34, 0xb9,
	0x12, 0x91, 0xdb, 0x3c, 0x6c, 0x1e, 0x39, 0x41, 0x3b, 0x9f, 0x0e, 0x22, 0xf6, 0x0a, 0x76, 0x2c,
	0x24, 0x4a, 0x93, 0x98, 0x38, 0xee, 0x03, 0x14, 0xb6, 0x57, 0x22, 0x22, 0xb7, 0x79, 0xd8, 0x3a,
	0x72, 0x02, 0xa7, 0x40, 0x06, 0x11, 0xb1, 0x77, 0x80, 0x17, 0x49, 0x16, 0x6b, 0xf2, 0xdc, 0x6a,
	0x17, 0x9c, 0x92, 0xa4, 0xcd, 0x3a, 0x73, 0x0e, 0x22, 0xac, 0xa5, 0xa1, 0xbc, 0x76, 0x1f, 0x28,
	0x5c, 0x8d, 0xd9, 0x16, 0x74, 0x2b, 0x32, 0x85, 0x39, 0x7b, 0x0f, 0x5b, 0x67, 0x82, 0xa6, 0xff,
	0xc1, 0xc0, 0x85, 0xed, 0xfb, 0x4a, 0xda, 0xe3, 0x05, 0xb8, 0x03, 0x2a, 0xb0, 0xf3, 0x64, 0x36,
	0x0d, 0xa5, 0xe4, 0xd1, 0x2a, 0x36, 0xec, 0x25, 0xec, 0x58, 0x88, 0xba, 0x6d, 0x7b, 0xe0, 0x7c,
	0x9b, 0x83, 0x8a, 0xd9, 0x09, 0x16, 0x00, 0x3b, 0x81, 0x6e, 0x41, 0x59, 0xbd, 0x2a, 0xb6, 0x0d,
	0x8f, 0xab, 0x1c, 0x7d, 0xfe, 0x73, 0xe8, 0x06, 0x9c, 0xc4, 0x2f, 0xfe, 0x6f, 0x1d, 0xca, 0x19,
	0xaa, 0x43, 0xad, 0x40, 0x8d, 0x73, 0xfd, 0xaa, 0x8e, 0xd6, 0xef, 0x03, 0x16, 0xc8, 0x50, 0x86,
	0x92, 0x56, 0x3a, 0xea, 0x1d, 0x74, 0x2b, 0x14, 0xdd, 0x13, 0x0f, 0x3a, 0x79, 0xe2, 0x86, 0xb9,
	0x73, 0x53, 0x39, 0x97, 0x73, 0x3c, 0x98, 0xc7, 0x6c, 0xb8, 0x38, 0x97, 0x81, 0xe0, 0x73, 0xd8,
	0x28, 0x66, 0x9f, 0x89, 0x47, 0x6a, 0x4f, 0x4b, 0xed, 0xb9, 0x87, 0x9e, 0xfc, 0x5e, 0x83, 0x76,
	0xe1, 0x8d, 0x97, 0xb0, 0xb9, 0x14, 0x6b, 0x7c, 0xd6, 0xd3, 0x4f, 0xa5, 0x57, 0xf7, 0x4e, 0x3c,
	0xf6, 0xb7, 0x2d, 0xba, 0x29, 0x0d, 0xfc, 0x00, 0x0f, 0x8d, 0xc4, 0xe2, 0x6e, 0x49, 0x5a, 0x7e,
	0x0e, 0xde, 0x9e, 0x7d, 0xb1, 0xd4, 0x1a, 0xc2, 0x46, 0x35, 0x9c, 0x78, 0x50, 0x32, 0xac, 0xf9,
	0xf7, 0x9e, 0xd6, 0xae, 0x97, 0xa2, 0x97, 0xb0, 0xb9, 0x14, 0x4f, 0xa3, 0xfc, 0xba, 0xcc, 0x1b,
	0xe5, 0xd7, 0xa6, 0x9b, 0x35, 0xf0, 0x02, 0x1e, 0x99, 0x69, 0xc4, 0x45, 0x89, 0x96, 0x60, 0x7b,
	0xfb, 0x35, 0xab, 0xa6, 0x9c, 0x19, 0x3e, 0x43, 0xce, 0x92, 0x6d, 0x43, 0xce, 0x9a, 0x58, 0x75,
	0x39, 0x46, 0x00, 0x8d, 0xcb, 0x59, 0x4e, 0xb2, 0x71, 0x39, 0x96, 0xcc, 0xb2, 0xc6, 0xdb, 0x37,
	0x5f, 0x5f, 0x7f, 0x17, 0xf2, 0x3a, 0x1b, 0xf5, 0xc6, 0xc9, 0xd4, 0x9f, 0x64, 0x23, 0x3e, 0x8b,
	0xb9, 0xe4, 0x74, 0x3c, 0x26, 0xe1, 0x8f, 0x49, 0x1c, 0xa7, 0xb3, 0xe4, 0xe7, 0x9d, 0x3f, 0xbe,
	0x11, 0x3c, 0x96, 0xf9, 0x9f, 0xec, 0x57, 0xbf, 0xe9, 0x51, 0x5b, 0xfd, 0xd1, 0xa7, 0x7f, 0x02,
	0x00, 0x00, 0xff, 0xff, 0x2d, 0x6a, 0x0f, 0xe7, 0xbf, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// VolumeClient is the client API for Volume service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type VolumeClient interface {
	// ListVolumesOnDisk returns the volume IDs (in \\.\Volume{GUID} format) for
	// all volumes on a Disk device
	ListVolumesOnDisk(ctx context.Context, in *ListVolumesOnDiskRequest, opts ...grpc.CallOption) (*ListVolumesOnDiskResponse, error)
	// MountVolume mounts the volume at the requested global staging path
	MountVolume(ctx context.Context, in *MountVolumeRequest, opts ...grpc.CallOption) (*MountVolumeResponse, error)
	// DismountVolume gracefully dismounts a volume
	DismountVolume(ctx context.Context, in *DismountVolumeRequest, opts ...grpc.CallOption) (*DismountVolumeResponse, error)
	// IsVolumeFormatted checks if a volume is formatted with NTFS
	IsVolumeFormatted(ctx context.Context, in *IsVolumeFormattedRequest, opts ...grpc.CallOption) (*IsVolumeFormattedResponse, error)
	// FormatVolume formats a volume with the provided file system
	FormatVolume(ctx context.Context, in *FormatVolumeRequest, opts ...grpc.CallOption) (*FormatVolumeResponse, error)
	// ResizeVolume performs resizing of the partition and file system for a block based volume
	ResizeVolume(ctx context.Context, in *ResizeVolumeRequest, opts ...grpc.CallOption) (*ResizeVolumeResponse, error)
	// VolumeStats gathers DiskSize, VolumeSize and VolumeUsedSize for a volume
	VolumeStats(ctx context.Context, in *VolumeStatsRequest, opts ...grpc.CallOption) (*VolumeStatsResponse, error)
}

type volumeClient struct {
	cc *grpc.ClientConn
}

func NewVolumeClient(cc *grpc.ClientConn) VolumeClient {
	return &volumeClient{cc}
}

func (c *volumeClient) ListVolumesOnDisk(ctx context.Context, in *ListVolumesOnDiskRequest, opts ...grpc.CallOption) (*ListVolumesOnDiskResponse, error) {
	out := new(ListVolumesOnDiskResponse)
	err := c.cc.Invoke(ctx, "/v1beta1.Volume/ListVolumesOnDisk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volumeClient) MountVolume(ctx context.Context, in *MountVolumeRequest, opts ...grpc.CallOption) (*MountVolumeResponse, error) {
	out := new(MountVolumeResponse)
	err := c.cc.Invoke(ctx, "/v1beta1.Volume/MountVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volumeClient) DismountVolume(ctx context.Context, in *DismountVolumeRequest, opts ...grpc.CallOption) (*DismountVolumeResponse, error) {
	out := new(DismountVolumeResponse)
	err := c.cc.Invoke(ctx, "/v1beta1.Volume/DismountVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volumeClient) IsVolumeFormatted(ctx context.Context, in *IsVolumeFormattedRequest, opts ...grpc.CallOption) (*IsVolumeFormattedResponse, error) {
	out := new(IsVolumeFormattedResponse)
	err := c.cc.Invoke(ctx, "/v1beta1.Volume/IsVolumeFormatted", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volumeClient) FormatVolume(ctx context.Context, in *FormatVolumeRequest, opts ...grpc.CallOption) (*FormatVolumeResponse, error) {
	out := new(FormatVolumeResponse)
	err := c.cc.Invoke(ctx, "/v1beta1.Volume/FormatVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volumeClient) ResizeVolume(ctx context.Context, in *ResizeVolumeRequest, opts ...grpc.CallOption) (*ResizeVolumeResponse, error) {
	out := new(ResizeVolumeResponse)
	err := c.cc.Invoke(ctx, "/v1beta1.Volume/ResizeVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volumeClient) VolumeStats(ctx context.Context, in *VolumeStatsRequest, opts ...grpc.CallOption) (*VolumeStatsResponse, error) {
	out := new(VolumeStatsResponse)
	err := c.cc.Invoke(ctx, "/v1beta1.Volume/VolumeStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VolumeServer is the server API for Volume service.
type VolumeServer interface {
	// ListVolumesOnDisk returns the volume IDs (in \\.\Volume{GUID} format) for
	// all volumes on a Disk device
	ListVolumesOnDisk(context.Context, *ListVolumesOnDiskRequest) (*ListVolumesOnDiskResponse, error)
	// MountVolume mounts the volume at the requested global staging path
	MountVolume(context.Context, *MountVolumeRequest) (*MountVolumeResponse, error)
	// DismountVolume gracefully dismounts a volume
	DismountVolume(context.Context, *DismountVolumeRequest) (*DismountVolumeResponse, error)
	// IsVolumeFormatted checks if a volume is formatted with NTFS
	IsVolumeFormatted(context.Context, *IsVolumeFormattedRequest) (*IsVolumeFormattedResponse, error)
	// FormatVolume formats a volume with the provided file system
	FormatVolume(context.Context, *FormatVolumeRequest) (*FormatVolumeResponse, error)
	// ResizeVolume performs resizing of the partition and file system for a block based volume
	ResizeVolume(context.Context, *ResizeVolumeRequest) (*ResizeVolumeResponse, error)
	// VolumeStats gathers DiskSize, VolumeSize and VolumeUsedSize for a volume
	VolumeStats(context.Context, *VolumeStatsRequest) (*VolumeStatsResponse, error)
}

// UnimplementedVolumeServer can be embedded to have forward compatible implementations.
type UnimplementedVolumeServer struct {
}

func (*UnimplementedVolumeServer) ListVolumesOnDisk(ctx context.Context, req *ListVolumesOnDiskRequest) (*ListVolumesOnDiskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVolumesOnDisk not implemented")
}
func (*UnimplementedVolumeServer) MountVolume(ctx context.Context, req *MountVolumeRequest) (*MountVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MountVolume not implemented")
}
func (*UnimplementedVolumeServer) DismountVolume(ctx context.Context, req *DismountVolumeRequest) (*DismountVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DismountVolume not implemented")
}
func (*UnimplementedVolumeServer) IsVolumeFormatted(ctx context.Context, req *IsVolumeFormattedRequest) (*IsVolumeFormattedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsVolumeFormatted not implemented")
}
func (*UnimplementedVolumeServer) FormatVolume(ctx context.Context, req *FormatVolumeRequest) (*FormatVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FormatVolume not implemented")
}
func (*UnimplementedVolumeServer) ResizeVolume(ctx context.Context, req *ResizeVolumeRequest) (*ResizeVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResizeVolume not implemented")
}
func (*UnimplementedVolumeServer) VolumeStats(ctx context.Context, req *VolumeStatsRequest) (*VolumeStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VolumeStats not implemented")
}

func RegisterVolumeServer(s *grpc.Server, srv VolumeServer) {
	s.RegisterService(&_Volume_serviceDesc, srv)
}

func _Volume_ListVolumesOnDisk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListVolumesOnDiskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolumeServer).ListVolumesOnDisk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1beta1.Volume/ListVolumesOnDisk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolumeServer).ListVolumesOnDisk(ctx, req.(*ListVolumesOnDiskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Volume_MountVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MountVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolumeServer).MountVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1beta1.Volume/MountVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolumeServer).MountVolume(ctx, req.(*MountVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Volume_DismountVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DismountVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolumeServer).DismountVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1beta1.Volume/DismountVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolumeServer).DismountVolume(ctx, req.(*DismountVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Volume_IsVolumeFormatted_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsVolumeFormattedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolumeServer).IsVolumeFormatted(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1beta1.Volume/IsVolumeFormatted",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolumeServer).IsVolumeFormatted(ctx, req.(*IsVolumeFormattedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Volume_FormatVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FormatVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolumeServer).FormatVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1beta1.Volume/FormatVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolumeServer).FormatVolume(ctx, req.(*FormatVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Volume_ResizeVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResizeVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolumeServer).ResizeVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1beta1.Volume/ResizeVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolumeServer).ResizeVolume(ctx, req.(*ResizeVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Volume_VolumeStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VolumeStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolumeServer).VolumeStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1beta1.Volume/VolumeStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolumeServer).VolumeStats(ctx, req.(*VolumeStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Volume_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1beta1.Volume",
	HandlerType: (*VolumeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListVolumesOnDisk",
			Handler:    _Volume_ListVolumesOnDisk_Handler,
		},
		{
			MethodName: "MountVolume",
			Handler:    _Volume_MountVolume_Handler,
		},
		{
			MethodName: "DismountVolume",
			Handler:    _Volume_DismountVolume_Handler,
		},
		{
			MethodName: "IsVolumeFormatted",
			Handler:    _Volume_IsVolumeFormatted_Handler,
		},
		{
			MethodName: "FormatVolume",
			Handler:    _Volume_FormatVolume_Handler,
		},
		{
			MethodName: "ResizeVolume",
			Handler:    _Volume_ResizeVolume_Handler,
		},
		{
			MethodName: "VolumeStats",
			Handler:    _Volume_VolumeStats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "volume/v1beta1/api.proto",
}
