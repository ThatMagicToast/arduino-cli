// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package commands

import (
	"context"

	"github.com/arduino/arduino-cli/commands/cache"
	"github.com/arduino/arduino-cli/commands/updatecheck"
	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
)

// NewArduinoCoreServer returns an implementation of the ArduinoCoreService gRPC service
// that uses the provided version string.
func NewArduinoCoreServer(version string) rpc.ArduinoCoreServiceServer {
	return &arduinoCoreServerImpl{
		versionString: version,
	}
}

type arduinoCoreServerImpl struct {
	rpc.UnsafeArduinoCoreServiceServer // Force compile error for unimplemented methods

	versionString string
}

// UpdateLibrariesIndex FIXMEDOC
func (s *arduinoCoreServerImpl) UpdateLibrariesIndex(req *rpc.UpdateLibrariesIndexRequest, stream rpc.ArduinoCoreService_UpdateLibrariesIndexServer) error {
	syncSend := NewSynchronizedSend(stream.Send)
	res, err := UpdateLibrariesIndex(stream.Context(), req,
		func(p *rpc.DownloadProgress) {
			syncSend.Send(&rpc.UpdateLibrariesIndexResponse{
				Message: &rpc.UpdateLibrariesIndexResponse_DownloadProgress{DownloadProgress: p},
			})
		},
	)
	if res != nil {
		syncSend.Send(&rpc.UpdateLibrariesIndexResponse{
			Message: &rpc.UpdateLibrariesIndexResponse_Result_{Result: res},
		})
	}
	return err
}

// Version FIXMEDOC
func (s *arduinoCoreServerImpl) Version(ctx context.Context, req *rpc.VersionRequest) (*rpc.VersionResponse, error) {
	return &rpc.VersionResponse{Version: s.versionString}, nil
}

// NewSketch FIXMEDOC
func (s *arduinoCoreServerImpl) NewSketch(ctx context.Context, req *rpc.NewSketchRequest) (*rpc.NewSketchResponse, error) {
	return NewSketch(ctx, req)
}

// LoadSketch FIXMEDOC
func (s *arduinoCoreServerImpl) LoadSketch(ctx context.Context, req *rpc.LoadSketchRequest) (*rpc.LoadSketchResponse, error) {
	resp, err := LoadSketch(ctx, req)
	return &rpc.LoadSketchResponse{Sketch: resp}, err
}

// SetSketchDefaults FIXMEDOC
func (s *arduinoCoreServerImpl) SetSketchDefaults(ctx context.Context, req *rpc.SetSketchDefaultsRequest) (*rpc.SetSketchDefaultsResponse, error) {
	return SetSketchDefaults(ctx, req)
}

// Upload FIXMEDOC
func (s *arduinoCoreServerImpl) Upload(req *rpc.UploadRequest, stream rpc.ArduinoCoreService_UploadServer) error {
	syncSend := NewSynchronizedSend(stream.Send)
	outStream := feedStreamTo(func(data []byte) {
		syncSend.Send(&rpc.UploadResponse{
			Message: &rpc.UploadResponse_OutStream{OutStream: data},
		})
	})
	errStream := feedStreamTo(func(data []byte) {
		syncSend.Send(&rpc.UploadResponse{
			Message: &rpc.UploadResponse_ErrStream{ErrStream: data},
		})
	})
	res, err := Upload(stream.Context(), req, outStream, errStream)
	outStream.Close()
	errStream.Close()
	if res != nil {
		syncSend.Send(&rpc.UploadResponse{
			Message: &rpc.UploadResponse_Result{
				Result: res,
			},
		})
	}
	return err
}

// UploadUsingProgrammer FIXMEDOC
func (s *arduinoCoreServerImpl) UploadUsingProgrammer(req *rpc.UploadUsingProgrammerRequest, stream rpc.ArduinoCoreService_UploadUsingProgrammerServer) error {
	syncSend := NewSynchronizedSend(stream.Send)
	outStream := feedStreamTo(func(data []byte) {
		syncSend.Send(&rpc.UploadUsingProgrammerResponse{
			Message: &rpc.UploadUsingProgrammerResponse_OutStream{
				OutStream: data,
			},
		})
	})
	errStream := feedStreamTo(func(data []byte) {
		syncSend.Send(&rpc.UploadUsingProgrammerResponse{
			Message: &rpc.UploadUsingProgrammerResponse_ErrStream{
				ErrStream: data,
			},
		})
	})
	err := UploadUsingProgrammer(stream.Context(), req, outStream, errStream)
	outStream.Close()
	errStream.Close()
	if err != nil {
		return err
	}
	return nil
}

// SupportedUserFields FIXMEDOC
func (s *arduinoCoreServerImpl) SupportedUserFields(ctx context.Context, req *rpc.SupportedUserFieldsRequest) (*rpc.SupportedUserFieldsResponse, error) {
	return SupportedUserFields(ctx, req)
}

// BurnBootloader FIXMEDOC
func (s *arduinoCoreServerImpl) BurnBootloader(req *rpc.BurnBootloaderRequest, stream rpc.ArduinoCoreService_BurnBootloaderServer) error {
	syncSend := NewSynchronizedSend(stream.Send)
	outStream := feedStreamTo(func(data []byte) {
		syncSend.Send(&rpc.BurnBootloaderResponse{
			Message: &rpc.BurnBootloaderResponse_OutStream{
				OutStream: data,
			},
		})
	})
	errStream := feedStreamTo(func(data []byte) {
		syncSend.Send(&rpc.BurnBootloaderResponse{
			Message: &rpc.BurnBootloaderResponse_ErrStream{
				ErrStream: data,
			},
		})
	})
	resp, err := BurnBootloader(stream.Context(), req, outStream, errStream)
	outStream.Close()
	errStream.Close()
	if err != nil {
		return err
	}
	return syncSend.Send(resp)
}

// ListProgrammersAvailableForUpload FIXMEDOC
func (s *arduinoCoreServerImpl) ListProgrammersAvailableForUpload(ctx context.Context, req *rpc.ListProgrammersAvailableForUploadRequest) (*rpc.ListProgrammersAvailableForUploadResponse, error) {
	return ListProgrammersAvailableForUpload(ctx, req)
}

// ArchiveSketch FIXMEDOC
func (s *arduinoCoreServerImpl) ArchiveSketch(ctx context.Context, req *rpc.ArchiveSketchRequest) (*rpc.ArchiveSketchResponse, error) {
	return ArchiveSketch(ctx, req)
}

// EnumerateMonitorPortSettings FIXMEDOC
func (s *arduinoCoreServerImpl) EnumerateMonitorPortSettings(ctx context.Context, req *rpc.EnumerateMonitorPortSettingsRequest) (*rpc.EnumerateMonitorPortSettingsResponse, error) {
	return EnumerateMonitorPortSettings(ctx, req)
}

// CheckForArduinoCLIUpdates FIXMEDOC
func (s *arduinoCoreServerImpl) CheckForArduinoCLIUpdates(ctx context.Context, req *rpc.CheckForArduinoCLIUpdatesRequest) (*rpc.CheckForArduinoCLIUpdatesResponse, error) {
	return updatecheck.CheckForArduinoCLIUpdates(ctx, req)
}

// CleanDownloadCacheDirectory FIXMEDOC
func (s *arduinoCoreServerImpl) CleanDownloadCacheDirectory(ctx context.Context, req *rpc.CleanDownloadCacheDirectoryRequest) (*rpc.CleanDownloadCacheDirectoryResponse, error) {
	return cache.CleanDownloadCacheDirectory(ctx, req)
}
