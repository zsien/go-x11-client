package input

import (
	"github.com/linuxdeepin/go-x11-client"
)

// #WREQ
func encodeOpenDevice(deviceId uint8) (b x.RequestBody) {
	b.AddBlock(1).
		Write1b(deviceId).
		WritePad(3).
		End()
	return
}

type OpenDeviceReply struct {
	XIReplyType uint8
	ClassInfos  []ClassInfo
}

func readOpenDeviceReply(r *x.Reader, v *OpenDeviceReply) error {
	r.Read1b()
	if r.Err() != nil {
		return r.Err()
	}

	v.XIReplyType = r.Read1b()
	if r.Err() != nil {
		return r.Err()
	}

	// seq
	r.Read2b()
	if r.Err() != nil {
		return r.Err()
	}

	// length
	r.Read4b()
	if r.Err() != nil {
		return r.Err()
	}

	classInfosLen := int(r.Read1b())
	if r.Err() != nil {
		return r.Err()
	}

	r.ReadPad(23)
	if r.Err() != nil {
		return r.Err()
	}

	if classInfosLen > 0 {
		var err error
		v.ClassInfos = make([]ClassInfo, classInfosLen)
		for i := 0; i < classInfosLen; i++ {
			v.ClassInfos[i], err = readClassInfo(r)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func FindTypeAndClass(deviceId uint8, classInfos []ClassInfo,
	classId uint8, eventCode uint8) (Type uint8, class EventClass) {
	for _, classInfo := range classInfos {
		if classInfo.ClassId == classId {
			Type = classInfo.EventTypeBase + eventCode
			class = EventClass(uint32(deviceId)<<8 | uint32(Type))
		}
	}
	return
}

type ClassInfo struct {
	ClassId       uint8
	EventTypeBase uint8
}

func readClassInfo(r *x.Reader) (ClassInfo, error) {
	var v ClassInfo
	v.ClassId = r.Read1b()
	if r.Err() != nil {
		return ClassInfo{}, r.Err()
	}

	v.EventTypeBase = r.Read1b()
	if r.Err() != nil {
		return ClassInfo{}, r.Err()
	}

	return v, nil
}

// #WREQ
func encodeCloseDevice(deviceId uint8) (b x.RequestBody) {
	b.AddBlock(1).
		Write1b(deviceId).
		WritePad(3).
		End()
	return
}

// #WREQ
func encodeSelectExtensionEvent(window x.Window, classes []EventClass) (b x.RequestBody) {
	b0 := b.AddBlock(2 + len(classes)).
		Write4b(uint32(window)).
		Write2b(uint16(len(classes))).
		WritePad(2)

	for _, class := range classes {
		b0.Write4b(uint32(class))
	}
	b0.End()
	return
}
