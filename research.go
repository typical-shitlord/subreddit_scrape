//This is very much a quick hack, but it works regardless
//
//TODO: check if file is empty
//      Support command line options
//      Properly sanitize input
//      Use a slice for appending data, making sure to check that it isn't duplicate
//      Make the UI show progress and not output until done
//      Use goroutines to check multiple headers at once (note: race conditions don't matter)
//      Figure out lowest possible ms delay required to get a header on start (note that this varies by PC and connection)

package main

import (
    "fmt"
    "os"
    "bufio"
	"net/http"
	"container/list"
	"errors"
	//"time" //see line 91
)

func catchDir(req *http.Request, via []*http.Request) error {
	return errors.New("Forward blocked")
}

func check(s string) int {
    // returns 0 on valid
    // 1 on private
    // -1 on not found/err

	gethttp := &http.Client{
	    CheckRedirect: catchDir,
	}

	subreddit := "http://www.reddit.com/r/" + s //Potential issue: unsanitized input
	httpStatus, _ := gethttp.Head(subreddit)
	code := httpStatus.StatusCode

    if code == 200 {
    	return 0
    } else if code == 403 {
    	return 1
    } else {
    	return -1
    }
}

func printStat(subreddit string) {
    code := check(subreddit)

	if code == 0 {
		fmt.Printf("[X] %s\n", subreddit)
	} else if code == 1 {
		fmt.Printf("[P] %s\n", subreddit)
	}
}

func main() {
	fmt.Println("Subreddit dictionary brute force")
	fmt.Println("[X] means public, [P] means private\n")

	dir, e := os.Open("list.txt") //hardcoded file being opened

	if e != nil {
		fmt.Println("error opening list")
		return
	}

    brute := list.New()
    placeholder := brute.PushFront("")
    scanner := bufio.NewScanner(dir)

    defer dir.Close()

    for scanner.Scan() {
	    brute.InsertAfter(scanner.Text(), placeholder)
    }

    for e := brute.Front(); e != nil; e = e.Next() {
		for i := brute.Front(); i != nil; i = i.Next() {
			if e.Value != i.Value {

				mess_word := fmt.Sprintf("%s%s", e.Value, i.Value)

				for x := brute.Front(); x != nil; x = x.Next() {

					if mess_word != fmt.Sprintf("%s", x.Value) {

						final_word := fmt.Sprintf("%s%s", mess_word, x.Value)
						printStat(final_word) //This should check if in slice/append if not

						//go printStat(final_word) //By setting the 3 below to the lowest possible number of ms to wait to get a 
						//time.Sleep(time.Millisecond * 10) //header we can drastically speed up lookup times with goroutines

						// Current way:        0m21.511s
						// 10 ms + goroutines: 0m10.509s
						// 3 ms + goroutines:  0m3.292s
					}
				}

			}

		}
	}

}
