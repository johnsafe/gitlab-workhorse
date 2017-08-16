// Code generated by protoc-gen-go. DO NOT EDIT.
// source: shared.proto

package gitaly

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Repository struct {
	StorageName  string `protobuf:"bytes,2,opt,name=storage_name,json=storageName" json:"storage_name,omitempty"`
	RelativePath string `protobuf:"bytes,3,opt,name=relative_path,json=relativePath" json:"relative_path,omitempty"`
	// Sets the GIT_OBJECT_DIRECTORY envvar on git commands to the value of this field.
	// It influences the object storage directory the SHA1 directories are created underneath.
	GitObjectDirectory string `protobuf:"bytes,4,opt,name=git_object_directory,json=gitObjectDirectory" json:"git_object_directory,omitempty"`
	// Sets the GIT_ALTERNATE_OBJECT_DIRECTORIES envvar on git commands to the values of this field.
	// It influences the list of Git object directories which can be used to search for Git objects.
	GitAlternateObjectDirectories []string `protobuf:"bytes,5,rep,name=git_alternate_object_directories,json=gitAlternateObjectDirectories" json:"git_alternate_object_directories,omitempty"`
}

func (m *Repository) Reset()                    { *m = Repository{} }
func (m *Repository) String() string            { return proto.CompactTextString(m) }
func (*Repository) ProtoMessage()               {}
func (*Repository) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{0} }

func (m *Repository) GetStorageName() string {
	if m != nil {
		return m.StorageName
	}
	return ""
}

func (m *Repository) GetRelativePath() string {
	if m != nil {
		return m.RelativePath
	}
	return ""
}

func (m *Repository) GetGitObjectDirectory() string {
	if m != nil {
		return m.GitObjectDirectory
	}
	return ""
}

func (m *Repository) GetGitAlternateObjectDirectories() []string {
	if m != nil {
		return m.GitAlternateObjectDirectories
	}
	return nil
}

type GitCommit struct {
	Id        string        `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Subject   []byte        `protobuf:"bytes,2,opt,name=subject,proto3" json:"subject,omitempty"`
	Body      []byte        `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	Author    *CommitAuthor `protobuf:"bytes,4,opt,name=author" json:"author,omitempty"`
	Committer *CommitAuthor `protobuf:"bytes,5,opt,name=committer" json:"committer,omitempty"`
	ParentIds []string      `protobuf:"bytes,6,rep,name=parent_ids,json=parentIds" json:"parent_ids,omitempty"`
}

func (m *GitCommit) Reset()                    { *m = GitCommit{} }
func (m *GitCommit) String() string            { return proto.CompactTextString(m) }
func (*GitCommit) ProtoMessage()               {}
func (*GitCommit) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{1} }

func (m *GitCommit) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *GitCommit) GetSubject() []byte {
	if m != nil {
		return m.Subject
	}
	return nil
}

func (m *GitCommit) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *GitCommit) GetAuthor() *CommitAuthor {
	if m != nil {
		return m.Author
	}
	return nil
}

func (m *GitCommit) GetCommitter() *CommitAuthor {
	if m != nil {
		return m.Committer
	}
	return nil
}

func (m *GitCommit) GetParentIds() []string {
	if m != nil {
		return m.ParentIds
	}
	return nil
}

