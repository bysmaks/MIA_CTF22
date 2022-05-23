package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	// CommandCodeMemoryAreaRead Command code: IO memory area read
	CommandCodeMemoryAreaRead uint16 = 0x0101

	// CommandCodeMemoryAreaWrite Command code: IO memory area write
	CommandCodeMemoryAreaWrite uint16 = 0x0102

	// CommandCodeMemoryAreaFill Command code: IO memory area fill
	CommandCodeMemoryAreaFill uint16 = 0x0103

	// CommandCodeMultipleMemoryAreaRead Command code: IO memory area multiple read
	CommandCodeMultipleMemoryAreaRead uint16 = 0x0104

	// CommandCodeMemoryAreaTransfer Command code: IO memory area transfer
	CommandCodeMemoryAreaTransfer uint16 = 0x0105

	// CommandCodeParameterAreaRead Command code: Parameter area read
	CommandCodeParameterAreaRead uint16 = 0x0201

	// CommandCodeParameterAreaWrite Command code: Parameter area write
	CommandCodeParameterAreaWrite uint16 = 0x0202

	// CommandCodeParameterAreaClear Command code: Parameter area clear
	CommandCodeParameterAreaClear uint16 = 0x0203

	// CommandCodeProgramAreaRead Command code: Program area read
	CommandCodeProgramAreaRead uint16 = 0x0301

	// CommandCodeProgramAreaWrite Command code: Program area write
	CommandCodeProgramAreaWrite uint16 = 0x0302

	// CommandCodeProgramAreaClear Command code: Program area clear
	CommandCodeProgramAreaClear uint16 = 0x0303

	// CommandCodeRun Command code: Set operating mode to run
	CommandCodeRun uint16 = 0x0401

	// CommandCodeStop Command code: Set operating mode to stop
	CommandCodeStop uint16 = 0x0402

	// CommandCodeCPUUnitDataRead Command code: CPU unit data read
	CommandCodeCPUUnitDataRead uint16 = 0x0501

	// CommandCodeConnectionDataRead Command code: connection data read
	CommandCodeConnectionDataRead uint16 = 0x0502

	// CommandCodeCPUUnitStatusRead Command code: CPU unit status read
	CommandCodeCPUUnitStatusRead uint16 = 0x0601

	// CommandCodeCycleTimeRead Command code: cycle time read
	CommandCodeCycleTimeRead uint16 = 0x0620

	// CommandCodeClockRead Command code: clock read
	CommandCodeClockRead uint16 = 0x701

	// CommandCodeClockWrite Command code: clock write
	CommandCodeClockWrite uint16 = 0x702

	// CommandCodeMessageReadClear Command code: message read/clear
	CommandCodeMessageReadClear uint16 = 0x0920

	// CommandCodeAccessRightAcquire Command code: access right acquire
	CommandCodeAccessRightAcquire uint16 = 0x0c01

	// CommandCodeAccessRightForcedAcquire Command code: accress right forced acquire
	CommandCodeAccessRightForcedAcquire uint16 = 0x0c02

	// CommandCodeAccessRightRelease Command code: access right release
	CommandCodeAccessRightRelease uint16 = 0x0c03

	// CommandCodeErrorClear Command code: error clear
	CommandCodeErrorClear uint16 = 0x2101

	// CommandCodeErrorLogRead Command code: error log read
	CommandCodeErrorLogRead uint16 = 0x2102

	// CommandCodeErrorLogClear Command code: error log clear
	CommandCodeErrorLogClear uint16 = 0x2103

	// CommandCodeFINSWriteAccessLogRead Command code: FINS write access log read
	CommandCodeFINSWriteAccessLogRead uint16 = 0x2140

	// CommandCodeFINSWriteAccessLogWrite Command code: FINS write access log write
	CommandCodeFINSWriteAccessLogWrite uint16 = 0x2141

	// CommandCodeFileNameRead Command code: file name read
	CommandCodeFileNameRead uint16 = 0x2101

	// CommandCodeSingleFileRead Command code: file read
	CommandCodeSingleFileRead uint16 = 0x2102

	// CommandCodeSingleFileWrite Command code: file write
	CommandCodeSingleFileWrite uint16 = 0x2103

	// CommandCodeFileMemoryFormat Command code: file memory format
	CommandCodeFileMemoryFormat uint16 = 0x2104

	// CommandCodeFileDelete Command code: file delete
	CommandCodeFileDelete uint16 = 0x2105

	// CommandCodeFileCopy Command code: file copy
	CommandCodeFileCopy uint16 = 0x2107

	// CommandCodeFileNameChange Command code: file name change
	CommandCodeFileNameChange uint16 = 0x2108

	// CommandCodeMemoryAreaFileTransfer Command code: memory area file transfer
	CommandCodeMemoryAreaFileTransfer uint16 = 0x210a

	// CommandCodeParameterAreaFileTransfer Command code: parameter area file transfer
	CommandCodeParameterAreaFileTransfer uint16 = 0x210b

	// CommandCodeProgramAreaFileTransfer Command code: program area file transfer
	CommandCodeProgramAreaFileTransfer uint16 = 0x210b

	// CommandCodeDirectoryCreateDelete Command code: directory create/delete
	CommandCodeDirectoryCreateDelete uint16 = 0x2115

	// CommandCodeMemoryCassetteTransfer Command code: memory cassette transfer (CP1H and CP1L CPU units only)
	CommandCodeMemoryCassetteTransfer uint16 = 0x2120

	// CommandCodeForcedSetReset Command code: forced set/reset
	CommandCodeForcedSetReset uint16 = 0x2301

	// CommandCodeForcedSetResetCancel Command code: forced set/reset cancel
	CommandCodeForcedSetResetCancel uint16 = 0x2302

	// CommandCodeConvertToCompoWayFCommand Command code: convert to CompoWay/F command
	CommandCodeConvertToCompoWayFCommand uint16 = 0x2803

	// CommandCodeConvertToModbusRTUCommand Command code: convert to Modbus-RTU command
	CommandCodeConvertToModbusRTUCommand uint16 = 0x2804

	// CommandCodeConvertToModbusASCIICommand Command code: convert to Modbus-ASCII command
	CommandCodeConvertToModbusASCIICommand uint16 = 0x2805
)

