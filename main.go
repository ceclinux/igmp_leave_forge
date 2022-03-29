package main

// https://www.omnisecu.com/tcpip/ipv4-protocol-and-ipv4-header.php
import (
	"flag"
	"log"
	"net"
	"syscall"
	"time"
)

func main() {
	num := flag.Int("n", 3, "number of iterations, -1 for infinite loop")
	src := flag.String("src", "", "source address")
	group := flag.String("group", "", "group address")
	flag.Parse()
	if *src == "" {
		log.Fatalf("source address cannot be empty")
	}
	if *group == "" {
		log.Fatalf("group address cannot be empty")
	}
	log.Printf("number of iterations", *num)
	log.Printf("source address", *src)
	log.Printf("group address", *group)

	srcIP := net.ParseIP(*src).To4()
	groupIP := net.ParseIP(*group).To4()
	for *num != 0 {
		igmp := []byte{
			0x17, //IGMPv2 leave
			0x00,
			0x00, // checksum initialized to 0
			0x00, // checksum initialized to 0
			groupIP[0],
			groupIP[1], groupIP[2], groupIP[3],
		}
		cs := checksum(igmp)
		igmp[2] = byte(cs)
		igmp[3] = byte(cs >> 8)

		h := []byte{
			0x45,
			0,      //Type of Service (ToS)
			0,      // total length(16 bits)
			20 + 8, // total length = 20 + len(igmp)
			0x6a,   // Identification(0x6a7f, fake)
			0x7f,   // Identification(0x6a7f, fake)
			0x40,   // Flags: 0x40, Don't fragment, first bit always 0,The next bit is called the DF (Don't Fragment) flag. DF flag set to "0" indicate that the IPv4 Datagram can be fragmented and DF set to 1 indicate "Don't Fragment" the IPv4 Datagram, The next bit is the MF (More Fragments) flag.  0b0100
			0,
			128,      // Time to live
			2,        // protocol(IGMP)
			0,        // header checksum(disabled)
			0,        // header checksum(disabled)
			srcIP[0], //source address
			srcIP[1],
			srcIP[2],
			srcIP[3],
			224, //destination address
			0,
			0,
			2,
		}

		t := append(h, igmp...)
		fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
		if err != nil {
			log.Fatal("Creating socket", err)
		}
		addr := syscall.SockaddrInet4{
			Port: 4444,
			Addr: [4]byte{224, 0, 0, 2},
		}
		err = syscall.Sendto(fd, t, 0, &addr)
		if err != nil {
			log.Fatal("Sendto:", err)
		}
		time.Sleep(1 * time.Second)
		*num -= 1
	}

}

// https://datatracker.ietf.org/doc/html/rfc2236#section-1.4
// The checksum is the 16-bit one's complement of the one's complement
// sum of the whole IGMP message (the entire IP payload).  For computing
// the checksum, the checksum field is set to zero.  When transmitting
// packets, the checksum MUST be computed and inserted into this field.
// When receiving packets, the checksum MUST be verified before
// processing a packet.
func checksum(b []byte) uint16 {
	var s uint32
	for i := 0; i < len(b); i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	// add back the carry
	s = s>>16 + s&0xffff
	s = s + s>>16
	return uint16(^s)
}
