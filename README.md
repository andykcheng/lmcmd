# LMCMD README

This is a small pet project that I find quite useful. It is using the Llama3 70B model to generate a command from the task that you want to perform.It will give you some explanation of the command first, and then output the command to run on terminal. 

Currently it supports linux and mac command line.

The first run will ask you for the API key from Together.ai, and save to the *.lmcmd.config* in the home folder.  Put this command on the PATH to run anywhere.

To use the program, just run:

lmcmd <Command that you want to generate>

For example:

lmcmd Find all files starting with abc, search from root

 It will output 

 ```
 Find all files starting with 'abc' from the root directory (/) and below. The '*' is a wildcard character that matches any characters following 'abc'.  

 find / -name 'abc*'  

 Command copied to clipboard  
 ```