// Data taken from Omron document Cat. No. W342-E1-15, pages 155-161
const (
	// EndCodeNormalCompletion End code: normal completion
	EndCodeNormalCompletion uint16 = 0x0000

	// EndCodeServiceInterrupted End code: normal completion; service was interrupted
	EndCodeServiceInterrupted uint16 = 0x0001

	// EndCodeLocalNodeNotInNetwork End code: local node error; local node not in network
	EndCodeLocalNodeNotInNetwork uint16 = 0x0101

	// EndCodeTokenTimeout End code: local node error; token timeout
	EndCodeTokenTimeout uint16 = 0x0102

	// EndCodeRetriesFailed End code: local node error; retries failed
	EndCodeRetriesFailed uint16 = 0x0103

	// EndCodeTooManySendFrames End code: local node error; too many send frames
	EndCodeTooManySendFrames uint16 = 0x0104

	// EndCodeNodeAddressRangeError End code: local node error; node address range error
	EndCodeNodeAddressRangeError uint16 = 0x0105

	// EndCodeNodeAddressRangeDuplication End code: local node error; node address range duplication
	EndCodeNodeAddressRangeDuplication uint16 = 0x0106

	// EndCodeDestinationNodeNotInNetwork End code: destination node error; destination node not in network
	EndCodeDestinationNodeNotInNetwork uint16 = 0x0201

	// EndCodeUnitMissing End code: destination node error; unit missing
	EndCodeUnitMissing uint16 = 0x0202

	// EndCodeThirdNodeMissing End code: destination node error; third node missing
	EndCodeThirdNodeMissing uint16 = 0x0203

	// EndCodeDestinationNodeBusy End code: destination node error; destination node busy
	EndCodeDestinationNodeBusy uint16 = 0x0204

	// EndCodeResponseTimeout End code: destination node error; response timeout
	EndCodeResponseTimeout uint16 = 0x0205

	// EndCodeCommunicationsControllerError End code: controller error; communication controller error
	EndCodeCommunicationsControllerError uint16 = 0x0301

	// EndCodeCPUUnitError End code: controller error; CPU unit error
	EndCodeCPUUnitError uint16 = 0x0302

	// EndCodeControllerError End code:  controller error; controller error
	EndCodeControllerError uint16 = 0x0303

	// EndCodeUnitNumberError End code: controller error; unit number error
	EndCodeUnitNumberError uint16 = 0x0304

	// EndCodeUndefinedCommand End code: service unsupported; undefined command
	EndCodeUndefinedCommand uint16 = 0x0401

	// EndCodeNotSupportedByModelVersion End code: service unsupported; not supported by model version
	EndCodeNotSupportedByModelVersion uint16 = 0x0402

	// EndCodeDestinationAddressSettingError End code: routing table error; destination address setting error
	EndCodeDestinationAddressSettingError uint16 = 0x0501

	// EndCodeNoRoutingTables End code: routing table error; no routing tables
	EndCodeNoRoutingTables uint16 = 0x0502

	// EndCodeRoutingTableError End code: routing table error; routing table error
	EndCodeRoutingTableError uint16 = 0x0503

	// EndCodeTooManyRelays End code: routing table error; too many relays
	EndCodeTooManyRelays uint16 = 0x0504

	// EndCodeCommandTooLong End code: command format error; command too long
	EndCodeCommandTooLong uint16 = 0x1001

	// EndCodeCommandTooShort End code: command format error; command too short
	EndCodeCommandTooShort uint16 = 0x1002

	// EndCodeElementsDataDontMatch End code: command format error; elements/data don't match
	EndCodeElementsDataDontMatch uint16 = 0x1003

	// EndCodeCommandFormatError End code: command format error; command format error
	EndCodeCommandFormatError uint16 = 0x1004

	// EndCodeHeaderError End code: command format error; header error
	EndCodeHeaderError uint16 = 0x1005

	// EndCodeAreaClassificationMissing End code: parameter error; classification missing
	EndCodeAreaClassificationMissing uint16 = 0x1101

	// EndCodeAccessSizeError End code: parameter error; access size error
	EndCodeAccessSizeError uint16 = 0x1102

	// EndCodeAddressRangeError End code: parameter error; address range error
	EndCodeAddressRangeError uint16 = 0x1103

	// EndCodeAddressRangeExceeded End code: parameter error; address range exceeded
	EndCodeAddressRangeExceeded uint16 = 0x1104

	// EndCodeProgramMissing End code: parameter error; program missing
	EndCodeProgramMissing uint16 = 0x1106

	// EndCodeRelationalError End code: parameter error; relational error
	EndCodeRelationalError uint16 = 0x1109

	// EndCodeDuplicateDataAccess End code: parameter error; duplicate data access
	EndCodeDuplicateDataAccess uint16 = 0x110a

	// EndCodeResponseTooBig End code: parameter error; response too big
	EndCodeResponseTooBig uint16 = 0x110b

	// EndCodeParameterError End code: parameter error
	EndCodeParameterError uint16 = 0x110c

	// EndCodeReadNotPossibleProtected End code: read not possible; protected
	EndCodeReadNotPossibleProtected uint16 = 0x2002

	// EndCodeReadNotPossibleTableMissing End code: read not possible; table missing
	EndCodeReadNotPossibleTableMissing uint16 = 0x2003

	// EndCodeReadNotPossibleDataMissing End code: read not possible; data missing
	EndCodeReadNotPossibleDataMissing uint16 = 0x2004

	// EndCodeReadNotPossibleProgramMissing End code: read not possible; program missing
	EndCodeReadNotPossibleProgramMissing uint16 = 0x2005

	// EndCodeReadNotPossibleFileMissing End code: read not possible; file missing
	EndCodeReadNotPossibleFileMissing uint16 = 0x2006

	// EndCodeReadNotPossibleDataMismatch End code: read not possible; data mismatch
	EndCodeReadNotPossibleDataMismatch uint16 = 0x2007

	// EndCodeWriteNotPossibleReadOnly End code: write not possible; read only
	EndCodeWriteNotPossibleReadOnly uint16 = 0x2101

	// EndCodeWriteNotPossibleProtected End code: write not possible; write protected
	EndCodeWriteNotPossibleProtected uint16 = 0x2102

	// EndCodeWriteNotPossibleCannotRegister End code: write not possible; cannot register
	EndCodeWriteNotPossibleCannotRegister uint16 = 0x2103

	// EndCodeWriteNotPossibleProgramMissing End code: write not possible; program missing
	EndCodeWriteNotPossibleProgramMissing uint16 = 0x2105

	// EndCodeWriteNotPossibleFileMissing End code: write not possible; file missing
	EndCodeWriteNotPossibleFileMissing uint16 = 0x2106

	// EndCodeWriteNotPossibleFileNameAlreadyExists End code: write not possible; file name already exists
	EndCodeWriteNotPossibleFileNameAlreadyExists uint16 = 0x2107

	// EndCodeWriteNotPossibleCannotChange End code: write not possible; cannot change
	EndCodeWriteNotPossibleCannotChange uint16 = 0x2108

	// EndCodeNotExecutableInCurrentModeNotPossibleDuringExecution End code: not executeable in current mode during execution
	EndCodeNotExecutableInCurrentModeNotPossibleDuringExecution uint16 = 0x2201

	// EndCodeNotExecutableInCurrentModeNotPossibleWhileRunning End code: not executeable in current mode while running
	EndCodeNotExecutableInCurrentModeNotPossibleWhileRunning uint16 = 0x2202

	// EndCodeNotExecutableInCurrentModeWrongPLCModeInProgram End code: not executeable in current mode; PLC is in PROGRAM mode
	EndCodeNotExecutableInCurrentModeWrongPLCModeInProgram uint16 = 0x2203

	// EndCodeNotExecutableInCurrentModeWrongPLCModeInDebug End code: not executeable in current mode; PLC is in DEBUG mode
	EndCodeNotExecutableInCurrentModeWrongPLCModeInDebug uint16 = 0x2204

	// EndCodeNotExecutableInCurrentModeWrongPLCModeInMonitor End code: not executeable in current mode; PLC is in MONITOR mode
	EndCodeNotExecutableInCurrentModeWrongPLCModeInMonitor uint16 = 0x2205

	// EndCodeNotExecutableInCurrentModeWrongPLCModeInRun End code: not executeable in current mode; PLC is in RUN mode
	EndCodeNotExecutableInCurrentModeWrongPLCModeInRun uint16 = 0x2206

	// EndCodeNotExecutableInCurrentModeSpecifiedNodeNotPollingNode End code: not executeable in current mode; specified node is not polling node
	EndCodeNotExecutableInCurrentModeSpecifiedNodeNotPollingNode uint16 = 0x2207

	// EndCodeNotExecutableInCurrentModeStepCannotBeExecuted End code: not executeable in current mode; step cannot be executed
	EndCodeNotExecutableInCurrentModeStepCannotBeExecuted uint16 = 0x2208

	// EndCodeNoSuchDeviceFileDeviceMissing End code: no such device; file device missing
	EndCodeNoSuchDeviceFileDeviceMissing uint16 = 0x2301

	// EndCodeNoSuchDeviceMemoryMissing End code: no such device; memory missing
	EndCodeNoSuchDeviceMemoryMissing uint16 = 0x2302

	// EndCodeNoSuchDeviceClockMissing End code: no such device; clock missing
	EndCodeNoSuchDeviceClockMissing uint16 = 0x2303

	// EndCodeCannotStartStopTableMissing End code: cannot start/stop; table missing
	EndCodeCannotStartStopTableMissing uint16 = 0x2401

	// EndCodeUnitErrorMemoryError End code: unit error; memory error
	EndCodeUnitErrorMemoryError uint16 = 0x2502

	// EndCodeUnitErrorIOError End code: unit error; IO error
	EndCodeUnitErrorIOError uint16 = 0x2503

	// EndCodeUnitErrorTooManyIOPoints End code: unit error; too many IO points
	EndCodeUnitErrorTooManyIOPoints uint16 = 0x2504

	// EndCodeUnitErrorCPUBusError End code: unit error; CPU bus error
	EndCodeUnitErrorCPUBusError uint16 = 0x2505

	// EndCodeUnitErrorIODuplication End code: unit error; IO duplication
	EndCodeUnitErrorIODuplication uint16 = 0x2506

	// EndCodeUnitErrorIOBusError End code: unit error; IO bus error
	EndCodeUnitErrorIOBusError uint16 = 0x2507

	// EndCodeUnitErrorSYSMACBUS2Error End code: unit error; SYSMAC BUS/2 error
	EndCodeUnitErrorSYSMACBUS2Error uint16 = 0x2509

	// EndCodeUnitErrorCPUBusUnitError End code: unit error; CPU bus unit error
	EndCodeUnitErrorCPUBusUnitError uint16 = 0x250a

	// EndCodeUnitErrorSYSMACBusNumberDuplication End code: unit error; SYSMAC bus number duplication
	EndCodeUnitErrorSYSMACBusNumberDuplication uint16 = 0x250d

	// EndCodeUnitErrorMemoryStatusError End code: unit error; memory status error
	EndCodeUnitErrorMemoryStatusError uint16 = 0x250f

	// EndCodeUnitErrorSYSMACBusTerminatorMissing End code: unit error; SYSMAC bus terminator missing
	EndCodeUnitErrorSYSMACBusTerminatorMissing uint16 = 0x2510

	// EndCodeCommandErrorNoProtection End code: command error; no protection
	EndCodeCommandErrorNoProtection uint16 = 0x2601

	// EndCodeCommandErrorIncorrectPassword End code: command error; incorrect password
	EndCodeCommandErrorIncorrectPassword uint16 = 0x2602

	// EndCodeCommandErrorProtected End code: command error; protected
	EndCodeCommandErrorProtected uint16 = 0x2604

	// EndCodeCommandErrorServiceAlreadyExecuting End code: command error; service already executing
	EndCodeCommandErrorServiceAlreadyExecuting uint16 = 0x2605

	// EndCodeCommandErrorServiceStopped End code: command error; service stopped
	EndCodeCommandErrorServiceStopped uint16 = 0x2606

	// EndCodeCommandErrorNoExecutionRight End code: command error; no execution right
	EndCodeCommandErrorNoExecutionRight uint16 = 0x2607

	// EndCodeCommandErrorSettingsNotComplete End code: command error; settings not complete
	EndCodeCommandErrorSettingsNotComplete uint16 = 0x2608

	// EndCodeCommandErrorNecessaryItemsNotSet End code: command error; necessary items not set
	EndCodeCommandErrorNecessaryItemsNotSet uint16 = 0x2609

	// EndCodeCommandErrorNumberAlreadyDefined End code: command error; number already defined
	EndCodeCommandErrorNumberAlreadyDefined uint16 = 0x260a

	// EndCodeCommandErrorErrorWillNotClear End code: command error; error will not clear
	EndCodeCommandErrorErrorWillNotClear uint16 = 0x260b

	// EndCodeAccessWriteErrorNoAccessRight End code: access write error; no access right
	EndCodeAccessWriteErrorNoAccessRight uint16 = 0x3001

	// EndCodeAbortServiceAborted End code: abort; service aborted
	EndCodeAbortServiceAborted uint16 = 0x4001
)

