package main

func main() {
	//cluster, err := gocb.Connect("couchbase://localhost")
	//if err != nil {
	//	fmt.Printf("error connecting to the cluster: %s", err.Error())
	//}
	//cluster.Authenticate(gocb.PasswordAuthenticator{
	//	Username: "Administrator",
	//	Password: "password",
	//})
	//bucket, err := cluster.OpenBucket("increment", "")
	//if err != nil {
	//	panic(err)
	//}
	//
	//inc, err := incrmntr.New(bucket, 999999999999999, 1, 1)
	//if err != nil {
	//	fmt.Printf("%s", err.Error())
	//}
	//for i := 0; i < 20000; i++ {
	//	err := inc.AddSafe("test")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}
}
