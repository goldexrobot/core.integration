package backend

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/goldexrobot/core.proto/api/evalpb"
	"github.com/goldexrobot/core.proto/api/integrationpb"
	"github.com/goldexrobot/core.proto/api/storagepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type Client struct {
	addr     string
	certFile string
	keyFile  string
	caFile   string

	conn *grpc.ClientConn
}

func NewClient(addr, cert, key, ca string) *Client {
	return &Client{
		addr:     addr,
		certFile: cert,
		keyFile:  key,
		caFile:   ca,
	}
}

func (c *Client) Connect() (err error) {
	// read cert pair
	tlsCert, err := tls.LoadX509KeyPair(c.certFile, c.keyFile)
	if err != nil {
		return fmt.Errorf("loading tls cert pair: %w", err)
	}

	// cas
	caPool := x509.NewCertPool()
	{
		f, err := os.OpenFile(c.caFile, os.O_RDONLY, 0400)
		if err != nil {
			return fmt.Errorf("opening ca certs file: %w", err)
		}
		b, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			return fmt.Errorf("reading ca certs file: %w", err)
		}
		if !caPool.AppendCertsFromPEM(b) {
			return fmt.Errorf("failed loading ca certs")
		}
	}

	// transport creds
	creds := credentials.NewTLS(
		&tls.Config{
			Certificates:             []tls.Certificate{tlsCert},
			RootCAs:                  caPool,
			PreferServerCipherSuites: true,
		},
	)

	// backoff config
	boConfig := backoff.DefaultConfig
	boConfig.MaxDelay = time.Second * 15

	// dial opts
	dopts := []grpc.DialOption{
		grpc.WithWriteBufferSize(64 * 1024),
		grpc.WithReadBufferSize(64 * 1024),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 60,
			Timeout:             time.Second * 30,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1 * 1024 * 1024),
		),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: boConfig,
		}),
		grpc.WithUserAgent("Integration Emulator"),
		grpc.WithTransportCredentials(creds),
	}

	// dial (non blocking)
	conn, err := grpc.Dial(c.addr, dopts...)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *Client) OccupiedCells(ctx context.Context) (domains map[string][]string, err error) {
	cli := storagepb.NewStorageClient(c.conn)
	res, err := cli.Occupied(ctx, &emptypb.Empty{}, grpc.WaitForReady(true))
	if err != nil {
		return
	}
	domains = make(map[string][]string, len(res.GetDomains()))
	for domain, cells := range res.GetDomains() {
		domains[storagepb.Domain(domain).String()] = cells.GetCells()
	}
	return
}

func (c *Client) NewEval(ctx context.Context) (id uint64, err error) {
	cli := evalpb.NewEvaluationClient(c.conn)
	res, err := cli.Begin(ctx, &emptypb.Empty{}, grpc.WaitForReady(true))
	if err != nil {
		return
	}
	id = res.GetEvalId()
	return
}

func (c *Client) EvaluateSpectrum(ctx context.Context, evalID uint64, spectrum map[string]float64) (r ResultFineness, rejected bool, err error) {
	cli := evalpb.NewEvaluationClient(c.conn)
	res, err := cli.Spectrum(ctx, &evalpb.SpectrumModel{EvalId: evalID, Spectrum: spectrum}, grpc.WaitForReady(true))
	if err != nil {
		return
	}
	if res.GetReject() != nil {
		rejected = true
		return
	}
	r = ResultFineness{
		Alloy:      res.GetFineness().GetAlloy(),
		Purity:     res.GetFineness().GetPurity(),
		Millesimal: int(res.GetFineness().GetMillesimal()),
		Carat:      res.GetFineness().GetCarat(),
		Confidence: res.GetFineness().GetConfidence(),
		Risky:      res.GetFineness().GetRisky(),
	}
	return
}

