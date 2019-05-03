package aiff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-audio/audio"
)

func (e *Encoder) WriteF32(buf *audio.Float32Buffer) error {
	if err := e.writeHeader(); err != nil {
		return err
	}
	if err := e.fixBeginning(); err != nil {
		return err
	}
	return e.addF32Buffer(buf)
}

func (e *Encoder) fixBeginning() error {
	if !e.pcmChunkStarted {
		// other chunks audio frames
		if err := e.AddBE([]byte("SSND")); err != nil {
			return fmt.Errorf("%v when writing SSND chunk ID header", err)
		}
		e.pcmChunkSizePos = e.WrittenBytes
		e.pcmChunkStarted = true
		// temporary blocksize uint32 chunksize :=
		//uint32((int(e.BitDepth)/8)*int(e.NumChans)*len(e.Frames)
		//+ 8)
		if err := e.AddBE(uint32(84)); err != nil {
			return fmt.Errorf("%v when writing SSND chunk size header", err)
		}
		if err := e.AddBE(uint32(0)); err != nil {
			return fmt.Errorf("%v when writing SSND offset",
				err)
		}
		if err := e.AddBE(uint32(0)); err != nil {
			return fmt.Errorf("%v when writing SSND block size", err)
		}
	}
	return nil
}

func (e *Encoder) addF32Buffer(buf *audio.Float32Buffer) error {
	if buf == nil {
		return fmt.Errorf("can't add a nil buffer")
	}

	frameCount := buf.NumFrames()
	// setup a buffer so we don't do many writes
	bb := bytes.NewBuffer(nil)
	var err error
	for i := 0; i < frameCount; i++ {
		for j := 0; j < buf.Format.NumChannels; j++ {
			v := buf.Data[i*buf.Format.NumChannels+j]
			switch e.BitDepth {
			case 32:
				if err = binary.Write(bb, binary.BigEndian, v); err != nil {
					return err
				}
			default:
				return fmt.Errorf("can't add frames of bit size %d", e.BitDepth)
			}
		}
		e.frames++
	}
	n, err := e.w.Write(bb.Bytes())
	e.WrittenBytes += n
	return err
}
