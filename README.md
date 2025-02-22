# Doro the pomodoro timer

A minimalistic command-line Pomodoro timer. It helps you
manage your work sessions efficiently by dividing your time into focused work
intervals and short breaks. The timer follows the traditional Pomodoro
technique, allowing you to customize the duration of your work and break
sessions.

## Installation

-   Clone the repository:

`git clone github.com/mikeyQwn/doro`

-   Build and install the application:

`go install ./cmd/doro`

-   Verify the installation by running:

`doro -v`

## Usage

To run Doro, use the following command:

`doro [flags]`

### Flags

-   -w: Duration of the work session in minutes (default: 25)
-   -s: Duration of the short break session in minutes (default: 5)
-   -l: Duration of the long break session in minutes (default: 30)
-   -h: Display help information

### Examples

1. Start a default 25-minute work session followed by a 5-minute break:

`doro -d`

2. Start a custom 30-minute work session followed by a 10-minute break:

`doro -w 30 -b 10`

### Notes

-   The timer will automatically alternate between work and break sessions until
    you manually stop it (e.g., by pressing `Ctrl+C` )

-   You can pause the timer at any time by pressinng `Space`

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file
for details.