const (
	// MemoryAreaCIOBit Memory area: CIO area; bit
	MemoryAreaCIOBit byte = 0x30

	// MemoryAreaWRBit Memory area: work area; bit
	MemoryAreaWRBit byte = 0x31

	// MemoryAreaHRBit Memory area: holding area; bit
	MemoryAreaHRBit byte = 0x32

	// MemoryAreaARBit Memory area: axuillary area; bit
	MemoryAreaARBit byte = 0x33

	// MemoryAreaCIOWord Memory area: CIO area; word
	MemoryAreaCIOWord byte = 0xb0

	// MemoryAreaWRWord Memory area: work area; word
	MemoryAreaWRWord byte = 0xb1

	// MemoryAreaHRWord Memory area: holding area; word
	MemoryAreaHRWord byte = 0xb2

	// MemoryAreaARWord Memory area: auxillary area; word
	MemoryAreaARWord byte = 0xb3

	// MemoryAreaTimerCounterCompletionFlag Memory area: counter completion flag
	MemoryAreaTimerCounterCompletionFlag byte = 0x09

	// MemoryAreaTimerCounterPV Memory area: counter PV
	MemoryAreaTimerCounterPV byte = 0x89

	// MemoryAreaDMBit Memory area: data area; bit
	MemoryAreaDMBit byte = 0x02

	// MemoryAreaDMWord Memory area: data area; word
	MemoryAreaDMWord byte = 0x82

	// MemoryAreaTaskBit Memory area: task flags; bit
	MemoryAreaTaskBit byte = 0x06

	// MemoryAreaTaskStatus Memory area: task flags; status
	MemoryAreaTaskStatus byte = 0x46

	// MemoryAreaIndexRegisterPV Memory area: CIO bit
	MemoryAreaIndexRegisterPV byte = 0xdc

	// MemoryAreaDataRegisterPV Memory area: CIO bit
	MemoryAreaDataRegisterPV byte = 0xbc

	// MemoryAreaClockPulsesConditionFlagsBit Memory area: CIO bit
	MemoryAreaClockPulsesConditionFlagsBit byte = 0x07
)

