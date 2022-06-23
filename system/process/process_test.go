package process

import (
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/isyscore/isc-gobase/system/common"
)

var mu sync.Mutex

func skipIfNotImplementedErr(t *testing.T, err error) {
	if err == common.ErrNotImplementedError {
		t.Skip("not implemented")
	}
}

func testGetProcess() Process {
	checkPid := os.Getpid() // process.test
	ret, _ := NewProcess(int32(checkPid))
	return *ret
}

func Test_Pids(t *testing.T) {
	ret, err := Pids()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("error %v", err)
	}
	if len(ret) == 0 {
		t.Errorf("could not get pids %v", ret)
	}
}

func Test_Pid_exists(t *testing.T) {
	checkPid := os.Getpid()

	ret, err := PidExists(int32(checkPid))
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("error %v", err)
	}

	if ret == false {
		t.Errorf("could not get process exists: %v", ret)
	}
}

func Test_NewProcess(t *testing.T) {
	checkPid := os.Getpid()

	ret, err := NewProcess(int32(checkPid))
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("error %v", err)
	}
	empty := &Process{}
	if runtime.GOOS != "windows" { // Windows pid is 0
		if empty == ret {
			t.Errorf("error %v", ret)
		}
	}

}

func Test_Process_MemoryInfo(t *testing.T) {
	p := testGetProcess()

	v, err := p.MemoryInfo()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("getting memory info error %v", err)
	}
	empty := MemoryInfoStat{}
	if v == nil || *v == empty {
		t.Errorf("could not get memory info %v", v)
	}
}

func Test_Process_CmdLine(t *testing.T) {
	p := testGetProcess()

	v, err := p.Cmdline()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("getting cmdline error %v", err)
	}
	if !strings.Contains(v, "process.test") {
		t.Errorf("invalid cmd line %v", v)
	}
}

func Test_Process_CmdLineSlice(t *testing.T) {
	p := testGetProcess()

	v, err := p.CmdlineSlice()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Fatalf("getting cmdline slice error %v", err)
	}
	if !reflect.DeepEqual(v, os.Args) {
		t.Errorf("returned cmdline slice not as expected:\nexp: %v\ngot: %v", os.Args, v)
	}
}

func Test_Process_Ppid(t *testing.T) {
	p := testGetProcess()

	v, err := p.Ppid()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("getting ppid error %v", err)
	}
	if v == 0 {
		t.Errorf("return value is 0 %v", v)
	}
	expected := os.Getppid()
	if v != int32(expected) {
		t.Errorf("return value is %v, expected %v", v, expected)
	}
}

func Test_Process_Status(t *testing.T) {
	p := testGetProcess()

	v, err := p.Status()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("getting status error %v", err)
	}
	if v != "R" && v != "S" {
		t.Errorf("could not get state %v", v)
	}
}

func Test_Process_Nice(t *testing.T) {
	p := testGetProcess()

	n, err := p.Nice()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("getting nice error %v", err)
	}
	if runtime.GOOS != "windows" && n != 0 && n != 20 && n != 8 {
		t.Errorf("invalid nice: %d", n)
	}
}

func Test_Process_NumThread(t *testing.T) {
	p := testGetProcess()

	n, err := p.NumThreads()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("getting NumThread error %v", err)
	}
	if n < 0 {
		t.Errorf("invalid NumThread: %d", n)
	}
}

func Test_Process_Name(t *testing.T) {
	p := testGetProcess()

	n, err := p.Name()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("getting name error %v", err)
	}
	if !strings.Contains(n, "process.test") {
		t.Errorf("invalid Exe %s", n)
	}
}

func Test_Process_Exe(t *testing.T) {
	p := testGetProcess()

	n, err := p.Exe()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("getting Exe error %v", err)
	}
	if !strings.Contains(n, "process.test") {
		t.Errorf("invalid Exe %s", n)
	}
}

func Test_Process_CpuPercent(t *testing.T) {
	p := testGetProcess()
	_, err := p.Percent(0)
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("error %v", err)
	}
	duration := time.Duration(1000) * time.Microsecond
	time.Sleep(duration)
	percent, err := p.Percent(0)
	if err != nil {
		t.Errorf("error %v", err)
	}

	numcpu := runtime.NumCPU()
	//	if percent < 0.0 || percent > 100.0*float64(numcpu) { // TODO
	if percent < 0.0 {
		t.Fatalf("CPUPercent value is invalid: %f, %d", percent, numcpu)
	}
}

func Test_Process_CpuPercentLoop(t *testing.T) {
	p := testGetProcess()
	numcpu := runtime.NumCPU()

	for i := 0; i < 2; i++ {
		duration := time.Duration(100) * time.Microsecond
		percent, err := p.Percent(duration)
		skipIfNotImplementedErr(t, err)
		if err != nil {
			t.Errorf("error %v", err)
		}
		//	if percent < 0.0 || percent > 100.0*float64(numcpu) { // TODO
		if percent < 0.0 {
			t.Fatalf("CPUPercent value is invalid: %f, %d", percent, numcpu)
		}
	}
}

