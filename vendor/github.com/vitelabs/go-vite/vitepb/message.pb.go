// Code generated by protoc-gen-go. DO NOT EDIT.
// source: vitepb/message.proto

package vitepb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type Handshake struct {
	Version              int64    `protobuf:"varint,1,opt,name=Version,proto3" json:"Version,omitempty"`
	NetId                int64    `protobuf:"varint,2,opt,name=NetId,proto3" json:"NetId,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=Name,proto3" json:"Name,omitempty"`
	ID                   []byte   `protobuf:"bytes,4,opt,name=ID,proto3" json:"ID,omitempty"`
	Timestamp            int64    `protobuf:"varint,5,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	Genesis              []byte   `protobuf:"bytes,6,opt,name=Genesis,proto3" json:"Genesis,omitempty"`
	Height               uint64   `protobuf:"varint,7,opt,name=Height,proto3" json:"Height,omitempty"`
	Head                 []byte   `protobuf:"bytes,8,opt,name=Head,proto3" json:"Head,omitempty"`
	FileAddress          []byte   `protobuf:"bytes,9,opt,name=FileAddress,proto3" json:"FileAddress,omitempty"`
	Key                  []byte   `protobuf:"bytes,10,opt,name=Key,proto3" json:"Key,omitempty"`
	Token                []byte   `protobuf:"bytes,11,opt,name=Token,proto3" json:"Token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Handshake) Reset()         { *m = Handshake{} }
func (m *Handshake) String() string { return proto.CompactTextString(m) }
func (*Handshake) ProtoMessage()    {}
func (*Handshake) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{0}
}

func (m *Handshake) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Handshake.Unmarshal(m, b)
}
func (m *Handshake) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Handshake.Marshal(b, m, deterministic)
}
func (m *Handshake) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Handshake.Merge(m, src)
}
func (m *Handshake) XXX_Size() int {
	return xxx_messageInfo_Handshake.Size(m)
}
func (m *Handshake) XXX_DiscardUnknown() {
	xxx_messageInfo_Handshake.DiscardUnknown(m)
}

var xxx_messageInfo_Handshake proto.InternalMessageInfo

func (m *Handshake) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Handshake) GetNetId() int64 {
	if m != nil {
		return m.NetId
	}
	return 0
}

func (m *Handshake) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Handshake) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *Handshake) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Handshake) GetGenesis() []byte {
	if m != nil {
		return m.Genesis
	}
	return nil
}

func (m *Handshake) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Handshake) GetHead() []byte {
	if m != nil {
		return m.Head
	}
	return nil
}

func (m *Handshake) GetFileAddress() []byte {
	if m != nil {
		return m.FileAddress
	}
	return nil
}

func (m *Handshake) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Handshake) GetToken() []byte {
	if m != nil {
		return m.Token
	}
	return nil
}

type SyncConnHandshake struct {
	ID                   []byte   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Timestamp            int64    `protobuf:"varint,2,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	Key                  []byte   `protobuf:"bytes,3,opt,name=Key,proto3" json:"Key,omitempty"`
	Token                []byte   `protobuf:"bytes,4,opt,name=Token,proto3" json:"Token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SyncConnHandshake) Reset()         { *m = SyncConnHandshake{} }
func (m *SyncConnHandshake) String() string { return proto.CompactTextString(m) }
func (*SyncConnHandshake) ProtoMessage()    {}
func (*SyncConnHandshake) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{1}
}

func (m *SyncConnHandshake) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncConnHandshake.Unmarshal(m, b)
}
func (m *SyncConnHandshake) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncConnHandshake.Marshal(b, m, deterministic)
}
func (m *SyncConnHandshake) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncConnHandshake.Merge(m, src)
}
func (m *SyncConnHandshake) XXX_Size() int {
	return xxx_messageInfo_SyncConnHandshake.Size(m)
}
func (m *SyncConnHandshake) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncConnHandshake.DiscardUnknown(m)
}

var xxx_messageInfo_SyncConnHandshake proto.InternalMessageInfo

func (m *SyncConnHandshake) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *SyncConnHandshake) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *SyncConnHandshake) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *SyncConnHandshake) GetToken() []byte {
	if m != nil {
		return m.Token
	}
	return nil
}

