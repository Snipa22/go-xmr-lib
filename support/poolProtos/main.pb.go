// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: main.proto

package poolProtos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type POOLTYPE int32

const (
	POOLTYPE_PPLNS POOLTYPE = 0
	POOLTYPE_PPS   POOLTYPE = 1
	POOLTYPE_PROP  POOLTYPE = 2
	POOLTYPE_SOLO  POOLTYPE = 3
)

// Enum value maps for POOLTYPE.
var (
	POOLTYPE_name = map[int32]string{
		0: "PPLNS",
		1: "PPS",
		2: "PROP",
		3: "SOLO",
	}
	POOLTYPE_value = map[string]int32{
		"PPLNS": 0,
		"PPS":   1,
		"PROP":  2,
		"SOLO":  3,
	}
)

func (x POOLTYPE) Enum() *POOLTYPE {
	p := new(POOLTYPE)
	*p = x
	return p
}

func (x POOLTYPE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (POOLTYPE) Descriptor() protoreflect.EnumDescriptor {
	return file_main_proto_enumTypes[0].Descriptor()
}

func (POOLTYPE) Type() protoreflect.EnumType {
	return &file_main_proto_enumTypes[0]
}

func (x POOLTYPE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *POOLTYPE) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = POOLTYPE(num)
	return nil
}

// Deprecated: Use POOLTYPE.Descriptor instead.
func (POOLTYPE) EnumDescriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{0}
}

type MESSAGETYPE int32

const (
	MESSAGETYPE_SHARE        MESSAGETYPE = 0
	MESSAGETYPE_BLOCK        MESSAGETYPE = 1
	MESSAGETYPE_INVALIDSHARE MESSAGETYPE = 2
)

// Enum value maps for MESSAGETYPE.
var (
	MESSAGETYPE_name = map[int32]string{
		0: "SHARE",
		1: "BLOCK",
		2: "INVALIDSHARE",
	}
	MESSAGETYPE_value = map[string]int32{
		"SHARE":        0,
		"BLOCK":        1,
		"INVALIDSHARE": 2,
	}
)

func (x MESSAGETYPE) Enum() *MESSAGETYPE {
	p := new(MESSAGETYPE)
	*p = x
	return p
}

func (x MESSAGETYPE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MESSAGETYPE) Descriptor() protoreflect.EnumDescriptor {
	return file_main_proto_enumTypes[1].Descriptor()
}

func (MESSAGETYPE) Type() protoreflect.EnumType {
	return &file_main_proto_enumTypes[1]
}

func (x MESSAGETYPE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *MESSAGETYPE) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = MESSAGETYPE(num)
	return nil
}

// Deprecated: Use MESSAGETYPE.Descriptor instead.
func (MESSAGETYPE) EnumDescriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{1}
}

type WSData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgType *MESSAGETYPE `protobuf:"varint,1,req,name=msgType,enum=support.MESSAGETYPE" json:"msgType,omitempty"`
	Key     *string      `protobuf:"bytes,2,req,name=key" json:"key,omitempty"`
	Msg     []byte       `protobuf:"bytes,3,req,name=msg" json:"msg,omitempty"`
	ExInt   *int32       `protobuf:"varint,4,req,name=exInt" json:"exInt,omitempty"`
}

func (x *WSData) Reset() {
	*x = WSData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WSData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WSData) ProtoMessage() {}

func (x *WSData) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WSData.ProtoReflect.Descriptor instead.
func (*WSData) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{0}
}

func (x *WSData) GetMsgType() MESSAGETYPE {
	if x != nil && x.MsgType != nil {
		return *x.MsgType
	}
	return MESSAGETYPE_SHARE
}

func (x *WSData) GetKey() string {
	if x != nil && x.Key != nil {
		return *x.Key
	}
	return ""
}

func (x *WSData) GetMsg() []byte {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *WSData) GetExInt() int32 {
	if x != nil && x.ExInt != nil {
		return *x.ExInt
	}
	return 0
}

type InvalidShare struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PaymentAddress *string `protobuf:"bytes,1,req,name=paymentAddress" json:"paymentAddress,omitempty"`
	PaymentID      *string `protobuf:"bytes,2,opt,name=paymentID" json:"paymentID,omitempty"`
	Identifier     *string `protobuf:"bytes,3,req,name=identifier" json:"identifier,omitempty"`
}

func (x *InvalidShare) Reset() {
	*x = InvalidShare{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvalidShare) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvalidShare) ProtoMessage() {}