type CommitAuthor struct {
	Name  []byte                     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Email []byte                     `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Date  *google_protobuf.Timestamp `protobuf:"bytes,3,opt,name=date" json:"date,omitempty"`
}

func (m *CommitAuthor) Reset()                    { *m = CommitAuthor{} }
func (m *CommitAuthor) String() string            { return proto.CompactTextString(m) }
func (*CommitAuthor) ProtoMessage()               {}
func (*CommitAuthor) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{2} }

func (m *CommitAuthor) GetName() []byte {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *CommitAuthor) GetEmail() []byte {
	if m != nil {
		return m.Email
	}
	return nil
}

func (m *CommitAuthor) GetDate() *google_protobuf.Timestamp {
	if m != nil {
		return m.Date
	}
	return nil
}

type ExitStatus struct {
	Value int32 `protobuf:"varint,1,opt,name=value" json:"value,omitempty"`
}

func (m *ExitStatus) Reset()                    { *m = ExitStatus{} }
func (m *ExitStatus) String() string            { return proto.CompactTextString(m) }
func (*ExitStatus) ProtoMessage()               {}
func (*ExitStatus) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{3} }

func (m *ExitStatus) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

func init() {
	proto.RegisterType((*Repository)(nil), "gitaly.Repository")
	proto.RegisterType((*GitCommit)(nil), "gitaly.GitCommit")
	proto.RegisterType((*CommitAuthor)(nil), "gitaly.CommitAuthor")
	proto.RegisterType((*ExitStatus)(nil), "gitaly.ExitStatus")
}

func init() { proto.RegisterFile("shared.proto", fileDescriptor7) }

var fileDescriptor7 = []byte{
	// 393 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0x4d, 0x8f, 0xd3, 0x40,
	0x0c, 0x55, 0xba, 0x69, 0x21, 0x6e, 0x40, 0x68, 0xd4, 0x43, 0xb4, 0xd2, 0x8a, 0x12, 0x2e, 0x7b,
	0x40, 0x59, 0x54, 0x7e, 0xc1, 0x0a, 0xd0, 0x0a, 0x0e, 0x80, 0x02, 0xf7, 0xc8, 0x6d, 0x4c, 0x62,
	0x94, 0x74, 0xa2, 0x19, 0x67, 0x45, 0xff, 0x22, 0x7f, 0x81, 0x3f, 0x83, 0xe2, 0x69, 0xc4, 0xc7,
	0x81, 0x9b, 0xfd, 0xfc, 0x9e, 0xf3, 0x5e, 0xc6, 0x90, 0xfa, 0x16, 0x1d, 0xd5, 0xc5, 0xe0, 0xac,
	0x58, 0xb3, 0x6a, 0x58, 0xb0, 0x3b, 0x5d, 0x3e, 0x6d, 0xac, 0x6d, 0x3a, 0xba, 0x51, 0x74, 0x3f,
	0x7e, 0xbd, 0x11, 0xee, 0xc9, 0x0b, 0xf6, 0x43, 0x20, 0xe6, 0x3f, 0x23, 0x80, 0x92, 0x06, 0xeb,
	0x59, 0xac, 0x3b, 0x99, 0x67, 0x90, 0x7a, 0xb1, 0x0e, 0x1b, 0xaa, 0x8e, 0xd8, 0x53, 0xb6, 0xd8,
	0x46, 0xd7, 0x49, 0xb9, 0x3e, 0x63, 0x1f, 0xb0, 0x27, 0xf3, 0x1c, 0x1e, 0x39, 0xea, 0x50, 0xf8,
	0x9e, 0xaa, 0x01, 0xa5, 0xcd, 0x2e, 0x94, 0x93, 0xce, 0xe0, 0x27, 0x94, 0xd6, 0xbc, 0x84, 0x4d,
	0xc3, 0x52, 0xd9, 0xfd, 0x37, 0x3a, 0x48, 0x55, 0xb3, 0xa3, 0xc3, 0xb4, 0x3f, 0x8b, 0x95, 0x6b,
	0x1a, 0x96, 0x8f, 0x3a, 0x7a, 0x33, 0x4f, 0xcc, 0x1d, 0x6c, 0x27, 0x05, 0x76, 0x42, 0xee, 0x88,
	0x42, 0xff, 0x6a, 0x99, 0x7c, 0xb6, 0xdc, 0x5e, 0x5c, 0x27, 0xe5, 0x55, 0xc3, 0x72, 0x3b, 0xd3,
	0xfe, 0x5e, 0xc3, 0xe4, 0xdf, 0xc7, 0x0f, 0xa3, 0x27, 0x8b, 0x32, 0x9e, 0xac, 0xe5, 0x3f, 0x22,
	0x48, 0xee, 0x58, 0x5e, 0xdb, 0xbe, 0x67, 0x31, 0x8f, 0x61, 0xc1, 0x75, 0x16, 0xa9, 0x85, 0x05,
	0xd7, 0x26, 0x83, 0x07, 0x7e, 0x54, 0xbd, 0xe6, 0x4c, 0xcb, 0xb9, 0x35, 0x06, 0xe2, 0xbd, 0xad,
	0x4f, 0x1a, 0x2d, 0x2d, 0xb5, 0x36, 0x2f, 0x60, 0x85, 0xa3, 0xb4, 0xd6, 0x69, 0x88, 0xf5, 0x6e,
	0x53, 0x84, 0x7f, 0x5c, 0x84, 0xed, 0xb7, 0x3a, 0x2b, 0xcf, 0x1c, 0xb3, 0x83, 0xe4, 0xa0, 0xb8,
	0x90, 0xcb, 0x96, 0xff, 0x11, 0xfc, 0xa6, 0x99, 0x2b, 0x80, 0x01, 0x1d, 0x1d, 0xa5, 0xe2, 0xda,
	0x67, 0x2b, 0x0d, 0x9b, 0x04, 0xe4, 0x5d, 0xed, 0xf3, 0x16, 0xd2, 0x3f, 0x95, 0x93, 0x49, 0x7d,
	0xa3, 0x28, 0x98, 0x9c, 0x6a, 0xb3, 0x81, 0x25, 0xf5, 0xc8, 0xdd, 0x39, 0x50, 0x68, 0x4c, 0x01,
	0x71, 0x8d, 0x42, 0x1a, 0x67, 0xbd, 0xbb, 0x2c, 0xc2, 0x51, 0x14, 0xf3, 0x51, 0x14, 0x5f, 0xe6,
	0xa3, 0x28, 0x95, 0x97, 0xe7, 0x00, 0x6f, 0xbf, 0xb3, 0x7c, 0x16, 0x94, 0xd1, 0x4f, 0x3b, 0xef,
	0xb1, 0x1b, 0xc3, 0x87, 0x96, 0x65, 0x68, 0xf6, 0x2b, 0x55, 0xbf, 0xfa, 0x15, 0x00, 0x00, 0xff,
	0xff, 0xe1, 0xc9, 0xfa, 0xcc, 0x78, 0x02, 0x00, 0x00,
}
