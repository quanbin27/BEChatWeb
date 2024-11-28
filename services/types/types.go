package types

import (
	"context"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/contacts"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/groups"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/messages"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/users"
	"gorm.io/gorm"
	"time"
)

type UserStore interface {
	GetUserByID(id int32) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
	UpdateInfo(userID int32, updatedData map[string]interface{}) error
	UpdatePassword(userID int32, password string) error
	GetNameByID(id int32) (string, error)
	GetUsersByIDs(userIDs []int32) ([]User, error)
}
type GroupStore interface {
	ChangeNameGroup(id int32, name string) error
	DeleteGroup(id int32) error
	GetGroupByID(id int32) (*Group, error)
	CreateGroup(group *Group) (int32, error)
	AddMemberToGroup(groupDetail *GroupDetail) error
	DeleteGroupDetails(groupID, userID int32) error
	UpdateGroup(group *Group) error
	GetListUserGroups(userID int32) ([]Group, error)
	GetRoleIDByUserAndGroup(userID, groupID int32) (int32, error)
	GetListMembers(groupID int32) ([]Member, error)
	ChangeAdmin(groupID int32, currentAdminID int32, newAdminID int32) error
}
type MessageStore interface {
	SendMessage(msg *Message) (int32, time.Time, error)
	GetMessages(groupID int32) ([]Message, error)
	GetLatestMessages(groupID int32) (Message, error)
	GetMessageByID(messageID int32, msg *Message) error
	DeleteMessage(msg *Message) error
}
type ContactStore interface {
	GetContacts(userID int32) ([]Contact, error)
	AddContact(contact *Contact) error
	RemoveContact(userID, contactUserID int32) error
	AcceptContact(userID, contactUserID int32) error
	RejectContact(userID, contactUserID int32) error
	GetPendingReceivedContacts(userID int32) ([]Contact, error)
	GetPendingSentContacts(userID int32) ([]Contact, error)
	GetFriendIDs(userID int32) ([]int32, error)
	GetContactsNotInGroup(userID int32, groupID int32) ([]Contact, error)
}
type UserService interface {
	CreateUser(ctx context.Context, user *users.RegisterRequest) error
	CreateJWT(ctx context.Context, login *users.LoginRequest) (string, error)
	UpdateUser(ctx context.Context, update *users.ChangeInfoRequest) error
	UpdatePassword(ctx context.Context, update *users.ChangePasswordRequest) error
	GetUserByID(ctx context.Context, id int32) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
type GroupService interface {
	CreateGroup(ctx context.Context, req *groups.CreateGroupRequest) error
	DeleteGroup(ctx context.Context, req *groups.DeleteGroupRequest) error
	ChangeNameGroup(ctx context.Context, req *groups.ChangeNameRequest) (string, error)
	GetGroupInfo(ctx context.Context, id int32) (*Group, error)
	AddMember(ctx context.Context, req *groups.AddMemberRequest) error
	KickMember(ctx context.Context, req *groups.KickMemberRequest) error
	ChangeAdmin(ctx context.Context, req *groups.ChangeAdminRequest) error
	GetListMembers(ctx context.Context, groupID int32) ([]Member, error)
	LeaveGroup(ctx context.Context, req *groups.LeaveGroupRequest) error
	GetListUserGroups(ctx context.Context, userID int32) ([]Group, error)
}
type MessageService interface {
	DeleteMessage(ctx context.Context, req *messages.DeleteMessageRequest) error
	GetLatestMessages(ctx context.Context, req *messages.GetLatestMessagesRequest) (Message, error)
	GetMessages(ctx context.Context, req *messages.GetMessagesRequest) ([]Message, error)
	SendMessage(ctx context.Context, req *messages.SendMessageRequest) (int32, time.Time, error)
}
type ContactService interface {
	AddContact(ctx context.Context, req *contacts.AddContactRequest) error
	RemoveContact(ctx context.Context, req *contacts.RemoveContactRequest) error
	AcceptContact(ctx context.Context, req *contacts.AcceptContactRequest) error
	GetContacts(ctx context.Context, userID int32) ([]*contacts.Contact, error)
	GetPendingSentContacts(ctx context.Context, userID int32) ([]*contacts.Contact, error)
	GetPendingReceivedContacts(ctx context.Context, userID int32) ([]*contacts.Contact, error)
	RejectContact(ctx context.Context, userID, contactUserID int32) error
	GetContactsNotInGroup(ctx context.Context, userID int32, groupID int32) ([]*contacts.Contact, error)
}

// ------ USER ------

type User struct {
	ID        int32          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Bio       string         `json:"bio"`
	Avatar    string         `gorm:"type:longtext"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
type Member struct {
	ID     int32  `gorm:"column:user_id" json:"id"` // Đổi column thành user_id
	Name   string `json:"name"`
	RoleID int32  `json:"role_id"`
}
type RegisterUserPayLoad struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=12"`
}
type LoginUserPayLoad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ChangeInfoPayLoad struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Bio   string `json:"bio" validate:"required"`
}
type ChangePasswordPayLoad struct {
	OldPassword        string `json:"old_password" validate:"required,min=3,max=12"`
	NewPassword        string `json:"new_password" validate:"required,min=3,max=12"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required,min=3,max=12"`
}

