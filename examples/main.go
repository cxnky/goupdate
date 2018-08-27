package main

import (
	"github.com/cxnky/goupdate"
	"fmt"
)

/**

 * Created by cxnky on 23/08/2018 at 15:40
 * goupdate_examples
 * https://github.com/cxnky/
 
**/

func main() {

	fmt.Println("NEW VERSION")

	updater := goupdate.CreateUpdater("https://connorwright.uk/goupdate/updater.json", "1.0.0.0", 5000, true)

	available, err := updater.CheckForUpdate()

	if err != nil {

		fmt.Println(err)

	}


	if available {

		fmt.Println("An update is available!")

		err := updater.PerformUpdate()

		if err != nil {

			fmt.Println(err)
			return

		}

	} else {

		fmt.Println("No updates are available!")

	}

}