type ChunkRequest struct {
	From                 uint64   `protobuf:"varint,1,opt,name=From,proto3" json:"From,omitempty"`
	To                   uint64   `protobuf:"varint,2,opt,name=To,proto3" json:"To,omitempty"`
	PrevHash             []byte   `protobuf:"bytes,3,opt,name=PrevHash,proto3" json:"PrevHash,omitempty"`
	EndHash              []byte   `protobuf:"bytes,4,opt,name=EndHash,proto3" json:"EndHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChunkRequest) Reset()         { *m = ChunkRequest{} }
func (m *ChunkRequest) String() string { return proto.CompactTextString(m) }
func (*ChunkRequest) ProtoMessage()    {}
func (*ChunkRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{2}
}

func (m *ChunkRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChunkRequest.Unmarshal(m, b)
}
func (m *ChunkRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChunkRequest.Marshal(b, m, deterministic)
}
func (m *ChunkRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChunkRequest.Merge(m, src)
}
func (m *ChunkRequest) XXX_Size() int {
	return xxx_messageInfo_ChunkRequest.Size(m)
}
func (m *ChunkRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ChunkRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ChunkRequest proto.InternalMessageInfo

func (m *ChunkRequest) GetFrom() uint64 {
	if m != nil {
		return m.From
	}
	return 0
}

func (m *ChunkRequest) GetTo() uint64 {
	if m != nil {
		return m.To
	}
	return 0
}

func (m *ChunkRequest) GetPrevHash() []byte {
	if m != nil {
		return m.PrevHash
	}
	return nil
}

func (m *ChunkRequest) GetEndHash() []byte {
	if m != nil {
		return m.EndHash
	}
	return nil
}

type ChunkResponse struct {
	From                 uint64   `protobuf:"varint,1,opt,name=From,proto3" json:"From,omitempty"`
	To                   uint64   `protobuf:"varint,2,opt,name=To,proto3" json:"To,omitempty"`
	PrevHash             []byte   `protobuf:"bytes,3,opt,name=PrevHash,proto3" json:"PrevHash,omitempty"`
	EndHash              []byte   `protobuf:"bytes,4,opt,name=EndHash,proto3" json:"EndHash,omitempty"`
	Size                 uint64   `protobuf:"varint,5,opt,name=Size,proto3" json:"Size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChunkResponse) Reset()         { *m = ChunkResponse{} }
func (m *ChunkResponse) String() string { return proto.CompactTextString(m) }
func (*ChunkResponse) ProtoMessage()    {}
func (*ChunkResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{3}
}

func (m *ChunkResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChunkResponse.Unmarshal(m, b)
}
func (m *ChunkResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChunkResponse.Marshal(b, m, deterministic)
}
func (m *ChunkResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChunkResponse.Merge(m, src)
}
func (m *ChunkResponse) XXX_Size() int {
	return xxx_messageInfo_ChunkResponse.Size(m)
}
func (m *ChunkResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ChunkResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ChunkResponse proto.InternalMessageInfo

func (m *ChunkResponse) GetFrom() uint64 {
	if m != nil {
		return m.From
	}
	return 0
}

func (m *ChunkResponse) GetTo() uint64 {
	if m != nil {
		return m.To
	}
	return 0
}

func (m *ChunkResponse) GetPrevHash() []byte {
	if m != nil {
		return m.PrevHash
	}
	return nil
}

func (m *ChunkResponse) GetEndHash() []byte {
	if m != nil {
		return m.EndHash
	}
	return nil
}

func (m *ChunkResponse) GetSize() uint64 {
	if m != nil {
		return m.Size
	}
	return 0
}

type State struct {
	Peers                []*State_Peer `protobuf:"bytes,1,rep,name=Peers,proto3" json:"Peers,omitempty"`
	Patch                bool          `protobuf:"varint,2,opt,name=Patch,proto3" json:"Patch,omitempty"`
	Head                 []byte        `protobuf:"bytes,3,opt,name=Head,proto3" json:"Head,omitempty"`
	Height               uint64        `protobuf:"varint,4,opt,name=Height,proto3" json:"Height,omitempty"`
	Timestamp            int64         `protobuf:"varint,10,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *State) Reset()         { *m = State{} }
func (m *State) String() string { return proto.CompactTextString(m) }
func (*State) ProtoMessage()    {}
func (*State) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{4}
}

func (m *State) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_State.Unmarshal(m, b)
}
func (m *State) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_State.Marshal(b, m, deterministic)
}
func (m *State) XXX_Merge(src proto.Message) {
	xxx_messageInfo_State.Merge(m, src)
}
func (m *State) XXX_Size() int {
	return xxx_messageInfo_State.Size(m)
}
func (m *State) XXX_DiscardUnknown() {
	xxx_messageInfo_State.DiscardUnknown(m)
}

var xxx_messageInfo_State proto.InternalMessageInfo

func (m *State) GetPeers() []*State_Peer {
	if m != nil {
		return m.Peers
	}
	return nil
}

func (m *State) GetPatch() bool {
	if m != nil {
		return m.Patch
	}
	return false
}

func (m *State) GetHead() []byte {
	if m != nil {
		return m.Head
	}
	return nil
}

func (m *State) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *State) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type State_Peer struct {
	ID                   []byte   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Add                  bool     `protobuf:"varint,2,opt,name=Add,proto3" json:"Add,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *State_Peer) Reset()         { *m = State_Peer{} }
func (m *State_Peer) String() string { return proto.CompactTextString(m) }
func (*State_Peer) ProtoMessage()    {}
func (*State_Peer) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{4, 0}
}

func (m *State_Peer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_State_Peer.Unmarshal(m, b)
}
func (m *State_Peer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_State_Peer.Marshal(b, m, deterministic)
}
func (m *State_Peer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_State_Peer.Merge(m, src)
}
func (m *State_Peer) XXX_Size() int {
	return xxx_messageInfo_State_Peer.Size(m)
}
func (m *State_Peer) XXX_DiscardUnknown() {
	xxx_messageInfo_State_Peer.DiscardUnknown(m)
}

var xxx_messageInfo_State_Peer proto.InternalMessageInfo

func (m *State_Peer) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *State_Peer) GetAdd() bool {
	if m != nil {
		return m.Add
	}
	return false
}

type HashHeight struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Height               uint64   `protobuf:"varint,2,opt,name=Height,proto3" json:"Height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HashHeight) Reset()         { *m = HashHeight{} }
func (m *HashHeight) String() string { return proto.CompactTextString(m) }
func (*HashHeight) ProtoMessage()    {}
func (*HashHeight) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{5}
}

func (m *HashHeight) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HashHeight.Unmarshal(m, b)
}
func (m *HashHeight) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HashHeight.Marshal(b, m, deterministic)
}
func (m *HashHeight) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HashHeight.Merge(m, src)
}
func (m *HashHeight) XXX_Size() int {
	return xxx_messageInfo_HashHeight.Size(m)
}
func (m *HashHeight) XXX_DiscardUnknown() {
	xxx_messageInfo_HashHeight.DiscardUnknown(m)
}

var xxx_messageInfo_HashHeight proto.InternalMessageInfo

func (m *HashHeight) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *HashHeight) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type HashHeightPoint struct {
	Point                *HashHeight `protobuf:"bytes,1,opt,name=Point,proto3" json:"Point,omitempty"`
	Size                 uint64      `protobuf:"varint,2,opt,name=Size,proto3" json:"Size,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *HashHeightPoint) Reset()         { *m = HashHeightPoint{} }
func (m *HashHeightPoint) String() string { return proto.CompactTextString(m) }
func (*HashHeightPoint) ProtoMessage()    {}
func (*HashHeightPoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{6}
}

func (m *HashHeightPoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HashHeightPoint.Unmarshal(m, b)
}
func (m *HashHeightPoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HashHeightPoint.Marshal(b, m, deterministic)
}
func (m *HashHeightPoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HashHeightPoint.Merge(m, src)
}
func (m *HashHeightPoint) XXX_Size() int {
	return xxx_messageInfo_HashHeightPoint.Size(m)
}
func (m *HashHeightPoint) XXX_DiscardUnknown() {
	xxx_messageInfo_HashHeightPoint.DiscardUnknown(m)
}

var xxx_messageInfo_HashHeightPoint proto.InternalMessageInfo

func (m *HashHeightPoint) GetPoint() *HashHeight {
	if m != nil {
		return m.Point
	}
	return nil
}

func (m *HashHeightPoint) GetSize() uint64 {
	if m != nil {
		return m.Size
	}
	return 0
}

type HashHeightList struct {
	Points               []*HashHeightPoint `protobuf:"bytes,1,rep,name=Points,proto3" json:"Points,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *HashHeightList) Reset()         { *m = HashHeightList{} }
func (m *HashHeightList) String() string { return proto.CompactTextString(m) }
func (*HashHeightList) ProtoMessage()    {}
func (*HashHeightList) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{7}
}

func (m *HashHeightList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HashHeightList.Unmarshal(m, b)
}
func (m *HashHeightList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HashHeightList.Marshal(b, m, deterministic)
}
func (m *HashHeightList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HashHeightList.Merge(m, src)
}
func (m *HashHeightList) XXX_Size() int {
	return xxx_messageInfo_HashHeightList.Size(m)
}
func (m *HashHeightList) XXX_DiscardUnknown() {
	xxx_messageInfo_HashHeightList.DiscardUnknown(m)
}

var xxx_messageInfo_HashHeightList proto.InternalMessageInfo

func (m *HashHeightList) GetPoints() []*HashHeightPoint {
	if m != nil {
		return m.Points
	}
	return nil
}

type GetHashHeightList struct {
	From                 []*HashHeight `protobuf:"bytes,1,rep,name=From,proto3" json:"From,omitempty"`
	Step                 uint64        `protobuf:"varint,2,opt,name=Step,proto3" json:"Step,omitempty"`
	To                   uint64        `protobuf:"varint,3,opt,name=To,proto3" json:"To,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *GetHashHeightList) Reset()         { *m = GetHashHeightList{} }
func (m *GetHashHeightList) String() string { return proto.CompactTextString(m) }
func (*GetHashHeightList) ProtoMessage()    {}
func (*GetHashHeightList) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{8}
}

func (m *GetHashHeightList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetHashHeightList.Unmarshal(m, b)
}
func (m *GetHashHeightList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetHashHeightList.Marshal(b, m, deterministic)
}
func (m *GetHashHeightList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetHashHeightList.Merge(m, src)
}
func (m *GetHashHeightList) XXX_Size() int {
	return xxx_messageInfo_GetHashHeightList.Size(m)
}
func (m *GetHashHeightList) XXX_DiscardUnknown() {
	xxx_messageInfo_GetHashHeightList.DiscardUnknown(m)
}

var xxx_messageInfo_GetHashHeightList proto.InternalMessageInfo

func (m *GetHashHeightList) GetFrom() []*HashHeight {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *GetHashHeightList) GetStep() uint64 {
	if m != nil {
		return m.Step
	}
	return 0
}

func (m *GetHashHeightList) GetTo() uint64 {
	if m != nil {
		return m.To
	}
	return 0
}

type GetSnapshotBlocks struct {
	From                 *HashHeight `protobuf:"bytes,1,opt,name=From,proto3" json:"From,omitempty"`
	Count                uint64      `protobuf:"varint,2,opt,name=Count,proto3" json:"Count,omitempty"`
	Forward              bool        `protobuf:"varint,3,opt,name=Forward,proto3" json:"Forward,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *GetSnapshotBlocks) Reset()         { *m = GetSnapshotBlocks{} }
func (m *GetSnapshotBlocks) String() string { return proto.CompactTextString(m) }
func (*GetSnapshotBlocks) ProtoMessage()    {}
func (*GetSnapshotBlocks) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{9}
}

func (m *GetSnapshotBlocks) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSnapshotBlocks.Unmarshal(m, b)
}
func (m *GetSnapshotBlocks) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSnapshotBlocks.Marshal(b, m, deterministic)
}
func (m *GetSnapshotBlocks) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSnapshotBlocks.Merge(m, src)
}
func (m *GetSnapshotBlocks) XXX_Size() int {
	return xxx_messageInfo_GetSnapshotBlocks.Size(m)
}
func (m *GetSnapshotBlocks) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSnapshotBlocks.DiscardUnknown(m)
}

var xxx_messageInfo_GetSnapshotBlocks proto.InternalMessageInfo

func (m *GetSnapshotBlocks) GetFrom() *HashHeight {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *GetSnapshotBlocks) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *GetSnapshotBlocks) GetForward() bool {
	if m != nil {
		return m.Forward
	}
	return false
}

type SnapshotBlocks struct {
	Blocks               []*SnapshotBlock `protobuf:"bytes,1,rep,name=Blocks,proto3" json:"Blocks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *SnapshotBlocks) Reset()         { *m = SnapshotBlocks{} }
func (m *SnapshotBlocks) String() string { return proto.CompactTextString(m) }
func (*SnapshotBlocks) ProtoMessage()    {}
func (*SnapshotBlocks) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{10}
}

func (m *SnapshotBlocks) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SnapshotBlocks.Unmarshal(m, b)
}
func (m *SnapshotBlocks) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SnapshotBlocks.Marshal(b, m, deterministic)
}
func (m *SnapshotBlocks) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SnapshotBlocks.Merge(m, src)
}
func (m *SnapshotBlocks) XXX_Size() int {
	return xxx_messageInfo_SnapshotBlocks.Size(m)
}
func (m *SnapshotBlocks) XXX_DiscardUnknown() {
	xxx_messageInfo_SnapshotBlocks.DiscardUnknown(m)
}