type ResponseTimeoutError struct {
	duration time.Duration
}

func (e ResponseTimeoutError) Error() string {
	return fmt.Sprintf("Response timeout of %d has been reached", e.duration)
}

type IncompatibleMemoryAreaError struct {
	area byte
}

func (e IncompatibleMemoryAreaError) Error() string {
	return fmt.Sprintf("The memory area is incompatible with the data type to be read: 0x%X", e.area)
}

// Driver errors

type BCDBadDigitError struct {
	v   string
	val uint64
}

func (e BCDBadDigitError) Error() string {
	return fmt.Sprintf("Bad digit in BCD decoding: %s = %d", e.v, e.val)
}

type BCDOverflowError struct{}

func (e BCDOverflowError) Error() string {
	return "Overflow occurred in BCD decoding"
}

// Header A FINS frame header
type Header struct {
	messageType      uint8
	responseRequired bool
	src              finsAddress
	dst              finsAddress
	serviceID        byte
	gatewayCount     uint8
}

const (
	// MessageTypeCommand Command message type
	MessageTypeCommand uint8 = iota

	// MessageTypeResponse Response message type
	MessageTypeResponse uint8 = iota
)

func defaultHeader(messageType uint8, responseRequired bool, src finsAddress, dst finsAddress, serviceID byte) Header {
	h := Header{}
	h.messageType = messageType
	h.responseRequired = responseRequired
	h.gatewayCount = 2
	h.src = src
	h.dst = dst
	h.serviceID = serviceID
	return h
}