func (x *InvalidShare) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvalidShare.ProtoReflect.Descriptor instead.
func (*InvalidShare) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{1}
}

func (x *InvalidShare) GetPaymentAddress() string {
	if x != nil && x.PaymentAddress != nil {
		return *x.PaymentAddress
	}
	return ""
}

func (x *InvalidShare) GetPaymentID() string {
	if x != nil && x.PaymentID != nil {
		return *x.PaymentID
	}
	return ""
}

func (x *InvalidShare) GetIdentifier() string {
	if x != nil && x.Identifier != nil {
		return *x.Identifier
	}
	return ""
}

type Share struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Shares         *int32    `protobuf:"varint,1,req,name=shares" json:"shares,omitempty"`
	PaymentAddress *string   `protobuf:"bytes,2,req,name=paymentAddress" json:"paymentAddress,omitempty"`
	FoundBlock     *bool     `protobuf:"varint,3,req,name=foundBlock" json:"foundBlock,omitempty"`
	PaymentID      *string   `protobuf:"bytes,4,opt,name=paymentID" json:"paymentID,omitempty"`
	TrustedShare   *bool     `protobuf:"varint,5,req,name=trustedShare" json:"trustedShare,omitempty"`
	PoolType       *POOLTYPE `protobuf:"varint,6,req,name=poolType,enum=support.POOLTYPE" json:"poolType,omitempty"`
	PoolID         *int32    `protobuf:"varint,7,req,name=poolID" json:"poolID,omitempty"`
	BlockDiff      *int64    `protobuf:"varint,8,req,name=blockDiff" json:"blockDiff,omitempty"`
	Bitcoin        *bool     `protobuf:"varint,9,req,name=bitcoin" json:"bitcoin,omitempty"`
	BlockHeight    *int32    `protobuf:"varint,10,req,name=blockHeight" json:"blockHeight,omitempty"`
	Timestamp      *int64    `protobuf:"varint,11,req,name=timestamp" json:"timestamp,omitempty"`
	Identifier     *string   `protobuf:"bytes,12,req,name=identifier" json:"identifier,omitempty"`
	BlobData       []byte    `protobuf:"bytes,13,opt,name=blobData" json:"blobData,omitempty"`
	MinerHex       *string   `protobuf:"bytes,14,opt,name=minerHex" json:"minerHex,omitempty"`
	SeedHash       []byte    `protobuf:"bytes,15,opt,name=seedHash" json:"seedHash,omitempty"`
}

func (x *Share) Reset() {
	*x = Share{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Share) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Share) ProtoMessage() {}

func (x *Share) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Share.ProtoReflect.Descriptor instead.
func (*Share) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{2}
}

func (x *Share) GetShares() int32 {
	if x != nil && x.Shares != nil {
		return *x.Shares
	}
	return 0
}

func (x *Share) GetPaymentAddress() string {
	if x != nil && x.PaymentAddress != nil {
		return *x.PaymentAddress
	}
	return ""
}

func (x *Share) GetFoundBlock() bool {
	if x != nil && x.FoundBlock != nil {
		return *x.FoundBlock
	}
	return false
}

func (x *Share) GetPaymentID() string {
	if x != nil && x.PaymentID != nil {
		return *x.PaymentID
	}
	return ""
}

func (x *Share) GetTrustedShare() bool {
	if x != nil && x.TrustedShare != nil {
		return *x.TrustedShare
	}
	return false
}

func (x *Share) GetPoolType() POOLTYPE {
	if x != nil && x.PoolType != nil {
		return *x.PoolType
	}
	return POOLTYPE_PPLNS
}

func (x *Share) GetPoolID() int32 {
	if x != nil && x.PoolID != nil {
		return *x.PoolID
	}
	return 0
}

func (x *Share) GetBlockDiff() int64 {
	if x != nil && x.BlockDiff != nil {
		return *x.BlockDiff
	}
	return 0
}

func (x *Share) GetBitcoin() bool {
	if x != nil && x.Bitcoin != nil {
		return *x.Bitcoin
	}
	return false
}

func (x *Share) GetBlockHeight() int32 {
	if x != nil && x.BlockHeight != nil {
		return *x.BlockHeight
	}
	return 0
}

func (x *Share) GetTimestamp() int64 {
	if x != nil && x.Timestamp != nil {
		return *x.Timestamp
	}
	return 0
}

