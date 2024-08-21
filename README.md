# LMCMD README

This is a small pet project that I find quite useful. It is using the Llama3 70B model to generate a command from the task that you want to perform.It will give you some explanation of the command first, and then output the command to run on terminal. 

Currently it supports linux and mac command line.

The first run will ask you for the API key from OpenAI, and save to the *.lmcmd.config* in the home folder.  Put this command on the PATH to run anywhere.

To use the program, just run:

lmcmd <Command that you want to generate>

For example:

lmcmd.go "create a file named haha.txt with content of 'hello world'"

 It will output 

 ```
Generating command for mac
Output from LLM: {
    "command": "echo 'hello world' > haha.txt",
    "explanation": "Create a file named haha.txt with the content 'hello world'"
}
Command: echo 'hello world' > haha.txt
Explanation: Create a file named haha.txt with the content 'hello world'
Command copied to clipboard.
 ```