func defaultCommandHeader(src finsAddress, dst finsAddress, serviceID byte) Header {
	h := defaultHeader(MessageTypeCommand, true, src, dst, serviceID)
	return h
}

func defaultResponseHeader(commandHeader Header) Header {
	h := defaultHeader(MessageTypeResponse, false, commandHeader.dst, commandHeader.src, commandHeader.serviceID)
	return h
}

// request A FINS command request
type request struct {
	header      Header
	commandCode uint16
	data        []byte
}

// response A FINS command response
type response struct {
	header      Header
	commandCode uint16
	endCode     uint16
	data        []byte
}

// memoryAddress A plc memory address to do a work
type memoryAddress struct {
	memoryArea byte
	address    uint16
	bitOffset  byte
}

func memAddr(memoryArea byte, address uint16) memoryAddress {
	return memAddrWithBitOffset(memoryArea, address, 0)
}

func memAddrWithBitOffset(memoryArea byte, address uint16, bitOffset byte) memoryAddress {
	return memoryAddress{memoryArea, address, bitOffset}
}

func readCommand(memoryAddr memoryAddress, itemCount uint16) []byte {
	commandData := make([]byte, 2, 8)
	binary.BigEndian.PutUint16(commandData[0:2], CommandCodeMemoryAreaRead)
	commandData = append(commandData, encodeMemoryAddress(memoryAddr)...)
	commandData = append(commandData, []byte{0, 0}...)
	binary.BigEndian.PutUint16(commandData[6:8], itemCount)
	return commandData
}

func writeCommand(memoryAddr memoryAddress, itemCount uint16, bytes []byte) []byte {
	commandData := make([]byte, 2, 8+len(bytes))
	binary.BigEndian.PutUint16(commandData[0:2], CommandCodeMemoryAreaWrite)
	commandData = append(commandData, encodeMemoryAddress(memoryAddr)...)
	commandData = append(commandData, []byte{0, 0}...)
	binary.BigEndian.PutUint16(commandData[6:8], itemCount)
	commandData = append(commandData, bytes...)
	return commandData
}

func clockReadCommand() []byte {
	commandData := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(commandData[0:2], CommandCodeClockRead)
	return commandData
}

func encodeMemoryAddress(memoryAddr memoryAddress) []byte {
	bytes := make([]byte, 4, 4)
	bytes[0] = memoryAddr.memoryArea
	binary.BigEndian.PutUint16(bytes[1:3], memoryAddr.address)
	bytes[3] = memoryAddr.bitOffset
	return bytes
}

