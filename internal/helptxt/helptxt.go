package helptxt

var Rm = `Usage: rm [-FLAGS] PATH

Remove a bin or part at a specified PATh.

Flags:
	-r		Recursive: remove bins and their contents recursively
	-v		Verbose: explain what is being done
	-qN		Quantity: delete a specified quantity (N) of parts sharing a name
	-h		Help: display this help message

A bin can not be removed if it contains bins/parts.
Parts can have the same name. When using a part's name in path with rm, the first returned part will be removed.
To remove a specific part, you will need to use that part's sku.`

var Mv = `Usage: mv [-FLAGS] SRC DEST

Rename a bin/part (SRC) to dest, or Move a bin/part (SRC) to a specified bin (DEST).

Flags:
	-v		Verbose: explain what is being done
	-h		Help: display this help message`

var Ls = `Usage: ls [-FLAGS] PATH

List bin(s)/parts(s). ls without a path will display the root bin.

Flags:
	-l		Long: display additional info about the bin(s)/part(s)
	-h		Help: display this help message

ls can accept a part sku to diplay the -l version of that part.`

var Mkprt = `Usage: mkprt [-FLAGS] PATH

Create a new part or parts in specified bin.

Flags:
	-v		Verbose: explain what is being done
	-qN		Quantity: create a specified quantity (N) of parts with the same name

Parts can only be contained in a bin. Trying to put a part in another part will return an error.
When a part is created, a sku is automatically generated for it. The sku is the abbreviated name followed by the part number`

var Mkbin = `Usage: mkbin [-FLAGS] PATH

Create a new bin.

Flags:
	-v		Verbose: explain what is being done
	-p		Parent: create parent bins as needed

The path will always start at the root bin. Bins can have the same name as the parent bin, 
but bins can not have the same name within the same parent.`

