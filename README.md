# go-dependency-tree
This program scans the current directory for a Go project and prints a graphviz-compatible format of the project's dependencies to the shell.


## Usage
> go build .
> 
> ./go-dependency-tree > tree.gv.txt

Now use dot or a graphviz program of your choosing to compile the tree into an image
For example:
> dot -Tsvg tree.gv.txt > dependency.svg

For this project, this outputs the following image:
![dep-tree](https://user-images.githubusercontent.com/8042849/170895153-d1f69beb-d928-4019-ad71-d043269f5cfd.svg)

## Warning
For extremely large projects, e.g. kubernetes, very common dependencies such as fmt, bufio, log, os, etc., will have a very large number of lines connected to them.

This may lead to an unreadable image using the command above, some arguments may need to be set for the graphviz command you are using.

For example:
> sfdp -x -Goverlap=false -Gsplines=true -Tsvg example.gv.txt > image.svg