func (c *Client) EvaluateHydro(ctx context.Context, evalID uint64, dry, wet float64) (rejected bool, err error) {
	cli := evalpb.NewEvaluationClient(c.conn)
	res1, err := cli.DryWeight(ctx, &evalpb.WeightModel{EvalId: evalID, Weight: dry}, grpc.WaitForReady(true))
	if err != nil {
		return
	}
	if res1.GetReject() != nil {
		rejected = true
		return
	}
	res2, err := cli.WetWeight(ctx, &evalpb.WeightModel{EvalId: evalID, Weight: wet}, grpc.WaitForReady(true))
	if err != nil {
		return
	}
	if res2.GetReject() != nil {
		rejected = true
		return
	}
	return
}

func (c *Client) FinalizeEvaluation(ctx context.Context, evalID uint64) (r ResultFineness, rejected bool, err error) {
	cli := evalpb.NewEvaluationClient(c.conn)
	res, err := cli.Finalize(ctx, &evalpb.FinalizeModel{EvalId: evalID}, grpc.WaitForReady(true))
	if err != nil {
		return
	}
	if res.GetReject() != nil {
		rejected = true
		return
	}
	r = ResultFineness{
		Alloy:      res.GetFineness().GetAlloy(),
		Purity:     res.GetFineness().GetPurity(),
		Millesimal: int(res.GetFineness().GetMillesimal()),
		Carat:      res.GetFineness().GetCarat(),
		Confidence: res.GetFineness().GetConfidence(),
		Risky:      res.GetFineness().GetRisky(),
	}
	return
}

func (c *Client) OccupyStorageCell(ctx context.Context, domain, cell, tx string) (forbidden bool, reason string, err error) {
	dom, ok := storagepb.Domain_value[domain]
	if !ok {
		err = fmt.Errorf("unknown domain %q", domain)
		return
	}

	cli := storagepb.NewStorageClient(c.conn)
	res, err := cli.Occupy(ctx, &storagepb.OccupyModel{Domain: storagepb.Domain(dom), Cell: cell, TransactionId: tx}, grpc.WaitForReady(true))
	if err != nil {
		return
	}
	switch res.GetCase().(type) {
	case *storagepb.OccupyResult_Success:
	case *storagepb.OccupyResult_Forbidden:
		forbidden = true
		reason = "forbidden by business backend"
	case *storagepb.OccupyResult_AlreadyOccupied:
		forbidden = true
		reason = "already occupied"
	default:
		forbidden = true
		reason = "not implemented"
	}
	return
}

func (c *Client) ReleaseStorageCell(ctx context.Context, domain, cell, tx string, strictDomainCheck bool) (forbidden bool, reason string, err error) {
	dom, ok := storagepb.Domain_value[domain]
	if !ok {
		err = fmt.Errorf("unknown domain %q", domain)
		return
	}

	cli := storagepb.NewStorageClient(c.conn)
	res, err := cli.Release(ctx, &storagepb.ReleaseModel{Domain: storagepb.Domain(dom), Cell: cell, TransactionId: tx, StrictDomain: strictDomainCheck}, grpc.WaitForReady(true))
	if err != nil {
		return
	}
	switch res.GetCase().(type) {
	case *storagepb.ReleaseResult_Success:
	case *storagepb.ReleaseResult_Forbidden:
		forbidden = true
		reason = "forbidden by business backend"
	case *storagepb.ReleaseResult_NotFound:
		forbidden = true
		reason = "not occupied"
	case *storagepb.ReleaseResult_WrongDomain:
		forbidden = true
		reason = "cell is occupied under another domain"
	default:
		forbidden = true
		reason = "not implemented"
	}
	return
}

func (c *Client) IntegrationUIMethod(ctx context.Context, method string, kv map[string]interface{}) (result map[string]interface{}, httpStatus int, err error) {
	bodyps, err := structpb.NewStruct(kv)
	if err != nil {
		err = fmt.Errorf("casting kv to proto struct: %w", err)
		return
	}

	cli := integrationpb.NewIntegrationClient(c.conn)
	res, err := cli.UIMethod(ctx, &integrationpb.UIMethodModel{Method: method, Body: bodyps}, grpc.WaitForReady(true))
	if err != nil {
		return
	}

	result = res.GetBody().AsMap()
	httpStatus = int(res.GetHttpStatus())
	return
}
