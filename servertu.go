package mbserver

import (
	"bytes"
	"log"

	"github.com/tarm/serial"
)

// ListenRTU starts the Modbus server listening to a serial device.
// For example:  err := s.ListenRTU(&serial.Config{Address: "/dev/ttyUSB0"})
func (s *Server) ListenRTU(serialConfig *serial.Config) (err error) {
	port, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Fatalf("failed to open %s: %v\n", serialConfig.Name, err)
	}
	s.ports = append(s.ports, port)
	go s.acceptSerialRequests(port)
	return err
}

func (s *Server) acceptSerialRequests(port *serial.Port) {
	var bb bytes.Buffer
	buffer := make([]byte, 8)
	for {
		bb.Reset()
		for {
			bytesRead, err := port.Read(buffer)
			if err != nil {
				log.Printf("serial read error %s\n", err)
				return
			}
			if bytesRead > 0 {
				bb.Write(buffer)
				if bb.Len() >= 8 {
					break
				}
			}
		}

		frame, err := NewRTUFrame(bb.Bytes())
		if err != nil {
			log.Printf("bad serial frame error %v\n", err)
			continue
		}

		request := &Request{port, frame}
		s.requestChan <- request
	}
}