func decodeMemoryAddress(data []byte) memoryAddress {
	return memoryAddress{data[0], binary.BigEndian.Uint16(data[1:3]), data[3]}
}

func decodeRequest(bytes []byte) request {
	return request{
		decodeHeader(bytes[0:10]),
		binary.BigEndian.Uint16(bytes[10:12]),
		bytes[12:],
	}
}

func decodeResponse(bytes []byte) response {
	return response{
		decodeHeader(bytes[0:10]),
		binary.BigEndian.Uint16(bytes[10:12]),
		binary.BigEndian.Uint16(bytes[12:14]),
		bytes[14:],
	}
}
func encodeResponse(resp response) []byte {
	bytes := make([]byte, 4, 4+len(resp.data))
	binary.BigEndian.PutUint16(bytes[0:2], resp.commandCode)
	binary.BigEndian.PutUint16(bytes[2:4], resp.endCode)
	bytes = append(bytes, resp.data...)
	bh := encodeHeader(resp.header)
	bh = append(bh, bytes...)
	return bh
}

const (
	icfBridgesBit          byte = 7
	icfMessageTypeBit      byte = 6
	icfResponseRequiredBit byte = 0
)

func decodeHeader(bytes []byte) Header {
	header := Header{}
	icf := bytes[0]
	if icf&1<<icfResponseRequiredBit == 0 {
		header.responseRequired = true
	}
	if icf&1<<icfMessageTypeBit == 0 {
		header.messageType = MessageTypeCommand
	} else {
		header.messageType = MessageTypeResponse
	}
	header.gatewayCount = bytes[2]
	header.dst = finsAddress{bytes[3], bytes[4], bytes[5]}
	header.src = finsAddress{bytes[6], bytes[7], bytes[8]}
	header.serviceID = bytes[9]

	return header
}

func encodeHeader(h Header) []byte {
	var icf byte
	icf = 1 << icfBridgesBit
	if h.responseRequired == false {
		icf |= 1 << icfResponseRequiredBit
	}
	if h.messageType == MessageTypeResponse {
		icf |= 1 << icfMessageTypeBit
	}
	bytes := []byte{
		icf, 0x00, h.gatewayCount,
		h.dst.network, h.dst.node, h.dst.unit,
		h.src.network, h.src.node, h.src.unit,
		h.serviceID}
	return bytes
}

func encodeBCD(x uint64) []byte {
	if x == 0 {
		return []byte{0x0f}
	}
	var n int
	for xx := x; xx > 0; n++ {
		xx = xx / 10
	}
	bcd := make([]byte, (n+1)/2)
	if n%2 == 1 {
		hi, lo := byte(x%10), byte(0x0f)
		bcd[(n-1)/2] = hi<<4 | lo
		x = x / 10
		n--
	}
	for i := n/2 - 1; i >= 0; i-- {
		hi, lo := byte((x/10)%10), byte(x%10)
		bcd[i] = hi<<4 | lo
		x = x / 100
	}
	return bcd
}

func timesTenPlusCatchingOverflow(x uint64, digit uint64) (uint64, error) {
	x5 := x<<2 + x
	if int64(x5) < 0 || x5<<1 > ^digit {
		return 0, BCDOverflowError{}
	}
	return x5<<1 + digit, nil
}

func decodeBCD(bcd []byte) (x uint64, err error) {
	for i, b := range bcd {
		hi, lo := uint64(b>>4), uint64(b&0x0f)
		if hi > 9 {
			return 0, BCDBadDigitError{"hi", hi}
		}
		x, err = timesTenPlusCatchingOverflow(x, hi)
		if err != nil {
			return 0, err
		}
		if lo == 0x0f && i == len(bcd)-1 {
			return x, nil
		}
		if lo > 9 {
			return 0, BCDBadDigitError{"lo", lo}
		}
		x, err = timesTenPlusCatchingOverflow(x, lo)
		if err != nil {
			return 0, err
		}
	}
	return x, nil
}

// finsAddress A FINS device address
type finsAddress struct {
	network byte
	node    byte
	unit    byte
}

// Address A full device address
type Address struct {
	finsAddress finsAddress
	udpAddress  *net.UDPAddr
}

func NewAddress(ip string, port int, network, node, unit byte) Address {
	return Address{
		udpAddress: &net.UDPAddr{
			IP:   net.ParseIP(ip),
			Port: port,
		},
		finsAddress: finsAddress{
			network: network,
			node:    node,
			unit:    unit,
		},
	}
}

const DEFAULT_RESPONSE_TIMEOUT = 20 // ms

// Client Omron FINS client
type Client struct {
	conn *net.UDPConn
	resp []chan response
	sync.Mutex
	dst               finsAddress
	src               finsAddress
	sid               byte
	closed            bool
	responseTimeoutMs time.Duration
	byteOrder         binary.ByteOrder
}

// NewClient creates a new Omron FINS client
func NewClient(localAddr, plcAddr Address) (*Client, error) {
	c := new(Client)
	c.dst = plcAddr.finsAddress
	c.src = localAddr.finsAddress
	c.responseTimeoutMs = DEFAULT_RESPONSE_TIMEOUT
	c.byteOrder = binary.BigEndian

	conn, err := net.DialUDP("udp", localAddr.udpAddress, plcAddr.udpAddress)
	if err != nil {
		return nil, err
	}
	c.conn = conn

	c.resp = make([]chan response, 256) //storage for all responses, sid is byte - only 256 values
	go c.listenLoop()
	return c, nil
}