var xxx_messageInfo_SnapshotBlocks proto.InternalMessageInfo

func (m *SnapshotBlocks) GetBlocks() []*SnapshotBlock {
	if m != nil {
		return m.Blocks
	}
	return nil
}

type GetAccountBlocks struct {
	Address              []byte      `protobuf:"bytes,1,opt,name=Address,proto3" json:"Address,omitempty"`
	From                 *HashHeight `protobuf:"bytes,2,opt,name=From,proto3" json:"From,omitempty"`
	Count                uint64      `protobuf:"varint,3,opt,name=Count,proto3" json:"Count,omitempty"`
	Forward              bool        `protobuf:"varint,4,opt,name=Forward,proto3" json:"Forward,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *GetAccountBlocks) Reset()         { *m = GetAccountBlocks{} }
func (m *GetAccountBlocks) String() string { return proto.CompactTextString(m) }
func (*GetAccountBlocks) ProtoMessage()    {}
func (*GetAccountBlocks) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{11}
}

func (m *GetAccountBlocks) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAccountBlocks.Unmarshal(m, b)
}
func (m *GetAccountBlocks) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAccountBlocks.Marshal(b, m, deterministic)
}
func (m *GetAccountBlocks) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAccountBlocks.Merge(m, src)
}
func (m *GetAccountBlocks) XXX_Size() int {
	return xxx_messageInfo_GetAccountBlocks.Size(m)
}
func (m *GetAccountBlocks) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAccountBlocks.DiscardUnknown(m)
}

var xxx_messageInfo_GetAccountBlocks proto.InternalMessageInfo

func (m *GetAccountBlocks) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *GetAccountBlocks) GetFrom() *HashHeight {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *GetAccountBlocks) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *GetAccountBlocks) GetForward() bool {
	if m != nil {
		return m.Forward
	}
	return false
}

type AccountBlocks struct {
	Blocks               []*AccountBlock `protobuf:"bytes,1,rep,name=Blocks,proto3" json:"Blocks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *AccountBlocks) Reset()         { *m = AccountBlocks{} }
func (m *AccountBlocks) String() string { return proto.CompactTextString(m) }
func (*AccountBlocks) ProtoMessage()    {}
func (*AccountBlocks) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{12}
}