// ------ GROUP ------

type Group struct {
	ID          int32      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `json:"name"`
	MemberCount int32      `json:"member_count"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
type CreateGroupPayload struct {
	Name      string  `json:"name"`
	MemberIds []int32 `json:"member_ids"`
}
type ChangeNameGroupPayload struct {
	Name string `json:"name"`
}
type AddMemberToGroupPayload struct {
	MemberIds []int32 `json:"member_ids"`
}
type KickMemberFromGroupPayload struct {
	MemberID int32 `json:"member_id"`
}
type ChangeAdminPayload struct {
	NewAdminID int32 `json:"new_admin_id"`
}

// ------ ROLE ------

type Role struct {
	ID   int32  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"index" json:"name"`
}

// ------GROUP DETAIL ------

type GroupDetail struct {
	UserID    int32      `gorm:"not null;primaryKey" json:"user_id"`
	GroupID   int32      `gorm:"not null;primaryKey" json:"group_id"`
	RoleID    int32      `gorm:"not null" json:"role_id"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	User      User       `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Group     Group      `gorm:"foreignKey:GroupID;references:ID" json:"group"`
	Role      Role       `gorm:"foreignKey:RoleID;references:ID" json:"role"`
}

// ------MESSAGE -------

type Message struct {
	ID             int32      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         int32      `gorm:"not null" json:"user_id"`
	GroupID        int32      `gorm:"not null" json:"group_id"`
	Content        string     `gorm:"type:text" json:"content"`
	ReplyMessageID *int32     `gorm:"index" json:"reply_message_id,omitempty"`
	ReplyMessage   *Message   `gorm:"foreignKey:ReplyMessageID;constraint:OnDelete:CASCADE" json:"reply_message,omitempty"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
type SendMessagePayload struct {
	GroupID        int32  `json:"group_id" validate:"required"`
	Content        string `json:"content" validate:"required"`
	MessageReplyID *int32 `json:"message_reply_id,omitempty"`
}
type DeleteMessagePayload struct {
	MessageID int32 `json:"message_id" validate:"required"`
}

// ------ CONTACT-------

type Contact struct {
	ID            int32     `gorm:"primaryKey;autoIncrement"`
	UserID        int32     `gorm:"not null;uniqueIndex:unique_friendship"`
	ContactUserID int32     `gorm:"not null;uniqueIndex:unique_friendship"`
	Status        string    `gorm:"type:enum('PENDING', 'ACCEPTED','REJECTED');default:'PENDING';not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	User          User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ContactUser   User      `gorm:"foreignKey:ContactUserID;constraint:OnDelete:CASCADE"`
}
type AddContactPayload struct {
	ContactUserID int32 `json:"contact_user_id" validate:"required"`
}

type RemoveContactPayload struct {
	ContactUserID int32 `json:"contact_user_id" validate:"required"`
}

type AcceptContactPayload struct {
	ContactUserID int32 `json:"contact_user_id" validate:"required"`
}
