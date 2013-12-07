package gomatrix

import (
	"errors"
	"net"
	"strconv"
)

type MatrixOrbital struct {
	tcpAddr             *net.TCPAddr
	tcpConn             *net.TCPConn
	maxColumns, maxRows int
}

func CreateMatrixOrbital(hostAddr string, port int, maxColumns int, maxRows int) (m *MatrixOrbital) {
	return newMatrixOrbital(hostAddr, port, maxColumns, maxRows)
}

func newMatrixOrbital(hostAddr string, port int, maxColumns int, maxRows int) (m *MatrixOrbital) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", hostAddr+":"+strconv.Itoa(port))
	if err != nil {
		// To-Do: Logging
		return nil
	}

	m = &MatrixOrbital{
		tcpAddr:    tcpAddr,
		maxColumns: maxColumns,
		maxRows:    maxRows,
	}

	return m
}

func (m *MatrixOrbital) Connect() (err error) {
	m.tcpConn, err = net.DialTCP("tcp", nil, m.tcpAddr)
	if err != nil {
		// To-Do: Logging
		return err
	}

	return nil
}

func (m *MatrixOrbital) Disconnect() (err error) {
	err = m.tcpConn.Close()
	if err != nil {
		return err
	}

	return err
}

func (m *MatrixOrbital) ClearScreen() (err error) {
	err = m.SendCmd(254, 88)
	return
}

func (m *MatrixOrbital) AutoTxKeyPressesOn() (err error) {
	err = m.SendCmd(254, 65)
	return
}

func (m *MatrixOrbital) AutoTxKeyPressesOff() (err error) {
	err = m.SendCmd(254, 79)
	return
}

func (m *MatrixOrbital) PollKeyPress() (key byte, err error) {
	if err = m.SendCmd(254, 38); err != nil {
		key, err = m.Read()
	}
	return
}

func (m *MatrixOrbital) ClearKeyBuffer() (err error) {
	err = m.SendCmd(254, 69)
	return
}

func (m *MatrixOrbital) OutputOn(number byte) (err error) {
	err = m.SendCmd(254, 87, number)
	return
}

func (m *MatrixOrbital) OutputOff(number byte) (err error) {
	err = m.SendCmd(254, 86, number)
	return
}

func (m *MatrixOrbital) GoHome() (err error) {
	err = m.SendCmd(254, 72)
	return
}

func (m *MatrixOrbital) MoveForward() (err error) {
	err = m.SendCmd(254, 77)
	return
}

func (m *MatrixOrbital) MoveBackward() (err error) {
	err = m.SendCmd(254, 76)
	return
}

func (m *MatrixOrbital) SetCursorPosition(col byte, row byte) (err error) {
	if int(col) > m.maxColumns || int(row) > m.maxRows {
		return errors.New("position is out of bound")
	}
	err = m.SendCmd(254, 71, col, row)
	return
}

func (m *MatrixOrbital) SendCmd(data ...byte) (err error) {
	_, err = m.tcpConn.Write(data)
	return
}

func (m *MatrixOrbital) Write(data string) (err error) {
	if len(data) > m.maxColumns*m.maxRows {
		return errors.New("data is larger than the available screen asset")
	}
	_, err = m.tcpConn.Write([]byte(data))
	return
}

func (m *MatrixOrbital) Read() (key byte, err error) {
	buf := make([]byte, 1)
	n, err := m.tcpConn.Read(buf)
	if n > 0 {
		key = buf[0]
	}
	return
}