func (m *AccountBlocks) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountBlocks.Unmarshal(m, b)
}
func (m *AccountBlocks) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountBlocks.Marshal(b, m, deterministic)
}
func (m *AccountBlocks) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountBlocks.Merge(m, src)
}
func (m *AccountBlocks) XXX_Size() int {
	return xxx_messageInfo_AccountBlocks.Size(m)
}
func (m *AccountBlocks) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountBlocks.DiscardUnknown(m)
}

var xxx_messageInfo_AccountBlocks proto.InternalMessageInfo

func (m *AccountBlocks) GetBlocks() []*AccountBlock {
	if m != nil {
		return m.Blocks
	}
	return nil
}

type NewSnapshotBlock struct {
	Block                *SnapshotBlock `protobuf:"bytes,1,opt,name=Block,proto3" json:"Block,omitempty"`
	TTL                  int32          `protobuf:"varint,2,opt,name=TTL,proto3" json:"TTL,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *NewSnapshotBlock) Reset()         { *m = NewSnapshotBlock{} }
func (m *NewSnapshotBlock) String() string { return proto.CompactTextString(m) }
func (*NewSnapshotBlock) ProtoMessage()    {}
func (*NewSnapshotBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{13}
}

func (m *NewSnapshotBlock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewSnapshotBlock.Unmarshal(m, b)
}
func (m *NewSnapshotBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewSnapshotBlock.Marshal(b, m, deterministic)
}
func (m *NewSnapshotBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewSnapshotBlock.Merge(m, src)
}
func (m *NewSnapshotBlock) XXX_Size() int {
	return xxx_messageInfo_NewSnapshotBlock.Size(m)
}
func (m *NewSnapshotBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_NewSnapshotBlock.DiscardUnknown(m)
}

var xxx_messageInfo_NewSnapshotBlock proto.InternalMessageInfo

func (m *NewSnapshotBlock) GetBlock() *SnapshotBlock {
	if m != nil {
		return m.Block
	}
	return nil
}

func (m *NewSnapshotBlock) GetTTL() int32 {
	if m != nil {
		return m.TTL
	}
	return 0
}

type NewAccountBlock struct {
	Block                *AccountBlock `protobuf:"bytes,1,opt,name=Block,proto3" json:"Block,omitempty"`
	TTL                  int32         `protobuf:"varint,2,opt,name=TTL,proto3" json:"TTL,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *NewAccountBlock) Reset()         { *m = NewAccountBlock{} }
func (m *NewAccountBlock) String() string { return proto.CompactTextString(m) }
func (*NewAccountBlock) ProtoMessage()    {}
func (*NewAccountBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{14}
}

func (m *NewAccountBlock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewAccountBlock.Unmarshal(m, b)
}
func (m *NewAccountBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewAccountBlock.Marshal(b, m, deterministic)
}
func (m *NewAccountBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewAccountBlock.Merge(m, src)
}
func (m *NewAccountBlock) XXX_Size() int {
	return xxx_messageInfo_NewAccountBlock.Size(m)
}
func (m *NewAccountBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_NewAccountBlock.DiscardUnknown(m)
}

var xxx_messageInfo_NewAccountBlock proto.InternalMessageInfo

func (m *NewAccountBlock) GetBlock() *AccountBlock {
	if m != nil {
		return m.Block
	}
	return nil
}

func (m *NewAccountBlock) GetTTL() int32 {
	if m != nil {
		return m.TTL
	}
	return 0
}

type NewAccountBlockBytes struct {
	Block                []byte   `protobuf:"bytes,1,opt,name=Block,proto3" json:"Block,omitempty"`
	TTL                  int32    `protobuf:"varint,2,opt,name=TTL,proto3" json:"TTL,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NewAccountBlockBytes) Reset()         { *m = NewAccountBlockBytes{} }
func (m *NewAccountBlockBytes) String() string { return proto.CompactTextString(m) }
func (*NewAccountBlockBytes) ProtoMessage()    {}
func (*NewAccountBlockBytes) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{15}
}

func (m *NewAccountBlockBytes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewAccountBlockBytes.Unmarshal(m, b)
}
func (m *NewAccountBlockBytes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewAccountBlockBytes.Marshal(b, m, deterministic)
}
func (m *NewAccountBlockBytes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewAccountBlockBytes.Merge(m, src)
}
func (m *NewAccountBlockBytes) XXX_Size() int {
	return xxx_messageInfo_NewAccountBlockBytes.Size(m)
}
func (m *NewAccountBlockBytes) XXX_DiscardUnknown() {
	xxx_messageInfo_NewAccountBlockBytes.DiscardUnknown(m)
}

var xxx_messageInfo_NewAccountBlockBytes proto.InternalMessageInfo

func (m *NewAccountBlockBytes) GetBlock() []byte {
	if m != nil {
		return m.Block
	}
	return nil
}

func (m *NewAccountBlockBytes) GetTTL() int32 {
	if m != nil {
		return m.TTL
	}
	return 0
}

type Trace struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Path                 [][]byte `protobuf:"bytes,2,rep,name=Path,proto3" json:"Path,omitempty"`
	TTL                  uint32   `protobuf:"varint,3,opt,name=TTL,proto3" json:"TTL,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Trace) Reset()         { *m = Trace{} }
func (m *Trace) String() string { return proto.CompactTextString(m) }
func (*Trace) ProtoMessage()    {}
func (*Trace) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{16}
}

func (m *Trace) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Trace.Unmarshal(m, b)
}
func (m *Trace) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Trace.Marshal(b, m, deterministic)
}
func (m *Trace) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Trace.Merge(m, src)
}
func (m *Trace) XXX_Size() int {
	return xxx_messageInfo_Trace.Size(m)
}
func (m *Trace) XXX_DiscardUnknown() {
	xxx_messageInfo_Trace.DiscardUnknown(m)
}