// Set byte order
// Default value: binary.BigEndian
func (c *Client) SetByteOrder(o binary.ByteOrder) {
	c.byteOrder = o
}

// Set response timeout duration (ms).
// Default value: 20ms.
// A timeout of zero can be used to block indefinitely.
func (c *Client) SetTimeoutMs(t uint) {
	c.responseTimeoutMs = time.Duration(t)
}

// Close Closes an Omron FINS connection
func (c *Client) Close() {
	c.closed = true
	c.conn.Close()
}

// ReadWords Reads words from the PLC data area
func (c *Client) ReadWords(memoryArea byte, address uint16, readCount uint16) ([]uint16, error) {
	if checkIsWordMemoryArea(memoryArea) == false {
		return nil, IncompatibleMemoryAreaError{memoryArea}
	}
	command := readCommand(memAddr(memoryArea, address), readCount)
	r, e := c.sendCommand(command)
	e = checkResponse(r, e)
	if e != nil {
		return nil, e
	}

	data := make([]uint16, readCount, readCount)
	for i := 0; i < int(readCount); i++ {
		data[i] = c.byteOrder.Uint16(r.data[i*2 : i*2+2])
	}

	return data, nil
}

// ReadBytes Reads bytes from the PLC data area
func (c *Client) ReadBytes(memoryArea byte, address uint16, readCount uint16) ([]byte, error) {
	if checkIsWordMemoryArea(memoryArea) == false {
		return nil, IncompatibleMemoryAreaError{memoryArea}
	}
	command := readCommand(memAddr(memoryArea, address), readCount)
	r, e := c.sendCommand(command)
	e = checkResponse(r, e)
	if e != nil {
		return nil, e
	}

	return r.data, nil
}

// ReadString Reads a string from the PLC data area
func (c *Client) ReadString(memoryArea byte, address uint16, readCount uint16) (string, error) {
	data, e := c.ReadBytes(memoryArea, address, readCount)
	if e != nil {
		return "", e
	}
	n := bytes.IndexByte(data, 0)
	if n == -1 {
		n = len(data)
	}
	return string(data[:n]), nil
}

// ReadBits Reads bits from the PLC data area
func (c *Client) ReadBits(memoryArea byte, address uint16, bitOffset byte, readCount uint16) ([]bool, error) {
	if checkIsBitMemoryArea(memoryArea) == false {
		return nil, IncompatibleMemoryAreaError{memoryArea}
	}
	command := readCommand(memAddrWithBitOffset(memoryArea, address, bitOffset), readCount)
	r, e := c.sendCommand(command)
	e = checkResponse(r, e)
	if e != nil {
		return nil, e
	}

	data := make([]bool, readCount, readCount)
	for i := 0; i < int(readCount); i++ {
		data[i] = r.data[i]&0x01 > 0
	}

	return data, nil
}

// ReadClock Reads the PLC clock
func (c *Client) ReadClock() (*time.Time, error) {
	r, e := c.sendCommand(clockReadCommand())
	e = checkResponse(r, e)
	if e != nil {
		return nil, e
	}
	year, _ := decodeBCD(r.data[0:1])
	if year < 50 {
		year += 2000
	} else {
		year += 1900
	}
	month, _ := decodeBCD(r.data[1:2])
	day, _ := decodeBCD(r.data[2:3])
	hour, _ := decodeBCD(r.data[3:4])
	minute, _ := decodeBCD(r.data[4:5])
	second, _ := decodeBCD(r.data[5:6])

	t := time.Date(
		int(year), time.Month(month), int(day), int(hour), int(minute), int(second),
		0, // nanosecond
		time.Local,
	)
	return &t, nil
}

// WriteWords Writes words to the PLC data area
func (c *Client) WriteWords(memoryArea byte, address uint16, data []uint16) error {
	if checkIsWordMemoryArea(memoryArea) == false {
		return IncompatibleMemoryAreaError{memoryArea}
	}
	l := uint16(len(data))
	bts := make([]byte, 2*l, 2*l)
	for i := 0; i < int(l); i++ {
		c.byteOrder.PutUint16(bts[i*2:i*2+2], data[i])
	}
	command := writeCommand(memAddr(memoryArea, address), l, bts)

	return checkResponse(c.sendCommand(command))
}

// WriteString Writes a string to the PLC data area
func (c *Client) WriteString(memoryArea byte, address uint16, s string) error {
	if checkIsWordMemoryArea(memoryArea) == false {
		return IncompatibleMemoryAreaError{memoryArea}
	}
	bts := make([]byte, 2*len(s), 2*len(s))
	copy(bts, s)

	command := writeCommand(memAddr(memoryArea, address), uint16((len(s)+1)/2), bts) //TODO: test on real PLC

	return checkResponse(c.sendCommand(command))
}

