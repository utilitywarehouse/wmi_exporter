// +build windows

package collector

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

func init() {
	registerCollector("remote_fx", NewRemoteFx, "RemoteFX Network", "RemoteFX Graphics")
}

// A RemoteFxNetworkCollector is a Prometheus collector for
// WMI Win32_PerfRawData_Counters_RemoteFXNetwork & Win32_PerfRawData_Counters_RemoteFXGraphics metrics
// https://wutils.com/wmi/root/cimv2/win32_perfrawdata_counters_remotefxnetwork/
// https://wutils.com/wmi/root/cimv2/win32_perfrawdata_counters_remotefxgraphics/

type RemoteFxCollector struct {
	// net
	BaseTCPRTT               *prometheus.Desc
	BaseUDPRTT               *prometheus.Desc
	CurrentTCPBandwidth      *prometheus.Desc
	CurrentTCPRTT            *prometheus.Desc
	CurrentUDPBandwidth      *prometheus.Desc
	CurrentUDPRTT            *prometheus.Desc
	FECRate                  *prometheus.Desc
	LossRate                 *prometheus.Desc
	RetransmissionRate       *prometheus.Desc
	TCPReceivedRate          *prometheus.Desc
	TCPSentRate              *prometheus.Desc
	TotalReceivedRate        *prometheus.Desc
	TotalSentRate            *prometheus.Desc
	TotalReceivedBytes       *prometheus.Desc
	TotalSentBytes           *prometheus.Desc
	UDPPacketsReceivedPersec *prometheus.Desc
	UDPPacketsSentPersec     *prometheus.Desc
	UDPReceivedRate          *prometheus.Desc
	UDPSentRate              *prometheus.Desc

	//gfx
	AverageEncodingTime                                *prometheus.Desc
	FrameQuality                                       *prometheus.Desc
	FramesSkippedPerSecondInsufficientClientResources  *prometheus.Desc
	FramesSkippedPerSecondInsufficientNetworkResources *prometheus.Desc
	FramesSkippedPerSecondInsufficientServerResources  *prometheus.Desc
	GraphicsCompressionratio                           *prometheus.Desc
	InputFramesPerSecond                               *prometheus.Desc
	OutputFramesPerSecond                              *prometheus.Desc
	SourceFramesPerSecond                              *prometheus.Desc
}

