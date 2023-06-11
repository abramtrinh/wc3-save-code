#wc3-save-code

wc3-save-code is a code generator for WC3 custom map saves. It is also contains functions to find the preimage of lookup2 (a non-cryptographic hash function).

## Table of Content:

* [Overview](#overview)
* [Usage](#usage)
* [Map Protection](#map-protection)
* [Notes](#notes)
* [To-do](#to-do)

## Overview:

**This repo contains 4 parts.**

* **breakingHash** - Contains functions that find the preimage of a hash (hashed by lookup2).
&emsp;&emsp;**BruteForceLookup2**
&emsp;&emsp;&emsp;&emsp;2^n time complexity preimage attack that returns a alphanumeric string of length 10. (Can be concurrent with goroutines option)
&emsp;&emsp;**UnHash**
&emsp;&emsp;&emsp;&emsp;Constant time complexity preimage attack that returns a byte slice of length 12.


* **lookup2** - Go port of the 32-bit C hash, [lookup2.c](https://burtleburtle.net/bob/hash/evahash.html).


* **saveCode** - Contains a save code generator.


* **saveFunctions** - Contains functions that use player names to generate values.
&emsp;&emsp;**StringHash**
&emsp;&emsp;&emsp;&emsp;Reverse engineered implementation of WC3's StringHash from the common.j file which is based off of lookup2.c.

## Usage:
```Go
package main

import (
    "fmt"

	"github.com/abramtrinh/wc3-save-code/breakingHash"
	"github.com/abramtrinh/wc3-save-code/lookup2"
    "github.com/abramtrinh/wc3-save-code/saveCode"

)

func main() {

	// Hashed value is: 1067317194
	hash := lookup2.Hash("Hello World", 0)

	// String: �O[♠�▬�%☼dwQ
	// String in byte slice form: [201 79 91 6 255 22 182 37 15 100 119 81]
	unhashString := string(breakingHash.UnHash(hash))

	// Hashed value is: 1067317194
	lookup2.Hash(unhashString, 0)

	// String: YBXgebLrvs
	// Time elapsed 4m34.2813502s
	// 5 is the number of goroutines spun up.
	bruteString, _ := breakingHash.BruteForceLookup2(hash, 5)

	// Hashed value is: 1067317194
	lookup2.Hash(bruteString, 0)

    // Generates a save code for player name "Hello" with rank of 10.
    save := saveCode.NewSave("Hello", 6)
	fmt.Println(save.SaveCode(10))

}
```

## Map Protection:
* **Obfuscation and minimization** - Makes the decompiled script file tedious to read and edit. Not impossible to edit but makes it take longer.
* **Using functions like SetRandomSeed and GetRandomInt** - Makes it so the only feasible option is to directly edit the map since you can't exactly recreate the function.
* **Periodically poll for state of certain objects** - A bit annoying to deal with due it being spread throughout the script file. Easier to notice if map is not obfuscated since you can look for the function names.
&emsp;&emsp;**e.g.** The map turns off defeat conditions and then polls if said conditions are off. If they are, the map kicks you. The map gives you max damage and then polls if you (unknowingly) kill a certain unit. If you do, the map kicks you.
* **Checksum or map signature** - Straight up can't edit the map. If you edit anything, the map fails to load.
* **Newline** - In certain editors, if you edit a map that contains a \n character, the map fails to load since the \n character is evaluated literally. It is not escaped.

## Notes:
* The simplest way to generate a save code is to directly edit the decompiled WC3 maps and script file. As a result, you would have to distribute the map. You would also need to edit the file every time you need to make any changes to the saved values. Also note that certain maps will fail to load if edited due to a change in the file signature.

* If I still really wanted to make a save-code generator, focus on maps that don't use SetRandomSeed and GetRandomInt functions. The implementation of those two functions are not known and I have not been able to figure it out.

* Use BruteForceLookup2 if you want to return an alphanumeric string and don't mind waiting. (longest I've waited is 35 minutes)
* Use UnHash if you want speed but don't mind the fact it contains UTF-8 control characters.


* There are two implementation of random strings for the BruteForceLookup2. Benchmarks are provided for fun. Use the mask version for speed.
	* About 115 ns/op vs. 49 ns/op difference.


* The main save function in saveFunctionsJASS.go "works" if you set the values of GetRandomInt and SetRandomSeed to what the WC3 map uses.


## To-do:
* Rewrite saveFunctionsJASS.go as idomatic Go instead of mimicing the decompiled structure of it.
* Implement a second UnHash() that returns alphanumeric strings instead of weird UTF-8.