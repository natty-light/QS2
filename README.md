This project is the result of following along to Thorsten Bell's [Writing An Interpreter In Go](https://interpreterbook.com).


There are some deviations, including:
- Mutable variables and variable assignment
- Line numbers in error messages
- For loops
- Support for && and || in logical expressions

In the macro system, line numbers are not propagated to the newly created tokens
