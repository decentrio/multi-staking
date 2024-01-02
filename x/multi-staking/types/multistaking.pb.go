// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: multistaking/v1/multistaking.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// AddBondTokenProposal is a gov v1beta1 Content type to add a token as a bond
// denom
type AddBondTokenProposal struct {
	// title of the proposal
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// description of the proposal
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// the token
	BondToken string `protobuf:"bytes,3,opt,name=bond_token,json=bondToken,proto3" json:"bond_token,omitempty"`
	// the bond token weight
	TokenWeight *cosmossdk_io_math.LegacyDec `protobuf:"bytes,4,opt,name=token_weight,json=tokenWeight,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"token_weight,omitempty"`
}

func (m *AddBondTokenProposal) Reset()      { *m = AddBondTokenProposal{} }
func (*AddBondTokenProposal) ProtoMessage() {}
func (*AddBondTokenProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f12af1dde3773b8, []int{0}
}
func (m *AddBondTokenProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AddBondTokenProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AddBondTokenProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AddBondTokenProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddBondTokenProposal.Merge(m, src)
}
func (m *AddBondTokenProposal) XXX_Size() int {
	return m.Size()
}
func (m *AddBondTokenProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_AddBondTokenProposal.DiscardUnknown(m)
}

var xxx_messageInfo_AddBondTokenProposal proto.InternalMessageInfo

// ChangeBondTokenWeightProposal is a gov v1beta1 Content type to change the
// weight of bond token
type ChangeBondTokenWeightProposal struct {
	// title of the proposal
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// description of the proposal
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// the token
	BondToken string `protobuf:"bytes,3,opt,name=bond_token,json=bondToken,proto3" json:"bond_token,omitempty"`
	// the bond token weight
	TokenWeight *cosmossdk_io_math.LegacyDec `protobuf:"bytes,4,opt,name=token_weight,json=tokenWeight,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"token_weight,omitempty"`
}

func (m *ChangeBondTokenWeightProposal) Reset()      { *m = ChangeBondTokenWeightProposal{} }
func (*ChangeBondTokenWeightProposal) ProtoMessage() {}
func (*ChangeBondTokenWeightProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f12af1dde3773b8, []int{1}
}
func (m *ChangeBondTokenWeightProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChangeBondTokenWeightProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChangeBondTokenWeightProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChangeBondTokenWeightProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangeBondTokenWeightProposal.Merge(m, src)
}
func (m *ChangeBondTokenWeightProposal) XXX_Size() int {
	return m.Size()
}
func (m *ChangeBondTokenWeightProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangeBondTokenWeightProposal.DiscardUnknown(m)
}

var xxx_messageInfo_ChangeBondTokenWeightProposal proto.InternalMessageInfo

func init() {
	proto.RegisterType((*AddBondTokenProposal)(nil), "multistaking.v1.AddBondTokenProposal")
	proto.RegisterType((*ChangeBondTokenWeightProposal)(nil), "multistaking.v1.ChangeBondTokenWeightProposal")
}

func init() {
	proto.RegisterFile("multistaking/v1/multistaking.proto", fileDescriptor_1f12af1dde3773b8)
}

var fileDescriptor_1f12af1dde3773b8 = []byte{
	// 356 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x92, 0x3d, 0x4b, 0xfb, 0x50,
	0x14, 0xc6, 0x93, 0xff, 0x8b, 0xd0, 0x5b, 0x41, 0x08, 0x1d, 0x62, 0xa5, 0x49, 0xe9, 0x24, 0x42,
	0x12, 0x82, 0x93, 0x6e, 0xb6, 0x1d, 0x1d, 0x4a, 0x15, 0x04, 0x97, 0x92, 0x97, 0xcb, 0xcd, 0xa5,
	0xc9, 0x3d, 0x21, 0x39, 0x8d, 0xf6, 0x1b, 0x38, 0x3a, 0x3a, 0xf6, 0x43, 0xf4, 0x43, 0x88, 0x53,
	0x71, 0x12, 0x07, 0x91, 0x76, 0x11, 0x5c, 0xfd, 0x00, 0xd2, 0xdc, 0x56, 0xaa, 0x1f, 0xc1, 0xed,
	0x3c, 0xe7, 0x79, 0xf8, 0x71, 0x0e, 0x3c, 0xa4, 0x95, 0x8c, 0x62, 0xe4, 0x39, 0x7a, 0x43, 0x2e,
	0x98, 0x53, 0xb8, 0xce, 0xa6, 0xb6, 0xd3, 0x0c, 0x10, 0xb4, 0x9d, 0x6f, 0xbb, 0xc2, 0xad, 0xd7,
	0x18, 0x30, 0x28, 0x3d, 0x67, 0x39, 0xc9, 0x58, 0x7d, 0x37, 0x80, 0x3c, 0x81, 0x7c, 0x20, 0x0d,
	0x29, 0xa4, 0xd5, 0x7a, 0x57, 0x49, 0xed, 0x24, 0x0c, 0xdb, 0x20, 0xc2, 0x73, 0x18, 0x52, 0xd1,
	0xcb, 0x20, 0x85, 0xdc, 0x8b, 0xb5, 0x1a, 0xf9, 0x8f, 0x1c, 0x63, 0xaa, 0xab, 0x4d, 0x75, 0xbf,
	0xd2, 0x97, 0x42, 0x6b, 0x92, 0x6a, 0x48, 0xf3, 0x20, 0xe3, 0x29, 0x72, 0x10, 0xfa, 0x9f, 0xd2,
	0xdb, 0x5c, 0x69, 0x0d, 0x42, 0x7c, 0x10, 0xe1, 0x00, 0x97, 0x34, 0xfd, 0x6f, 0x19, 0xa8, 0xf8,
	0x6b, 0xbc, 0xd6, 0x23, 0xdb, 0xa5, 0x33, 0xb8, 0xa2, 0x9c, 0x45, 0xa8, 0xff, 0x5b, 0x06, 0xda,
	0xd6, 0xf3, 0x8b, 0xb9, 0x27, 0xef, 0xca, 0xc3, 0xa1, 0xcd, 0xc1, 0x49, 0x3c, 0x8c, 0xec, 0x53,
	0xca, 0xbc, 0x60, 0xdc, 0xa5, 0xc1, 0xe3, 0xd4, 0x22, 0xab, 0xb3, 0xbb, 0x34, 0xe8, 0x57, 0x4b,
	0xc4, 0x45, 0x49, 0x38, 0x3e, 0xb8, 0x99, 0x98, 0xca, 0xdd, 0xc4, 0x54, 0xde, 0x26, 0xa6, 0xf2,
	0x30, 0xb5, 0xea, 0xab, 0x20, 0x83, 0xc2, 0x2e, 0x5c, 0x9f, 0xa2, 0xe7, 0xda, 0x1d, 0x10, 0x48,
	0x05, 0xb6, 0x3e, 0x54, 0xd2, 0xe8, 0x44, 0x9e, 0x60, 0xf4, 0xeb, 0x61, 0x49, 0xf9, 0xd5, 0x6f,
	0xb7, 0xcf, 0xee, 0xe7, 0x86, 0x3a, 0x9b, 0x1b, 0xea, 0xeb, 0xdc, 0x50, 0x6f, 0x17, 0x86, 0x32,
	0x5b, 0x18, 0xca, 0xd3, 0xc2, 0x50, 0x2e, 0x8f, 0x18, 0xc7, 0x68, 0xe4, 0xdb, 0x01, 0x24, 0x4e,
	0x46, 0xbd, 0x98, 0x03, 0xd2, 0x20, 0x92, 0x55, 0xb3, 0xd6, 0xdd, 0xbb, 0xfe, 0xa1, 0x71, 0x9c,
	0xd2, 0xdc, 0xdf, 0x2a, 0x0b, 0x74, 0xf8, 0x19, 0x00, 0x00, 0xff, 0xff, 0xe4, 0xe2, 0x38, 0xb4,
	0xa8, 0x02, 0x00, 0x00,
}

