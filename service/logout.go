package service

type User struct {
	Name     string `json:"name,omitempty"`
	GoToPath string `json:"go_to_path,omitempty"`
	Age      int    `json:"age,omitempty"`
	YourHome string `json:"your_home,omitempty"`
}
