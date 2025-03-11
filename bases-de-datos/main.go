package main

import "bd/storage"

func main()  {
	storage.NewMysqlDB()

	storage.NewMysqlProduct(storage.Pool())
}