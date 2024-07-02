package main

import (
	"fmt"
	"net/http"

	"backend/pkg"
	db "backend/pkg/db/sqlite"
)

func main() {
	// Init Database
	db.DatabaseInit()
	defer db.DB.Close()

	// Serve static files from the current directory
	http.Handle("/", http.FileServer(http.Dir("../public")))
	http.Handle("/avatars/", http.StripPrefix("/avatars", http.FileServer(http.Dir("./avatars"))))
	http.Handle("/postsImages/", http.StripPrefix("/postsImages", http.FileServer(http.Dir("./postsImages"))))
	http.Handle("/commentsImages/", http.StripPrefix("/commentsImages", http.FileServer(http.Dir("./commentsImages"))))

	//auth
	http.HandleFunc("/login", pkg.LoginHandler)
	http.HandleFunc("/register", pkg.RegisterHandler)
	http.HandleFunc("/session", pkg.SessionHandler)

	//profile
	http.HandleFunc("/privacy", pkg.PrivacyHandler)
	http.HandleFunc("/user", pkg.GetUserInfoHandler)
	http.HandleFunc("/editProfile", pkg.ProfileEditorHandler)
	http.HandleFunc("/uploadImage", pkg.UpdateImageHandler)

	//posts
	http.HandleFunc("/posts", pkg.PostsHandler)
	http.HandleFunc("/newPost", pkg.NewPostHandler)
	http.HandleFunc("/newComment", pkg.NewCommentHandler)
	http.HandleFunc("/comment", pkg.CommentHandler)

	//search
	http.HandleFunc("/allusers", pkg.GetAllUsersHandler)

	//followers
	http.HandleFunc("/api/followers", pkg.FollowHandler)
	http.HandleFunc("/messages", pkg.GetMessages)
	
	//groups
	http.HandleFunc("/groups", pkg.GetGroupsHandler)
	http.HandleFunc("/newGroup", pkg.NewGroupHandler)
	

	// websocket
	http.HandleFunc("/ws", pkg.WsHandler)

	// Start the server on port 8080
	fmt.Println("SN is running on http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
