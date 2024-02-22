// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: multistaking/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	types "github.com/cosmos/cosmos-sdk/x/staking/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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

type GenesisState struct {
	MultiStakingLocks          []MultiStakingLock          `protobuf:"bytes,1,rep,name=multi_staking_locks,json=multiStakingLocks,proto3" json:"multi_staking_locks"`
	MultiStakingUnlocks        []MultiStakingUnlock        `protobuf:"bytes,2,rep,name=multi_staking_unlocks,json=multiStakingUnlocks,proto3" json:"multi_staking_unlocks"`
	MultiStakingCoinInfo       []MultiStakingCoinInfo      `protobuf:"bytes,3,rep,name=multi_staking_coin_info,json=multiStakingCoinInfo,proto3" json:"multi_staking_coin_info"`
	ValidatorMultiStakingCoins []ValidatorMultiStakingCoin `protobuf:"bytes,4,rep,name=validator_multi_staking_coins,json=validatorMultiStakingCoins,proto3" json:"validator_multi_staking_coins"`
	StakingGenesisState        types.GenesisState          `protobuf:"bytes,5,opt,name=staking_genesis_state,json=stakingGenesisState,proto3" json:"staking_genesis_state"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f95a201ebed173c, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetMultiStakingLocks() []MultiStakingLock {
	if m != nil {
		return m.MultiStakingLocks
	}
	return nil
}

func (m *GenesisState) GetMultiStakingUnlocks() []MultiStakingUnlock {
	if m != nil {
		return m.MultiStakingUnlocks
	}
	return nil
}

func (m *GenesisState) GetMultiStakingCoinInfo() []MultiStakingCoinInfo {
	if m != nil {
		return m.MultiStakingCoinInfo
	}
	return nil
}

func (m *GenesisState) GetValidatorMultiStakingCoins() []ValidatorMultiStakingCoin {
	if m != nil {
		return m.ValidatorMultiStakingCoins
	}
	return nil
}

func (m *GenesisState) GetStakingGenesisState() types.GenesisState {
	if m != nil {
		return m.StakingGenesisState
	}
	return types.GenesisState{}
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "multistaking.v1.GenesisState")
}

func init() { proto.RegisterFile("multistaking/v1/genesis.proto", fileDescriptor_8f95a201ebed173c) }

var fileDescriptor_8f95a201ebed173c = []byte{
	// 425 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0x41, 0x6f, 0xd3, 0x30,
	0x1c, 0xc5, 0x13, 0x3a, 0x38, 0x64, 0x48, 0x88, 0x6c, 0xd3, 0xb2, 0x48, 0xcb, 0x06, 0x1b, 0xd2,
	0x84, 0x44, 0xac, 0x8c, 0x13, 0x57, 0x38, 0x20, 0x24, 0xb8, 0x50, 0x01, 0x12, 0x12, 0x44, 0x4e,
	0x70, 0x53, 0xab, 0xb1, 0xff, 0x55, 0xff, 0x4e, 0x44, 0xbf, 0x05, 0x1f, 0xab, 0xc7, 0x1e, 0x39,
	0x21, 0xd4, 0x7e, 0x10, 0x50, 0x6c, 0x47, 0x6d, 0x52, 0xe8, 0xcd, 0xcf, 0xef, 0xe5, 0xfd, 0x62,
	0xfb, 0xef, 0x9d, 0x8b, 0xaa, 0x54, 0x1c, 0x15, 0x9d, 0x70, 0x59, 0x90, 0x3a, 0x21, 0x05, 0x93,
	0x0c, 0x39, 0xc6, 0xd3, 0x19, 0x28, 0xf0, 0x1f, 0x6c, 0xdb, 0x71, 0x9d, 0x84, 0x67, 0x05, 0x40,
	0x51, 0x32, 0xa2, 0xed, 0xac, 0x1a, 0x11, 0x2a, 0xe7, 0x26, 0x1b, 0x5e, 0xf4, 0x2d, 0xc5, 0x05,
	0x43, 0x45, 0xc5, 0xd4, 0x06, 0x8e, 0x0b, 0x28, 0x40, 0x2f, 0x49, 0xb3, 0xb2, 0xbb, 0x67, 0x39,
	0xa0, 0x00, 0x4c, 0x8d, 0x61, 0x84, 0xb5, 0x22, 0xa3, 0x48, 0x46, 0x91, 0x91, 0x3a, 0xc9, 0x98,
	0xa2, 0x09, 0xc9, 0x81, 0x4b, 0xeb, 0x5f, 0x5b, 0x7f, 0xf3, 0xfb, 0x26, 0xd2, 0x39, 0x43, 0x78,
	0xd5, 0x3f, 0xa2, 0xd6, 0x69, 0x7b, 0x28, 0x13, 0x3a, 0xb5, 0x55, 0x02, 0x4d, 0x04, 0xad, 0xf1,
	0xf8, 0xcf, 0xc0, 0xbb, 0xff, 0xda, 0xf4, 0x0d, 0x15, 0x55, 0xcc, 0xff, 0xe4, 0x1d, 0x75, 0x0a,
	0xd2, 0x12, 0xf2, 0x09, 0x06, 0xee, 0xe5, 0xe0, 0xe6, 0xf0, 0xf6, 0x51, 0xdc, 0xbb, 0xb0, 0xf8,
	0x5d, 0xa3, 0x87, 0x46, 0xbf, 0x85, 0x7c, 0xf2, 0xf2, 0x60, 0xf1, 0xeb, 0xc2, 0x79, 0xff, 0x50,
	0xf4, 0xf6, 0xd1, 0xff, 0xe2, 0x9d, 0x74, 0x8b, 0x2b, 0x69, 0xaa, 0xef, 0xe8, 0xea, 0xab, 0xbd,
	0xd5, 0x1f, 0x74, 0xd6, 0x96, 0x1f, 0x89, 0x1d, 0x07, 0xfd, 0xcc, 0x3b, 0xed, 0xd6, 0x37, 0x17,
	0x99, 0x72, 0x39, 0x82, 0x60, 0xa0, 0x01, 0x4f, 0xf6, 0x02, 0x5e, 0x01, 0x97, 0x6f, 0xe4, 0x08,
	0x2c, 0xe2, 0x58, 0xfc, 0xc3, 0xf3, 0xd1, 0x3b, 0xaf, 0x69, 0xc9, 0xbf, 0x51, 0x05, 0xb3, 0x74,
	0x97, 0x86, 0xc1, 0x81, 0x26, 0x3d, 0xdd, 0x21, 0x7d, 0x6c, 0xbf, 0xea, 0x23, 0x2d, 0x2e, 0xac,
	0xff, 0x17, 0x40, 0xff, 0xab, 0x77, 0xd2, 0x42, 0xec, 0xc3, 0x37, 0x50, 0xc5, 0x82, 0xbb, 0x97,
	0xee, 0xcd, 0xe1, 0xed, 0x75, 0x6c, 0x67, 0x6a, 0x83, 0xd3, 0x53, 0x12, 0x6f, 0xbf, 0x6a, 0x7b,
	0x71, 0x36, 0xd3, 0xb1, 0x86, 0x8b, 0x55, 0xe4, 0x2e, 0x57, 0x91, 0xfb, 0x7b, 0x15, 0xb9, 0x3f,
	0xd6, 0x91, 0xb3, 0x5c, 0x47, 0xce, 0xcf, 0x75, 0xe4, 0x7c, 0x7e, 0x51, 0x70, 0x35, 0xae, 0xb2,
	0x38, 0x07, 0x41, 0x66, 0x8c, 0x96, 0x1c, 0x14, 0xcb, 0xc7, 0x66, 0xbe, 0x9e, 0xb5, 0x03, 0xf7,
	0xbd, 0xa7, 0xd5, 0x7c, 0xca, 0x30, 0xbb, 0xa7, 0xa7, 0xeb, 0xf9, 0xdf, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x02, 0xcc, 0xb7, 0xf3, 0x80, 0x03, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.StakingGenesisState.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if len(m.ValidatorMultiStakingCoins) > 0 {
		for iNdEx := len(m.ValidatorMultiStakingCoins) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ValidatorMultiStakingCoins[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.MultiStakingCoinInfo) > 0 {
		for iNdEx := len(m.MultiStakingCoinInfo) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MultiStakingCoinInfo[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.MultiStakingUnlocks) > 0 {
		for iNdEx := len(m.MultiStakingUnlocks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MultiStakingUnlocks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.MultiStakingLocks) > 0 {
		for iNdEx := len(m.MultiStakingLocks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MultiStakingLocks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.MultiStakingLocks) > 0 {
		for _, e := range m.MultiStakingLocks {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.MultiStakingUnlocks) > 0 {
		for _, e := range m.MultiStakingUnlocks {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.MultiStakingCoinInfo) > 0 {
		for _, e := range m.MultiStakingCoinInfo {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.ValidatorMultiStakingCoins) > 0 {
		for _, e := range m.ValidatorMultiStakingCoins {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.StakingGenesisState.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MultiStakingLocks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MultiStakingLocks = append(m.MultiStakingLocks, MultiStakingLock{})
			if err := m.MultiStakingLocks[len(m.MultiStakingLocks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MultiStakingUnlocks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MultiStakingUnlocks = append(m.MultiStakingUnlocks, MultiStakingUnlock{})
			if err := m.MultiStakingUnlocks[len(m.MultiStakingUnlocks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MultiStakingCoinInfo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MultiStakingCoinInfo = append(m.MultiStakingCoinInfo, MultiStakingCoinInfo{})
			if err := m.MultiStakingCoinInfo[len(m.MultiStakingCoinInfo)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorMultiStakingCoins", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorMultiStakingCoins = append(m.ValidatorMultiStakingCoins, ValidatorMultiStakingCoin{})
			if err := m.ValidatorMultiStakingCoins[len(m.ValidatorMultiStakingCoins)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StakingGenesisState", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.StakingGenesisState.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