// NewRemoteFx ...
func NewRemoteFx() (Collector, error) {
	const subsystem = "remote_fx"
	return &RemoteFxCollector{
		// net
		BaseTCPRTT: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_base_tcp_rrt"),
			"Base TCP round-trip time (RTT) detected in milliseconds",
			[]string{"session_name"},
			nil,
		),
		BaseUDPRTT: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_base_udp_rrt"),
			"Base UDP round-trip time (RTT) detected in milliseconds.",
			[]string{"session_name"},
			nil,
		),
		CurrentTCPBandwidth: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_current_tcp_bandwidth"),
			"TCP Bandwidth detected in thousands of bits per second (1000 bps).",
			[]string{"session_name"},
			nil,
		),
		CurrentTCPRTT: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_current_tcp_rtt"),
			"Average TCP round-trip time (RTT) detected in milliseconds.",
			[]string{"session_name"},
			nil,
		),
		CurrentUDPBandwidth: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_current_udp_bandwidth"),
			"UDP Bandwidth detected in thousands of bits per second (1000 bps).",
			[]string{"session_name"},
			nil,
		),
		CurrentUDPRTT: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_current_udp_rtt"),
			"Average UDP round-trip time (RTT) detected in milliseconds.",
			[]string{"session_name"},
			nil,
		),
		FECRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_fec_rate"),
			"Forward Error Correction (FEC) percentage",
			[]string{"session_name"},
			nil,
		),
		LossRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_loss_rate"),
			"Loss percentage",
			[]string{"session_name"},
			nil,
		),
		RetransmissionRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_retransmission_rate"),
			"Percentage of packets that have been retransmitted",
			[]string{"session_name"},
			nil,
		),
		TCPReceivedRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_tcp_received_rate"),
			"Rate in bits per second (bps) at which data is received over TCP.",
			[]string{"session_name"},
			nil,
		),
		TCPSentRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_tcp_sent_rate"),
			"Rate in bits per second (bps) at which data is sent over TCP.",
			[]string{"session_name"},
			nil,
		),
		TotalReceivedRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_total_received_rate"),
			"Rate in bits per second (bps) at which data is received.",
			[]string{"session_name"},
			nil,
		),
		TotalSentRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_total_sent_rate"),
			"Rate in bits per second (bps) at which data is sent.",
			[]string{"session_name"},
			nil,
		),
		TotalReceivedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_total_received_bytes"),
			"(TotalReceivedBytes)",
			[]string{"session_name"},
			nil,
		),
		TotalSentBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_total_sent_bytes"),
			"(TotalSentBytes)",
			[]string{"session_name"},
			nil,
		),
		UDPPacketsReceivedPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_udp_packets_received_persec"),
			"Rate in packets per second at which packets are received over UDP.",
			[]string{"session_name"},
			nil,
		),
		UDPPacketsSentPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_udp_packets_sent_persec"),
			"Rate in packets per second at which packets are sent over UDP.",
			[]string{"session_name"},
			nil,
		),
		UDPReceivedRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_udp_received_rate"),
			"Rate in bits per second (bps) at which data is received over UDP.",
			[]string{"session_name"},
			nil,
		),
		UDPSentRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "net_udp_sent_rate"),
			"Rate in bits per second (bps) at which data is sent over UDP.",
			[]string{"session_name"},
			nil,
		),

		//gfx
		AverageEncodingTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "gfx_average_encoding_time"),
			"Average frame encoding time in milliseconds",
			[]string{"session_name"},
			nil,
		),
		FrameQuality: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "gfx_frame_quality"),
			"Quality of the output frame expressed as a percentage of the quality of the source frame.",
			[]string{"session_name"},
			nil,
		),
		FramesSkippedPerSecondInsufficientClientResources: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "gfx_frames_skipped_persec_insufficient_clt_res"),
			"Number of frames skipped per second due to insufficient client resources.",
			[]string{"session_name"},
			nil,
		),
		FramesSkippedPerSecondInsufficientNetworkResources: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "gfx_frames_skipped_persec_insufficient_net_res"),
			"Number of frames skipped per second due to insufficient network resources.",
			[]string{"session_name"},
			nil,
		),
		FramesSkippedPerSecondInsufficientServerResources: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "gfx_frames_skipped_persec_insufficient_srv_res"),
			"Number of frames skipped per second due to insufficient server resources.",
			[]string{"session_name"},
			nil,
		),
		GraphicsCompressionratio: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "gfx_graphics_compression_ratio"),
			"Ratio of the number of bytes encoded to the number of bytes input.",
			[]string{"session_name"},
			nil,
		),
		InputFramesPerSecond: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "gfx_input_frames_persec"),
			"Number of sources frames provided as input to RemoteFX graphics per second.",
			[]string{"session_name"},
			nil,
		),
		OutputFramesPerSecond: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "gfx_output_frames_persec"),
			"Number of frames sent to the client per second.",
			[]string{"session_name"},
			nil,
		),
		SourceFramesPerSecond: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "gfx_source_frames_persec"),
			"Number of frames composed by the source (DWM) per second.",
			[]string{"session_name"},
			nil,
		),
	}, nil
}

// Collect sends the metric values for each metric
// to the provided prometheus Metric channel.
func (c *RemoteFxCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	if desc, err := c.collectRemoteFXNetworkCount(ctx, ch); err != nil {
		log.Error("failed collecting terminal services session count metrics:", desc, err)
		return err
	}
	if desc, err := c.collectRemoteFXGraphicsCounters(ctx, ch); err != nil {
		log.Error("failed collecting terminal services session count metrics:", desc, err)
		return err
	}
	return nil
}

type perflibRemoteFxNetwork struct {
	Name                     string
	BaseTCPRTT               float64 `perflib:"Base TCP RTT"`
	BaseUDPRTT               float64 `perflib:"Base UDP RTT"`
	CurrentTCPBandwidth      float64 `perflib:"Current TCP Bandwidth"`
	CurrentTCPRTT            float64 `perflib:"Current TCP RTT"`
	CurrentUDPBandwidth      float64 `perflib:"Current UDP Bandwidth"`
	CurrentUDPRTT            float64 `perflib:"Current UDP RTT"`
	FECRate                  float64 `perflib:"FEC Rate"`
	LossRate                 float64 `perflib:"Loss Rate"`
	RetransmissionRate       float64 `perflib:"Retransmission Rate"`
	TCPReceivedRate          float64 `perflib:"TCP Received Rate"`
	TCPSentRate              float64 `perflib:"TCP Sent Rate"`
	TotalReceivedRate        float64 `perflib:"Total Received Rate"`
	TotalSentRate            float64 `perflib:"Total Sent Rate"`
	TotalReceivedBytes       float64 `perflib:"Total Received Bytes"`
	TotalSentBytes           float64 `perflib:"Total Sent Bytes"`
	UDPPacketsReceivedPersec float64 `perflib:"UDP Packets Received/sec"`
	UDPPacketsSentPersec     float64 `perflib:"UDP Packets Sent/sec"`
	UDPReceivedRate          float64 `perflib:"UDP Received Rate"`
	UDPSentRate              float64 `perflib:"UDP Sent Rate"`
}