func Test_Process_CreateTime(t *testing.T) {
	if os.Getenv("CIRCLECI") == "true" {
		t.Skip("Skip CI")
	}

	p := testGetProcess()

	c, err := p.CreateTime()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("error %v", err)
	}

	if c < 1420000000 {
		t.Errorf("process created time is wrong.")
	}

	gotElapsed := time.Since(time.Unix(int64(c/1000), 0))
	maxElapsed := time.Duration(20 * time.Second)

	if gotElapsed >= maxElapsed {
		t.Errorf("this process has not been running for %v", gotElapsed)
	}
}

func Test_Parent(t *testing.T) {
	p := testGetProcess()

	c, err := p.Parent()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Fatalf("error %v", err)
	}
	if c == nil {
		t.Fatalf("could not get parent")
	}
	if c.Pid == 0 {
		t.Fatalf("wrong parent pid")
	}
}

func Test_Connections(t *testing.T) {
	p := testGetProcess()

	addr, err := net.ResolveTCPAddr("tcp", "localhost:0") // dynamically get a random open port from OS
	if err != nil {
		t.Fatalf("unable to resolve localhost: %v", err)
	}
	l, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		t.Fatalf("unable to listen on %v: %v", addr, err)
	}
	defer l.Close()

	tcpServerAddr := l.Addr().String()
	tcpServerAddrIP := strings.Split(tcpServerAddr, ":")[0]
	tcpServerAddrPort, err := strconv.ParseUint(strings.Split(tcpServerAddr, ":")[1], 10, 32)
	if err != nil {
		t.Fatalf("unable to parse tcpServerAddr port: %v", err)
	}

	serverEstablished := make(chan struct{})
	go func() { // TCP listening goroutine
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		serverEstablished <- struct{}{}
		_, err = ioutil.ReadAll(conn)
		if err != nil {
			panic(err)
		}
	}()

	conn, err := net.Dial("tcp", tcpServerAddr)
	if err != nil {
		t.Fatalf("unable to dial %v: %v", tcpServerAddr, err)
	}
	defer conn.Close()

	// Rarely the call to net.Dial returns before the server connection is
	// established. Wait so that the test doesn't fail.
	<-serverEstablished

	c, err := p.Connections()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Fatalf("error %v", err)
	}
	if len(c) == 0 {
		t.Fatal("no connections found")
	}

	serverConnections := 0
	for _, connection := range c {
		if connection.Laddr.IP == tcpServerAddrIP && connection.Laddr.Port == uint32(tcpServerAddrPort) && connection.Raddr.Port != 0 {
			if connection.Status != "ESTABLISHED" {
				t.Fatalf("expected server connection to be ESTABLISHED, have %+v", connection)
			}
			serverConnections++
		}
	}

	clientConnections := 0
	for _, connection := range c {
		if connection.Raddr.IP == tcpServerAddrIP && connection.Raddr.Port == uint32(tcpServerAddrPort) {
			if connection.Status != "ESTABLISHED" {
				t.Fatalf("expected client connection to be ESTABLISHED, have %+v", connection)
			}
			clientConnections++
		}
	}

	if serverConnections != 1 { // two established connections, one for the server, the other for the client
		t.Fatalf("expected 1 server connection, have %d.\nDetails: %+v", serverConnections, c)
	}

	if clientConnections != 1 { // two established connections, one for the server, the other for the client
		t.Fatalf("expected 1 server connection, have %d.\nDetails: %+v", clientConnections, c)
	}
}

func Test_AllProcesses_cmdLine(t *testing.T) {
	procs, err := Processes()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Fatalf("getting processes error %v", err)
	}
	for _, proc := range procs {
		var exeName string
		var cmdLine string

		exeName, _ = proc.Exe()
		cmdLine, err = proc.Cmdline()
		if err != nil {
			cmdLine = "Error: " + err.Error()
		}

		t.Logf("Process #%v: Name: %v / CmdLine: %v\n", proc.Pid, exeName, cmdLine)
	}
}

func Test_AllProcesses_environ(t *testing.T) {
	procs, err := Processes()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Fatalf("getting processes error %v", err)
	}
	for _, proc := range procs {
		exeName, _ := proc.Exe()
		environ, err := proc.Environ()
		if err != nil {
			environ = []string{"Error: " + err.Error()}
		}

		t.Logf("Process #%v: Name: %v / Environment Variables: %v\n", proc.Pid, exeName, environ)
	}
}

func Test_AllProcesses_Cwd(t *testing.T) {
	procs, err := Processes()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Fatalf("getting processes error %v", err)
	}
	for _, proc := range procs {
		exeName, _ := proc.Exe()
		cwd, err := proc.Cwd()
		if err != nil {
			cwd = "Error: " + err.Error()
		}

		t.Logf("Process #%v: Name: %v / Current Working Directory: %s\n", proc.Pid, exeName, cwd)
	}
}

func BenchmarkNewProcess(b *testing.B) {
	checkPid := os.Getpid()
	for i := 0; i < b.N; i++ {
		NewProcess(int32(checkPid))
	}
}

func BenchmarkProcessName(b *testing.B) {
	p := testGetProcess()
	for i := 0; i < b.N; i++ {
		p.Name()
	}
}

func BenchmarkProcessPpid(b *testing.B) {
	p := testGetProcess()
	for i := 0; i < b.N; i++ {
		p.Ppid()
	}
}