var xxx_messageInfo_Trace proto.InternalMessageInfo

func (m *Trace) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *Trace) GetPath() [][]byte {
	if m != nil {
		return m.Path
	}
	return nil
}

func (m *Trace) GetTTL() uint32 {
	if m != nil {
		return m.TTL
	}
	return 0
}

func init() {
	proto.RegisterType((*Handshake)(nil), "vitepb.Handshake")
	proto.RegisterType((*SyncConnHandshake)(nil), "vitepb.SyncConnHandshake")
	proto.RegisterType((*ChunkRequest)(nil), "vitepb.ChunkRequest")
	proto.RegisterType((*ChunkResponse)(nil), "vitepb.ChunkResponse")
	proto.RegisterType((*State)(nil), "vitepb.State")
	proto.RegisterType((*State_Peer)(nil), "vitepb.State.Peer")
	proto.RegisterType((*HashHeight)(nil), "vitepb.HashHeight")
	proto.RegisterType((*HashHeightPoint)(nil), "vitepb.HashHeightPoint")
	proto.RegisterType((*HashHeightList)(nil), "vitepb.HashHeightList")
	proto.RegisterType((*GetHashHeightList)(nil), "vitepb.GetHashHeightList")
	proto.RegisterType((*GetSnapshotBlocks)(nil), "vitepb.GetSnapshotBlocks")
	proto.RegisterType((*SnapshotBlocks)(nil), "vitepb.SnapshotBlocks")
	proto.RegisterType((*GetAccountBlocks)(nil), "vitepb.GetAccountBlocks")
	proto.RegisterType((*AccountBlocks)(nil), "vitepb.AccountBlocks")
	proto.RegisterType((*NewSnapshotBlock)(nil), "vitepb.NewSnapshotBlock")
	proto.RegisterType((*NewAccountBlock)(nil), "vitepb.NewAccountBlock")
	proto.RegisterType((*NewAccountBlockBytes)(nil), "vitepb.NewAccountBlockBytes")
	proto.RegisterType((*Trace)(nil), "vitepb.Trace")
}

