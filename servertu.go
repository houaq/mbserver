package mbserver

import (
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
	port.Flush()
	s.ports = append(s.ports, port)
	go s.acceptSerialRequests(port)
	return err
}

func readFullFrame(port *serial.Port, expectLen int) ([]byte, error) {
	buffer := make([]byte, 32)
	frame := make([]byte, 0)
	for len(frame) < expectLen {
		n, err := port.Read(buffer)
		if err != nil {
			return nil, err
		}

		chunk := buffer[:n]
		log.Printf("serial read: [% 0x]\n", chunk)
		frame = append(frame, chunk...)
	}
	return frame, nil
}

func (s *Server) acceptSerialRequests(port *serial.Port) {
	for {
		packet, err := readFullFrame(port, 8)
		if err != nil {
			log.Printf("serial read error %s\n", err)
			continue
		}

		frame, err := NewRTUFrame(packet)
		if err != nil {
			log.Printf("bad serial frame error %v\n", err)
			port.Flush()
			continue
		}

		request := &Request{port, frame}
		s.requestChan <- request
	}
}
