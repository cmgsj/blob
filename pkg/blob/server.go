package blob

import (
	"context"
	"io"
	"io/fs"
	"sync"

	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ServiceName = blobv1.BlobService_ServiceDesc.ServiceName
)

type Server struct {
	blobv1.UnimplementedBlobServiceServer
	mu sync.RWMutex
	fs afero.Fs
}

func NewServer() blobv1.BlobServiceServer {
	return &Server{
		fs: afero.NewMemMapFs(),
	}
}

func (s *Server) ListFiles(ctx context.Context, req *blobv1.ListFilesRequest) (*blobv1.ListFilesResponse, error) {
	info, err := s.fs.Stat(req.GetPath())
	if err != nil {
		return nil, err
	}
	var infos []fs.FileInfo
	if info.IsDir() {
		infos, err = afero.ReadDir(s.fs, info.Name())
		if err != nil {
			return nil, err
		}
	} else {
		infos = []fs.FileInfo{info}
	}
	var files []*blobv1.File
	for _, info := range infos {
		file, err := s.blobv1File(info)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return &blobv1.ListFilesResponse{Files: files}, nil
}

func (s *Server) GetFile(ctx context.Context, req *blobv1.GetFileRequest) (*blobv1.GetFileResponse, error) {
	info, err := s.fs.Stat(req.GetFileName())
	if err != nil {
		return nil, err
	}
	file, err := s.blobv1File(info)
	if err != nil {
		return nil, err
	}
	return &blobv1.GetFileResponse{File: file}, nil
}

func (s *Server) WriteFile(ctx context.Context, req *blobv1.WriteFileRequest) (*blobv1.WriteFileResponse, error) {
	if req.GetIsDir() {
		err := s.fs.Mkdir(req.GetFileName(), fs.ModePerm)
		if err != nil {
			return nil, err
		}
	} else {
		file, err := s.fs.Create(req.GetFileName())
		if err != nil {
			return nil, err
		}
		_, err = file.Write(req.GetContent())
		if err != nil {
			return nil, err
		}
		err = file.Close()
		if err != nil {
			return nil, err
		}
	}
	return &blobv1.WriteFileResponse{}, nil
}

func (s *Server) RenameFile(ctx context.Context, req *blobv1.RenameFileRequest) (*blobv1.RenameFileResponse, error) {
	err := s.fs.Rename(req.GetFileName(), req.GetNewFileName())
	if err != nil {
		return nil, err
	}
	return &blobv1.RenameFileResponse{}, nil
}

func (s *Server) DeleteFile(ctx context.Context, req *blobv1.DeleteFileRequest) (*blobv1.DeleteFileResponse, error) {
	err := s.fs.Remove(req.GetFileName())
	if err != nil {
		return nil, err
	}
	return &blobv1.DeleteFileResponse{}, nil
}

func (s *Server) blobv1File(info fs.FileInfo) (*blobv1.File, error) {
	file, err := s.fs.Open(info.Name())
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	return &blobv1.File{
		Name:    info.Name(),
		Size:    info.Size(),
		ModTime: timestamppb.New(info.ModTime()),
		IsDir:   info.IsDir(),
		Content: content,
	}, nil
}