func init() { proto.RegisterFile("vitepb/message.proto", fileDescriptor_2a6a8486deb9ab39) }

var fileDescriptor_2a6a8486deb9ab39 = []byte{
	// 739 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0x49, 0x6f, 0x13, 0x4b,
	0x10, 0xd6, 0x6c, 0x8e, 0x5d, 0x59, 0x9e, 0xd3, 0xf2, 0x7b, 0x6f, 0x64, 0x38, 0x58, 0x73, 0x40,
	0x16, 0x8b, 0x23, 0xc1, 0x85, 0x0b, 0x20, 0x27, 0x21, 0x71, 0x44, 0x64, 0x4c, 0xdb, 0xe2, 0x1a,
	0x4d, 0x3c, 0xa5, 0x78, 0xe4, 0x78, 0xda, 0x4c, 0x77, 0x12, 0x05, 0x89, 0x1b, 0x57, 0x7e, 0x11,
	0x7f, 0x0e, 0xf5, 0x32, 0x9b, 0x63, 0x23, 0x2e, 0xdc, 0xaa, 0xba, 0xaa, 0xeb, 0xab, 0xaf, 0xea,
	0x9b, 0x1e, 0x68, 0xdd, 0xc6, 0x02, 0x97, 0x97, 0x07, 0x0b, 0xe4, 0x3c, 0xbc, 0xc2, 0xde, 0x32,
	0x65, 0x82, 0x91, 0x9a, 0x3e, 0x6d, 0xb7, 0x4d, 0x34, 0x9c, 0x4e, 0xd9, 0x4d, 0x22, 0x2e, 0x2e,
	0xaf, 0xd9, 0x74, 0xae, 0x73, 0xda, 0x8f, 0x4c, 0x8c, 0x27, 0xe1, 0x92, 0xcf, 0x58, 0x25, 0x18,
	0xfc, 0xb0, 0xa1, 0x31, 0x08, 0x93, 0x88, 0xcf, 0xc2, 0x39, 0x12, 0x1f, 0xb6, 0x3e, 0x63, 0xca,
	0x63, 0x96, 0xf8, 0x56, 0xc7, 0xea, 0x3a, 0x34, 0x73, 0x49, 0x0b, 0xbc, 0x21, 0x8a, 0xb3, 0xc8,
	0xb7, 0xd5, 0xb9, 0x76, 0x08, 0x01, 0x77, 0x18, 0x2e, 0xd0, 0x77, 0x3a, 0x56, 0xb7, 0x41, 0x95,
	0x4d, 0xf6, 0xc0, 0x3e, 0x3b, 0xf6, 0xdd, 0x8e, 0xd5, 0xdd, 0xa1, 0xf6, 0xd9, 0x31, 0x79, 0x0c,
	0x8d, 0x49, 0xbc, 0x40, 0x2e, 0xc2, 0xc5, 0xd2, 0xf7, 0xd4, 0xed, 0xe2, 0x40, 0x22, 0x9e, 0x62,
	0x82, 0x3c, 0xe6, 0x7e, 0x4d, 0x5d, 0xc9, 0x5c, 0xf2, 0x1f, 0xd4, 0x06, 0x18, 0x5f, 0xcd, 0x84,
	0xbf, 0xd5, 0xb1, 0xba, 0x2e, 0x35, 0x9e, 0xc4, 0x1c, 0x60, 0x18, 0xf9, 0x75, 0x95, 0xae, 0x6c,
	0xd2, 0x81, 0xed, 0x93, 0xf8, 0x1a, 0xfb, 0x51, 0x94, 0x22, 0xe7, 0x7e, 0x43, 0x85, 0xca, 0x47,
	0xa4, 0x09, 0xce, 0x07, 0xbc, 0xf7, 0x41, 0x45, 0xa4, 0x29, 0x19, 0x4d, 0xd8, 0x1c, 0x13, 0x7f,
	0x5b, 0x9d, 0x69, 0x27, 0x88, 0x61, 0x7f, 0x7c, 0x9f, 0x4c, 0x8f, 0x58, 0x92, 0x14, 0x63, 0xd1,
	0x94, 0xac, 0xf5, 0x94, 0xec, 0x55, 0x4a, 0x06, 0xca, 0x59, 0x03, 0xe5, 0x96, 0xa1, 0x66, 0xb0,
	0x73, 0x34, 0xbb, 0x49, 0xe6, 0x14, 0xbf, 0xdc, 0x20, 0x57, 0xc4, 0x4e, 0x52, 0xb6, 0x50, 0x38,
	0x2e, 0x55, 0xb6, 0x44, 0x9e, 0x30, 0x05, 0xe1, 0x52, 0x7b, 0xc2, 0x48, 0x1b, 0xea, 0xa3, 0x14,
	0x6f, 0x07, 0x21, 0x9f, 0x19, 0x80, 0xdc, 0x97, 0xa3, 0x7c, 0x9f, 0x44, 0x2a, 0xa4, 0x71, 0x32,
	0x37, 0xf8, 0x06, 0xbb, 0x06, 0x89, 0x2f, 0x59, 0xc2, 0xf1, 0xef, 0x41, 0xc9, 0xca, 0xe3, 0xf8,
	0x2b, 0xaa, 0x45, 0xbb, 0x54, 0xd9, 0xc1, 0x4f, 0x0b, 0xbc, 0xb1, 0x08, 0x05, 0x92, 0x2e, 0x78,
	0x23, 0xc4, 0x94, 0xfb, 0x56, 0xc7, 0xe9, 0x6e, 0xbf, 0x24, 0x3d, 0x2d, 0xcd, 0x9e, 0x8a, 0xf6,
	0x64, 0x88, 0xea, 0x04, 0x39, 0xb2, 0x51, 0x28, 0xa6, 0x33, 0xd5, 0x50, 0x9d, 0x6a, 0x27, 0xdf,
	0xbd, 0x53, 0xda, 0x7d, 0xa1, 0x13, 0xb7, 0xa2, 0x93, 0xca, 0x92, 0x60, 0x65, 0x49, 0xed, 0x2e,
	0xb8, 0x12, 0xe8, 0xc1, 0x6a, 0x9b, 0xe0, 0xf4, 0xa3, 0xc8, 0xa0, 0x4a, 0x33, 0x78, 0x0d, 0x20,
	0x99, 0x95, 0xd4, 0x27, 0x69, 0x5b, 0xa6, 0x03, 0xc9, 0xb9, 0xe8, 0xc0, 0x2e, 0x77, 0x10, 0x7c,
	0x84, 0x7f, 0x8a, 0x9b, 0x23, 0x16, 0x27, 0x42, 0x0d, 0x40, 0x1a, 0xea, 0x7e, 0x69, 0x00, 0x45,
	0x1e, 0xd5, 0x09, 0xf9, 0x20, 0xed, 0xd2, 0x20, 0xfb, 0xb0, 0x57, 0x24, 0x9e, 0xc7, 0x5c, 0x90,
	0x03, 0xa8, 0xa9, 0xf4, 0x6c, 0xa2, 0xff, 0x3f, 0x2c, 0xa8, 0xe2, 0xd4, 0xa4, 0x05, 0x17, 0xb0,
	0x7f, 0x8a, 0x62, 0xa5, 0xca, 0x93, 0x5c, 0x0e, 0xce, 0x86, 0xa6, 0xb4, 0x44, 0x64, 0x4f, 0x02,
	0x97, 0x79, 0x4f, 0x02, 0x97, 0x46, 0x36, 0x4e, 0x26, 0x9b, 0x60, 0xae, 0x00, 0xc6, 0xe6, 0xad,
	0x39, 0x94, 0x4f, 0x0d, 0x2f, 0x01, 0x58, 0xbf, 0x05, 0x68, 0x81, 0x77, 0x24, 0xdf, 0x2f, 0x83,
	0xa0, 0x1d, 0xa9, 0xb6, 0x13, 0x96, 0xde, 0x85, 0xa9, 0x5e, 0x7c, 0x9d, 0x66, 0x6e, 0xf0, 0x0e,
	0xf6, 0x56, 0x90, 0x5e, 0x40, 0x4d, 0x5b, 0x86, 0xcc, 0xbf, 0xb9, 0xc4, 0xca, 0x79, 0xd4, 0x24,
	0x05, 0xdf, 0x2d, 0x68, 0x9e, 0xa2, 0xe8, 0xeb, 0x67, 0xd3, 0xd4, 0xf0, 0x61, 0x2b, 0x7b, 0x49,
	0xf4, 0x9a, 0x33, 0x37, 0xe7, 0x61, 0xff, 0x29, 0x0f, 0x67, 0x03, 0x0f, 0xb7, 0xca, 0xe3, 0x0d,
	0xec, 0x56, 0x5b, 0x78, 0xbe, 0x42, 0xa3, 0x95, 0x41, 0x95, 0xd3, 0x72, 0x16, 0x9f, 0xa0, 0x39,
	0xc4, 0xbb, 0x0a, 0x43, 0xf2, 0x0c, 0x3c, 0x65, 0x98, 0x99, 0x6f, 0x98, 0x83, 0xce, 0x91, 0xaa,
	0x9f, 0x4c, 0xce, 0x15, 0x2d, 0x8f, 0x4a, 0x53, 0x6a, 0x77, 0x88, 0x77, 0x65, 0x34, 0xf2, 0xb4,
	0x5a, 0x71, 0x7d, 0x4b, 0x1b, 0x0b, 0xbe, 0x85, 0xd6, 0x4a, 0xc1, 0xc3, 0x7b, 0x81, 0xea, 0x43,
	0x2f, 0xaa, 0xee, 0x6c, 0xbe, 0xdf, 0x07, 0x6f, 0x92, 0x86, 0x53, 0x5c, 0xfb, 0x05, 0x12, 0x70,
	0x47, 0xa1, 0x90, 0x8f, 0x85, 0x23, 0xcf, 0xa4, 0x9d, 0x95, 0x90, 0x1b, 0xd8, 0x55, 0x25, 0x2e,
	0x6b, 0xea, 0x97, 0xf7, 0xea, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6d, 0x8a, 0x8c, 0xed, 0x4b,
	0x07, 0x00, 0x00,
}