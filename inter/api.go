package inter

import "fmt"
type Handler interface {
	Packager
}
type Kjcx interface {
	Add()
}
type client struct {
	Packager
}
func (mb *client)Add(){
	mb.Packager.Encode()
	fmt.Println("client add")
}
func NewClient(handler Handler) Kjcx{
	return &client{Packager:handler}
}
type Test interface {
	Packager
}
 type Packager interface {
 	Encode()
 	Decode()
 }
type ServerHandler struct {
	ServerPackager
}
type ServerPackager struct {
	SaveId string
}
func (s *ServerPackager)Encode(){
	fmt.Println("server encode")
}
func (s *ServerPackager)Decode(){
	fmt.Println("server decode")
}
type ClientHandler struct {
	ClientPackager
}

type ClientPackager struct {
	SaveId string
}
func (c *ClientPackager)Encode(){
	fmt.Println("client encode")
}
func (c *ClientPackager)Decode(){
	fmt.Println("client decode")
}