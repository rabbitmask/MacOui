package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {

	targets := flag.String("f", "", "f targets.txt")
	target := flag.String("t", "", "t target")
	flag.Usage = func() {
		fmt.Println("MacOui -f targets.txt")
		fmt.Println("MacOui -t 00-00-00... or 00:00:00...")
	}

	flag.Parse()

	if *targets != "" {
		get_targets(*targets)
	}
	if *target != "" {
		fmt.Println(find_oui(*target))
	}
}



func mac_re(str string,ch chan string) {
	ip,_:=regexp.Compile("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}")
	mac,_:=regexp.Compile("([A-Fa-f0-9]{2}[:-]){5}[A-Fa-f0-9]{2}")

	if len(mac.FindString(str))>0{
		ch <- ip.FindString(str) + "\t" + mac.FindString(str) + "\t" + find_oui(mac.FindString(str))
	} else{
		ch <- ""
	}

}


func find_oui(mac string) string {
	macindex,_:=regexp.Compile("([A-Fa-f0-9]{2}[:-]){2}[A-Fa-f0-9]{2}")
	mac=macindex.FindString(mac)
	mac=strings.ToUpper(mac)
	mac=strings.Replace(mac,"-",":",-1)

	file, err := os.Open("manuf")
	if err != nil {
		fmt.Println("You need update your database:")
		fmt.Println("curl -O https://gitlab.com/wireshark/wireshark/-/raw/master/manuf")
		os.Exit(1)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var res string
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if strings.Contains(str, mac)==true{

			strSplit := strings.Split(str, "\t")
			if len(strSplit)==3{
				res=("[ "+strSplit[1]+" ]\t"+strings.Replace(strSplit[2],"\n","",-1))
			}else {
				res=("[ "+strings.Replace(strSplit[1],"\n","",-1)+" ]")
			}





			break
		}

	}
	return res

}

func get_targets(targets string){
	file, err := os.Open(targets)
	checkErr(err)
	defer file.Close()
	reader := bufio.NewReader(file)

	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		ch := make(chan string, 1)
		go mac_re(str,ch)
		res:=<-ch
		if len(res)>0{
			fmt.Println(res)
		}
		}

}




func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}