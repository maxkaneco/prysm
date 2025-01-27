// Package params defines all custom parameter configurations
// for running end to end tests.
package params

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/ethereum/go-ethereum/common"
)

// params struct defines the parameters needed for running E2E tests to properly handle test sharding.
type params struct {
	TestPath                  string
	LogPath                   string
	TestShardIndex            int
	BeaconNodeCount           int
	LighthouseBeaconNodeCount int
	ContractAddress           common.Address
	Ports                     *ports
}

type ports struct {
	BootNodePort                    int
	BootNodeMetricsPort             int
	Eth1Port                        int
	Eth1RPCPort                     int
	Eth1WSPort                      int
	PrysmBeaconNodeRPCPort          int
	PrysmBeaconNodeUDPPort          int
	PrysmBeaconNodeTCPPort          int
	PrysmBeaconNodeGatewayPort      int
	PrysmBeaconNodeMetricsPort      int
	PrysmBeaconNodePprofPort        int
	LighthouseBeaconNodeP2PPort     int
	LighthouseBeaconNodeHTTPPort    int
	LighthouseBeaconNodeMetricsPort int
	ValidatorMetricsPort            int
	ValidatorGatewayPort            int
	JaegerTracingPort               int
}

// TestParams is the globally accessible var for getting config elements.
var TestParams *params

// BootNodeLogFileName is the file name used for the beacon chain node logs.
var BootNodeLogFileName = "bootnode.log"

// TracingRequestSinkFileName is the file name for writing raw trace requests.
var TracingRequestSinkFileName = "tracing-http-requests.log.gz"

// BeaconNodeLogFileName is the file name used for the beacon chain node logs.
var BeaconNodeLogFileName = "beacon-%d.log"

// ValidatorLogFileName is the file name used for the validator client logs.
var ValidatorLogFileName = "vals-%d.log"

// StandardBeaconCount is a global constant for the count of beacon nodes of standard E2E tests.
var StandardBeaconCount = 2

// StandardLighthouseNodeCount is a global constant for the count of lighthouse beacon nodes of standard E2E tests.
var StandardLighthouseNodeCount = 2

// DepositCount is the amount of deposits E2E makes on a separate validator client.
var DepositCount = uint64(64)

// Base port values.
const (
	portSpan = 50

	BootNodePort        = 2150
	BootNodeMetricsPort = BootNodePort + portSpan

	Eth1Port    = 3150
	Eth1RPCPort = Eth1Port + portSpan
	Eth1WSPort  = Eth1Port + 2*portSpan

	PrysmBeaconNodeRPCPort     = 4150
	PrysmBeaconNodeUDPPort     = PrysmBeaconNodeRPCPort + portSpan
	PrysmBeaconNodeTCPPort     = PrysmBeaconNodeRPCPort + 2*portSpan
	PrysmBeaconNodeGatewayPort = PrysmBeaconNodeRPCPort + 3*portSpan
	PrysmBeaconNodeMetricsPort = PrysmBeaconNodeRPCPort + 4*portSpan
	PrysmBeaconNodePprofPort   = PrysmBeaconNodeRPCPort + 5*portSpan

	LighthouseBeaconNodeP2PPort     = 5150
	LighthouseBeaconNodeHTTPPort    = LighthouseBeaconNodeP2PPort + portSpan
	LighthouseBeaconNodeMetricsPort = LighthouseBeaconNodeP2PPort + 2*portSpan

	ValidatorGatewayPort = 6150
	ValidatorMetricsPort = ValidatorGatewayPort + portSpan

	JaegerTracingPort = 9150
)