func (x *Share) GetIdentifier() string {
	if x != nil && x.Identifier != nil {
		return *x.Identifier
	}
	return ""
}

func (x *Share) GetBlobData() []byte {
	if x != nil {
		return x.BlobData
	}
	return nil
}

func (x *Share) GetMinerHex() string {
	if x != nil && x.MinerHex != nil {
		return *x.MinerHex
	}
	return ""
}

func (x *Share) GetSeedHash() []byte {
	if x != nil {
		return x.SeedHash
	}
	return nil
}

type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash       *string   `protobuf:"bytes,1,req,name=hash" json:"hash,omitempty"`
	Difficulty *int64    `protobuf:"varint,2,req,name=difficulty" json:"difficulty,omitempty"`
	Shares     *int64    `protobuf:"varint,3,req,name=shares" json:"shares,omitempty"`
	Timestamp  *int64    `protobuf:"varint,4,req,name=timestamp" json:"timestamp,omitempty"`
	PoolType   *POOLTYPE `protobuf:"varint,5,req,name=poolType,enum=support.POOLTYPE" json:"poolType,omitempty"`
	Unlocked   *bool     `protobuf:"varint,6,req,name=unlocked" json:"unlocked,omitempty"`
	Valid      *bool     `protobuf:"varint,7,req,name=valid" json:"valid,omitempty"`
	Value      *int64    `protobuf:"varint,8,opt,name=value" json:"value,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{3}
}

func (x *Block) GetHash() string {
	if x != nil && x.Hash != nil {
		return *x.Hash
	}
	return ""
}

func (x *Block) GetDifficulty() int64 {
	if x != nil && x.Difficulty != nil {
		return *x.Difficulty
	}
	return 0
}

func (x *Block) GetShares() int64 {
	if x != nil && x.Shares != nil {
		return *x.Shares
	}
	return 0
}

func (x *Block) GetTimestamp() int64 {
	if x != nil && x.Timestamp != nil {
		return *x.Timestamp
	}
	return 0
}

func (x *Block) GetPoolType() POOLTYPE {
	if x != nil && x.PoolType != nil {
		return *x.PoolType
	}
	return POOLTYPE_PPLNS
}

func (x *Block) GetUnlocked() bool {
	if x != nil && x.Unlocked != nil {
		return *x.Unlocked
	}
	return false
}

func (x *Block) GetValid() bool {
	if x != nil && x.Valid != nil {
		return *x.Valid
	}
	return false
}

func (x *Block) GetValue() int64 {
	if x != nil && x.Value != nil {
		return *x.Value
	}
	return 0
}

type FoundBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source     *string `protobuf:"bytes,1,req,name=source" json:"source,omitempty"`
	HashedData []byte  `protobuf:"bytes,2,req,name=hashedData" json:"hashedData,omitempty"`
	Pool       *string `protobuf:"bytes,3,req,name=pool" json:"pool,omitempty"`
}

func (x *FoundBlock) Reset() {
	*x = FoundBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FoundBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FoundBlock) ProtoMessage() {}

func (x *FoundBlock) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FoundBlock.ProtoReflect.Descriptor instead.
func (*FoundBlock) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{4}
}

func (x *FoundBlock) GetSource() string {
	if x != nil && x.Source != nil {
		return *x.Source
	}
	return ""
}

func (x *FoundBlock) GetHashedData() []byte {
	if x != nil {
		return x.HashedData
	}
	return nil
}

func (x *FoundBlock) GetPool() string {
	if x != nil && x.Pool != nil {
		return *x.Pool
	}
	return ""
}

type BlockTemplate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source            *string `protobuf:"bytes,1,req,name=source" json:"source,omitempty"`
	Difficulty        *int64  `protobuf:"varint,2,req,name=difficulty" json:"difficulty,omitempty"`
	Height            *int64  `protobuf:"varint,3,req,name=height" json:"height,omitempty"`
	ReservedOffset    *int64  `protobuf:"varint,4,req,name=reserved_offset,json=reservedOffset" json:"reserved_offset,omitempty"`
	BlocktemplateBlob *string `protobuf:"bytes,5,req,name=blocktemplate_blob,json=blocktemplateBlob" json:"blocktemplate_blob,omitempty"`
	PrevHash          *string `protobuf:"bytes,6,req,name=prev_hash,json=prevHash" json:"prev_hash,omitempty"`
	SeedHash          *string `protobuf:"bytes,7,req,name=seed_hash,json=seedHash" json:"seed_hash,omitempty"`
}

func (x *BlockTemplate) Reset() {
	*x = BlockTemplate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockTemplate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockTemplate) ProtoMessage() {}

func (x *BlockTemplate) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockTemplate.ProtoReflect.Descriptor instead.
func (*BlockTemplate) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{5}
}

func (x *BlockTemplate) GetSource() string {
	if x != nil && x.Source != nil {
		return *x.Source
	}
	return ""
}

func (x *BlockTemplate) GetDifficulty() int64 {
	if x != nil && x.Difficulty != nil {
		return *x.Difficulty
	}
	return 0
}

func (x *BlockTemplate) GetHeight() int64 {
	if x != nil && x.Height != nil {
		return *x.Height
	}
	return 0
}

func (x *BlockTemplate) GetReservedOffset() int64 {
	if x != nil && x.ReservedOffset != nil {
		return *x.ReservedOffset
	}
	return 0
}

func (x *BlockTemplate) GetBlocktemplateBlob() string {
	if x != nil && x.BlocktemplateBlob != nil {
		return *x.BlocktemplateBlob
	}
	return ""
}

func (x *BlockTemplate) GetPrevHash() string {
	if x != nil && x.PrevHash != nil {
		return *x.PrevHash
	}
	return ""
}

func (x *BlockTemplate) GetSeedHash() string {
	if x != nil && x.SeedHash != nil {
		return *x.SeedHash
	}
	return ""
}

type MinerData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source         *int32  `protobuf:"varint,1,req,name=source" json:"source,omitempty"`
	MinerId        []byte  `protobuf:"bytes,2,req,name=miner_id,json=minerId" json:"miner_id,omitempty"`
	Difficulty     *int32  `protobuf:"varint,3,req,name=difficulty" json:"difficulty,omitempty"`
	ConnectionTime *int32  `protobuf:"varint,4,opt,name=connectionTime" json:"connectionTime,omitempty"`
	UserAgent      *string `protobuf:"bytes,5,opt,name=user_agent,json=userAgent" json:"user_agent,omitempty"`
}

func (x *MinerData) Reset() {
	*x = MinerData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinerData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinerData) ProtoMessage() {}

func (x *MinerData) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinerData.ProtoReflect.Descriptor instead.
func (*MinerData) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{6}
}

func (x *MinerData) GetSource() int32 {
	if x != nil && x.Source != nil {
		return *x.Source
	}
	return 0
}

func (x *MinerData) GetMinerId() []byte {
	if x != nil {
		return x.MinerId
	}
	return nil
}

func (x *MinerData) GetDifficulty() int32 {
	if x != nil && x.Difficulty != nil {
		return *x.Difficulty
	}
	return 0
}

func (x *MinerData) GetConnectionTime() int32 {
	if x != nil && x.ConnectionTime != nil {
		return *x.ConnectionTime
	}
	return 0
}

func (x *MinerData) GetUserAgent() string {
	if x != nil && x.UserAgent != nil {
		return *x.UserAgent
	}
	return ""
}

var File_main_proto protoreflect.FileDescriptor

var file_main_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x75,
	0x70, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x72, 0x0a, 0x06, 0x57, 0x53, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x2e, 0x0a, 0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x02, 0x28, 0x0e,
	0x32, 0x14, 0x2e, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x4d, 0x45, 0x53, 0x53, 0x41,
	0x47, 0x45, 0x54, 0x59, 0x50, 0x45, 0x52, 0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x02, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x03, 0x20, 0x02, 0x28, 0x0c, 0x52, 0x03,
	0x6d, 0x73, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x78, 0x49, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x02,
	0x28, 0x05, 0x52, 0x05, 0x65, 0x78, 0x49, 0x6e, 0x74, 0x22, 0x74, 0x0a, 0x0c, 0x49, 0x6e, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x53, 0x68, 0x61, 0x72, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x70, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x02, 0x28,
	0x09, 0x52, 0x0e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12,
	0x1e, 0x0a, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x03, 0x20,
	0x02, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x22,
	0xdc, 0x03, 0x0a, 0x05, 0x53, 0x68, 0x61, 0x72, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x68, 0x61,
	0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x02, 0x28, 0x05, 0x52, 0x06, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x73, 0x12, 0x26, 0x0a, 0x0e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x02, 0x28, 0x09, 0x52, 0x0e, 0x70, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x66, 0x6f, 0x75,
	0x6e, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x03, 0x20, 0x02, 0x28, 0x08, 0x52, 0x0a, 0x66,
	0x6f, 0x75, 0x6e, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x72, 0x75, 0x73, 0x74,
	0x65, 0x64, 0x53, 0x68, 0x61, 0x72, 0x65, 0x18, 0x05, 0x20, 0x02, 0x28, 0x08, 0x52, 0x0c, 0x74,
	0x72, 0x75, 0x73, 0x74, 0x65, 0x64, 0x53, 0x68, 0x61, 0x72, 0x65, 0x12, 0x2d, 0x0a, 0x08, 0x70,
	0x6f, 0x6f, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x02, 0x28, 0x0e, 0x32, 0x11, 0x2e,
	0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x50, 0x4f, 0x4f, 0x4c, 0x54, 0x59, 0x50, 0x45,
	0x52, 0x08, 0x70, 0x6f, 0x6f, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6f,
	0x6f, 0x6c, 0x49, 0x44, 0x18, 0x07, 0x20, 0x02, 0x28, 0x05, 0x52, 0x06, 0x70, 0x6f, 0x6f, 0x6c,
	0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x44, 0x69, 0x66, 0x66, 0x18,
	0x08, 0x20, 0x02, 0x28, 0x03, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x44, 0x69, 0x66, 0x66,
	0x12, 0x18, 0x0a, 0x07, 0x62, 0x69, 0x74, 0x63, 0x6f, 0x69, 0x6e, 0x18, 0x09, 0x20, 0x02, 0x28,
	0x08, 0x52, 0x07, 0x62, 0x69, 0x74, 0x63, 0x6f, 0x69, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x0a, 0x20, 0x02, 0x28, 0x05, 0x52,
	0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x0b, 0x20, 0x02, 0x28, 0x03, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x0c, 0x20, 0x02, 0x28, 0x09, 0x52, 0x0a,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x6c,
	0x6f, 0x62, 0x44, 0x61, 0x74, 0x61, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x62, 0x6c,
	0x6f, 0x62, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x69, 0x6e, 0x65, 0x72, 0x48,
	0x65, 0x78, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x69, 0x6e, 0x65, 0x72, 0x48,
	0x65, 0x78, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x65, 0x64, 0x48, 0x61, 0x73, 0x68, 0x18, 0x0f,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x73, 0x65, 0x65, 0x64, 0x48, 0x61, 0x73, 0x68, 0x22, 0xe8,
	0x01, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x1e, 0x0a, 0x0a,
	0x64, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x18, 0x02, 0x20, 0x02, 0x28, 0x03,
	0x52, 0x0a, 0x64, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x73, 0x18, 0x03, 0x20, 0x02, 0x28, 0x03, 0x52, 0x06, 0x73, 0x68,
	0x61, 0x72, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x04, 0x20, 0x02, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x2d, 0x0a, 0x08, 0x70, 0x6f, 0x6f, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x18, 0x05,
	0x20, 0x02, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x50,
	0x4f, 0x4f, 0x4c, 0x54, 0x59, 0x50, 0x45, 0x52, 0x08, 0x70, 0x6f, 0x6f, 0x6c, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x18, 0x06, 0x20,
	0x02, 0x28, 0x08, 0x52, 0x08, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x18, 0x07, 0x20, 0x02, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x58, 0x0a, 0x0a, 0x46, 0x6f, 0x75,
	0x6e, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x1e, 0x0a, 0x0a, 0x68, 0x61, 0x73, 0x68, 0x65, 0x64, 0x44, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20,
	0x02, 0x28, 0x0c, 0x52, 0x0a, 0x68, 0x61, 0x73, 0x68, 0x65, 0x64, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x6f, 0x6f, 0x6c, 0x18, 0x03, 0x20, 0x02, 0x28, 0x09, 0x52, 0x04, 0x70,
	0x6f, 0x6f, 0x6c, 0x22, 0xf1, 0x01, 0x0a, 0x0d, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x65, 0x6d,
	0x70, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18,
	0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x1e, 0x0a,
	0x0a, 0x64, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x18, 0x02, 0x20, 0x02, 0x28,
	0x03, 0x52, 0x0a, 0x64, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x12, 0x16, 0x0a,
	0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x02, 0x28, 0x03, 0x52, 0x06, 0x68,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x64, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x04, 0x20, 0x02, 0x28, 0x03, 0x52, 0x0e,
	0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x2d,
	0x0a, 0x12, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f,
	0x62, 0x6c, 0x6f, 0x62, 0x18, 0x05, 0x20, 0x02, 0x28, 0x09, 0x52, 0x11, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x6c, 0x6f, 0x62, 0x12, 0x1b, 0x0a,
	0x09, 0x70, 0x72, 0x65, 0x76, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x06, 0x20, 0x02, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x72, 0x65, 0x76, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65,
	0x65, 0x64, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x07, 0x20, 0x02, 0x28, 0x09, 0x52, 0x08, 0x73,
	0x65, 0x65, 0x64, 0x48, 0x61, 0x73, 0x68, 0x22, 0xa5, 0x01, 0x0a, 0x09, 0x4d, 0x69, 0x6e, 0x65,
	0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18,
	0x01, 0x20, 0x02, 0x28, 0x05, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x19, 0x0a,
	0x08, 0x6d, 0x69, 0x6e, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x02, 0x28, 0x0c, 0x52,
	0x07, 0x6d, 0x69, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x69, 0x66, 0x66,
	0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x18, 0x03, 0x20, 0x02, 0x28, 0x05, 0x52, 0x0a, 0x64, 0x69,
	0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x12, 0x26, 0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x75, 0x73, 0x65, 0x72, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x2a,
	0x32, 0x0a, 0x08, 0x50, 0x4f, 0x4f, 0x4c, 0x54, 0x59, 0x50, 0x45, 0x12, 0x09, 0x0a, 0x05, 0x50,
	0x50, 0x4c, 0x4e, 0x53, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x50, 0x50, 0x53, 0x10, 0x01, 0x12,
	0x08, 0x0a, 0x04, 0x50, 0x52, 0x4f, 0x50, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x53, 0x4f, 0x4c,
	0x4f, 0x10, 0x03, 0x2a, 0x35, 0x0a, 0x0b, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x54, 0x59,
	0x50, 0x45, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x48, 0x41, 0x52, 0x45, 0x10, 0x00, 0x12, 0x09, 0x0a,
	0x05, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x49, 0x4e, 0x56, 0x41,
	0x4c, 0x49, 0x44, 0x53, 0x48, 0x41, 0x52, 0x45, 0x10, 0x02, 0x42, 0x10, 0x5a, 0x0e, 0x67, 0x69,
	0x74, 0x2e, 0x6a, 0x61, 0x67, 0x74, 0x65, 0x63, 0x68, 0x2e, 0x69, 0x6f,
}

var (
	file_main_proto_rawDescOnce sync.Once
	file_main_proto_rawDescData = file_main_proto_rawDesc
)

func file_main_proto_rawDescGZIP() []byte {
	file_main_proto_rawDescOnce.Do(func() {
		file_main_proto_rawDescData = protoimpl.X.CompressGZIP(file_main_proto_rawDescData)
	})
	return file_main_proto_rawDescData
}

var file_main_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_main_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_main_proto_goTypes = []interface{}{
	(POOLTYPE)(0),         // 0: support.POOLTYPE
	(MESSAGETYPE)(0),      // 1: support.MESSAGETYPE
	(*WSData)(nil),        // 2: support.WSData
	(*InvalidShare)(nil),  // 3: support.InvalidShare
	(*Share)(nil),         // 4: support.Share
	(*Block)(nil),         // 5: support.Block
	(*FoundBlock)(nil),    // 6: support.FoundBlock
	(*BlockTemplate)(nil), // 7: support.BlockTemplate
	(*MinerData)(nil),     // 8: support.MinerData
}
var file_main_proto_depIdxs = []int32{
	1, // 0: support.WSData.msgType:type_name -> support.MESSAGETYPE
	0, // 1: support.Share.poolType:type_name -> support.POOLTYPE
	0, // 2: support.Block.poolType:type_name -> support.POOLTYPE
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_main_proto_init() }
func file_main_proto_init() {
	if File_main_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_main_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WSData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvalidShare); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Share); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FoundBlock); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockTemplate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinerData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_main_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_main_proto_goTypes,
		DependencyIndexes: file_main_proto_depIdxs,
		EnumInfos:         file_main_proto_enumTypes,
		MessageInfos:      file_main_proto_msgTypes,
	}.Build()
	File_main_proto = out.File
	file_main_proto_rawDesc = nil
	file_main_proto_goTypes = nil
	file_main_proto_depIdxs = nil
}