func (m *AddBondTokenProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AddBondTokenProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AddBondTokenProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.TokenWeight != nil {
		{
			size := m.TokenWeight.Size()
			i -= size
			if _, err := m.TokenWeight.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintMultistaking(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.BondToken) > 0 {
		i -= len(m.BondToken)
		copy(dAtA[i:], m.BondToken)
		i = encodeVarintMultistaking(dAtA, i, uint64(len(m.BondToken)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintMultistaking(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintMultistaking(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ChangeBondTokenWeightProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChangeBondTokenWeightProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChangeBondTokenWeightProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.TokenWeight != nil {
		{
			size := m.TokenWeight.Size()
			i -= size
			if _, err := m.TokenWeight.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintMultistaking(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.BondToken) > 0 {
		i -= len(m.BondToken)
		copy(dAtA[i:], m.BondToken)
		i = encodeVarintMultistaking(dAtA, i, uint64(len(m.BondToken)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintMultistaking(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintMultistaking(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMultistaking(dAtA []byte, offset int, v uint64) int {
	offset -= sovMultistaking(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *AddBondTokenProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovMultistaking(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovMultistaking(uint64(l))
	}
	l = len(m.BondToken)
	if l > 0 {
		n += 1 + l + sovMultistaking(uint64(l))
	}
	if m.TokenWeight != nil {
		l = m.TokenWeight.Size()
		n += 1 + l + sovMultistaking(uint64(l))
	}
	return n
}

func (m *ChangeBondTokenWeightProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovMultistaking(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovMultistaking(uint64(l))
	}
	l = len(m.BondToken)
	if l > 0 {
		n += 1 + l + sovMultistaking(uint64(l))
	}
	if m.TokenWeight != nil {
		l = m.TokenWeight.Size()
		n += 1 + l + sovMultistaking(uint64(l))
	}
	return n
}

func sovMultistaking(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMultistaking(x uint64) (n int) {
	return sovMultistaking(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AddBondTokenProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMultistaking
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AddBondTokenProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AddBondTokenProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMultistaking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMultistaking
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMultistaking
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMultistaking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMultistaking
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMultistaking
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BondToken", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMultistaking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMultistaking
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMultistaking
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BondToken = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenWeight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMultistaking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMultistaking
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMultistaking
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v cosmossdk_io_math.LegacyDec
			m.TokenWeight = &v
			if err := m.TokenWeight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMultistaking(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMultistaking
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ChangeBondTokenWeightProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMultistaking
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ChangeBondTokenWeightProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChangeBondTokenWeightProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMultistaking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMultistaking
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMultistaking
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMultistaking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMultistaking
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMultistaking
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BondToken", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMultistaking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMultistaking
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMultistaking
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BondToken = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenWeight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMultistaking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMultistaking
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMultistaking
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v cosmossdk_io_math.LegacyDec
			m.TokenWeight = &v
			if err := m.TokenWeight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMultistaking(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMultistaking
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMultistaking(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMultistaking
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMultistaking
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMultistaking
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthMultistaking
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMultistaking
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMultistaking
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMultistaking        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMultistaking          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMultistaking = fmt.Errorf("proto: unexpected end of group")
)