// Init initializes the E2E config, properly handling test sharding.
func Init(beaconNodeCount int) error {
	testPath := bazel.TestTmpDir()
	logPath, ok := os.LookupEnv("TEST_UNDECLARED_OUTPUTS_DIR")
	if !ok {
		return errors.New("expected TEST_UNDECLARED_OUTPUTS_DIR to be defined")
	}
	testTotalShardsStr, ok := os.LookupEnv("TEST_TOTAL_SHARDS")
	if !ok {
		testTotalShardsStr = "1"
	}
	testTotalShards, err := strconv.Atoi(testTotalShardsStr)
	if err != nil {
		return err
	}
	testShardIndexStr, ok := os.LookupEnv("TEST_SHARD_INDEX")
	if !ok {
		testShardIndexStr = "0"
	}
	testShardIndex, err := strconv.Atoi(testShardIndexStr)
	if err != nil {
		return err
	}

	var existingRegistrations []int
	bootnodePort, err := port(BootNodePort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	bootnodeMetricsPort, err := port(BootNodeMetricsPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	eth1Port, err := port(Eth1Port, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	eth1RPCPort, err := port(Eth1RPCPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	eth1WSPort, err := port(Eth1WSPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	beaconNodeRPCPort, err := port(PrysmBeaconNodeRPCPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	beaconNodeUDPPort, err := port(PrysmBeaconNodeUDPPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	beaconNodeTCPPort, err := port(PrysmBeaconNodeTCPPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	beaconNodeGatewayPort, err := port(PrysmBeaconNodeGatewayPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	beaconNodeMetricsPort, err := port(PrysmBeaconNodeMetricsPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	beaconNodePprofPort, err := port(PrysmBeaconNodePprofPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	validatorGatewayPort, err := port(ValidatorGatewayPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	validatorMetricsPort, err := port(ValidatorMetricsPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	jaegerTracingPort, err := port(JaegerTracingPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	testPorts := &ports{
		BootNodePort:               bootnodePort,
		BootNodeMetricsPort:        bootnodeMetricsPort,
		Eth1Port:                   eth1Port,
		Eth1RPCPort:                eth1RPCPort,
		Eth1WSPort:                 eth1WSPort,
		PrysmBeaconNodeRPCPort:     beaconNodeRPCPort,
		PrysmBeaconNodeUDPPort:     beaconNodeUDPPort,
		PrysmBeaconNodeTCPPort:     beaconNodeTCPPort,
		PrysmBeaconNodeGatewayPort: beaconNodeGatewayPort,
		PrysmBeaconNodeMetricsPort: beaconNodeMetricsPort,
		PrysmBeaconNodePprofPort:   beaconNodePprofPort,
		ValidatorMetricsPort:       validatorMetricsPort,
		ValidatorGatewayPort:       validatorGatewayPort,
		JaegerTracingPort:          jaegerTracingPort,
	}

	TestParams = &params{
		TestPath:        filepath.Join(testPath, fmt.Sprintf("shard-%d", testShardIndex)),
		LogPath:         logPath,
		TestShardIndex:  testShardIndex,
		BeaconNodeCount: beaconNodeCount,
		Ports:           testPorts,
	}
	return nil
}

// InitMultiClient initializes the multiclient E2E config, properly handling test sharding.
func InitMultiClient(beaconNodeCount int, lighthouseNodeCount int) error {
	testPath := bazel.TestTmpDir()
	logPath, ok := os.LookupEnv("TEST_UNDECLARED_OUTPUTS_DIR")
	if !ok {
		return errors.New("expected TEST_UNDECLARED_OUTPUTS_DIR to be defined")
	}
	testTotalShardsStr, ok := os.LookupEnv("TEST_TOTAL_SHARDS")
	if !ok {
		testTotalShardsStr = "1"
	}
	testTotalShards, err := strconv.Atoi(testTotalShardsStr)
	if err != nil {
		return err
	}
	testShardIndexStr, ok := os.LookupEnv("TEST_SHARD_INDEX")
	if !ok {
		testShardIndexStr = "0"
	}
	testShardIndex, err := strconv.Atoi(testShardIndexStr)
	if err != nil {
		return err
	}

	var existingRegistrations []int
	bootnodePort, err := port(BootNodePort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	bootnodeMetricsPort, err := port(BootNodeMetricsPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	eth1Port, err := port(Eth1Port, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	eth1RPCPort, err := port(Eth1RPCPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	eth1WSPort, err := port(Eth1WSPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	prysmBeaconNodeRPCPort, err := port(PrysmBeaconNodeRPCPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	prysmBeaconNodeUDPPort, err := port(PrysmBeaconNodeUDPPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	prysmBeaconNodeTCPPort, err := port(PrysmBeaconNodeTCPPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	prysmBeaconNodeGatewayPort, err := port(PrysmBeaconNodeGatewayPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	prysmBeaconNodeMetricsPort, err := port(PrysmBeaconNodeMetricsPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	prysmBeaconNodePprofPort, err := port(PrysmBeaconNodePprofPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	lighthouseBeaconNodeP2PPort, err := port(LighthouseBeaconNodeP2PPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	lighthouseBeaconNodeHTTPPort, err := port(LighthouseBeaconNodeHTTPPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	lighthouseBeaconNodeMetricsPort, err := port(LighthouseBeaconNodeMetricsPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	validatorGatewayPort, err := port(ValidatorGatewayPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	validatorMetricsPort, err := port(ValidatorMetricsPort, testTotalShards, testShardIndex, &existingRegistrations)
	if err != nil {
		return err
	}
	testPorts := &ports{
		BootNodePort:                    bootnodePort,
		BootNodeMetricsPort:             bootnodeMetricsPort,
		Eth1Port:                        eth1Port,
		Eth1RPCPort:                     eth1RPCPort,
		Eth1WSPort:                      eth1WSPort,
		PrysmBeaconNodeRPCPort:          prysmBeaconNodeRPCPort,
		PrysmBeaconNodeUDPPort:          prysmBeaconNodeUDPPort,
		PrysmBeaconNodeTCPPort:          prysmBeaconNodeTCPPort,
		PrysmBeaconNodeGatewayPort:      prysmBeaconNodeGatewayPort,
		PrysmBeaconNodeMetricsPort:      prysmBeaconNodeMetricsPort,
		PrysmBeaconNodePprofPort:        prysmBeaconNodePprofPort,
		LighthouseBeaconNodeP2PPort:     lighthouseBeaconNodeP2PPort,
		LighthouseBeaconNodeHTTPPort:    lighthouseBeaconNodeHTTPPort,
		LighthouseBeaconNodeMetricsPort: lighthouseBeaconNodeMetricsPort,
		ValidatorMetricsPort:            validatorMetricsPort,
		ValidatorGatewayPort:            validatorGatewayPort,
	}

	TestParams = &params{
		TestPath:                  filepath.Join(testPath, fmt.Sprintf("shard-%d", testShardIndex)),
		LogPath:                   logPath,
		TestShardIndex:            testShardIndex,
		BeaconNodeCount:           beaconNodeCount,
		LighthouseBeaconNodeCount: lighthouseNodeCount,
		Ports:                     testPorts,
	}
	return nil
}

// port returns a safe port number based on the seed and shard data.
func port(seed, shardCount, shardIndex int, existingRegistrations *[]int) (int, error) {
	portToRegister := seed + portSpan/shardCount*shardIndex
	for _, p := range *existingRegistrations {
		if portToRegister >= p && portToRegister <= p+(portSpan/shardCount)-1 {
			return 0, fmt.Errorf("port %d overlaps with already registered port %d", seed, p)
		}
	}
	*existingRegistrations = append(*existingRegistrations, portToRegister)

	// Calculation example: 3 shards, seed 2000, port span 50.
	// Shard 0: 2000 + (50 / 3 * 0) = 2000 (we can safely use ports 2000-2015)
	// Shard 1: 2000 + (50 / 3 * 1) = 2016 (we can safely use ports 2016-2031)
	// Shard 2: 2000 + (50 / 3 * 2) = 2032 (we can safely use ports 2032-2047, and in reality 2032-2049)
	return portToRegister, nil
}