func (c *RemoteFxCollector) collectRemoteFXNetworkCount(ctx *ScrapeContext, ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	dst := make([]perflibRemoteFxNetwork, 0)
	err := unmarshalObject(ctx.perfObjects["RemoteFX Network"], &dst)
	if err != nil {
		return nil, err
	}

	for _, d := range dst {
		// only connect metrics for remote named sessions
		n := strings.ToLower(d.Name)
		if n == "" || n == "services" || n == "console" {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			c.BaseTCPRTT,
			prometheus.GaugeValue,
			d.BaseTCPRTT,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.BaseUDPRTT,
			prometheus.GaugeValue,
			d.BaseUDPRTT,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.CurrentTCPBandwidth,
			prometheus.GaugeValue,
			d.CurrentTCPBandwidth,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.CurrentTCPRTT,
			prometheus.GaugeValue,
			d.CurrentTCPRTT,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.CurrentUDPBandwidth,
			prometheus.GaugeValue,
			d.CurrentUDPBandwidth,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.CurrentUDPRTT,
			prometheus.GaugeValue,
			d.CurrentUDPRTT,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.FECRate,
			prometheus.CounterValue,
			d.FECRate,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.LossRate,
			prometheus.CounterValue,
			d.LossRate,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.RetransmissionRate,
			prometheus.CounterValue,
			d.RetransmissionRate,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.TCPReceivedRate,
			prometheus.CounterValue,
			d.TCPReceivedRate,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.TCPSentRate,
			prometheus.CounterValue,
			d.TCPSentRate,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.TotalReceivedRate,
			prometheus.CounterValue,
			d.TotalReceivedRate,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.TotalSentRate,
			prometheus.CounterValue,
			d.TotalSentRate,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.TotalReceivedBytes,
			prometheus.CounterValue,
			d.TotalReceivedBytes,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.TotalSentBytes,
			prometheus.CounterValue,
			d.TotalSentBytes,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.UDPPacketsReceivedPersec,
			prometheus.CounterValue,
			d.UDPPacketsReceivedPersec,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.UDPPacketsSentPersec,
			prometheus.CounterValue,
			d.UDPPacketsSentPersec,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.UDPReceivedRate,
			prometheus.CounterValue,
			d.UDPReceivedRate,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.UDPSentRate,
			prometheus.CounterValue,
			d.UDPSentRate,
			d.Name,
		)
	}
	return nil, nil
}

type perflibRemoteFxGraphics struct {
	Name                                               string
	AverageEncodingTime                                float64 `perflib:"Average Encoding Time"`
	FrameQuality                                       float64 `perflib:"Frame Quality"`
	FramesSkippedPerSecondInsufficientClientResources  float64 `perflib:"Frames Skipped/Second - Insufficient Server Resources"`
	FramesSkippedPerSecondInsufficientNetworkResources float64 `perflib:"Frames Skipped/Second - Insufficient Network Resources"`
	FramesSkippedPerSecondInsufficientServerResources  float64 `perflib:"Frames Skipped/Second - Insufficient Client Resources"`
	GraphicsCompressionratio                           float64 `perflib:"Graphics Compression ratio"`
	InputFramesPerSecond                               float64 `perflib:"Input Frames/Second"`
	OutputFramesPerSecond                              float64 `perflib:"Output Frames/Second"`
	SourceFramesPerSecond                              float64 `perflib:"Source Frames/Second"`
}

func (c *RemoteFxCollector) collectRemoteFXGraphicsCounters(ctx *ScrapeContext, ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	dst := make([]perflibRemoteFxGraphics, 0)
	err := unmarshalObject(ctx.perfObjects["RemoteFX Graphics"], &dst)
	if err != nil {
		return nil, err
	}

	for _, d := range dst {
		// only connect metrics for remote named sessions
		n := strings.ToLower(d.Name)
		if n == "" || n == "services" || n == "console" {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			c.AverageEncodingTime,
			prometheus.GaugeValue,
			d.AverageEncodingTime,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.FrameQuality,
			prometheus.GaugeValue,
			d.FrameQuality,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.FramesSkippedPerSecondInsufficientClientResources,
			prometheus.CounterValue,
			d.FramesSkippedPerSecondInsufficientClientResources,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.FramesSkippedPerSecondInsufficientNetworkResources,
			prometheus.CounterValue,
			d.FramesSkippedPerSecondInsufficientNetworkResources,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.FramesSkippedPerSecondInsufficientServerResources,
			prometheus.CounterValue,
			d.FramesSkippedPerSecondInsufficientServerResources,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.GraphicsCompressionratio,
			prometheus.GaugeValue,
			d.GraphicsCompressionratio,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.InputFramesPerSecond,
			prometheus.CounterValue,
			d.InputFramesPerSecond,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.OutputFramesPerSecond,
			prometheus.CounterValue,
			d.OutputFramesPerSecond,
			d.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.SourceFramesPerSecond,
			prometheus.CounterValue,
			d.SourceFramesPerSecond,
			d.Name,
		)
	}

	return nil, nil
}
