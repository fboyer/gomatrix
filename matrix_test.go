package gomatrix

import (
	"flag"
	gc "launchpad.net/gocheck"
	"testing"
	"time"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

var live = flag.Bool("live", false, "Include live tests")

// These tests require a live environment
type LiveMatrixOrbitalSuite struct {
	display *MatrixOrbital
}

var _ = gc.Suite(&LiveMatrixOrbitalSuite{})

func (l *LiveMatrixOrbitalSuite) SetUpSuite(c *gc.C) {
	if !*live {
		c.Skip("-live not provided")
	}

	l.display = CreateMatrixOrbital("192.168.1.237", 4002, 20, 2)
	if l.display == nil {
		c.Skip("unable to create matrixorbital")
	}
}

func (l *LiveMatrixOrbitalSuite) SetUpTest(c *gc.C) {
	l.display.Connect()
}

func (l *LiveMatrixOrbitalSuite) TearDownTest(c *gc.C) {
	l.display.Disconnect()
}

func (l *LiveMatrixOrbitalSuite) TearDownSuite(c *gc.C) {
}

func (l *LiveMatrixOrbitalSuite) TestConnect(c *gc.C) {
	c.Assert(l.display.tcpConn, gc.NotNil)
}

func (l *LiveMatrixOrbitalSuite) TestDisconnect(c *gc.C) {
	err := l.display.Disconnect()
	c.Assert(err, gc.IsNil)
}

func (l *LiveMatrixOrbitalSuite) TestClearScreen(c *gc.C) {
	l.display.Write("This is a test")
	err := l.display.ClearScreen()
	c.Assert(err, gc.IsNil)
}

func (l *LiveMatrixOrbitalSuite) TestOutputOnOff(c *gc.C) {
	err0 := l.display.OutputOn(3)
	time.Sleep(1 * time.Second)
	err1 := l.display.OutputOff(3)
	c.Assert(err0, gc.IsNil)
	c.Assert(err1, gc.IsNil)
}

func (l *LiveMatrixOrbitalSuite) TestCursorMovement(c *gc.C) {
	err0 := l.display.MoveForward()
	l.display.Write("1")
	err1 := l.display.SetCursorPosition(5, 2)
	l.display.Write("2")
	l.display.MoveBackward()
	err2 := l.display.MoveBackward()
	l.display.Write("3")
	err3 := l.display.GoHome()
	l.display.Write("4")
	l.display.ClearScreen()
	c.Assert(err0, gc.IsNil)
	c.Assert(err1, gc.IsNil)
	c.Assert(err2, gc.IsNil)
	c.Assert(err3, gc.IsNil)
}

func (l *LiveMatrixOrbitalSuite) TestPollKeyPress(c *gc.C) {
	err0 := l.display.AutoTxKeyPressesOff()
	key, err1 := l.display.PollKeyPress()
	c.Assert(err0, gc.IsNil)
	c.Assert(key, gc.Equals, byte(0))
	c.Assert(err1, gc.IsNil)
}
