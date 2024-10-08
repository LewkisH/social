package app

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// only some info for searching purposes
type SearchData struct {
	ID           int
	FirstName    string
	LastName     string
	Avatar       string
	IsFollowing  int
	AreFollowing int
	Privacy      int
}

// for registering a user
type RegisterData struct {
	Email           string `json:"email"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	DateOfBirth     string `json:"dateOfBirth"`
	Avatar          string `json:"avatar"`
	AvatarName      string `json:"avatarName"`
	NickName        string `json:"nickName"`
	AboutMe         string `json:"aboutMe"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

// for auth checking
type User struct {
	ID           int            `json:"id"`
	Email        string         `json:"email"`
	FirstName    string         `json:"firstName"`
	Privacy      int            `json:"privacy"`
	LastName     string         `json:"lastName"`
	DateOfBirth  sql.NullString `json:"dateOfBirth"`
	Avatar       string         `json:"avatar"`
	NickName     sql.NullString `json:"nickName"`
	AboutMe      sql.NullString `json:"aboutMe"`
	Session      sql.NullString `json:"session"`
	IsFollowing  int            `json:"isFollowing"`
	AreFollowing int            `json:"areFollowing"`
	Followers    []SearchData   `json:"followers"`
	Following    []SearchData   `json:"following"`
	Posts        []Post         `json:"posts"`
}

type Session struct {
	Token   string `json:"token"`
	Expires int    `json:"expires"`
}

type WSMessage struct {
	Type       string `json:"messageType"`
	FromUserID int    `json:"fromUserID"`
	ToUserID   int    `json:"toUserID"`
	GroupID    int    `json:"groupID"`
}

type ChatIDResponse struct {
	Type   string `json:"type"`
	ChatID int    `json:"chatID"`
}

type PrivateMessage struct {
	AvatarPath   string `json:"AvatarPath"`
	Type         string `json:"type"`
	MsgID        int    `json:"msgID"`
	ChatID       int    `json:"chatID"`
	FromUserID   int    `json:"fromUserID"`
	FromUsername string `json:"fromUsername"`
	ToUserID     int    `json:"toUserID"`
	Content      string `json:"content"`
	Time         string `json:"time"`
}

type GroupMessage struct {
	AvatarPath string `json:"AvatarPath"`
	Type       string `json:"type"`
	MsgID      int    `json:"msgID"`
	ChatID     int    `json:"chatID"`
	FromUserID int    `json:"fromUserID"`
	Content    string `json:"content"`
	Time       string `json:"time"`
}

type GroupChatInfo struct {
	ChatID   int
	ChatSeen int
}

type ChatMessage struct {
	ID       int    `json:"messageID"`
	Content  string `json:"content"`
	User     string `json:"user"`
	Date     string `json:"date"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func (msg *ChatMessage) SetUsername(db *sql.DB) error {
	userID, err := strconv.Atoi(msg.User)
	if err != nil {
		fmt.Println("strconv error in method SetUsername")
		return err
	}
	var firstname, lastname, avatar string
	err = db.QueryRow("SELECT firstname, lastname, avatar FROM users where id = ?", userID).Scan(&firstname, &lastname, &avatar)
	if err != nil {
		return err
	}
	msg.Username = firstname + " " + lastname
	msg.Avatar = avatar
	return nil
}

type MessageGetter struct {
	ID     int       `json:"message_id"` // last existing message id if 0 then no messages exist
	ChatID int       `json:"chat_id"`    // chat id
	Date   time.Time `json:"date"`
}

type PostPreview struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userID"`
	Content   string    `json:"content"`
	Img       string    `json:"img"`
	CreatedAt string    `json:"createdAt"`
	Comments  []Comment `json:"comments"`
}

type Post struct {
	ID               int       `json:"id"`
	UserID           int       `json:"userID"`
	Content          string    `json:"content"`
	Img              string    `json:"img"`
	CreatedAt        string    `json:"createdAt"`
	Privacy          int       `json:"privacy"`
	GroupID          int       `json:"groupID"`
	CustomPrivacyIDs []int     `json:"customPrivacyIDs"`
	Comments         []Comment `json:"comments"`
}

type Comment struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userID"`
	PostID    int    `json:"postID"`
	Content   string `json:"content"`
	Img       string `json:"img"`
	CreatedAt string `json:"createdAt"`
	User      User   `json:"user"`
}

type Group struct {
	ID          int           `json:"id"`
	Name        string        `json:"title"`
	Description string        `json:"description"`
	OwnerID     int           `json:"ownerID"`
	CreatedAt   string        `json:"createdAt"`
	JoinStatus  int           `json:"joinStatus"`
	OwnerName   string        `json:"ownerName"`
	Posts       []PostPreview `json:"posts"`
	ChatID      int           `json:"chatid"`
	Media       sql.NullString        `json:"media"`
}

type Event struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Date          string `json:"date"`
	GroupID       int    `json:"groupID"`
	OwnerID       int    `json:"ownerID"`
	OwnerName     string `json:"ownerName"`
	Certainty     int    `json:"certainty"`
	Going         int    `json:"going"`
	GoingCount    int    `json:"goingCount"`
	NotSureCount  int    `json:"notSureCount"`
	NotGoingCount int    `json:"notGoingCount"`
}

type Notification struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Link      string `json:"link"`
	Seen      bool   `json:"seen"`
	CreatedAt string `json:"createdAt"`
	Type      string `json:"type"`
	FromID    int    `json:"fromID"`
}