// WriteBytes Writes bytes array to the PLC data area
func (c *Client) WriteBytes(memoryArea byte, address uint16, b []byte) error {
	if checkIsWordMemoryArea(memoryArea) == false {
		return IncompatibleMemoryAreaError{memoryArea}
	}
	command := writeCommand(memAddr(memoryArea, address), uint16(len(b)), b)
	return checkResponse(c.sendCommand(command))
}

// WriteBits Writes bits to the PLC data area
func (c *Client) WriteBits(memoryArea byte, address uint16, bitOffset byte, data []bool) error {
	if checkIsBitMemoryArea(memoryArea) == false {
		return IncompatibleMemoryAreaError{memoryArea}
	}
	l := uint16(len(data))
	bts := make([]byte, 0, l)
	var d byte
	for i := 0; i < int(l); i++ {
		if data[i] {
			d = 0x01
		} else {
			d = 0x00
		}
		bts = append(bts, d)
	}
	command := writeCommand(memAddrWithBitOffset(memoryArea, address, bitOffset), l, bts)

	return checkResponse(c.sendCommand(command))
}

// SetBit Sets a bit in the PLC data area
func (c *Client) SetBit(memoryArea byte, address uint16, bitOffset byte) error {
	return c.bitTwiddle(memoryArea, address, bitOffset, 0x01)
}

// ResetBit Resets a bit in the PLC data area
func (c *Client) ResetBit(memoryArea byte, address uint16, bitOffset byte) error {
	return c.bitTwiddle(memoryArea, address, bitOffset, 0x00)
}

// ToggleBit Toggles a bit in the PLC data area
func (c *Client) ToggleBit(memoryArea byte, address uint16, bitOffset byte) error {
	b, e := c.ReadBits(memoryArea, address, bitOffset, 1)
	if e != nil {
		return e
	}
	var t byte
	if b[0] {
		t = 0x00
	} else {
		t = 0x01
	}
	return c.bitTwiddle(memoryArea, address, bitOffset, t)
}

func (c *Client) bitTwiddle(memoryArea byte, address uint16, bitOffset byte, value byte) error {
	if checkIsBitMemoryArea(memoryArea) == false {
		return IncompatibleMemoryAreaError{memoryArea}
	}
	mem := memoryAddress{memoryArea, address, bitOffset}
	command := writeCommand(mem, 1, []byte{value})

	return checkResponse(c.sendCommand(command))
}

func checkResponse(r *response, e error) error {
	if e != nil {
		return e
	}
	if r.endCode != EndCodeNormalCompletion {
		return fmt.Errorf("error reported by destination, end code 0x%x", r.endCode)
	}
	return nil
}

func (c *Client) nextHeader() *Header {
	sid := c.incrementSid()
	header := defaultCommandHeader(c.src, c.dst, sid)
	return &header
}

func (c *Client) incrementSid() byte {
	c.Lock() //thread-safe sid incrementation
	c.sid++
	sid := c.sid
	c.Unlock()
	c.resp[sid] = make(chan response) //clearing cell of storage for new response
	return sid
}

func (c *Client) sendCommand(command []byte) (*response, error) {
	header := c.nextHeader()
	bts := encodeHeader(*header)
	bts = append(bts, command...)
	_, err := (*c.conn).Write(bts)
	if err != nil {
		return nil, err
	}

	// if response timeout is zero, block indefinitely
	if c.responseTimeoutMs > 0 {
		select {
		case resp := <-c.resp[header.serviceID]:
			return &resp, nil
		case <-time.After(c.responseTimeoutMs * time.Millisecond):
			return nil, ResponseTimeoutError{c.responseTimeoutMs}
		}
	} else {
		resp := <-c.resp[header.serviceID]
		return &resp, nil
	}
}

func (c *Client) listenLoop() {
	for {
		buf := make([]byte, 2048)
		n, err := bufio.NewReader(c.conn).Read(buf)
		if err != nil {
			// do not complain when connection is closed by user
			if !c.closed {
				log.Fatal(err)
			}
			break
		}

		if n > 0 {
			ans := decodeResponse(buf[:n])
			c.resp[ans.header.serviceID] <- ans
		} else {
			log.Println("cannot read response: ", buf)
		}
	}
}

func checkIsWordMemoryArea(memoryArea byte) bool {
	if memoryArea == MemoryAreaDMWord ||
		memoryArea == MemoryAreaARWord ||
		memoryArea == MemoryAreaHRWord ||
		memoryArea == MemoryAreaWRWord {
		return true
	}
	return false
}

func checkIsBitMemoryArea(memoryArea byte) bool {
	if memoryArea == MemoryAreaDMBit ||
		memoryArea == MemoryAreaARBit ||
		memoryArea == MemoryAreaHRBit ||
		memoryArea == MemoryAreaWRBit {
		return true
	}
	return false
}

func main() {

	clientAddr := NewAddress("127.0.0.1", 10000, 0, 34, 0)
	port, _ := strconv.Atoi(os.Args[2])
	plcAddr := NewAddress(os.Args[1], port, 0, 0, 0)

	c, err := NewClient(clientAddr, plcAddr)
	if err != nil {
		log.Println(err)
	}
	defer c.Close()

	b, err := c.ReadBytes(MemoryAreaDMWord, 0, 19)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Result:", string(b))
}
